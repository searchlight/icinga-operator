package check_node_count

import (
	"fmt"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

type Request struct {
	Count int
}

func CheckNodeCount(req *Request) (util.IcingaState, interface{}) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return util.Unknown, err
	}

	nodeList, err := kubeClient.Client.Core().
		Nodes().List(
		kapi.ListOptions{
			LabelSelector: labels.Everything(),
		},
	)
	if err != nil {
		return util.Unknown, err
	}

	if len(nodeList.Items) == req.Count {
		return util.Ok, "Found all nodes"
	} else {
		return util.Critical, fmt.Sprintf("Found %d node(s) instead of %d", len(nodeList.Items), req.Count)
	}
}

func NewCmd() *cobra.Command {
	var req Request

	c := &cobra.Command{
		Use:     "check_node_count",
		Short:   "Count Kubernetes Nodes",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "count")
			util.Output(CheckNodeCount(&req))
		},
	}

	c.Flags().IntVarP(&req.Count, "count", "c", 0, "Number of expected Kubernetes Node")
	return c
}
