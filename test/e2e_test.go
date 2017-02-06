package e2e

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/appscode/go/crypto/rand"
	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/k8s-addons/pkg/events"
	testing_lib "github.com/appscode/k8s-addons/pkg/testing"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/client"
	"github.com/appscode/searchlight/pkg/client/icinga"
	config "github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/pkg/controller/types"
	"github.com/appscode/searchlight/test/general"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/test/plugin/component_status"
	"github.com/appscode/searchlight/test/plugin/kube_event"
	"github.com/appscode/searchlight/test/plugin/node_count"
	"github.com/appscode/searchlight/test/plugin/node_status"
	"github.com/appscode/searchlight/test/plugin/pod_status"
	"github.com/appscode/searchlight/test/util"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v1/json"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	ext "k8s.io/kubernetes/pkg/apis/extensions"
)

type testData struct {
	data         map[string]interface{}
	expectedCode int
	deleteObject bool
}

func getKubernetesClient() *config.KubeClient {
	kubeClient, err := config.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return kubeClient
}

func execCheckCommand(cmdName string, cmdArgs ...string) (code int) {
	cmd := fmt.Sprintf("../dist/%s/%s-%s-%s", cmdName, cmdName, GOHOSTOS, GOHOSTARCH)
	command := exec.Command(cmd, cmdArgs...)
	fmt.Println(fmt.Sprintf("Running: %v", command.Args))
	cmdOut, err := command.Output()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code = status.ExitStatus()
			}
		}
	}
	fmt.Println(string(cmdOut))
	return
}

func runKubeD(context *client.Context) {
	fmt.Println("-- TestE2E: Waiting for kubed")
	w := &app.Watcher{
		Watcher: acw.Watcher{
			Client:                  context.KubeClient.Client,
			AppsCodeExtensionClient: context.KubeClient.AppscodeExtensionClient,
			SyncPeriod:              time.Minute * 2,
		},
		IcingaClient: context.IcingaClient,
	}

	w.Watcher.Dispatch = w.Dispatch
	w.Run()
}

func TestKubeD(t *testing.T) {
	os.Setenv("E2E_ICINGA_SECRET", "appscode-icinga.kube-system")

	context := &client.Context{}
	kubeClient := getKubernetesClient()
	context.KubeClient = kubeClient

	/*
		TODO: Pass Secret Name
	*/
	icingaClient, err := icinga.NewIcingaClient(kubeClient.Client, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	context.IcingaClient = icingaClient

	go runKubeD(context)

	time.Sleep(time.Minute * 5)
}

func TestComponentStatus(t *testing.T) {
	fmt.Println("== Testing >", host.CheckComponentStatus)

	statusCode := execCheckCommand("hyperalert", "check_component_status")
	kubeClient := getKubernetesClient()
	expectedStatusCode := component_status.GetStatusCodeForComponentStatus(kubeClient)
	assert.EqualValues(t, expectedStatusCode, statusCode)
}

func TestJsonPath(t *testing.T) {
	fmt.Println("== Testing >", host.CheckJsonPath)

	url := "https://api.appscode.info"
	uri := "/health"
	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"url":     url + uri,
				"query":   ".status",
				"warning": `.status!="OK"`,
			},
			expectedCode: 0,
		},
		testData{
			data: map[string]interface{}{
				"url":     url + uri,
				"query":   ".version",
				"warning": `.version=="0.3.0"`,
			},
			expectedCode: 1,
		},
		testData{
			data: map[string]interface{}{
				"url":     url + uri,
				"query":   ".version",
				"warning": `.version=="0.2.6"`, "critical": `.version=="0.3.0"`,
			},
			expectedCode: 2,
		},
		testData{
			data: map[string]interface{}{
				"url":     url + "/healthz",
				"query":   ".version",
				"warning": `.version=="0.2.6"`,
			},
			expectedCode: 3,
		},
	}

	for _, testData := range testDataList {
		argList := []string{
			"check_json_path",
		}
		for key, val := range testData.data {
			argList = append(argList, fmt.Sprintf(`--%s=%s`, key, val))
		}
		statusCode := execCheckCommand("hyperalert", argList...)
		assert.EqualValues(t, testData.expectedCode, statusCode)
	}
}

func TestNodeCount(t *testing.T) {
	fmt.Println("== Testing >", host.CheckNodeCount)

	kubeClient := getKubernetesClient()
	actualNodeCount := node_count.GetKubernetesNodeCount(kubeClient)

	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"count": actualNodeCount,
			},
			expectedCode: 0,
		},
		testData{
			data: map[string]interface{}{
				"count": actualNodeCount + 1,
			},
			expectedCode: 2,
		},
		testData{
			data:         map[string]interface{}{},
			expectedCode: 3,
		},
	}

	for _, testData := range testDataList {
		argList := []string{
			"check_node_count",
		}
		for key, val := range testData.data {
			argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
		}
		statusCode := execCheckCommand("hyperalert", argList...)
		assert.EqualValues(t, testData.expectedCode, statusCode)
	}
}

func TestNodeStatus(t *testing.T) {
	fmt.Println("== Testing >", host.CheckNodeStatus)

	kubeClient := getKubernetesClient()
	actualNodeName := node_status.GetKubernetesNodeName(kubeClient)
	hostname := actualNodeName + "@default"

	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"host": hostname,
			},
			expectedCode: 0,
		},
		testData{
			data: map[string]interface{}{
				// make node name invalid using random 2 character.
				// Added as prefix because 1st part of hostname is nodename. (<node-name>@<alert-namespace>)
				"host": rand.Characters(2) + hostname,
			},
			expectedCode: 3,
		},
	}

	for _, testData := range testDataList {
		argList := []string{
			"check_node_status",
		}
		for key, val := range testData.data {
			argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
		}
		statusCode := execCheckCommand("hyperalert", argList...)
		assert.EqualValues(t, testData.expectedCode, statusCode)
	}
}

func TestPodExists(t *testing.T) {
	fmt.Println("== Testing >", host.CheckCommandPodExists)

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
				expectedCode: 0,
			},
			testData{
				// To check for specific number of pods
				data: map[string]interface{}{
					"host":  name,
					"count": count,
				},
				expectedCode: 0,
			},
			testData{
				// To check for critical when pod number mismatch
				data: map[string]interface{}{
					"host":  name,
					"count": count + 1,
				},
				expectedCode: 2,
				deleteObject: true,
			},
		}

		for _, testData := range testDataList {
			argList := []string{
				"check_pod_exists",
			}
			for key, val := range testData.data {
				argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
			}
			statusCode := execCheckCommand("hyperalert", argList...)
			assert.EqualValues(t, testData.expectedCode, statusCode)
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
	fmt.Println("== Testing >", host.CheckCommandPodStatus)

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
				expectedCode: expectedCode,
			},
		}

		for _, testData := range testDataList {
			argList := []string{
				"check_pod_status",
			}
			for key, val := range testData.data {
				argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
			}
			statusCode := execCheckCommand("hyperalert", argList...)
			assert.EqualValues(t, testData.expectedCode, statusCode)
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

func TestKubeEvent(t *testing.T) {
	fmt.Println("== Testing >", host.CheckCommandKubeEvent)

	kubeClient := getKubernetesClient()

	checkInterval, _ := time.ParseDuration("30s")
	clockSkew, _ := time.ParseDuration("0s")
	exitCode := kube_event.GetStatusCodeForEventCount(kubeClient, checkInterval, clockSkew)

	testDataList := []testData{
		testData{
			data: map[string]interface{}{
				"check_interval": checkInterval,
			},
			expectedCode: exitCode,
		},
	}
	for _, testData := range testDataList {
		argList := []string{
			"check_kube_event",
		}
		for key, val := range testData.data {
			argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
		}
		statusCode := execCheckCommand("hyperalert", argList...)
		assert.EqualValues(t, testData.expectedCode, statusCode)
	}
}

func TestKubeExec(t *testing.T) {
	fmt.Println("== Testing >", host.CheckCommandKubeExec)

	kubeClient := getKubernetesClient()

	testkubeExec := func(dataConfig *dataConfig) {
		// This will create object & return icinga_host name
		name, _ := getTestData(kubeClient, dataConfig)
		time.Sleep(time.Second * 30)

		objectType, objectName, namespace := plugin.GetKubeObjectInfo(name)
		objectList, err := host.GetObjectList(kubeClient.Client, host.CheckCommandKubeExec, host.HostTypePod, namespace, objectType, objectName, "")
		if err != nil {
			log.Fatal(err)
		}

		testDataList := make([]testData, 0)

		for _, object := range objectList {
			testDataList = append(testDataList, testData{
				data: map[string]interface{}{
					"host": object.Name,
					"arg":  "exit 0",
				},
				expectedCode: 0,
			})
			testDataList = append(testDataList, testData{
				data: map[string]interface{}{
					"host": object.Name,
					"arg":  "exit 5",
				},
				expectedCode: 2,
			})
		}

		for _, testData := range testDataList {
			argList := []string{
				"check_kube_exec",
			}
			for key, val := range testData.data {
				argList = append(argList, fmt.Sprintf("--%s=%v", key, val))
			}
			statusCode := execCheckCommand("hyperalert", argList...)
			assert.EqualValues(t, testData.expectedCode, statusCode)
		}
	}

	ns := "e2e-1"
	dataConfig := &dataConfig{
		Namespace: ns,
	}

	fmt.Println(">> Creating namespace", ns)
	createNewNamespace(kubeClient, ns)
	fmt.Println()

	fmt.Println(">> Testing plugings for", host.TypeReplicationcontrollers)
	dataConfig.ObjectType = host.TypeReplicationcontrollers
	testkubeExec(dataConfig)

	fmt.Println(">> Deleting namespace", ns)
	deleteNewNamespace(kubeClient, ns)

	fmt.Println()
}

func TestGeneralAlert(t *testing.T) {
	var err error
	context := &client.Context{}
	kubeClient := getKubernetesClient()
	context.KubeClient = kubeClient

	/*
		TODO: Pass Secret Name
	*/
	icingaClient, err := icinga.NewIcingaClient(kubeClient.Client, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	context.IcingaClient = icingaClient

	go runKubeD(context)
	time.Sleep(time.Minute * 1)

	ns := "e2e"
	fmt.Println(">> Creating namespace", ns)
	createNewNamespace(kubeClient, ns)
	fmt.Println()

	replicaSet := &ext.ReplicaSet{}
	replicaSet.Namespace = ns
	fmt.Println("--> Creating replicaSet")
	err = testing_lib.CreateKubernetesObject(kubeClient.Client, replicaSet)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	alert := &aci.Alert{
		ObjectMeta: kapi.ObjectMeta{
			Namespace: ns,
			Labels: map[string]string{
				"alert.appscode.com/objectType": host.TypeReplicasets,
				"alert.appscode.com/objectName": replicaSet.Name,
			},
		},
		Spec: aci.AlertSpec{
			CheckCommand: host.CheckCommandPodExists,
			IcingaParam: &aci.IcingaParam{
				CheckIntervalSec: 30,
			},
		},
	}

	err = createAlertObject(kubeClient, alert)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("--> Creating alert on replicaSet:", alert.Name)

	time.Sleep(time.Minute * 1)

	fmt.Println("--> Count alert Service")
	err = util.CountAlertService(context, alert, false)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("--> Update alert")
	err = general.UpdateAlert(kubeClient, alert)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	time.Sleep(time.Minute * 1)

	fmt.Println("--> Check alert")
	err = general.CheckAlertServiceData(context, alert)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("--> Delete alert")
	err = deleteAlertObject(kubeClient, alert)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	time.Sleep(time.Minute * 1)

	fmt.Println("--> Count alert Service")
	err = util.CountAlertService(context, alert, true)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(">> Deleting namespace", ns)
	deleteNewNamespace(kubeClient, ns)
}

func TestAcknowledge(t *testing.T) {
	os.Setenv("E2E_ICINGA_SECRET", "appscode-icinga.kube-system")

	context := &client.Context{}
	kubeClient := getKubernetesClient()
	context.KubeClient = kubeClient

	/*
		TODO: Pass Secret Name
	*/
	icingaClient, err := icinga.NewIcingaClient(kubeClient.Client, "")
	if err != nil {
		log.Fatalln(err)
	}
	context.IcingaClient = icingaClient

	go runKubeD(context)

	ns := rand.WithUniqSuffix("e2e")
	fmt.Println(">> Creating namespace", ns)
	createNewNamespace(kubeClient, ns)
	fmt.Println()

	alert := &aci.Alert{
		ObjectMeta: kapi.ObjectMeta{
			Namespace: ns,
			Labels: map[string]string{
				"alert.appscode.com/objectType": host.TypeCluster,
			},
		},
		Spec: aci.AlertSpec{
			CheckCommand: host.CheckNodeCount,
			IcingaParam: &aci.IcingaParam{
				CheckIntervalSec: 30,
			},
			Vars: map[string]interface{}{
				"count": 0,
			},
		},
	}

	err = createAlertObject(kubeClient, alert)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(">> Creating Alert", alert.Name)
	fmt.Println(">> Waiting to detect problem in Icinga")
	time.Sleep(time.Minute * 2)

	message := &types.AlertEventMessage{
		Comment:  "Acknowledged by Event",
		UserName: "shahriar",
	}
	messageByte, err := json.Marshal(message)
	if err != nil {
		log.Fatalln(err)
	}

	timestamp := unversioned.NewTime(time.Now().UTC())
	event := &kapi.Event{
		ObjectMeta: kapi.ObjectMeta{
			Name:      "acknowledgement",
			Namespace: ns,
		},
		InvolvedObject: kapi.ObjectReference{
			Kind:      events.ObjectKindAlert.String(),
			Namespace: ns,
			Name:      alert.Name,
		},
		Source: kapi.EventSource{
			Component: "searchlight",
			Host:      "node_count@" + ns,
		},
		Reason:  events.EventReasonAlertAcknowledgement.String(),
		Message: string(messageByte),
		Type:    kapi.EventTypeNormal,

		Count:          1,
		FirstTimestamp: timestamp,
		LastTimestamp:  timestamp,
	}

	fmt.Println(">> Creating event to acknowledge")
	_, err = kubeClient.Client.Core().Events(ns).Create(event)
	if err != nil {
		log.Fatalln(err)
	}

	time.Sleep(time.Minute * 1)

	objectList := []*host.KubeObjectInfo{{Name: "node_count@" + ns, IP: "127.0.0.1"}}
	in := host.IcingaServiceSearchQuery(alert.Name, objectList)
	var respService host.ResponseObject
	if _, err := icingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
		log.Fatalln(err)
	}

	if len(respService.Results) == 0 {
		log.Fatalln("No Icinga Service found")
	}

	assert.NotEqual(t, 0.0, respService.Results[0].Attrs.Acknowledgement)
}
