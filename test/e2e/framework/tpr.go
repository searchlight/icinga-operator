package framework

import (
	tapi "github.com/appscode/searchlight/api"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func (f *Framework) EventuallyClusterAlert() GomegaAsyncAssertion {
	return Eventually(func() error {
		_, err := f.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(tapi.ResourceNameClusterAlert+"."+tapi.GroupName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// TPR group registration has 10 sec delay inside Kuberneteas api server. So, needs the extra check.
		_, err = f.extClient.ClusterAlerts(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		return err
	})
}

func (f *Framework) EventuallyNodeAlert() GomegaAsyncAssertion {
	return Eventually(func() error {
		_, err := f.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(tapi.ResourceNameNodeAlert+"."+tapi.GroupName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// TPR group registration has 10 sec delay inside Kuberneteas api server. So, needs the extra check.
		_, err = f.extClient.NodeAlerts(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		return err
	})
}

func (f *Framework) EventuallyPodAlert() GomegaAsyncAssertion {
	return Eventually(func() error {
		_, err := f.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(tapi.ResourceNamePodAlert+"."+tapi.GroupName, metav1.GetOptions{})
		if err != nil {
			return err
		}
		// TPR group registration has 10 sec delay inside Kuberneteas api server. So, needs the extra check.
		_, err = f.extClient.PodAlerts(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		return err
	})
}
