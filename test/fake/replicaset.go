package fake

import (
	"fmt"
	"os"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/k8s"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func getFakeReplicaSet(namespace string) *extensions.ReplicaSet {
	replicasetName := rand.WithUniqSuffix("fake-replicaset")
	labelMap := map[string]string{
		"fake-server-replicaset": replicasetName,
	}
	// Fake ReplicaSet object
	replicaset := &extensions.ReplicaSet{
		TypeMeta: unversioned.TypeMeta{
			Kind:       "ReplicaSet",
			APIVersion: "extensions/v1beta1",
		},
		ObjectMeta: kapi.ObjectMeta{
			Name:      replicasetName,
			Namespace: namespace,
		},
		Spec: extensions.ReplicaSetSpec{
			Selector: unversioned.SetAsLabelSelector(labelMap),
			Replicas: 3,
			Template: kapi.PodTemplateSpec{
				ObjectMeta: kapi.ObjectMeta{
					Labels: labelMap,
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
	return replicaset
}

func CreateFakeReplicaSet(kubeClient *k8s.KubeClient, namespace string) *extensions.ReplicaSet {
	// Fake ReplicaSet object
	replicaset := getFakeReplicaSet(namespace)
	// Create Fake ReplicaSet
	if _, err := kubeClient.Client.Extensions().ReplicaSets(replicaset.Namespace).Create(replicaset); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create Pods under ReplicationController
	createFakePodsUnderObjects(kubeClient, replicaset.Name, replicaset.Namespace, replicaset.Spec.Selector.MatchLabels, replicaset.Spec.Replicas)
	return replicaset
}
