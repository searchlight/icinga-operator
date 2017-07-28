package icinga

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/appscode/envconfig"
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/log"
	"github.com/cloudflare/cfssl/cli"
	"github.com/cloudflare/cfssl/cli/genkey"
	"github.com/cloudflare/cfssl/cli/sign"
	"github.com/cloudflare/cfssl/config"
	"github.com/cloudflare/cfssl/csr"
	"github.com/cloudflare/cfssl/initca"
	"github.com/cloudflare/cfssl/signer"
	"gopkg.in/ini.v1"
)

const (
	ICINGA_ADDRESS         = "ICINGA_ADDRESS" // host:port
	ICINGA_API_USER        = "ICINGA_API_USER"
	ICINGA_API_PASSWORD    = "ICINGA_API_PASSWORD"
	ICINGA_CA_CERT         = "ICINGA_CA_CERT"
	ICINGA_SERVER_KEY      = "ICINGA_SERVER_KEY"
	ICINGA_SERVER_CERT     = "ICINGA_SERVER_CERT"
	ICINGA_IDO_HOST        = "ICINGA_IDO_HOST"
	ICINGA_IDO_PORT        = "ICINGA_IDO_PORT"
	ICINGA_IDO_DB          = "ICINGA_IDO_DB"
	ICINGA_IDO_USER        = "ICINGA_IDO_USER"
	ICINGA_IDO_PASSWORD    = "ICINGA_IDO_PASSWORD"
	ICINGA_WEB_HOST        = "ICINGA_WEB_HOST"
	ICINGA_WEB_PORT        = "ICINGA_WEB_PORT"
	ICINGA_WEB_DB          = "ICINGA_WEB_DB"
	ICINGA_WEB_USER        = "ICINGA_WEB_USER"
	ICINGA_WEB_PASSWORD    = "ICINGA_WEB_PASSWORD"
	ICINGA_WEB_UI_PASSWORD = "ICINGA_WEB_UI_PASSWORD"
)

var (
	// Key -> Required (true) | Optional (false)
	icingaKeys = map[string]bool{
		ICINGA_ADDRESS:         false,
		ICINGA_CA_CERT:         true,
		ICINGA_API_USER:        true,
		ICINGA_API_PASSWORD:    true,
		ICINGA_SERVER_KEY:      false,
		ICINGA_SERVER_CERT:     false,
		ICINGA_IDO_HOST:        true,
		ICINGA_IDO_PORT:        true,
		ICINGA_IDO_DB:          true,
		ICINGA_IDO_USER:        true,
		ICINGA_IDO_PASSWORD:    true,
		ICINGA_WEB_HOST:        true,
		ICINGA_WEB_PORT:        true,
		ICINGA_WEB_DB:          true,
		ICINGA_WEB_USER:        true,
		ICINGA_WEB_PASSWORD:    true,
		ICINGA_WEB_UI_PASSWORD: true,
	}
)

func init() {
	ini.PrettyFormat = false
}

type Configurator struct {
	ConfigRoot       string
	IcingaSecretName string
	Expiry           time.Duration
}

func (c *Configurator) ConfigFile() string {
	return filepath.Join(c.ConfigRoot, "searchlight/config.ini")
}

func (c *Configurator) PKIDir() string {
	return filepath.Join(c.ConfigRoot, "searchlight/pki")
}

func (c *Configurator) certFile(name string) string {
	return filepath.Join(c.PKIDir(), strings.ToLower(name)+".crt")
}

func (c *Configurator) keyFile(name string) string {
	return filepath.Join(c.PKIDir(), strings.ToLower(name)+".key")
}

// Returns PHID, cert []byte, key []byte, error
func (c *Configurator) initCA() error {
	certReq := &csr.CertificateRequest{
		CN: "searchlight-operator",
		Hosts: []string{
			"127.0.0.1",
		},
		KeyRequest: csr.NewBasicKeyRequest(),
		CA: &csr.CAConfig{
			PathLength: 2,
			Expiry:     c.Expiry.String(),
		},
	}

	cert, _, key, err := initca.New(certReq)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.certFile("ca"), cert, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.keyFile("ca"), key, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c *Configurator) createClientCert(csrReq *csr.CertificateRequest) error {
	g := &csr.Generator{Validator: genkey.Validator}
	csrPem, key, err := g.ProcessRequest(csrReq)
	if err != nil {
		return err
	}

	var cfg cli.Config
	cfg.CAKeyFile = c.keyFile("ca")
	cfg.CAFile = c.certFile("ca")
	cfg.CFG = &config.Config{
		Signing: &config.Signing{
			Profiles: map[string]*config.SigningProfile{},
			Default:  config.DefaultConfig(),
		},
	}
	cfg.CFG.Signing.Default.Expiry = c.Expiry
	cfg.CFG.Signing.Default.ExpiryString = c.Expiry.String()

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

	err = ioutil.WriteFile(c.certFile(csrReq.CN), cert, 0644)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(c.keyFile(csrReq.CN), key, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (c *Configurator) generateCertificates() error {
	err := os.MkdirAll(c.PKIDir(), 0755)
	if err != nil {
		return err
	}
	err = c.initCA()
	if err != nil {
		return err
	}
	log.Infoln("Created CA cert")

	var csrReq csr.CertificateRequest
	csrReq.KeyRequest = &csr.BasicKeyRequest{A: "rsa", S: 2048} // ECC does not work with icinga api
	csrReq.CN = "icinga"
	csrReq.Hosts = []string{"127.0.0.1"} // Add all local IPs
	return c.createClientCert(&csrReq)
}

func (c *Configurator) LoadConfig(userInput envconfig.LoaderFunc) (*Config, error) {
	if _, err := os.Stat(c.ConfigFile()); os.IsNotExist(err) {
		// auto generate the file
		cfg := ini.Empty()
		sec := cfg.Section("")
		sec.NewKey(ICINGA_ADDRESS, "127.0.0.1:5665")
		sec.NewKey(ICINGA_API_USER, "icingaapi")
		if v, ok := userInput(ICINGA_API_PASSWORD); ok {
			sec.NewKey(ICINGA_API_PASSWORD, v)
		} else {
			sec.NewKey(ICINGA_API_PASSWORD, rand.GeneratePassword())
		}

		caCert, caCertOK := userInput(ICINGA_CA_CERT)
		serverCert, serverCertOK := userInput(ICINGA_SERVER_CERT)
		serverKey, serverKeyOK := userInput(ICINGA_SERVER_KEY)
		if caCertOK && serverCertOK && serverKeyOK {
			err = os.MkdirAll(c.PKIDir(), 0755)
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(c.certFile("ca"), []byte(caCert), 0755)
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(c.certFile("icinga"), []byte(serverCert), 0755)
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(c.keyFile("icinga"), []byte(serverKey), 0644)
			if err != nil {
				return nil, err
			}
		} else if !caCertOK && !serverCertOK && !serverKeyOK {
			err = c.generateCertificates()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("Only some certs were provided.")
		}
		sec.NewKey(ICINGA_CA_CERT, c.certFile("ca"))
		sec.NewKey(ICINGA_SERVER_CERT, c.certFile("icinga"))
		sec.NewKey(ICINGA_SERVER_KEY, c.keyFile("icinga"))

		sec.NewKey(ICINGA_IDO_HOST, "127.0.0.1")
		sec.NewKey(ICINGA_IDO_PORT, "5432")
		sec.NewKey(ICINGA_IDO_DB, "icingaidodb")
		sec.NewKey(ICINGA_IDO_USER, "icingaido")
		if v, ok := userInput(ICINGA_IDO_PASSWORD); ok {
			sec.NewKey(ICINGA_IDO_PASSWORD, v)
		} else {
			sec.NewKey(ICINGA_IDO_PASSWORD, rand.GeneratePassword())
		}
		sec.NewKey(ICINGA_WEB_HOST, "127.0.0.1")
		sec.NewKey(ICINGA_WEB_PORT, "5432")
		sec.NewKey(ICINGA_WEB_DB, "icingawebdb")
		sec.NewKey(ICINGA_WEB_USER, "icingaweb")
		if v, ok := userInput(ICINGA_WEB_PASSWORD); ok {
			sec.NewKey(ICINGA_WEB_PASSWORD, v)
		} else {
			sec.NewKey(ICINGA_WEB_PASSWORD, rand.GeneratePassword())
		}
		if v, ok := userInput(ICINGA_WEB_UI_PASSWORD); ok {
			sec.NewKey(ICINGA_WEB_UI_PASSWORD, v)
		} else {
			sec.NewKey(ICINGA_WEB_UI_PASSWORD, rand.GeneratePassword())
		}

		err = os.MkdirAll(filepath.Dir(c.ConfigFile()), 0755)
		if err != nil {
			return nil, err
		}
		err = cfg.SaveTo(c.ConfigFile())
		if err != nil {
			return nil, err
		}
	}

	cfg, err := ini.Load(c.ConfigFile())
	if err != nil {
		return nil, err
	}
	sec := cfg.Section("")
	for key, required := range icingaKeys {
		if required && !sec.HasKey(key) {
			return nil, fmt.Errorf("No Icinga config found for key %s", key)
		}
	}

	addr := "127.0.0.1:5665"
	if key, err := sec.GetKey(ICINGA_ADDRESS); err == nil {
		addr = key.Value()
	}
	ctx := &Config{
		Endpoint: fmt.Sprintf("https://%s/v1", addr),
	}
	if key, err := sec.GetKey(ICINGA_API_USER); err == nil {
		ctx.BasicAuth.Username = key.Value()
	}
	if key, err := sec.GetKey(ICINGA_API_PASSWORD); err == nil {
		ctx.BasicAuth.Password = key.Value()
	}

	if caCert, err := ioutil.ReadFile(c.certFile("ca")); err == nil {
		ctx.CACert = caCert
	}

	return ctx, nil
}
