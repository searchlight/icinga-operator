package operator

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initNamespaceWatcher() {
	op.nsInformer = op.kubeInformerFactory.Core().V1().Namespaces().Informer()
	op.nsInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			if ns, ok := obj.(*core.Namespace); ok {
				op.ExtClient.MonitoringV1alpha1().ClusterAlerts(ns.Name).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{})
				op.ExtClient.MonitoringV1alpha1().NodeAlerts(ns.Name).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{})
				op.ExtClient.MonitoringV1alpha1().PodAlerts(ns.Name).DeleteCollection(&metav1.DeleteOptions{}, metav1.ListOptions{})
			}
		},
	})
	op.nsLister = op.kubeInformerFactory.Core().V1().Namespaces().Lister()
}
