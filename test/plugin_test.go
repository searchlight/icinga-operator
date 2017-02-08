package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/appscode/searchlight/pkg/client"
	config "github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
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
	kapi "k8s.io/kubernetes/pkg/api"
)

func TestComponentStatus(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckComponentStatus)

	expectedIcingaState := component_status.GetStatusCodeForComponentStatus()
	icingaState, _ := check_component_status.CheckComponentStatus()
	assert.EqualValues(t, expectedIcingaState, icingaState)
}

func TestJsonPath(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckJsonPath)

	testDataList := json_path.GetTestData()
	for _, testData := range testDataList {
		var req check_json_path.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_json_path.CheckJsonPath(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}
}

func TestKubeEvent(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckCommandKubeEvent)

	checkInterval, _ := time.ParseDuration("2m")
	clockSkew, _ := time.ParseDuration("0s")
	testDataList := kube_event.GetTestData(checkInterval, clockSkew)
	for _, testData := range testDataList {
		var req check_kube_event.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_kube_event.CheckKubeEvent(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}
}

func TestKubeExec(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckCommandKubeExec)

	context := &client.Context{}
	kubeClient, err := config.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	context.KubeClient = kubeClient

	// Run KubeD
	watcher := runKubeD(context)
	fmt.Println("--> Running kubeD")

	replicaSet := mini.CreateReplicaSet(watcher, kapi.NamespaceDefault)

	objectList, err := host.GetObjectList(kubeClient.Client, host.CheckCommandKubeExec, host.HostTypePod,
		replicaSet.Namespace, host.TypeReplicasets, replicaSet.Name, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	testDataList := kube_exec.GetTestData(objectList)
	for _, testData := range testDataList {
		var req check_kube_exec.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_kube_exec.CheckKubeExec(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}

	mini.DeleteReplicaSet(watcher, replicaSet)
}

func TestNodeCount(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckNodeCount)

	testDataList := node_count.GetTestData()
	for _, testData := range testDataList {
		var req check_node_count.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_node_count.CheckNodeCount(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}
}

func TestNodeStatus(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckNodeStatus)

	testDataList := node_status.GetTestData()
	for _, testData := range testDataList {
		var req check_node_status.Request
		plugin.FillStruct(testData.Data, &req)

		icingaState, _ := check_node_status.CheckNodeStatus(&req)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}
}

func TestPodExistsPodStatus(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckCommandPodExists, host.CheckCommandPodStatus)

	context := &client.Context{}
	kubeClient, err := config.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	context.KubeClient = kubeClient

	// Run KubeD
	watcher := runKubeD(context)
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
		testDataList := pod_status.GetTestData(watcher, objectType, objectName, namespace)
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
	replicationController := mini.CreateReplicationController(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(t, host.TypeReplicationcontrollers, replicationController.Name, replicationController.Namespace, int(replicationController.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(t, watcher, host.TypeReplicationcontrollers, replicationController.Name, replicationController.Namespace)
	fmt.Println("---- >> Deleting")
	mini.DeleteReplicationController(watcher, replicationController)

	// Daemonsets
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeDaemonsets)
	fmt.Println("---- >> Creating")
	daemonSet := mini.CreateDaemonSet(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(t, host.TypeDaemonsets, daemonSet.Name, daemonSet.Namespace, int(daemonSet.Status.DesiredNumberScheduled))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(t, watcher, host.TypeDaemonsets, daemonSet.Name, daemonSet.Namespace)
	fmt.Println("---- >> Deleting")
	mini.DeleteDaemonSet(watcher, daemonSet)

	// Deployments
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeDeployments)
	fmt.Println("---- >> Creating")
	deployment := mini.CreateDeployment(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(t, host.TypeDeployments, deployment.Name, deployment.Namespace, int(deployment.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(t, watcher, host.TypeDeployments, deployment.Name, deployment.Namespace)
	fmt.Println("---- >> Deleting")
	mini.DeleteDeployment(watcher, deployment)

	// StatefulSet
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeStatefulSet)
	fmt.Println("---- >> Creating")
	statefulSet := mini.CreateStatefulSet(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(t, host.TypeStatefulSet, statefulSet.Name, statefulSet.Namespace, int(statefulSet.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(t, watcher, host.TypeStatefulSet, statefulSet.Name, statefulSet.Namespace)
	fmt.Println(fmt.Sprintf(`---- >> Skip deleting "%s" for further test`, host.TypeStatefulSet))

	// Replicasets
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeReplicasets)
	fmt.Println("---- >> Creating")
	replicaSet := mini.CreateReplicaSet(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(t, host.TypeReplicasets, replicaSet.Name, replicaSet.Namespace, int(replicaSet.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(t, watcher, host.TypeReplicasets, replicaSet.Name, replicaSet.Namespace)
	fmt.Println(fmt.Sprintf(`---- >> Skip deleting "%s" for further test`, host.TypeReplicasets))

	// Services
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeServices)
	fmt.Println("---- >> Creating", host.TypeServices)
	service := mini.CreateService(watcher, replicaSet.Namespace, replicaSet.Spec.Template.Labels)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	checkPodExists(t, host.TypeServices, service.Name, service.Namespace, int(replicaSet.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(t, watcher, host.TypeServices, service.Name, service.Namespace)
	fmt.Println("---- >> Deleting", host.TypeServices)
	mini.DeleteService(watcher, service)

	// Cluster
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeCluster)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	totalPod := pod_exists.GetPodCount(watcher, kapi.NamespaceDefault)
	checkPodExists(t, "", "", kapi.NamespaceDefault, totalPod)
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	checkPodStatus(t, watcher, "", "", kapi.NamespaceDefault)

	// Delete skiped objects
	fmt.Println("-- >> Deleting", host.TypeStatefulSet)
	mini.DeleteStatefulSet(watcher, statefulSet)
	fmt.Println("-- >> Deleting", host.TypeReplicasets)
	mini.DeleteReplicaSet(watcher, replicaSet)
}
