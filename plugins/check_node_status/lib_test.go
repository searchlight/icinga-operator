package check_node_status

import (
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var _ = Describe("check_node_status", func() {
	var node *core.Node
	var client corev1.NodeInterface
	var opts options

	BeforeEach(func() {
		client = cs.CoreV1().Nodes()
		node = &core.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node",
			},
		}
		opts = options{
			nodeName: node.Name,
		}
	})

	AfterEach(func() {
		if client != nil {
			client.Delete(node.Name, &metav1.DeleteOptions{})
		}
	})

	Describe("there is a ready node", func() {
		Context("with no other problems", func() {
			It("should be OK", func() {
				_, err := client.Create(node)
				Expect(err).ShouldNot(HaveOccurred())

				node.Status.Conditions = []core.NodeCondition{
					{
						Type:   core.NodeReady,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(node)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
		})
		Context("with other problems", func() {
			It("such as OutOfDisk", func() {
				_, err := client.Create(node)
				Expect(err).ShouldNot(HaveOccurred())

				node.Status.Conditions = []core.NodeCondition{
					{
						Type:   core.NodeReady,
						Status: core.ConditionTrue,
					},
					{
						Type:   core.NodeOutOfDisk,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(node)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
			It("such as OutOfDisk, MemoryPressure", func() {
				_, err := client.Create(node)
				Expect(err).ShouldNot(HaveOccurred())

				node.Status.Conditions = []core.NodeCondition{
					{
						Type:   core.NodeReady,
						Status: core.ConditionTrue,
					},
					{
						Type:   core.NodeOutOfDisk,
						Status: core.ConditionTrue,
					},
					{
						Type:   core.NodeMemoryPressure,
						Status: core.ConditionTrue,
					},
					{
						Type:   core.NodeDiskPressure,
						Status: core.ConditionTrue,
					},
					{
						Type:   core.NodeNetworkUnavailable,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(node)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
	})

	Describe("there is a not ready node", func() {
		JustBeforeEach(func() {
			client = cs.CoreV1().Nodes()
			opts = options{
				nodeName: node.Name,
			}
		})
		Context("with no other problems", func() {
			It("should be Critical", func() {
				_, err := client.Create(node)
				Expect(err).ShouldNot(HaveOccurred())

				node.Status.Conditions = []core.NodeCondition{
					{
						Type:   core.NodeReady,
						Status: core.ConditionFalse,
					},
				}
				_, err = client.Update(node)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
		Context("with other problems", func() {
			It("such as OutOfDisk", func() {
				_, err := client.Create(node)
				Expect(err).ShouldNot(HaveOccurred())

				node.Status.Conditions = []core.NodeCondition{
					{
						Type:   core.NodeReady,
						Status: core.ConditionFalse,
					},
					{
						Type:   core.NodeOutOfDisk,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(node)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
			It("such as OutOfDisk, MemoryPressure", func() {
				_, err := client.Create(node)
				Expect(err).ShouldNot(HaveOccurred())

				node.Status.Conditions = []core.NodeCondition{
					{
						Type:   core.NodeReady,
						Status: core.ConditionFalse,
					},
					{
						Type:   core.NodeOutOfDisk,
						Status: core.ConditionTrue,
					},
					{
						Type:   core.NodeMemoryPressure,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(node)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
	})

	Describe("Check validation", func() {
		var (
			cmd *cobra.Command
		)

		JustBeforeEach(func() {
			cmd = new(cobra.Command)
			cmd.Flags().String(plugins.FlagHost, "", "")
			cmd.Flags().String(plugins.FlagKubeConfig, "", "")
			cmd.Flags().String(plugins.FlagKubeConfigContext, "", "")
		})

		Context("for invalid", func() {
			It("with invalid part", func() {
				opts := options{}
				cmd.Flags().Set(plugins.FlagHost, "demo@node")
				err := opts.complete(cmd)
				Expect(err).Should(HaveOccurred())
			})
			It("with invalid type", func() {
				opts := options{}
				cmd.Flags().Set(plugins.FlagHost, "demo@cluster")
				err := opts.complete(cmd)
				Expect(err).ShouldNot(HaveOccurred())
				err = opts.validate()
				Expect(err).Should(HaveOccurred())
			})
		})
		Context("for valid", func() {
			It("with valid name", func() {
				opts := options{}
				cmd.Flags().Set(plugins.FlagHost, "demo@node@name")
				err := opts.complete(cmd)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(opts.nodeName).Should(BeIdenticalTo("name"))
				Expect(opts.host.Type).Should(BeIdenticalTo("node"))
				err = opts.validate()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
