package check_ca_cert

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins"
	"github.com/spf13/cobra"
	"gomodules.xyz/cert"
)

type plugin struct {
	options options
}

var _ plugins.PluginInterface = &plugin{}

func newPlugin(opts options) *plugin {
	return &plugin{opts}
}

type options struct {
	warning  time.Duration
	critical time.Duration
}

func (o *options) complete(cmd *cobra.Command) error {
	return nil
}

func (o *options) validate() error {
	return nil
}

func (p *plugin) loadCACert() ([]*x509.Certificate, error) {
	caCert := "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	data, err := os.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	return cert.ParseCertsPEM(data)
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

func (p *plugin) Check() (icinga.State, interface{}) {
	certs, err := p.loadCACert()
	if err != nil {
		return icinga.Unknown, err.Error()
	}

	for _, cert := range certs {
		if state, remaining := p.checkNotAfter(cert); state != icinga.OK {
			return state, fmt.Errorf(`certificate will be expired within %v hours`, remaining.Hours())
		}
	}
	return icinga.OK, nil
}

func NewCmd() *cobra.Command {
	var opts options

	c := &cobra.Command{
		Use:   "check_ca_cert",
		Short: "Check Certificate expire date",

		Run: func(cmd *cobra.Command, args []string) {
			if err := opts.complete(cmd); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			if err := opts.validate(); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			icinga.Output(newPlugin(opts).Check())
		},
	}

	c.Flags().DurationVarP(&opts.warning, "warning", "w", time.Hour*360, `Remaining duration for warning state. [Default: 360h]`)
	c.Flags().DurationVarP(&opts.critical, "critical", "c", time.Hour*120, `Remaining duration for critical state. [Default: 120h]`)
	return c
}
