package mini

import (
	"fmt"
	"sync"

	"github.com/appscode/go/crypto/rand"
	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type alertThirdPartyResource struct {
	once sync.Once
}

var alertResource = alertThirdPartyResource{}

func createAlertThirdPartyResource(watcher *app.Watcher) (err error) {
	alertResource.once.Do(
		func() {
			_, err = watcher.Client.Extensions().ThirdPartyResources().Get("alert.appscode.com")
			if err == nil {
				return
			}

			fmt.Println("== > Creating ThirdPartyResource")
			thirdPartyResource := &extensions.ThirdPartyResource{
				TypeMeta: unversioned.TypeMeta{
					APIVersion: "extensions/v1beta1",
					Kind:       "ThirdPartyResource",
				},
				ObjectMeta: kapi.ObjectMeta{
					Name: "alert.appscode.com",
				},
				Versions: []extensions.APIVersion{
					extensions.APIVersion{
						Name: "v1beta1",
					},
				},
			}
			_, err = watcher.Client.Extensions().ThirdPartyResources().Create(thirdPartyResource)
			if err != nil {
				return
			}

			_, err = watcher.Client.Extensions().ThirdPartyResources().Get("alert.appscode.com")
			if err != nil {
				return
			}
		},
	)
	return
}

func getAlert(namespace string) *aci.Alert {
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

func CreateAlert(watcher *app.Watcher, namespace string, labelMap map[string]string, checkCommand string) (*aci.Alert, error) {
	// Add Alert ThirdPartyResource
	if err := createAlertThirdPartyResource(watcher); err != nil {
		return nil, err
	}

	alert := getAlert(namespace)
	alert.Spec = aci.AlertSpec{
		CheckCommand: checkCommand,
		IcingaParam: &aci.IcingaParam{
			CheckIntervalSec: 30,
		},
	}

	for key, val := range labelMap {
		alert.ObjectMeta.Labels[fmt.Sprintf("alert.appscode.com/%s", key)] = val
	}

	// Create Fake 1st Alert
	if _, err := watcher.AppsCodeExtensionClient.Alert(alert.Namespace).Create(alert); err != nil {
		return nil, err
	}

	return alert, nil
}

func DeleteAlert(watcher *app.Watcher, alert *aci.Alert) error {
	// Delete Alert
	if err := watcher.AppsCodeExtensionClient.Alert(alert.Namespace).Delete(alert.Name); err != nil {
		return err
	}
	return nil
}
