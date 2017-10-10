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

func (op *Operator) initPodAlertWatcher() {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (rt.Object, error) {
			return op.ExtClient.PodAlerts(core.NamespaceAll).List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.ExtClient.PodAlerts(core.NamespaceAll).Watch(options)
		},
	}

	// create the workqueue
	op.paQueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "pod-alert")

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the PodAlert than the version which was responsible for triggering the update.
	op.paIndexer, op.paInformer = cache.NewIndexerInformer(lw, &api.PodAlert{}, op.Opt.ResyncPeriod, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				op.paQueue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				op.paQueue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				op.paQueue.Add(key)
			}
		},
	}, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	op.paLister = slite_listers.NewPodAlertLister(op.paIndexer)
}

func (op *Operator) runPodAlertWatcher() {
	for op.processNextPodAlert() {
	}
}

func (op *Operator) processNextPodAlert() bool {
	// Wait until there is a new item in the working queue
	key, quit := op.paQueue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two deployments with the same key are never processed in
	// parallel.
	defer op.paQueue.Done(key)

	// Invoke the method containing the business logic
	err := op.runPodAlertInjector(key.(string))
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		op.paQueue.Forget(key)
		return true
	}
	log.Errorf("Failed to process PodAlert %v. Reason: %s", key, err)

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if op.paQueue.NumRequeues(key) < op.Opt.MaxNumRequeues {
		glog.Infof("Error syncing deployment %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		op.paQueue.AddRateLimited(key)
		return true
	}

	op.paQueue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping deployment %q out of the queue: %v", key, err)
	return true
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) runPodAlertInjector(key string) error {
	obj, exists, err := op.paIndexer.GetByKey(key)
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
			return op.ExtClient.PodAlerts(core.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.ExtClient.PodAlerts(core.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&api.PodAlert{},
		op.Opt.ResyncPeriod,
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
					log.Errorln(errors.New("Invalid PodAlert object"))
					return
				}
				newAlert, ok := new.(*api.PodAlert)
				if !ok {
					log.Errorln(errors.New("Invalid PodAlert object"))
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
