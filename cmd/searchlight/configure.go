package main

import (
	"time"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
)

func NewCmdGenerate() *cobra.Command {
	mgr := &icinga.Configurator{
		Expiry: 10 * 365 * 24 * time.Hour,
	}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate icinga config",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.SetLogLevel(4)
			flags.EnsureRequiredFlags(cmd, "folder", "namespace", "cluster", "master-external-ip", "master-internal-ip")

			err := mgr.GenerateCertificates()
			if err != nil {
				return err
			}
			_, err = mgr.LoadIcingaConfig()
			return err
		},
	}

	cmd.Flags().StringVarP(&mgr.ConfigRoot, "config-dir", "s", mgr.ConfigRoot, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")

	return cmd
}
