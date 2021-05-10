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

package framework

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/gomega"
	"gomodules.xyz/x/crypto/rand"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	TestSourceDataVolumeName = "source-data"
	TestSourceDataMountPath  = "/source/data"
)

func (f *Invocation) Pod() *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("pod"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: getPodSpec(),
	}
}

func (f *Framework) CreatePod(obj *core.Pod) (*core.Pod, error) {
	return f.kubeClient.CoreV1().Pods(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (f *Framework) DeletePod(meta metav1.ObjectMeta) error {
	return f.kubeClient.CoreV1().Pods(meta.Namespace).Delete(context.TODO(), meta.Name, *deleteInForeground())
}

func (f *Framework) EventuallyPodRunning(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() *core.PodList {
			obj, err := f.kubeClient.CoreV1().Pods(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			podList, err := f.GetPodList(obj)
			Expect(err).NotTo(HaveOccurred())
			return podList
		},
		time.Minute*5,
		time.Second*5,
	)
}

func (f *Invocation) PodTemplate() core.PodTemplateSpec {
	return core.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: getPodSpec(),
	}
}

func getPodSpec() core.PodSpec {
	return core.PodSpec{
		Containers: []core.Container{
			{
				Name:            "busybox",
				Image:           "busybox",
				ImagePullPolicy: core.PullIfNotPresent,
				Command: []string{
					"sleep",
					"1d",
				},
				VolumeMounts: []core.VolumeMount{
					{
						Name:      TestSourceDataVolumeName,
						MountPath: TestSourceDataMountPath,
					},
				},
			},
		},
		Volumes: []core.Volume{
			{
				Name: TestSourceDataVolumeName,
				VolumeSource: core.VolumeSource{
					EmptyDir: &core.EmptyDirVolumeSource{},
				},
			},
		},
	}
}

func (f *Framework) GetPodList(actual interface{}) (*core.PodList, error) {
	switch obj := actual.(type) {
	case *core.Pod:
		return f.listPods(obj.Namespace, obj.Labels)
	case *extensions.ReplicaSet:
		return f.listPods(obj.Namespace, obj.Spec.Selector.MatchLabels)
	case *extensions.Deployment:
		return f.listPods(obj.Namespace, obj.Spec.Selector.MatchLabels)
	case *apps.Deployment:
		return f.listPods(obj.Namespace, obj.Spec.Selector.MatchLabels)
	case *apps.StatefulSet:
		return f.listPods(obj.Namespace, obj.Spec.Selector.MatchLabels)
	default:
		return nil, fmt.Errorf("Unknown object type")
	}
}

func (f *Framework) listPods(namespace string, label map[string]string) (*core.PodList, error) {
	return f.kubeClient.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(label).String(),
	})
}
