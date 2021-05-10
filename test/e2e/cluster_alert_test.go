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
	"strconv"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/test/e2e/framework"
	. "go.searchlight.dev/icinga-operator/test/e2e/matcher"

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
			f.EventuallyClusterAlertIcingaService(alert.ObjectMeta).
				Should(HaveIcingaObject(icingaServiceState))

			By("Delete clusteralert")
			err = f.DeleteClusterAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyClusterAlertIcingaService(alert.ObjectMeta).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}

		Context("check_component_status", func() {
			BeforeEach(func() {
				alert.Spec.Check = api.CheckComponentStatus
				icingaServiceState = IcingaServiceState{OK: 1}
			})

			It("should manage icinga service for OK State", shouldManageIcingaService)
		})

		Context("check_node_exists", func() {
			BeforeEach(func() {
				alert.Spec.Check = api.CheckNodeExists
				totalNode, _ = f.CountNode()
			})

			Context("State OK", func() {
				BeforeEach(func() {
					alert.Spec.Vars["count"] = strconv.Itoa(int(totalNode))
					icingaServiceState = IcingaServiceState{OK: 1}
				})

				It("should manage icinga service for OK State", shouldManageIcingaService)
			})

			Context("State Critical", func() {
				BeforeEach(func() {
					alert.Spec.Vars["count"] = strconv.Itoa(int(totalNode + 1))
					icingaServiceState = IcingaServiceState{Critical: 1}
				})

				It("should manage icinga service for Critical State", shouldManageIcingaService)
			})

		})

		Context("check_pod_exists", func() {

			AfterEach(func() {
				f.DeleteReplicaSet(rs)
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
				f.EventuallyClusterAlertIcingaService(alert.ObjectMeta).
					Should(HaveIcingaObject(icingaServiceState))

				By("Delete clusteralert")
				err = f.DeleteClusterAlert(alert.ObjectMeta)
				Expect(err).NotTo(HaveOccurred())

				By("Wait for icinga services to be deleted")
				f.EventuallyClusterAlertIcingaService(alert.ObjectMeta).
					Should(HaveIcingaObject(IcingaServiceState{}))
			}

			Context("State OK", func() {
				BeforeEach(func() {
					alert.Spec.Vars["count"] = strconv.Itoa(int(*rs.Spec.Replicas))
					icingaServiceState = IcingaServiceState{OK: 1}
				})

				It("should manage icinga service for OK State", shouldManageIcingaService)
			})

			Context("State Critical", func() {
				BeforeEach(func() {
					alert.Spec.Vars["count"] = strconv.Itoa(int(*rs.Spec.Replicas + 1))
					icingaServiceState = IcingaServiceState{Critical: 1}
				})

				It("should manage icinga service for Critical State", shouldManageIcingaService)
			})

		})
	})
})
