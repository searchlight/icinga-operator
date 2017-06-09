package main

import (
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/appscode/log"
	_ "github.com/appscode/searchlight/api/install"
	acs "github.com/appscode/searchlight/client/clientset"
	_ "github.com/appscode/searchlight/client/clientset/fake"
	"github.com/appscode/searchlight/pkg/client/icinga"
	acw "github.com/appscode/searchlight/pkg/watcher"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	_ "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset/fake"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

var (
	masterURL      string
	kubeconfigPath string

	icingaSecretName      string
	icingaSecretNamespace string = kapi.NamespaceDefault

	address string

	kubeClient clientset.Interface
	extClient  acs.ExtensionInterface
)

func NewCmdRun() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run operator",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	cmd.Flags().StringVar(&masterURL, "master", masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVarP(&icingaSecretName, "icinga-secret-name", "s", icingaSecretName, "Icinga secret name")
	cmd.Flags().StringVarP(&icingaSecretNamespace, "icinga-secret-namespace", "n", icingaSecretNamespace, "Icinga secret namespace")

	cmd.Flags().StringVar(&address, "address", address, "Address to listen on for web interface and telemetry.")

	return cmd
}

func run() {
	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	kubeClient = clientset.NewForConfigOrDie(config)
	extClient = acs.NewForConfigOrDie(config)

	w := &acw.Watcher{
		KubeClient: kubeClient,
		ExtClient:  extClient,
		SyncPeriod: time.Minute * 2,
	}
	if icingaSecretName == "" {
		log.Fatalln("Missing icinga secret")
	}
	icingaClient, err := icinga.NewIcingaClient(w.KubeClient, icingaSecretName, icingaSecretNamespace)
	if err != nil {
		log.Fatalln(err)
	}
	w.IcingaClient = icingaClient

	log.Infoln("Starting Searchlight operator...")
	go w.Run()

	http.Handle("/metrics", promhttp.Handler())
	log.Infoln("Listening on", address)
	log.Fatal(http.ListenAndServe(address, nil))
}
