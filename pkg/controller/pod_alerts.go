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
func (c *Controller) WatchPodAlerts() {
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return c.ExtClient.PodAlerts(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.ExtClient.PodAlerts(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&tapi.PodAlert{},
		c.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.PodAlert); ok {
					if ok, err := alert.Spec.IsValid(); !ok {
						c.recorder.Eventf(
							alert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToCreate,
							`Fail to be create PodAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					c.EnsurePodAlert(nil, alert)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldAlert, ok := old.(*tapi.PodAlert)
				if !ok {
					log.Errorln(errors.New("Invalid PodAlert object"))
					return
				}
				newAlert, ok := new.(*tapi.PodAlert)
				if !ok {
					log.Errorln(errors.New("Invalid PodAlert object"))
					return
				}
				if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
					if ok, err := newAlert.Spec.IsValid(); !ok {
						c.recorder.Eventf(
							newAlert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be update PodAlert: "%v". Reason: %v`,
							newAlert.Name,
							err,
						)
						return
					}
					c.EnsurePodAlert(oldAlert, newAlert)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if alert, ok := obj.(*tapi.PodAlert); ok {
					if ok, err := alert.Spec.IsValid(); !ok {
						c.recorder.Eventf(
							alert,
							apiv1.EventTypeWarning,
							eventer.EventReasonFailedToDelete,
							`Fail to be delete PodAlert: "%v". Reason: %v`,
							alert.Name,
							err,
						)
						return
					}
					c.EnsurePodAlertDeleted(alert)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (c *Controller) EnsurePodAlert(old, new *tapi.PodAlert) {
	var oldOpt, newOpt *metav1.ListOptions
	if old != nil {
		oldSelector, err := metav1.LabelSelectorAsSelector(&old.Spec.Selector)
		if err != nil {
			return
		}
		oldOpt = &metav1.ListOptions{LabelSelector: oldSelector.String()}
	}

	newSelector, err := metav1.LabelSelectorAsSelector(&new.Spec.Selector)
	if err != nil {
		return
	}
	newOpt = &metav1.ListOptions{LabelSelector: newSelector.String()}

	{
		oldObjs := make(map[string]apiv1.Pod)
		if oldOpt != nil {
			if resources, err := c.KubeClient.CoreV1().Pods(new.Namespace).List(*oldOpt); err == nil {
				for _, resource := range resources.Items {
					oldObjs[resource.Name] = resource
				}
			}
		}

		if resources, err := c.KubeClient.CoreV1().Pods(new.Namespace).List(*newOpt); err == nil {
			for _, resource := range resources.Items {
				delete(oldObjs, resource.Name)
				go c.EnsurePod(&resource, old, new)
			}
		}
		for _, resource := range oldObjs {
			go c.EnsurePodDeleted(&resource, old)
		}
	}
}

func (c *Controller) EnsurePodAlertDeleted(alert *tapi.PodAlert) {
	sel, err := metav1.LabelSelectorAsSelector(&alert.Spec.Selector)
	if err != nil {
		return
	}
	opt := metav1.ListOptions{LabelSelector: sel.String()}

	if resources, err := c.KubeClient.CoreV1().Pods(alert.Namespace).List(opt); err == nil {
		for _, resource := range resources.Items {
			go c.EnsurePodDeleted(&resource, alert)
		}
	}
}
