package util

import (
	"errors"
	"fmt"
	"os"

	"github.com/appscode/searchlight/pkg/controller/host"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientset "k8s.io/client-go/kubernetes"
)

func GetLabels(client clientset.Interface, namespace, objectType, objectName string) (labels.Selector, error) {
	var labelMap map[string]string
	switch objectType {
	case host.TypeServices:
		service, err := client.CoreV1().Services(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		labelMap = service.Spec.Selector
	case host.TypeReplicationcontrollers:
		rc, err := client.CoreV1().ReplicationControllers(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		labelMap = rc.Spec.Selector
	case host.TypeDaemonsets:
		daemonSet, err := client.ExtensionsV1beta1().DaemonSets(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		labelMap = daemonSet.Spec.Selector.MatchLabels
	case host.TypeReplicasets:
		replicaSet, err := client.ExtensionsV1beta1().ReplicaSets(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		labelMap = replicaSet.Spec.Selector.MatchLabels
	case host.TypeStatefulSet:
		stateFulSet, err := client.AppsV1beta1().StatefulSets(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		labelMap = stateFulSet.Spec.Selector.MatchLabels
	case host.TypeDeployments:
		deployment, err := client.ExtensionsV1beta1().Deployments(namespace).Get(objectName, metav1.GetOptions{})
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
