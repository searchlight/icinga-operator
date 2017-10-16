package framework

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	kutilext "github.com/appscode/kutil/extensions/v1beta1"
	. "github.com/onsi/gomega"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Invocation) ReplicaSet() *extensions.ReplicaSet {
	return &extensions.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("replicaset"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: extensions.ReplicaSetSpec{
			Replicas: types.Int32P(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": f.app,
				},
			},
			Template: f.PodTemplate(),
		},
	}
}

func (f *Framework) GetReplicaSet(meta metav1.ObjectMeta) (*extensions.ReplicaSet, error) {
	return f.kubeClient.ExtensionsV1beta1().ReplicaSets(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) CreateReplicaSet(obj *extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	return f.kubeClient.ExtensionsV1beta1().ReplicaSets(obj.Namespace).Create(obj)
}

func (f *Framework) TryPatchReplicaSet(meta metav1.ObjectMeta, transformer func(*extensions.ReplicaSet) *extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	return kutilext.TryPatchReplicaSet(f.kubeClient, meta, transformer)
}

func (f *Framework) EventuallyDeleteReplicaSet(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	_, err := f.kubeClient.ExtensionsV1beta1().ReplicaSets(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return Eventually(func() bool { return true })
	}

	rs, err := f.TryPatchReplicaSet(meta, func(in *extensions.ReplicaSet) *extensions.ReplicaSet {
		in.Spec.Replicas = types.Int32P(0)
		return in
	})
	Expect(err).NotTo(HaveOccurred())

	return Eventually(
		func() bool {
			podList, err := f.GetPodList(rs)
			Expect(err).NotTo(HaveOccurred())
			if len(podList.Items) != 0 {
				return false
			}

			err = f.kubeClient.ExtensionsV1beta1().ReplicaSets(meta.Namespace).Delete(meta.Name, deleteInForeground())
			Expect(err).NotTo(HaveOccurred())
			return true
		},
		time.Minute*5,
		time.Second*5,
	)
}

func (f *Framework) EventuallyReplicaSet(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() *apiv1.PodList {
			obj, err := f.kubeClient.ExtensionsV1beta1().ReplicaSets(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			podList, err := f.GetPodList(obj)
			Expect(err).NotTo(HaveOccurred())
			return podList
		},
		time.Minute*5,
		time.Second*5,
	)
}
