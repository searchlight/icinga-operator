/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package check_json_path

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	"github.com/Knetic/govaluate"
	"github.com/spf13/cobra"
	"gomodules.xyz/envconfig"
	"gomodules.xyz/x/flags"
	"gomodules.xyz/x/net/httpclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/util/jsonpath"
	"kmodules.xyz/client-go/tools/clientcmd"
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
	client, err := clientcmd.ClientFromContext(opts.kubeconfigPath, opts.contextName)
	if err != nil {
		return nil, err
	}
	return newPlugin(client.CoreV1().Secrets(opts.namespace), opts), nil
}

type options struct {
	kubeconfigPath string
	contextName    string
	// http url
	url string
	// auth secret information
	secretName string
	namespace  string
	// Check condition
	warning  string
	critical string
	// IcingaHost
	host *icinga.IcingaHost
}

func (o *options) complete(cmd *cobra.Command) (err error) {
	hostname, err := cmd.Flags().GetString(plugins.FlagHost)
	if err != nil {
		return err
	}
	o.host, err = icinga.ParseHost(hostname)
	if err != nil {
		return errors.New("invalid icinga host.name")
	}
	o.namespace = o.host.AlertNamespace

	o.kubeconfigPath, err = cmd.Flags().GetString(plugins.FlagKubeConfig)
	if err != nil {
		return
	}
	o.contextName, err = cmd.Flags().GetString(plugins.FlagKubeConfigContext)
	if err != nil {
		return
	}
	return nil
}

func (o *options) validate() error {
	if o.host.Type != icinga.TypeCluster {
		return errors.New("invalid icinga host type")
	}
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
		secret, err := p.client.Get(context.TODO(), opts.secretName, metav1.GetOptions{})
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

const (
	flagURL = "url"
)

func NewCmd() *cobra.Command {
	var opts options

	c := &cobra.Command{
		Use:   "check_json_path",
		Short: "Check Json Object",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, plugins.FlagHost, flagURL)
			flags.EnsureAlterableFlags(cmd, "warning", "critical")

			if err := opts.complete(cmd); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
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

	c.Flags().StringP(plugins.FlagHost, "H", "", "Icinga host name")
	c.Flags().StringVarP(&opts.url, flagURL, "u", "", "URL to get data")
	c.Flags().StringVarP(&opts.secretName, "secretName", "s", "", `Kubernetes secret name`)
	c.Flags().StringVarP(&opts.warning, "warning", "w", "", `Warning jsonpath query which returns [true/false]`)
	c.Flags().StringVarP(&opts.critical, "critical", "c", "", `Critical jsonpath query which returns [true/false]`)
	return c
}
