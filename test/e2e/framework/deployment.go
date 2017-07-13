package framework

import (
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func (f *Invocation) DeploymentApp() *apps.Deployment {
	return &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("deployment"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: types.Int32P(1),
			Template: f.PodTemplate(),
		},
	}
}

func (f *Framework) GetDeploymentApp(meta metav1.ObjectMeta) (*apps.Deployment, error) {
	return f.kubeClient.AppsV1beta1().Deployments(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) CreateDeploymentApp(obj *apps.Deployment) error {
	_, err := f.kubeClient.AppsV1beta1().Deployments(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) DeleteDeploymentApp(meta metav1.ObjectMeta) error {
	return f.kubeClient.AppsV1beta1().Deployments(meta.Namespace).Delete(meta.Name, deleteInForeground())
}

func (f *Framework) EventuallyDeploymentApp(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(func() *apps.Deployment {
		obj, err := f.kubeClient.AppsV1beta1().Deployments(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		return obj
	})
}

func (f *Invocation) DeploymentExtension() *extensions.Deployment {
	return &extensions.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("deployment"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: extensions.DeploymentSpec{
			Replicas: types.Int32P(1),
			Template: f.PodTemplate(),
		},
	}
}

func (f *Framework) GetDeploymentExtension(meta metav1.ObjectMeta) (*extensions.Deployment, error) {
	return f.kubeClient.ExtensionsV1beta1().Deployments(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) CreateDeploymentExtension(obj *extensions.Deployment) error {
	_, err := f.kubeClient.ExtensionsV1beta1().Deployments(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) DeleteDeploymentExtension(meta metav1.ObjectMeta) error {
	return f.kubeClient.ExtensionsV1beta1().Deployments(meta.Namespace).Delete(meta.Name, deleteInForeground())
}

func (f *Framework) EventuallyDeploymentExtension(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(func() *extensions.Deployment {
		obj, err := f.kubeClient.ExtensionsV1beta1().Deployments(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		return obj
	})
}
