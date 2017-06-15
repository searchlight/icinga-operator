package host

import (
	"os"

	"github.com/appscode/errors"
	aci "github.com/appscode/searchlight/api"
	acs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/controller/types"
	"github.com/appscode/searchlight/pkg/events"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

const (
	TypeServices               = "services"
	TypeReplicationcontrollers = "replicationcontrollers"
	TypeDaemonsets             = "daemonsets"
	TypeStatefulSet            = "statefulsets"
	TypeReplicasets            = "replicasets"
	TypeDeployments            = "deployments"
	TypePods                   = "pods"
	TypeNodes                  = "nodes"
	TypeCluster                = "cluster"
)

func getLabels(client clientset.Interface, namespace, objectType, objectName string) (labels.Selector, error) {
	label := labels.NewSelector()
	labelsMap := make(map[string]string, 0)
	if objectType == TypeServices {
		service, err := client.CoreV1().Services(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, errors.New().WithCause(err).Err()
		}
		labelsMap = service.Spec.Selector

	} else if objectType == TypeReplicationcontrollers {
		rc, err := client.CoreV1().ReplicationControllers(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, errors.New().WithCause(err).Err()
		}
		labelsMap = rc.Spec.Selector
	} else if objectType == TypeDaemonsets {
		daemonSet, err := client.ExtensionsV1beta1().DaemonSets(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, errors.New().WithCause(err).Err()
		}
		labelsMap = daemonSet.Spec.Selector.MatchLabels
	} else if objectType == TypeReplicasets {
		replicaSet, err := client.ExtensionsV1beta1().ReplicaSets(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, errors.New().WithCause(err).Err()
		}
		labelsMap = replicaSet.Spec.Selector.MatchLabels
	} else if objectType == TypeStatefulSet {
		statefulSet, err := client.Apps().StatefulSets(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, errors.New().WithCause(err).Err()
		}
		labelsMap = statefulSet.Spec.Selector.MatchLabels
	} else if objectType == TypeDeployments {
		deployment, err := client.ExtensionsV1beta1().Deployments(namespace).Get(objectName, metav1.GetOptions{})
		if err != nil {
			return nil, errors.New().WithCause(err).Err()
		}
		labelsMap = deployment.Spec.Selector.MatchLabels
	} else {
		return label, errors.New("Invalid kubernetes object type").Err()
	}

	for key, value := range labelsMap {
		s := sets.NewString(value)
		ls, err := labels.NewRequirement(key, selection.Equals, s.List())
		if err != nil {
			return nil, errors.New().WithCause(err).Err()
		}
		label = label.Add(*ls)
	}

	return label, nil
}

func GetPodList(client clientset.Interface, namespace, objectType, objectName string) ([]*KubeObjectInfo, error) {
	var podList []*KubeObjectInfo

	label, err := getLabels(client, namespace, objectType, objectName)
	if err != nil {
		return nil, errors.New().WithCause(err).Err()
	}

	pods, err := client.CoreV1().Pods(namespace).List(metav1.ListOptions{LabelSelector: label.String()})
	if err != nil {
		return nil, errors.New().WithCause(err).Err()
	}

	for _, pod := range pods.Items {
		podList = append(podList, &KubeObjectInfo{Name: pod.Name + "@" + namespace, IP: pod.Status.PodIP, GroupName: objectName, GroupType: objectType})
	}

	return podList, nil
}

func GetPod(client clientset.Interface, namespace, objectType, objectName, podName string) ([]*KubeObjectInfo, error) {
	var podList []*KubeObjectInfo
	pod, err := client.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New().WithCause(err).Err()
	}
	podList = append(podList, &KubeObjectInfo{Name: pod.Name + "@" + namespace, IP: pod.Status.PodIP, GroupName: objectName, GroupType: objectType})
	return podList, nil
}

func GetNodeList(client clientset.Interface, alertNamespace string) ([]*KubeObjectInfo, error) {
	var nodeList []*KubeObjectInfo
	nodes, err := client.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: labels.Everything().String()})
	if err != nil {
		return nodeList, errors.New().WithCause(err).Err()
	}
	for _, node := range nodes.Items {
		nodeIP := "127.0.0.1"
		for _, ip := range node.Status.Addresses {
			if ip.Type == internalIP {
				nodeIP = ip.Address
				break
			}
		}
		nodeList = append(nodeList, &KubeObjectInfo{Name: node.Name + "@" + alertNamespace, IP: nodeIP, GroupName: TypeNodes, GroupType: ""})
	}
	return nodeList, nil
}

func GetNode(client clientset.Interface, nodeName, alertNamespace string) ([]*KubeObjectInfo, error) {
	var nodeList []*KubeObjectInfo
	node := &apiv1.Node{}
	node, err := client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return nodeList, errors.New().WithCause(err).Err()
	}
	nodeIP := "127.0.0.1"
	for _, ip := range node.Status.Addresses {
		if ip.Type == internalIP {
			nodeIP = ip.Address
			break
		}
	}
	nodeList = append(nodeList, &KubeObjectInfo{Name: node.Name + "@" + alertNamespace, IP: nodeIP, GroupName: TypeNodes, GroupType: ""})
	return nodeList, nil
}

func GetAlertList(acExtClient acs.ExtensionInterface, kubeClient clientset.Interface, namespace string, ls labels.Selector) ([]aci.Alert, error) {
	alerts := make([]aci.Alert, 0)
	if namespace != "" {
		alertList, err := acExtClient.Alert(namespace).List(metav1.ListOptions{LabelSelector: ls.String()})
		if err != nil {
			return nil, errors.New().WithCause(err).Err()
		}
		if len(alertList.Items) > 0 {
			alerts = append(alerts, alertList.Items...)
		}
	} else {
		alertList, err := acExtClient.Alert(apiv1.NamespaceAll).List(metav1.ListOptions{LabelSelector: ls.String()})
		if err != nil {
			return nil, errors.New().WithCause(err).Err()
		}
		if len(alertList.Items) > 0 {
			alerts = append(alerts, alertList.Items...)
		}
	}

	return alerts, nil
}

func GetAlert(acExtClient acs.ExtensionInterface, namespace, name string) (*aci.Alert, error) {
	return acExtClient.Alert(namespace).Get(name)
}

const (
	ObjectType = "alert.appscode.com/objectType"
	ObjectName = "alert.appscode.com/objectName"
)

func GetLabelSelector(objectType, objectName string) (labels.Selector, error) {
	lb := labels.NewSelector()
	if objectType != "" {
		lsot, err := labels.NewRequirement(ObjectType, selection.Equals, sets.NewString(objectType).List())
		if err != nil {
			return lb, errors.New().WithCause(err).Err()
		}
		lb = lb.Add(*lsot)
	}

	if objectName != "" {
		lson, err := labels.NewRequirement(ObjectName, selection.Equals, sets.NewString(objectName).List())
		if err != nil {
			return lb, errors.New().WithCause(err).Err()
		}
		lb = lb.Add(*lson)
	}

	return lb, nil
}

type labelMap map[string]string

func (s labelMap) ObjectType() string {
	v, _ := s[ObjectType]
	return v
}

func (s labelMap) ObjectName() string {
	v, _ := s[ObjectName]
	return v
}

func GetObjectInfo(label map[string]string) (objectType string, objectName string) {
	opts := labelMap(label)
	objectType = opts.ObjectType()
	objectName = opts.ObjectName()
	return
}

func CheckAlertConfig(oldConfig, newConfig *aci.Alert) error {
	oldOpts := labelMap(oldConfig.ObjectMeta.Labels)
	newOpts := labelMap(newConfig.ObjectMeta.Labels)

	if newOpts.ObjectType() != oldOpts.ObjectType() {
		return errors.New("Kubernetes ObjectType mismatch")
	}

	if newOpts.ObjectName() != oldOpts.ObjectName() {
		return errors.New("Kubernetes ObjectName mismatch")
	}

	if newConfig.Spec.CheckCommand != oldConfig.Spec.CheckCommand {
		return errors.New("CheckCommand mismatch")
	}

	return nil
}

func IsIcingaApp(ancestors []*types.Ancestors, namespace string) bool {
	icingaServiceNamespace := os.Getenv("ICINGA_SERVICE_NAMESPACE")
	if icingaServiceNamespace != namespace {
		return false
	}

	icingaService := os.Getenv("ICINGA_SERVICE_NAME")

	for _, ancestor := range ancestors {
		if ancestor.Type == events.Service.String() {
			for _, service := range ancestor.Names {
				if service == icingaService {
					return true
				}
			}
		}
	}
	return false
}
