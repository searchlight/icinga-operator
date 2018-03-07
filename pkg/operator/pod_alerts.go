package operator

import (
	"strings"

	"github.com/appscode/go/log"
	"github.com/appscode/go/sets"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initPodAlertWatcher() {
	op.paInformer = op.monInformerFactory.Monitoring().V1alpha1().PodAlerts().Informer()
	op.paQueue = queue.New("PodAlert", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcilePodAlert)
	op.paInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			alert := obj.(*api.PodAlert)
			if err := op.isValid(alert); err == nil {
				queue.Enqueue(op.paQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*api.PodAlert)
			nu := newObj.(*api.PodAlert)

			if err := op.isValid(nu); err != nil {
				return
			}
			if !equalPodAlert(old, nu) {
				queue.Enqueue(op.paQueue.GetQueue(), nu)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.paQueue.GetQueue(), obj)
		},
	})
	op.paLister = op.monInformerFactory.Monitoring().V1alpha1().PodAlerts().Lister()
}

func (op *Operator) reconcilePodAlert(key string) error {
	obj, exists, err := op.paInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		log.Warningf("PodAlert %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		return op.EnsurePodAlertDeleted(namespace, name)
	}

	alert := obj.(*api.PodAlert)
	log.Infof("Sync/Add/Update for PodAlert %s\n", alert.GetName())

	op.EnsurePodAlert(alert)
	op.EnsurePodAlertDeleted(alert.Namespace, alert.Name)
	return nil
}

func (op *Operator) EnsurePodAlert(alert *api.PodAlert) error {
	if alert.Spec.PodName != nil {
		pod, err := op.podLister.Pods(alert.Namespace).Get(*alert.Spec.PodName)
		if err != nil {
			return err
		}
		key, err := cache.MetaNamespaceKeyFunc(pod)
		if err == nil {
			op.podQueue.GetQueue().Add(key)
		}
	}

	sel, err := metav1.LabelSelectorAsSelector(alert.Spec.Selector)
	if err != nil {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonAlertInvalid,
			"Reason: %s",
			err.Error(),
		)
		return err
	}
	pods, err := op.podLister.Pods(alert.Namespace).List(sel)
	if err != nil {
		return err
	}
	for i := range pods {
		pod := pods[i]
		key, err := cache.MetaNamespaceKeyFunc(pod)
		if err == nil {
			op.nodeQueue.GetQueue().Add(key)
		}
	}
	return nil
}

func GetAppliedPodAlerts(a map[string]string, key string) bool {
	if a == nil {
		return false
	}
	if val, ok := a[annotationAlertsName]; ok {
		names := strings.Split(val, ",")
		return sets.NewString(names...).Has(key)
	}
	return false
}

func (op *Operator) EnsurePodAlertDeleted(alertNamespace, alertName string) error {
	pods, err := op.podLister.Pods(alertNamespace).List(labels.Everything())
	if err != nil {
		return err
	}
	for _, pod := range pods {
		if GetAppliedNodeAlerts(pod.Annotations, alertName) {
			key, err := cache.MetaNamespaceKeyFunc(pod)
			if err == nil {
				op.podQueue.GetQueue().Add(key)
			}
		}
	}
	return nil
}

func (op *Operator) EnsureIcingaPodAlert(alert *api.PodAlert, pod *core.Pod) (err error) {
	err = op.podHost.Reconcile(alert.DeepCopy(), pod.DeepCopy())
	if err != nil {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToSync,
			`Reason: %v`,
			alert.Name,
			err,
		)
	}
	return err
}

func (op *Operator) EnsureIcingaPodAlertDeleted(alert *api.PodAlert, pod *core.Pod) (err error) {
	err = op.podHost.Delete(alert, pod)
	if err != nil && alert != nil {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			`Fail to delete Icinga objects of PodAlert "%s@%s" for Pod "%s@%s". Reason: %v`,
			alert.Name, alert.Namespace, pod.Name, pod.Namespace,
			err,
		)
	}
	return err
}
