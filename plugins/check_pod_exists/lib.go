package check_pod_exists

import (
	"fmt"
	"os"

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

	Namespace string
	Selector  string
	PodName   string
	Count     int
}

type objectInfo struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Status    string `json:"status,omitempty"`
}

type serviceOutput struct {
	Objects []*objectInfo `json:"objects,omitempty"`
	Message string        `json:"message,omitempty"`
}

func CheckPodExists(req *Request, isCountSet bool) (icinga.State, interface{}) {
	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		return icinga.Unknown, err
	}
	kubeClient := kubernetes.NewForConfigOrDie(config)

	total_pod := 0
	if req.PodName != "" {
		_, err := kubeClient.CoreV1().Pods(req.Namespace).Get(req.PodName, metav1.GetOptions{})
		if err != nil {
			return icinga.Unknown, err
		}
		total_pod = 1
	} else {
		podList, err := kubeClient.CoreV1().Pods(req.Namespace).List(metav1.ListOptions{
			LabelSelector: req.Selector,
		})
		if err != nil {
			return icinga.Unknown, err
		}
		total_pod = len(podList.Items)
	}

	if isCountSet {
		if req.Count != total_pod {
			return icinga.Critical, fmt.Sprintf("Found %d pod(s) instead of %d", total_pod, req.Count)
		} else {
			return icinga.OK, "Found all pods"
		}
	} else {
		if total_pod == 0 {
			return icinga.Critical, "No pod found"
		} else {
			return icinga.OK, fmt.Sprintf("Found %d pods(s)", total_pod)
		}
	}
}

func NewCmd() *cobra.Command {
	var req Request
	var icingaHost string

	cmd := &cobra.Command{
		Use:   "check_pod_exists",
		Short: "Check Kubernetes Pod(s)",
		Run: func(c *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(c, "host")

			host, err := icinga.ParseHost(icingaHost)
			if err != nil {
				fmt.Fprintln(os.Stdout, icinga.Warning, "Invalid icinga host.name")
				os.Exit(3)
			}
			if host.Type != icinga.TypeCluster {
				fmt.Fprintln(os.Stdout, icinga.Warning, "Invalid icinga host type")
				os.Exit(3)
			}
			req.Namespace = host.AlertNamespace

			isCountSet := c.Flag("count").Changed
			icinga.Output(CheckPodExists(&req, isCountSet))
		},
	}
	cmd.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	cmd.Flags().StringVarP(&req.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.")
	cmd.Flags().StringVarP(&req.PodName, "podName", "p", "", "Name of pod whose existence is checked")
	cmd.Flags().IntVarP(&req.Count, "count", "c", 0, "Number of Kubernetes pods")
	return cmd
}
