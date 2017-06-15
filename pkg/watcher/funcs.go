package watcher

import (
	acs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/events"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func DaemonSetListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.ExtensionsV1beta1().DaemonSets(apiv1.NamespaceAll).List(opts)
	}
}

func DaemonSetWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.ExtensionsV1beta1().DaemonSets(apiv1.NamespaceAll).Watch(options)
	}
}

func ReplicaSetListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.ExtensionsV1beta1().ReplicaSets(apiv1.NamespaceAll).List(opts)
	}
}

func ReplicaSetWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.ExtensionsV1beta1().ReplicaSets(apiv1.NamespaceAll).Watch(options)
	}
}

func StatefulSetListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.AppsV1beta1().StatefulSets(apiv1.NamespaceAll).List(opts)
	}
}

func StatefulSetWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.AppsV1beta1().StatefulSets(apiv1.NamespaceAll).Watch(options)
	}
}

func AlertListFunc(c acs.ExtensionInterface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.Alert(apiv1.NamespaceAll).List(opts)
	}
}

func AlertWatchFunc(c acs.ExtensionInterface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.Alert(apiv1.NamespaceAll).Watch(options)
	}
}

func AlertEventListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		sets := fields.Set{
			api.EventTypeField:         apiv1.EventTypeNormal,
			api.EventReasonField:       events.EventReasonAlertAcknowledgement.String(),
			api.EventInvolvedKindField: events.ObjectKindAlert.String(),
		}
		fieldSelector := fields.SelectorFromSet(sets)
		opts.FieldSelector = fieldSelector.String()
		return c.CoreV1().Events(apiv1.NamespaceAll).List(opts)
	}
}

func AlertEventWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		sets := fields.Set{
			api.EventTypeField:         apiv1.EventTypeNormal,
			api.EventReasonField:       events.EventReasonAlertAcknowledgement.String(),
			api.EventInvolvedKindField: events.ObjectKindAlert.String(),
		}
		fieldSelector := fields.SelectorFromSet(sets)
		options.FieldSelector = fieldSelector.String()
		return c.CoreV1().Events(apiv1.NamespaceAll).Watch(options)
	}
}

func NamespaceListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.CoreV1().Namespaces().List(opts)
	}
}

func NamespaceWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.CoreV1().Namespaces().Watch(options)
	}
}

func PodListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.CoreV1().Pods(apiv1.NamespaceAll).List(opts)
	}
}

func PodWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.CoreV1().Pods(apiv1.NamespaceAll).Watch(options)
	}
}

func ServiceListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.CoreV1().Services(apiv1.NamespaceAll).List(opts)
	}
}

func ServiceWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.CoreV1().Services(apiv1.NamespaceAll).Watch(options)
	}
}

func ReplicationControllerWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.CoreV1().ReplicationControllers(apiv1.NamespaceAll).Watch(options)
	}
}

func ReplicationControllerListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.CoreV1().ReplicationControllers(apiv1.NamespaceAll).List(opts)
	}
}

func EndpointListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.CoreV1().Endpoints(apiv1.NamespaceAll).List(opts)
	}
}

func EndpointWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.CoreV1().Endpoints(apiv1.NamespaceAll).Watch(options)
	}
}

func NodeListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.CoreV1().Nodes().List(opts)
	}
}

func NodeWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.CoreV1().Nodes().Watch(options)
	}
}

func DeploymentListFunc(c clientset.Interface) func(metav1.ListOptions) (runtime.Object, error) {
	return func(opts metav1.ListOptions) (runtime.Object, error) {
		return c.Extensions().Deployments(apiv1.NamespaceAll).List(opts)
	}
}

func DeploymentWatchFunc(c clientset.Interface) func(options metav1.ListOptions) (watch.Interface, error) {
	return func(options metav1.ListOptions) (watch.Interface, error) {
		return c.Extensions().Deployments(apiv1.NamespaceAll).Watch(options)
	}
}
