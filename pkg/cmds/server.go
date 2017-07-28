package cmds

import (
	"fmt"
	"os"

	"github.com/appscode/go/net"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/analytics"
	"github.com/appscode/searchlight/pkg/hostfacts"
	"github.com/spf13/cobra"
)

func NewCmdServer(version string) *cobra.Command {
	_, nodeIP, err := net.NodeIP()
	if err != nil {
		log.Fatalln(err)
	}
	log.Infoln("Detected IP for hostfacts server:", nodeIP.String())
	srv := hostfacts.Server{
		Address:  fmt.Sprintf("%s:%d", nodeIP, 56977),
		Username: os.Getenv("HOSTFACTS_AUTH_USERNAME"),
		Password: os.Getenv("HOSTFACTS_AUTH_PASSWORD"),
		Token:    os.Getenv("HOSTFACTS_AUTH_TOKEN"),
	}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run server",
		PreRun: func(cmd *cobra.Command, args []string) {
			//if opt.EnableAnalytics {
			//	analytics.Enable()
			//}
			analytics.SendEvent("operator", "started", version)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			analytics.SendEvent("operator", "stopped", version)
		},
		Run: func(cmd *cobra.Command, args []string) {
			srv.ListenAndServe()
		},
	}

	cmd.Flags().StringVar(&srv.Address, "address", srv.Address, "Http server address")
	cmd.Flags().StringVar(&srv.CACertFile, "caCertFile", srv.CACertFile, "File containing CA certificate")
	cmd.Flags().StringVar(&srv.CertFile, "certFile", srv.CertFile, "File container server TLS certificate")
	cmd.Flags().StringVar(&srv.KeyFile, "keyFile", srv.KeyFile, "File containing server TLS private key")

	cmd.Flags().StringVar(&srv.Username, "username", srv.Username, "Username used for basic authentication")
	cmd.Flags().StringVar(&srv.Password, "password", srv.Password, "Password used for basic authentication")
	cmd.Flags().StringVar(&srv.Token, "token", srv.Token, "Token used for bearer authentication")
	// cmd.Flags().BoolVar(&opt.EnableAnalytics, "analytics", opt.EnableAnalytics, "Send analytical event to Google Analytics")
	return cmd
}
