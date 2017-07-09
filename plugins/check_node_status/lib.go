package check_node_status

import (
	"fmt"
	"os"
	"strings"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type Request struct {
	Name string
}

func CheckNodeStatus(req *Request) (icinga.State, interface{}) {
	kubeClient, err := util.NewClient()
	if err != nil {
		return icinga.UNKNOWN, err
	}

	node, err := kubeClient.Client.CoreV1().Nodes().Get(req.Name, metav1.GetOptions{})
	if err != nil {
		return icinga.UNKNOWN, err
	}

	if node == nil {
		return icinga.CRITICAL, "Node not found"
	}

	for _, condition := range node.Status.Conditions {
		if condition.Type == apiv1.NodeReady && condition.Status == apiv1.ConditionFalse {
			return icinga.CRITICAL, "Node is not Ready"
		}
	}

	return icinga.OK, "Node is Ready"
}

func NewCmd() *cobra.Command {
	var req Request
	var icingaHost string
	c := &cobra.Command{
		Use:     "check_node_status",
		Short:   "Check Kubernetes Node",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host")
			parts := strings.Split(icingaHost, "@")
			if len(parts) != 2 {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host.name")
				os.Exit(3)
			}

			req.Name = parts[0]
			icinga.Output(CheckNodeStatus(&req))
		},
	}

	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	return c
}
