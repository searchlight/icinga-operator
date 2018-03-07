package operator

import (
	"encoding/json"
	"reflect"

	"github.com/appscode/go/log"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initPodWatcher() {
	op.pInformer = op.kubeInformerFactory.Core().V1().Pods().Informer()
	op.pQueue = queue.New("Pod", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcilePod)
	op.pInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if pod, ok := obj.(*core.Pod); ok {
				if pod.Status.PodIP == "" {
					log.Warningf("Skipping pod %s@%s, since it has no IP", pod.Name, pod.Namespace)
					return
				}
			}
			queue.Enqueue(op.pQueue.GetQueue(), obj)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldPod, ok := old.(*core.Pod)
			if !ok {
				log.Errorln("invalid Pod object")
				return
			}
			newPod, ok := new.(*core.Pod)
			if !ok {
				log.Errorln("invalid Pod object")
				return
			}
			if !reflect.DeepEqual(oldPod.Labels, newPod.Labels) || oldPod.Status.PodIP != newPod.Status.PodIP {
				queue.Enqueue(op.pQueue.GetQueue(), new)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.pQueue.GetQueue(), obj)
		},
	})
	op.pLister = op.kubeInformerFactory.Core().V1().Pods().Lister()
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) reconcilePod(key string) error {
	obj, exists, err := op.pInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		log.Debugf("Pod %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}

		if err := op.ForceDeleteIcingaObjectsForPod(namespace, name); err != nil {
			log.Errorf("Failed to delete alert for Pod %s@%s", name, namespace)
		}
	} else {
		pod := obj.(*core.Pod)
		if err := op.EnsurePod(pod); err != nil {
			log.Errorf("Failed to patch alert for Pod %s@%s", pod.Name, pod.Namespace)
		}
	}
	return nil
}

func (op *Operator) EnsurePod(pod *core.Pod) error {
	log.Infof("Sync/Add/Update for Pod %s\n", pod.GetName())

	oldAlerts := make([]*api.PodAlert, 0)

	oldAlertNames := make([]string, 0)
	if val, ok := pod.Annotations[annotationAlertsName]; ok {
		if err := json.Unmarshal([]byte(val), &oldAlertNames); err != nil {
			return err
		}
	}
	for _, l := range oldAlertNames {
		oldAlerts = append(oldAlerts, &api.PodAlert{
			ObjectMeta: metav1.ObjectMeta{
				Name:      l,
				Namespace: pod.Namespace,
			},
		})
	}

	newAlerts, err := util.FindPodAlert(op.paLister, pod.ObjectMeta)
	if err != nil {
		return err
	}

	type change struct {
		old *api.PodAlert
		new *api.PodAlert
	}
	diff := make(map[string]*change)
	for i := range oldAlerts {
		diff[oldAlerts[i].Name] = &change{old: oldAlerts[i]}
	}

	alertNames := make([]string, 0)

	for i := range newAlerts {
		alertNames = append(alertNames, newAlerts[i].Name)
		if ch, ok := diff[newAlerts[i].Name]; ok {
			ch.new = newAlerts[i]
		} else {
			diff[newAlerts[i].Name] = &change{new: newAlerts[i]}
		}
	}

	for alert := range diff {
		ch := diff[alert]
		if ch.old != nil && ch.new == nil {
			op.EnsureIcingaPodAlertDeleted(ch.old, pod)
		} else {
			op.EnsureIcingaPodAlert(ch.new, pod)
		}
	}

	_, vr, err := core_util.PatchPod(op.KubeClient, pod, func(in *core.Pod) *core.Pod {
		if len(newAlerts) > 0 {
			if in.Annotations == nil {
				in.Annotations = make(map[string]string, 0)
			}
			data, _ := json.Marshal(alertNames)
			in.Annotations[annotationAlertsName] = string(data)
		} else {
			delete(in.Annotations, annotationAlertsName)
		}

		return in
	})
	if err != nil {
		log.Errorf("Failed to %v Pod %s@%s.", vr, pod.Name, pod.Namespace)
	}

	return nil
}

func (op *Operator) ForceDeleteIcingaObjectsForPod(namespace, name string) (err error) {
	defer func() {
		if err != nil {
			log.Errorln(err)
			return
		}
	}()

	err = op.podHost.ForceDeleteIcingaHost(icinga.IcingaHost{
		Type:           icinga.TypePod,
		AlertNamespace: namespace,
		ObjectName:     name,
	})
	return
}
