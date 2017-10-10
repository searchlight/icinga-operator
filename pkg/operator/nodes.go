package operator

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/appscode/go/log"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rt "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	core_listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

func (op *Operator) initNodeWatcher() {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (rt.Object, error) {
			return op.KubeClient.CoreV1().Nodes().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Nodes().Watch(options)
		},
	}

	// create the workqueue
	op.nQueue = workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "node")

	// Bind the workqueue to a cache with the help of an informer. This way we make sure that
	// whenever the cache is updated, the pod key is added to the workqueue.
	// Note that when we finally process the item from the workqueue, we might see a newer version
	// of the Node than the version which was responsible for triggering the update.
	op.nIndexer, op.nInformer = cache.NewIndexerInformer(lw, &core.Node{}, op.Opt.ResyncPeriod, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				op.nQueue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				op.nQueue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			// IndexerInformer uses a delta queue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				op.nQueue.Add(key)
			}
		},
	}, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
	op.nLister = core_listers.NewNodeLister(op.nIndexer)
}

func (op *Operator) runNodeWatcher() {
	for op.processNextNode() {
	}
}

func (op *Operator) processNextNode() bool {
	// Wait until there is a new item in the working queue
	key, quit := op.nQueue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two deployments with the same key are never processed in
	// parallel.
	defer op.nQueue.Done(key)

	// Invoke the method containing the business logic
	err := op.runNodeInjector(key.(string))
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		op.nQueue.Forget(key)
		return true
	}
	log.Errorf("Failed to process Node %v. Reason: %s", key, err)

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if op.nQueue.NumRequeues(key) < op.Opt.MaxNumRequeues {
		glog.Infof("Error syncing deployment %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		op.nQueue.AddRateLimited(key)
		return true
	}

	op.nQueue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping deployment %q out of the queue: %v", key, err)
	return true
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) runNodeInjector(key string) error {
	obj, exists, err := op.nIndexer.GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		// Below we will warm up our cache with a Node, so that we will see a delete for one d
		fmt.Printf("Node %s does not exist anymore\n", key)
	} else {
		pod := obj.(*core.Node)
		fmt.Printf("Sync/Add/Update for Node %s\n", pod.GetName())

	}
	return nil
}

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchNodes() {
	defer runtime.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (rt.Object, error) {
			return op.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Nodes().Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&core.Node{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if resource, ok := obj.(*core.Node); ok {
					log.Infof("Node %s@%s added", resource.Name, resource.Namespace)

					alerts, err := util.FindNodeAlert(op.ExtClient, resource.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					if len(alerts) == 0 {
						log.Errorf("No NodeAlert found for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					for i := range alerts {
						err = op.EnsureNode(resource, nil, alerts[i])
						if err != nil {
							log.Errorf("Failed to add icinga2 alert for Node %s@%s.", resource.Name, resource.Namespace)
							// return
						}
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldNode, ok := old.(*core.Node)
				if !ok {
					log.Errorln(errors.New("Invalid Node object"))
					return
				}
				newNode, ok := new.(*core.Node)
				if !ok {
					log.Errorln(errors.New("Invalid Node object"))
					return
				}
				if !reflect.DeepEqual(oldNode.Labels, newNode.Labels) {
					oldAlerts, err := util.FindNodeAlert(op.ExtClient, oldNode.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", oldNode.Name, oldNode.Namespace)
						return
					}
					newAlerts, err := util.FindNodeAlert(op.ExtClient, newNode.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", newNode.Name, newNode.Namespace)
						return
					}

					type change struct {
						old *api.NodeAlert
						new *api.NodeAlert
					}
					diff := make(map[string]*change)
					for i := range oldAlerts {
						diff[oldAlerts[i].Name] = &change{old: oldAlerts[i]}
					}
					for i := range newAlerts {
						if ch, ok := diff[newAlerts[i].Name]; ok {
							ch.new = newAlerts[i]
						} else {
							diff[newAlerts[i].Name] = &change{new: newAlerts[i]}
						}
					}
					for alert := range diff {
						ch := diff[alert]
						if ch.old == nil && ch.new != nil {
							go op.EnsureNode(newNode, nil, ch.new)
						} else if ch.old != nil && ch.new == nil {
							go op.EnsureNodeDeleted(newNode, ch.old)
						} else if ch.old != nil && ch.new != nil && !reflect.DeepEqual(ch.old.Spec, ch.new.Spec) {
							go op.EnsureNode(newNode, ch.old, ch.new)
						}
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if resource, ok := obj.(*core.Node); ok {
					log.Infof("Node %s@%s deleted", resource.Name, resource.Namespace)

					alerts, err := util.FindNodeAlert(op.ExtClient, resource.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					if len(alerts) == 0 {
						log.Errorf("No NodeAlert found for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					for i := range alerts {
						err = op.EnsureNodeDeleted(resource, alerts[i])
						if err != nil {
							log.Errorf("Failed to delete icinga2 alert for Node %s@%s.", resource.Name, resource.Namespace)
							// return
						}
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (op *Operator) EnsureNode(node *core.Node, old, new *api.NodeAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				new.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulSync,
				`Applied NodeAlert: "%v"`,
				new.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				new.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`Fail to be apply NodeAlert: "%v". Reason: %v`,
				new.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()

	if old == nil {
		err = op.nodeHost.Create(*new, *node)
	} else {
		err = op.nodeHost.Update(*new, *node)
	}
	return
}

func (op *Operator) EnsureNodeDeleted(node *core.Node, alert *api.NodeAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulDelete,
				`Deleted NodeAlert: "%v"`,
				alert.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				`Fail to be delete NodeAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()
	err = op.nodeHost.Delete(alert.Namespace, alert.Name, *node)
	return
}
