package check_ca_cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins"
	"github.com/spf13/cobra"
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

func (p *plugin) loadCACert() (*x509.Certificate, error) {
	caCert := "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	data, err := ioutil.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse certificate")
	}
	return x509.ParseCertificate(block.Bytes)
}

func (p *plugin) Check() (icinga.State, interface{}) {
	crt, err := p.loadCACert()
	if err != nil {
		return icinga.Unknown, err.Error()
	}

	remaining := crt.NotAfter.Sub(time.Now())

	if remaining.Seconds() < p.options.critical.Seconds() {
		return icinga.Critical, fmt.Sprintf("Certificate will be expired within %v hours", remaining.Hours())
	}

	if remaining.Seconds() < p.options.warning.Seconds() {
		return icinga.Warning, fmt.Sprintf("Certificate will be expired within %v hours", remaining.Hours())
	}

	return icinga.OK, fmt.Sprintf("Certificate is valid more than %v days", remaining.Hours()/24.0)
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
