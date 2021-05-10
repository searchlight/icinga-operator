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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/tools/queue"
)

func (op *Operator) initPodAlertWatcher() {
	op.paInformer = op.monInformerFactory.Monitoring().V1alpha1().PodAlerts().Informer()
	op.paQueue = queue.New("PodAlert", op.MaxNumRequeues, op.NumThreads, op.reconcilePodAlert)
	op.paInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			alert := obj.(*api.PodAlert)
			if op.isValid(alert) {
				queue.Enqueue(op.paQueue.GetQueue(), obj)
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*api.PodAlert)
			nu := newObj.(*api.PodAlert)

			if reflect.DeepEqual(old.Spec, nu.Spec) {
				return
			}
			if op.isValid(nu) {
				queue.Enqueue(op.paQueue.GetQueue(), nu)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.paQueue.GetQueue(), obj)
		},
	})
	op.paLister = op.monInformerFactory.Monitoring().V1alpha1().PodAlerts().Lister()
}

func (op *Operator) reconcilePodAlert(key string) error {
	obj, exists, err := op.paInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		klog.Warningf("PodAlert %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		return op.ensurePodAlertDeleted(namespace, name)
	}

	alert := obj.(*api.PodAlert).DeepCopy()
	klog.Infof("Sync/Add/Update for PodAlert %s\n", alert.GetName())

	op.ensurePodAlert(alert)
	op.ensurePodAlertDeleted(alert.Namespace, alert.Name)
	return nil
}

func (op *Operator) ensurePodAlert(alert *api.PodAlert) error {
	if alert.Spec.PodName != nil {
		pod, err := op.podLister.Pods(alert.Namespace).Get(*alert.Spec.PodName)
		if err != nil {
			return err
		}
		key, err := cache.MetaNamespaceKeyFunc(pod)
		if err == nil {
			op.podQueue.GetQueue().Add(key)
		}
		return nil
	}

	sel, err := metav1.LabelSelectorAsSelector(alert.Spec.Selector)
	if err != nil {
		return err
	}
	pods, err := op.podLister.Pods(alert.Namespace).List(sel)
	if err != nil {
		return err
	}
	for i := range pods {
		pod := pods[i]
		key, err := cache.MetaNamespaceKeyFunc(pod)
		if err == nil {
			op.podQueue.GetQueue().Add(key)
		}
	}
	return nil
}

func alertAppliedToPod(a map[string]string, key string) bool {
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

func (op *Operator) ensurePodAlertDeleted(alertNamespace, alertName string) error {
	pods, err := op.podLister.Pods(alertNamespace).List(labels.Everything())
	if err != nil {
		return err
	}
	for _, pod := range pods {
		if alertAppliedToPod(pod.Annotations, alertName) {
			key, err := cache.MetaNamespaceKeyFunc(pod)
			if err == nil {
				op.podQueue.GetQueue().Add(key)
			}
		}
	}
	return nil
}
