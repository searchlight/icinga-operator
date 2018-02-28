package operator

import (
	"fmt"
	"reflect"

	"github.com/appscode/go/log"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initPodAlertWatcher() {
	op.paInformer = op.searchlightInformerFactory.Monitoring().V1alpha1().PodAlerts().Informer()
	op.paQueue = queue.New("PodAlert", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcilePodAlert)
	op.paInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if alert, ok := obj.(*api.PodAlert); ok {
				if ok, err := alert.IsValid(); !ok {
					op.recorder.Eventf(
						alert.ObjectReference(),
						core.EventTypeWarning,
						eventer.EventReasonFailedToCreate,
						`Reason: %v`,
						alert.Name,
						err,
					)
					return
				}
				queue.Enqueue(op.paQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			oldAlert, ok := oldObj.(*api.PodAlert)
			if !ok {
				log.Errorln("invalid PodAlert object")
				return
			}
			newAlert, ok := newObj.(*api.PodAlert)
			if !ok {
				log.Errorln("invalid PodAlert object")
				return
			}
			if ok, err := newAlert.IsValid(); !ok {
				op.recorder.Eventf(
					newAlert.ObjectReference(),
					core.EventTypeWarning,
					eventer.EventReasonFailedToDelete,
					`Reason: %v`,
					newAlert.Name,
					err,
				)
				return
			}
			if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
				queue.Enqueue(op.paQueue.GetQueue(), newObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.paQueue.GetQueue(), obj)
		},
	})
	op.paLister = op.searchlightInformerFactory.Monitoring().V1alpha1().PodAlerts().Lister()
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) reconcilePodAlert(key string) error {
	obj, exists, err := op.paInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		alert := obj.(*api.PodAlert)
		// Below we will warm up our cache with a PodAlert, so that we will see a delete for one d
		fmt.Printf("PodAlert %s does not exist anymore\n", key)
		err := op.EnsurePodAlertDeleted(alert)
		if err != nil {
			log.Errorf("Failed to delete PodAlert %s@%s", alert.Name, alert.Namespace)
		}
	} else {
		alert := obj.(*api.PodAlert)
		fmt.Printf("Sync/Add/Update for PodAlert %s\n", alert.GetName())

		if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonBadNotifier,
				`Bad notifier config for PodAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			return err
		}

		err := op.EnsurePodAlert(alert)
		if err != nil {
			log.Errorf("Failed to patch PodAlert %s@%s", alert.Name, alert.Namespace)
		}

	}
	return nil
}

func (op *Operator) EnsurePodAlert(new *api.PodAlert) error {
	newSel, err := metav1.LabelSelectorAsSelector(&new.Spec.Selector)
	if err != nil {
		return err
	}
	if new.Spec.PodName != "" {
		if resource, err := op.KubeClient.CoreV1().Pods(new.Namespace).Get(new.Spec.PodName, metav1.GetOptions{}); err == nil {
			if newSel.Matches(labels.Set(resource.Labels)) {
				if resource.Status.PodIP == "" {
					log.Warningf("Skipping pod %s@%s, since it has no IP", resource.Name, resource.Namespace)
				}
				go op.EnsureIcingaPodAlert(resource, new)
			}
		}
	} else {
		if resources, err := op.KubeClient.CoreV1().Pods(new.Namespace).List(metav1.ListOptions{LabelSelector: newSel.String()}); err == nil {
			for i := range resources.Items {
				resource := resources.Items[i]
				if resource.Status.PodIP == "" {
					log.Warningf("Skipping pod %s@%s, since it has no IP", resource.Name, resource.Namespace)
					continue
				}
				go op.EnsureIcingaPodAlert(&resource, new)
			}
		}
	}

	return nil
}

func (op *Operator) EnsurePodAlertDeleted(alert *api.PodAlert) error {
	sel, err := metav1.LabelSelectorAsSelector(&alert.Spec.Selector)
	if err != nil {
		return err
	}
	if alert.Spec.PodName != "" {
		if resource, err := op.KubeClient.CoreV1().Pods(alert.Namespace).Get(alert.Spec.PodName, metav1.GetOptions{}); err == nil {
			if sel.Matches(labels.Set(resource.Labels)) {
				go op.EnsureIcingaPodAlertDeleted(resource, alert)
			}
		}
	} else {
		if resources, err := op.KubeClient.CoreV1().Pods(alert.Namespace).List(metav1.ListOptions{LabelSelector: sel.String()}); err == nil {
			for i := range resources.Items {
				go op.EnsureIcingaPodAlertDeleted(&resources.Items[i], alert)
			}
		}
	}

	return nil
}

func (op *Operator) EnsureIcingaPodAlert(pod *core.Pod, new *api.PodAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				new.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulSync,
				`Applied PodAlert: "%v"`,
				new.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				new.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`Fail to be apply PodAlert: "%v". Reason: %v`,
				new.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()

	err = op.podHost.Create(new, pod)
	return
}

func (op *Operator) EnsureIcingaPodAlertDeleted(pod *core.Pod, alert *api.PodAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulDelete,
				`Deleted PodAlert: "%v"`,
				alert.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				`Fail to be delete PodAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()

	err = op.podHost.Delete(alert.Namespace, alert.Name, *pod)
	return
}
