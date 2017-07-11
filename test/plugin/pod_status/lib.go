package pod_status

import (
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/test/plugin"
	"k8s.io/apimachinery/pkg/labels"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func getStatusCodeForPodStatus(w *controller.Controller, objectType, objectName, namespace string) (icinga.State, error) {
	var err error
	if objectType == icinga.TypePod {
		pod, err := w.Storage.PodStore.Pods(namespace).Get(objectName)
		if err != nil {
			return icinga.UNKNOWN, err
		}
		if !(pod.Status.Phase == apiv1.PodSucceeded || pod.Status.Phase == apiv1.PodRunning) {
			return icinga.CRITICAL, nil
		}

	} else {
		labelSelector := labels.Everything()
		if objectType != "" {
			labelSelector, err = icinga.GetLabels(w.KubeClient, namespace, objectType, objectName)
			if err != nil {
				return icinga.UNKNOWN, err
			}
		}

		podList, err := w.Storage.PodStore.Pods(namespace).List(labelSelector)
		if err != nil {
			return icinga.UNKNOWN, err
		}

		for _, pod := range podList {
			if !(pod.Status.Phase == apiv1.PodSucceeded || pod.Status.Phase == apiv1.PodRunning) {
				return icinga.CRITICAL, nil
			}
		}
	}
	return icinga.OK, nil
}

func GetTestData(watcher *controller.Controller, objectType, objectName, namespace string) ([]plugin.TestData, error) {
	expectedStatusCode, err := getStatusCodeForPodStatus(watcher, objectType, objectName, namespace)
	if err != nil {
		return nil, err
	}
	testDataList := []plugin.TestData{
		{
			Data: map[string]interface{}{
				"ObjectType": objectType,
				"ObjectName": objectName,
				"Namespace":  namespace,
			},
			ExpectedIcingaState: expectedStatusCode,
		},
	}
	return testDataList, nil
}
