package operator

import (
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
	rt "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func (op *Operator) initClusterAlertWatcher() {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (rt.Object, error) {
			return op.ExtClient.ClusterAlerts(core.NamespaceAll).List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.ExtClient.ClusterAlerts(core.NamespaceAll).Watch(options)
		},
	}

	// create the workqueue
	op.caQueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "cluster-alert")

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the ClusterAlert than the version which was responsible for triggering the update.
	op.caIndexer, op.caInformer = cache.NewIndexerInformer(lw, &api.ClusterAlert{}, op.Opt.ResyncPeriod, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if alert, ok := obj.(*api.ClusterAlert); ok {
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
						op.caQueue.Add(key)
					}
				}
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldAlert, ok := old.(*api.ClusterAlert)
			if !ok {
				log.Errorln("Invalid ClusterAlert object")
				return
			}
			newAlert, ok := new.(*api.ClusterAlert)
			if !ok {
				log.Errorln("Invalid ClusterAlert object")
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
					op.caQueue.Add(key)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				op.caQueue.Add(key)
			}
		},
	}, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	op.caLister = slite_listers.NewClusterAlertLister(op.caIndexer)
}

func (op *Operator) runClusterAlertWatcher() {
	for op.processNextClusterAlert() {
	}
}

func (op *Operator) processNextClusterAlert() bool {
	// Wait until there is a new item in the working queue
	key, quit := op.caQueue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two deployments with the same key are never processed in
	// parallel.
	defer op.caQueue.Done(key)

	// Invoke the method containing the business logic
	err := op.runClusterAlertInjector(key.(string))
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		op.caQueue.Forget(key)
		return true
	}
	log.Errorf("Failed to process ClusterAlert %v. Reason: %s", key, err)

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if op.caQueue.NumRequeues(key) < op.Opt.MaxNumRequeues {
		glog.Infof("Error syncing deployment %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		op.caQueue.AddRateLimited(key)
		return true
	}

	op.caQueue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping deployment %q out of the queue: %v", key, err)
	return true
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) runClusterAlertInjector(key string) error {
	obj, exists, err := op.caIndexer.GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a ClusterAlert, so that we will see a delete for one d
		fmt.Printf("ClusterAlert %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		return op.clusterHost.Delete(namespace, name)
	} else {
		alert := obj.(*api.ClusterAlert)
		fmt.Printf("Sync/Add/Update for ClusterAlert %s\n", alert.GetName())

		if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonBadNotifier,
				`Bad notifier config for ClusterAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			return err
		}
		op.EnsureClusterAlert(nil, alert)
	}
	return nil
}

func (op *Operator) EnsureClusterAlert(old, new *api.ClusterAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				new.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulSync,
				`Applied ClusterAlert: "%v"`,
				new.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				new.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`Fail to be apply ClusterAlert: "%v". Reason: %v`,
				new.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()

	if old == nil {
		err = op.clusterHost.Create(*new)
	} else {
		err = op.clusterHost.Update(*new)
	}
	return
}
