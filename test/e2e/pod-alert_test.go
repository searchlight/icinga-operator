package e2e_test

import (
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/test/e2e/framework"
	. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var _ = Describe("PodAlert", func() {

	var (
		err   error
		f     *framework.Invocation
		rs    *extensions.ReplicaSet
		alert *tapi.PodAlert
	)

	BeforeEach(func() {
		f = root.Invoke()
		alert = f.PodAlert()
	})
	JustBeforeEach(func() {
		rs = f.ReplicaSet()
		alert.Spec.Selector = *(rs.Spec.Selector)
	})

	var (
		shouldManageIcingaServiceForOriginalReplicas = func() {
			By("Creating replica set " + rs.Name + "@" + rs.Namespace)
			_, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for Running pods")
			f.EventuallyReplicaSetRunning(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			By("Create matching pod alert")
			err := f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			By("Delete pod alerts")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForNewReplica = func() {
			By("Creating replica set " + rs.Name + "@" + rs.Namespace)
			_, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for Running pods")
			f.EventuallyReplicaSetRunning(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			By("Create matching pod alert")
			err := f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Waiting for icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			// replica = 3

			By("Delete pod alerts")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		// Reducing replica removed icinga service

		// Changing RS label removed alert
		// Changing alert selector, removed alert
	)

	Describe("Test pod_status", func() {
		AfterEach(func() {
			f.DeleteReplicaSet(rs.ObjectMeta)
			f.DeletePodAlert(alert.ObjectMeta)
		})

		Context("check_pod_status", func() {
			BeforeEach(func() {
				alert.Spec.Check = tapi.CheckPodStatus
				// vars
			})

			It("should manage icinga service for original replicas", shouldManageIcingaServiceForOriginalReplicas)
			It("should manage icinga service for new replica", shouldManageIcingaServiceForNewReplica)
		})

		Context("check_pod_exec", func() {
			BeforeEach(func() {
				alert.Spec.Check = tapi.CheckPodExec
				// vars
			})

			It("should manage icinga service for original replicas", shouldManageIcingaServiceForOriginalReplicas)
			It("should manage icinga service for new replica", shouldManageIcingaServiceForNewReplica)
		})
	})
})
