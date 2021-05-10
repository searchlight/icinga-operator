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

package icinga

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gomodules.xyz/blobfs"
	"gomodules.xyz/cert"
	"gomodules.xyz/cert/certstore"
	"gomodules.xyz/envconfig"
	"gomodules.xyz/x/crypto/rand"
	ini "gopkg.in/ini.v1"
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

func (c *Configurator) LoadConfig(userInput envconfig.LoaderFunc) (*Config, error) {
	fs := blobfs.NewOsFs()
	pkidir := filepath.Join(c.ConfigRoot, "searchlight/pki")
	store, err := certstore.New(fs, pkidir)
	if err != nil {
		return nil, err
	}
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
			err = fs.WriteFile(context.TODO(), store.CertFile("ca"), []byte(caCert)) // 0755
			if err != nil {
				return nil, err
			}
			err = fs.WriteFile(context.TODO(), store.CertFile("icinga"), []byte(serverCert)) // 0755
			if err != nil {
				return nil, err
			}
			err = fs.WriteFile(context.TODO(), store.KeyFile("icinga"), []byte(serverKey)) // 0644
			if err != nil {
				return nil, err
			}
		} else if !caCertOK && !serverCertOK && !serverKeyOK {
			err = store.NewCA()
			if err != nil {
				return nil, err
			}
			sans := cert.AltNames{
				DNSNames: []string{"icinga"},
				IPs:      []net.IP{net.ParseIP("127.0.0.1")},
			}
			serverCert, serverKey, err := store.NewServerCertPairBytes(sans)
			if err != nil {
				return nil, err
			}
			err = fs.WriteFile(context.TODO(), store.CertFile("icinga"), []byte(serverCert)) // 0755
			if err != nil {
				return nil, err
			}
			err = fs.WriteFile(context.TODO(), store.KeyFile("icinga"), []byte(serverKey)) // 0644
			if err != nil {
				return nil, err
			}
		} else {
			return nil, errors.New("only some certs were provided")
		}
		sec.NewKey(ICINGA_CA_CERT, store.CertFile("ca"))
		sec.NewKey(ICINGA_SERVER_CERT, store.CertFile("icinga"))
		sec.NewKey(ICINGA_SERVER_KEY, store.KeyFile("icinga"))

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
			return nil, fmt.Errorf("no Icinga config found for key %s", key)
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

	if store.IsExists("ca") {
		ctx.CACert = store.CACertBytes()
	}

	return ctx, nil
}
