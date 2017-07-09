package check_component_status

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type objectInfo struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

type serviceOutput struct {
	Objects []*objectInfo `json:"objects,omitempty"`
	Message string        `json:"message,omitempty"`
}

func CheckComponentStatus() (icinga.State, interface{}) {
	kubeClient, err := util.NewClient()
	if err != nil {
		return icinga.UNKNOWN, err
	}

	components, err := kubeClient.Client.CoreV1().ComponentStatuses().List(metav1.ListOptions{
		LabelSelector: labels.Everything().String(),
	},
	)
	if err != nil {
		return icinga.UNKNOWN, err
	}

	objectInfoList := make([]*objectInfo, 0)
	for _, component := range components.Items {
		for _, condition := range component.Conditions {
			if condition.Type == apiv1.ComponentHealthy && condition.Status == apiv1.ConditionFalse {
				objectInfoList = append(objectInfoList,
					&objectInfo{
						Name:   component.Name,
						Status: "Unhealthy",
					},
				)
			}
		}
	}

	if len(objectInfoList) == 0 {
		return icinga.OK, "All components are healthy"
	} else {
		output := &serviceOutput{
			Objects: objectInfoList,
			Message: fmt.Sprintf("%d unhealthy component(s)", len(objectInfoList)),
		}
		outputByte, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return icinga.UNKNOWN, err
		}
		return icinga.CRITICAL, outputByte
	}
}

func NewCmd() *cobra.Command {
	c := &cobra.Command{
		Use:     "check_component_status",
		Short:   "Check Kubernetes Component Status",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			icinga.Output(CheckComponentStatus())

		},
	}
	return c
}
