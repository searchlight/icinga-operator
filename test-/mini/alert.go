package mini

import (
	"fmt"
	"sync"
	"time"

	"github.com/appscode/go/crypto/rand"
	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type alertThirdPartyResource struct {
	once sync.Once
}

var alertResource = alertThirdPartyResource{}

func createAlertThirdPartyResource(w *controller.Controller) (err error) {
	alertResource.once.Do(
		func() {
			_, err = w.KubeClient.ExtensionsV1beta1().ThirdPartyResources().Get("alert.monitoring.appscode.com", metav1.GetOptions{})
			if err == nil {
				return
			}

			fmt.Println("== > Creating ThirdPartyResource")
			thirdPartyResource := &extensions.ThirdPartyResource{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "extensions/v1beta1",
					Kind:       "ThirdPartyResource",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "alert.monitoring.appscode.com",
				},
				Versions: []extensions.APIVersion{
					{
						Name: aci.V1alpha1SchemeGroupVersion.Version,
					},
				},
			}
			_, err = w.KubeClient.ExtensionsV1beta1().ThirdPartyResources().Create(thirdPartyResource)
			if err != nil {
				return
			}

			try := 0
			for {
				_, err = w.ExtClient.PodAlerts(apiv1.NamespaceDefault).List(metav1.ListOptions{LabelSelector: labels.Everything().String()})

				if err != nil {
					fmt.Println(err.Error())
				} else {
					return
				}
				if try > 5 {
					return
				}

				fmt.Println("--> Waiting for 30 second more in check process")
				time.Sleep(time.Second * 30)
				try++
			}
		},
	)
	return
}

func getAlert(namespace string) *aci.PodAlert {
	fakeAlert := &aci.PodAlert{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Alert",
			APIVersion: "monitoring.appscode.com/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("alert"),
			Namespace: namespace,
			Labels: map[string]string{
				"alert.appscode.com/objectType": "cluster",
			},
		},
		Spec: aci.PodAlertSpec{},
	}
	return fakeAlert
}

func CreateAlert(watcher *controller.Controller, namespace string, labelMap map[string]string, check string) (*aci.PodAlert, error) {
	// Add Alert ThirdPartyResource
	if err := createAlertThirdPartyResource(watcher); err != nil {
		return nil, err
	}

	alert := getAlert(namespace)
	alert.Spec = aci.PodAlertSpec{
		Check:         check,
		CheckInterval: metav1.Duration{Duration: 30 * time.Second},
	}

	for key, val := range labelMap {
		alert.ObjectMeta.Labels[fmt.Sprintf("monitoring.appscode.com/%s", key)] = val
	}

	// Create Fake 1st Alert
	if _, err := watcher.ExtClient.PodAlerts(alert.Namespace).Create(alert); err != nil {
		return nil, err
	}

	return alert, nil
}

func DeleteAlert(watcher *controller.Controller, alert *aci.PodAlert) error {
	// Delete Alert
	if err := watcher.ExtClient.PodAlerts(alert.Namespace).Delete(alert.Name); err != nil {
		return err
	}
	return nil
}
