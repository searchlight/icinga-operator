package e2e

import (
	"strings"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/types"
	kutil_core "github.com/appscode/kutil/core/v1"
	ext_util "github.com/appscode/kutil/extensions/v1beta1"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/client/clientset/versioned/typed/monitoring/v1alpha1/util"
	"github.com/appscode/searchlight/test/e2e/framework"
	. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("PodAlert", func() {
	var (
		err             error
		f               *framework.Invocation
		rs              *extensions.ReplicaSet
		ss              *apps.StatefulSet
		pod             *core.Pod
		alert           *api.PodAlert
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
			By("Create ReplicaSet: " + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			By("Create matching podalert: " + alert.Name)
			err = f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{OK: *rs.Spec.Replicas}))

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForNewPod = func() {
			By("Create ReplicaSet: " + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			By("Create matching podalert: " + alert.Name)
			err = f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{OK: *rs.Spec.Replicas}))

			By("Increase replica")
			rs, _, err := ext_util.PatchReplicaSet(f.KubeClient(), rs, func(in *extensions.ReplicaSet) *extensions.ReplicaSet {
				in.Spec.Replicas = types.Int32P(3)
				return in
			})
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{OK: *rs.Spec.Replicas}))

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForDeletedPod = func() {
			By("Create ReplicaSet: " + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			By("Create matching podalert: " + alert.Name)
			err = f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{OK: *rs.Spec.Replicas}))

			By("Decreate replica")
			rs, _, err := ext_util.PatchReplicaSet(f.KubeClient(), rs, func(in *extensions.ReplicaSet) *extensions.ReplicaSet {
				in.Spec.Replicas = types.Int32P(1)
				return in
			})
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{OK: *rs.Spec.Replicas}))

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldManageIcingaServiceForLabelChanged = func() {
			By("Create ReplicaSet: " + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

			By("Create matching podalert: " + alert.Name)
			err = f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{OK: *rs.Spec.Replicas}))

			alert, err = f.GetPodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			oldAlertSpec := alert.Spec

			By("Change LabelSelector")
			alert, _, err = util.PatchPodAlert(f.MonitoringClient(), alert, func(in *api.PodAlert) *api.PodAlert {
				in.Spec.Selector.MatchLabels = map[string]string{
					"app": rand.WithUniqSuffix("searchlight-e2e"),
				}
				return in
			})
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
			By("Create Pod: " + pod.Name)
			pod, err = f.CreatePod(pod)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for Running pods")
			f.EventuallyPodRunning(pod.ObjectMeta).Should(HaveRunningPods(1))

			By("Create matching podalert: " + alert.Name)
			err = f.CreatePodAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{OK: 1}))

			By("Delete podalert")
			err = f.DeletePodAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		shouldHandleIcingaServiceForCriticalState = func() {
			By("Create ReplicaSet: " + rs.Name)
			rs, err = f.CreateReplicaSet(rs)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for all pods")
			f.EventuallyReplicaSet(rs.ObjectMeta).Should(HavePods(*rs.Spec.Replicas))

			By("Create matching podalert: " + alert.Name)
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
			go f.DeleteReplicaSet(rs)
			go f.DeletePod(pod.ObjectMeta)
		})

		// Check "pod-status" and basic searchlight functionality
		Context("check_pod_status", func() {
			BeforeEach(func() {
				alert.Spec.Check = api.CheckPodStatus
				alert.Spec.Selector = rs.Spec.Selector
			})

			It("should manage icinga service for Alert.Spec.Selector", shouldManageIcingaServiceForLabelSelector)
			It("should manage icinga service for new Pod", shouldManageIcingaServiceForNewPod)
			It("should manage icinga service for deleted Pod", shouldManageIcingaServiceForDeletedPod)
			It("should manage icinga service for Alert.Spec.Selector changed", shouldManageIcingaServiceForLabelChanged)

			Context("PodName", func() {
				BeforeEach(func() {
					alert.Spec.PodName = &pod.Name
					alert.Spec.Selector = nil
				})

				It("should manage icinga service for Alert.Spec.PodName", shouldManageIcingaServiceForPodName)
			})

			Context("invalid image", func() {
				BeforeEach(func() {
					rs.Spec.Template.Spec.Containers[0].Image = "invalid-image"
				})
				It("should handle icinga service for Critical State", shouldHandleIcingaServiceForCriticalState)
			})

			Context("change labels", func() {
				BeforeEach(func() {
					alert.Spec.Selector = metav1.SetAsLabelSelector(pod.Labels)
				})

				It("should manage icinga service for Pod label changed", func() {
					By("Create Pod: " + pod.Name)
					pod, err = f.CreatePod(pod)
					Expect(err).NotTo(HaveOccurred())

					By("Wait for Running pods")
					f.EventuallyPodRunning(pod.ObjectMeta).Should(HaveRunningPods(1))

					By("Create matching podalert: " + alert.Name)
					err = f.CreatePodAlert(alert)
					Expect(err).NotTo(HaveOccurred())

					By("Check icinga services")
					f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
						Should(HaveIcingaObject(IcingaServiceState{OK: 1}))

					newAlert := alert.DeepCopy()
					newAlert.Name = newAlert.Name + "-new"
					newAlert.Spec.Selector.MatchLabels["app"] = newAlert.Spec.Selector.MatchLabels["app"] + "-new"

					By("Create podalert: " + newAlert.Name)
					err = f.CreatePodAlert(newAlert)
					Expect(err).NotTo(HaveOccurred())

					By("Patch Pod: " + pod.Name)
					_, _, err = kutil_core.PatchPod(f.KubeClient(), pod, func(in *core.Pod) *core.Pod {
						in.Labels["app"] = newAlert.Spec.Selector.MatchLabels["app"]
						return in
					})
					Expect(err).NotTo(HaveOccurred())

					By("Check icinga services " + newAlert.Name)
					f.EventuallyPodAlertIcingaService(newAlert.ObjectMeta, newAlert.Spec).
						Should(HaveIcingaObject(IcingaServiceState{OK: 1}))

					By("Wait for icinga services to be deleted " + alert.Name)
					f.EventuallyPodAlertIcingaService(alert.ObjectMeta, alert.Spec).
						Should(HaveIcingaObject(IcingaServiceState{}))

					By("Delete podalert " + alert.Name)
					err = f.DeletePodAlert(newAlert.ObjectMeta)
					Expect(err).NotTo(HaveOccurred())

					By("Wait for icinga services to be deleted " + newAlert.Name)
					f.EventuallyPodAlertIcingaService(newAlert.ObjectMeta, newAlert.Spec).
						Should(HaveIcingaObject(IcingaServiceState{}))

					By("Delete podalert " + alert.Name)
					err = f.DeletePodAlert(alert.ObjectMeta)
					Expect(err).NotTo(HaveOccurred())

				})
			})
		})

		// Check "volume"
		Context("check_pod_volume", func() {
			AfterEach(func() {
				go f.DeleteStatefulSet(ss)
			})
			BeforeEach(func() {
				if strings.ToLower(f.Provider) == "minikube" {
					skippingMessage = `"check_pod_volume" will not work in minikube"`
				}

				ss.Spec.Template.Spec.Containers[0].Command = []string{
					"/bin/sh",
					"-c",
					"dd if=/dev/zero of=/source/data/data bs=1024 count=52500 && sleep 1d",
				}
				alert.Spec.Check = api.CheckPodVolume
				alert.Spec.Selector = ss.Spec.Selector
				alert.Spec.Vars["volume_name"] = framework.TestSourceDataVolumeName
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

					By("Create matching podalert: " + alert.Name)
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
					icingaServiceState = IcingaServiceState{OK: *ss.Spec.Replicas}
					alert.Spec.Vars["warning"] = "100.0"
				})

				It("should manage icinga service for OK State", forStatefulSet)
			})

			Context("State Warning", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Warning: *ss.Spec.Replicas}
					alert.Spec.Vars["warning"] = "1.0"
				})

				It("should manage icinga service for Warning State", forStatefulSet)
			})

			Context("State Critical", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Critical: *ss.Spec.Replicas}
					alert.Spec.Vars["critical"] = "1.0"
				})

				It("should manage icinga service for Critical State", forStatefulSet)
			})

		})

	})
})
