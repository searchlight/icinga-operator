package node_status

import (
	"errors"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/test/plugin"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func getKubernetesNodeName(kubeClient *k8s.KubeClient) (string, error) {
	nodeList, err := kubeClient.Client.Core().Nodes().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
	if err != nil {
		return "", err
	}

	if len(nodeList.Items) == 0 {
		return "", errors.New("No node found")
	}
	return nodeList.Items[0].Name, nil
}

func GetTestData(kubeClient *k8s.KubeClient) ([]plugin.TestData, error) {
	actualNodeName, err := getKubernetesNodeName(kubeClient)
	if err != nil {
		return nil, err
	}

	testDataList := []plugin.TestData{
		{
			Data: map[string]interface{}{
				"Name": actualNodeName,
			},
			ExpectedIcingaState: 0,
		},
		{
			Data: map[string]interface{}{
				// make node name invalid using random 2 character.
				"Name": actualNodeName + rand.Characters(2),
			},
			ExpectedIcingaState: 3,
		},
	}
	return testDataList, nil
}
