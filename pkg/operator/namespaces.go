package operator

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rt "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initNamespaceWatcher() {
	lw := &cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (rt.Object, error) {
			return op.KubeClient.CoreV1().Namespaces().List(options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return op.KubeClient.CoreV1().Namespaces().Watch(options)
		},
	}

	op.nsIndexer, op.nsInformer = cache.NewIndexerInformer(lw, &core.Namespace{}, op.Opt.ResyncPeriod, cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			if ns, ok := obj.(*core.Namespace); ok {
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
	}, cache.Indexers{})
}
