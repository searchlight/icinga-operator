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
				if alerts, err := op.ExtClient.MonitoringV1alpha1().ClusterAlerts(ns.Name).List(metav1.ListOptions{}); err == nil {
					for _, alert := range alerts.Items {
						op.ExtClient.MonitoringV1alpha1().ClusterAlerts(alert.Namespace).Delete(alert.Name, &metav1.DeleteOptions{})
					}
				}
				if alerts, err := op.ExtClient.MonitoringV1alpha1().NodeAlerts(ns.Name).List(metav1.ListOptions{}); err == nil {
					for _, alert := range alerts.Items {
						op.ExtClient.MonitoringV1alpha1().NodeAlerts(alert.Namespace).Delete(alert.Name, &metav1.DeleteOptions{})
					}
				}
				if alerts, err := op.ExtClient.MonitoringV1alpha1().PodAlerts(ns.Name).List(metav1.ListOptions{}); err == nil {
					for _, alert := range alerts.Items {
						op.ExtClient.MonitoringV1alpha1().PodAlerts(alert.Namespace).Delete(alert.Name, &metav1.DeleteOptions{})
					}
				}
			}
		},
	})
}
