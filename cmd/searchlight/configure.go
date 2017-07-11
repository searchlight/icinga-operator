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
		Short: "Generate certificates for Kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.SetLogLevel(4)
			flags.EnsureRequiredFlags(cmd, "folder", "namespace", "cluster", "master-external-ip", "master-internal-ip")

			return mgr.GenerateCertificates()
		},
	}

	cmd.Flags().StringVarP(&configDir, "config-dir", "s", configDir, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")
	cmd.Flags().StringVar(&mgr.ConfigRoot, "folder", mgr.ConfigRoot, "Folder where certs are stored")

	return cmd
}
