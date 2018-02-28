package operator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/appscode/go/log"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
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
		pod := obj.(*core.Pod)
		// Below we will warm up our cache with a Pod, so that we will see a delete for one d
		fmt.Printf("Pod %s does not exist anymore\n", key)

		if err := op.EnsurePodDeleted(pod); err != nil {
			log.Errorf("Failed to delete alert for Pod %s@%s", pod.Name, pod.Namespace)
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
	fmt.Printf("Sync/Add/Update for Pod %s\n", pod.GetName())

	oldAlerts := make([]*api.PodAlert, 0)
	if names, ok := pod.Annotations[annotationAlertsName]; ok {
		list := strings.Split(names, ",")
		for _, l := range list {
			oldAlerts = append(oldAlerts, &api.PodAlert{
				ObjectMeta: metav1.ObjectMeta{
					Name:      strings.Trim(l, " "),
					Namespace: pod.Namespace,
				},
			})
		}
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

	newAlertNameList := make([]string, 0)

	for i := range newAlerts {
		newAlertNameList = append(newAlertNameList, newAlerts[i].Name)
		if ch, ok := diff[newAlerts[i].Name]; ok {
			ch.new = newAlerts[i]
		} else {
			diff[newAlerts[i].Name] = &change{new: newAlerts[i]}
		}
	}

	for alert := range diff {
		ch := diff[alert]
		if ch.old == nil && ch.new != nil {
			go op.EnsureIcingaPodAlert(pod, ch.new)
		} else if ch.old != nil && ch.new == nil {
			go op.EnsureIcingaPodAlertDeleted(pod, ch.old)
		} else if ch.old != nil && ch.new != nil && !reflect.DeepEqual(ch.old.Spec, ch.new.Spec) {
			go op.EnsureIcingaPodAlert(pod, ch.new)
		}
	}

	_, vr, err := core_util.PatchPod(op.KubeClient, pod, func(in *core.Pod) *core.Pod {
		if in.Annotations == nil {
			in.Annotations = make(map[string]string, 0)
		}
		in.Annotations[annotationAlertsName] = strings.Join(newAlertNameList, ", ")
		return in
	})
	if err != nil {
		log.Errorf("Failed to %v Pod %s@%s.", vr, pod.Name, pod.Namespace)
	}

	return nil
}

func (op *Operator) EnsurePodDeleted(pod *core.Pod) error {
	alerts, err := util.FindPodAlert(op.paLister, pod.ObjectMeta)
	if err != nil {
		log.Errorf("Error while searching PodAlert for Pod %s@%s.", pod.Name, pod.Namespace)
		return err
	}
	if len(alerts) == 0 {
		log.Errorf("No PodAlert found for Pod %s@%s.", pod.Name, pod.Namespace)
		return err
	}
	for i := range alerts {
		op.EnsureIcingaPodAlertDeleted(pod, alerts[i])
	}

	return nil
}
