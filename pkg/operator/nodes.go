package operator

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/appscode/go/log"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initNodeWatcher() {
	op.nInformer = op.kubeInformerFactory.Core().V1().Nodes().Informer()
	op.nQueue = queue.New("Node", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcileNode)
	op.nInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			queue.Enqueue(op.nQueue.GetQueue(), obj)
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldNode, ok := old.(*core.Node)
			if !ok {
				return
			}
			newNode, ok := new.(*core.Node)
			if !ok {
				return
			}
			if !reflect.DeepEqual(oldNode.Labels, newNode.Labels) {
				queue.Enqueue(op.nQueue.GetQueue(), new)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.nQueue.GetQueue(), obj)
		},
	})
	op.nLister = op.kubeInformerFactory.Core().V1().Nodes().Lister()
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) reconcileNode(key string) error {
	obj, exists, err := op.nInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		log.Debugf("Node %s does not exist anymore\n", key)
		_, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}

		if err := op.ForceDeleteIcingaObjectsForNode(name); err != nil {
			log.Errorf("Failed to delete alert for Node %s", name)
		}
	} else {
		node := obj.(*core.Node)
		fmt.Printf("Sync/Add/Update for Node %s\n", node.GetName())
		if err := op.EnsureNode(node); err != nil {
			log.Errorf("Failed to patch alert for Node %s@%s", node.Name, node.Namespace)
		}
	}
	return nil
}

func (op *Operator) EnsureNode(node *core.Node) error {
	fmt.Printf("Sync/Add/Update for Node %s\n", node.GetName())

	oldAlerts := make([]*api.NodeAlert, 0)

	oldAlertNames := make([]string, 0)
	if val, ok := node.Annotations[annotationAlertsName]; ok {
		if err := json.Unmarshal([]byte(val), &oldAlertNames); err != nil {
			return err
		}
	}
	for _, l := range oldAlertNames {
		namespace, name, err := cache.SplitMetaNamespaceKey(l)
		if err != nil {
			return err
		}

		oldAlerts = append(oldAlerts, &api.NodeAlert{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		})
	}

	newAlerts, err := util.FindNodeAlert(op.naLister, node.ObjectMeta)
	if err != nil {
		return err
	}

	type change struct {
		old *api.NodeAlert
		new *api.NodeAlert
	}
	diff := make(map[string]*change)
	for i := range oldAlerts {
		diff[oldAlerts[i].Name] = &change{old: oldAlerts[i]}
	}

	alertNames := make([]string, 0)

	for i := range newAlerts {
		alertNames = append(alertNames, fmt.Sprintf("%s/%s", newAlerts[i].Namespace, newAlerts[i].Name))
		if ch, ok := diff[newAlerts[i].Name]; ok {
			ch.new = newAlerts[i]
		} else {
			diff[newAlerts[i].Name] = &change{new: newAlerts[i]}
		}
	}

	for alert := range diff {
		ch := diff[alert]
		if ch.old != nil && ch.new == nil {
			op.EnsureIcingaNodeAlertDeleted(ch.old, node)
		} else {
			op.EnsureIcingaNodeAlert(ch.new, node)
		}
	}

	_, vt, err := core_util.PatchNode(op.KubeClient, node, func(in *core.Node) *core.Node {
		if len(newAlerts) > 0 {
			if in.Annotations == nil {
				in.Annotations = make(map[string]string, 0)
			}
			data, _ := json.Marshal(alertNames)
			in.Annotations[annotationAlertsName] = string(data)
		} else {
			delete(in.Annotations, annotationAlertsName)
		}
		return in
	})
	if err != nil {
		log.Errorf("Failed to %v Node %s", vt, node.Name)
	}
	return nil
}

func (op *Operator) ForceDeleteIcingaObjectsForNode(name string) error {
	namespaceList, err := op.KubeClient.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, ns := range namespaceList.Items {
		nh := icinga.IcingaHost{
			ObjectName:     name,
			Type:           icinga.TypeNode,
			AlertNamespace: ns.Name,
		}
		err := op.nodeHost.ForceDeleteIcingaHost(nh)

		if err != nil {
			host, _ := nh.Name()
			log.Errorf(`Failed to delete Icinga Host "%s" for Node %s`, host, name)
		}
	}
	return nil
}
