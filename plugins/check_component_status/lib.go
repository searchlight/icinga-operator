package check_component_status

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type Request struct {
	Selector      string
	ComponentName string
}

type objectInfo struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

type serviceOutput struct {
	Objects []*objectInfo `json:"objects,omitempty"`
	Message string        `json:"message,omitempty"`
}

func CheckComponentStatus(req *Request) (icinga.State, interface{}) {
	kubeClient, err := util.NewClient()
	if err != nil {
		return icinga.UNKNOWN, err
	}

	var components []apiv1.ComponentStatus
	if req.ComponentName != "" {
		comp, err := kubeClient.Client.CoreV1().ComponentStatuses().Get(req.ComponentName, metav1.GetOptions{})
		if err != nil {
			return icinga.UNKNOWN, err
		}
		components = []apiv1.ComponentStatus{*comp}
	} else {
		comps, err := kubeClient.Client.CoreV1().ComponentStatuses().List(metav1.ListOptions{
			LabelSelector: req.Selector,
		})
		if err != nil {
			return icinga.UNKNOWN, err
		}
		components = comps.Items
	}

	objectInfoList := make([]*objectInfo, 0)
	for _, component := range components {
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
	var req Request

	cmd := &cobra.Command{
		Use:     "check_component_status",
		Short:   "Check Kubernetes Component Status",
		Example: "",

		Run: func(c *cobra.Command, args []string) {
			icinga.Output(CheckComponentStatus(&req))
		},
	}
	cmd.Flags().StringVarP(&req.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.")
	cmd.Flags().StringVarP(&req.ComponentName, "componentName", "n", "", "Name of component which should be ready")
	return cmd
}
