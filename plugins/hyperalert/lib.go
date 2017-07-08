package hyperalert

import (
	"github.com/appscode/searchlight/plugins/check_component_status"
	"github.com/appscode/searchlight/plugins/check_influx_query"
	"github.com/appscode/searchlight/plugins/check_json_path"
	"github.com/appscode/searchlight/plugins/check_kube_event"
	"github.com/appscode/searchlight/plugins/check_kube_exec"
	"github.com/appscode/searchlight/plugins/check_node_count"
	"github.com/appscode/searchlight/plugins/check_node_status"
	"github.com/appscode/searchlight/plugins/check_pod_exists"
	"github.com/appscode/searchlight/plugins/check_pod_status"
	"github.com/appscode/searchlight/plugins/check_prometheus_metric"
	"github.com/appscode/searchlight/plugins/check_volume"
	"github.com/appscode/searchlight/plugins/notifier"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	c := &cobra.Command{
		Use:   "hyperalert",
		Short: "AppsCode Icinga2 plugin",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	c.AddCommand(check_component_status.NewCmd())
	c.AddCommand(check_influx_query.NewCmd())
	c.AddCommand(check_json_path.NewCmd())
	c.AddCommand(check_node_count.NewCmd())
	c.AddCommand(check_node_status.NewCmd())
	c.AddCommand(check_pod_exists.NewCmd())
	c.AddCommand(check_pod_status.NewCmd())
	c.AddCommand(check_prometheus_metric.NewCmd())
	c.AddCommand(check_volume.NewCmd())
	c.AddCommand(check_kube_event.NewCmd())
	c.AddCommand(check_kube_exec.NewCmd())
	c.AddCommand(notifier.NewCmd())
	return c
}
