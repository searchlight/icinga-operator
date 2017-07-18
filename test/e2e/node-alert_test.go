package e2e_test

import (
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/test/e2e/framework"
	. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = FDescribe("NodeAlert", func() {
	var (
		err       error
		f         *framework.Invocation
		alert     *tapi.NodeAlert
		totalNode int32
	)

	BeforeEach(func() {
		f = root.Invoke()
		alert = f.NodeAlert()
	})

	var (
		shouldManageIcingaServiceForOkState = func() {
			By("Create matching nodealert")
			err = f.CreateNodeAlert(alert)
			Expect(err).NotTo(HaveOccurred())

			totalNode, err = f.CountNode()
			Expect(err).NotTo(HaveOccurred())

			By("Check icinga services")
			f.EventuallyNodeAlertIcingaService(alert.ObjectMeta, alert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: totalNode}))

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
				alert.Spec.Check = tapi.CheckNodeStatus
			})

			It("should manage icinga service for Ok State", shouldManageIcingaServiceForOkState)
		})
	})
})
