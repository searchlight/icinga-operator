package node_count

import (
	"fmt"
	"os"

	"github.com/appscode/searchlight/pkg/client/k8s"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"github.com/appscode/searchlight/test/plugin"
)

func getKubernetesNodeCount(kubeClient *k8s.KubeClient) int {

	nodeList, err := kubeClient.Client.Core().
		Nodes().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return len(nodeList.Items)
}

func GetTestData() []plugin.TestData {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	actualNodeCount := getKubernetesNodeCount(kubeClient)

	testDataList := []plugin.TestData{
		plugin.TestData{
			Data: map[string]interface{}{
				"Count": actualNodeCount,
			},
			ExpectedIcingaState: 0,
		},
		plugin.TestData{
			Data: map[string]interface{}{
				"Count": actualNodeCount + 1,
			},
			ExpectedIcingaState: 2,
		},
	}

	return testDataList
}