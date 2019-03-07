package check_pod_exists

import (
	"errors"
	"fmt"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"kmodules.xyz/client-go/tools/clientcmd"
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
	namespace  string
	selector   string
	podName    string
	count      int
	isCountSet bool
	// IcingaHost
	host *icinga.IcingaHost
}

func (o *options) complete(cmd *cobra.Command) (err error) {
	hostname, err := cmd.Flags().GetString(plugins.FlagHost)
	if err != nil {
		return err
	}
	o.host, err = icinga.ParseHost(hostname)
	if err != nil {
		return errors.New("invalid icinga host.name")
	}

	o.namespace = o.host.AlertNamespace
	o.isCountSet = cmd.Flag(flagCount).Changed

	o.kubeconfigPath, err = cmd.Flags().GetString(plugins.FlagKubeConfig)
	if err != nil {
		return
	}
	o.contextName, err = cmd.Flags().GetString(plugins.FlagKubeConfigContext)
	if err != nil {
		return
	}
	return nil
}

func (o *options) validate() error {
	if o.host.Type != icinga.TypeCluster {
		return errors.New("invalid icinga host type")
	}
	return nil
}

func (p *plugin) Check() (icinga.State, interface{}) {
	opts := p.options

	totalPod := 0
	if opts.podName != "" {
		_, err := p.client.Get(opts.podName, metav1.GetOptions{})
		if err != nil {
			return icinga.Unknown, err
		}
		totalPod = 1
	} else {
		podList, err := p.client.List(metav1.ListOptions{
			LabelSelector: opts.selector,
		})
		if err != nil {
			return icinga.Unknown, err
		}
		totalPod = len(podList.Items)
	}

	if opts.isCountSet {
		if opts.count != totalPod {
			return icinga.Critical, fmt.Sprintf("Found %d pod(s) instead of %d", totalPod, opts.count)
		} else {
			return icinga.OK, "Found all pods"
		}
	} else {
		if totalPod == 0 {
			return icinga.Critical, "No pod found"
		} else {
			return icinga.OK, fmt.Sprintf("Found %d pods(s)", totalPod)
		}
	}
}

const (
	flagCount = "count"
)

func NewCmd() *cobra.Command {
	var opts options
	cmd := &cobra.Command{
		Use:   "check_pod_exists",
		Short: "Check Kubernetes Pod(s)",

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
	cmd.Flags().StringP(plugins.FlagHost, "H", "", "Icinga host name")
	cmd.Flags().StringVarP(&opts.selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.")
	cmd.Flags().StringVarP(&opts.podName, "podName", "p", "", "Name of pod whose existence is checked")
	cmd.Flags().IntVarP(&opts.count, flagCount, "c", 0, "Number of Kubernetes pods")
	return cmd
}
