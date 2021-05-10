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
	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/test/e2e/framework"
	. "go.searchlight.dev/icinga-operator/test/e2e/matcher"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NodeAlert", func() {
	var (
		err                error
		f                  *framework.Invocation
		alert              *api.NodeAlert
		totalNode          int32
		icingaServiceState IcingaServiceState
		skippingMessage    string
	)

	BeforeEach(func() {
		f = root.Invoke()
		alert = f.NodeAlert()
		skippingMessage = ""
		totalNode, _ = f.CountNode()
	})

	var (
		shouldManageIcingaService = func() {
			if skippingMessage != "" {
				Skip(skippingMessage)
			}

			By("Create matching nodealert: " + alert.Name)
			err = f.CreateNodeAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyNodeAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(icingaServiceState))

			By("Delete nodealert")
			err = f.DeleteNodeAlert(alert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())

			By("Wait for icinga services to be deleted")
			f.EventuallyNodeAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		}
	)

	Describe("Test", func() {
		Context("check_node_status", func() {
			BeforeEach(func() {
				icingaServiceState = IcingaServiceState{OK: totalNode}
				alert.Spec.Check = api.CheckNodeStatus
			})

			It("should manage icinga service for OK State", shouldManageIcingaService)
		})

		// Check "node_volume"
		Context("node_volume", func() {
			BeforeEach(func() {
				skippingMessage = `"node_volume will not work without hostfact"`
				alert.Spec.Check = api.CheckNodeVolume
			})

			Context("State OK", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{OK: totalNode}
					alert.Spec.Vars["warning"] = "100.0"
				})

				It("should manage icinga service for OK State", shouldManageIcingaService)
			})

			Context("State Warning", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Warning: totalNode}
					alert.Spec.Vars["warning"] = "1.0"
				})

				It("should manage icinga service for Warning State", shouldManageIcingaService)
			})

			Context("State Critical", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Critical: totalNode}
					alert.Spec.Vars["critical"] = "1.0"
				})

				It("should manage icinga service for Critical State", shouldManageIcingaService)
			})
		})

	})
})
