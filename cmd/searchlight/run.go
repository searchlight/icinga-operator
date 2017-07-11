package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/appscode/log"
	_ "github.com/appscode/searchlight/api/install"
	tcs "github.com/appscode/searchlight/client/clientset"
	_ "github.com/appscode/searchlight/client/clientset/fake"
	"github.com/appscode/searchlight/pkg/analytics"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL       string
	kubeconfigPath  string
	address         string = ":56790"
	enableAnalytics bool   = true

	kubeClient clientset.Interface
	extClient  tcs.ExtensionInterface
)

func NewCmdRun() *cobra.Command {
	mgr := &icinga.Configurator{
		Expiry: 10 * 365 * 24 * time.Hour,
	}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run operator",
		PreRun: func(cmd *cobra.Command, args []string) {
			if enableAnalytics {
				analytics.Enable()
			}
			analytics.SendEvent("operator", "started", Version)
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			analytics.SendEvent("operator", "stopped", Version)
		},
		Run: func(cmd *cobra.Command, args []string) {
			config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
			if err != nil {
				log.Fatalf("Could not get Kubernetes config: %s", err)
			}

			kubeClient = clientset.NewForConfigOrDie(config)
			extClient = tcs.NewForConfigOrDie(config)

			err = mgr.GenerateCertificates()
			if err != nil {
				log.Fatalln(err)
			}
			cfg, err := mgr.LoadIcingaConfig()
			if err != nil {
				log.Fatalln(err)
			}
			icingaClient := icinga.NewClient(*cfg)
			for {
				icingaClient.Check()

				time.Sleep(1 * time.Second)
			}

			ctrl := controller.New(kubeClient, extClient, icingaClient)
			if err := ctrl.Setup(); err != nil {
				log.Fatalln(err)
			}

			log.Infoln("Starting Searchlight operator...")
			go ctrl.Run()

			if enableAnalytics {
				analytics.Enable()
			}
			analytics.SendEvent("operator", "started", Version)

			http.Handle("/metrics", promhttp.Handler())
			log.Infoln("Listening on", address)
			log.Fatal(http.ListenAndServe(address, nil))
		},
	}

	cmd.Flags().StringVar(&masterURL, "master", masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVarP(&mgr.ConfigRoot, "config-dir", "s", mgr.ConfigRoot, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")
	cmd.Flags().StringVar(&address, "address", address, "Address to listen on for web interface and telemetry.")
	cmd.Flags().BoolVar(&enableAnalytics, "analytics", enableAnalytics, "Send analytical event to Google Analytics")

	return cmd
}
