package check_pod_status

import (
	"encoding/json"
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
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type Request struct {
	ObjectType string
	ObjectName string
	Namespace  string
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

func CheckPodStatus(req *Request) (util.IcingaState, interface{}) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		return util.Unknown, err
	}

	objectInfoList := make([]*objectInfo, 0)
	if req.ObjectType == host.TypePods {
		pod, err := kubeClient.Client.CoreV1().Pods(req.Namespace).Get(req.ObjectName, metav1.GetOptions{})
		if err != nil {
			return util.Unknown, err
		}

		if !(pod.Status.Phase == apiv1.PodSucceeded || pod.Status.Phase == apiv1.PodRunning) {
			objectInfoList = append(objectInfoList, &objectInfo{Name: pod.Name, Status: string(pod.Status.Phase), Namespace: pod.Namespace})
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
		})
		if err != nil {
			return util.Unknown, err
		}

		for _, pod := range podList.Items {
			if !(pod.Status.Phase == apiv1.PodSucceeded || pod.Status.Phase == apiv1.PodRunning) {
				objectInfoList = append(objectInfoList, &objectInfo{Name: pod.Name, Status: string(pod.Status.Phase), Namespace: pod.Namespace})
			}
		}
	}

	if len(objectInfoList) == 0 {
		return util.Ok, "All pods are Ready"
	} else {
		output := &serviceOutput{
			Objects: objectInfoList,
			Message: fmt.Sprintf("Found %d not running pods(s)", len(objectInfoList)),
		}
		outputByte, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return util.Unknown, err
		}
		return util.Critical, outputByte
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
				fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
				os.Exit(3)
			}

			name := parts[0]
			namespace := parts[1]

			objectType := ""
			objectName := ""
			if name != host.CheckCommandPodStatus {
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

			util.Output(CheckPodStatus(&req))
		},
	}
	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	return c
}
