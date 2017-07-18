package framework

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	. "github.com/onsi/gomega"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
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

func (f *Framework) CreateDeploymentApp(obj *apps.Deployment) error {
	_, err := f.kubeClient.AppsV1beta1().Deployments(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) UpdateDeploymentApp(obj *apps.Deployment) (*apps.Deployment, error) {
	return f.kubeClient.AppsV1beta1().Deployments(obj.Namespace).Update(obj)
}

func (f *Framework) EventuallyDeleteDeploymentApp(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	deployment, err := f.kubeClient.AppsV1beta1().Deployments(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return Eventually(func() bool { return true })
	}
	deployment.Spec.Replicas = types.Int32P(0)
	deployment, err = f.UpdateDeploymentApp(deployment)
	Expect(err).NotTo(HaveOccurred())

	return Eventually(
		func() bool {
			podList, err := f.GetPodList(deployment)
			Expect(err).NotTo(HaveOccurred())
			if len(podList.Items) != 0 {
				return false
			}

			err = f.kubeClient.AppsV1beta1().Deployments(meta.Namespace).Delete(meta.Name, deleteInForeground())
			Expect(err).NotTo(HaveOccurred())
			return true
		},
		time.Minute*5,
		time.Second*5,
	)
}

func (f *Framework) EventuallyDeploymentApp(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(func() *apiv1.PodList {
		obj, err := f.kubeClient.AppsV1beta1().Deployments(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		Expect(err).NotTo(HaveOccurred())
		podList, err := f.GetPodList(obj)
		Expect(err).NotTo(HaveOccurred())
		return podList
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

func (f *Framework) CreateDeploymentExtension(obj *extensions.Deployment) error {
	_, err := f.kubeClient.ExtensionsV1beta1().Deployments(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) UpdateDeploymentExtension(obj *extensions.Deployment) (*extensions.Deployment, error) {
	return f.kubeClient.ExtensionsV1beta1().Deployments(obj.Namespace).Update(obj)
}

func (f *Framework) EventuallyDeleteDeploymentExtension(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	deployment, err := f.kubeClient.ExtensionsV1beta1().Deployments(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		return Eventually(func() bool { return true })
	}
	deployment.Spec.Replicas = types.Int32P(0)
	deployment, err = f.UpdateDeploymentExtension(deployment)
	Expect(err).NotTo(HaveOccurred())

	return Eventually(
		func() bool {
			podList, err := f.GetPodList(deployment)
			Expect(err).NotTo(HaveOccurred())
			if len(podList.Items) != 0 {
				return false
			}

			err = f.kubeClient.ExtensionsV1beta1().Deployments(meta.Namespace).Delete(meta.Name, deleteInForeground())
			Expect(err).NotTo(HaveOccurred())
			return true
		},
		time.Minute*5,
		time.Second*5,
	)
}

func (f *Framework) EventuallyDeploymentExtension(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() *apiv1.PodList {
			obj, err := f.kubeClient.ExtensionsV1beta1().Deployments(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			podList, err := f.GetPodList(obj)
			Expect(err).NotTo(HaveOccurred())
			return podList
		},
		time.Minute*5,
		time.Second*5,
	)
}
