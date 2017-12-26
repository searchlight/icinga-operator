package framework

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
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

func (f *Framework) DeleteReplicaSet(obj *extensions.ReplicaSet) error {
	return f.kubeClient.ExtensionsV1beta1().ReplicaSets(obj.Namespace).Delete(obj.Name, deleteInForeground())
}

func (f *Framework) EventuallyReplicaSet(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() *core.PodList {
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
