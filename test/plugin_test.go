package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/appscode/searchlight/plugins/check_component_status"
	"github.com/appscode/searchlight/plugins/check_json_path"
	"github.com/appscode/searchlight/plugins/check_kube_event"
	"github.com/appscode/searchlight/plugins/check_kube_exec"
	"github.com/appscode/searchlight/plugins/check_node_count"
	"github.com/appscode/searchlight/plugins/check_node_status"
	"github.com/appscode/searchlight/plugins/check_pod_exists"
	"github.com/appscode/searchlight/plugins/check_pod_status"
	"github.com/appscode/searchlight/test/mini"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/test/plugin/component_status"
	"github.com/appscode/searchlight/test/plugin/json_path"
	"github.com/appscode/searchlight/test/plugin/kube_event"
	"github.com/appscode/searchlight/test/plugin/kube_exec"
	"github.com/appscode/searchlight/test/plugin/node_count"
	"github.com/appscode/searchlight/test/plugin/node_status"
	"github.com/appscode/searchlight/test/plugin/pod_exists"
	"github.com/appscode/searchlight/test/plugin/pod_status"
	"github.com/stretchr/testify/assert"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func TestComponentStatus(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckComponentStatus)

	kubeClient, err := getKubeClient()
	if !assert.Nil(t, err) {
		return
	}
	expectedIcingaState, err := component_status.GetStatusCodeForComponentStatus(kubeClient)
	if !assert.Nil(t, err) {
		return
	}
	icingaState, _ := check_component_status.CheckComponentStatus()
	assert.EqualValues(t, expectedIcingaState, icingaState)
}

func TestJsonPath(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckJsonPath)

	testDataList, err := json_path.GetTestData()
	if !assert.Nil(t, err) {
		return
	}

	for _, testData := range testDataList {
		var req check_json_path.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_json_path.CheckJsonPath(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}
}

func TestKubeEvent(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckCommandKubeEvent)

	kubeClient, err := getKubeClient()
	if !assert.Nil(t, err) {
		return
	}

	checkInterval, _ := time.ParseDuration("2m")
	clockSkew, _ := time.ParseDuration("0s")
	testDataList, err := kube_event.GetTestData(kubeClient, checkInterval, clockSkew)
	if !assert.Nil(t, err) {
		return
	}

	for _, testData := range testDataList {
		var req check_kube_event.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_kube_event.CheckKubeEvent(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}
}

func TestKubeExec(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckCommandKubeExec)

	// Run KubeD
	watcher, err := runKubeD(false)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("--> Running kubeD")

	replicaSet, err := mini.CreateReplicaSet(watcher, apiv1.NamespaceDefault)
	if !assert.Nil(t, err) {
		return
	}

	objectList, err := host.GetObjectList(watcher.KubeClient, host.CheckCommandKubeExec, host.HostTypePod,
		replicaSet.Namespace, host.TypeReplicasets, replicaSet.Name, "")
	if !assert.Nil(t, err) {
		return
	}

	testDataList, err := kube_exec.GetTestData(objectList)
	if !assert.Nil(t, err) {
		return
	}
	for _, testData := range testDataList {
		var req check_kube_exec.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_kube_exec.CheckKubeExec(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}

	err = mini.DeleteReplicaSet(watcher, replicaSet)
	if !assert.Nil(t, err) {
		return
	}
}

func TestNodeCount(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckNodeCount)

	kubeClient, err := getKubeClient()
	if !assert.Nil(t, err) {
		return
	}

	testDataList, err := node_count.GetTestData(kubeClient)
	if !assert.Nil(t, err) {
		return
	}

	for _, testData := range testDataList {
		var req check_node_count.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_node_count.CheckNodeCount(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}
}

func TestNodeStatus(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckNodeStatus)

	kubeClient, err := getKubeClient()
	if !assert.Nil(t, err) {
		return
	}

	testDataList, err := node_status.GetTestData(kubeClient)
	if !assert.Nil(t, err) {
		return
	}

	for _, testData := range testDataList {
		var req check_node_status.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_node_status.CheckNodeStatus(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}
}

func TestPodExistsPodStatus(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckCommandPodExists, host.CheckCommandPodStatus)

	// Run KubeD
	watcher, err := runKubeD(false)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("--> Running kubeD")

	checkPodExists := func(objectType, objectName, namespace string, count int) {
		testDataList := pod_exists.GetTestData(objectType, objectName, namespace, count)
		for _, testData := range testDataList {
			var req check_pod_exists.Request
			plugin.FillStruct(testData.Data, &req)
			isCountSet := false
			if req.Count != 0 {
				isCountSet = true
			}
			icingaState, _ := check_pod_exists.CheckPodExists(&req, isCountSet)
			assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
		}
	}

	checkPodStatus := func(objectType, objectName, namespace string) {
		testDataList, err := pod_status.GetTestData(watcher, objectType, objectName, namespace)
		if !assert.Nil(t, err) {
			return
		}
		for _, testData := range testDataList {
			var req check_pod_status.Request
			plugin.FillStruct(testData.Data, &req)
			icingaState, _ := check_pod_status.CheckPodStatus(&req)
			assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
		}
	}

	// Replicationcontrollers
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeReplicationcontrollers)
	fmt.Println("---- >> Creating")
	replicationController, err := mini.CreateReplicationController(watcher, apiv1.NamespaceDefault)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(host.TypeReplicationcontrollers, replicationController.Name, replicationController.Namespace, int(replicationController.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(host.TypeReplicationcontrollers, replicationController.Name, replicationController.Namespace)
	fmt.Println("---- >> Deleting")
	err = mini.DeleteReplicationController(watcher, replicationController)
	if !assert.Nil(t, err) {
		return
	}

	// Daemonsets
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeDaemonsets)
	fmt.Println("---- >> Creating")
	daemonSet, err := mini.CreateDaemonSet(watcher, apiv1.NamespaceDefault)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(host.TypeDaemonsets, daemonSet.Name, daemonSet.Namespace, int(daemonSet.Status.DesiredNumberScheduled))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(host.TypeDaemonsets, daemonSet.Name, daemonSet.Namespace)
	fmt.Println("---- >> Deleting")
	err = mini.DeleteDaemonSet(watcher, daemonSet)
	if !assert.Nil(t, err) {
		return
	}

	// Deployments
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeDeployments)
	fmt.Println("---- >> Creating")
	deployment, err := mini.CreateDeployment(watcher, apiv1.NamespaceDefault)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(host.TypeDeployments, deployment.Name, deployment.Namespace, int(deployment.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(host.TypeDeployments, deployment.Name, deployment.Namespace)
	fmt.Println("---- >> Deleting")
	err = mini.DeleteDeployment(watcher, deployment)
	if !assert.Nil(t, err) {
		return
	}

	// StatefulSet
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeStatefulSet)
	fmt.Println("---- >> Creating")
	statefulSet, err := mini.CreateStatefulSet(watcher, apiv1.NamespaceDefault)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(host.TypeStatefulSet, statefulSet.Name, statefulSet.Namespace, int(statefulSet.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(host.TypeStatefulSet, statefulSet.Name, statefulSet.Namespace)
	fmt.Println(fmt.Sprintf(`---- >> Skip deleting "%s" for further test`, host.TypeStatefulSet))

	// Replicasets
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeReplicasets)
	fmt.Println("---- >> Creating")
	replicaSet, err := mini.CreateReplicaSet(watcher, apiv1.NamespaceDefault)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(host.TypeReplicasets, replicaSet.Name, replicaSet.Namespace, int(replicaSet.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(host.TypeReplicasets, replicaSet.Name, replicaSet.Namespace)
	fmt.Println(fmt.Sprintf(`---- >> Skip deleting "%s" for further test`, host.TypeReplicasets))

	// Services
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeServices)
	fmt.Println("---- >> Creating", host.TypeServices)
	service, err := mini.CreateService(watcher, replicaSet.Namespace, replicaSet.Spec.Template.Labels)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(host.TypeServices, service.Name, service.Namespace, int(replicaSet.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(host.TypeServices, service.Name, service.Namespace)
	fmt.Println("---- >> Deleting", host.TypeServices)
	err = mini.DeleteService(watcher, service)
	if !assert.Nil(t, err) {
		return
	}

	// Cluster
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeCluster)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	totalPod, err := pod_exists.GetPodCount(watcher, apiv1.NamespaceDefault)
	if !assert.Nil(t, err) {
		return
	}
	checkPodExists("", "", apiv1.NamespaceDefault, totalPod)
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus("", "", apiv1.NamespaceDefault)

	// Delete skiped objects
	fmt.Println("-- >> Deleting", host.TypeStatefulSet)
	err = mini.DeleteStatefulSet(watcher, statefulSet)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("-- >> Deleting", host.TypeReplicasets)
	err = mini.DeleteReplicaSet(watcher, replicaSet)
	if !assert.Nil(t, err) {
		return
	}
}
