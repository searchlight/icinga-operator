package operator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/appscode/go/log"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	slite_listers "github.com/appscode/searchlight/listers/monitoring/v1alpha1"
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
	"k8s.io/client-go/util/workqueue"
)

func (op *Operator) initNodeAlertWatcher() {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (rt.Object, error) {
			return op.ExtClient.NodeAlerts(core.NamespaceAll).List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.ExtClient.NodeAlerts(core.NamespaceAll).Watch(options)
		},
	}

	// create the workqueue
	op.naQueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "node-alert")

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the NodeAlert than the version which was responsible for triggering the update.
	op.naIndexer, op.naInformer = cache.NewIndexerInformer(lw, &api.NodeAlert{}, op.Opt.ResyncPeriod, cache.ResourceEventHandlerFuncs{
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
				} else {
					key, err := cache.MetaNamespaceKeyFunc(obj)
					if err == nil {
						op.naQueue.Add(key)
					}
				}
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldAlert, ok := old.(*api.NodeAlert)
			if !ok {
				log.Errorln("Invalid NodeAlert object")
				return
			}
			newAlert, ok := new.(*api.NodeAlert)
			if !ok {
				log.Errorln("Invalid NodeAlert object")
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
				key, err := cache.MetaNamespaceKeyFunc(new)
				if err == nil {
					op.naQueue.Add(key)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				op.naQueue.Add(key)
			}
		},
	}, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	op.naLister = slite_listers.NewNodeAlertLister(op.naIndexer)
}

func (op *Operator) runNodeAlertWatcher() {
	for op.processNextNodeAlert() {
	}
}

func (op *Operator) processNextNodeAlert() bool {
	// Wait until there is a new item in the working queue
	key, quit := op.naQueue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two deployments with the same key are never processed in
	// parallel.
	defer op.naQueue.Done(key)

	// Invoke the method containing the business logic
	err := op.runNodeAlertInjector(key.(string))
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		op.naQueue.Forget(key)
		return true
	}
	log.Errorf("Failed to process NodeAlert %v. Reason: %s", key, err)

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if op.naQueue.NumRequeues(key) < op.Opt.MaxNumRequeues {
		glog.Infof("Error syncing deployment %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		op.naQueue.AddRateLimited(key)
		return true
	}

	op.naQueue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping deployment %q out of the queue: %v", key, err)
	return true
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) runNodeAlertInjector(key string) error {
	obj, exists, err := op.naIndexer.GetByKey(key)
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
			return op.ExtClient.NodeAlerts(core.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.ExtClient.NodeAlerts(core.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&api.NodeAlert{},
		op.Opt.ResyncPeriod,
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
