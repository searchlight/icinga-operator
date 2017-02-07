package pod_exists

import (
	"fmt"
	"os"

	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/plugins/check_pod_exists"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/stretchr/testify/assert"
	"k8s.io/kubernetes/pkg/labels"
	"testing"
)

func GetPodCount(watcher *app.Watcher, namespace string) int {
	podList, err := watcher.Storage.PodStore.Pods(namespace).List(labels.Everything())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return len(podList)
}

func CheckPodExists(t *testing.T, objectType, objectName, namespace string, count int) {
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

	for _, testData := range testDataList {
		var req check_pod_exists.Request
		plugin.FillStruct(testData.Data, &req)
		isCountSet := false
		if req.Count != 0 {
			isCountSet = true
		}
		icingaState, _ := check_pod_exists.CheckPodExists(&req, isCountSet)
		assert.EqualValues(t, testData.ExpectedIcingaState, icingaState)
	}
}
