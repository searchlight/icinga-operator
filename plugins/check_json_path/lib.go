package check_json_path

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/appscode/envconfig"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/net/httpclient"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/jsonpath"
)

type plugin struct {
	client  corev1.SecretInterface
	options options
}

var _ plugins.PluginInterface = &plugin{}

func newPlugin(client corev1.SecretInterface, opts options) *plugin {
	return &plugin{client, opts}
}

func newPluginFromConfig(opts options) (*plugin, error) {
	config, err := clientcmd.BuildConfigFromFlags(opts.masterURL, opts.kubeconfigPath)
	if err != nil {
		return nil, err
	}
	client := kubernetes.NewForConfigOrDie(config).CoreV1().Secrets(opts.namespace)
	return newPlugin(client, opts), nil
}

type options struct {
	masterURL      string
	kubeconfigPath string
	// Icinga host name
	hostname string
	// http url
	url string
	// auth secret information
	secretName string
	namespace  string
	// Check condition
	warning  string
	critical string
}

func (o *options) validate() error {
	host, err := icinga.ParseHost(o.hostname)
	if err != nil {
		return errors.New("invalid icinga host.name")
	}
	if host.Type != icinga.TypeCluster {
		return errors.New("invalid icinga host type")
	}
	o.namespace = host.AlertNamespace
	return nil
}

type authInfo struct {
	username           string `envconfig:"USERNAME"`
	password           string `envconfig:"PASSWORD"`
	token              string `envconfig:"TOKEN"`
	caCertData         string `envconfig:"CA_CERT_DATA"`
	clientCertData     string `envconfig:"CLIENT_CERT_DATA"`
	clientKeyData      string `envconfig:"CLIENT_KEY_DATA"`
	insecureSkipVerify bool   `envconfig:"INSECURE_SKIP_VERIFY"`
}

func (p *plugin) getData() (interface{}, error) {
	opts := p.options

	var hc *httpclient.Client
	hc = httpclient.Default().WithBaseURL(opts.url).WithTimeout(time.Second * 10)

	if opts.secretName != "" {
		secret, err := p.client.Get(opts.secretName, metav1.GetOptions{})
		if err != nil {
			return "", err
		}
		var au authInfo
		err = envconfig.Load("", &au, func(key string) (string, bool) {
			v, ok := secret.Data[key]
			if !ok {
				return "", false
			}
			return string(v), true
		})
		if err != nil {
			return "", err
		}
		hc = hc.WithBasicAuth(au.username, au.password).WithBearerToken(au.token)
		if au.caCertData != "" {
			if au.clientCertData != "" && au.clientKeyData != "" {
				hc = hc.WithTLSConfig([]byte(au.caCertData), []byte(au.clientCertData), []byte(au.clientKeyData))
			} else {
				hc = hc.WithTLSConfig([]byte(au.caCertData))
			}
		}
		if au.insecureSkipVerify {
			hc = hc.WithInsecureSkipVerify()
		}
	}

	var respJson interface{}
	resp, err := hc.Call("GET", "", nil, &respJson, true)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("invalid status_code %v", resp.StatusCode)
	}

	return respJson, nil
}

func (p *plugin) checkResult(response interface{}, query string) (bool, error) {
	j := jsonpath.New("check")
	if err := j.Parse(query); err != nil {
		return false, err
	}
	buf := new(bytes.Buffer)
	if err := j.Execute(buf, response); err != nil {
		return false, err
	}

	expr, err := govaluate.NewEvaluableExpression(buf.String())
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	param := make(map[string]interface{})
	for _, v := range expr.Vars() {
		param[v] = v
	}
	res, err := expr.Evaluate(param)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	v, ok := res.(bool)
	return v && ok, nil
}

func (p *plugin) Check() (icinga.State, interface{}) {
	opts := p.options
	jsonInterface, err := p.getData()
	if err != nil {
		return icinga.Unknown, err
	}

	if opts.critical != "" {
		isCritical, err := p.checkResult(jsonInterface, opts.critical)
		if err != nil {
			return icinga.Unknown, err
		}
		if isCritical {
			return icinga.Critical, fmt.Sprintf("%v", opts.critical)
		}
	}
	if opts.warning != "" {
		isWarning, err := p.checkResult(jsonInterface, opts.warning)
		if err != nil {
			return icinga.Unknown, err
		}
		if isWarning {
			return icinga.Warning, fmt.Sprintf("%v", opts.warning)
		}
	}
	return icinga.OK, "response looks good"
}

func NewCmd() *cobra.Command {
	var opts options

	c := &cobra.Command{
		Use:   "check_json_path",
		Short: "Check Json Object",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host", "url")
			flags.EnsureAlterableFlags(cmd, "warning", "critical")

			if err := opts.validate(); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			plugin, err := newPluginFromConfig(opts)
			if err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			icinga.Output(plugin.Check())
		},
	}

	c.Flags().StringVar(&opts.masterURL, "master", opts.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	c.Flags().StringVar(&opts.kubeconfigPath, "kubeconfig", opts.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	c.Flags().StringVarP(&opts.hostname, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&opts.url, "url", "u", "", "URL to get data")
	c.Flags().StringVarP(&opts.secretName, "secretName", "s", "", `Kubernetes secret name`)
	c.Flags().StringVarP(&opts.warning, "warning", "w", "", `Warning jsonpath query which returns [true/false]`)
	c.Flags().StringVarP(&opts.critical, "critical", "c", "", `Critical jsonpath query which returns [true/false]`)
	return c
}
