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
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initClusterAlertWatcher() {
	op.caInformer = op.searchlightInformerFactory.Monitoring().V1alpha1().ClusterAlerts().Informer()
	op.caQueue = queue.New("ClusterAlert", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcileClusterAlert)
	op.caInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
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
				}
				queue.Enqueue(op.caQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			oldAlert, ok := oldObj.(*api.ClusterAlert)
			if !ok {
				log.Errorln("invalid ClusterAlert object")
				return
			}
			newAlert, ok := newObj.(*api.ClusterAlert)
			if !ok {
				log.Errorln("invalid ClusterAlert object")
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
				queue.Enqueue(op.caQueue.GetQueue(), newObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.caQueue.GetQueue(), obj)
		},
	})
	op.caLister = op.searchlightInformerFactory.Monitoring().V1alpha1().ClusterAlerts().Lister()
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) reconcileClusterAlert(key string) error {
	obj, exists, err := op.caInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		alert := obj.(*api.ClusterAlert)
		// Below we will warm up our cache with a ClusterAlert, so that we will see a delete for one d
		fmt.Printf("ClusterAlert %s does not exist anymore\n", key)
		op.EnsureIcingaClusterAlertDeleted(alert)
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
		op.EnsureIcingaClusterAlert(alert)
	}
	return nil
}

func (op *Operator) EnsureIcingaClusterAlert(alert *api.ClusterAlert) {
	var err error
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulSync,
				`Applied ClusterAlert: "%v"`,
				alert.Name,
			)
		} else {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`Fail to apply ClusterAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			log.Errorln(err)
		}
	}()
	err = op.clusterHost.Create(alert.DeepCopy())
	return
}

func (op *Operator) EnsureIcingaClusterAlertDeleted(alert *api.ClusterAlert) {
	var err error
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulDelete,
				`Deleted ClusterAlert: "%v"`,
				alert.Name,
			)
		} else {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				`Fail to delete ClusterAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			log.Errorln(err)
		}
	}()
	err = op.clusterHost.Delete(alert.DeepCopy())
	return
}
