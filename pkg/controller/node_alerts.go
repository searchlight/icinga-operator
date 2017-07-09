package controller

import (
	"errors"
	"reflect"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/log"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/eventer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (c *Controller) WatchNodeAlerts() {
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return c.ExtClient.NodeAlerts(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.ExtClient.NodeAlerts(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.NodeAlert{},
		c.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.NodeAlert); ok {
					if ok, err := alert.Spec.IsValid(); !ok {
						c.recorder.Eventf(
							alert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToCreate,
							`Fail to be create NodeAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					c.EnsureNodeAlert(nil, alert)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldAlert, ok := old.(*tapi.NodeAlert)
				if !ok {
					log.Errorln(errors.New("Invalid NodeAlert object"))
					return
				}
				newAlert, ok := new.(*tapi.NodeAlert)
				if !ok {
					log.Errorln(errors.New("Invalid NodeAlert object"))
					return
				}
				if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
					if ok, err := newAlert.Spec.IsValid(); !ok {
						c.recorder.Eventf(
							newAlert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be update NodeAlert: "%v". Reason: %v`,
							newAlert.Name,
							err,
						)
						return
					}
					c.EnsureNodeAlert(oldAlert, newAlert)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.NodeAlert); ok {
					if ok, err := alert.Spec.IsValid(); !ok {
						c.recorder.Eventf(
							alert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be delete NodeAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					c.EnsureNodeAlertDeleted(alert)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (c *Controller) EnsureNodeAlert(old, new *tapi.NodeAlert) {
	var oldOpt, newOpt *metav1.ListOptions
	if old != nil {
		oldOpt = &metav1.ListOptions{LabelSelector: labels.SelectorFromSet(old.Spec.Selector).String()}
	}

	newOpt = &metav1.ListOptions{LabelSelector: labels.SelectorFromSet(new.Spec.Selector).String()}

	{
		oldObjs := make(map[string]apiv1.Node)
		if oldOpt != nil {
			if resources, err := c.KubeClient.CoreV1().Nodes().List(*oldOpt); err == nil {
				for _, resource := range resources.Items {
					oldObjs[resource.Name] = resource
				}
			}
		}

		if resources, err := c.KubeClient.CoreV1().Nodes().List(*newOpt); err == nil {
			for _, resource := range resources.Items {
				delete(oldObjs, resource.Name)
				go c.EnsureNode(&resource, old, new)
			}
		}
		for _, resource := range oldObjs {
			go c.EnsureNodeDeleted(&resource, old)
		}
	}
}

func (c *Controller) EnsureNodeAlertDeleted(alert *tapi.NodeAlert) {
	opt := metav1.ListOptions{LabelSelector: labels.SelectorFromSet(alert.Spec.Selector).String()}

	if resources, err := c.KubeClient.CoreV1().Nodes().List(opt); err == nil {
		for _, resource := range resources.Items {
			go c.EnsureNodeDeleted(&resource, alert)
		}
	}
}
