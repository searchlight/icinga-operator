package check_pod_exists

import (
	"fmt"
	"os"
	"strings"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type Request struct {
	ObjectType string
	ObjectName string
	Namespace  string
	Count      int
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

func CheckPodExists(req *Request, isCountSet bool) (util.IcingaState, interface{}) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return util.Unknown, err
	}

	total_pod := 0
	if req.ObjectType == host.TypePods {
		pod, err := kubeClient.Client.CoreV1().Pods(req.Namespace).Get(req.ObjectName, metav1.GetOptions{})
		if err != nil {
			return util.Unknown, err
		}
		if pod != nil {
			total_pod = 1
		}
	} else {
		labelSelector := labels.Everything()
		if req.ObjectType != "" {
			if labelSelector, err = util.GetLabels(kubeClient.Client, req.Namespace, req.ObjectType, req.ObjectName); err != nil {
				return util.Unknown, err
			}
		}

		podList, err := kubeClient.Client.CoreV1().Pods(req.Namespace).List(metav1.ListOptions{
			LabelSelector: labelSelector.String(),
		},
		)
		if err != nil {
			return util.Unknown, err
		}

		total_pod = len(podList.Items)
	}

	if isCountSet {
		if req.Count != total_pod {
			return util.Critical, fmt.Sprintf("Found %d pod(s) instead of %d", total_pod, req.Count)
		} else {
			return util.Ok, "Found all pods"
		}
	} else {
		if total_pod == 0 {
			return util.Critical, "No pod found"
		} else {
			return util.Ok, fmt.Sprintf("Found %d pods(s)", total_pod)
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
				fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
				os.Exit(3)
			}

			name := parts[0]
			namespace := parts[1]

			objectType := ""
			objectName := ""
			if name != host.CheckCommandPodExists {
				parts = strings.Split(name, "|")
				if len(parts) == 1 {
					objectType = host.TypePods
					objectName = parts[0]
				} else if len(parts) == 2 {
					objectType = parts[0]
					objectName = parts[1]
				} else {
					fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
					os.Exit(3)
				}
			}

			req.ObjectType = objectType
			req.ObjectName = objectName
			req.Namespace = namespace

			isCountSet := cmd.Flag("count").Changed
			util.Output(CheckPodExists(&req, isCountSet))
		},
	}
	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	c.Flags().IntVarP(&req.Count, "count", "c", 0, "Number of Kubernetes Node")
	return c
}
