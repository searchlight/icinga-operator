package cmds

import (
	"flag"
	"io"
	"log"
	"os"

	v "github.com/appscode/go/version"
	"github.com/appscode/kutil/tools/cli"
	"github.com/appscode/searchlight/client/clientset/versioned/scheme"
	"github.com/appscode/searchlight/pkg/cmds/server"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	genericapiserver "k8s.io/apiserver/pkg/server"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
)

func NewCmdSearchlight() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "searchlight [command]",
		Short: `Searchlight by AppsCode - Alerts for Kubernetes`,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			c.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
			})
			cli.SendAnalytics(c, v.Version.Version)

			scheme.AddToScheme(clientsetscheme.Scheme)
		},
	}
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	// ref: https://github.com/kubernetes/kubernetes/issues/17162#issuecomment-225596212
	flag.CommandLine.Parse([]string{})
	rootCmd.PersistentFlags().BoolVar(&cli.EnableAnalytics, "enable-analytics", cli.EnableAnalytics, "send usage events to Google Analytics")

	stopCh := genericapiserver.SetupSignalHandler()
	rootCmd.AddCommand(NewCmdRun(os.Stdout, os.Stderr, stopCh))
	rootCmd.AddCommand(NewCmdConfigure())
	rootCmd.AddCommand(v.NewCmdVersion())

	return rootCmd
}

func NewCmdRun(out, errOut io.Writer, stopCh <-chan struct{}) *cobra.Command {
	o := server.NewSearchlightOptions(out, errOut)

	cmd := &cobra.Command{
		Use:               "run",
		Short:             "Launch Searchlight operator",
		Long:              "Launch Searchlight operator",
		DisableAutoGenTag: true,
		PreRun: func(c *cobra.Command, args []string) {
			cli.SendPeriodicAnalytics(c, v.Version.Version)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			glog.Infof("Starting operator version %s+%s ...", v.Version.Version, v.Version.CommitHash)

			if err := o.Complete(cmd); err != nil {
				return err
			}
			if err := o.Validate(args); err != nil {
				return err
			}
			if err := o.Run(stopCh); err != nil {
				return err
			}
			return nil
		},
	}

	o.AddFlags(cmd.Flags())

	return cmd
}
