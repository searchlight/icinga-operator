package pod_exists

import (
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/test/plugin"
	"k8s.io/apimachinery/pkg/labels"
)

func GetPodCount(w *controller.Controller, namespace string) (int, error) {
	podList, err := w.Storage.PodStore.Pods(namespace).List(labels.Everything())
	if err != nil {
		return 0, err
	}
	return len(podList), nil
}

func GetTestData(objectType, objectName, namespace string, count int) []plugin.TestData {
	testDataList := []plugin.TestData{
		{
			// To check for any pods
			Data: map[string]interface{}{
				"ObjectType": objectType,
				"ObjectName": objectName,
				"Namespace":  namespace,
			},
			ExpectedIcingaState: 0,
		},
		{
			// To check for specific number of pods
			Data: map[string]interface{}{
				"ObjectType": objectType,
				"ObjectName": objectName,
				"Namespace":  namespace,
				"Count":      count,
			},
			ExpectedIcingaState: 0,
		},
		{
			// To check for critical when pod number mismatch
			Data: map[string]interface{}{
				"ObjectType": objectType,
				"ObjectName": objectName,
				"Namespace":  namespace,
				"Count":      count + 1,
			},
			ExpectedIcingaState: 2,
		},
	}
	return testDataList
}
