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

package hyperalert

import (
	"flag"

	"go.searchlight.dev/icinga-operator/client/clientset/versioned/scheme"
	"go.searchlight.dev/icinga-operator/plugins"
	"go.searchlight.dev/icinga-operator/plugins/analytics_id"
	"go.searchlight.dev/icinga-operator/plugins/check_ca_cert"
	"go.searchlight.dev/icinga-operator/plugins/check_cert"
	"go.searchlight.dev/icinga-operator/plugins/check_component_status"
	"go.searchlight.dev/icinga-operator/plugins/check_env"
	"go.searchlight.dev/icinga-operator/plugins/check_event"
	"go.searchlight.dev/icinga-operator/plugins/check_json_path"
	"go.searchlight.dev/icinga-operator/plugins/check_node_exists"
	"go.searchlight.dev/icinga-operator/plugins/check_node_status"
	"go.searchlight.dev/icinga-operator/plugins/check_pod_exec"
	"go.searchlight.dev/icinga-operator/plugins/check_pod_exists"
	"go.searchlight.dev/icinga-operator/plugins/check_pod_status"
	"go.searchlight.dev/icinga-operator/plugins/check_volume"
	"go.searchlight.dev/icinga-operator/plugins/check_webhook"
	"go.searchlight.dev/icinga-operator/plugins/notifier"

	"github.com/spf13/cobra"
	"gomodules.xyz/kglog"
	v "gomodules.xyz/x/version"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hyperalert",
		Short: "AppsCode Icinga2 plugin",
		PersistentPreRun: func(c *cobra.Command, args []string) {
			scheme.AddToScheme(clientsetscheme.Scheme)
		},
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	kglog.ParseFlags()
	cmd.PersistentFlags().String(plugins.FlagKubeConfig, "", "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.PersistentFlags().String(plugins.FlagKubeConfigContext, "", "Use the context in kubeconfig")
	cmd.PersistentFlags().Int(plugins.FlagCheckInterval, 30, "Icinga check_interval in second. [Format: 30, 300]")

	// CheckCluster
	cmd.AddCommand(check_component_status.NewCmd())
	cmd.AddCommand(check_json_path.NewCmd())
	cmd.AddCommand(check_node_exists.NewCmd())
	cmd.AddCommand(check_pod_exists.NewCmd())
	cmd.AddCommand(check_event.NewCmd())
	cmd.AddCommand(check_ca_cert.NewCmd())
	cmd.AddCommand(check_cert.NewCmd())
	cmd.AddCommand(check_env.NewCmd())
	cmd.AddCommand(check_webhook.NewCmd())

	// CheckNode
	cmd.AddCommand(check_node_status.NewCmd())

	// CheckPod
	cmd.AddCommand(check_pod_status.NewCmd())
	cmd.AddCommand(check_pod_exec.NewCmd())

	// Combined
	cmd.AddCommand(check_volume.NewCmd())

	// Notifier
	cmd.AddCommand(notifier.NewCmd())

	cmd.AddCommand(analytics_id.NewCmd())
	cmd.AddCommand(v.NewCmdVersion())

	return cmd
}
