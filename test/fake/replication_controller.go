package fake

import (
	"fmt"
	"os"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/k8s"
	kapi "k8s.io/kubernetes/pkg/api"
)

func getFakeReplicationController(namespace string) *kapi.ReplicationController {
	rcName := rand.WithUniqSuffix("fake-rc")
	// Fake ReplicationController object
	rc := &kapi.ReplicationController{
		ObjectMeta: kapi.ObjectMeta{
			Name:      rcName,
			Namespace: namespace,
		},
		Spec: kapi.ReplicationControllerSpec{
			Selector: map[string]string{
				"fake-server-rc": rcName,
			},
			Replicas: 3,
			Template: &kapi.PodTemplateSpec{
				ObjectMeta: kapi.ObjectMeta{
					Labels: map[string]string{
						"fake-server-rc": rcName,
					},
				},
				Spec: kapi.PodSpec{
					Containers: []kapi.Container{
						kapi.Container{
							Name:  "test",
							Image: "fake/image",
						},
					},
				},
			},
		},
	}
	return rc
}

func CreateFakeReplicationController(kubeClient *k8s.KubeClient, namespace string) *kapi.ReplicationController {
	// Fake ReplicationController object
	rc := getFakeReplicationController(namespace)

	// Create Fake ReplicationController
	if _, err := kubeClient.Client.Core().ReplicationControllers(rc.Namespace).Create(rc); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create Pods under ReplicationController
	createFakePodsUnderObjects(kubeClient, rc.Name, rc.Namespace, rc.Spec.Selector, rc.Spec.Replicas)
	return rc
}
