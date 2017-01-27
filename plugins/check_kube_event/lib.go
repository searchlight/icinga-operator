package check_kube_event

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/fields"
)

type eventInfo struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Kind      string `json:"kind,omitempty"`
	Count     int32  `json:"count,omitempty"`
	Reason    string `json:"reason,omitempty"`
	Message   string `json:"message,omitempty"`
}

type serviceOutput struct {
	Events  []*eventInfo `json:"events,omitempty"`
	Message string       `json:"message,omitempty"`
}

func checkKubeEvent(req *request) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	eventInfoList := make([]*eventInfo, 0)
	field := fields.OneTermEqualSelector(kapi.EventTypeField, kapi.EventTypeWarning)

	checkTime := time.Now().Add(-(req.checkInterval + req.clockSkew))

	eventList, err := kubeClient.Client.Core().Events(kapi.NamespaceAll).List(
		kapi.ListOptions{
			FieldSelector: field,
		},
	)
	if err != nil {
		fmt.Fprintln(os.Stdout, util.State[3], err)
		os.Exit(3)
	}

	for _, event := range eventList.Items {
		if checkTime.Before(event.LastTimestamp.Time) {
			eventInfoList = append(eventInfoList,
				&eventInfo{
					Name:      event.InvolvedObject.Name,
					Namespace: event.InvolvedObject.Namespace,
					Kind:      event.InvolvedObject.Kind,
					Count:     event.Count,
					Reason:    event.Reason,
					Message:   event.Message,
				},
			)
		}
	}

	if len(eventInfoList) == 0 {
		fmt.Fprintln(os.Stdout, util.State[0], "All events look Normal")
		os.Exit(0)
	} else {
		output := &serviceOutput{
			Events:  eventInfoList,
			Message: fmt.Sprintf("Found %d Warning event(s)", len(eventInfoList)),
		}
		outputByte, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			fmt.Fprintln(os.Stdout, util.State[3], err)
			os.Exit(3)
		}
		fmt.Fprintln(os.Stdout, util.State[1], string(outputByte))
		os.Exit(1)
	}
}

type request struct {
	checkInterval time.Duration
	clockSkew     time.Duration
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "check_kube_event",
		Short:   "Check kubernetes events for all namespaces",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "check_interval")
			checkKubeEvent(&req)
		},
	}

	c.Flags().DurationVarP(&req.checkInterval, "check_interval", "c", time.Second*0, "Icinga check_interval in duration. [Format: 30s, 5m]")
	c.Flags().DurationVarP(&req.clockSkew, "clock_skew", "s", time.Second*30, "Add skew with check_interval in duration. [Default: 30s]")
	return c
}
