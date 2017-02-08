package pod_exists

import (
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/test/plugin"
	"k8s.io/kubernetes/pkg/labels"
)

func GetPodCount(watcher *app.Watcher, namespace string) (int, error) {
	podList, err := watcher.Storage.PodStore.Pods(namespace).List(labels.Everything())
	if err != nil {
		return 0, err
	}
	return len(podList), nil
}

func GetTestData(objectType, objectName, namespace string, count int) []plugin.TestData {
	testDataList := []plugin.TestData{
		plugin.TestData{
			// To check for any pods
			Data: map[string]interface{}{
				"ObjectType": objectType,
				"ObjectName": objectName,
				"Namespace":  namespace,
			},
			ExpectedIcingaState: 0,
		},
		plugin.TestData{
			// To check for specific number of pods
			Data: map[string]interface{}{
				"ObjectType": objectType,
				"ObjectName": objectName,
				"Namespace":  namespace,
				"Count":      count,
			},
			ExpectedIcingaState: 0,
		},
		plugin.TestData{
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
