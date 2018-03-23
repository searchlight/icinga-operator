package check_pod_status

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

var _ = Describe("check_pod_status", func() {
	var pod *core.Pod
	var client corev1.PodInterface
	var opts options

	BeforeEach(func() {
		client = cs.CoreV1().Pods("demo")
		pod = &core.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "pod",
			},
		}
		opts = options{
			podName: pod.Name,
		}
	})

	AfterEach(func() {
		if client != nil {
			client.Delete(pod.Name, &metav1.DeleteOptions{})
		}
	})

	Describe("there is a ready pod", func() {
		Context("with no other problems", func() {
			It("should be OK", func() {
				_, err := client.Create(pod)
				Expect(err).ShouldNot(HaveOccurred())
				pod.Status.Phase = core.PodRunning
				pod.Status.Conditions = []core.PodCondition{
					{
						Type:   core.PodReady,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(pod)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
		})
	})

	Describe("there is a not ready pod", func() {
		Context("with no other problems", func() {
			It("should be Critical", func() {
				_, err := client.Create(pod)
				Expect(err).ShouldNot(HaveOccurred())

				pod.Status.Phase = core.PodRunning
				pod.Status.Conditions = []core.PodCondition{
					{
						Type:   core.PodReady,
						Status: core.ConditionFalse,
					},
				}
				_, err = client.Update(pod)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
	})

	Describe("there is a not running pod", func() {
		Context("succeeded", func() {
			It("should be Critical", func() {
				_, err := client.Create(pod)
				Expect(err).ShouldNot(HaveOccurred())

				pod.Status.Phase = core.PodSucceeded
				_, err = client.Update(pod)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
		Context("failed", func() {
			It("should be Critical", func() {
				_, err := client.Create(pod)
				Expect(err).ShouldNot(HaveOccurred())

				pod.Status.Phase = core.PodFailed
				_, err = client.Update(pod)
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
				cmd.Flags().Set(plugins.FlagHost, "demo@pod")
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
				cmd.Flags().Set(plugins.FlagHost, "demo@pod@name")
				err := opts.complete(cmd)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(opts.podName).Should(BeIdenticalTo("name"))
				Expect(opts.namespace).Should(BeIdenticalTo("demo"))
				Expect(opts.host.Type).Should(BeIdenticalTo("pod"))
				err = opts.validate()
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
