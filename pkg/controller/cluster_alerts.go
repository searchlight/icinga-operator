package controller

import (
	"errors"
	"reflect"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/log"
	tapi "github.com/appscode/searchlight/api"
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
		c.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if resource, ok := obj.(*tapi.ClusterAlert); ok {
					c.EnsureClusterAlert(nil, resource)
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldObj, ok := old.(*tapi.ClusterAlert)
				if !ok {
					log.Errorln(errors.New("Invalid ClusterAlert object"))
					return
				}
				newObj, ok := new.(*tapi.ClusterAlert)
				if !ok {
					log.Errorln(errors.New("Invalid ClusterAlert object"))
					return
				}
				if !reflect.DeepEqual(oldObj, newObj) {
					c.EnsureClusterAlert(oldObj, newObj)
				}
			},
			DeleteFunc: func(obj interface{}) {
				if resource, ok := obj.(*tapi.ClusterAlert); ok {
					c.EnsureClusterAlertDeleted(resource)
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (c *Controller) EnsureClusterAlert(old, new *tapi.ClusterAlert) {
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
				go c.EnsureLocalhost(&resource, old, new)
			}
		}
		for _, resource := range oldObjs {
			go c.EnsureLocalhostDeleted(&resource, old)
		}
	}
}

func (c *Controller) EnsureClusterAlertDeleted(alert *tapi.ClusterAlert) {
	sel, err := metav1.LabelSelectorAsSelector(&alert.Spec.Selector)
	if err != nil {
		return
	}
	opt := metav1.ListOptions{LabelSelector: sel.String()}

	if resources, err := c.KubeClient.CoreV1().Pods(alert.Namespace).List(opt); err == nil {
		for _, resource := range resources.Items {
			go c.EnsureLocalhostDeleted(&resource, alert)
		}
	}
}
