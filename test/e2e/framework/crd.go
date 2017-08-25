package framework

import (
	tapi "github.com/appscode/searchlight/api"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func (f *Framework) EventuallyClusterAlert() GomegaAsyncAssertion {
	client := f.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	name := tapi.ResourceTypeClusterAlert + "." + tapi.V1alpha1SchemeGroupVersion.Group
	return Eventually(func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// CRD group registration has 10 sec delay inside Kuberneteas api server. So, needs the extra check.
		_, err = f.extClient.ClusterAlerts(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		return err
	})
}

func (f *Framework) EventuallyNodeAlert() GomegaAsyncAssertion {
	client := f.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	name := tapi.ResourceTypeNodeAlert + "." + tapi.V1alpha1SchemeGroupVersion.Group
	return Eventually(func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// CRD group registration has 10 sec delay inside Kuberneteas api server. So, needs the extra check.
		_, err = f.extClient.NodeAlerts(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		return err
	})
}

func (f *Framework) EventuallyPodAlert() GomegaAsyncAssertion {
	client := f.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	name := tapi.ResourceTypePodAlert + "." + tapi.V1alpha1SchemeGroupVersion.Group
	return Eventually(func() error {
		_, err := client.Get(name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// CRD group registration has 10 sec delay inside Kuberneteas api server. So, needs the extra check.
		_, err = f.extClient.PodAlerts(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		return err
	})
}
