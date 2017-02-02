package fake

import (
	"fmt"
	"os"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/k8s"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
)

func getFakePod(namespace string) *kapi.Pod {
	pod := &kapi.Pod{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: kapi.ObjectMeta{
			Namespace: namespace,
		},
		Spec: kapi.PodSpec{
			Containers: []kapi.Container{
				kapi.Container{
					Name:  "test",
					Image: "fake/image",
				},
			},
		},
		Status: kapi.PodStatus{
			PodIP: "127.0.0.1",
		},
	}
	return pod
}

func createFakePodsUnderObjects(kubeClient *k8s.KubeClient, name, namespace string, labels map[string]string, replica int32) {
	pod := getFakePod(namespace)
	pod.Labels = labels

	for i := int32(0); i < replica; i++ {
		pod.Name = rand.WithUniqSuffix(name)
		if _, err := kubeClient.Client.Core().Pods(pod.Namespace).Create(pod); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
