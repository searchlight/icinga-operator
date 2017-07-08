package controller

import (
	"errors"
	"reflect"

	acrt "github.com/appscode/go/runtime"
	"github.com/appscode/log"
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
	_, ctrl := cache.NewInformer(
		lw,
		&apiv1.Pod{},
		c.syncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if resource, ok := obj.(*apiv1.Pod); ok {
					log.Infof("Pod %s@%s added", resource.Name, resource.Namespace)
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
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
