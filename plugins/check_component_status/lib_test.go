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

package check_component_status

import (
	"context"

	"go.searchlight.dev/icinga-operator/pkg/icinga"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var _ = Describe("check_component_status", func() {
	var component, component2 *core.ComponentStatus
	var client corev1.ComponentStatusInterface
	var opts options

	BeforeEach(func() {
		client = cs.CoreV1().ComponentStatuses()
		component = &core.ComponentStatus{
			ObjectMeta: metav1.ObjectMeta{
				Name:   "component",
				Labels: make(map[string]string),
			},
		}
		opts = options{
			componentName: component.Name,
		}
	})

	AfterEach(func() {
		if client != nil {
			client.Delete(context.TODO(), component.Name, metav1.DeleteOptions{})
		}
	})

	Describe("Check component status", func() {
		Context("with name", func() {
			It("healthy", func() {
				_, err := client.Create(context.TODO(), component, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				component.Conditions = []core.ComponentCondition{
					{
						Type:   core.ComponentHealthy,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(context.TODO(), component, metav1.UpdateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
			It("unhealthy", func() {
				_, err := client.Create(context.TODO(), component, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				component.Conditions = []core.ComponentCondition{
					{
						Type:   core.ComponentHealthy,
						Status: core.ConditionFalse,
					},
				}
				_, err = client.Update(context.TODO(), component, metav1.UpdateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
		Context("with selector", func() {
			JustBeforeEach(func() {
				component.Labels["app/searchlight"] = "ac"
				opts.componentName = ""
				opts.selector = labels.SelectorFromSet(component.Labels).String()

				component2 = &core.ComponentStatus{
					ObjectMeta: metav1.ObjectMeta{
						Name: "component-2",
						Labels: map[string]string{
							"app/searchlight": "ac",
						},
					},
				}
			})
			AfterEach(func() {
				if client != nil {
					client.Delete(context.TODO(), component2.Name, metav1.DeleteOptions{})
				}
			})
			It("healthy", func() {
				_, err := client.Create(context.TODO(), component, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				_, err = client.Create(context.TODO(), component2, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				component.Conditions = []core.ComponentCondition{
					{
						Type:   core.ComponentHealthy,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(context.TODO(), component, metav1.UpdateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
			It("unhealthy", func() {
				_, err := client.Create(context.TODO(), component, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				_, err = client.Create(context.TODO(), component2, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				component.Conditions = []core.ComponentCondition{
					{
						Type:   core.ComponentHealthy,
						Status: core.ConditionFalse,
					},
				}
				_, err = client.Update(context.TODO(), component, metav1.UpdateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
	})
})
