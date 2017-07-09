package icinga

import (
	"os"

	"github.com/appscode/errors"
	"github.com/appscode/kubed/pkg/events"
	tapi "github.com/appscode/searchlight/api"
	tcs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/controller/types"
	//"github.com/appscode/searchlight/pkg/events"
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

func GetPod(client clientset.Interface, namespace, objectType, objectName, podName string) ([]*KHost, error) {
	var podList []*KHost
	pod, err := client.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.FromErr(err).Err()
	}
	podList = append(podList, &KHost{Name: pod.Name + "@" + namespace, IP: pod.Status.PodIP, GroupName: objectName, GroupType: objectType})
	return podList, nil
}

func GetNode(client clientset.Interface, nodeName, alertNamespace string) ([]*KHost, error) {
	var nodeList []*KHost
	node := &apiv1.Node{}
	node, err := client.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return nodeList, errors.FromErr(err).Err()
	}
	nodeIP := "127.0.0.1"
	for _, ip := range node.Status.Addresses {
		if ip.Type == internalIP {
			nodeIP = ip.Address
			break
		}
	}
	nodeList = append(nodeList, &KHost{Name: node.Name + "@" + alertNamespace, IP: nodeIP, GroupName: TypeNodes, GroupType: ""})
	return nodeList, nil
}

func GetAlert(acExtClient tcs.ExtensionInterface, namespace, name string) (*tapi.PodAlert, error) {
	return acExtClient.PodAlerts(namespace).Get(name)
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
			return lb, errors.FromErr(err).Err()
		}
		lb = lb.Add(*lsot)
	}

	if objectName != "" {
		lson, err := labels.NewRequirement(ObjectName, selection.Equals, sets.NewString(objectName).List())
		if err != nil {
			return lb, errors.FromErr(err).Err()
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

func CheckAlertConfig(oldConfig, newConfig *tapi.PodAlert) error {
	oldOpts := labelMap(oldConfig.ObjectMeta.Labels)
	newOpts := labelMap(newConfig.ObjectMeta.Labels)

	if newOpts.ObjectType() != oldOpts.ObjectType() {
		return errors.New("Kubernetes ObjectType mismatch")
	}

	if newOpts.ObjectName() != oldOpts.ObjectName() {
		return errors.New("Kubernetes ObjectName mismatch")
	}

	if newConfig.Spec.Check != oldConfig.Spec.Check {
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
