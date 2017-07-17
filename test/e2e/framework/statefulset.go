package framework

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
)

func (f *Invocation) StatefulSet() *apps.StatefulSet {
	resource := &apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("statefulset"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: apps.StatefulSetSpec{
			Replicas:    types.Int32P(1),
			Template:    f.PodTemplate(),
			ServiceName: TEST_HEADLESS_SERVICE,
		},
	}
	return resource
}

func (f *Framework) GetStatefulSet(meta metav1.ObjectMeta) (*apps.StatefulSet, error) {
	return f.kubeClient.AppsV1beta1().StatefulSets(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) CreateStatefulSet(obj *apps.StatefulSet) (*apps.StatefulSet, error) {
	return f.kubeClient.AppsV1beta1().StatefulSets(obj.Namespace).Create(obj)
}

func (f *Framework) UpdateStatefulSet(obj *apps.StatefulSet) (*apps.StatefulSet, error) {
	return f.kubeClient.AppsV1beta1().StatefulSets(obj.Namespace).Update(obj)
}

func (f *Framework) DeleteStatefulSet(meta metav1.ObjectMeta) error {
	return f.kubeClient.AppsV1beta1().StatefulSets(meta.Namespace).Delete(meta.Name, deleteInBackground())
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
		time.Minute*2,
		time.Second*2,
	)
}
