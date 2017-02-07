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
	"github.com/appscode/searchlight/test/plugin/pod_status"
	"github.com/appscode/searchlight/util"
	"github.com/stretchr/testify/assert"
	kapi "k8s.io/kubernetes/pkg/api"
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

	mini.DeleteReplicaSet(kubeClient, replicaSet)
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

func TestPodExists(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckCommandPodExists)

	kubeClient := getKubernetesClient()
	testPodExists := func(dataConfig *dataConfig) {
		// This will create object & return icinga_host name
		// and number of pods under it
		name, count := getTestData(kubeClient, dataConfig)
		time.Sleep(time.Second * 30)

		testDataList := []testData{
			testData{
				// To check for any pods
				data: map[string]interface{}{
					"host": name,
				},
				expectedIcingaState: 0,
			},
			testData{
				// To check for specific number of pods
				data: map[string]interface{}{
					"host":  name,
					"count": count,
				},
				expectedIcingaState: 0,
			},
			testData{
				// To check for critical when pod number mismatch
				data: map[string]interface{}{
					"host":  name,
					"count": count + 1,
				},
				expectedIcingaState: 2,
				deleteObject:        true,
			},
		}

		for _, testData := range testDataList {
			argList := []string{
				"check_pod_exists",
			}
			for key, val := range testData.data {
				argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
			}
			//statusCode := execCheckCommand("hyperalert", argList...)
			//assert.EqualValues(t, testData.expectedCode, statusCode)
		}
	}

	ns := "e2e"
	dataConfig := &dataConfig{
		Namespace: ns,
	}

	fmt.Println(">> Creating namespace", ns)
	createNewNamespace(kubeClient, ns)
	fmt.Println()

	fmt.Println(">> Testing plugings for", host.TypeReplicationcontrollers)
	dataConfig.ObjectType = host.TypeReplicationcontrollers
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeReplicasets)
	dataConfig.ObjectType = host.TypeReplicasets
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeDaemonsets)
	dataConfig.ObjectType = host.TypeDaemonsets
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeDeployments)
	dataConfig.ObjectType = host.TypeDeployments
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeServices)
	dataConfig.ObjectType = host.TypeServices
	testPodExists(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeCluster)
	dataConfig.ObjectType = host.TypeCluster
	dataConfig.CheckCommand = host.CheckCommandPodExists
	testPodExists(dataConfig)

	fmt.Println(">> Deleting namespace", ns)
	deleteNewNamespace(kubeClient, ns)

	fmt.Println()
}

func TestPodStatus(t *testing.T) {
	fmt.Println("== Plugin Testing >", host.CheckCommandPodStatus)

	kubeClient := getKubernetesClient()

	testPodStatus := func(dataConfig *dataConfig) {
		// This will create object & return icinga_host name
		name, _ := getTestData(kubeClient, dataConfig)
		time.Sleep(time.Second * 30)

		// This will check pod status under specific object
		// and will return 2 (critical) if any pod is not running
		expectedCode := pod_status.GetStatusCodeForPodStatus(kubeClient, name)

		testDataList := []testData{
			testData{
				data: map[string]interface{}{
					"host": name,
				},
				expectedIcingaState: expectedCode,
			},
		}

		for _, testData := range testDataList {
			argList := []string{
				"check_pod_status",
			}
			for key, val := range testData.data {
				argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
			}
			//statusCode := execCheckCommand("hyperalert", argList...)
			//assert.EqualValues(t, testData.expectedCode, statusCode)
		}
	}

	ns := "e2e"
	dataConfig := &dataConfig{
		Namespace: ns,
	}

	fmt.Println(">> Creating namespace", ns)
	createNewNamespace(kubeClient, ns)
	fmt.Println()

	fmt.Println(">> Testing plugings for", host.TypeReplicationcontrollers)
	dataConfig.ObjectType = host.TypeReplicationcontrollers
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeReplicasets)
	dataConfig.ObjectType = host.TypeReplicasets
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeDaemonsets)
	dataConfig.ObjectType = host.TypeDaemonsets
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeDeployments)
	dataConfig.ObjectType = host.TypeDeployments
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeServices)
	dataConfig.ObjectType = host.TypeServices
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypePods)
	dataConfig.ObjectType = host.TypePods
	testPodStatus(dataConfig)

	fmt.Println(">> Testing plugings for", host.TypeCluster)
	dataConfig.ObjectType = host.TypeCluster
	dataConfig.CheckCommand = host.CheckCommandPodStatus
	testPodStatus(dataConfig)

	fmt.Println(">> Deleting namespace", ns)
	deleteNewNamespace(kubeClient, ns)

	fmt.Println()
}
