package watcher

import (
	acs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/events"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/watch"
)

func DaemonSetListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Extensions().DaemonSets(apiv1.NamespaceAll).List(opts)
	}
}

func DaemonSetWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Extensions().DaemonSets(apiv1.NamespaceAll).Watch(options)
	}
}

func ReplicaSetListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Extensions().ReplicaSets(apiv1.NamespaceAll).List(opts)
	}
}

func ReplicaSetWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Extensions().ReplicaSets(apiv1.NamespaceAll).Watch(options)
	}
}

func StatefulSetListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Apps().StatefulSets(apiv1.NamespaceAll).List(opts)
	}
}

func StatefulSetWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Apps().StatefulSets(apiv1.NamespaceAll).Watch(options)
	}
}

func AlertListFunc(c acs.ExtensionInterface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Alert(apiv1.NamespaceAll).List(opts)
	}
}

func AlertWatchFunc(c acs.ExtensionInterface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Alert(apiv1.NamespaceAll).Watch(options)
	}
}

func AlertEventListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		sets := fields.Set{
			apiv1.EventTypeField:         apiv1.EventTypeNormal,
			apiv1.EventReasonField:       events.EventReasonAlertAcknowledgement.String(),
			apiv1.EventInvolvedKindField: events.ObjectKindAlert.String(),
		}
		fieldSelector := fields.SelectorFromSet(sets)

		opts.FieldSelector = fieldSelector
		return c.Core().Events(apiv1.NamespaceAll).List(opts)
	}
}

func AlertEventWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		sets := fields.Set{
			apiv1.EventTypeField:         apiv1.EventTypeNormal,
			apiv1.EventReasonField:       events.EventReasonAlertAcknowledgement.String(),
			apiv1.EventInvolvedKindField: events.ObjectKindAlert.String(),
		}
		fieldSelector := fields.SelectorFromSet(sets)

		options.FieldSelector = fieldSelector
		return c.Core().Events(apiv1.NamespaceAll).Watch(options)
	}
}

func NamespaceListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Core().Namespaces().List(opts)
	}
}

func NamespaceWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Core().Namespaces().Watch(options)
	}
}

func PodListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Core().Pods(apiv1.NamespaceAll).List(opts)
	}
}

func PodWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Core().Pods(apiv1.NamespaceAll).Watch(options)
	}
}

func ServiceListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Core().Services(apiv1.NamespaceAll).List(opts)
	}
}

func ServiceWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Core().Services(apiv1.NamespaceAll).Watch(options)
	}
}

func ReplicationControllerWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Core().ReplicationControllers(apiv1.NamespaceAll).Watch(options)
	}
}

func ReplicationControllerListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Core().ReplicationControllers(apiv1.NamespaceAll).List(opts)
	}
}

func EndpointListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Core().Endpoints(apiv1.NamespaceAll).List(opts)
	}
}

func EndpointWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Core().Endpoints(apiv1.NamespaceAll).Watch(options)
	}
}

func NodeListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Core().Nodes().List(opts)
	}
}

func NodeWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Core().Nodes().Watch(options)
	}
}

func DeploymentListFunc(c clientset.Interface) func(apiv1.ListOptions) (runtime.Object, error) {
	return func(opts apiv1.ListOptions) (runtime.Object, error) {
		return c.Extensions().Deployments(apiv1.NamespaceAll).List(opts)
	}
}

func DeploymentWatchFunc(c clientset.Interface) func(options apiv1.ListOptions) (watch.Interface, error) {
	return func(options apiv1.ListOptions) (watch.Interface, error) {
		return c.Extensions().Deployments(apiv1.NamespaceAll).Watch(options)
	}
}
