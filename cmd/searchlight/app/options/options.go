package options

import (
	"github.com/spf13/pflag"
	kapi "k8s.io/kubernetes/pkg/api"
)

type Config struct {
	Master                string
	KubeConfig            string
	IcingaSecretName      string
	IcingaSecretNamespace string
}

func NewConfig() *Config {
	return &Config{
		Master:                "",
		KubeConfig:            "",
		IcingaSecretName:      "",
		IcingaSecretNamespace: kapi.NamespaceDefault,
	}
}

func (s *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&s.Master, "master", s.Master, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	fs.StringVar(&s.KubeConfig, "kubeconfig", s.KubeConfig, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	fs.StringVarP(&s.IcingaSecretName, "icinga-secret-name", "s", s.IcingaSecretName, "Icinga secret name")
	fs.StringVarP(&s.IcingaSecretNamespace, "icinga-secret-namespace", "n", s.IcingaSecretNamespace, "Icinga secret namespace")
}
