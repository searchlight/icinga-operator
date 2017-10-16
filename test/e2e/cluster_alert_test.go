package e2e_test

import (
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/test/e2e/framework"
	. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/labels"
)

var _ = Describe("ClusterAlert", func() {
	var (
		err                error
		f                  *framework.Invocation
		rs                 *extensions.ReplicaSet
		alert              *api.ClusterAlert
		totalNode          int32
		icingaServiceState IcingaServiceState
	)

	BeforeEach(func() {
		f = root.Invoke()
		alert = f.ClusterAlert()
	})

	Describe("Test", func() {

		var shouldManageIcingaService = func() {
			By("Create matching clusteralert: " + alert.Name)
			err = f.CreateClusterAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyClusterAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(icingaServiceState))

			By("Delete clusteralert")
			err = f.DeleteClusterAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyClusterAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		Context("check_component_status", func() {
			BeforeEach(func() {
				alert.Spec.Check = api.CheckComponentStatus
				icingaServiceState = IcingaServiceState{Ok: 1}
			})

			It("should manage icinga service for Ok State", shouldManageIcingaService)
		})

		Context("check_node_exists", func() {
			BeforeEach(func() {
				alert.Spec.Check = api.CheckNodeExists
				totalNode, _ = f.CountNode()
			})

			Context("State OK", func() {
				BeforeEach(func() {
					alert.Spec.Vars["count"] = totalNode
					icingaServiceState = IcingaServiceState{Ok: 1}
				})

				It("should manage icinga service for Ok State", shouldManageIcingaService)
			})

			Context("State Critical", func() {
				BeforeEach(func() {
					alert.Spec.Vars["count"] = totalNode + 1
					icingaServiceState = IcingaServiceState{Critical: 1}
				})

				It("should manage icinga service for Critical State", shouldManageIcingaService)
			})

		})

		Context("check_pod_exists", func() {

			AfterEach(func() {
				go f.EventuallyDeleteReplicaSet(rs.ObjectMeta)
			})

			BeforeEach(func() {
				rs = f.ReplicaSet()
				alert.Spec.Check = api.CheckPodExists
				alert.Spec.Vars["selector"] = labels.SelectorFromSet(rs.Labels).String()
			})

			var shouldManageIcingaService = func() {
				By("Create ReplicaSet " + rs.Name + "@" + rs.Namespace)
				rs, err = f.CreateReplicaSet(rs)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for Running pods")
				f.EventuallyReplicaSet(rs.ObjectMeta).Should(HaveRunningPods(*rs.Spec.Replicas))

				By("Create matching clusteralert: " + alert.Name)
				err = f.CreateClusterAlert(alert)
				Expect(err).NotTo(HaveOccurred())

				By("Check icinga services")
				f.EventuallyClusterAlertIcingaService(alert.ObjectMeta, alert.Spec).
					Should(HaveIcingaObject(icingaServiceState))

				By("Delete clusteralert")
				err = f.DeleteClusterAlert(alert.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for icinga services to be deleted")
				f.EventuallyClusterAlertIcingaService(alert.ObjectMeta, alert.Spec).
					Should(HaveIcingaObject(IcingaServiceState{}))
			}

			Context("State OK", func() {
				BeforeEach(func() {
					alert.Spec.Vars["count"] = *rs.Spec.Replicas
					icingaServiceState = IcingaServiceState{Ok: 1}
				})

				It("should manage icinga service for Ok State", shouldManageIcingaService)
			})

			Context("State Critical", func() {
				BeforeEach(func() {
					alert.Spec.Vars["count"] = *rs.Spec.Replicas + 1
					icingaServiceState = IcingaServiceState{Critical: 1}
				})

				It("should manage icinga service for Critical State", shouldManageIcingaService)
			})

		})
	})
})
