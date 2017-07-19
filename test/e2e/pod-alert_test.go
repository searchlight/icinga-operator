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
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"strings"
)

var _ = Describe("PodAlert", func() {
	var (
		err             error
		f               *framework.Invocation
		rs              *extensions.ReplicaSet
		ss              *apps.StatefulSet
		pod             *apiv1.Pod
		alert           *tapi.PodAlert
		skippingMessage string
	)

	BeforeEach(func() {
		f = root.Invoke()
		rs = f.ReplicaSet()
		ss = f.StatefulSet()
		pod = f.Pod()
		alert = f.PodAlert()
		skippingMessage = ""
	})

	var (
		shouldManageIcingaServiceForLabelSelector = func() {
			By("Create ReplicaSet :" + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching podalert :"+ alert.Name)
			err = f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForNewPod = func() {
			By("Create ReplicaSet :" + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching podalert :"+ alert.Name)
			err = f.CreatePodAlert(alert)
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

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *rs.Spec.Replicas}))

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForDeletedPod = func() {
			By("Create ReplicaSet :" + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching podalert :"+ alert.Name)
			err = f.CreatePodAlert(alert)
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

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForLabelChanged = func() {
			By("Create ReplicaSet :" + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching podalert :"+ alert.Name)
			err = f.CreatePodAlert(alert)
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

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())
		}

		shouldManageIcingaServiceForPodName = func() {
			By("Create Pod :" + pod.Name)
			pod, err = f.CreatePod(pod)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyPodRunning(pod.ObjectMeta).Should(HaveRunningPods(1))

			alert.Spec.PodName = pod.Name

			By("Create matching podalert :"+ alert.Name)
			err = f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: 1}))

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldHandleIcingaServiceForCriticalState = func() {
			By("Create ReplicaSet :" + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for all pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HavePods(*rs.Spec.Replicas))

			alert.Spec.Selector = *(rs.Spec.Selector)

			By("Create matching podalert :"+ alert.Name)
			err = f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Critical: *rs.Spec.Replicas}))

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}
	)

	Describe("Test", func() {
		AfterEach(func() {
			go f.EventuallyDeleteReplicaSet(rs.ObjectMeta).Should(BeTrue())
			go f.DeletePod(pod.ObjectMeta)
		})

		// Check "pod_status" and basic searchlight functionality
		Context("check_pod_status", func() {
			BeforeEach(func() {
				alert.Spec.Check = tapi.CheckPodStatus
			})

			It("should manage icinga service for Alert.Spec.Selector", shouldManageIcingaServiceForLabelSelector)
			It("should manage icinga service for new Pod", shouldManageIcingaServiceForNewPod)
			It("should manage icinga service for deleted Pod", shouldManageIcingaServiceForDeletedPod)
			It("should manage icinga service for Alert.Spec.Selector changed", shouldManageIcingaServiceForLabelChanged)
			It("should manage icinga service for Alert.Spec.PodName", shouldManageIcingaServiceForPodName)

			Context("invalid image", func() {
				BeforeEach(func() {
					rs.Spec.Template.Spec.Containers[0].Image = "invalid-image"
				})
				It("should handle icinga service for Critical State", shouldHandleIcingaServiceForCriticalState)
			})
		})

		// Check "volume"
		Context("check_volume", func() {
			AfterEach(func() {
				go f.EventuallyDeleteStatefulSet(ss.ObjectMeta).Should(BeTrue())
			})
			BeforeEach(func() {
				if strings.ToLower(f.Provider) == "minikube" {
					skippingMessage = `"check_volume" will not work in minikube"`
				}

				ss.Spec.Template.Spec.Containers[0].Command = []string{
					"/bin/sh",
					"-c",
					"dd if=/dev/zero of=/source/data/data bs=1024 count=52500 && sleep 1d",
				}
				alert.Spec.Check = tapi.CheckVolume
				alert.Spec.Vars = map[string]interface{}{
					"volume_name": framework.TestSourceDataVolumeName,
				}
			})

			var icingaServiceState IcingaServiceState
			var (
				forStatefulSet = func() {
					if skippingMessage != "" {
						Skip(skippingMessage)
					}

					By("Create StatefulSet: " + ss.Name)
					ss, err = f.CreateStatefulSet(ss)
					Expect(err).NotTo(HaveOccurred())

					By("Wait for Running pods")
					f.EventuallyStatefulSet(ss.ObjectMeta).Should(HaveRunningPods(*ss.Spec.Replicas))

					alert.Spec.Selector = *(ss.Spec.Selector)

					By("Create matching podalert :"+ alert.Name)
					err = f.CreatePodAlert(alert)
					Expect(err).NotTo(HaveOccurred())

					By("Check icinga services")
					f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
						Should(HaveIcingaObject(icingaServiceState))

					By("Delete podalert")
					err = f.DeletePodAlert(alert.ObjectMeta)
					Expect(err).NotTo(HaveOccurred())

					By("Wait for icinga services to be deleted")
					f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
						Should(HaveIcingaObject(IcingaServiceState{}))
				}
			)

			Context("State OK", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Ok: *ss.Spec.Replicas}
					alert.Spec.Vars["warning"] = 100.0
				})

				It("should manage icinga service for Ok State", forStatefulSet)
			})

			Context("State Warning", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Warning: *ss.Spec.Replicas}
					alert.Spec.Vars["warning"] = 1.0
				})

				It("should manage icinga service for Warning State", forStatefulSet)
			})

			Context("State Critical", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Critical: *ss.Spec.Replicas}
					alert.Spec.Vars["critical"] = 1.0
				})

				It("should manage icinga service for Critical State", forStatefulSet)
			})

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
				BeforeEach(func() {
					alert.Spec.Vars["argv"] = "exit 0"
				})

				It("should manage icinga service for Ok State", shouldManageIcingaServiceForLabelSelector)
			})

			Context("exit 2", func() {
				BeforeEach(func() {
					alert.Spec.Vars["argv"] = "exit 2"
				})

				It("should handle icinga service for Critical State", shouldHandleIcingaServiceForCriticalState)
			})

		})

	})
})
