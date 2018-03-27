package e2e

import (
	"fmt"
	"net/http/httptest"
	"net/url"

	"github.com/appscode/go/types"
	kutil_ext "github.com/appscode/kutil/extensions/v1beta1"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins/notifier"
	"github.com/appscode/searchlight/test/e2e/framework"
	. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core_v1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
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
				rs.Spec.Replicas = types.Int32P(*rs.Spec.Replicas - 1)

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
				sms.Author = "e2e"
				By("Check received notification message")
				f.EventuallyHTTPServerResponse(serverURL).Should(BeIdenticalTo(sms.Render()))

				By("Acknowledge notification")
				f.AcknowledgeClusterAlertNotification(clusterAlert.ObjectMeta, hostname)

				sms.NotificationType = string(api.NotificationAcknowledgement)
				By("Check received notification message")
				f.EventuallyHTTPServerResponse(serverURL).Should(BeIdenticalTo(sms.Render()))

				By("Patch ReplicaSet to increate replicas")
				rs, _, err = kutil_ext.PatchReplicaSet(f.KubeClient(), rs, func(set *extensions.ReplicaSet) *extensions.ReplicaSet {
					set.Spec.Replicas = types.Int32P(*rs.Spec.Replicas + 1)
					return set
				})

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
