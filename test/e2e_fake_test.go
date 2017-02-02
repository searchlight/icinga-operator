package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/appscode/go/flags"
	"github.com/appscode/k8s-addons/pkg/events"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/client"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/test/fake"
	"github.com/appscode/searchlight/util"
	"github.com/stretchr/testify/assert"
	kapi "k8s.io/kubernetes/pkg/api"
)

func init() {
	flags.SetLogLevel(10)
}

func TestFakeKubeD(t *testing.T) {
	context := &client.Context{}
	kubeClient := newFakeKubeClient()
	context.KubeClient = kubeClient

	fmt.Println("--> Creating fake Secret for Icinga2 API")
	// Secret will be created with information of Icinga2 running in Docker for travisCI test
	secretMap := map[string]string{
		icinga.IcingaAPIUser: "",
		icinga.IcingaAPIPass: "",
		icinga.IcingaService: "",
	}
	fakeIcingaSecretName := fake.CreateFakeIcingaSecret(kubeClient, "default", secretMap)
	icingaClient, err := icinga.NewIcingaClient(kubeClient.Client, fakeIcingaSecretName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	context.IcingaClient = icingaClient

	fmt.Println("--> Creating 1st fake Alert on fake ReplicaSet")
	labelMap := map[string]string{
		"objectType": host.TypeReplicasets,
		"objectName": "test",
	}
	fake.CreateFakeAlert(kubeClient, "default", labelMap, host.CheckCommandVolume)

	//fake.CreateFakeReplicationController(kubeClient, "default")

	time.Sleep(time.Minute * 5)
}

func TestFakeKubernetesClient(t *testing.T) {
	// Fake KubeClient
	kubeClient := newFakeKubeClient()
	// Create Fake ReplicationController
	fakeRC := fake.CreateFakeReplicationController(kubeClient, "default")
	// Getting ReplicationController selector
	labelSelector, err := util.GetLabels(kubeClient, fakeRC.Namespace, host.TypeReplicationcontrollers, fakeRC.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	podList, err := kubeClient.Client.Core().Pods(fakeRC.Namespace).List(kapi.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	assert.EqualValues(t, fakeRC.Spec.Replicas, len(podList.Items))
}

func TestMultipleAlerts(t *testing.T) {
	context := &client.Context{}

	kubeClient := newFakeKubeClient()
	context.KubeClient = kubeClient

	// Secret will be created with information of Icinga2 running in Docker for travisCI test
	secretMap := map[string]string{
		icinga.IcingaAPIUser: "",
		icinga.IcingaAPIPass: "",
		icinga.IcingaService: "",
	}
	fakeIcingaSecretName := fake.CreateFakeIcingaSecret(kubeClient, "default", secretMap)
	icingaClient, err := icinga.NewIcingaClient(kubeClient.Client, fakeIcingaSecretName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	context.IcingaClient = icingaClient

	w := &app.Watcher{
		Watcher: acw.Watcher{
			Client:                  kubeClient.Client,
			AppsCodeExtensionClient: kubeClient.AppscodeExtensionClient,
			SyncPeriod:              time.Minute * 2,
		},
		IcingaClient: icingaClient,
	}
	// KubeD setup
	w.Run()
	fmt.Println("--> Running Fake kubeD")

	// Create Fake ReplicaSet
	fmt.Println("--> Creating fake ReplicaSet")
	fakeReplicaSet := fake.CreateFakeReplicaSet(kubeClient, "default")

	fmt.Println("--> Creating 1st fake Alert on fake ReplicaSet")
	labelMap := map[string]string{
		"objectType": host.TypeReplicasets,
		"objectName": fakeReplicaSet.Name,
	}
	fakeAlert := fake.CreateFakeAlert(kubeClient, fakeReplicaSet.Namespace, labelMap, host.CheckCommandVolume)
	// Dispatch with fakeEvent
	dispatch(w, events.Alert, events.Added, fakeAlert)

	// Check Icinga Objects for 1st Alert.
	fmt.Println("--> Checking Icinga Objects for 1st Alert")
	fake.CheckIcingaObjects(context, fakeAlert, false, false)
	fmt.Println("++> Check Successful")

	fmt.Println("--> Creating 2nd fake Alert on fake ReplicaSet")
	secondFakeAlert := fake.CreateFakeAlert(kubeClient, fakeReplicaSet.Namespace, labelMap, host.CheckCommandVolume)
	// Dispatch with fakeEvent
	dispatch(w, events.Alert, events.Added, secondFakeAlert)

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("--> Checking Icinga Objects for 2nd Alert")
	fake.CheckIcingaObjects(context, secondFakeAlert, false, false)
	fmt.Println("++> Check Successful")

	// Delete 1st Alert
	fmt.Println("--> Deleting 1st Alert")
	if err := kubeClient.AppscodeExtensionClient.Alert(fakeAlert.Namespace).Delete(fakeAlert.Name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dispatch(w, events.Alert, events.Deleted, fakeAlert)

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("--> Checking Icinga Objects for 1st Alert")
	fake.CheckIcingaObjects(context, fakeAlert, true, false)
	fmt.Println("++> Check Successful")

	// Delete 2nd Alert
	fmt.Println("--> Deleting 2nd Alert")
	if err := kubeClient.AppscodeExtensionClient.Alert(secondFakeAlert.Namespace).Delete(secondFakeAlert.Name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Dispatch with fakeEvent
	dispatch(w, events.Alert, events.Deleted, secondFakeAlert)

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("--> Checking Icinga Objects for 1st Alert")
	fake.CheckIcingaObjects(context, secondFakeAlert, true, true)
	fmt.Println("++> Check Successful")
}
