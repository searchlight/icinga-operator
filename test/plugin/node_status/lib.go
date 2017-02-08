package node_status

import (
	"fmt"
	"os"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/test/plugin"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func getKubernetesNodeName(kubeClient *k8s.KubeClient) string {
	nodeList, err := kubeClient.Client.Core().Nodes().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(nodeList.Items) == 0 {
		fmt.Println("No node found")
		os.Exit(1)
	}
	return nodeList.Items[0].Name
}

func GetTestData() []plugin.TestData {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	actualNodeName := getKubernetesNodeName(kubeClient)

	testDataList := []plugin.TestData{
		plugin.TestData{
			Data: map[string]interface{}{
				"Name": actualNodeName,
			},
			ExpectedIcingaState: 0,
		},
		plugin.TestData{
			Data: map[string]interface{}{
				// make node name invalid using random 2 character.
				"Name": actualNodeName + rand.Characters(2),
			},
			ExpectedIcingaState: 3,
		},
	}
	return testDataList
}
