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
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

type request struct {
	host string
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

func checkPodStatus(namespace, objectType, objectName string) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	objectInfoList := make([]*objectInfo, 0)
	if objectType == host.TypePods {
		pod, err := kubeClient.Client.Core().Pods(namespace).Get(objectName)
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}

		if !(pod.Status.Phase == kapi.PodSucceeded || pod.Status.Phase == kapi.PodRunning) {
			objectInfoList = append(objectInfoList, &objectInfo{Name: pod.Name, Status: string(pod.Status.Phase), Namespace: pod.Namespace})
		}
	} else {
		labelSelector := labels.Everything()
		if objectType != "" {
			if labelSelector, err = util.GetLabels(kubeClient.Client, namespace, objectType, objectName); err != nil {
				fmt.Fprintln(os.Stdout, util.State[3], err)
				os.Exit(3)
			}
		}

		podList, err := kubeClient.Client.Core().
			Pods(namespace).List(
			kapi.ListOptions{
				LabelSelector: labelSelector,
			},
		)
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}

		for _, pod := range podList.Items {
			if !(pod.Status.Phase == kapi.PodSucceeded || pod.Status.Phase == kapi.PodRunning) {
				objectInfoList = append(objectInfoList, &objectInfo{Name: pod.Name, Status: string(pod.Status.Phase), Namespace: pod.Namespace})
			}
		}
	}

	if len(objectInfoList) == 0 {
		fmt.Fprintln(os.Stdout, util.State[0], "All pods are Ready")
		os.Exit(0)
	} else {
		output := &serviceOutput{
			Objects: objectInfoList,
			Message: fmt.Sprintf("Found %d not running pods(s)", len(objectInfoList)),
		}
		outputByte, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}
		fmt.Fprintln(os.Stdout, util.State[2], string(outputByte))
		os.Exit(2)
	}
}

func NewCmd() *cobra.Command {
	var req request
	c := &cobra.Command{
		Use:     "check_pod_status",
		Short:   "Check Kubernetes Pod(s) status",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host")

			parts := strings.Split(req.host, "@")
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

			checkPodStatus(namespace, objectType, objectName)
		},
	}
	c.Flags().StringVarP(&req.host, "host", "H", "", "Icinga host name")
	return c
}
