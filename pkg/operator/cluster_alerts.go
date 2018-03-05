package operator

import (
	"fmt"
	"reflect"

	"github.com/appscode/go/log"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	mon_api "github.com/appscode/searchlight/apis/monitoring"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	mon_util "github.com/appscode/searchlight/client/clientset/versioned/typed/monitoring/v1alpha1/util"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initClusterAlertWatcher() {
	op.caInformer = op.searchlightInformerFactory.Monitoring().V1alpha1().ClusterAlerts().Informer()
	op.caQueue = queue.New("ClusterAlert", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcileClusterAlert)
	op.caInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if alert, ok := obj.(*api.ClusterAlert); ok {
				if !op.validateClusterAlert(alert) {
					log.Errorf(`Invalid ClusterAlert "%s@%s"`, alert.Name, alert.Namespace)
					return
				}
				queue.Enqueue(op.caQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			oldAlert, ok := oldObj.(*api.ClusterAlert)
			if !ok {
				return
			}
			newAlert, ok := newObj.(*api.ClusterAlert)
			if !ok {
				return
			}
			// DeepEqual old & new
			// DeepEqual MapperConfiguration of old & new
			// Patch PodAlert with necessary annotation
			newAlert, err := op.processClusterAlertUpdate(oldAlert, newAlert)
			if err != nil {
				log.Error(err)
			} else {
				queue.Enqueue(op.caQueue.GetQueue(), newAlert)
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
		log.Debugf("ClusterAlert %s does not exist anymore\n", key)
	} else {
		alert := obj.(*api.ClusterAlert)
		if alert.DeletionTimestamp != nil {
			if core_util.HasFinalizer(alert.ObjectMeta, mon_api.GroupName) {
				// Delete all Icinga objects created for this ClusterAlert
				if err := op.EnsureIcingaClusterAlertDeleted(alert); err != nil {
					log.Errorf("Failed to delete ClusterAlert %s@%s", alert.Name, alert.Namespace)
					return err
				}

				_, _, err = mon_util.PatchClusterAlert(op.ExtClient.MonitoringV1alpha1(), alert, func(in *api.ClusterAlert) *api.ClusterAlert {
					in.ObjectMeta = core_util.RemoveFinalizer(in.ObjectMeta, mon_api.GroupName)
					return in
				})
				return err
			}
		} else {
			fmt.Printf("Sync/Add/Update for ClusterAlert %s\n", alert.GetName())

			alert, _, err = mon_util.PatchClusterAlert(op.ExtClient.MonitoringV1alpha1(), alert, func(in *api.ClusterAlert) *api.ClusterAlert {
				in.ObjectMeta = core_util.AddFinalizer(in.ObjectMeta, mon_api.GroupName)
				return in
			})

			if err := op.EnsureIcingaClusterAlert(alert); err != nil {
				log.Errorf("Failed to patch ClusterAlert %s@%s", alert.Name, alert.Namespace)
			}
		}
	}
	return nil
}

func (op *Operator) EnsureIcingaClusterAlert(alert *api.ClusterAlert) (err error) {
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
	return err
}

func (op *Operator) EnsureIcingaClusterAlertDeleted(alert *api.ClusterAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulDelete,
				`Deleted Icinga objects of ClusterAlert "%s@%s"`,
				alert.Name, alert.Namespace,
			)
		} else {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				`Fail to delete Icinga objects of ClusterAlert "%s@%s. Reason: %v`,
				alert.Name, alert.Namespace,
				err,
			)
			log.Errorln(err)
		}
	}()
	err = op.clusterHost.Delete(alert.DeepCopy())
	return err
}

func (op *Operator) processClusterAlertUpdate(oldAlert, newAlert *api.ClusterAlert) (*api.ClusterAlert, error) {
	// Check for changes in Spec
	if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
		if !op.validateClusterAlert(newAlert) {
			return nil, errors.Errorf(`Invalid ClusterAlert "%s@%s"`, newAlert.Name, newAlert.Namespace)
		}
	}

	return newAlert, nil
}

func (op *Operator) validateClusterAlert(alert *api.ClusterAlert) bool {
	// Validate IcingaCommand & it's variables.
	// And also check supported IcingaState
	if ok, err := alert.IsValid(); !ok {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToCreate,
			`Reason: %v`,
			err,
		)
		return false
	}

	// Validate Notifiers configurations
	if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonBadNotifier,
			`Bad notifier config for ClusterAlert: "%s@%s". Reason: %v`,
			alert.Name, alert.Namespace,
			err,
		)
		return false
	}

	return true
}
