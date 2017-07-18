package e2e_test

import (
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/test/e2e/framework"
	. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var _ = Describe("PodAlert", func() {

	var (
		err   error
		f     *framework.Invocation
		rs    *extensions.ReplicaSet
		pod   *apiv1.Pod
		alert *tapi.PodAlert
	)

	BeforeEach(func() {
		f = root.Invoke()
		rs = f.ReplicaSet()
		pod = f.Pod()
		alert = f.PodAlert()
	})

	var (
		shouldManageIcingaServiceForLabelSelector = func() {
			By("Create ReplicaSet " + rs.Name + "@" + rs.Namespace)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching pod alert")
			err := f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			By("Delete podalerts")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForNewPod = func() {
			By("Create ReplicaSet " + rs.Name + "@" + rs.Namespace)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching pod alert")
			err := f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			rs, err = f.GetReplicaSet(rs.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Increase replica")
			rs.Spec.Replicas = types.Int32P(3)
			rs, err = f.UpdateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			By("Delete podalerts")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForDeletedPod = func() {
			By("Create ReplicaSet " + rs.Name + "@" + rs.Namespace)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching pod alert")
			err := f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			rs, err = f.GetReplicaSet(rs.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Decreate replica")
			rs.Spec.Replicas = types.Int32P(1)
			rs, err = f.UpdateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			By("Delete podalerts")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForLabelChanged = func() {
			By("Create ReplicaSet " + rs.Name + "@" + rs.Namespace)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching pod alert")
			err := f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			alert, err = f.GetPodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			oldAlertSpec := alert.Spec

			By("Change LabelSelector")
			alert.Spec.Selector.MatchLabels = map[string]string{
				"app": rand.WithUniqSuffix("searchlight-e2e"),
			}

			alert, err = f.UpdatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, oldAlertSpec).
				Should(HaveIcingaObject(IcingaServiceState{}))
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))

			By("Delete podalerts")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())
		}

		shouldManageIcingaServiceForPodName = func() {
			By("Create Pod " + pod.Name + "@" + pod.Namespace)
			pod, err = f.CreatePod(pod)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyPodRunning(pod.ObjectMeta).Should(HaveRunningPods(1))

			alert.Spec.PodName = pod.Name

			By("Create matching pod alert")
			err := f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: 1}))

			By("Delete podalerts")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldHandleIcingaServiceForCriticalState = func() {
			rs.Spec.Template.Spec.Containers[0].Image = "invalid-image"
			By("Create ReplicaSet " + rs.Name + "@" + rs.Namespace)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for all pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HavePods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching pod alert")
			err := f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Critical: *rs.Spec.Replicas}))

			By("Delete podalerts")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}
	)

	Describe("Test", func() {
		AfterEach(func() {
			f.DeleteReplicaSet(rs.ObjectMeta)
			f.DeletePod(pod.ObjectMeta)
		})

		// Check "pod_status" and basic searchlight functionality
		Context("check_pod_status", func() {
			JustBeforeEach(func() {
				alert.Spec.Check = tapi.CheckPodStatus
			})

			It("should manage icinga service for Alert.Spec.Selector", shouldManageIcingaServiceForLabelSelector)
			It("should manage icinga service for new Pod", shouldManageIcingaServiceForNewPod)
			It("should manage icinga service for deleted Pod", shouldManageIcingaServiceForDeletedPod)
			It("should manage icinga service for Alert.Spec.Selector changed", shouldManageIcingaServiceForLabelChanged)
			It("should manage icinga service for Alert.Spec.PodName", shouldManageIcingaServiceForPodName)
			It("should handle icinga service for Critical State", shouldHandleIcingaServiceForCriticalState)
		})

		// Check "volume"
		Context("check_volume", func() {
			JustBeforeEach(func() {
				alert.Spec.Check = tapi.CheckVolume
				alert.Spec.Vars = map[string]interface{}{
					"volume_name": framework.TestSourceDataVolumeName,
				}
			})

			It("should manage icinga service for Ok State", shouldManageIcingaServiceForLabelSelector)
		})

		// Check "kube_exec"
		Context("check_kube_exec", func() {

			BeforeEach(func() {
				alert.Spec.Check = tapi.CheckPodExec
				alert.Spec.Vars = map[string]interface{}{
					"container": "busybox",
					"cmd":       "/bin/sh",
				}
			})

			Context("exit 0", func() {
				JustBeforeEach(func() {
					alert.Spec.Vars["argv"] = "exit 0"
				})

				It("should manage icinga service for Ok State", shouldManageIcingaServiceForLabelSelector)
			})

			Context("exit 2", func() {
				JustBeforeEach(func() {
					alert.Spec.Vars["argv"] = "exit 2"
				})

				It("should handle icinga service for Critical State", shouldHandleIcingaServiceForCriticalState)
			})

		})

	})
})
