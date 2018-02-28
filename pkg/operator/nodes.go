package operator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/appscode/go/log"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
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
				log.Errorln("invalid Node object")
				return
			}
			newNode, ok := new.(*core.Node)
			if !ok {
				log.Errorln("invalid Node object")
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
		node := obj.(*core.Node)
		// Below we will warm up our cache with a Node, so that we will see a delete for one d
		fmt.Printf("Node %s does not exist anymore\n", key)

		if err := op.EnsureNodeDeleted(node); err != nil {
			log.Errorf("Failed to delete alert for Node %s@%s", node.Name, node.Namespace)
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
	if names, ok := node.Annotations[annotationAlertsName]; ok {
		list := strings.Split(names, ",")
		for _, l := range list {
			oldAlerts = append(oldAlerts, &api.NodeAlert{
				ObjectMeta: metav1.ObjectMeta{
					Name:      strings.Trim(l, " "),
					Namespace: node.Namespace,
				},
			})
		}
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

	newAlertNameList := make([]string, 0)

	for i := range newAlerts {
		newAlertNameList = append(newAlertNameList, newAlerts[i].Name)
		if ch, ok := diff[newAlerts[i].Name]; ok {
			ch.new = newAlerts[i]
		} else {
			diff[newAlerts[i].Name] = &change{new: newAlerts[i]}
		}
	}

	for alert := range diff {
		ch := diff[alert]
		if ch.old == nil && ch.new != nil {
			go op.EnsureIcingaNodeAlert(ch.new, node)
		} else if ch.old != nil && ch.new == nil {
			go op.EnsureIcingaNodeAlertDeleted(ch.old, node)
		} else if ch.old != nil && ch.new != nil && !reflect.DeepEqual(ch.old.Spec, ch.new.Spec) {
			go op.EnsureIcingaNodeAlert(ch.new, node)
		}
	}

	_, vr, err := core_util.PatchNode(op.KubeClient, node, func(in *core.Node) *core.Node {
		if in.Annotations == nil {
			in.Annotations = make(map[string]string, 0)
		}
		in.Annotations[annotationAlertsName] = strings.Join(newAlertNameList, ", ")
		return in
	})
	if err != nil {
		log.Errorf("Failed to %v Node %s", vr, node.Name)
	}

	return nil
}

func (op *Operator) EnsureNodeDeleted(node *core.Node) error {
	alerts, err := util.FindNodeAlert(op.naLister, node.ObjectMeta)
	if err != nil {
		log.Errorf("Error while searching NodeAlert for Node %s@%s.", node.Name, node.Namespace)
		return err
	}
	if len(alerts) == 0 {
		log.Errorf("No NodeAlert found for Node %s@%s.", node.Name, node.Namespace)
		return err
	}
	for _, alert := range alerts {
		go op.EnsureIcingaNodeAlertDeleted(alert, node)
	}

	return nil
}
