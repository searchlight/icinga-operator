package operator

import (
	"reflect"

	"github.com/appscode/go/log"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	"github.com/appscode/searchlight/apis/monitoring"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
)

func (op *Operator) initClusterAlertWatcher() {
	op.caInformer = op.monInformerFactory.Monitoring().V1alpha1().ClusterAlerts().Informer()
	op.caQueue = queue.New("ClusterAlert", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcileClusterAlert)
	op.caInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			alert := obj.(*api.ClusterAlert)
			if op.isValid(alert) {
				queue.Enqueue(op.caQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*api.ClusterAlert)
			nu := newObj.(*api.ClusterAlert)

			if op.processClusterAlertUpdate(old, nu) {
				queue.Enqueue(op.caQueue.GetQueue(), nu)
				return
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.caQueue.GetQueue(), obj)
		},
	})
	op.caLister = op.monInformerFactory.Monitoring().V1alpha1().ClusterAlerts().Lister()
}

func (op *Operator) processClusterAlertUpdate(old, nu *api.ClusterAlert) bool {
	if nu.DeletionTimestamp != nil && core_util.HasFinalizer(nu.ObjectMeta, monitoring.GroupName) {
		return true
	}
	if !reflect.DeepEqual(old.Spec, nu.Spec) && op.isValid(nu) {
		return true
	}
	return false
}

func (op *Operator) reconcileClusterAlert(key string) error {
	obj, exists, err := op.caInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}
	if !exists {
		log.Debugf("ClusterAlert %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		return op.clusterHost.Delete(namespace, name)
	}

	alert := obj.(*api.ClusterAlert)
	log.Infof("Sync/Add/Update for ClusterAlert %s\n", alert.GetName())

	err = op.clusterHost.Reconcile(alert.DeepCopy())
	if err != nil {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToSync,
			`Reason: %v`,
			err,
		)
	}
	return err
}
