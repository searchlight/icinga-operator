package main

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/appscode/log"
	_ "github.com/appscode/searchlight/api/install"
	tcs "github.com/appscode/searchlight/client/clientset"
	_ "github.com/appscode/searchlight/client/clientset/fake"
	"github.com/appscode/searchlight/pkg/analytics"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL      string
	kubeconfigPath string

	icingaSecretName      string
	icingaSecretNamespace string = apiv1.NamespaceDefault

	address string = ":56790"

	kubeClient clientset.Interface
	extClient  tcs.ExtensionInterface

	enableAnalytics bool = true
)

func NewCmdRun() *cobra.Command {
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

			ctrl := controller.New(kubeClient, extClient)
			if err := ctrl.Setup(); err != nil {
				log.Fatalln(err)
			}
			//if icingaSecretName == "" {
			//	log.Fatalln("Missing icinga secret")
			//}
			//icingaClient, err := icinga.NewIcingaClient(w.KubeClient, icingaSecretName, icingaSecretNamespace)
			//if err != nil {
			//	log.Fatalln(err)
			//}
			//w.IcingaClient = icingaClient

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
	cmd.Flags().StringVarP(&icingaSecretName, "icinga-secret-name", "s", icingaSecretName, "Icinga secret name")
	cmd.Flags().StringVarP(&icingaSecretNamespace, "icinga-secret-namespace", "n", icingaSecretNamespace, "Icinga secret namespace")

	cmd.Flags().StringVar(&address, "address", address, "Address to listen on for web interface and telemetry.")

	// Analytics flags
	cmd.Flags().BoolVar(&enableAnalytics, "analytics", enableAnalytics, "Send analytical event to Google Analytics")

	return cmd
}
