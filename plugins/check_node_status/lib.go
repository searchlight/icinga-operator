package check_node_status

import (
	"fmt"
	"os"
	"strings"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
)

type Request struct {
	Name string
}

func CheckNodeStatus(req *Request) (util.IcingaState, interface{}) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return util.Unknown, err
	}

	node, err := kubeClient.Client.Core().Nodes().Get(req.Name)
	if err != nil {
		return util.Unknown, err
	}

	if node == nil {
		return util.Critical, "Node not found"
	}

	for _, condition := range node.Status.Conditions {
		if condition.Type == kapi.NodeReady && condition.Status == kapi.ConditionFalse {
			return util.Critical, "Node is not Ready"
		}
	}

	return util.Ok, "Node is Ready"
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
				fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
				os.Exit(3)
			}

			req.Name = parts[0]
			CheckNodeStatus(&req)
		},
	}

	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	return c
}
