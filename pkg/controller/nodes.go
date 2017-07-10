package controller

import (
	"errors"
	"reflect"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/log"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (c *Controller) WatchNodes() {
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return c.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.KubeClient.CoreV1().Nodes().Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&apiv1.Node{},
		c.SyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if resource, ok := obj.(*apiv1.Node); ok {
					log.Infof("Node %s@%s added", resource.Name, resource.Namespace)

					alerts, err := util.FindNodeAlert(c.ExtClient, resource.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					if len(alerts) == 0 {
						log.Errorf("No NodeAlert found for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					for i := range alerts {
						err = c.EnsureNode(resource, nil, alerts[i])
						if err != nil {
							log.Errorf("Failed to add icinga2 alert for Node %s@%s.", resource.Name, resource.Namespace)
							// return
						}
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldNode, ok := old.(*apiv1.Node)
				if !ok {
					log.Errorln(errors.New("Invalid Node object"))
					return
				}
				newNode, ok := new.(*apiv1.Node)
				if !ok {
					log.Errorln(errors.New("Invalid Node object"))
					return
				}
				if !reflect.DeepEqual(oldNode.Labels, newNode.Labels) {
					oldAlerts, err := util.FindNodeAlert(c.ExtClient, oldNode.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", oldNode.Name, oldNode.Namespace)
						return
					}
					newAlerts, err := util.FindNodeAlert(c.ExtClient, newNode.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", newNode.Name, newNode.Namespace)
						return
					}

					type change struct {
						old *tapi.NodeAlert
						new *tapi.NodeAlert
					}
					diff := make(map[string]*change)
					for i := range oldAlerts {
						diff[oldAlerts[i].Name] = &change{old: oldAlerts[i]}
					}
					for i := range newAlerts {
						if ch, ok := diff[newAlerts[i].Name]; ok {
							ch.new = newAlerts[i]
						} else {
							diff[newAlerts[i].Name] = &change{new: newAlerts[i]}
						}
					}
					for alert := range diff {
						ch := diff[alert]
						if ch.old == nil && ch.new != nil {
							go c.EnsureNode(newNode, nil, ch.new)
						} else if ch.old != nil && ch.new == nil {
							go c.EnsureNodeDeleted(newNode, ch.old)
						} else if ch.old != nil && ch.new != nil && !reflect.DeepEqual(ch.old.Spec, ch.new.Spec) {
							go c.EnsureNode(newNode, ch.old, ch.new)
						}
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if resource, ok := obj.(*apiv1.Node); ok {
					log.Infof("Node %s@%s deleted", resource.Name, resource.Namespace)

					alerts, err := util.FindNodeAlert(c.ExtClient, resource.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					if len(alerts) == 0 {
						log.Errorf("No NodeAlert found for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					for i := range alerts {
						err = c.EnsureNodeDeleted(resource, alerts[i])
						if err != nil {
							log.Errorf("Failed to delete icinga2 alert for Node %s@%s.", resource.Name, resource.Namespace)
							// return
						}
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (c *Controller) EnsureNode(node *apiv1.Node, old, new *tapi.NodeAlert) (err error) {
	defer func() {
		if err == nil {
			c.recorder.Eventf(
				new,
				apiv1.EventTypeWarning,
				eventer.EventReasonSuccessfulSync,
				`Applied NodeAlert: "%v"`,
				new.Name,
			)
			return
		} else {
			c.recorder.Eventf(
				new,
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`Fail to be apply NodeAlert: "%v". Reason: %v`,
				new.Name,
				err,
			)
			return
		}
	}()

	if old == nil {
		err = c.nodeHost.Create(*new, *node)
	} else {
		err = c.nodeHost.Update(*new, *node)
	}
	return
}

func (c *Controller) EnsureNodeDeleted(node *apiv1.Node, alert *tapi.NodeAlert) (err error) {
	defer func() {
		if err == nil {
			c.recorder.Eventf(
				alert,
				apiv1.EventTypeWarning,
				eventer.EventReasonSuccessfulDelete,
				`Deleted NodeAlert: "%v"`,
				alert.Name,
			)
			return
		} else {
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
	}()
	err = c.nodeHost.Delete(*alert, *node)
	return
}
