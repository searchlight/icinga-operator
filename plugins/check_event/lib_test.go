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

package check_event

import (
	"context"

	"go.searchlight.dev/icinga-operator/pkg/icinga"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gomodules.xyz/x/crypto/rand"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

/*
Event test List method doesn't support field selector
*/
var _ = XDescribe("check_event", func() {
	var pod *core.Pod
	var event *core.Event
	var client corev1.EventInterface
	var podClient corev1.PodInterface
	var reference core.ObjectReference
	var opts options

	BeforeEach(func() {
		pod = &core.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      rand.Characters(6),
				Namespace: "demo",
				UID:       types.UID(rand.Characters(20)),
			},
		}
		podClient = cs.CoreV1().Pods(pod.Namespace)
		client = cs.CoreV1().Events(pod.Namespace)
		opts = options{
			namespace:               pod.Namespace,
			checkIntervalSecs:       60,
			involvedObjectName:      pod.Name,
			involvedObjectNamespace: pod.Namespace,
			involvedObjectKind:      "Pod",
		}
		reference = core.ObjectReference{
			Kind:      "Pod",
			Namespace: pod.Namespace,
			Name:      pod.Name,
		}
		event = &core.Event{
			ObjectMeta: metav1.ObjectMeta{
				Name: rand.Characters(10),
			},
			Reason:        "test",
			Message:       "unit test",
			LastTimestamp: metav1.Now(),
		}
	})

	AfterEach(func() {
		if podClient != nil {
			podClient.Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		}
		if client != nil {
			client.Delete(context.TODO(), event.Name, metav1.DeleteOptions{})
		}
	})

	Describe("Check Events", func() {
		Context("with warning", func() {
			It("should be Warning", func() {
				pod, err := podClient.Create(context.TODO(), pod, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				reference.UID = pod.UID
				event.InvolvedObject = reference
				event.Type = core.EventTypeWarning
				_, err = client.Create(context.TODO(), event, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				opts.involvedObjectUID = string(pod.UID)
				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.Warning))
			})

		})
		Context("without warning", func() {
			It("should be Ok", func() {
				pod, err := podClient.Create(context.TODO(), pod, metav1.CreateOptions{})
				Expect(err).ShouldNot(HaveOccurred())

				reference.UID = pod.UID
				event.InvolvedObject = reference
				event.Type = core.EventTypeNormal
				_, err = client.Create(context.TODO(), event, metav1.CreateOptions{})

				opts.involvedObjectUID = string(pod.UID)
				state, _ := newPlugin(client, opts).Check()
				Expect(state).Should(BeIdenticalTo(icinga.OK))
			})

		})
	})
})
