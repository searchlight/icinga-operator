package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/appscode/searchlight/pkg/controller/host"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/labels"
)

func GetLabels(client clientset.Interface, namespace, objectType, objectName string) (labels.Selector, error) {
	var labelMap map[string]string
	switch objectType {
	case host.TypeServices:
		service, err := client.Core().Services(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = service.Spec.Selector
	case host.TypeReplicationcontrollers:
		rc, err := client.Core().ReplicationControllers(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = rc.Spec.Selector
	case host.TypeDaemonsets:
		daemonSet, err := client.Extensions().DaemonSets(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = daemonSet.Spec.Selector.MatchLabels
	case host.TypeReplicasets:
		replicaSet, err := client.Extensions().ReplicaSets(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = replicaSet.Spec.Selector.MatchLabels
	case host.TypeStatefulSet:
		stateFulSet, err := client.Apps().StatefulSets(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = stateFulSet.Spec.Selector.MatchLabels
	case host.TypeDeployments:
		deployment, err := client.Extensions().Deployments(namespace).Get(objectName)
		if err != nil {
			return nil, err
		}
		labelMap = deployment.Spec.Selector.MatchLabels
	default:
		return nil, errors.New("Invalid kubernetes object type")
	}
	return labels.SelectorFromSet(labelMap), nil
}

func Output(icingaState IcingaState, message interface{}) {
	fmt.Fprintln(os.Stdout, State[int(icingaState)], message)
	os.Exit(int(icingaState))
}
