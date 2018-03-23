package server

import (
	"flag"
	"time"

	"github.com/appscode/go/log"
	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	"github.com/appscode/kutil/meta"
	cs "github.com/appscode/searchlight/client/clientset/versioned"
	"github.com/appscode/searchlight/pkg/admission/plugin"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/operator"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type OperatorOptions struct {
	ConfigRoot       string
	ConfigSecretName string
	OpsAddress       string
	ResyncPeriod     time.Duration
	MaxNumRequeues   int
	NumThreads       int
	// V logging level, the value of the -v flag
	verbosity string
}

func NewOperatorOptions() *OperatorOptions {
	return &OperatorOptions{
		ConfigRoot:       "/srv",
		ConfigSecretName: "searchlight-operator",
		OpsAddress:       ":56790",
		ResyncPeriod:     5 * time.Minute,
		MaxNumRequeues:   5,
		NumThreads:       1,
		verbosity:        "3",
	}
}

func (s *OperatorOptions) AddGoFlags(fs *flag.FlagSet) {
	fs.StringVar(&s.ConfigRoot, "config-dir", s.ConfigRoot, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")
	fs.StringVar(&s.ConfigSecretName, "config-secret-name", s.ConfigSecretName, "Name of Kubernetes secret used to pass icinga credentials.")
	fs.StringVar(&s.OpsAddress, "ops-address", s.OpsAddress, "Address to listen on for web interface and telemetry.")
	fs.DurationVar(&s.ResyncPeriod, "resync-period", s.ResyncPeriod, "If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out.")
}

func (s *OperatorOptions) AddFlags(fs *pflag.FlagSet) {
	pfs := flag.NewFlagSet("searchlight", flag.ExitOnError)
	s.AddGoFlags(pfs)
	fs.AddGoFlagSet(pfs)
}

func (s *OperatorOptions) ApplyTo(cfg *operator.OperatorConfig) error {
	var err error

	cfg.ConfigRoot = s.ConfigRoot
	cfg.ConfigSecretName = s.ConfigSecretName
	cfg.OpsAddress = s.OpsAddress
	cfg.ResyncPeriod = s.ResyncPeriod
	cfg.MaxNumRequeues = s.MaxNumRequeues
	cfg.NumThreads = s.NumThreads
	cfg.Verbosity = s.verbosity

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.ExtClient, err = cs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.CRDClient, err = crd_cs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	cfg.AdmissionHooks = []hooks.AdmissionHook{&plugin.CRDValidator{}}

	secret, err := cfg.KubeClient.CoreV1().Secrets(meta.Namespace()).Get(s.ConfigSecretName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to load secret: %s", s.ConfigSecretName)
	}

	mgr := &icinga.Configurator{
		ConfigRoot:       s.ConfigRoot,
		IcingaSecretName: s.ConfigSecretName,
		Expiry:           10 * 365 * 24 * time.Hour,
	}
	data, err := mgr.LoadConfig(func(key string) (value string, found bool) {
		var bytes []byte
		bytes, found = secret.Data[key]
		value = string(bytes)
		return
	})
	if err != nil {
		log.Fatalln(err)
	}

	cfg.IcingaClient = icinga.NewClient(*data)
	for {
		if cfg.IcingaClient.Check().Get(nil).Do().Status == 200 {
			log.Infoln("connected to icinga api")
			break
		}
		log.Infoln("Waiting for icinga to start")
		time.Sleep(2 * time.Second)
	}

	return nil
}
