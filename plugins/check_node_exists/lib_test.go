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

package check_node_exists

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

var _ = Describe("check_node_exists", func() {
	var node, node2 *core.Node
	var client corev1.NodeInterface
	var opts options

	BeforeEach(func() {
		client = cs.CoreV1().Nodes()
		node = &core.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node",
				Labels: map[string]string{
					"app/searchlight": "node",
				},
			},
		}
	})

	AfterEach(func() {
		if client != nil {
			client.Delete(context.TODO(), node.Name, metav1.DeleteOptions{})
		}
	})

	Describe("when a single node exists", func() {
		Context("with node name", func() {
			JustBeforeEach(func() {
				opts = options{
					nodeName: node.Name,
				}
			})
			It("should be OK", func() {
				_, err := client.Create(context.TODO(), node, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
		})
	})
	Describe("when two node exist", func() {
		JustBeforeEach(func() {
			node2 = &core.Node{
				ObjectMeta: metav1.ObjectMeta{
					Name: "node-2",
					Labels: map[string]string{
						"app/searchlight": "node",
					},
				},
			}
		})
		AfterEach(func() {
			if client != nil {
				client.Delete(context.TODO(), node2.Name, metav1.DeleteOptions{})
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
					_, err := client.Create(context.TODO(), node, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), node2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.count = opts.count + 1
					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.Critical))
				})
				It("less than actual", func() {
					_, err := client.Create(context.TODO(), node, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), node2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.count = opts.count - 1
					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.Critical))
				})
				It("similar to actual", func() {
					_, err := client.Create(context.TODO(), node, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), node2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.OK))
				})
			})
			Context("without count", func() {
				It("should be OK", func() {
					_, err := client.Create(context.TODO(), node, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), node2, metav1.CreateOptions{})
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
						selector:   labels.SelectorFromSet(node.Labels).String(),
					}
				})
				It("greater than actual", func() {
					_, err := client.Create(context.TODO(), node, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), node2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.count = opts.count + 1
					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.Critical))
				})
				It("less than actual", func() {
					_, err := client.Create(context.TODO(), node, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), node2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					opts.count = opts.count - 1
					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.Critical))
				})
				It("similar to actual", func() {
					_, err := client.Create(context.TODO(), node, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), node2, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())

					state, _ := newPlugin(client, opts).Check()
					Expect(state).Should(BeIdenticalTo(icinga.OK))
				})
			})
			Context("without count", func() {
				JustBeforeEach(func() {
					opts = options{
						selector: labels.SelectorFromSet(node.Labels).String(),
					}
				})
				It("should be OK", func() {
					_, err := client.Create(context.TODO(), node, metav1.CreateOptions{})
					Expect(err).ShouldNot(HaveOccurred())
					_, err = client.Create(context.TODO(), node2, metav1.CreateOptions{})
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
			cmd.Flags().String(plugins.FlagKubeConfig, "", "")
			cmd.Flags().String(plugins.FlagKubeConfigContext, "", "")
		})
		Context("valid", func() {
			It("-", func() {
				opts := options{}
				cmd.Flags().Set(flagCount, "2")
				err := opts.complete(cmd)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(opts.isCountSet).Should(BeIdenticalTo(true))
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
