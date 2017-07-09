package controller

import (
	"errors"
	"reflect"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/log"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
func (c *Controller) WatchPods() {
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return c.KubeClient.CoreV1().Pods(apiv1.NamespaceAll).List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return c.KubeClient.CoreV1().Pods(apiv1.NamespaceAll).Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&apiv1.Pod{},
		c.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if resource, ok := obj.(*apiv1.Pod); ok {
					log.Infof("Pod %s@%s added", resource.Name, resource.Namespace)

					alerts, err := util.FindPodAlert(c.ExtClient, resource.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching PodAlert for Pod %s@%s.", resource.Name, resource.Namespace)
						return
					}
					if len(alerts) == 0 {
						log.Errorf("No PodAlert found for Pod %s@%s.", resource.Name, resource.Namespace)
						return
					}
					for _, alert := range alerts {
						err = c.EnsurePod(resource, nil, alert)
						if err != nil {
							log.Errorf("Failed to add icinga2 alert for Pod %s@%s.", resource.Name, resource.Namespace)
							// return
						}
					}
				}
			},
			UpdateFunc: func(old, new interface{}) {
				oldObj, ok := old.(*apiv1.Pod)
				if !ok {
					log.Errorln(errors.New("Invalid Pod object"))
					return
				}
				newObj, ok := new.(*apiv1.Pod)
				if !ok {
					log.Errorln(errors.New("Invalid Pod object"))
					return
				}
				if !reflect.DeepEqual(oldObj.Labels, newObj.Labels) {
					oldAlerts, err := util.FindPodAlert(c.ExtClient, oldObj.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching PodAlert for Pod %s@%s.", oldObj.Name, oldObj.Namespace)
						return
					}
					newAlerts, err := util.FindPodAlert(c.ExtClient, newObj.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching PodAlert for Pod %s@%s.", newObj.Name, newObj.Namespace)
						return
					}

					type change struct {
						old *tapi.PodAlert
						new *tapi.PodAlert
					}
					diff := make(map[string]*change)
					for _, alert := range oldAlerts {
						diff[alert.Name] = &change{old: alert}
					}
					for _, alert := range newAlerts {
						if ch, ok := diff[alert.Name]; ok {
							ch.new = alert
						} else {
							diff[alert.Name] = &change{new: alert}
						}
					}
					for _, ch := range diff {
						if ch.old == nil && ch.new != nil {
							c.EnsurePod(newObj, nil, ch.new)
						} else if ch.old != nil && ch.new == nil {
							c.EnsurePodDeleted(newObj, ch.old)
						} else if ch.old != nil && ch.new != nil && !reflect.DeepEqual(ch.old.Spec, ch.new.Spec) {
							c.EnsurePod(newObj, ch.old, ch.new)
						}
					}
				}
			},
			DeleteFunc: func(obj interface{}) {
				if resource, ok := obj.(*apiv1.Pod); ok {
					log.Infof("Pod %s@%s deleted", resource.Name, resource.Namespace)

					alerts, err := util.FindPodAlert(c.ExtClient, resource.ObjectMeta)
					if err != nil {
						log.Errorf("Error while searching PodAlert for Pod %s@%s.", resource.Name, resource.Namespace)
						return
					}
					if len(alerts) == 0 {
						log.Errorf("No PodAlert found for Pod %s@%s.", resource.Name, resource.Namespace)
						return
					}
					for _, alert := range alerts {
						err = c.EnsurePodDeleted(resource, alert)
						if err != nil {
							log.Errorf("Failed to delete icinga2 alert for Pod %s@%s.", resource.Name, resource.Namespace)
							// return
						}
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}

func (c *Controller) EnsurePod(resource *apiv1.Pod, old, new *tapi.PodAlert) (err error) {
	return nil
}

func (c *Controller) EnsurePodDeleted(resource *apiv1.Pod, alert *tapi.PodAlert) (err error) {
	return nil
}
