package mini

import (
	"fmt"
	"os"
	"time"

	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/util"
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

func DeleteReplicaSet(watcher *app.Watcher, replicaSet *extensions.ReplicaSet) {
	// Create ReplicationController
	replicaSet.Spec.Replicas = 0
	if _, err := watcher.Client.Extensions().ReplicaSets(replicaSet.Namespace).Update(replicaSet); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	labelSelector, err := util.GetLabels(watcher.Client, replicaSet.Namespace, host.TypeReplicasets, replicaSet.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		podList, err := watcher.Storage.PodStore.List(labelSelector)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if len(podList) == 0 {
			break
		}

		if check > 6 {
			fmt.Println("Fail to delete ReplicaSet Pods")
			os.Exit(1)
		}
		check++
	}

	// Create ReplicationController
	if err := watcher.Client.Extensions().ReplicaSets(replicaSet.Namespace).Delete(replicaSet.Name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
