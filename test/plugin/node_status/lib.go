package node_status

import (
	"errors"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/appscode/searchlight/test/plugin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func getKubernetesNodeName(kubeClient *util.KubeClient) (string, error) {
	nodeList, err := kubeClient.Client.CoreV1().Nodes().List(metav1.ListOptions{
		LabelSelector: labels.Everything().String(),
	})
	if err != nil {
		return "", err
	}

	if len(nodeList.Items) == 0 {
		return "", errors.New("No node found")
	}
	return nodeList.Items[0].Name, nil
}

func GetTestData(kubeClient *util.KubeClient) ([]plugin.TestData, error) {
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
