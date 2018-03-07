package operator

import (
	"strings"

	"github.com/appscode/go/log"
	"github.com/appscode/go/sets"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initNodeAlertWatcher() {
	op.naInformer = op.monInformerFactory.Monitoring().V1alpha1().NodeAlerts().Informer()
	op.naQueue = queue.New("NodeAlert", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcileNodeAlert)
	op.naInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			alert := obj.(*api.NodeAlert)
			if err := op.isValid(alert); err == nil {
				queue.Enqueue(op.naQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*api.NodeAlert)
			nu := newObj.(*api.NodeAlert)

			if err := op.isValid(nu); err != nil {
				return
			}
			if !equalNodeAlert(old, nu) {
				queue.Enqueue(op.naQueue.GetQueue(), nu)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.naQueue.GetQueue(), obj)
		},
	})
	op.naLister = op.monInformerFactory.Monitoring().V1alpha1().NodeAlerts().Lister()
}

func (op *Operator) reconcileNodeAlert(key string) error {
	obj, exists, err := op.naInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		log.Warningf("NodeAlert %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		return op.EnsureNodeAlertDeleted(namespace, name)
	}

	alert := obj.(*api.NodeAlert)
	log.Infof("Sync/Add/Update for NodeAlert %s\n", key)

	op.EnsureNodeAlert(alert)
	op.EnsureNodeAlertDeleted(alert.Namespace, alert.Name)
	return nil
}

func (op *Operator) EnsureNodeAlert(alert *api.NodeAlert) error {
	if alert.Spec.NodeName != nil {
		node, err := op.nodeLister.Get(*alert.Spec.NodeName)
		if err != nil {
			return err
		}
		key, err := cache.MetaNamespaceKeyFunc(node)
		if err == nil {
			op.nodeQueue.GetQueue().Add(key)
		}
	}

	sel := labels.SelectorFromSet(alert.Spec.Selector)
	nodes, err := op.nodeLister.List(sel)
	if err != nil {
		return err
	}
	for i := range nodes {
		node := nodes[i]
		key, err := cache.MetaNamespaceKeyFunc(node)
		if err == nil {
			op.nodeQueue.GetQueue().Add(key)
		}
	}
	return nil
}

func GetAppliedNodeAlerts(a map[string]string, key string) bool {
	if a == nil {
		return false
	}
	if val, ok := a[annotationAlertsName]; ok {
		names := strings.Split(val, ",")
		return sets.NewString(names...).Has(key)
	}
	return false
}

func (op *Operator) EnsureNodeAlertDeleted(alertNamespace, alertName string) error {
	nodes, err := op.nodeLister.List(labels.Everything())
	if err != nil {
		return err
	}
	alertKey := alertNamespace + "/" + alertName
	for _, node := range nodes {
		if GetAppliedNodeAlerts(node.Annotations, alertKey) {
			key, err := cache.MetaNamespaceKeyFunc(node)
			if err == nil {
				op.nodeQueue.GetQueue().Add(key)
			}
		}
	}
	return nil
}

func (op *Operator) EnsureIcingaNodeAlert(alert *api.NodeAlert, node *core.Node) (err error) {
	err = op.nodeHost.Reconcile(alert.DeepCopy(), node.DeepCopy())
	if err != nil {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToSync,
			`Reason: %v`,
			err,
		)
	}
	return
}

func (op *Operator) EnsureIcingaNodeAlertDeleted(alert *api.NodeAlert, node *core.Node) (err error) {
	err = op.nodeHost.Delete(alert, node)
	if err != nil && alert != nil {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToDelete,
			`Reason: %v`,
			err,
		)
	}
	return
}
