package check_ca_cert

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
)

type request struct {
	warning  time.Duration
	critical time.Duration
}

func loadCACert() (*x509.Certificate, error) {
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

func checkCertificate(req *request) (icinga.State, string) {
	crt, err := loadCACert()
	if err != nil {
		return icinga.Unknown, err.Error()
	}

	remaining := crt.NotAfter.Sub(time.Now())

	if remaining.Seconds() < req.critical.Seconds() {
		return icinga.Critical, fmt.Sprintf("Certificate will be expired within %v hours", remaining.Hours())
	}

	if remaining.Seconds() < req.warning.Seconds() {
		return icinga.Warning, fmt.Sprintf("Certificate will be expired within %v hours", remaining.Hours())
	}

	return icinga.OK, fmt.Sprintf("Certificate is valid more than %v days", remaining.Hours()/24.0)
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "check_ca_cert",
		Short:   "Check Certificate expire date",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			icinga.Output(checkCertificate(&req))
		},
	}

	c.Flags().DurationVarP(&req.warning, "warning", "w", time.Hour*360, `Remaining duration for warning state. [Default: 360h]`)
	c.Flags().DurationVarP(&req.critical, "critical", "c", time.Hour*120, `Remaining duration for critical state. [Default: 120h]`)
	return c
}
