package operator

import (
	"reflect"
	"strings"

	"github.com/appscode/go/log"
	utilerrors "github.com/appscode/go/util/errors"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initNodeWatcher() {
	op.nodeInformer = op.kubeInformerFactory.Core().V1().Nodes().Informer()
	op.nodeQueue = queue.New("Node", op.MaxNumRequeues, op.NumThreads, op.reconcileNode)
	op.nodeInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			queue.Enqueue(op.nodeQueue.GetQueue(), obj)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*core.Node)
			nu := newObj.(*core.Node)
			if !reflect.DeepEqual(old.Labels, nu.Labels) {
				queue.Enqueue(op.nodeQueue.GetQueue(), newObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.nodeQueue.GetQueue(), obj)
		},
	})
	op.nodeLister = op.kubeInformerFactory.Core().V1().Nodes().Lister()
}

func (op *Operator) reconcileNode(key string) error {
	obj, exists, err := op.nodeInformer.GetIndexer().GetByKey(key)
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

		return op.forceDeleteIcingaObjectsForNode(name)
	}

	log.Infof("Sync/Add/Update for Node %s\n", key)

	node := obj.(*core.Node).DeepCopy()
	err = op.ensureNode(node)
	if err != nil {
		log.Errorf("failed to reconcile alert for node %s. reason: %s", key, err)
	}
	return err
}

func (op *Operator) ensureNode(node *core.Node) error {
	var errlist []error

	oldAlerts := sets.NewString()
	if val, ok := node.Annotations[api.AnnotationKeyAlerts]; ok {
		keys := strings.Split(val, ",")
		oldAlerts.Insert(keys...)
	}

	newAlerts, err := findNodeAlert(op.kubeClient, op.naLister, node.ObjectMeta)
	if err != nil {
		return err
	}
	newKeys := make([]string, len(newAlerts))
	for i := range newAlerts {
		alert := newAlerts[i]

		err = op.nodeHost.Apply(alert, node)
		if err != nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`failed to apply to node %s. Reason: %v`,
				node.Name, err,
			)
			errlist = append(errlist, err)
		}

		key, _ := cache.MetaNamespaceKeyFunc(alert)
		newKeys[i] = key
		if oldAlerts.Has(key) {
			oldAlerts.Delete(key)
		}
	}

	for _, key := range oldAlerts.List() {
		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			// ignore
			continue
		}

		err = op.nodeHost.Delete(namespace, name, node)
		if err != nil {
			if alert, e2 := op.naLister.NodeAlerts(namespace).Get(name); e2 == nil {
				op.recorder.Eventf(
					alert.ObjectReference(),
					core.EventTypeWarning,
					eventer.EventReasonFailedToDelete,
					`failed to delete for node %s. Reason: %s`,
					node.Name, err,
				)
			}
			errlist = append(errlist, err)
		}
	}

	_, _, err = core_util.PatchNode(op.kubeClient, node, func(in *core.Node) *core.Node {
		if in.Annotations == nil {
			in.Annotations = make(map[string]string, 0)
		}
		if len(newKeys) > 0 {
			in.Annotations[api.AnnotationKeyAlerts] = strings.Join(newKeys, ",")
		} else {
			delete(in.Annotations, api.AnnotationKeyAlerts)
		}
		return in
	})
	if err != nil {
		errlist = append(errlist, err)
	}
	return utilerrors.NewAggregate(errlist)
}

func (op *Operator) forceDeleteIcingaObjectsForNode(name string) error {
	namespaces, err := op.nsLister.List(labels.Everything())
	if err != nil {
		return err
	}

	var errlist []error
	for _, ns := range namespaces {
		h := icinga.IcingaHost{
			ObjectName:     name,
			Type:           icinga.TypeNode,
			AlertNamespace: ns.Name,
		}
		if err := op.nodeHost.ForceDeleteIcingaHost(h); err != nil {
			errlist = append(errlist, err)
		}
	}
	return utilerrors.NewAggregate(errlist)
}
