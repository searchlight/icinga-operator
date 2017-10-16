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
func (op *Operator) WatchNodes() {
	defer runtime.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (rt.Object, error) {
			return op.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Nodes().Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&apiv1.Node{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if resource, ok := obj.(*apiv1.Node); ok {
					log.Infof("Node %s@%s added", resource.Name, resource.Namespace)

					alerts, err := util.FindNodeAlert(op.ExtClient, resource.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					if len(alerts) == 0 {
						log.Errorf("No NodeAlert found for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					for i := range alerts {
						err = op.EnsureNode(resource, nil, alerts[i])
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
					oldAlerts, err := util.FindNodeAlert(op.ExtClient, oldNode.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", oldNode.Name, oldNode.Namespace)
						return
					}
					newAlerts, err := util.FindNodeAlert(op.ExtClient, newNode.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", newNode.Name, newNode.Namespace)
						return
					}

					type change struct {
						old *api.NodeAlert
						new *api.NodeAlert
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
							go op.EnsureNode(newNode, nil, ch.new)
						} else if ch.old != nil && ch.new == nil {
							go op.EnsureNodeDeleted(newNode, ch.old)
						} else if ch.old != nil && ch.new != nil && !reflect.DeepEqual(ch.old.Spec, ch.new.Spec) {
							go op.EnsureNode(newNode, ch.old, ch.new)
						}
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if resource, ok := obj.(*apiv1.Node); ok {
					log.Infof("Node %s@%s deleted", resource.Name, resource.Namespace)

					alerts, err := util.FindNodeAlert(op.ExtClient, resource.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching NodeAlert for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					if len(alerts) == 0 {
						log.Errorf("No NodeAlert found for Node %s@%s.", resource.Name, resource.Namespace)
						return
					}
					for i := range alerts {
						err = op.EnsureNodeDeleted(resource, alerts[i])
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

func (op *Operator) EnsureNode(node *apiv1.Node, old, new *api.NodeAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				new.ObjectReference(),
				apiv1.EventTypeNormal,
				eventer.EventReasonSuccessfulSync,
				`Applied NodeAlert: "%v"`,
				new.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				new.ObjectReference(),
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`Fail to be apply NodeAlert: "%v". Reason: %v`,
				new.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()

	if old == nil {
		err = op.nodeHost.Create(*new, *node)
	} else {
		err = op.nodeHost.Update(*new, *node)
	}
	return
}

func (op *Operator) EnsureNodeDeleted(node *apiv1.Node, alert *api.NodeAlert) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				apiv1.EventTypeNormal,
				eventer.EventReasonSuccessfulDelete,
				`Deleted NodeAlert: "%v"`,
				alert.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				alert.ObjectReference(),
				apiv1.EventTypeWarning,
				eventer.EventReasonFailedToDelete,
				`Fail to be delete NodeAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()
	err = op.nodeHost.Delete(*alert, *node)
	return
}
