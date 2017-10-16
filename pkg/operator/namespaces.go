package operator

import (
	acrt "github.com/appscode/go/runtime"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// Blocks caller. Intended to be called as a Go routine.
// ref: https://github.com/kubernetes/kubernetes/issues/46736
func (op *Operator) WatchNamespaces() {
	defer acrt.HandleCrash()

	lw := &cache.ListWatch{
		ListFunc: func(opts metav1.ListOptions) (runtime.Object, error) {
			return op.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Namespaces().Watch(metav1.ListOptions{})
		},
	}
	_, ctrl := cache.NewInformer(lw,
		&apiv1.Namespace{},
		op.Opt.ResyncPeriod,
		cache.ResourceEventHandlerFuncs{
			DeleteFunc: func(obj interface{}) {
				if ns, ok := obj.(*apiv1.Namespace); ok {
					if alerts, err := op.ExtClient.ClusterAlerts(ns.Name).List(metav1.ListOptions{}); err == nil {
						for _, alert := range alerts.Items {
							op.ExtClient.ClusterAlerts(alert.Namespace).Delete(alert.Name, &metav1.DeleteOptions{})
						}
					}
					if alerts, err := op.ExtClient.NodeAlerts(ns.Name).List(metav1.ListOptions{}); err == nil {
						for _, alert := range alerts.Items {
							op.ExtClient.NodeAlerts(alert.Namespace).Delete(alert.Name, &metav1.DeleteOptions{})
						}
					}
					if alerts, err := op.ExtClient.PodAlerts(ns.Name).List(metav1.ListOptions{}); err == nil {
						for _, alert := range alerts.Items {
							op.ExtClient.PodAlerts(alert.Namespace).Delete(alert.Name, &metav1.DeleteOptions{})
						}
					}
				}
			},
		},
	)
	ctrl.Run(wait.NeverStop)
}
