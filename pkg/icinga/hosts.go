package icinga

import (
	"github.com/appscode/errors"
	clientset "k8s.io/client-go/kubernetes"
)

const (
	CheckComponentStatus  = "component_status"
	CheckJsonPath         = "json_path"
	CheckNodeCount        = "node_count"
	CheckNodeStatus       = "node_status"
	CheckCommandPodStatus = "pod_status"
	CheckCommandPodExists = "pod_exists"
	CheckCommandKubeEvent = "kube_event"
	CheckCommandKubeExec  = "kube_exec"
	CheckCommandVolume    = "volume"
)

func GetObjectList(kubeClient clientset.Interface, check, hostType, namespace, objectType, objectName, specificObject string) ([]*KHost, error) {
	switch hostType {
	case HostTypePod:
		switch objectType {
		case TypeServices, TypeReplicationcontrollers, TypeDaemonsets, TypeStatefulSet, TypeReplicasets, TypeDeployments:
			if specificObject == "" {
				return GetPodList(kubeClient, namespace, objectType, objectName)
			} else {
				return GetPod(kubeClient, namespace, objectType, objectName, specificObject)
			}
		case TypePods:
			return GetPod(kubeClient, namespace, objectType, objectName, objectName)
		default:
			return nil, errors.New("Invalid kubernetes object type").Err()
		}
	case HostTypeNode:
		switch objectType {
		case TypeCluster:
			if specificObject == "" {
				return GetNodeList(kubeClient, namespace)
			} else {
				return GetNode(kubeClient, specificObject, namespace)
			}
		case TypeNodes:
			return GetNode(kubeClient, objectName, namespace)

		default:
			return nil, errors.New("Invalid object type").Err()
		}
	case HostTypeLocalhost:
		hostName := check
		if objectType != TypeCluster {
			hostName = objectType + "|" + objectName
		}
		return []*KHost{{Name: hostName + "@" + namespace, IP: "127.0.0.1"}}, nil
	default:
		return nil, errors.New("Invalid Icinga HostType").Err()
	}
}
