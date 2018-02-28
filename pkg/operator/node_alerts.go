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
		alert := obj.(*api.NodeAlert)
		// Below we will warm up our cache with a NodeAlert, so that we will see a delete for one d
		fmt.Printf("NodeAlert %s does not exist anymore\n", key)
		op.EnsureNodeAlertDeleted(alert)
	} else {
		alert := obj.(*api.NodeAlert)
		fmt.Printf("Sync/Add/Update for NodeAlert %s\n", alert.GetName())

		if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonBadNotifier,
				`Bad notifier config for NodeAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			return err
		}
		op.EnsureNodeAlert(alert)
	}
	return nil
}

func (op *Operator) EnsureNodeAlert(alert *api.NodeAlert) {
	sel := labels.SelectorFromSet(alert.Spec.Selector)
	if alert.Spec.NodeName != "" {
		if node, err := op.KubeClient.CoreV1().Nodes().Get(alert.Spec.NodeName, metav1.GetOptions{}); err == nil {
			if sel.Matches(labels.Set(node.Labels)) {
				go op.EnsureIcingaNodeAlert(alert, node)
			}
		}
	} else {
		if nodes, err := op.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: sel.String()}); err == nil {
			for i := range nodes.Items {
				node := nodes.Items[i]
				go op.EnsureIcingaNodeAlert(alert, &node)
			}
		}
	}
	return
}

func (op *Operator) EnsureNodeAlertDeleted(alert *api.NodeAlert) {
	sel := labels.SelectorFromSet(alert.Spec.Selector)
	if alert.Spec.NodeName != "" {
		if node, err := op.KubeClient.CoreV1().Nodes().Get(alert.Spec.NodeName, metav1.GetOptions{}); err == nil {
			if sel.Matches(labels.Set(node.Labels)) {
				op.EnsureIcingaNodeAlertDeleted(alert, node)
			}
		}
	} else {
		if nodes, err := op.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: sel.String()}); err == nil {
			for i := range nodes.Items {
				node := nodes.Items[i]
				go op.EnsureIcingaNodeAlertDeleted(alert, &node)
			}
		}
	}
	return
}

func (op *Operator) EnsureIcingaNodeAlert(alert *api.NodeAlert, node *core.Node) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulSync,
				`Applied NodeAlert: "%v"`,
				alert.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`Fail to be apply NodeAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()

	err = op.nodeHost.Create(alert, node)
	return
}

func (op *Operator) EnsureIcingaNodeAlertDeleted(alert *api.NodeAlert, node *core.Node) (err error) {
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
	err = op.nodeHost.Delete(alert, node)
	return
}
