package cmds

import (
	"flag"
	"fmt"
	"log"

	v "github.com/appscode/go/version"
	"github.com/appscode/searchlight/client/clientset/versioned/scheme"
	"github.com/appscode/searchlight/pkg/hostfacts"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/client-go/tools/cli"
)

func NewCmdHostfacts() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hostfacts [command]",
		Short: `Hostfacts by AppsCode - Expose node metrics`,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			c.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
			})
			cli.SendAnalytics(c, v.Version.Version)

			scheme.AddToScheme(clientsetscheme.Scheme)
		},
	}
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	logs.ParseFlags()
	cmd.PersistentFlags().BoolVar(&cli.EnableAnalytics, "enable-analytics", cli.EnableAnalytics, "send usage events to Google Analytics")

	cmd.AddCommand(NewCmdServer())
	cmd.AddCommand(v.NewCmdVersion())
	return cmd
}

func NewCmdServer() *cobra.Command {
	srv := hostfacts.Server{
		Address: fmt.Sprintf(":%d", 56977),
	}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run server",
		PreRun: func(c *cobra.Command, args []string) {
			cli.SendPeriodicAnalytics(c, v.Version.Version)
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
	return cmd
}
