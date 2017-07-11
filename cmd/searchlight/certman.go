package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/appscode/go/flags"
	"github.com/appscode/log"
	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/signer"
	"github.com/spf13/cobra"
)

func NewCmdGenerate() *cobra.Command {
	mgr := &ConfigLoader{
		Expiry: 10 * 365 * 24 * time.Hour,
	}
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate certificates for Kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.SetLogLevel(4)
			flags.EnsureRequiredFlags(cmd, "folder", "namespace", "cluster", "master-external-ip", "master-internal-ip")

			return mgr.GenClusterCerts()
		},
	}

	cmd.Flags().StringVarP(&configDir, "config-dir", "s", configDir, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")
	cmd.Flags().StringVar(&mgr.Folder, "folder", mgr.Folder, "Folder where certs are stored")

	return cmd
}

type ConfigLoader struct {
	Folder string
	Expiry time.Duration
}

func (opt *ConfigLoader) certFile(name string) string {
	return fmt.Sprintf("%s/%s.crt", opt.Folder, strings.ToLower(name))
}

func (opt *ConfigLoader) keyFile(name string) string {
	return fmt.Sprintf("%s/%s.key", opt.Folder, strings.ToLower(name))
}

// Returns PHID, cert []byte, key []byte, error
func (opt *ConfigLoader) initCA() error {
	certReq := &csr.CertificateRequest{
		CN: "searchlight-operator",
		Hosts: []string{
			"127.0.0.1",
		},
		KeyRequest: csr.NewBasicKeyRequest(),
		CA: &csr.CAConfig{
			PathLength: 2,
			Expiry:     opt.Expiry.String(),
		},
	}

	cert, _, key, err := initca.New(certReq)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(opt.certFile("ca"), cert, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(opt.keyFile("ca"), key, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (opt *ConfigLoader) createClientCert(csrReq *csr.CertificateRequest) error {
	g := &csr.Generator{Validator: genkey.Validator}
	csrPem, key, err := g.ProcessRequest(csrReq)
	if err != nil {
		return err
	}

	var cfg cli.Config
	cfg.CAKeyFile = opt.keyFile("ca")
	cfg.CAFile = opt.certFile("ca")
	cfg.CFG = &config.Config{
		Signing: &config.Signing{
			Profiles: map[string]*config.SigningProfile{},
			Default:  config.DefaultConfig(),
		},
	}
	cfg.CFG.Signing.Default.Expiry = opt.Expiry
	cfg.CFG.Signing.Default.ExpiryString = opt.Expiry.String()

	s, err := sign.SignerFromConfig(cfg)
	if err != nil {
		return err
	}
	var cert []byte
	signReq := signer.SignRequest{
		Request: string(csrPem),
		Hosts:   signer.SplitHosts(cfg.Hostname),
		Profile: cfg.Profile,
		Label:   cfg.Label,
	}

	cert, err = s.Sign(signReq)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(opt.certFile(csrReq.CN), cert, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(opt.keyFile(csrReq.CN), key, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (opt *ConfigLoader) GenClusterCerts() error {
	var csrReq csr.CertificateRequest
	csrReq.KeyRequest = csr.NewBasicKeyRequest() // &csr.BasicKeyRequest{A: "rsa", S: 2048}

	err := os.MkdirAll(opt.Folder, 0755)
	if err != nil {
		return err
	}

	err = opt.initCA()
	if err != nil {
		return err
	}
	log.Infoln("Created CA cert")

	csrReq.CN = "icinga"
	csrReq.Hosts = []string{"127.0.0.1"} // Add all local IPs
	err = opt.createClientCert(&csrReq)
	if err != nil {
		return err
	}
	log.Infoln("Created icinga cert")
	return nil
}
