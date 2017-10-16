package check_node_status

import (
	"fmt"
	"os"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Request struct {
	masterURL      string
	kubeconfigPath string

	Name string
}

func CheckNodeStatus(req *Request) (icinga.State, interface{}) {
	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		return icinga.UNKNOWN, err
	}
	kubeClient := kubernetes.NewForConfigOrDie(config)

	node, err := kubeClient.CoreV1().Nodes().Get(req.Name, metav1.GetOptions{})
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

			host, err := icinga.ParseHost(icingaHost)
			if err != nil {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host.name")
				os.Exit(3)
			}
			if host.Type != icinga.TypeNode {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host type")
				os.Exit(3)
			}
			req.Name = host.ObjectName
			icinga.Output(CheckNodeStatus(&req))
		},
	}

	c.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	c.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")

	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	return c
}
