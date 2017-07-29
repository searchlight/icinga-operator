package operator

import (
	"errors"
	"reflect"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/log"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/eventer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (op *Operator) WatchClusterAlerts() {
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.ExtClient.ClusterAlerts(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.ExtClient.ClusterAlerts(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.ClusterAlert{},
		op.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.ClusterAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						op.recorder.Eventf(
							alert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToCreate,
							`Fail to be create ClusterAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					op.EnsureClusterAlert(nil, alert)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldAlert, ok := old.(*tapi.ClusterAlert)
				if !ok {
					log.Errorln(errors.New("Invalid ClusterAlert object"))
					return
				}
				newAlert, ok := new.(*tapi.ClusterAlert)
				if !ok {
					log.Errorln(errors.New("Invalid ClusterAlert object"))
					return
				}
				if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
					if ok, err := newAlert.IsValid(); !ok {
						op.recorder.Eventf(
							newAlert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be update ClusterAlert: "%v". Reason: %v`,
							newAlert.Name,
							err,
						)
						return
					}
					op.EnsureClusterAlert(oldAlert, newAlert)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.ClusterAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						op.recorder.Eventf(
							alert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be delete ClusterAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					op.EnsureClusterAlertDeleted(alert)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (op *Operator) EnsureClusterAlert(old, new *tapi.ClusterAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				new,
				apiv1.EventTypeNormal,
				eventer.EventReasonSuccessfulSync,
				`Applied ClusterAlert: "%v"`,
				new.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				new,
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

func (op *Operator) EnsureClusterAlertDeleted(alert *tapi.ClusterAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert,
				apiv1.EventTypeNormal,
				eventer.EventReasonSuccessfulDelete,
				`Deleted ClusterAlert: "%v"`,
				alert.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				alert,
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
