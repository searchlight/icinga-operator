package e2e_test

import (
	"github.com/appscode/go/types"
	tapi "github.com/appscode/searchlight/api"
	. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

var _ = Describe("PodAlert", func() {

	var (
		replicaSet *extensions.ReplicaSet
		podAlert   *tapi.PodAlert
		err        error
	)

	It("Create ReplicaSet", func() {
		replicaSet = root.Invoke().ReplicaSet()
		replicaSet.Spec.Replicas = types.Int32P(2)
		replicaSet, err = root.CreateReplicaSet(replicaSet)
		Expect(err).NotTo(HaveOccurred())
		By("Waiting for Running pods")
		root.EventuallyReplicaSetRunning(replicaSet.ObjectMeta).Should(HaveRunningPods(*replicaSet.Spec.Replicas))
	})

	Describe("Test pod_status", func() {
		It("Create PodAlert", func() {
			podAlert = root.Invoke().PodAlert()
			podAlert.Spec.Selector = *(replicaSet.Spec.Selector)
			podAlert.Spec.Check = tapi.CheckPodStatus
			err := root.CreatePodAlert(podAlert)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Check Icinga Object", func() {
			root.EventuallyPodAlertIcingaService(podAlert.ObjectMeta, podAlert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{Ok: *replicaSet.Spec.Replicas}))
		})

		It("Delete PodAlert", func() {
			err := root.DeletePodAlert(podAlert.ObjectMeta)
			Expect(err).NotTo(HaveOccurred())
		})

		It("Check Icinga Object", func() {
			root.EventuallyPodAlertIcingaService(podAlert.ObjectMeta, podAlert.Spec).
				Should(HaveIcingaObject(IcingaServiceState{}))
		})
	})

	It("Delete ReplicaSet", func() {
		err := root.DeleteReplicaSet(replicaSet.ObjectMeta)
		Expect(err).NotTo(HaveOccurred())
	})
})
