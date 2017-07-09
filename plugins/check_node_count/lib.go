package check_node_count

import (
	"fmt"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type Request struct {
	Count int
}

func CheckNodeCount(req *Request) (icinga.State, interface{}) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return icinga.UNKNOWN, err
	}

	nodeList, err := kubeClient.Client.CoreV1().Nodes().List(metav1.ListOptions{
		LabelSelector: labels.Everything().String(),
	},
	)
	if err != nil {
		return icinga.UNKNOWN, err
	}

	if len(nodeList.Items) == req.Count {
		return icinga.OK, "Found all nodes"
	} else {
		return icinga.CRITICAL, fmt.Sprintf("Found %d node(s) instead of %d", len(nodeList.Items), req.Count)
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
			icinga.Output(CheckNodeCount(&req))
		},
	}

	c.Flags().IntVarP(&req.Count, "count", "c", 0, "Number of expected Kubernetes Node")
	return c
}
