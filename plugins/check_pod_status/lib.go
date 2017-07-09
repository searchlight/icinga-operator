package check_pod_status

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/appscode/go/flags"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type Request struct {
	Namespace string
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

func CheckPodStatus(req *Request) (icinga.State, interface{}) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return icinga.UNKNOWN, err
	}

	objectInfoList := make([]*objectInfo, 0)
	if req.PodName != "" {
		pod, err := kubeClient.Client.CoreV1().Pods(req.Namespace).Get(req.PodName, metav1.GetOptions{})
		if err != nil {
			return icinga.UNKNOWN, err
		}

		if !(pod.Status.Phase == apiv1.PodSucceeded || pod.Status.Phase == apiv1.PodRunning) {
			objectInfoList = append(objectInfoList, &objectInfo{Name: pod.Name, Status: string(pod.Status.Phase), Namespace: pod.Namespace})
		}
	} else {
		podList, err := kubeClient.Client.CoreV1().Pods(req.Namespace).List(metav1.ListOptions{LabelSelector: req.Selector})
		if err != nil {
			return icinga.UNKNOWN, err
		}

		for _, pod := range podList.Items {
			if !(pod.Status.Phase == apiv1.PodSucceeded || pod.Status.Phase == apiv1.PodRunning) {
				objectInfoList = append(objectInfoList, &objectInfo{Name: pod.Name, Status: string(pod.Status.Phase), Namespace: pod.Namespace})
			}
		}
	}

	if len(objectInfoList) == 0 {
		return icinga.OK, "All pods are Ready"
	} else {
		output := &serviceOutput{
			Objects: objectInfoList,
			Message: fmt.Sprintf("Found %d not running pods(s)", len(objectInfoList)),
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
	var icingaHost string

	c := &cobra.Command{
		Use:     "check_pod_status",
		Short:   "Check Kubernetes Pod(s) status",
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
			if name != string(tapi.CheckPodStatus) {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host.name")
				os.Exit(3)
			}
			req.Namespace = namespace

			icinga.Output(CheckPodStatus(&req))
		},
	}
	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.")
	c.Flags().StringVarP(&req.PodName, "pod_name", "p", "", "Name of pod whose status is checked")
	return c
}
