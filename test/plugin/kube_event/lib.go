package kube_event

import (
	"fmt"
	"os"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/fields"
)

func getStatusCodeForEventCount(kubeClient *k8s.KubeClient, checkInterval, clockSkew time.Duration) util.IcingaState {

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
			fmt.Println(err)
			os.Exit(1)
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
		fmt.Println(err)
		os.Exit(1)
	}

	for _, event := range eventList.Items {
		if checkTime.Before(event.LastTimestamp.Time) {
			count = count + 1
		}
	}

	if count > 0 {
		return util.Warning
	}
	return util.Ok
}

func GetTestData(checkInterval, clockSkew time.Duration) []plugin.TestData {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	expectedIcingaState := getStatusCodeForEventCount(kubeClient, checkInterval, clockSkew)
	testDataList := []plugin.TestData{
		plugin.TestData{
			Data: map[string]interface{}{
				"CheckInterval": checkInterval,
				"ClockSkew":     clockSkew,
			},
			ExpectedIcingaState: expectedIcingaState,
		},
	}

	return testDataList
}
