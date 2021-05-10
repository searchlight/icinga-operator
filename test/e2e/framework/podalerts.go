/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"context"
	"time"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/test/e2e/matcher"

	. "github.com/onsi/gomega"
	"gomodules.xyz/x/crypto/rand"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Invocation) PodAlert() *api.PodAlert {
	return &api.PodAlert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("podalert"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: api.PodAlertSpec{
			CheckInterval: metav1.Duration{Duration: time.Second * 5},
			Vars:          make(map[string]string),
		},
	}
}

func (f *Framework) CreatePodAlert(obj *api.PodAlert) error {
	_, err := f.extClient.MonitoringV1alpha1().PodAlerts(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

func (f *Framework) GetPodAlert(meta metav1.ObjectMeta) (*api.PodAlert, error) {
	return f.extClient.MonitoringV1alpha1().PodAlerts(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
}

func (f *Framework) DeletePodAlert(meta metav1.ObjectMeta) error {
	return f.extClient.MonitoringV1alpha1().PodAlerts(meta.Namespace).Delete(context.TODO(), meta.Name, metav1.DeleteOptions{})
}

func (f *Framework) getPodAlertObjects(meta metav1.ObjectMeta, podAlertSpec api.PodAlertSpec) ([]icinga.IcingaHost, error) {
	names := make([]string, 0)

	if podAlertSpec.PodName != nil {
		pod, err := f.kubeClient.CoreV1().Pods(meta.Namespace).Get(context.TODO(), *podAlertSpec.PodName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		names = append(names, pod.Name)
	} else if podAlertSpec.Selector != nil {
		sel, err := metav1.LabelSelectorAsSelector(podAlertSpec.Selector)
		if err != nil {
			return nil, err
		}
		podList, err := f.kubeClient.CoreV1().Pods(meta.Namespace).List(
			context.TODO(),
			metav1.ListOptions{
				LabelSelector: sel.String(),
			},
		)
		if err != nil {
			return nil, err
		}

		for _, pod := range podList.Items {
			names = append(names, pod.Name)
		}

	}

	objectList := make([]icinga.IcingaHost, 0)
	for _, name := range names {
		objectList = append(objectList,
			icinga.IcingaHost{
				Type:           icinga.TypePod,
				AlertNamespace: meta.Namespace,
				ObjectName:     name,
			})
	}

	return objectList, nil
}

func (f *Framework) EventuallyPodAlertIcingaService(meta metav1.ObjectMeta, podAlertSpec api.PodAlertSpec) GomegaAsyncAssertion {
	objectList, err := f.getPodAlertObjects(meta, podAlertSpec)
	Expect(err).NotTo(HaveOccurred())

	in := icinga.NewPodHost(f.icingaClient, "6").
		IcingaServiceSearchQuery(meta.Name, objectList...)

	return Eventually(
		func() matcher.IcingaServiceState {
			var respService icinga.ResponseObject
			status, err := f.icingaClient.Service("").Get([]string{}, in).Do().Into(&respService)
			if status == 0 {
				return matcher.IcingaServiceState{Unknown: 1.0}
			}
			Expect(err).NotTo(HaveOccurred())

			var icingaServiceState matcher.IcingaServiceState
			for _, service := range respService.Results {
				if service.Attrs.LastState == 0.0 {
					icingaServiceState.OK++
				}
				if service.Attrs.LastState == 1.0 {
					icingaServiceState.Warning++
				}
				if service.Attrs.LastState == 2.0 {
					icingaServiceState.Critical++
				}
				if service.Attrs.LastState == 3.0 {
					icingaServiceState.Unknown++
				}
			}
			return icingaServiceState
		},
		time.Minute*5,
		time.Second*5,
	)
}

func (f *Framework) CleanPodAlert() {
	caList, err := f.extClient.MonitoringV1alpha1().PodAlerts(f.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}
	for _, e := range caList.Items {
		f.extClient.MonitoringV1alpha1().PodAlerts(f.namespace).Delete(context.TODO(), e.Name, metav1.DeleteOptions{})
	}
}
