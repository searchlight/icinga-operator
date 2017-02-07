package mini

import (
	"fmt"
	"os"

	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"time"
)

func CreateReplicationController(watcher *app.Watcher, namespace string) *kapi.ReplicationController {
	replicationController := &kapi.ReplicationController{}
	replicationController.Namespace = namespace
	if err := testing.CreateKubernetesObject(watcher.Client, replicationController); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		nReplicationController, err := watcher.Storage.RcStore.ReplicationControllers(replicationController.Namespace).Get(replicationController.Name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if nReplicationController.Status.ReadyReplicas == nReplicationController.Status.Replicas {
			break
		}

		if check > 6 {
			fmt.Println("Fail to create ReplicationController")
			os.Exit(1)
		}
		check++
	}

	return replicationController
}

func DeleteReplicationController(watcher *app.Watcher, replicationController *kapi.ReplicationController) {
	// Update ReplicationController
	replicationController.Spec.Replicas = 0
	if _, err := watcher.Client.Core().ReplicationControllers(replicationController.Namespace).Update(replicationController); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	labelSelector, err := util.GetLabels(watcher.Client, replicationController.Namespace, host.TypeReplicationcontrollers, replicationController.Name)
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
			fmt.Println("Fail to delete ReplicationController Pods")
			os.Exit(1)
		}
		check++
	}

	// Delete ReplicationController
	if err := watcher.Client.Core().ReplicationControllers(replicationController.Namespace).Delete(replicationController.Name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
