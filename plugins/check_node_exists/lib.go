package check_node_exists

import (
	"fmt"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Request struct {
	masterURL      string
	kubeconfigPath string

	Selector string
	NodeName string
	Count    int
}

func CheckNodeExists(req *Request, isCountSet bool) (icinga.State, interface{}) {
	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		return icinga.UNKNOWN, err
	}
	kubeClient := kubernetes.NewForConfigOrDie(config)

	total_node := 0
	if req.NodeName != "" {
		node, err := kubeClient.CoreV1().Nodes().Get(req.NodeName, metav1.GetOptions{})
		if err != nil {
			return icinga.UNKNOWN, err
		}
		if node != nil {
			total_node = 1
		}
	} else {
		nodeList, err := kubeClient.CoreV1().Nodes().List(metav1.ListOptions{
			LabelSelector: req.Selector,
		},
		)
		if err != nil {
			return icinga.UNKNOWN, err
		}

		total_node = len(nodeList.Items)
	}

	if isCountSet {
		if req.Count != total_node {
			return icinga.CRITICAL, fmt.Sprintf("Found %d node(s) instead of %d", total_node, req.Count)
		} else {
			return icinga.OK, "Found all nodes"
		}
	} else {
		if total_node == 0 {
			return icinga.CRITICAL, "No node found"
		} else {
			return icinga.OK, fmt.Sprintf("Found %d node(s)", total_node)
		}
	}
}

func NewCmd() *cobra.Command {
	var req Request

	cmd := &cobra.Command{
		Use:     "check_node_exists",
		Short:   "Count Kubernetes Nodes",
		Example: "",

		Run: func(c *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(c, "count")
			isCountSet := c.Flag("count").Changed
			icinga.Output(CheckNodeExists(&req, isCountSet))
		},
	}

	cmd.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")

	cmd.Flags().StringVarP(&req.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.")
	cmd.Flags().StringVarP(&req.NodeName, "nodeName", "n", "", "Name of node whose existence is checked")
	cmd.Flags().IntVarP(&req.Count, "count", "c", 0, "Number of expected Kubernetes Node")
	return cmd
}
