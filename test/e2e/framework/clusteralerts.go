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

func (f *Invocation) ClusterAlert() *api.ClusterAlert {
	return &api.ClusterAlert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("clusteralert"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: api.ClusterAlertSpec{
			CheckInterval: metav1.Duration{Duration: time.Second * 5},
			AlertInterval: metav1.Duration{Duration: time.Minute * 5},
			Vars:          make(map[string]string),
		},
	}
}

func (f *Framework) CreateClusterAlert(obj *api.ClusterAlert) error {
	_, err := f.extClient.MonitoringV1alpha1().ClusterAlerts(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

func (f *Framework) GetClusterAlert(meta metav1.ObjectMeta) (*api.ClusterAlert, error) {
	return f.extClient.MonitoringV1alpha1().ClusterAlerts(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
}

func (f *Framework) DeleteClusterAlert(meta metav1.ObjectMeta) error {
	return f.extClient.MonitoringV1alpha1().ClusterAlerts(meta.Namespace).Delete(context.TODO(), meta.Name, metav1.DeleteOptions{})
}

func (f *Framework) getClusterAlertObjects(meta metav1.ObjectMeta) icinga.IcingaHost {
	return icinga.IcingaHost{
		Type:           icinga.TypeCluster,
		AlertNamespace: meta.Namespace,
	}
}

func (f *Framework) EventuallyClusterAlertIcingaService(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	icingaHost := f.getClusterAlertObjects(meta)

	in := icinga.NewClusterHost(f.icingaClient, "6").
		IcingaServiceSearchQuery(meta.Name, icingaHost)

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

type notificationObject struct {
	Results []struct {
		Attrs struct {
			NotificationNumber float64 `json:"notification_number"`
		} `json:"attrs"`
	} `json:"results"`
}

func (f *Framework) EventuallyClusterAlertIcingaNotification(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	icingaHost := f.getClusterAlertObjects(meta)
	host, err := icingaHost.Name()
	Expect(err).NotTo(HaveOccurred())

	return Eventually(
		func() float64 {
			var resp notificationObject
			status, err := f.icingaClient.Notifications(host).Get([]string{meta.GetName(), meta.GetName()}, "").Do().Into(&resp)
			if status == 0 {
				return -1
			}
			Expect(err).NotTo(HaveOccurred())

			if len(resp.Results) != 1 {
				return -1
			}
			return resp.Results[0].Attrs.NotificationNumber
		},
		time.Minute*5,
		time.Second*5,
	)
}

func (f *Framework) CleanClusterAlert() {
	caList, err := f.extClient.MonitoringV1alpha1().ClusterAlerts(f.namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return
	}
	for _, e := range caList.Items {
		f.extClient.MonitoringV1alpha1().ClusterAlerts(f.namespace).Delete(context.TODO(), e.Name, metav1.DeleteOptions{})
	}
}
