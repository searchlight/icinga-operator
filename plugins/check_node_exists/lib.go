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

package check_node_exists

import (
	"context"
	"fmt"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	"github.com/spf13/cobra"
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
	// http url
	selector   string
	nodeName   string
	count      int
	isCountSet bool
}

func (o *options) complete(cmd *cobra.Command) (err error) {
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
	return nil
}

func (p *plugin) Check() (icinga.State, interface{}) {
	opts := p.options

	totalNode := 0
	if opts.nodeName != "" {
		node, err := p.client.Get(context.TODO(), opts.nodeName, metav1.GetOptions{})
		if err != nil {
			return icinga.Unknown, err
		}
		if node != nil {
			totalNode = 1
		}
	} else {
		nodeList, err := p.client.List(context.TODO(),
			metav1.ListOptions{
				LabelSelector: opts.selector,
			},
		)
		if err != nil {
			return icinga.Unknown, err
		}

		totalNode = len(nodeList.Items)
	}

	if opts.isCountSet {
		if opts.count != totalNode {
			return icinga.Critical, fmt.Sprintf("Found %d node(s) instead of %d", totalNode, opts.count)
		} else {
			return icinga.OK, "Found all nodes"
		}
	} else {
		if totalNode == 0 {
			return icinga.Critical, "No node found"
		} else {
			return icinga.OK, fmt.Sprintf("Found %d node(s)", totalNode)
		}
	}
}

const (
	flagCount = "count"
)

func NewCmd() *cobra.Command {
	var opts options

	cmd := &cobra.Command{
		Use:   "check_node_exists",
		Short: "Count Kubernetes Nodes",

		Run: func(cmd *cobra.Command, args []string) {
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

	cmd.Flags().StringVarP(&opts.selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.")
	cmd.Flags().StringVarP(&opts.nodeName, "nodeName", "n", "", "Name of node whose existence is checked")
	cmd.Flags().IntVarP(&opts.count, flagCount, "c", 0, "Number of expected Kubernetes Node")
	return cmd
}
