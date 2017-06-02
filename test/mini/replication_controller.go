package mini

import (
	"errors"
	"time"

	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/pkg/testing"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
)

func CreateReplicationController(watcher *app.Watcher, namespace string) (*kapi.ReplicationController, error) {
	replicationController := &kapi.ReplicationController{}
	replicationController.Namespace = namespace
	if err := testing.CreateKubernetesObject(watcher.Client, replicationController); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		nReplicationController, err := watcher.Storage.RcStore.ReplicationControllers(replicationController.Namespace).Get(replicationController.Name)
		if err != nil {
			return nil, err
		}
		if nReplicationController.Status.ReadyReplicas == nReplicationController.Status.Replicas {
			break
		}

		if check > 6 {
			return nil, errors.New("Fail to create ReplicationController")
		}
		check++
	}

	return replicationController, nil
}

func DeleteReplicationController(watcher *app.Watcher, replicationController *kapi.ReplicationController) error {
	replicationController, err := watcher.Client.Core().ReplicationControllers(replicationController.Namespace).Get(replicationController.Name)
	if err != nil {
		return err
	}
	// Update ReplicationController
	replicationController.Spec.Replicas = 0
	if _, err := watcher.Client.Core().ReplicationControllers(replicationController.Namespace).Update(replicationController); err != nil {
		return err
	}

	labelSelector, err := util.GetLabels(watcher.Client, replicationController.Namespace, host.TypeReplicationcontrollers, replicationController.Name)
	if err != nil {
		return err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		podList, err := watcher.Storage.PodStore.List(labelSelector)
		if err != nil {
			return err
		}
		if len(podList) == 0 {
			break
		}

		if check > 6 {
			return errors.New("Fail to delete ReplicationController Pods")
		}
		check++
	}

	// Delete ReplicationController
	if err := watcher.Client.Core().ReplicationControllers(replicationController.Namespace).Delete(replicationController.Name, nil); err != nil {
		return err
	}
	return nil
}
