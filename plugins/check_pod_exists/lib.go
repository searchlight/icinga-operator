package check_pod_exists

import (
	"fmt"
	"os"
	"strings"

	"github.com/appscode/go/flags"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Request struct {
	Namespace string
	Count     int
	Selector  string
	PodName   string
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
	kubeClient, err := util.NewClient()
	if err != nil {
		return icinga.UNKNOWN, err
	}

	total_pod := 0
	if req.PodName != "" {
		pod, err := kubeClient.Client.CoreV1().Pods(req.Namespace).Get(req.PodName, metav1.GetOptions{})
		if err != nil {
			return icinga.UNKNOWN, err
		}
		if pod != nil {
			total_pod = 1
		}
	} else {
		podList, err := kubeClient.Client.CoreV1().Pods(req.Namespace).List(metav1.ListOptions{
			LabelSelector: req.Selector,
		},
		)
		if err != nil {
			return icinga.UNKNOWN, err
		}

		total_pod = len(podList.Items)
	}

	if isCountSet {
		if req.Count != total_pod {
			return icinga.CRITICAL, fmt.Sprintf("Found %d pod(s) instead of %d", total_pod, req.Count)
		} else {
			return icinga.OK, "Found all pods"
		}
	} else {
		if total_pod == 0 {
			return icinga.CRITICAL, "No pod found"
		} else {
			return icinga.OK, fmt.Sprintf("Found %d pods(s)", total_pod)
		}
	}
}

func NewCmd() *cobra.Command {
	var req Request
	var icingaHost string

	c := &cobra.Command{
		Use:     "check_pod_exists",
		Short:   "Check Kubernetes Pod(s)",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host")

			parts := strings.Split(icingaHost, "@")
			if len(parts) != 2 {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host.name")
				os.Exit(3)
			}

			name := parts[0]
			namespace := parts[1]
			if name != string(tapi.CheckPodExists) {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host.name")
				os.Exit(3)
			}
			req.Namespace = namespace

			isCountSet := cmd.Flag("count").Changed
			icinga.Output(CheckPodExists(&req, isCountSet))
		},
	}
	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	c.Flags().IntVarP(&req.Count, "count", "c", 0, "Number of Kubernetes Node")
	c.Flags().StringP("selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.")
	c.Flags().StringVarP(&req.PodName, "pod_name", "p", "", "Name of pod whose existence is checked")
	return c
}
