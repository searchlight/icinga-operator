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

package check_pod_status

import (
	"context"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

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
			client.Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		}
	})

	Describe("there is a ready pod", func() {
		Context("with no other problems", func() {
			It("should be OK", func() {
				_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())
				pod.Status.Phase = core.PodRunning
				pod.Status.Conditions = []core.PodCondition{
					{
						Type:   core.PodReady,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(context.TODO(), pod, metav1.UpdateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
		})
	})

	Describe("there is a not ready pod", func() {
		Context("with no other problems", func() {
			It("should be Critical", func() {
				_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				pod.Status.Phase = core.PodRunning
				pod.Status.Conditions = []core.PodCondition{
					{
						Type:   core.PodReady,
						Status: core.ConditionFalse,
					},
				}
				_, err = client.Update(context.TODO(), pod, metav1.UpdateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
	})

	Describe("there is a not running pod", func() {
		Context("succeeded", func() {
			It("should be Critical", func() {
				_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				pod.Status.Phase = core.PodSucceeded
				_, err = client.Update(context.TODO(), pod, metav1.UpdateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
		Context("failed", func() {
			It("should be Critical", func() {
				_, err := client.Create(context.TODO(), pod, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				pod.Status.Phase = core.PodFailed
				_, err = client.Update(context.TODO(), pod, metav1.UpdateOptions{})
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
