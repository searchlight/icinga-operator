package operator

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/appscode/go/log"
	"github.com/appscode/go/sets"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	mon_api "github.com/appscode/searchlight/apis/monitoring"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	mon_util "github.com/appscode/searchlight/client/clientset/versioned/typed/monitoring/v1alpha1/util"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// nodeAlertMapperConfiguration
type naMapperConf struct {
	Selector map[string]string `json:"selector,omitempty"`
	NodeName string            `json:"nodeName,omitempty"`
}

func (op *Operator) initNodeAlertWatcher() {
	op.naInformer = op.searchlightInformerFactory.Monitoring().V1alpha1().NodeAlerts().Informer()
	op.naQueue = queue.New("NodeAlert", op.options.MaxNumRequeues, op.options.NumThreads, op.reconcileNodeAlert)
	op.naInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if alert, ok := obj.(*api.NodeAlert); ok {
				if !op.validateNodeAlert(alert) {
					log.Errorf(`Invalid NodeAlert "%s@%s"`, alert.Name, alert.Namespace)
					return
				}
				queue.Enqueue(op.naQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			oldAlert, ok := old.(*api.NodeAlert)
			if !ok {
				return
			}
			newAlert, ok := new.(*api.NodeAlert)
			if !ok {
				return
			}
			// DeepEqual old & new
			// DeepEqual MapperConfiguration of old & new
			// Patch PodAlert with necessary annotation
			newAlert, err := op.processNodeAlertUpdate(oldAlert, newAlert)
			if err != nil {
				log.Error(err)
			} else {
				queue.Enqueue(op.naQueue.GetQueue(), newAlert)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.naQueue.GetQueue(), obj)
		},
	})
	op.naLister = op.searchlightInformerFactory.Monitoring().V1alpha1().NodeAlerts().Lister()
}

// syncToStdout is the business logic of the controller. In this controller it simply prints
// information about the deployment to stdout. In case an error happened, it has to simply return the error.
// The retry logic should not be part of the business logic.
func (op *Operator) reconcileNodeAlert(key string) error {
	obj, exists, err := op.naInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		log.Debugf("NodeAlert %s does not exist anymore\n", key)
	} else {
		alert := obj.(*api.NodeAlert)
		if alert.DeletionTimestamp != nil {
			if core_util.HasFinalizer(alert.ObjectMeta, mon_api.GroupName) {
				// Delete all Icinga objects created for this NodeAlert
				if err := op.EnsureNodeAlertDeleted(alert); err != nil {
					log.Errorf("Failed to delete NodeAlert %s@%s. Reason: %v", alert.Name, alert.Namespace, err)
					return err
				}
				// Remove Finalizer
				_, _, err = mon_util.PatchNodeAlert(op.ExtClient.MonitoringV1alpha1(), alert, func(in *api.NodeAlert) *api.NodeAlert {
					in.ObjectMeta = core_util.RemoveFinalizer(in.ObjectMeta, mon_api.GroupName)
					return in
				})
				return err
			}
		} else {
			fmt.Printf("Sync/Add/Update for NodeAlert %s\n", alert.GetName())

			alert, _, err = mon_util.PatchNodeAlert(op.ExtClient.MonitoringV1alpha1(), alert, func(in *api.NodeAlert) *api.NodeAlert {
				in.ObjectMeta = core_util.AddFinalizer(in.ObjectMeta, mon_api.GroupName)
				return in
			})

			if err := op.EnsureNodeAlert(alert); err != nil {
				log.Errorf("Failed to sync NodeAlert %s@%s. Reason: %v", alert.Name, alert.Namespace, err)
			}
		}
	}
	return nil
}

func (op *Operator) EnsureNodeAlert(alert *api.NodeAlert) error {

	var oldMc *naMapperConf
	if val, ok := alert.Annotations[annotationLastConfiguration]; ok {
		if err := json.Unmarshal([]byte(val), &oldMc); err != nil {
			return err
		}
	}

	oldMappedNode, err := op.getMappedNodeList(alert.Namespace, oldMc)
	if err != nil {
		return err
	}

	newMC := &naMapperConf{
		Selector: alert.Spec.Selector,
		NodeName: alert.Spec.NodeName,
	}
	newMappedNode, err := op.getMappedNodeList(alert.Namespace, newMC)
	if err != nil {
		return err
	}

	for key, node := range newMappedNode {
		delete(oldMappedNode, node.Name)

		op.setNodeAlertNamesInAnnotation(node, alert)

		go op.EnsureIcingaNodeAlert(alert, newMappedNode[key])
	}

	for _, node := range oldMappedNode {
		op.EnsureIcingaNodeAlertDeleted(alert, node)
	}

	return nil
}

func (op *Operator) EnsureNodeAlertDeleted(alert *api.NodeAlert) error {
	mc := &naMapperConf{
		Selector: alert.Spec.Selector,
		NodeName: alert.Spec.NodeName,
	}
	mappedPod, err := op.getMappedNodeList(alert.Namespace, mc)
	if err != nil {
		return err
	}

	for _, node := range mappedPod {
		op.EnsureIcingaNodeAlertDeleted(alert, node)
	}

	return nil
}

func (op *Operator) EnsureIcingaNodeAlert(alert *api.NodeAlert, node *core.Node) (err error) {
	defer func() {
		if err == nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeNormal,
				eventer.EventReasonSuccessfulSync,
				`Applied NodeAlert: "%v"`,
				alert.Name,
			)
			return
		} else {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`Fail to be apply NodeAlert: "%v". Reason: %v`,
				alert.Name,
				err,
			)
			log.Errorln(err)
			return
		}
	}()
	err = op.nodeHost.Create(alert.DeepCopy(), node.DeepCopy())
	return
}

func (op *Operator) EnsureIcingaNodeAlertDeleted(alert *api.NodeAlert, node *core.Node) (err error) {
	defer func() {
		if err == nil {
			if alert != nil {
				op.recorder.Eventf(
					alert.ObjectReference(),
					core.EventTypeNormal,
					eventer.EventReasonSuccessfulDelete,
					`Deleted Icinga objects of NodeAlert "%s@%s" for Node "%s"`,
					alert.Name, alert.Namespace, node.Name,
				)
			}
			return
		} else {
			if alert != nil {
				op.recorder.Eventf(
					alert.ObjectReference(),
					core.EventTypeWarning,
					eventer.EventReasonFailedToDelete,
					`Fail to delete Icinga objects of NodeAlert "%s@%s" for Node "%s". Reason: %v`,
					alert.Name, alert.Namespace, node.Name,
					err,
				)
			}
			log.Errorln(err)
			return
		}
	}()
	err = op.nodeHost.Delete(alert, node)
	return
}

func (op *Operator) processNodeAlertUpdate(oldAlert, newAlert *api.NodeAlert) (*api.NodeAlert, error) {
	// Check for changes in Spec
	if !reflect.DeepEqual(oldAlert.Spec, newAlert.Spec) {
		if !op.validateNodeAlert(newAlert) {
			return nil, errors.Errorf(`Invalid NodeAlert "%s@%s"`, newAlert.Name, newAlert.Namespace)
		}

		// We need Selector/NodeName from oldAlert while processing this update operation.
		// Because we need to remove Icinga objects for oldAlert.
		oldMC := &naMapperConf{
			Selector: oldAlert.Spec.Selector,
			NodeName: oldAlert.Spec.NodeName,
		}
		newMC := &naMapperConf{
			Selector: newAlert.Spec.Selector,
			NodeName: newAlert.Spec.NodeName,
		}

		// We will store Selector/PodName from oldAlert in annotation
		if !reflect.DeepEqual(oldMC, newMC) {
			var err error
			// Patch NodeAlert with Selector/PodName from oldAlert (oldMC)
			newAlert, _, err = mon_util.PatchNodeAlert(op.ExtClient.MonitoringV1alpha1(), newAlert, func(in *api.NodeAlert) *api.NodeAlert {
				if in.Annotations == nil {
					in.Annotations = make(map[string]string, 0)
				}
				data, _ := json.Marshal(oldMC)
				in.Annotations[annotationLastConfiguration] = string(data)
				return in
			})
			if err != nil {
				op.recorder.Eventf(
					newAlert.ObjectReference(),
					core.EventTypeWarning,
					eventer.EventReasonFailedToUpdate,
					`Reason: %v`,
					err,
				)
				return nil, errors.WithMessage(err,
					fmt.Sprintf(`Failed to patch PodAlert "%s@%s"`, newAlert.Name, newAlert.Namespace),
				)
			}
		}
	}

	return newAlert, nil
}

func (op *Operator) getMappedNodeList(namespace string, mc *naMapperConf) (map[string]*core.Node, error) {
	mappedPodList := make(map[string]*core.Node)

	if mc == nil {
		return mappedPodList, nil
	}

	sel := labels.SelectorFromSet(mc.Selector)

	if mc.NodeName != "" {
		if node, err := op.KubeClient.CoreV1().Nodes().Get(mc.NodeName, metav1.GetOptions{}); err == nil {
			if sel.Matches(labels.Set(node.Labels)) {
				mappedPodList[node.Name] = node
			}
		}
	} else {
		if nodeList, err := op.KubeClient.CoreV1().Nodes().List(metav1.ListOptions{LabelSelector: sel.String()}); err != nil {
			return nil, err
		} else {
			for i, node := range nodeList.Items {
				mappedPodList[node.Name] = &nodeList.Items[i]
			}
		}
	}

	return mappedPodList, nil
}

func (op *Operator) setNodeAlertNamesInAnnotation(node *core.Node, alert *api.NodeAlert) {
	_, _, err := core_util.PatchNode(op.KubeClient, node, func(in *core.Node) *core.Node {
		if in.Annotations == nil {
			in.Annotations = make(map[string]string, 0)
		}

		alertNames := make([]string, 0)
		if val, ok := alert.Annotations[annotationAlertsName]; ok {
			if err := json.Unmarshal([]byte(val), &alertNames); err != nil {
				log.Errorf("Failed to patch Node %s.", node.Name)
			}
		}
		ss := sets.NewString(alertNames...)
		ss.Insert(alert.Name)
		alertNames = ss.List()
		data, _ := json.Marshal(alertNames)
		in.Annotations[annotationAlertsName] = string(data)
		return in
	})

	if err != nil {
		log.Errorf("Failed to patch Node %s.", node.Name)
	}
}

func (op *Operator) validateNodeAlert(alert *api.NodeAlert) bool {
	// Validate IcingaCommand & it's variables.
	// And also check supported IcingaState
	if ok, err := alert.IsValid(); !ok {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonFailedToCreate,
			`Reason: %v`,
			err,
		)
		return false
	}

	// Validate Notifiers configurations
	if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonBadNotifier,
			`Bad notifier config for NodeAlert: "%s@%s". Reason: %v`,
			alert.Name, alert.Namespace,
			err,
		)
		return false
	}

	return true
}
