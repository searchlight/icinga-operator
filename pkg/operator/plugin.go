package operator

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/appscode/go/ioutil"
	"github.com/appscode/go/log"
	"github.com/appscode/kutil/meta"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/plugin"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/remotecommand"
)

func (op *Operator) initPluginWatcher() {
	op.pluginInformer = op.monInformerFactory.Monitoring().V1alpha1().SearchlightPlugins().Informer()
	op.pluginQueue = queue.New("SearchlightPlugin", op.MaxNumRequeues, op.NumThreads, op.reconcilePlugin)
	op.pluginInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			queue.Enqueue(op.pluginQueue.GetQueue(), obj)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*api.SearchlightPlugin)
			nu := newObj.(*api.SearchlightPlugin)

			if reflect.DeepEqual(old.Spec, nu.Spec) {
				return
			}
			queue.Enqueue(op.pluginQueue.GetQueue(), nu)
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.pluginQueue.GetQueue(), obj)
		},
	})
	op.pluginLister = op.monInformerFactory.Monitoring().V1alpha1().SearchlightPlugins().Lister()
}

func (op *Operator) reconcilePlugin(key string) error {
	obj, exists, err := op.pluginInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		log.Warningf("SearchlightPlugin %s does not exist anymore\n", key)

		_, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}

		fmt.Println("deleting CheckCommand for ", name)
		return op.ensureCheckCommandDeleted(name)
	}

	plugin := obj.(*api.SearchlightPlugin).DeepCopy()
	log.Infof("Sync/Add/Update for SearchlightPlugin %s\n", plugin.GetName())

	return op.ensureCheckCommand(plugin)
}

func (op *Operator) ensureCheckCommand(wp *api.SearchlightPlugin) error {

	ic := api.IcingaCommand{
		Name: wp.Name,
		Vars: make(map[string]data.CommandVar),
	}

	for _, item := range wp.Spec.Arguments.Vars {
		ic.Vars[item] = data.CommandVar{}
	}

	ic.States = wp.Spec.State

	for _, t := range wp.Spec.AlertKinds {
		if t == api.ResourceKindClusterAlert {
			api.ClusterCommands[api.CheckCluster(wp.Name)] = ic
		} else if t == api.ResourceKindNodeAlert {
			api.NodeCommands[api.CheckNode(wp.Name)] = ic
		} else if t == api.ResourceKindPodAlert {
			api.PodCommands[api.CheckPod(wp.Name)] = ic
		}
	}

	return op.addPluginSupport(wp)
}

func (op *Operator) ensureCheckCommandDeleted(name string) error {

	pod, err := op.GetIcingaPod()
	if err != nil {
		return errors.WithMessage(err, "failed to get Icinga2 Pod name.")
	}

	// Pause all ClusterAlerts for this plugin
	caList, err := op.extClient.MonitoringV1alpha1().ClusterAlerts(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, ca := range caList.Items {
		if ca.Spec.Check == api.CheckCluster(name) {
			ca.Spec.Paused = true
			if _, err := op.extClient.MonitoringV1alpha1().ClusterAlerts(ca.Namespace).Update(&ca); err != nil {
				return err
			}
		}
	}

	// Pause all PodAlerts for this plugin
	paList, err := op.extClient.MonitoringV1alpha1().PodAlerts(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, pa := range paList.Items {
		if pa.Spec.Check == api.CheckPod(name) {
			pa.Spec.Paused = true
			if _, err := op.extClient.MonitoringV1alpha1().PodAlerts(pa.Namespace).Update(&pa); err != nil {
				return err
			}
		}
	}

	// Pause all NodeAlerts for this plugin
	naList, err := op.extClient.MonitoringV1alpha1().NodeAlerts(core.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, na := range naList.Items {
		if na.Spec.Check == api.CheckNode(name) {
			na.Spec.Paused = true
			if _, err := op.extClient.MonitoringV1alpha1().NodeAlerts(na.Namespace).Update(&na); err != nil {
				return err
			}
		}
	}

	// Pausing Alerts may take times. In that case, removing CheckCommand will cause panic in Icinga API
	// That's why deleting all Icinga2 Service Objects with check_command matched with this plugin.
	// We can confirm that removing CheckCommand config is safe now
	if err := icinga.NewClusterHost(op.icingaClient, "").DeleteForCheckCommand(name); err != nil {
		return err
	}
	if err := icinga.NewNodeHost(op.icingaClient, "").DeleteForCheckCommand(name); err != nil {
		return err
	}
	if err := icinga.NewPodHost(op.icingaClient, "").DeleteForCheckCommand(name); err != nil {
		return err
	}

	// Delete IcingaCommand definition from Maps
	delete(api.ClusterCommands, api.CheckCluster(name))
	delete(api.NodeCommands, api.CheckNode(name))
	delete(api.PodCommands, api.CheckPod(name))

	// Remove CheckCommand config file from custom.d folder
	if err := ioutil.EnsureDirectory(filepath.Join(op.ConfigRoot, "custom.d")); err != nil {
		return err
	}
	path := filepath.Join(op.ConfigRoot, "custom.d", fmt.Sprintf("%s.conf", name))
	if err := os.Remove(path); err != nil {
		return err
	}

	// Restart Icinga2
	if err := op.restartIcinga2(pod); err != nil {
		return err
	}

	return nil
}

func (op *Operator) addPluginSupport(wp *api.SearchlightPlugin) error {

	checkCommandString := plugin.GenerateCheckCommand(wp)

	pod, err := op.GetIcingaPod()
	if err != nil {
		return errors.WithMessage(err, "failed to get Icinga2 Pod name.")
	}

	if err := ioutil.EnsureDirectory(filepath.Join(op.ConfigRoot, "custom.d")); err != nil {
		return err
	}

	path := filepath.Join(op.ConfigRoot, "custom.d", fmt.Sprintf("%s.conf", wp.Name))

	if !ioutil.WriteString(path, checkCommandString) {
		return fmt.Errorf(`failed to write CheckCommand "%s" in %s`, wp.Name, path)
	}

	return op.restartIcinga2(pod)
}

func (op *Operator) restartIcinga2(pod *core.Pod) error {
	podExecOptions := &core.PodExecOptions{
		Container: "icinga",
		Command:   []string{"sh", "-c", "kill -9 $(cat /run/icinga2/icinga2.pid)"},
		Stdout:    true,
		Stderr:    true,
	}

	_, err := op.executeCommand(pod, podExecOptions)
	if err != nil {
		return err
	}

	return nil
}

func (op *Operator) GetIcingaPod() (*core.Pod, error) {
	namespace := meta.Namespace()
	podName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
		},
	}, nil
}

func (op *Operator) executeCommand(pod *core.Pod, podExecOptions *core.PodExecOptions) (string, error) {
	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	req := op.kubeClient.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")

	req.VersionedParams(podExecOptions, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(op.clientConfig, "POST", req.URL())
	if err != nil {
		return "", err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &execOut,
		Stderr: &execErr,
	})

	if err != nil {
		return "", err
	}

	if execErr.Len() > 0 {
		return "", errors.New("failed to exec to restart Icinga2")
	}

	return execOut.String(), nil
}
