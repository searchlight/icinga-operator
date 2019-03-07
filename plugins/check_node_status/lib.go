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
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"kmodules.xyz/client-go/tools/clientcmd"
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
	client, err := clientcmd.ClientFromContext(opts.kubeconfigPath, opts.contextName)
	if err != nil {
		return nil, err
	}
	return newPlugin(client.CoreV1().Nodes(), opts), nil
}

type options struct {
	kubeconfigPath string
	contextName    string
	// options for Secret
	nodeName string
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
	o.nodeName = o.host.ObjectName

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
	if o.host.Type != icinga.TypeNode {
		return errors.New("invalid icinga host type")
	}
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
	} else if msg.Ready == core.ConditionUnknown {
		state = icinga.Unknown
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
