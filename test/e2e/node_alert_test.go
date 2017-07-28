package e2e_test

import (
	"strings"

	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/test/e2e/framework"
	. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NodeAlert", func() {
	var (
		err                error
		f                  *framework.Invocation
		alert              *tapi.NodeAlert
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
				icingaServiceState = IcingaServiceState{Ok: totalNode}
				alert.Spec.Check = tapi.CheckNodeStatus
			})

			It("should manage icinga service for Ok State", shouldManageIcingaService)
		})

		// Check "node_volume"
		Context("node_volume", func() {
			BeforeEach(func() {
				if strings.ToLower(f.Provider) == "minikube" {
					skippingMessage = `"node_volume will not work in minikube"`
				}
				alert.Spec.Check = tapi.CheckNodeVolume
			})

			Context("State OK", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Ok: totalNode}
					alert.Spec.Vars["warning"] = 100.0
				})

				It("should manage icinga service for Ok State", shouldManageIcingaService)
			})

			Context("State Warning", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Warning: totalNode}
					alert.Spec.Vars["warning"] = 1.0
				})

				It("should manage icinga service for Warning State", shouldManageIcingaService)
			})

			Context("State Critical", func() {
				BeforeEach(func() {
					icingaServiceState = IcingaServiceState{Critical: totalNode}
					alert.Spec.Vars["critical"] = 1.0
				})

				It("should manage icinga service for Critical State", shouldManageIcingaService)
			})
		})

	})
})
