package main

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/log"
	_ "github.com/appscode/searchlight/api/install"
	tcs "github.com/appscode/searchlight/client/clientset"
	_ "github.com/appscode/searchlight/client/clientset/fake"
	"github.com/appscode/searchlight/pkg/analytics"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"gopkg.in/ini.v1"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	masterURL      string
	kubeconfigPath string
	configDir      string

	icingaSecretName      string
	icingaSecretNamespace string = apiv1.NamespaceDefault

	address string = ":56790"

	kubeClient clientset.Interface
	extClient  tcs.ExtensionInterface

	enableAnalytics bool = true
)

const (
	ICINGA_WEB_HOST           = "ICINGA_WEB_HOST"
	ICINGA_WEB_PORT           = "ICINGA_WEB_PORT"
	ICINGA_WEB_DB             = "ICINGA_WEB_DB"
	ICINGA_WEB_USER           = "ICINGA_WEB_USER"
	ICINGA_WEB_PASSWORD       = "ICINGA_WEB_PASSWORD"
	ICINGA_WEB_ADMIN_PASSWORD = "ICINGA_WEB_ADMIN_PASSWORD"
	ICINGA_IDO_HOST           = "ICINGA_IDO_HOST"
	ICINGA_IDO_PORT           = "ICINGA_IDO_PORT"
	ICINGA_IDO_DB             = "ICINGA_IDO_DB"
	ICINGA_IDO_USER           = "ICINGA_IDO_USER"
	ICINGA_IDO_PASSWORD       = "ICINGA_IDO_PASSWORD"
	ICINGA_API_USER           = "ICINGA_API_USER"
	ICINGA_API_PASSWORD       = "ICINGA_API_PASSWORD"
	ICINGA_ADDRESS            = "ICINGA_ADDRESS"

	ICINGA_CA_CERT     = "ICINGA_CA_CERT"
	ICINGA_SERVER_KEY  = "ICINGA_SERVER_KEY"
	ICINGA_SERVER_CERT = "ICINGA_SERVER_CERT"
)

//ca.crt: $ICINGA_CA_CERT
//icinga.key: $ICINGA_SERVER_KEY
//icinga.crt: $ICINGA_SERVER_CERT

var (
	icingaKeys = []string{
		ICINGA_WEB_HOST,
		ICINGA_WEB_PORT,
		ICINGA_WEB_DB,
		ICINGA_WEB_USER,
		ICINGA_WEB_PASSWORD,
		ICINGA_WEB_ADMIN_PASSWORD,
		ICINGA_IDO_HOST,
		ICINGA_IDO_PORT,
		ICINGA_IDO_DB,
		ICINGA_IDO_USER,
		ICINGA_IDO_PASSWORD,
		ICINGA_API_USER,
		ICINGA_API_PASSWORD,
		ICINGA_ADDRESS,
		ICINGA_CA_CERT,
		ICINGA_SERVER_KEY,
		ICINGA_SERVER_CERT,
	}
)

//ICINGA_WEB_HOST=127.0.0.1
//ICINGA_WEB_PORT=5432
//ICINGA_WEB_DB=icingawebdb
//ICINGA_WEB_USER=icingaweb
//ICINGA_WEB_PASSWORD=12345678
//ICINGA_WEB_ADMIN_PASSWORD=admin
//ICINGA_IDO_HOST=127.0.0.1
//ICINGA_IDO_PORT=5432
//ICINGA_IDO_DB=icingaidodb
//ICINGA_IDO_USER=icingaido
//ICINGA_IDO_PASSWORD=12345678
//ICINGA_API_USER=icingaapi
//ICINGA_API_PASSWORD=12345678
//ICINGA_ADDRESS=searchlight-icinga.kube-system

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
			if icingaSecretName == "" {
				log.Fatalln("Missing icinga secret")
			}

			pkiDir := filepath.Join(configDir, "pki")
			err = os.MkdirAll(pkiDir, 0755)

			// gen pki

			cfgFile := filepath.Join(configDir, ".icinga/config")

			if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
				// auto generate the file
				cfg := ini.Empty()
				sec := cfg.Section("")
				sec.NewKey(ICINGA_WEB_HOST, "127.0.0.1")
				sec.NewKey(ICINGA_WEB_PORT, "5432")
				sec.NewKey(ICINGA_WEB_DB, "icingawebdb")
				sec.NewKey(ICINGA_WEB_USER, "icingaweb")
				sec.NewKey(ICINGA_WEB_PASSWORD, rand.GeneratePassword())
				sec.NewKey(ICINGA_WEB_ADMIN_PASSWORD, rand.GeneratePassword())
				sec.NewKey(ICINGA_IDO_HOST, "127.0.0.1")
				sec.NewKey(ICINGA_IDO_PORT, "5432")
				sec.NewKey(ICINGA_IDO_DB, "icingaidodb")
				sec.NewKey(ICINGA_IDO_USER, "icingaido")
				sec.NewKey(ICINGA_IDO_PASSWORD, rand.GeneratePassword())
				sec.NewKey(ICINGA_API_USER, "icingaapi")
				sec.NewKey(ICINGA_API_PASSWORD, rand.GeneratePassword())
				sec.NewKey(ICINGA_ADDRESS, "127.0.0.1")

				err = os.MkdirAll(filepath.Dir(cfgFile), 0755)
				if err != nil {
					log.Errorln(err)
				}
				err = cfg.SaveTo(configDir)
				if err != nil {
					log.Errorln(err)
				}
			}

			cfg, err := ini.Load(cfgFile)
			if err != nil {
				log.Errorln(err)
			}
			sec := cfg.Section("")
			for _, key := range icingaKeys {
				if !sec.HasKey(key) {
					log.Fatalf("No Icinga config found for key %s", key)
				}
			}

			// path/to/whatever does not exist

			icingaClient, err := icinga.NewClient(kubeClient, icingaSecretName, icingaSecretNamespace)
			if err != nil {
				log.Fatalln(err)
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
	cmd.Flags().StringVarP(&configDir, "config-dir", "s", configDir, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")
	cmd.Flags().StringVar(&address, "address", address, "Address to listen on for web interface and telemetry.")
	cmd.Flags().BoolVar(&enableAnalytics, "analytics", enableAnalytics, "Send analytical event to Google Analytics")

	return cmd
}
