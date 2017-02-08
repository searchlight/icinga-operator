package pod_status

import (
	"fmt"
	"os"

	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func getStatusCodeForPodStatus(watcher *app.Watcher, objectType, objectName, namespace string) util.IcingaState {
	var err error
	if objectType == host.TypePods {
		pod, err := watcher.Storage.PodStore.Pods(namespace).Get(objectName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !(pod.Status.Phase == kapi.PodSucceeded || pod.Status.Phase == kapi.PodRunning) {
			return util.Critical
		}

	} else {
		labelSelector := labels.Everything()
		if objectType != "" {
			labelSelector, err = util.GetLabels(watcher.Client, namespace, objectType, objectName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		podList, err := watcher.Storage.PodStore.Pods(namespace).List(labelSelector)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, pod := range podList {
			if !(pod.Status.Phase == kapi.PodSucceeded || pod.Status.Phase == kapi.PodRunning) {
				return util.Critical
			}
		}
	}
	return util.Ok
}

func GetTestData(watcher *app.Watcher, objectType, objectName, namespace string) []plugin.TestData {
	testDataList := []plugin.TestData{
		plugin.TestData{
			Data: map[string]interface{}{
				"ObjectType": objectType,
				"ObjectName": objectName,
				"Namespace":  namespace,
			},
			ExpectedIcingaState: getStatusCodeForPodStatus(watcher, objectType, objectName, namespace),
		},
	}
	return testDataList
}
