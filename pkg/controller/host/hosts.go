package host

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/errors"
	"github.com/appscode/searchlight/pkg/client/icinga"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
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

func CreateIcingaHost(icingaClient *icinga.IcingaClient, objectList []*KubeObjectInfo, alertNamespace string) error {
	for _, object := range objectList {
		hostName := object.Name
		resp := icingaClient.Objects().Hosts(hostName).Get([]string{}).Do()
		if resp.Status == 200 {
			continue
		}
		var obj IcingaObject
		obj.Templates = []string{"generic-host"}
		mp := make(map[string]interface{})
		mp["address"] = object.IP

		obj.Attrs = mp
		jsonStr, err := json.Marshal(obj)
		if err != nil {
			return errors.New().WithCause(err).Err()
		}

		resp = icingaClient.Objects().Hosts(hostName).Create([]string{}, string(jsonStr)).Do()
		if resp.Err != nil {
			return errors.New().WithCause(resp.Err).Err()
		}

		if resp.Status != 200 {
			return errors.New("Can't create Icinga host").Err()
		}
	}
	return nil
}

func DeleteIcingaHost(icingaClient *icinga.IcingaClient, host string) error {
	param := map[string]string{
		"cascade": "1",
	}

	in := fmt.Sprintf(`{"filter": "match(\"%s\",host.name)"}`, host)
	var respService ResponseObject
	if _, err := icingaClient.Objects().Service("").Update([]string{}, in).Do().Into(&respService); err != nil {
		return errors.New("Can't get Icinga service").Err()
	}

	if len(respService.Results) <= 1 {
		resp := icingaClient.Objects().Hosts("").Delete([]string{}, in).Params(param).Do()
		if resp.Err != nil {
			return errors.New("Can't delete Icinga host").Err()
		}
	}
	return nil
}

func GetObjectList(kubeClient clientset.Interface, checkCommand, hostType, namespace, objectType, objectName, specificObject string) ([]*KubeObjectInfo, error) {
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
		hostName := checkCommand
		if objectType != TypeCluster {
			hostName = objectType + "|" + objectName
		}
		return []*KubeObjectInfo{{Name: hostName + "@" + namespace, IP: "127.0.0.1"}}, nil
	default:
		return nil, errors.New("Invalid Icinga HostType").Err()
	}
}
