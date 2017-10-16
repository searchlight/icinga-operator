package check_pod_status

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

	Host string
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
	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		return icinga.UNKNOWN, err
	}
	kubeClient := kubernetes.NewForConfigOrDie(config)

	host, err := icinga.ParseHost(req.Host)
	if err != nil {
		fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host.name")
		os.Exit(3)
	}
	if host.Type != icinga.TypePod {
		fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host type")
		os.Exit(3)
	}

	pod, err := kubeClient.CoreV1().Pods(host.AlertNamespace).Get(host.ObjectName, metav1.GetOptions{})
	if err != nil {
		return icinga.UNKNOWN, err
	}

	if ok, err := PodRunningAndReady(*pod); !ok {
		return icinga.CRITICAL, err
	}
	return icinga.OK, pod.Status.Phase
}

// ref: https://github.com/coreos/prometheus-operator/blob/c79166fcff3dae7bb8bc1e6bddc81837c2d97c04/pkg/k8sutil/k8sutil.go#L64
// PodRunningAndReady returns whether a pod is running and each container has
// passed it's ready state.
func PodRunningAndReady(pod apiv1.Pod) (bool, error) {
	switch pod.Status.Phase {
	case apiv1.PodFailed, apiv1.PodSucceeded:
		return false, fmt.Errorf("pod completed")
	case apiv1.PodRunning:
		for _, cond := range pod.Status.Conditions {
			if cond.Type != apiv1.PodReady {
				continue
			}
			return cond.Status == apiv1.ConditionTrue, nil
		}
		return false, fmt.Errorf("pod ready condition not found")
	}
	return false, nil
}

func NewCmd() *cobra.Command {
	var req Request
	c := &cobra.Command{
		Use:     "check_pod_status",
		Short:   "Check Kubernetes Pod(s) status",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host")
			icinga.Output(CheckPodStatus(&req))
		},
	}

	c.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	c.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")

	c.Flags().StringVarP(&req.Host, "host", "H", "", "Icinga host name")
	return c
}
