package component_status

import (
	"fmt"
	"os"

	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func GetStatusCodeForComponentStatus(kubeClient *k8s.KubeClient) util.IcingaState {
	components, err := kubeClient.Client.Core().ComponentStatuses().
		List(kapi.ListOptions{LabelSelector: labels.Everything()})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, component := range components.Items {
		for _, condition := range component.Conditions {
			if condition.Type == kapi.ComponentHealthy && condition.Status == kapi.ConditionFalse {
				return util.Critical
			}
		}
	}
	return util.Ok
}
