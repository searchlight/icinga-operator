package check_cert

import (
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/cert"
)

type plugin struct {
	client corev1.SecretInterface
	option *Option
}

var _ plugins.PluginInterface = &plugin{}

func NewPlugin(client corev1.SecretInterface, option *Option) *plugin {
	return &plugin{client, option}
}

func newPluginFromConfig(option *Option) (*plugin, error) {
	config, err := clientcmd.BuildConfigFromFlags(option.masterURL, option.kubeconfigPath)
	if err != nil {
		return nil, err
	}
	client := kubernetes.NewForConfigOrDie(config).CoreV1().Secrets(option.Namespace)
	return NewPlugin(client, option), nil
}

type Option struct {
	masterURL      string
	kubeconfigPath string
	hostname       string

	Namespace  string
	Selector   string
	SecretName string
	SecretKey  []string

	Warning  time.Duration
	Critical time.Duration
}

func (o *Option) validate() {
	host, err := icinga.ParseHost(o.hostname)
	if err != nil {
		fmt.Fprintln(os.Stdout, icinga.Warning, "Invalid icinga host.name")
		os.Exit(3)
	}
	if host.Type != icinga.TypeCluster {
		fmt.Fprintln(os.Stdout, icinga.Warning, "Invalid icinga host type")
		os.Exit(3)
	}
	o.Namespace = host.AlertNamespace
}

func (p *plugin) getCertSecrets() ([]core.Secret, error) {
	option := p.option
	if option.SecretName != "" {
		var secret *core.Secret
		secret, err := p.client.Get(option.SecretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return []core.Secret{*secret}, nil
	} else {
		secretList, err := p.client.List(metav1.ListOptions{
			LabelSelector: option.Selector,
		})
		if err != nil {
			return nil, err
		}
		return secretList.Items, nil
	}
	return nil, nil
}

func (p *plugin) checkNotAfter(cert *x509.Certificate) (icinga.State, time.Duration, bool) {
	option := p.option
	remaining := cert.NotAfter.Sub(time.Now())
	if remaining.Seconds() < option.Critical.Seconds() {
		return icinga.Critical, remaining, false
	}

	if remaining.Seconds() < option.Warning.Seconds() {
		return icinga.Warning, remaining, false
	}

	return icinga.OK, remaining, true
}

func (p *plugin) checkCert(data []byte, secret *core.Secret, key string) (icinga.State, error) {
	certs, err := cert.ParseCertsPEM(data)
	if err != nil {
		return icinga.Unknown, fmt.Errorf(
			`failed to parse certificate for key "%s" in Secret "%s/%s"`,
			key, secret.Namespace, secret.Name,
		)
	}

	for _, cert := range certs {
		if state, remaining, ok := p.checkNotAfter(cert); !ok {
			return state, fmt.Errorf(
				`certificate for key "%s" in Secret "%s/%s" will be expired within %v hours`,
				key, secret.Namespace, secret.Name, remaining.Hours(),
			)
		}
	}
	return icinga.OK, nil
}

func (p *plugin) checkCertPerSecretKey(secret *core.Secret) (icinga.State, error) {
	option := p.option
	for _, key := range option.SecretKey {
		data, ok := secret.Data[key]
		if !ok {
			return icinga.Warning, fmt.Errorf(`key "%s" not found in Secret "%s/%s"`, key, secret.Namespace, secret.Name)
		}

		if state, err := p.checkCert(data, secret, key); err != nil {
			return state, err
		}
	}

	if len(option.SecretKey) == 0 && secret.Type == core.SecretTypeTLS {
		data, ok := secret.Data["tls.crt"]
		if !ok {
			return icinga.Warning, fmt.Errorf(`key "tls.crt" not found in Secret "%s/%s"`, secret.Namespace, secret.Name)
		}

		if state, err := p.checkCert(data, secret, "tls.crt"); err != nil {
			return state, err
		}
	}

	return icinga.OK, nil
}

func (p *plugin) Check() (icinga.State, interface{}) {

	secretList, err := p.getCertSecrets()
	if err != nil {
		return icinga.Unknown, err
	}

	for _, secret := range secretList {
		if state, err := p.checkCertPerSecretKey(&secret); err != nil {
			return state, err
		}
	}

	return icinga.OK, fmt.Sprintf("Certificate expirity check is succeeded")
}

func NewCmd() *cobra.Command {
	var option *Option

	cmd := &cobra.Command{
		Use:   "check_cert",
		Short: "Check Certificate expire date",

		Run: func(cmd *cobra.Command, args []string) {
			option.validate()
			plugin, err := newPluginFromConfig(option)
			if err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			icinga.Output(plugin.Check())
		},
	}

	cmd.Flags().StringVar(&option.masterURL, "master", option.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&option.kubeconfigPath, "kubeconfig", option.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVarP(&option.hostname, "host", "H", "", "Icinga host name")
	cmd.Flags().StringVarP(&option.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='")
	cmd.Flags().StringVarP(&option.SecretName, "secretName", "s", "", "Name of secret from where certificates are checked")
	cmd.Flags().StringSliceVarP(&option.SecretKey, "secretKey", "k", nil, "Name of secret key where certificates are kept")
	cmd.Flags().DurationVarP(&option.Warning, "warning", "w", time.Hour*360, `Remaining duration for Warning state. [Default: 360h]`)
	cmd.Flags().DurationVarP(&option.Critical, "critical", "c", time.Hour*120, `Remaining duration for Critical state. [Default: 120h]`)
	return cmd
}
