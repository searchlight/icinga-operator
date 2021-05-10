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
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var _ = Describe("check_json_path", func() {
	var secret *core.Secret
	var client corev1.SecretInterface
	var ts *httptest.Server
	BeforeEach(func() {

	})

	AfterEach(func() {
		if client != nil {
			client.Delete(context.TODO(), secret.Name, metav1.DeleteOptions{})
		}
		if ts != nil {
			ts.Close()
		}
	})

	Describe("when server return", func() {
		Context("404", func() {
			JustBeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, `{"status": "404"}`)
				}))
			})
			It("check for 404", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.status}==404",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
	})

	Describe("when server return json data", func() {
		Context("check string", func() {
			JustBeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, jsonDataAsServerOutput)
				}))
			})
			It("1st Book Category", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.Book[0].Category}==reference",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
			It("1st Book Category", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.Book[0].Category}==novel",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
		})
		Context("check float", func() {
			JustBeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, jsonDataAsServerOutput)
				}))
			})
			It("1st Book Price", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.Book[0].Price}==8.95",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
			It("1st Book Price > 10", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.Book[0].Price} > 10",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
			It("1st Book Price < 10", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.Book[0].Price} < 10",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
			It("1st Book Price < 10", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.Book[0].Price} < 5",
					warning:  "{.Book[0].Price} < 10",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Warning))
			})
		})
		Context("check boolean", func() {
			JustBeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, jsonDataAsServerOutput)
				}))
			})
			It("1st Bicycle IsNew", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.Bicycle[0].IsNew} != true",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
			It("2ns Bicycle IsNew", func() {
				opts := options{
					url: ts.URL,

					critical: "{.Bicycle[1].IsNew} == false",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
		Context("check multiple", func() {
			JustBeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintln(w, jsonDataAsServerOutput)
				}))
			})
			It("logical operator", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.Bicycle[0].Color} == red && {.Bicycle[0].Price} < 20",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
			It("-", func() {
				opts := options{
					url:      ts.URL,
					critical: "{.Bicycle[0].Color} != {.Bicycle[1].Color}",
				}
				state, _ := newPlugin(nil, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
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

var jsonDataAsServerOutput = `
{  
   "Book":[  
      {  
         "Category":"reference",
         "Author":"Nigel Rees",
         "Title":"Sayings of the Centurey",
         "Price":8.95
      }
   ],
   "Bicycle":[  
      {  
         "Color":"red",
         "Price":19.95,
         "IsNew":true
      },
      {  
         "Color":"green",
         "Price":20.01,
         "IsNew":false
      }
   ]
}
`
