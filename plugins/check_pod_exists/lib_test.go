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

package check_pod_exists

import (
	"context"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var _ = Describe("check_pod_exists", func() {
	var pod, pod2 *core.Pod
	var client corev1.PodInterface
	var opts options

	BeforeEach(func() {
		client = cs.CoreV1().Pods("demo")
		pod = &core.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pod",
				Labels: map[string]string{
					"app/searchlight": "pod",
				},
			},
		}
	})

	AfterEach(func() {
		if client != nil {
			client.Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		}
	})

	Describe("when a single pod exists", func() {
		Context("with pod name", func() {
			JustBeforeEach(func() {
				opts = options{
					podName: pod.Name,
				}
			})
			It("should be OK", func() {
				_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
		})
	})
	Describe("when two pod exist", func() {
		JustBeforeEach(func() {
			pod2 = &core.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "pod-2",
					Labels: map[string]string{
						"app/searchlight": "pod",
					},
				},
			}
		})
		AfterEach(func() {
			if client != nil {
				client.Delete(context.TODO(), pod2.Name, metav1.DeleteOptions{})
			}
		})
		Context("without selector", func() {
			Context("with count", func() {
				JustBeforeEach(func() {
					opts = options{
						count:      2,
						isCountSet: true,
					}
				})
				It("greater than actual", func() {
					_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), pod2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.count = opts.count + 1
					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.Critical))
				})
				It("less than actual", func() {
					_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), pod2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.count = opts.count - 1
					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.Critical))
				})
				It("similar to actual", func() {
					_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), pod2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.OK))
				})
			})
			Context("without count", func() {
				It("should be OK", func() {
					_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), pod2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.OK))
				})
			})
		})

		Context("with selector", func() {
			Context("with count", func() {
				JustBeforeEach(func() {
					opts = options{
						count:      2,
						isCountSet: true,
						selector:   labels.SelectorFromSet(pod.Labels).String(),
					}
				})
				It("greater than actual", func() {
					_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), pod2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.count = opts.count + 1
					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.Critical))
				})
				It("less than actual", func() {
					_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), pod2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.count = opts.count - 1
					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.Critical))
				})
				It("similar to actual", func() {
					_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), pod2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.OK))
				})
			})
			Context("without count", func() {
				JustBeforeEach(func() {
					opts = options{
						selector: labels.SelectorFromSet(pod.Labels).String(),
					}
				})
				It("should be OK", func() {
					_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), pod2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.OK))
				})
			})
		})
	})
	Describe("test options", func() {
		var (
			cmd *cobra.Command
		)

		JustBeforeEach(func() {
			cmd = new(cobra.Command)
			cmd.Flags().Int(flagCount, 0, "")
			cmd.Flags().String(plugins.FlagHost, "demo@cluster", "")
			cmd.Flags().String(plugins.FlagKubeConfig, "", "")
			cmd.Flags().String(plugins.FlagKubeConfigContext, "", "")
		})
		Context("valid", func() {
			It("-", func() {
				opts := options{}
				cmd.Flags().Set(flagCount, "2")
				cmd.Execute()
				err := opts.complete(cmd)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(opts.isCountSet).Should(BeIdenticalTo(true))
				Expect(opts.namespace).Should(BeIdenticalTo("demo"))
				err = opts.validate()
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("-", func() {
				opts := options{}
				err := opts.complete(cmd)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(opts.isCountSet).Should(BeIdenticalTo(false))
				err = opts.validate()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
