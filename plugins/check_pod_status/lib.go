package check_pod_status

import (
	"fmt"

	"errors"
	"github.com/appscode/go/flags"
	"github.com/appscode/kutil/tools/clientcmd"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type plugin struct {
	client  corev1.PodInterface
	options options
}

var _ plugins.PluginInterface = &plugin{}

func newPlugin(client corev1.PodInterface, opts options) *plugin {
	return &plugin{client, opts}
}

func newPluginFromConfig(opts options) (*plugin, error) {
	client, err := clientcmd.ClientFromContext(opts.kubeconfigPath, opts.contextName)
	if err != nil {
		return nil, err
	}
	return newPlugin(client.CoreV1().Pods(opts.namespace), opts), nil
}

type options struct {
	kubeconfigPath string
	contextName    string
	// options
	podName   string
	namespace string
	// IcingaHost
	host *icinga.IcingaHost
}

func (o *options) complete(cmd *cobra.Command) error {
	hostname, err := cmd.Flags().GetString(plugins.FlagHost)
	if err != nil {
		return err
	}
	o.host, err = icinga.ParseHost(hostname)
	if err != nil {
		return errors.New("invalid icinga host.name")
	}
	o.podName = o.host.ObjectName
	o.namespace = o.host.AlertNamespace

	o.kubeconfigPath, err = cmd.Flags().GetString(plugins.FlagKubeConfig)
	if err != nil {
		return err
	}
	o.contextName, err = cmd.Flags().GetString(plugins.FlagKubeConfigContext)
	if err != nil {
		return err
	}
	return nil
}

func (o *options) validate() error {
	if o.host.Type != icinga.TypePod {
		return errors.New("invalid icinga host type")
	}
	return nil
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

func (p *plugin) Check() (icinga.State, interface{}) {
	opts := p.options
	host := opts.host

	pod, err := p.client.Get(host.ObjectName, metav1.GetOptions{})
	if err != nil {
		return icinga.Unknown, err
	}

	if ok, err := p.podRunningAndReady(*pod); !ok {
		return icinga.Critical, err
	}
	return icinga.OK, pod.Status.Phase
}

// ref: https://github.com/coreos/prometheus-operator/blob/c79166fcff3dae7bb8bc1e6bddc81837c2d97c04/pkg/k8sutil/k8sutil.go#L64
// podRunningAndReady returns whether a pod is running and each container has
// passed it's ready state.
func (p *plugin) podRunningAndReady(pod core.Pod) (bool, error) {
	switch pod.Status.Phase {
	case core.PodFailed, core.PodSucceeded:
		return false, fmt.Errorf("pod completed")
	case core.PodRunning:
		for _, cond := range pod.Status.Conditions {
			if cond.Type != core.PodReady {
				continue
			}
			return cond.Status == core.ConditionTrue, nil
		}
		return false, fmt.Errorf("pod ready condition not found")
	}
	return false, nil
}

func NewCmd() *cobra.Command {
	var opts options
	c := &cobra.Command{
		Use:   "check_pod_status",
		Short: "Check Kubernetes Pod(s) status",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, plugins.FlagHost)

			if err := opts.complete(cmd); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			if err := opts.validate(); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			plugin, err := newPluginFromConfig(opts)
			if err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			icinga.Output(plugin.Check())
		},
	}

	c.Flags().StringP(plugins.FlagHost, "H", "", "Icinga host name")
	return c
}
