package cmds

import (
	_ "net/http/pprof"
	"time"

	"github.com/appscode/log"
	tcs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/migrator"
	"github.com/appscode/searchlight/pkg/operator"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewCmdOperator() *cobra.Command {
	opt := operator.Options{
		ConfigRoot:       "/srv",
		ConfigSecretName: "searchlight-operator",
		APIAddress:       ":8080",
		WebAddress:       ":56790",
		ResyncPeriod:     5 * time.Minute,
	}
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run operator",
		Run: func(cmd *cobra.Command, args []string) {
			run(opt)
		},
	}

	cmd.Flags().StringVar(&opt.Master, "master", opt.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&opt.KubeConfig, "kubeconfig", opt.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVar(&opt.ConfigRoot, "config-dir", opt.ConfigRoot, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")
	cmd.Flags().StringVar(&opt.ConfigSecretName, "config-secret-name", opt.ConfigSecretName, "Name of Kubernetes secret used to pass icinga credentials.")
	cmd.Flags().StringVar(&opt.APIAddress, "api.address", opt.APIAddress, "The address of the Searchlight API Server")
	cmd.Flags().StringVar(&opt.WebAddress, "web.address", opt.WebAddress, "Address to listen on for web interface and telemetry.")
	cmd.Flags().DurationVar(&opt.ResyncPeriod, "resync-period", opt.ResyncPeriod, "If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out.")

	return cmd
}

func run(opt operator.Options) {
	config, err := clientcmd.BuildConfigFromFlags(opt.Master, opt.KubeConfig)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	kubeClient := clientset.NewForConfigOrDie(config)
	apiExtKubeClient := apiextensionsclient.NewForConfigOrDie(config)
	extClient := tcs.NewForConfigOrDie(config)

	secret, err := kubeClient.CoreV1().Secrets(util.OperatorNamespace()).Get(opt.ConfigSecretName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to load secret: %s", err)
	}

	mgr := &icinga.Configurator{
		ConfigRoot:       opt.ConfigRoot,
		IcingaSecretName: opt.ConfigSecretName,
		Expiry:           10 * 365 * 24 * time.Hour,
	}
	cfg, err := mgr.LoadConfig(func(key string) (value string, found bool) {
		var bytes []byte
		bytes, found = secret.Data[key]
		value = string(bytes)
		return
	})
	if err != nil {
		log.Fatalln(err)
	}
	icingaClient := icinga.NewClient(*cfg)
	for {
		if icingaClient.Check().Get(nil).Do().Status == 200 {
			log.Infoln("connected to icinga api")
			break
		}
		log.Infoln("Waiting for icinga to start")
		time.Sleep(2 * time.Second)
	}

	op := operator.New(kubeClient, apiExtKubeClient, extClient, icingaClient, opt)
	if err := op.Setup(); err != nil {
		log.Fatalln(err)
	}

	if err = migrator.NewMigrator(kubeClient, apiExtKubeClient, extClient).RunMigration(); err != nil {
		log.Fatalln(err)
	}

	log.Infoln("Starting Searchlight operator...")
	op.RunAndHold()
}
