package hyperalert

import (
	"github.com/appscode/searchlight/plugins/check_ca_cert"
	"github.com/appscode/searchlight/plugins/check_component_status"
	"github.com/appscode/searchlight/plugins/check_env"
	"github.com/appscode/searchlight/plugins/check_event"
	"github.com/appscode/searchlight/plugins/check_influx_query"
	"github.com/appscode/searchlight/plugins/check_json_path"
	"github.com/appscode/searchlight/plugins/check_node_exists"
	"github.com/appscode/searchlight/plugins/check_node_status"
	"github.com/appscode/searchlight/plugins/check_pod_exec"
	"github.com/appscode/searchlight/plugins/check_pod_exists"
	"github.com/appscode/searchlight/plugins/check_pod_status"
	"github.com/appscode/searchlight/plugins/check_volume"
	"github.com/appscode/searchlight/plugins/notifier"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hyperalert",
		Short: "AppsCode Icinga2 plugin",
		Run: func(c *cobra.Command, args []string) {
			c.Help()
		},
	}

	// CheckCluster
	cmd.AddCommand(check_component_status.NewCmd())
	cmd.AddCommand(check_json_path.NewCmd())
	cmd.AddCommand(check_node_exists.NewCmd())
	cmd.AddCommand(check_pod_exists.NewCmd())
	cmd.AddCommand(check_event.NewCmd())
	cmd.AddCommand(check_ca_cert.NewCmd())
	cmd.AddCommand(check_env.NewCmd())

	// CheckNode
	cmd.AddCommand(check_node_status.NewCmd())

	// CheckPod
	cmd.AddCommand(check_pod_status.NewCmd())
	cmd.AddCommand(check_pod_exec.NewCmd())

	// Combined
	cmd.AddCommand(check_volume.NewCmd())
	cmd.AddCommand(check_influx_query.NewCmd())

	// Notifier
	cmd.AddCommand(notifier.NewCmd())

	return cmd
}
