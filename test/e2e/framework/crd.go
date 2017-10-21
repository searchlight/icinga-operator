package framework

import (
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Framework) EventuallyClusterAlert() GomegaAsyncAssertion {
	client := f.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	name := api.ResourceTypeClusterAlert + "." + api.SchemeGroupVersion.Group
	return Eventually(func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// CRD group registration has 10 sec delay inside Kuberneteas api server. So, needs the extra check.
		_, err = f.extClient.ClusterAlerts(core.NamespaceDefault).List(metav1.ListOptions{})
		return err
	})
}

func (f *Framework) EventuallyNodeAlert() GomegaAsyncAssertion {
	client := f.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	name := api.ResourceTypeNodeAlert + "." + api.SchemeGroupVersion.Group
	return Eventually(func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// CRD group registration has 10 sec delay inside Kuberneteas api server. So, needs the extra check.
		_, err = f.extClient.NodeAlerts(core.NamespaceDefault).List(metav1.ListOptions{})
		return err
	})
}

func (f *Framework) EventuallyPodAlert() GomegaAsyncAssertion {
	client := f.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	name := api.ResourceTypePodAlert + "." + api.SchemeGroupVersion.Group
	return Eventually(func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// CRD group registration has 10 sec delay inside Kuberneteas api server. So, needs the extra check.
		_, err = f.extClient.PodAlerts(core.NamespaceDefault).List(metav1.ListOptions{})
		return err
	})
}
