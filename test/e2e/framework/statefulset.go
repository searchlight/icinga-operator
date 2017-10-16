package framework

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	kutilapps "github.com/appscode/kutil/apps/v1beta1"
	. "github.com/onsi/gomega"
	apps "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
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
			Replicas: types.Int32P(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": f.app,
				},
			},
			Template:    f.PodTemplate(),
			ServiceName: TEST_HEADLESS_SERVICE,
		},
	}

	ss.Spec.Template.Spec.Volumes = []apiv1.Volume{}
	ss.Spec.VolumeClaimTemplates = []apiv1.PersistentVolumeClaim{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: TestSourceDataVolumeName,
				Annotations: map[string]string{
					"volume.beta.kubernetes.io/storage-class": f.storageClass,
				},
			},
			Spec: apiv1.PersistentVolumeClaimSpec{
				StorageClassName: types.StringP(f.storageClass),
				AccessModes: []apiv1.PersistentVolumeAccessMode{
					apiv1.ReadWriteOnce,
				},
				Resources: apiv1.ResourceRequirements{
					Requests: apiv1.ResourceList{
						apiv1.ResourceStorage: resource.MustParse("5Gi"),
					},
				},
			},
		},
	}
	return ss
}

func (f *Framework) GetStatefulSet(meta metav1.ObjectMeta) (*apps.StatefulSet, error) {
	return f.kubeClient.AppsV1beta1().StatefulSets(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) CreateStatefulSet(obj *apps.StatefulSet) (*apps.StatefulSet, error) {
	return f.kubeClient.AppsV1beta1().StatefulSets(obj.Namespace).Create(obj)
}

func (f *Framework) TryPatchStatefulSet(meta metav1.ObjectMeta, transformer func(*apps.StatefulSet) *apps.StatefulSet) (*apps.StatefulSet, error) {
	return kutilapps.TryPatchStatefulSet(f.kubeClient, meta, transformer)
}

func (f *Framework) EventuallyDeleteStatefulSet(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	ss, err := f.TryPatchStatefulSet(meta, func(in *apps.StatefulSet) *apps.StatefulSet {
		in.Spec.Replicas = types.Int32P(0)
		return in
	})
	if kerr.IsNotFound(err) {
		return Eventually(func() bool { return true })
	}
	Expect(err).NotTo(HaveOccurred())

	return Eventually(
		func() bool {
			podList, err := f.GetPodList(ss)
			Expect(err).NotTo(HaveOccurred())
			if len(podList.Items) != 0 {
				return false
			}

			err = f.kubeClient.AppsV1beta1().StatefulSets(meta.Namespace).Delete(meta.Name, deleteInBackground())
			Expect(err).NotTo(HaveOccurred())
			return true
		},
		time.Minute*5,
		time.Second*5,
	)
}

func (f *Framework) EventuallyStatefulSet(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() *apiv1.PodList {
			obj, err := f.kubeClient.AppsV1beta1().StatefulSets(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			podList, err := f.GetPodList(obj)
			Expect(err).NotTo(HaveOccurred())
			return podList
		},
		time.Minute*5,
		time.Second*5,
	)
}
