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

package e2e

import (
	"context"
	"fmt"
	"net/http/httptest"
	"net/url"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins/notifier"
	"go.searchlight.dev/icinga-operator/test/e2e/framework"
	. "go.searchlight.dev/icinga-operator/test/e2e/matcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gomodules.xyz/pointer"
	core_v1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kutil_ext "kmodules.xyz/client-go/extensions/v1beta1"
)

var _ = Describe("notification", func() {
	var (
		f            *framework.Invocation
		rs           *extensions.ReplicaSet
		clusterAlert *api.ClusterAlert
		secret       *core_v1.Secret
		server       *httptest.Server
		serverURL    string
		webhookURL   string
		icingaHost   *icinga.IcingaHost
		hostname     string
	)

	BeforeEach(func() {
		if root.Provider != "minikube" {
			Skip("notification test is only allowed in minikube")
		}

		f = root.Invoke()
		rs = f.ReplicaSet()
		clusterAlert = f.ClusterAlert()
		secret = f.GetWebHookSecret()
		server = framework.StartServer()
		serverURL = server.URL
		url, _ := url.Parse(serverURL)
		webhookURL = fmt.Sprintf("http://10.0.2.2:%s", url.Port())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Test", func() {
		Context("check notification", func() {
			BeforeEach(func() {
				rs.Spec.Replicas = pointer.Int32P(*rs.Spec.Replicas - 1)

				secret.StringData["WEBHOOK_URL"] = webhookURL
				secret.StringData["WEBHOOK_TO"] = "test"

				clusterAlert.Spec.Check = api.CheckPodExists
				clusterAlert.Spec.Vars["count"] = fmt.Sprintf("%v", *rs.Spec.Replicas+1)
				clusterAlert.Spec.NotifierSecretName = secret.Name
				clusterAlert.Spec.Receivers = []api.Receiver{
					{
						State:    "Critical",
						To:       []string{"shahriar"},
						Notifier: "webhook",
					},
				}

				icingaHost = &icinga.IcingaHost{
					Type:           icinga.TypeCluster,
					AlertNamespace: clusterAlert.Namespace,
				}
				hostname, _ = icingaHost.Name()
			})
			AfterEach(func() {
				server.Close()
			})
			It("with webhook receiver", func() {
				By("Create notifier secret: " + secret.Name)
				err := f.CreateWebHookSecret(secret)
				Expect(err).NotTo(HaveOccurred())

				By("Create ReplicaSet: " + rs.Name)
				rs, err = f.CreateReplicaSet(rs)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for Running pods")
				f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

				By("Create cluster alert: " + clusterAlert.Name)
				err = f.CreateClusterAlert(clusterAlert)
				Expect(err).NotTo(HaveOccurred())

				By("Check icinga services")
				f.EventuallyClusterAlertIcingaService(clusterAlert.ObjectMeta).
					Should(HaveIcingaObject(IcingaServiceState{Critical: 1}))

				By("Force check now")
				f.ForceCheckClusterAlert(clusterAlert.ObjectMeta, hostname, 5)

				By("Count icinga notification")
				f.EventuallyClusterAlertIcingaNotification(clusterAlert.ObjectMeta).Should(BeNumerically(">", 0.0))

				hostname, err := icingaHost.Name()
				Expect(err).NotTo(HaveOccurred())
				sms := &notifier.SMS{
					AlertName:        clusterAlert.Name,
					Hostname:         hostname,
					ServiceState:     "Critical",
					NotificationType: string(api.NotificationProblem),
				}
				By("Check received notification message")
				f.EventuallyHTTPServerResponse(serverURL).Should(BeIdenticalTo(sms.Render()))

				By("Send custom notification")
				f.SendClusterAlertCustomNotification(clusterAlert.ObjectMeta, hostname)

				sms.NotificationType = string(api.NotificationCustom)
				sms.Comment = "test"
				// Used in regular expression to match any author
				sms.Author = "(.*)"

				By("Check received notification message")
				f.EventuallyHTTPServerResponse(serverURL).Should(ReceiveNotificationWithExp(sms.Render()))

				By("Acknowledge notification")
				f.AcknowledgeClusterAlertNotification(clusterAlert.ObjectMeta, hostname)

				sms.NotificationType = string(api.NotificationAcknowledgement)
				By("Check received notification message")
				f.EventuallyHTTPServerResponse(serverURL).Should(ReceiveNotificationWithExp(sms.Render()))

				By("Patch ReplicaSet to increate replicas")
				rs, _, err = kutil_ext.PatchReplicaSet(context.TODO(), f.KubeClient(), rs, func(set *extensions.ReplicaSet) *extensions.ReplicaSet {
					set.Spec.Replicas = pointer.Int32P(*rs.Spec.Replicas + 1)
					return set
				}, metav1.PatchOptions{})

				By("Check icinga services")
				f.EventuallyClusterAlertIcingaService(clusterAlert.ObjectMeta).
					Should(HaveIcingaObject(IcingaServiceState{OK: 1}))

				sms.Comment = ""
				sms.Author = ""
				sms.NotificationType = string(api.NotificationRecovery)

				By("Check received notification message")
				f.EventuallyHTTPServerResponse(serverURL).Should(BeIdenticalTo(sms.Render()))
			})
		})
	})
})
