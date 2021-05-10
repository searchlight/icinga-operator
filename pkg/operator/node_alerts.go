/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

import (
	"reflect"
	"strings"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"

	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/tools/queue"
)

func (op *Operator) initNodeAlertWatcher() {
	op.naInformer = op.monInformerFactory.Monitoring().V1alpha1().NodeAlerts().Informer()
	op.naQueue = queue.New("NodeAlert", op.MaxNumRequeues, op.NumThreads, op.reconcileNodeAlert)
	op.naInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			alert := obj.(*api.NodeAlert)
			if op.isValid(alert) {
				queue.Enqueue(op.naQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*api.NodeAlert)
			nu := newObj.(*api.NodeAlert)

			if reflect.DeepEqual(old.Spec, nu.Spec) {
				return
			}
			if op.isValid(nu) {
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
		klog.Warningf("NodeAlert %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		return op.ensureNodeAlertDeleted(namespace, name)
	}

	alert := obj.(*api.NodeAlert).DeepCopy()
	klog.Infof("Sync/Add/Update for NodeAlert %s\n", key)

	op.ensureNodeAlert(alert)
	op.ensureNodeAlertDeleted(alert.Namespace, alert.Name)
	return nil
}

func (op *Operator) ensureNodeAlert(alert *api.NodeAlert) error {
	if alert.Spec.NodeName != nil {
		node, err := op.nodeLister.Get(*alert.Spec.NodeName)
		if err != nil {
			return err
		}
		key, err := cache.MetaNamespaceKeyFunc(node)
		if err == nil {
			op.nodeQueue.GetQueue().Add(key)
		}
		return nil
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

func alertAppliedToNode(a map[string]string, key string) bool {
	if a == nil {
		return false
	}
	if val, ok := a[api.AnnotationKeyAlerts]; ok {
		names := strings.Split(val, ",")
		for _, name := range names {
			if name == key {
				return true
			}
		}
	}
	return false
}

func (op *Operator) ensureNodeAlertDeleted(alertNamespace, alertName string) error {
	nodes, err := op.nodeLister.List(labels.Everything())
	if err != nil {
		return err
	}
	alertKey := alertNamespace + "/" + alertName
	for _, node := range nodes {
		if alertAppliedToNode(node.Annotations, alertKey) {
			key, err := cache.MetaNamespaceKeyFunc(node)
			if err == nil {
				op.nodeQueue.GetQueue().Add(key)
			}
		}
	}
	return nil
}
