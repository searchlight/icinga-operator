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

package check_cert

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"time"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	"github.com/spf13/cobra"
	"gomodules.xyz/cert"
	"gomodules.xyz/x/flags"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
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
	// options for Secret
	namespace  string
	selector   string
	secretName string
	secretKey  []string
	// Certificate expirity duration
	warning  time.Duration
	critical time.Duration
	// IcingaHost
	host *icinga.IcingaHost
}

func (o *options) complete(cmd *cobra.Command) error {
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
		return err
	}
	o.contextName, err = cmd.Flags().GetString(plugins.FlagKubeConfigContext)
	if err != nil {
		return err
	}
	return nil
}

func (o *options) validate() error {
	if o.host.Type != icinga.TypeCluster {
		return errors.New("invalid icinga host type")
	}
	return nil
}

func (p *plugin) getCertSecrets() ([]core.Secret, error) {
	opts := p.options
	if opts.secretName != "" {
		var secret *core.Secret
		secret, err := p.client.Get(context.TODO(), opts.secretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return []core.Secret{*secret}, nil
	}

	secretList, err := p.client.List(context.TODO(), metav1.ListOptions{
		LabelSelector: opts.selector,
	})
	if err != nil {
		return nil, err
	}
	return secretList.Items, nil
}

func (p *plugin) checkNotAfter(cert *x509.Certificate) (icinga.State, time.Duration) {
	remaining := cert.NotAfter.Sub(time.Now())
	if remaining.Seconds() < p.options.critical.Seconds() {
		return icinga.Critical, remaining
	}

	if remaining.Seconds() < p.options.warning.Seconds() {
		return icinga.Warning, remaining
	}

	return icinga.OK, remaining
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
		if state, remaining := p.checkNotAfter(cert); state != icinga.OK {
			return state, fmt.Errorf(
				`certificate found in key "%s" in Secret "%s/%s" will be expired within %v hours`,
				key, secret.Namespace, secret.Name, remaining.Hours(),
			)
		}
	}
	return icinga.OK, nil
}

func (p *plugin) checkCertPerSecretKey(secret *core.Secret) (icinga.State, error) {
	opts := p.options
	for _, key := range opts.secretKey {
		data, ok := secret.Data[key]
		if !ok {
			return icinga.Warning, fmt.Errorf(`key "%s" not found in Secret "%s/%s"`, key, secret.Namespace, secret.Name)
		}

		if state, err := p.checkCert(data, secret, key); err != nil {
			return state, err
		}
	}

	if len(opts.secretKey) == 0 && secret.Type == core.SecretTypeTLS {
		data, ok := secret.Data[core.TLSCertKey]
		if !ok {
			return icinga.Warning, fmt.Errorf(`key "%s" not found in Secret "%s/%s"`, core.TLSCertKey, secret.Namespace, secret.Name)
		}

		if state, err := p.checkCert(data, secret, core.TLSCertKey); err != nil {
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
	var opts options

	cmd := &cobra.Command{
		Use:   "check_cert",
		Short: "Check Certificate expire date",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, plugins.FlagHost)

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

	cmd.Flags().StringP(plugins.FlagHost, "H", "", "Icinga host name")
	cmd.Flags().StringVarP(&opts.selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='")
	cmd.Flags().StringVarP(&opts.secretName, "secretName", "s", "", "Name of secret from where certificates are checked")
	cmd.Flags().StringSliceVarP(&opts.secretKey, "secretKey", "k", nil, "Name of secret key where certificates are kept")
	cmd.Flags().DurationVarP(&opts.warning, "warning", "w", time.Hour*360, `Remaining duration for Warning state. [Default: 360h]`)
	cmd.Flags().DurationVarP(&opts.critical, "critical", "c", time.Hour*120, `Remaining duration for Critical state. [Default: 120h]`)
	return cmd
}
