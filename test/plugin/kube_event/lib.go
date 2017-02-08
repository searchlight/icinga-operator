package kube_event

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/fields"
)

func getStatusCodeForEventCount(kubeClient *k8s.KubeClient, checkInterval, clockSkew time.Duration) (util.IcingaState, error) {

	now := time.Now()
	// Create some fake event
	for i := 0; i < 5; i++ {
		_, err := kubeClient.Client.Core().Events(kapi.NamespaceDefault).Create(&kapi.Event{
			ObjectMeta: kapi.ObjectMeta{
				Name: rand.WithUniqSuffix("event"),
			},
			Type:           kapi.EventTypeWarning,
			FirstTimestamp: unversioned.NewTime(now),
			LastTimestamp:  unversioned.NewTime(now),
		})
		if err != nil {
			return util.Unknown, err
		}
	}

	count := 0
	field := fields.OneTermEqualSelector(kapi.EventTypeField, kapi.EventTypeWarning)
	checkTime := time.Now().Add(-(checkInterval + clockSkew))
	eventList, err := kubeClient.Client.Core().Events(kapi.NamespaceAll).List(
		kapi.ListOptions{
			FieldSelector: field,
		},
	)
	if err != nil {
		return util.Unknown, err
	}

	for _, event := range eventList.Items {
		if checkTime.Before(event.LastTimestamp.Time) {
			count = count + 1
		}
	}

	if count > 0 {
		return util.Warning, nil
	}
	return util.Ok, nil
}

func GetTestData(kubeClient *k8s.KubeClient, checkInterval, clockSkew time.Duration) ([]plugin.TestData, error) {
	expectedIcingaState, err := getStatusCodeForEventCount(kubeClient, checkInterval, clockSkew)
	if err != nil {
		return nil, err
	}
	testDataList := []plugin.TestData{
		plugin.TestData{
			Data: map[string]interface{}{
				"CheckInterval": checkInterval,
				"ClockSkew":     clockSkew,
			},
			ExpectedIcingaState: expectedIcingaState,
		},
	}

	return testDataList, nil
}
