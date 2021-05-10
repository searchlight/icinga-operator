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
	"time"

	. "github.com/onsi/gomega"
	"gomodules.xyz/pointer"
	"gomodules.xyz/x/crypto/rand"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Invocation) StatefulSet() *apps.StatefulSet {
	ss := &apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("statefulset"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: apps.StatefulSetSpec{
			Replicas: pointer.Int32P(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": f.app,
				},
			},
			Template:    f.PodTemplate(),
			ServiceName: TEST_HEADLESS_SERVICE,
		},
	}

	ss.Spec.Template.Spec.Volumes = []core.Volume{}
	ss.Spec.VolumeClaimTemplates = []core.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: TestSourceDataVolumeName,
				Annotations: map[string]string{
					"volume.beta.kubernetes.io/storage-class": f.storageClass,
				},
			},
			Spec: core.PersistentVolumeClaimSpec{
				StorageClassName: pointer.StringP(f.storageClass),
				AccessModes: []core.PersistentVolumeAccessMode{
					core.ReadWriteOnce,
				},
				Resources: core.ResourceRequirements{
					Requests: core.ResourceList{
						core.ResourceStorage: resource.MustParse("5Gi"),
					},
				},
			},
		},
	}
	return ss
}

func (f *Framework) GetStatefulSet(meta metav1.ObjectMeta) (*apps.StatefulSet, error) {
	return f.kubeClient.AppsV1beta1().StatefulSets(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
}

func (f *Framework) CreateStatefulSet(obj *apps.StatefulSet) (*apps.StatefulSet, error) {
	return f.kubeClient.AppsV1beta1().StatefulSets(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
}

func (f *Framework) DeleteStatefulSet(obj *apps.StatefulSet) error {
	return f.kubeClient.AppsV1beta1().StatefulSets(obj.Namespace).Delete(context.TODO(), obj.Name, *deleteInForeground())
}

func (f *Framework) EventuallyStatefulSet(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() *core.PodList {
			obj, err := f.kubeClient.AppsV1beta1().StatefulSets(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			podList, err := f.GetPodList(obj)
			Expect(err).NotTo(HaveOccurred())
			return podList
		},
		time.Minute*5,
		time.Second*5,
	)
}
