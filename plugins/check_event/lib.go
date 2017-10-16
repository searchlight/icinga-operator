package check_event

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

func CheckKubeEvent(req *Request) (icinga.State, interface{}) {
	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		return icinga.UNKNOWN, err
	}
	kubeClient := kubernetes.NewForConfigOrDie(config)

	checkTime := time.Now().Add(-(req.CheckInterval + req.ClockSkew))
	eventInfoList := make([]*eventInfo, 0)

	var objName, objNamespace, objKind, objUID *string
	if req.InvolvedObjectName != "" {
		objName = &req.InvolvedObjectName
	}
	if req.InvolvedObjectNamespace != "" {
		objNamespace = &req.InvolvedObjectNamespace
	}
	if req.InvolvedObjectKind != "" {
		objKind = &req.InvolvedObjectKind
	}
	if req.InvolvedObjectUID != "" {
		objUID = &req.InvolvedObjectUID
	}
	fs := fields.AndSelectors(
		fields.OneTermEqualSelector("type", apiv1.EventTypeWarning),
		kubeClient.CoreV1().Events(req.Namespace).GetFieldSelector(objName, objNamespace, objKind, objUID),
	)
	eventList, err := kubeClient.CoreV1().Events(req.Namespace).List(metav1.ListOptions{
		FieldSelector: fs.String(),
	})
	if err != nil {
		return icinga.UNKNOWN, err
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
		return icinga.OK, "All events look Normal"
	} else {
		output := &serviceOutput{
			Events:  eventInfoList,
			Message: fmt.Sprintf("Found %d Warning event(s)", len(eventInfoList)),
		}
		outputByte, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return icinga.UNKNOWN, err
		}
		return icinga.WARNING, string(outputByte)
	}
}

type Request struct {
	masterURL      string
	kubeconfigPath string

	Namespace     string
	CheckInterval time.Duration
	ClockSkew     time.Duration

	InvolvedObjectName      string
	InvolvedObjectNamespace string
	InvolvedObjectKind      string
	InvolvedObjectUID       string
}

func NewCmd() *cobra.Command {
	var req Request
	var icingaHost string

	cmd := &cobra.Command{
		Use:     "check_event",
		Short:   "Check kubernetes events for all namespaces",
		Example: "",

		Run: func(c *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(c, "check_interval")

			host, err := icinga.ParseHost(icingaHost)
			if err != nil {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host.name")
				os.Exit(3)
			}
			if host.Type != icinga.TypeCluster {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host type")
				os.Exit(3)
			}
			req.Namespace = host.AlertNamespace

			icinga.Output(CheckKubeEvent(&req))
		},
	}

	cmd.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")

	cmd.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	cmd.Flags().DurationVarP(&req.CheckInterval, "checkInterval", "c", time.Second*0, "Icinga check_interval in duration. [Format: 30s, 5m]")
	cmd.Flags().DurationVarP(&req.ClockSkew, "clockSkew", "s", time.Second*30, "Add skew with check_interval in duration. [Default: 30s]")

	cmd.Flags().StringVar(&req.InvolvedObjectName, "involvedObjectName", "", "Involved object name used to select events")
	cmd.Flags().StringVar(&req.InvolvedObjectNamespace, "involvedObjectNamespace", "", "Involved object namespace used to select events")
	cmd.Flags().StringVar(&req.InvolvedObjectKind, "involvedObjectKind", "", "Involved object kind used to select events")
	cmd.Flags().StringVar(&req.InvolvedObjectUID, "involvedObjectUID", "", "Involved object uid used to select events")

	return cmd
}
