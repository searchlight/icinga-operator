package check_component_status

import (
	"github.com/appscode/searchlight/pkg/icinga"
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
			client.Delete(component.Name, &metav1.DeleteOptions{})
		}
	})

	Describe("Check component status", func() {
		Context("with name", func() {
			It("healthy", func() {
				_, err := client.Create(component)
				Expect(err).ShouldNot(HaveOccurred())

				component.Conditions = []core.ComponentCondition{
					{
						Type:   core.ComponentHealthy,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(component)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
			It("unhealthy", func() {
				_, err := client.Create(component)
				Expect(err).ShouldNot(HaveOccurred())

				component.Conditions = []core.ComponentCondition{
					{
						Type:   core.ComponentHealthy,
						Status: core.ConditionFalse,
					},
				}
				_, err = client.Update(component)
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
					client.Delete(component2.Name, &metav1.DeleteOptions{})
				}
			})
			It("healthy", func() {
				_, err := client.Create(component)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = client.Create(component2)
				Expect(err).ShouldNot(HaveOccurred())

				component.Conditions = []core.ComponentCondition{
					{
						Type:   core.ComponentHealthy,
						Status: core.ConditionTrue,
					},
				}
				_, err = client.Update(component)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})
			It("unhealthy", func() {
				_, err := client.Create(component)
				Expect(err).ShouldNot(HaveOccurred())

				_, err = client.Create(component2)
				Expect(err).ShouldNot(HaveOccurred())

				component.Conditions = []core.ComponentCondition{
					{
						Type:   core.ComponentHealthy,
						Status: core.ConditionFalse,
					},
				}
				_, err = client.Update(component)
				Expect(err).ShouldNot(HaveOccurred())

				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Critical))
			})
		})
	})
})
