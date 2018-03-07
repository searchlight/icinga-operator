package operator

import (
	"reflect"
	"strings"

	"github.com/appscode/go/log"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initNodeWatcher() {
	op.nInformer = op.kubeInformerFactory.Core().V1().Nodes().Informer()
	op.nQueue = queue.New("Node", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcileNode)
	op.nInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			queue.Enqueue(op.nQueue.GetQueue(), obj)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*core.Node)
			nu := newObj.(*core.Node)
			if !reflect.DeepEqual(old.Labels, nu.Labels) {
				queue.Enqueue(op.nQueue.GetQueue(), newObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.nQueue.GetQueue(), obj)
		},
	})
	op.nLister = op.kubeInformerFactory.Core().V1().Nodes().Lister()
}

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
		log.Infof("Sync/Add/Update for Node %s\n", key)

		node := obj.(*core.Node)
		if err := op.EnsureNode(node); err != nil {
			log.Errorf("Failed to patch alert for Node %s@%s", node.Name, node.Namespace)
		}
	}
	return nil
}

func (op *Operator) EnsureNode(node *core.Node) error {
	oldAlerts := sets.NewString()
	if val, ok := node.Annotations[annotationAlertsName]; ok {
		keys := strings.Split(val, ",")
		oldAlerts.Insert(keys...)
	}

	newAlerts, err := util.FindNodeAlert(op.naLister, node.ObjectMeta)
	if err != nil {
		return err
	}
	newKeys := make([]string, len(newAlerts))
	for i := range newAlerts {
		alert := newAlerts[i]
		op.EnsureIcingaNodeAlert(alert, node)

		key, err := cache.MetaNamespaceKeyFunc(alert)
		if err != nil {
			return err
		}
		newKeys[i] = key
		if oldAlerts.Has(key) {
			oldAlerts.Delete(key)
		}
	}

	for _, key := range oldAlerts.List() {
		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		op.EnsureIcingaNodeAlertDeleted(&api.NodeAlert{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}, node)
	}

	_, vt, err := core_util.PatchNode(op.KubeClient, node, func(in *core.Node) *core.Node {
		if in.Annotations == nil {
			in.Annotations = make(map[string]string, 0)
		}
		if len(newKeys) > 0 {
			in.Annotations[annotationAlertsName] = strings.Join(newKeys, ",")
		} else {
			delete(in.Annotations, annotationAlertsName)
		}
		return in
	})
	if err != nil {
		log.Errorf("Failed to %v Node %s", vt, node.Name)
	}
	return err
}

func (op *Operator) ForceDeleteIcingaObjectsForNode(name string) error {
	namespaces, err := op.nsLister.List(labels.Everything())
	if err != nil {
		return err
	}
	for _, ns := range namespaces {
		h := icinga.IcingaHost{
			ObjectName:     name,
			Type:           icinga.TypeNode,
			AlertNamespace: ns.Name,
		}
		err := op.nodeHost.ForceDeleteIcingaHost(h)
		if err != nil {
			host, _ := h.Name()
			log.Errorf(`Failed to delete Icinga Host "%s" for Node %s`, host, name)
		}
	}
	return nil
}
