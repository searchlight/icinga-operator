package mini

import (
	"errors"
	"time"

	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/labels"
)

func CreateReplicaSet(watcher *app.Watcher, namespace string) (*extensions.ReplicaSet, error) {
	replicaSet := &extensions.ReplicaSet{}
	replicaSet.Namespace = namespace
	if err := testing.CreateKubernetesObject(watcher.Client, replicaSet); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		nReplicaset, err := watcher.Storage.ReplicaSetStore.ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name)
		if err != nil {
			return nil, err
		}
		if nReplicaset.Status.ReadyReplicas == nReplicaset.Status.Replicas {
			break
		}

		if check > 6 {
			return nil, errors.New("Fail to create ReplicaSet")
		}
		check++
	}

	return replicaSet, nil
}

func GetSinglePod(watcher *app.Watcher, replicaSet *extensions.ReplicaSet) (*kapi.Pod, error) {
	podList, err := watcher.Storage.PodStore.List(labels.Set(replicaSet.Spec.Selector.MatchLabels).AsSelector())
	if err != nil {
		return nil, err
	}
	if len(podList) == 0 {
		return nil, errors.New("Pod Not Fount")
	}

	return podList[0], nil
}

func DeleteReplicaSet(watcher *app.Watcher, replicaSet *extensions.ReplicaSet) error {
	// Update ReplicaSet
	replicaSet.Spec.Replicas = 0
	if _, err := watcher.Client.Extensions().ReplicaSets(replicaSet.Namespace).Update(replicaSet); err != nil {
		return err
	}

	labelSelector, err := util.GetLabels(watcher.Client, replicaSet.Namespace, host.TypeReplicasets, replicaSet.Name)
	if err != nil {
		return err
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		podList, err := watcher.Storage.PodStore.List(labelSelector)
		if err != nil {
			return err
		}
		if len(podList) == 0 {
			break
		}

		if check > 6 {
			return errors.New("Fail to delete ReplicaSet Pods")
		}
		check++
	}

	// Delete ReplicaSet
	if err := watcher.Client.Extensions().ReplicaSets(replicaSet.Namespace).Delete(replicaSet.Name, nil); err != nil {
		return err
	}
	return nil
}
