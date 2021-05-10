/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package check_pod_status

import (
	"context"
	"errors"
	"fmt"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	"github.com/spf13/cobra"
	"gomodules.xyz/x/flags"
	core "k8s.io/api/core/v1"
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

func (p *plugin) Check() (icinga.State, interface{}) {
	opts := p.options

	pod, err := p.client.Get(context.TODO(), opts.podName, metav1.GetOptions{})
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
func (p *plugin) podRunningAndReady(pod core.Pod) (bool, interface{}) {
	switch pod.Status.Phase {
	case core.PodFailed, core.PodSucceeded:
		return false, fmt.Errorf("pod completed")
	case core.PodRunning:
		for _, cond := range pod.Status.Conditions {
			if cond.Type != core.PodReady {
				continue
			}
			if cond.Status == core.ConditionTrue {
				return true, nil
			}

			return false, cond.Message
		}
		return false, fmt.Errorf("pod ready condition not found")
	default:
		return false, fmt.Errorf(`pod is in "%s" phase`, pod.Status.Phase)
	}
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
