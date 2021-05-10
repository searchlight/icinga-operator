/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/client/clientset/versioned/typed/monitoring/v1alpha1/util"
	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/pkg/plugin"

	"github.com/golang/glog"
	"github.com/pkg/errors"
	"gomodules.xyz/x/ioutil"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/tools/queue"
)

func (op *Operator) initPluginWatcher() {
	op.pluginInformer = op.monInformerFactory.Monitoring().V1alpha1().SearchlightPlugins().Informer()
	op.pluginQueue = queue.New("SearchlightPlugin", op.MaxNumRequeues, op.NumThreads, op.reconcilePlugin)
	op.pluginInformer.AddEventHandler(queue.NewEventHandler(op.pluginQueue.GetQueue(), func(oldObj, newObj interface{}) bool {
		old := oldObj.(*api.SearchlightPlugin)
		nu := newObj.(*api.SearchlightPlugin)
		return !reflect.DeepEqual(old.Spec, nu.Spec)
	}))
	op.pluginLister = op.monInformerFactory.Monitoring().V1alpha1().SearchlightPlugins().Lister()
}

func (op *Operator) reconcilePlugin(key string) error {
	obj, exists, err := op.pluginInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		klog.Warningf("SearchlightPlugin %s does not exist anymore\n", key)

		_, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}

		fmt.Println("deleting CheckCommand for ", name)
		return op.ensureCheckCommandDeleted(name)
	}

	searchlightPlugin := obj.(*api.SearchlightPlugin).DeepCopy()
	klog.Infof("Sync/Add/Update for SearchlightPlugin %s\n", searchlightPlugin.GetName())

	return op.ensureCheckCommand(searchlightPlugin)
}

func (op *Operator) ensureCheckCommand(wp *api.SearchlightPlugin) error {

	ic := api.IcingaCommand{
		Name: wp.Name,
		Vars: &api.PluginVars{
			Fields:   make(map[string]api.PluginVarField),
			Required: make([]string, 0),
		},
	}

	ic.Vars = wp.Spec.Arguments.Vars
	ic.States = wp.Spec.States

	for _, t := range wp.Spec.AlertKinds {
		if t == api.ResourceKindClusterAlert {
			api.ClusterCommands.Insert(wp.Name, ic)
		} else if t == api.ResourceKindNodeAlert {
			api.NodeCommands.Insert(wp.Name, ic)
		} else if t == api.ResourceKindPodAlert {
			api.PodCommands.Insert(wp.Name, ic)
		}
	}

	return op.addPluginSupport(wp)
}

func (op *Operator) ensureCheckCommandDeleted(name string) error {
	pod, err := op.GetIcingaPod()
	if err != nil {
		return errors.WithMessage(err, "failed to get Icinga2 Pod name.")
	}

	var errs []error
	{
		// Pause all ClusterAlerts for this plugin
		err = cache.ListAll(op.caInformer.GetIndexer(), labels.Everything(), func(obj interface{}) {
			_, _, err = util.PatchClusterAlert(context.TODO(), op.extClient.MonitoringV1alpha1(), obj.(*api.ClusterAlert), func(alert *api.ClusterAlert) *api.ClusterAlert {
				pause := alert.Spec.Check == name
				alert.Spec.Paused = pause
				return alert
			}, metav1.PatchOptions{})
			if err != nil {
				errs = append(errs, err)
			}
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	{
		// Pause all PodAlerts for this plugin
		err = cache.ListAll(op.paInformer.GetIndexer(), labels.Everything(), func(obj interface{}) {
			_, _, err = util.PatchPodAlert(context.TODO(), op.extClient.MonitoringV1alpha1(), obj.(*api.PodAlert), func(alert *api.PodAlert) *api.PodAlert {
				pause := alert.Spec.Check == name
				alert.Spec.Paused = pause
				return alert
			}, metav1.PatchOptions{})
			if err != nil {
				errs = append(errs, err)
			}
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	{
		// Pause all NodeAlerts for this plugin
		err = cache.ListAll(op.naInformer.GetIndexer(), labels.Everything(), func(obj interface{}) {
			_, _, err = util.PatchNodeAlert(context.TODO(), op.extClient.MonitoringV1alpha1(), obj.(*api.NodeAlert), func(alert *api.NodeAlert) *api.NodeAlert {
				pause := alert.Spec.Check == name
				alert.Spec.Paused = pause
				return alert
			}, metav1.PatchOptions{})
			if err != nil {
				errs = append(errs, err)
			}
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	// Pausing Alerts may take times. In that case, removing CheckCommand will cause panic in Icinga API
	// That's why deleting all Icinga2 Service Objects with check_command matched with this plugin.
	// We can confirm that removing CheckCommand config is safe now
	if err := icinga.NewClusterHost(op.icingaClient, "").DeleteChecks(name); err != nil {
		return err
	}
	if err := icinga.NewNodeHost(op.icingaClient, "").DeleteChecks(name); err != nil {
		return err
	}
	if err := icinga.NewPodHost(op.icingaClient, "").DeleteChecks(name); err != nil {
		return err
	}

	// Delete IcingaCommand definition from Maps
	api.ClusterCommands.Delete(name)
	api.NodeCommands.Delete(name)
	api.PodCommands.Delete(name)

	// Remove CheckCommand config file from custom.d folder
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

func (op *Operator) createBuiltinSearchlightPlugin() error {
	plugins := []*api.SearchlightPlugin{
		plugin.GetComponentStatusPlugin(),
		plugin.GetJsonPathPlugin(),
		plugin.GetNodeExistsPlugin(),
		plugin.GetPodExistsPlugin(),
		plugin.GetEventPlugin(),
		plugin.GetCACertPlugin(),
		plugin.GetCertPlugin(),
		plugin.GetNodeStatusPlugin(),
		plugin.GetNodeVolumePlugin(),
		plugin.GetPodStatusPlugin(),
		plugin.GetPodVolumePlugin(),
		plugin.GetPodExecPlugin(),
	}

	var errs []error
	for _, p := range plugins {
		_, _, err := util.CreateOrPatchSearchlightPlugin(context.TODO(), op.extClient.MonitoringV1alpha1(), p.ObjectMeta, func(sp *api.SearchlightPlugin) *api.SearchlightPlugin {
			sp.Spec = p.Spec
			return sp
		}, metav1.PatchOptions{})
		if err != nil {
			errs = append(errs, err)
		}
	}

	return utilerrors.NewAggregate(errs)
}
