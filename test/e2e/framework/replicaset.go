package framework

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
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

func (f *Framework) UpdateReplicaSet(obj *extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	return f.kubeClient.ExtensionsV1beta1().ReplicaSets(obj.Namespace).Update(obj)
}

func (f *Framework) EventuallyDeleteReplicaSet(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	rs, err := f.kubeClient.ExtensionsV1beta1().ReplicaSets(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return Eventually(func() bool { return true })
	}
	rs.Spec.Replicas = types.Int32P(0)
	rs, err = f.UpdateReplicaSet(rs)
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
