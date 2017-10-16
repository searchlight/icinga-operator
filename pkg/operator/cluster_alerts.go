package operator

import (
	"errors"
	"reflect"

	"github.com/appscode/go/log"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/util"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rt "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchClusterAlerts() {
	defer runtime.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (rt.Object, error) {
			return op.ExtClient.ClusterAlerts(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.ExtClient.ClusterAlerts(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&api.ClusterAlert{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*api.ClusterAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						op.recorder.Eventf(
							alert.ObjectReference(),
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToCreate,
							`Fail to be create ClusterAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
						op.recorder.Eventf(
							alert.ObjectReference(),
							apiv1.EventTypeWarning,
							eventer.EventReasonBadNotifier,
							`Bad notifier config for ClusterAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
					}
					op.EnsureClusterAlert(nil, alert)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldAlert, ok := old.(*api.ClusterAlert)
				if !ok {
					log.Errorln(errors.New("Invalid ClusterAlert object"))
					return
				}
				newAlert, ok := new.(*api.ClusterAlert)
				if !ok {
					log.Errorln(errors.New("Invalid ClusterAlert object"))
					return
				}
				if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
					if ok, err := newAlert.IsValid(); !ok {
						op.recorder.Eventf(
							newAlert.ObjectReference(),
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be update ClusterAlert: "%v". Reason: %v`,
							newAlert.Name,
							err,
						)
						return
					}
					if err := util.CheckNotifiers(op.KubeClient, newAlert); err != nil {
						op.recorder.Eventf(
							newAlert.ObjectReference(),
							apiv1.EventTypeWarning,
							eventer.EventReasonBadNotifier,
							`Bad notifier config for ClusterAlert: "%v". Reason: %v`,
							newAlert.Name,
							err,
						)
					}
					op.EnsureClusterAlert(oldAlert, newAlert)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if alert, ok := obj.(*api.ClusterAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						op.recorder.Eventf(
							alert.ObjectReference(),
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be delete ClusterAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
						op.recorder.Eventf(
							alert.ObjectReference(),
							apiv1.EventTypeWarning,
							eventer.EventReasonBadNotifier,
							`Bad notifier config for ClusterAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
					}
					op.EnsureClusterAlertDeleted(alert)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (op *Operator) EnsureClusterAlert(old, new *api.ClusterAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				new.ObjectReference(),
				apiv1.EventTypeNormal,
				eventer.EventReasonSuccessfulSync,
				`Applied ClusterAlert: "%v"`,
				new.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				new.ObjectReference(),
				apiv1.EventTypeWarning,
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

func (op *Operator) EnsureClusterAlertDeleted(alert *api.ClusterAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				apiv1.EventTypeNormal,
				eventer.EventReasonSuccessfulDelete,
				`Deleted ClusterAlert: "%v"`,
				alert.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				alert.ObjectReference(),
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				`Fail to be delete ClusterAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()
	err = op.clusterHost.Delete(*alert)
	return
}
