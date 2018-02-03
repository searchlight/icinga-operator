package operator

import (
	"errors"
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
	rt "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
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
		// Below we will warm up our cache with a PodAlert, so that we will see a delete for one d
		fmt.Printf("PodAlert %s does not exist anymore\n", key)
	} else {
		a := obj.(*api.PodAlert)
		fmt.Printf("Sync/Add/Update for PodAlert %s\n", a.GetName())

	}
	return nil
}

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchPodAlerts() {
	defer runtime.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (rt.Object, error) {
			return op.ExtClient.MonitoringV1alpha1().PodAlerts(core.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.ExtClient.MonitoringV1alpha1().PodAlerts(core.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&api.PodAlert{},
		op.options.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*api.PodAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						op.recorder.Eventf(
							alert.ObjectReference(),
							core.EventTypeWarning,
							eventer.EventReasonFailedToCreate,
							`Fail to be create PodAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
						op.recorder.Eventf(
							alert.ObjectReference(),
							core.EventTypeWarning,
							eventer.EventReasonBadNotifier,
							`Bad notifier config for PodAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
					}
					op.EnsurePodAlert(nil, alert)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldAlert, ok := old.(*api.PodAlert)
				if !ok {
					log.Errorln(errors.New("invalid PodAlert object"))
					return
				}
				newAlert, ok := new.(*api.PodAlert)
				if !ok {
					log.Errorln(errors.New("invalid PodAlert object"))
					return
				}
				if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
					if ok, err := newAlert.IsValid(); !ok {
						op.recorder.Eventf(
							newAlert.ObjectReference(),
							core.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be update PodAlert: "%v". Reason: %v`,
							newAlert.Name,
							err,
						)
						return
					}
					if err := util.CheckNotifiers(op.KubeClient, newAlert); err != nil {
						op.recorder.Eventf(
							newAlert.ObjectReference(),
							core.EventTypeWarning,
							eventer.EventReasonBadNotifier,
							`Bad notifier config for PodAlert: "%v". Reason: %v`,
							newAlert.Name,
							err,
						)
					}
					op.EnsurePodAlert(oldAlert, newAlert)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if alert, ok := obj.(*api.PodAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						op.recorder.Eventf(
							alert.ObjectReference(),
							core.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be delete PodAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
						op.recorder.Eventf(
							alert.ObjectReference(),
							core.EventTypeWarning,
							eventer.EventReasonBadNotifier,
							`Bad notifier config for PodAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
					}
					op.EnsurePodAlertDeleted(alert)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (op *Operator) EnsurePodAlert(old, new *api.PodAlert) {
	oldObjs := make(map[string]*core.Pod)

	if old != nil {
		oldSel, err := metav1.LabelSelectorAsSelector(&old.Spec.Selector)
		if err != nil {
			return
		}
		if old.Spec.PodName != "" {
			if resource, err := op.KubeClient.CoreV1().Pods(old.Namespace).Get(old.Spec.PodName, metav1.GetOptions{}); err == nil {
				if oldSel.Matches(labels.Set(resource.Labels)) {
					oldObjs[resource.Name] = resource
				}
			}
		} else {
			if resources, err := op.KubeClient.CoreV1().Pods(old.Namespace).List(metav1.ListOptions{LabelSelector: oldSel.String()}); err == nil {
				for i := range resources.Items {
					oldObjs[resources.Items[i].Name] = &resources.Items[i]
				}
			}
		}
	}

	newSel, err := metav1.LabelSelectorAsSelector(&new.Spec.Selector)
	if err != nil {
		return
	}
	if new.Spec.PodName != "" {
		if resource, err := op.KubeClient.CoreV1().Pods(new.Namespace).Get(new.Spec.PodName, metav1.GetOptions{}); err == nil {
			if newSel.Matches(labels.Set(resource.Labels)) {
				delete(oldObjs, resource.Name)
				if resource.Status.PodIP == "" {
					log.Warningf("Skipping pod %s@%s, since it has no IP", resource.Name, resource.Namespace)
				}
				go op.EnsurePod(resource, old, new)
			}
		}
	} else {
		if resources, err := op.KubeClient.CoreV1().Pods(new.Namespace).List(metav1.ListOptions{LabelSelector: newSel.String()}); err == nil {
			for i := range resources.Items {
				resource := resources.Items[i]
				delete(oldObjs, resource.Name)
				if resource.Status.PodIP == "" {
					log.Warningf("Skipping pod %s@%s, since it has no IP", resource.Name, resource.Namespace)
					continue
				}
				go op.EnsurePod(&resource, old, new)
			}
		}
	}
	for i := range oldObjs {
		go op.EnsurePodDeleted(oldObjs[i], old)
	}
}

func (op *Operator) EnsurePodAlertDeleted(alert *api.PodAlert) {
	sel, err := metav1.LabelSelectorAsSelector(&alert.Spec.Selector)
	if err != nil {
		return
	}
	if alert.Spec.PodName != "" {
		if resource, err := op.KubeClient.CoreV1().Pods(alert.Namespace).Get(alert.Spec.PodName, metav1.GetOptions{}); err == nil {
			if sel.Matches(labels.Set(resource.Labels)) {
				go op.EnsurePodDeleted(resource, alert)
			}
		}
	} else {
		if resources, err := op.KubeClient.CoreV1().Pods(alert.Namespace).List(metav1.ListOptions{LabelSelector: sel.String()}); err == nil {
			for i := range resources.Items {
				go op.EnsurePodDeleted(&resources.Items[i], alert)
			}
		}
	}
}
