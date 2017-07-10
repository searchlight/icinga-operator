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
		c.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.NodeAlert); ok {
					if ok, err := alert.IsValid(); !ok {
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
					if ok, err := newAlert.IsValid(); !ok {
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
					if ok, err := alert.IsValid(); !ok {
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
	oldObjs := make(map[string]*apiv1.Node)

	if old != nil {
		oldSel := labels.SelectorFromSet(old.Spec.Selector)
		if old.Spec.NodeName != "" {
			if resource, err := c.KubeClient.CoreV1().Nodes().Get(old.Spec.NodeName, metav1.GetOptions{}); err == nil {
				if oldSel.Matches(labels.Set(resource.Labels)) {
					oldObjs[resource.Name] = resource
				}
			}
		} else {
			if resources, err := c.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: oldSel.String()}); err == nil {
				for _, resource := range resources.Items {
					oldObjs[resource.Name] = &resource
				}
			}
		}
	}

	newSel := labels.SelectorFromSet(new.Spec.Selector)
	if new.Spec.NodeName != "" {
		if resource, err := c.KubeClient.CoreV1().Nodes().Get(new.Spec.NodeName, metav1.GetOptions{}); err == nil {
			if newSel.Matches(labels.Set(resource.Labels)) {
				delete(oldObjs, resource.Name)
				go c.EnsureNode(resource, old, new)
			}
		}
	} else {
		if resources, err := c.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: newSel.String()}); err == nil {
			for _, resource := range resources.Items {
				delete(oldObjs, resource.Name)
				go c.EnsureNode(&resource, old, new)
			}
		}
	}
	for i := range oldObjs {
		go c.EnsureNodeDeleted(oldObjs[i], old)
	}
}

func (c *Controller) EnsureNodeAlertDeleted(alert *tapi.NodeAlert) {
	sel := labels.SelectorFromSet(alert.Spec.Selector)
	if alert.Spec.NodeName != "" {
		if resource, err := c.KubeClient.CoreV1().Nodes().Get(alert.Spec.NodeName, metav1.GetOptions{}); err == nil {
			if sel.Matches(labels.Set(resource.Labels)) {
				go c.EnsureNodeDeleted(resource, alert)
			}
		}
	} else {
		if resources, err := c.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: sel.String()}); err == nil {
			for _, resource := range resources.Items {
				go c.EnsureNodeDeleted(&resource, alert)
			}
		}
	}
}
