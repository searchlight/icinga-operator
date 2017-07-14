package e2e_test

import (
	//"github.com/appscode/go/types"
	//tapi "github.com/appscode/searchlight/api"
	//. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
	"fmt"
)

var _ = Describe("PodAlert", func() {
	fmt.Println("----")
	//replicaSet := root.Invoke().ReplicaSet()
	//replicaSet.Spec.Replicas = types.Int32P(2)
	//err := root.CreateReplicaSet(replicaSet)
	//Expect(err).NotTo(HaveOccurred())
	//
	//By("Waiting for Running pods")
	//root.EventuallyReplicaSetRunning(replicaSet.ObjectMeta).Should(HaveRunningPods(*replicaSet.Spec.Replicas))
	//
	//replicaSet, err = root.GetReplicaSet(replicaSet.ObjectMeta)
	//Expect(err).NotTo(HaveOccurred())

	/*Describe(string(tapi.CheckPodStatus), func() {
		podAlert := root.Invoke().PodAlert()
		podAlert.Spec.Selector = *replicaSet.Spec.Selector
		podAlert.Spec.Check = tapi.CheckPodStatus
		err = root.CreatePodAlert(podAlert)
		Expect(err).NotTo(HaveOccurred())

	})*/
})
