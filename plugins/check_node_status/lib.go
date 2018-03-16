package check_node_status

import (
	"encoding/json"
	"errors"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type plugin struct {
	client  corev1.NodeInterface
	options options
}

var _ plugins.PluginInterface = &plugin{}

func newPlugin(client corev1.NodeInterface, opts options) *plugin {
	return &plugin{client, opts}
}

func newPluginFromConfig(opts options) (*plugin, error) {
	config, err := clientcmd.BuildConfigFromFlags(opts.masterURL, opts.kubeconfigPath)
	if err != nil {
		return nil, err
	}
	client := kubernetes.NewForConfigOrDie(config).CoreV1().Nodes()
	return newPlugin(client, opts), nil
}

type options struct {
	masterURL      string
	kubeconfigPath string
	// Icinga host name
	hostname string
	// options for Secret
	nodeName string
}

func (o *options) validate() error {
	host, err := icinga.ParseHost(o.hostname)
	if err != nil {
		return errors.New("invalid icinga host.name")
	}
	if host.Type != icinga.TypeNode {
		return errors.New("invalid icinga host type")
	}
	o.nodeName = host.ObjectName
	return nil
}

type message struct {
	Ready              core.ConditionStatus `json:"ready,omitempty"`
	OutOfDisk          core.ConditionStatus `json:"outOfDisk,omitempty"`
	MemoryPressure     core.ConditionStatus `json:"memoryPressure,omitempty"`
	DiskPressure       core.ConditionStatus `json:"diskPressure,omitempty"`
	NetworkUnavailable core.ConditionStatus `json:"networkUnavailable,omitempty"`
}

func (p *plugin) Check() (icinga.State, interface{}) {
	node, err := p.client.Get(p.options.nodeName, metav1.GetOptions{})
	if err != nil {
		return icinga.Unknown, err
	}

	msg := message{}
	for _, condition := range node.Status.Conditions {
		switch condition.Type {
		case core.NodeReady:
			msg.Ready = condition.Status
		case core.NodeOutOfDisk:
			msg.OutOfDisk = condition.Status
		case core.NodeMemoryPressure:
			msg.MemoryPressure = condition.Status
		case core.NodeDiskPressure:
			msg.DiskPressure = condition.Status
		case core.NodeNetworkUnavailable:
			msg.NetworkUnavailable = condition.Status
		}
	}

	var state icinga.State
	if msg.Ready == core.ConditionFalse {
		state = icinga.Critical
	} else if msg.OutOfDisk == core.ConditionTrue ||
		msg.MemoryPressure == core.ConditionTrue ||
		msg.DiskPressure == core.ConditionTrue ||
		msg.NetworkUnavailable == core.ConditionTrue {
		state = icinga.Critical
	}

	output, err := json.MarshalIndent(msg, "", " ")
	if err != nil {
		return icinga.Unknown, err
	}

	return state, string(output)
}

func NewCmd() *cobra.Command {
	var opts options

	c := &cobra.Command{
		Use:   "check_node_status",
		Short: "Check Kubernetes Node",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host")

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

	c.Flags().StringVar(&opts.masterURL, "master", opts.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	c.Flags().StringVar(&opts.kubeconfigPath, "kubeconfig", opts.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	c.Flags().StringVarP(&opts.hostname, "host", "H", "", "Icinga host name")
	return c
}
