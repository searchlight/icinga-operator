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

func (op *Operator) initNodeAlertWatcher() {
	op.naInformer = op.searchlightInformerFactory.Monitoring().V1alpha1().NodeAlerts().Informer()
	op.naQueue = queue.New("NodeAlert", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcileNodeAlert)
	op.naInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if alert, ok := obj.(*api.NodeAlert); ok {
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
				queue.Enqueue(op.naQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldAlert, ok := old.(*api.NodeAlert)
			if !ok {
				log.Errorln("invalid NodeAlert object")
				return
			}
			newAlert, ok := new.(*api.NodeAlert)
			if !ok {
				log.Errorln("invalid NodeAlert object")
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
				queue.Enqueue(op.naQueue.GetQueue(), new)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.naQueue.GetQueue(), obj)
		},
	})
	op.naLister = op.searchlightInformerFactory.Monitoring().V1alpha1().NodeAlerts().Lister()
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) reconcileNodeAlert(key string) error {
	obj, exists, err := op.naInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a NodeAlert, so that we will see a delete for one d
		fmt.Printf("NodeAlert %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		return op.clusterHost.Delete(namespace, name)
	} else {
		a := obj.(*api.NodeAlert)
		fmt.Printf("Sync/Add/Update for NodeAlert %s\n", a.GetName())

	}
	return nil
}

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchNodeAlerts() {
	defer runtime.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (rt.Object, error) {
			return op.ExtClient.MonitoringV1alpha1().NodeAlerts(core.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.ExtClient.MonitoringV1alpha1().NodeAlerts(core.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&api.NodeAlert{},
		op.options.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*api.NodeAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						op.recorder.Eventf(
							alert.ObjectReference(),
							core.EventTypeWarning,
							eventer.EventReasonFailedToCreate,
							`Fail to be create NodeAlert: "%v". Reason: %v`,
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
							`Bad notifier config for NodeAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
					}
					op.EnsureNodeAlert(nil, alert)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldAlert, ok := old.(*api.NodeAlert)
				if !ok {
					log.Errorln(errors.New("invalid NodeAlert object"))
					return
				}
				newAlert, ok := new.(*api.NodeAlert)
				if !ok {
					log.Errorln(errors.New("invalid NodeAlert object"))
					return
				}
				if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
					if ok, err := newAlert.IsValid(); !ok {
						op.recorder.Eventf(
							newAlert.ObjectReference(),
							core.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be update NodeAlert: "%v". Reason: %v`,
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
							`Bad notifier config for NodeAlert: "%v". Reason: %v`,
							newAlert.Name,
							err,
						)
					}
					op.EnsureNodeAlert(oldAlert, newAlert)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if alert, ok := obj.(*api.NodeAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						op.recorder.Eventf(
							alert.ObjectReference(),
							core.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be delete NodeAlert: "%v". Reason: %v`,
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
							`Bad notifier config for NodeAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
					}
					op.EnsureNodeAlertDeleted(alert)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (op *Operator) EnsureNodeAlert(old, new *api.NodeAlert) {
	oldObjs := make(map[string]*core.Node)

	if old != nil {
		oldSel := labels.SelectorFromSet(old.Spec.Selector)
		if old.Spec.NodeName != "" {
			if resource, err := op.KubeClient.CoreV1().Nodes().Get(old.Spec.NodeName, metav1.GetOptions{}); err == nil {
				if oldSel.Matches(labels.Set(resource.Labels)) {
					oldObjs[resource.Name] = resource
				}
			}
		} else {
			if resources, err := op.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: oldSel.String()}); err == nil {
				for i := range resources.Items {
					oldObjs[resources.Items[i].Name] = &resources.Items[i]
				}
			}
		}
	}

	newSel := labels.SelectorFromSet(new.Spec.Selector)
	if new.Spec.NodeName != "" {
		if resource, err := op.KubeClient.CoreV1().Nodes().Get(new.Spec.NodeName, metav1.GetOptions{}); err == nil {
			if newSel.Matches(labels.Set(resource.Labels)) {
				delete(oldObjs, resource.Name)
				go op.EnsureNode(resource, old, new)
			}
		}
	} else {
		if resources, err := op.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: newSel.String()}); err == nil {
			for i := range resources.Items {
				resource := resources.Items[i]
				delete(oldObjs, resource.Name)
				go op.EnsureNode(&resource, old, new)
			}
		}
	}
	for i := range oldObjs {
		go op.EnsureNodeDeleted(oldObjs[i], old)
	}
}

func (op *Operator) EnsureNodeAlertDeleted(alert *api.NodeAlert) {
	sel := labels.SelectorFromSet(alert.Spec.Selector)
	if alert.Spec.NodeName != "" {
		if resource, err := op.KubeClient.CoreV1().Nodes().Get(alert.Spec.NodeName, metav1.GetOptions{}); err == nil {
			if sel.Matches(labels.Set(resource.Labels)) {
				go op.EnsureNodeDeleted(resource, alert)
			}
		}
	} else {
		if resources, err := op.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: sel.String()}); err == nil {
			for i := range resources.Items {
				go op.EnsureNodeDeleted(&resources.Items[i], alert)
			}
		}
	}
}
