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
	"context"
	"reflect"
	"strings"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/pkg/eventer"
	"go.searchlight.dev/icinga-operator/pkg/icinga"

	"github.com/golang/glog"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	core_util "kmodules.xyz/client-go/core/v1"
	"kmodules.xyz/client-go/tools/queue"
)

func (op *Operator) initPodWatcher() {
	op.podInformer = op.kubeInformerFactory.Core().V1().Pods().Informer()
	op.podQueue = queue.New("Pod", op.MaxNumRequeues, op.NumThreads, op.reconcilePod)
	op.podInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*core.Pod)
			if pod.Status.PodIP != "" {
				klog.Warningf("Skipping pod %s/%s, since it has no IP", pod.Namespace, pod.Name)
				return
			}
			queue.Enqueue(op.podQueue.GetQueue(), obj)
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			old := oldObj.(*core.Pod)
			nu := newObj.(*core.Pod)
			if !reflect.DeepEqual(old.Labels, nu.Labels) || old.Status.PodIP != nu.Status.PodIP {
				queue.Enqueue(op.podQueue.GetQueue(), newObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			queue.Enqueue(op.podQueue.GetQueue(), obj)
		},
	})
	op.podLister = op.kubeInformerFactory.Core().V1().Pods().Lister()
}

func (op *Operator) reconcilePod(key string) error {
	obj, exists, err := op.podInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exists {
		klog.V(5).Infof("Pod %s does not exist anymore\n", key)

		namespace, name, err := cache.SplitMetaNamespaceKey(key)
		if err != nil {
			return err
		}
		return op.podHost.ForceDeleteIcingaHost(icinga.IcingaHost{
			Type:           icinga.TypePod,
			AlertNamespace: namespace,
			ObjectName:     name,
		})
	}

	klog.Infof("Sync/Add/Update for Pod %s\n", key)
	pod := obj.(*core.Pod).DeepCopy()
	err = op.ensurePod(pod)
	if err != nil {
		klog.Errorf("failed to reconcile alert for pod %s. reason: %s", key, err)
	}
	return err
}

func (op *Operator) ensurePod(pod *core.Pod) error {
	var errlist []error

	oldAlerts := sets.NewString()
	if val, ok := pod.Annotations[api.AnnotationKeyAlerts]; ok {
		names := strings.Split(val, ",")
		oldAlerts.Insert(names...)
	}

	newAlerts, err := findPodAlert(op.kubeClient, op.paLister, pod.ObjectMeta)
	if err != nil {
		return err
	}
	newNames := make([]string, len(newAlerts))
	for i := range newAlerts {
		alert := newAlerts[i]

		err = op.podHost.Apply(alert, pod)
		if err != nil {
			op.recorder.Eventf(
				alert.ObjectReference(),
				core.EventTypeWarning,
				eventer.EventReasonFailedToSync,
				`failed to  apply to pod %s/%s. Reason: %s`,
				pod.Namespace, pod.Name, err,
			)
			errlist = append(errlist, err)
		}

		newNames[i] = alert.Name
		if oldAlerts.Has(alert.Name) {
			oldAlerts.Delete(alert.Name)
		}
	}

	for _, name := range oldAlerts.List() {
		err = op.podHost.Delete(pod.Namespace, name, pod)
		if err != nil {
			if alert, e2 := op.paLister.PodAlerts(pod.Namespace).Get(name); e2 == nil {
				op.recorder.Eventf(
					alert.ObjectReference(),
					core.EventTypeWarning,
					eventer.EventReasonFailedToDelete,
					`failed to  delete for pod %s/%s. Reason: %s`,
					pod.Namespace, pod.Name, err,
				)
			}
			errlist = append(errlist, err)
		}
	}

	_, _, err = core_util.PatchPod(context.TODO(), op.kubeClient, pod, func(in *core.Pod) *core.Pod {
		if in.Annotations == nil {
			in.Annotations = make(map[string]string, 0)
		}
		if len(newNames) > 0 {
			in.Annotations[api.AnnotationKeyAlerts] = strings.Join(newNames, ",")
		} else {
			delete(in.Annotations, api.AnnotationKeyAlerts)
		}
		return in
	}, metav1.PatchOptions{})
	if err != nil {
		errlist = append(errlist, err)
	}
	return utilerrors.NewAggregate(errlist)
}
