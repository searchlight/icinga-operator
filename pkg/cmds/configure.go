package cmds

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
)

func NewCmdConfigure() *cobra.Command {
	mgr := &icinga.Configurator{
		Expiry: 10 * 365 * 24 * time.Hour,
	}
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Generate icinga configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.SetLogLevel(4)

			cfg, err := mgr.LoadConfig(func(key string) (string, bool) {
				return "", false
			})
			if err != nil {
				return err
			}
			bytes, err := json.MarshalIndent(cfg, "", " ")
			if err != nil {
				return err
			}
			fmt.Println(string(bytes))
			return err
		},
	}

	cmd.Flags().StringVarP(&mgr.ConfigRoot, "config-dir", "s", mgr.ConfigRoot, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")

	return cmd
}
