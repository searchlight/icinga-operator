package cmds

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/appscode/go/analytics"
	v "github.com/appscode/go/version"
	"github.com/appscode/searchlight/client/clientset/versioned/scheme"
	"github.com/appscode/searchlight/pkg/hostfacts"
	"github.com/jpillora/go-ogle-analytics"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

func NewCmdHostfacts() *cobra.Command {
	var (
		enableAnalytics = true
	)
	cmd := &cobra.Command{
		Use:   "hostfacts [command]",
		Short: `Hostfacts by AppsCode - Expose node metrics`,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			c.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
			})
			if enableAnalytics && gaTrackingCode != "" {
				if client, err := ga.NewClient(gaTrackingCode); err == nil {
					client.ClientID(analytics.ClientID())
					parts := strings.Split(c.CommandPath(), " ")
					client.Send(ga.NewEvent(parts[0], strings.Join(parts[1:], "/")).Label(v.Version.Version))
				}
			}
			scheme.AddToScheme(clientsetscheme.Scheme)
		},
	}
	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})
	cmd.PersistentFlags().BoolVar(&enableAnalytics, "enable-analytics", enableAnalytics, "send usage events to Google Analytics")

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
