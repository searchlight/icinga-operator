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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	"gomodules.xyz/cert"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func generateCertificate(ttl time.Duration) ([]byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"AppsCode Inc."},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(ttl),
	}

	certByte, err := x509.CreateCertificate(rand.Reader, &template, &template, priv.Public(), priv)
	if err != nil {
		return nil, err
	}

	certificate, err := x509.ParseCertificate(certByte)
	if err != nil {
		return nil, err
	}
	return cert.EncodeCertPEM(certificate), nil
}

var _ = Describe("check_cert", func() {
	var secret, secret2 *core.Secret
	var client corev1.SecretInterface
	var opts options

	BeforeEach(func() {
		secret = &core.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "secret-1",
				Namespace: "demo",
				Labels:    make(map[string]string),
			},
			Data: make(map[string][]byte),
		}
	})

	AfterEach(func() {
		if client != nil {
			client.Delete(context.TODO(), secret.Name, metav1.DeleteOptions{})
		}
	})

	Describe("when there is a valid certificate", func() {
		JustBeforeEach(func() {
			client = cs.CoreV1().Secrets(secret.Namespace)
			opts = options{
				namespace:  secret.Namespace,
				secretName: secret.Name,
				warning:    time.Hour * 7 * 24,
				critical:   time.Hour * 5 * 24,
			}
		})
		Context("with nearly expire", func() {
			It("with in 4 days", func() {
				cert, err := generateCertificate(time.Hour * 4 * 24)
				Expect(err).ShouldNot(HaveOccurred())
				secret.Data["client.cert"] = cert
				_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				opts.secretKey = []string{"client.cert"}

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
			It("with in 6 days", func() {
				cert, err := generateCertificate(time.Hour * 6 * 24)
				Expect(err).ShouldNot(HaveOccurred())
				secret.Data["client.cert"] = cert
				_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				opts.secretKey = []string{"client.cert"}

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Warning))
			})
			It("with in 10 days", func() {
				cert, err := generateCertificate(time.Hour * 10 * 24)
				Expect(err).ShouldNot(HaveOccurred())
				secret.Data["client.cert"] = cert
				_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				opts.secretKey = []string{"client.cert"}

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
		})

		Context("with label", func() {
			JustBeforeEach(func() {
				secret.Labels["app/searchlight"] = "cert"
				opts.secretName = ""
				opts.selector = labels.SelectorFromSet(secret.Labels).String()
			})
			It("with in 10 days", func() {
				cert, err := generateCertificate(time.Hour * 10 * 24)
				Expect(err).ShouldNot(HaveOccurred())
				secret.Data["client.cert"] = cert
				_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				opts.secretKey = []string{"client.cert"}

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})

			Context("for multiple secret", func() {
				JustBeforeEach(func() {
					secret.Labels["app/searchlight"] = "cert"
					opts.secretName = ""
					opts.selector = labels.SelectorFromSet(secret.Labels).String()

					// Another Secret
					secret2 = &core.Secret{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "secret-2",
							Namespace: "demo",
							Labels: map[string]string{
								"app/searchlight": "cert",
							},
						},
						Data: make(map[string][]byte),
					}
				})
				AfterEach(func() {
					client.Delete(context.TODO(), secret2.Name, metav1.DeleteOptions{})
				})
				It("both contain valid certificate", func() {
					cert, err := generateCertificate(time.Hour * 10 * 24)
					Expect(err).ShouldNot(HaveOccurred())
					secret.Data["client.cert"] = cert
					_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					cert, err = generateCertificate(time.Hour * 10 * 24)
					Expect(err).ShouldNot(HaveOccurred())
					secret2.Data["client.cert"] = cert
					_, err = client.Create(context.TODO(), secret2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.secretKey = []string{"client.cert"}

					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.OK))
				})

				It("one contains nearly expire certificate", func() {
					cert, err := generateCertificate(time.Hour * 10 * 24)
					Expect(err).ShouldNot(HaveOccurred())
					secret.Data["client.cert"] = cert
					_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					cert, err = generateCertificate(time.Hour * 2 * 24)
					Expect(err).ShouldNot(HaveOccurred())
					secret2.Data["client.cert"] = cert
					_, err = client.Create(context.TODO(), secret2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.secretKey = []string{"client.cert"}

					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.Critical))
				})
			})
		})
	})

	Describe("when there is an invalid secret", func() {
		JustBeforeEach(func() {
			client = cs.CoreV1().Secrets(secret.Namespace)
			opts = options{
				namespace:  secret.Namespace,
				secretName: secret.Name,
			}
		})
		Context("with invalid certificate", func() {
			It("should be Unknown", func() {
				cert, err := generateCertificate(time.Hour * 10 * 24)
				cert[0] = cert[0] + 1
				Expect(err).ShouldNot(HaveOccurred())
				secret.Data["client.cert"] = cert
				_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				opts.secretKey = []string{"client.cert"}

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Unknown))
			})
		})
		Context("with invalid secret key", func() {
			It("should be warning", func() {
				_, err := client.Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				opts.secretKey = []string{"unknown.cert"}

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Warning))
			})
		})
	})

	Describe("when there is a tls secret", func() {
		JustBeforeEach(func() {
			client = cs.CoreV1().Secrets(secret.Namespace)
			secret.Type = core.SecretTypeTLS
			opts = options{
				namespace:  secret.Namespace,
				secretName: secret.Name,
			}
		})
		Context("with valid certificate", func() {
			It("with key", func() {
				cert, err := generateCertificate(time.Hour * 10 * 24)
				Expect(err).ShouldNot(HaveOccurred())
				secret.Data["tls.crt"] = cert
				_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				opts.secretKey = []string{"tls.crt"}

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
			It("without key", func() {
				cert, err := generateCertificate(time.Hour * 10 * 24)
				Expect(err).ShouldNot(HaveOccurred())
				secret.Data["tls.crt"] = cert
				_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
			It("with invalid key", func() {
				cert, err := generateCertificate(time.Hour * 10 * 24)
				Expect(err).ShouldNot(HaveOccurred())
				secret.Data["tls-invalid.crt"] = cert
				_, err = client.Create(context.TODO(), secret, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Warning))
			})
		})
	})

	Describe("test options", func() {
		var (
			cmd *cobra.Command
		)

		JustBeforeEach(func() {
			cmd = new(cobra.Command)
			cmd.Flags().String(plugins.FlagHost, "", "")
			cmd.Flags().String(plugins.FlagKubeConfig, "", "")
			cmd.Flags().String(plugins.FlagKubeConfigContext, "", "")
		})
		Context("valid", func() {
			It("host", func() {
				opts := options{}
				cmd.Flags().Set(plugins.FlagHost, "demo@cluster")
				err := opts.complete(cmd)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(opts.namespace).Should(BeIdenticalTo("demo"))
				err = opts.validate()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		Context("invalid", func() {
			It("hostname", func() {
				opts := options{}
				cmd.Flags().Set(plugins.FlagHost, "demo@cluster@demo")
				err := opts.complete(cmd)
				Expect(err).Should(HaveOccurred())
			})
			It("host type", func() {
				opts := options{}
				cmd.Flags().Set(plugins.FlagHost, "demo@pod@demo")
				err := opts.complete(cmd)
				Expect(err).ShouldNot(HaveOccurred())
				err = opts.validate()
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
