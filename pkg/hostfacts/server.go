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

package hostfacts

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-macaron/auth"
	"github.com/go-macaron/toolbox"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	macaron "gopkg.in/macaron.v1"
	"k8s.io/klog/v2"
)

type Server struct {
	Address    string
	CACertFile string
	CertFile   string
	KeyFile    string
	Username   string
	Password   string
	Token      string
}

func (s Server) ListenAndServe() {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())

	// auth
	if s.Username != "" && s.Password != "" {
		m.Use(auth.Basic(s.Username, s.Password))
	} else if s.Token != "" {
		m.Use(auth.Bearer(s.Token))
	}

	m.Use(toolbox.Toolboxer(m))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		IndentJSON: true,
	}))

	m.Get("/cpu", func(ctx *macaron.Context) {
		r, _ := cpu.Info()
		ctx.JSON(200, r)
	})

	m.Get("/virt_mem", func(ctx *macaron.Context) {
		r, _ := mem.VirtualMemory()
		ctx.JSON(200, r)
	})
	m.Get("/swap_mem", func(ctx *macaron.Context) {
		r, _ := mem.SwapMemory()
		ctx.JSON(200, r)
	})

	m.Get("/disks", func(ctx *macaron.Context) {
		r, _ := disk.Partitions(true)
		ctx.JSON(200, r)
	})

	m.Get("/du", func(ctx *macaron.Context) {
		paths := ctx.QueryStrings("p")
		du := make([]*disk.UsageStat, len(paths))
		for i, p := range paths {
			du[i], _ = disk.Usage(p)
		}
		ctx.JSON(200, du)
	})

	m.Get("/load", func(ctx *macaron.Context) {
		l, _ := load.Avg()
		ctx.JSON(200, l)
	})

	m.Get("/metrics", promhttp.Handler().ServeHTTP)

	macaron.Env = macaron.PROD
	klog.Infof("listening on %s (%s)\n", s.Address, macaron.Env)

	srv := &http.Server{
		Addr:         s.Address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      m,
	}
	if s.CACertFile == "" && s.CertFile == "" && s.KeyFile == "" {
		klog.Fatalln(srv.ListenAndServe())
	} else {
		/*
			Ref:
			 - https://blog.cloudflare.com/exposing-go-on-the-internet/
			 - http://www.bite-code.com/2015/06/25/tls-mutual-auth-in-golang/
			 - http://www.hydrogen18.com/blog/your-own-pki-tls-golang.html
		*/
		tlsConfig := &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			SessionTicketsDisabled:   true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			ClientAuth: tls.VerifyClientCertIfGiven,
			NextProtos: []string{"h2", "http/1.1"},
		}
		if s.CACertFile != "" {
			caCert, err := ioutil.ReadFile(s.CACertFile)
			if err != nil {
				klog.Fatal(err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.ClientCAs = caCertPool
		}
		tlsConfig.BuildNameToCertificate()

		srv.TLSConfig = tlsConfig
		klog.Fatalln(srv.ListenAndServeTLS(s.CertFile, s.KeyFile))
	}
}
