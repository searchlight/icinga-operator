package kube_event

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/appscode/searchlight/test/plugin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/pkg/api"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func getStatusCodeForEventCount(kubeClient *util.KubeClient, checkInterval, clockSkew time.Duration) (icinga.State, error) {

	now := time.Now()
	// Create some fake event
	for i := 0; i < 5; i++ {
		_, err := kubeClient.Client.CoreV1().Events(apiv1.NamespaceDefault).Create(&apiv1.Event{
			ObjectMeta: metav1.ObjectMeta{
				Name: rand.WithUniqSuffix("event"),
			},
			Type:           apiv1.EventTypeWarning,
			FirstTimestamp: metav1.NewTime(now),
			LastTimestamp:  metav1.NewTime(now),
		})
		if err != nil {
			return icinga.UNKNOWN, err
		}
	}

	count := 0
	field := fields.OneTermEqualSelector(api.EventTypeField, apiv1.EventTypeWarning)
	checkTime := time.Now().Add(-(checkInterval + clockSkew))
	eventList, err := kubeClient.Client.CoreV1().Events(apiv1.NamespaceAll).List(metav1.ListOptions{
		FieldSelector: field.String(),
	})
	if err != nil {
		return icinga.UNKNOWN, err
	}

	for _, event := range eventList.Items {
		if checkTime.Before(event.LastTimestamp.Time) {
			count = count + 1
		}
	}

	if count > 0 {
		return icinga.WARNING, nil
	}
	return icinga.OK, nil
}

func GetTestData(kubeClient *util.KubeClient, checkInterval, clockSkew time.Duration) ([]plugin.TestData, error) {
	expectedIcingaState, err := getStatusCodeForEventCount(kubeClient, checkInterval, clockSkew)
	if err != nil {
		return nil, err
	}
	testDataList := []plugin.TestData{
		{
			Data: map[string]interface{}{
				"CheckInterval": checkInterval,
				"ClockSkew":     clockSkew,
			},
			ExpectedIcingaState: expectedIcingaState,
		},
	}

	return testDataList, nil
}
