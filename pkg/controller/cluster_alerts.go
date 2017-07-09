package controller

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
func (c *Controller) WatchClusterAlerts() {
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return c.ExtClient.ClusterAlerts(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.ExtClient.ClusterAlerts(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.ClusterAlert{},
		c.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.ClusterAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						c.recorder.Eventf(
							alert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToCreate,
							`Fail to be create ClusterAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					c.EnsureClusterAlert(nil, alert)
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
						c.recorder.Eventf(
							newAlert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be update ClusterAlert: "%v". Reason: %v`,
							newAlert.Name,
							err,
						)
						return
					}
					c.EnsureClusterAlert(oldAlert, newAlert)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.ClusterAlert); ok {
					if ok, err := alert.IsValid(); !ok {
						c.recorder.Eventf(
							alert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be delete ClusterAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					c.EnsureClusterAlertDeleted(alert)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (c *Controller) EnsureClusterAlert(old, new *tapi.ClusterAlert) (err error) {
	defer func() {
		if err == nil {
			c.recorder.Eventf(
				new,
				apiv1.EventTypeWarning,
				eventer.EventReasonSuccessfulSync,
				`Applied ClusterAlert: "%v"`,
				new.Name,
			)
			return
		} else {
			c.recorder.Eventf(
				new,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`Fail to be apply ClusterAlert: "%v". Reason: %v`,
				new.Name,
				err,
			)
			return
		}
	}()

	if old == nil {
		err = c.clusterHost.Create(*new)
	} else {
		err = c.clusterHost.Update(*new)
	}
	return
}

func (c *Controller) EnsureClusterAlertDeleted(alert *tapi.ClusterAlert) (err error) {
	defer func() {
		if err == nil {
			c.recorder.Eventf(
				alert,
				apiv1.EventTypeWarning,
				eventer.EventReasonSuccessfulDelete,
				`Deleted ClusterAlert: "%v"`,
				alert.Name,
			)
			return
		} else {
			c.recorder.Eventf(
				alert,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				`Fail to be delete ClusterAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			return
		}
	}()
	err = c.clusterHost.Delete(*alert)
	return
}
