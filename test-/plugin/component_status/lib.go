package component_status

import (
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func GetStatusCodeForComponentStatus(kubeClient *util.KubeClient) (icinga.State, error) {
	components, err := kubeClient.Client.CoreV1().ComponentStatuses().List(metav1.ListOptions{LabelSelector: labels.Everything().String()})
	if err != nil {
		return icinga.UNKNOWN, err
	}

	for _, component := range components.Items {
		for _, condition := range component.Conditions {
			if condition.Type == apiv1.ComponentHealthy && condition.Status == apiv1.ConditionFalse {
				return icinga.CRITICAL, nil
			}
		}
	}
	return icinga.OK, nil
}
