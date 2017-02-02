package fake

import (
	"fmt"
	"github.com/appscode/searchlight/pkg/client/k8s"

	"github.com/appscode/go/crypto/rand"
	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/searchlight/pkg/client"
	"github.com/appscode/searchlight/pkg/controller/host"
	testutil "github.com/appscode/searchlight/test/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"os"
)

func getFakeAlert(namespace string) *aci.Alert {
	fakeAlert := &aci.Alert{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Alert",
			APIVersion: "appscode.com/v1beta1",
		},
		ObjectMeta: kapi.ObjectMeta{
			Name:      rand.WithUniqSuffix("alert"),
			Namespace: namespace,
			Labels: map[string]string{
				"alert.appscode.com/objectType": "cluster",
			},
		},
		Spec: aci.AlertSpec{},
	}
	return fakeAlert
}

func CreateFakeAlert(kubeClient *k8s.KubeClient, namespace string, labelMap map[string]string, checkCommand string) *aci.Alert {
	fakeAlert := getFakeAlert(namespace)
	fakeAlert.Spec = aci.AlertSpec{
		CheckCommand: checkCommand,
		IcingaParam: &aci.IcingaParam{
			CheckIntervalSec: 30,
		},
	}

	for key, val := range labelMap {
		fakeAlert.ObjectMeta.Labels[fmt.Sprintf("alert.appscode.com/%s", key)] = val
	}

	// Create Fake 1st Alert
	if _, err := kubeClient.AppscodeExtensionClient.Alert(fakeAlert.Namespace).Create(fakeAlert); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return fakeAlert
}

func CheckIcingaObjects(context *client.Context, fakeAlert *aci.Alert, expectZeroService, expectZeroHost bool) {
	// Count Icinga Service for 1st Alert. Should be found
	fmt.Println("--> Counting Icinga Service")
	if err := testutil.CountIcingaService(context, fakeAlert, expectZeroService); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	objectType, objectName := host.GetObjectInfo(fakeAlert.Labels)

	// Count Icinga Host in Icinga2. Should be found
	fmt.Println("--> Counting Icinga Host")
	if err := testutil.CountIcingaHost(context, fakeAlert.Spec.CheckCommand, fakeAlert.Namespace, objectType, objectName, expectZeroHost); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
