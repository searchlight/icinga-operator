package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client"
	config "github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/plugins/check_component_status"
	"github.com/appscode/searchlight/plugins/check_json_path"
	"github.com/appscode/searchlight/plugins/check_kube_event"
	"github.com/appscode/searchlight/plugins/check_kube_exec"
	"github.com/appscode/searchlight/plugins/check_node_count"
	"github.com/appscode/searchlight/plugins/check_node_status"
	"github.com/appscode/searchlight/test/mini"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/test/plugin/component_status"
	"github.com/appscode/searchlight/test/plugin/json_path"
	"github.com/appscode/searchlight/test/plugin/kube_event"
	"github.com/appscode/searchlight/test/plugin/node_count"
	"github.com/appscode/searchlight/test/plugin/node_status"
	"github.com/appscode/searchlight/test/plugin/pod_exists"
	"github.com/appscode/searchlight/util"
	"github.com/stretchr/testify/assert"
	kapi "k8s.io/kubernetes/pkg/api"
	"github.com/appscode/searchlight/test/plugin/pod_status"
)

type testData struct {
	data                map[string]interface{}
	expectedIcingaState util.IcingaState
	deleteObject        bool
}

func getKubernetesClient() *config.KubeClient {
	kubeClient, err := config.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return kubeClient
}

func TestComponentStatus(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckComponentStatus)

	kubeClient, err := config.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	expectedIcingaState := component_status.GetStatusCodeForComponentStatus(kubeClient)
	icingaState, _ := check_component_status.CheckComponentStatus()
	assert.EqualValues(t, expectedIcingaState, icingaState)
}

func TestJsonPath(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckJsonPath)

	url := "https://api.github.com"
	uri := "/orgs/appscode"

	repoNumber, err := json_path.GetPublicRepoNumber()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"Url":     url + uri,
				"Query":   ".",
				"Warning": fmt.Sprintf(`.public_repos!=%v`, repoNumber),
			},
			expectedIcingaState: 0,
		},
		testData{
			data: map[string]interface{}{
				"Url":     url + uri,
				"Query":   ".",
				"Warning": fmt.Sprintf(`.public_repos==%v`, repoNumber),
			},
			expectedIcingaState: 1,
		},
		testData{
			data: map[string]interface{}{
				"Url":      url + uri,
				"Query":    ".",
				"Warning":  fmt.Sprintf(`.public_repos==%v`, repoNumber-1),
				"Critical": fmt.Sprintf(`.public_repos==%v`, repoNumber),
			},
			expectedIcingaState: 2,
		},
		testData{
			data: map[string]interface{}{
				"Url":     url + uri + "fake",
				"Query":   ".",
				"Warning": fmt.Sprintf(`.public_repos==%v`, repoNumber-1),
			},
			expectedIcingaState: 3,
		},
	}

	for _, testData := range testDataList {
		var req check_json_path.Request
		plugin.FillStruct(testData.data, &req)

		icingaState, _ := check_json_path.CheckJsonPath(&req)
		assert.EqualValues(t, testData.expectedIcingaState, icingaState)
	}
}

func TestKubeEvent(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckCommandKubeEvent)

	kubeClient, err := config.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	checkInterval, _ := time.ParseDuration("2m")
	clockSkew, _ := time.ParseDuration("0s")
	expectedIcingaState := kube_event.GetStatusCodeForEventCount(kubeClient, checkInterval, clockSkew)

	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"CheckInterval": checkInterval,
				"ClockSkew":     clockSkew,
			},
			expectedIcingaState: expectedIcingaState,
		},
	}
	for _, testData := range testDataList {
		var req check_kube_event.Request
		plugin.FillStruct(testData.data, &req)

		icingaState, _ := check_kube_event.CheckKubeEvent(&req)
		assert.EqualValues(t, testData.expectedIcingaState, icingaState)
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

	testDataList := make([]testData, 0)
	for _, object := range objectList {
		_, objectName, namespace := plugin.GetKubeObjectInfo(object.Name)
		testDataList = append(testDataList, testData{
			data: map[string]interface{}{
				"Pod":       objectName,
				"Namespace": namespace,
				"Command":   "/bin/sh",
				"Arg":       "exit 0",
			},
			expectedIcingaState: 0,
		})
		testDataList = append(testDataList, testData{
			data: map[string]interface{}{
				"Pod":       objectName,
				"Namespace": namespace,
				"Command":   "/bin/sh",
				"Arg":       "exit 5",
			},
			expectedIcingaState: 2,
		})
	}

	for _, testData := range testDataList {
		var req check_kube_exec.Request
		plugin.FillStruct(testData.data, &req)

		icingaState, _ := check_kube_exec.CheckKubeExec(&req)
		assert.EqualValues(t, testData.expectedIcingaState, icingaState)
	}

	mini.DeleteReplicaSet(watcher, replicaSet)
}

func TestNodeCount(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckNodeCount)

	kubeClient, err := config.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	actualNodeCount := node_count.GetKubernetesNodeCount(kubeClient)

	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"Count": actualNodeCount,
			},
			expectedIcingaState: 0,
		},
		testData{
			data: map[string]interface{}{
				"Count": actualNodeCount + 1,
			},
			expectedIcingaState: 2,
		},
	}

	for _, testData := range testDataList {
		var req check_node_count.Request
		plugin.FillStruct(testData.data, &req)

		icingaState, _ := check_node_count.CheckNodeCount(&req)
		assert.EqualValues(t, testData.expectedIcingaState, icingaState)
	}
}

func TestNodeStatus(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckNodeStatus)

	kubeClient, err := config.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	actualNodeName := node_status.GetKubernetesNodeName(kubeClient)

	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"Name": actualNodeName,
			},
			expectedIcingaState: 0,
		},
		testData{
			data: map[string]interface{}{
				// make node name invalid using random 2 character.
				"Name": actualNodeName + rand.Characters(2),
			},
			expectedIcingaState: 3,
		},
	}

	for _, testData := range testDataList {
		var req check_node_status.Request
		plugin.FillStruct(testData.data, &req)

		icingaState, _ := check_node_status.CheckNodeStatus(&req)
		assert.EqualValues(t, testData.expectedIcingaState, icingaState)
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

	// Replicationcontrollers
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeReplicationcontrollers)
	fmt.Println("---- >> Creating")
	replicationController := mini.CreateReplicationController(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	pod_exists.CheckPodExists(t, host.TypeReplicationcontrollers, replicationController.Name, replicationController.Namespace, int(replicationController.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	pod_status.CheckPodStatus(t, watcher, host.TypeReplicationcontrollers, replicationController.Name, replicationController.Namespace)
	fmt.Println("---- >> Deleting")
	mini.DeleteReplicationController(watcher, replicationController)

	// Daemonsets
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeDaemonsets)
	fmt.Println("---- >> Creating")
	daemonSet := mini.CreateDaemonSet(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	pod_exists.CheckPodExists(t, host.TypeDaemonsets, daemonSet.Name, daemonSet.Namespace, int(daemonSet.Status.DesiredNumberScheduled))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	pod_status.CheckPodStatus(t, watcher, host.TypeDaemonsets, daemonSet.Name, daemonSet.Namespace)
	fmt.Println("---- >> Deleting")
	mini.DeleteDaemonSet(watcher, daemonSet)

	// Deployments
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeDeployments)
	fmt.Println("---- >> Creating")
	deployment := mini.CreateDeployment(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	pod_exists.CheckPodExists(t, host.TypeDeployments, deployment.Name, deployment.Namespace, int(deployment.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	pod_status.CheckPodStatus(t, watcher, host.TypeDeployments, deployment.Name, deployment.Namespace)
	fmt.Println("---- >> Deleting")
	mini.DeleteDeployment(watcher, deployment)

	// StatefulSet
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeStatefulSet)
	fmt.Println("---- >> Creating")
	statefulSet := mini.CreateStatefulSet(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	pod_exists.CheckPodExists(t, host.TypeStatefulSet, statefulSet.Name, statefulSet.Namespace, int(statefulSet.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	pod_status.CheckPodStatus(t, watcher, host.TypeStatefulSet, statefulSet.Name, statefulSet.Namespace)
	fmt.Println(fmt.Sprintf(`---- >> Skip deleting "%s" for further test`, host.TypeStatefulSet))

	// Replicasets
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeReplicasets)
	fmt.Println("---- >> Creating")
	replicaSet := mini.CreateReplicaSet(watcher, kapi.NamespaceDefault)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	pod_exists.CheckPodExists(t, host.TypeReplicasets, replicaSet.Name, replicaSet.Namespace, int(replicaSet.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	pod_status.CheckPodStatus(t, watcher, host.TypeReplicasets, replicaSet.Name, replicaSet.Namespace)
	fmt.Println(fmt.Sprintf(`---- >> Skip deleting "%s" for further test`, host.TypeReplicasets))

	// Services
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeServices)
	fmt.Println("---- >> Creating", host.TypeServices)
	service := mini.CreateService(watcher, replicaSet.Namespace, replicaSet.Spec.Template.Labels)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	pod_exists.CheckPodExists(t, host.TypeServices, service.Name, service.Namespace, int(replicaSet.Spec.Replicas))
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	pod_status.CheckPodStatus(t, watcher, host.TypeServices, service.Name, service.Namespace)
	fmt.Println("---- >> Deleting", host.TypeServices)
	mini.DeleteService(watcher, service)

	// Cluster
	fmt.Println()
	fmt.Println("-- >> Testing plugings for", host.TypeCluster)
	fmt.Println("---- >> Testing", host.CheckCommandPodExists)
	totalPod := pod_exists.GetPodCount(watcher, kapi.NamespaceDefault)
	pod_exists.CheckPodExists(t, "", "", kapi.NamespaceDefault, totalPod)
	fmt.Println("---- >> Testing", host.CheckCommandPodStatus)
	pod_status.CheckPodStatus(t, watcher, "", "", kapi.NamespaceDefault)

	// Delete skiped objects
	fmt.Println("-- >> Deleting", host.TypeStatefulSet)
	mini.DeleteStatefulSet(watcher, statefulSet)
	fmt.Println("-- >> Deleting", host.TypeReplicasets)
	mini.DeleteReplicaSet(watcher, replicaSet)
}
