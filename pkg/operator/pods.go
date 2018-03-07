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
			pod := obj.(*core.Pod)
			if pod.Status.PodIP != "" {
				log.Warningf("Skipping pod %s/%s, since it has no IP", pod.Namespace, pod.Name)
				return
			}
			queue.Enqueue(op.pQueue.GetQueue(), obj)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*core.Pod)
			nu := newObj.(*core.Pod)
			if !reflect.DeepEqual(old.Labels, nu.Labels) || old.Status.PodIP != nu.Status.PodIP {
				queue.Enqueue(op.pQueue.GetQueue(), newObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.pQueue.GetQueue(), obj)
		},
	})
	op.pLister = op.kubeInformerFactory.Core().V1().Pods().Lister()
}

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
		err = op.podHost.ForceDeleteIcingaHost(icinga.IcingaHost{
			Type:           icinga.TypePod,
			AlertNamespace: namespace,
			ObjectName:     name,
		})
		if err != nil {
			log.Errorf("Failed to delete alert for Pod %s@%s", name, namespace)
		}
	} else {
		log.Infof("Sync/Add/Update for Pod %s\n", key)

		pod := obj.(*core.Pod)
		if err := op.EnsurePod(pod); err != nil {
			log.Errorf("Failed to patch alert for Pod %s@%s", pod.Name, pod.Namespace)
		}
	}
	return nil
}

func (op *Operator) EnsurePod(pod *core.Pod) error {
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
