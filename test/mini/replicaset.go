package mini

import (
	"fmt"
	"os"
	"time"

	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func CreateReplicaSet(watcher *app.Watcher, namespace string) *extensions.ReplicaSet {

	replicaSet := &extensions.ReplicaSet{}
	replicaSet.Namespace = namespace
	if err := testing.CreateKubernetesObject(watcher.Client, replicaSet); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		nReplicaset, err := watcher.Storage.ReplicaSetStore.ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if nReplicaset.Status.ReadyReplicas == nReplicaset.Status.Replicas {
			break
		}

		if check > 6 {
			fmt.Println("Fail to create ReplicaSet")
			os.Exit(1)
		}
		check++
	}

	return replicaSet
}

func DeleteReplicaSet(kubeClient *k8s.KubeClient, replicaset *extensions.ReplicaSet) {
	// Create ReplicationController
	if err := kubeClient.Client.Extensions().ReplicaSets(replicaset.Namespace).Delete(replicaset.Name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
