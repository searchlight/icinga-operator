package mini

import (
	"errors"
	"time"

	"github.com/appscode/go/types"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func checkReplicaSet(w *controller.Controller, replicaSet *extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	check := 0
	for {
		time.Sleep(time.Second * 30)
		nReplicaset, err := w.Storage.ReplicaSetStore.ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name)
		if err != nil {
			return nil, err
		}
		if nReplicaset.Status.ReadyReplicas == nReplicaset.Status.Replicas {
			return nReplicaset, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create ReplicaSet")
		}
		check++
	}
}

func CreateReplicaSet(w *controller.Controller, namespace string) (*extensions.ReplicaSet, error) {
	replicaSet := &extensions.ReplicaSet{}
	replicaSet.Namespace = namespace
	if err := CreateKubernetesObject(w.KubeClient, replicaSet); err != nil {
		return nil, err
	}

	return checkReplicaSet(w, replicaSet)
}

func ReCreateReplicaSet(w *controller.Controller, replicaSet *extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	newReplicaSet := &extensions.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      replicaSet.Name,
			Namespace: replicaSet.Namespace,
		},
		Spec: extensions.ReplicaSetSpec{
			Replicas: replicaSet.Spec.Replicas,
		},
	}
	if err := CreateKubernetesObject(w.KubeClient, newReplicaSet); err != nil {
		return nil, err
	}

	return checkReplicaSet(w, newReplicaSet)
}

func GetLastReplica(w *controller.Controller, replicaSet *extensions.ReplicaSet) (*apiv1.Pod, error) {
	podList, err := w.Storage.PodStore.List(labels.Set(replicaSet.Spec.Selector.MatchLabels).AsSelector())
	if err != nil {
		return nil, err
	}
	if len(podList) == 0 {
		return nil, errors.New("Pod Not Fount")
	}

	var lastCreationTime metav1.Time
	var lastPod *apiv1.Pod

	for _, pod := range podList {
		if lastCreationTime.Before(pod.CreationTimestamp) {
			lastCreationTime = pod.CreationTimestamp
			lastPod = pod
		}

	}
	return lastPod, nil
}

func DeleteReplicaSet(w *controller.Controller, replicaSet *extensions.ReplicaSet) error {
	// Update ReplicaSet
	replicaSet, err := w.KubeClient.ExtensionsV1beta1().ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	replicaSet.Spec.Replicas = types.Int32P(0)
	if _, err := w.KubeClient.ExtensionsV1beta1().ReplicaSets(replicaSet.Namespace).Update(replicaSet); err != nil {
		return err
	}

	labelSelector, err := icinga.GetLabels(w.KubeClient, replicaSet.Namespace, icinga.TypeReplicasets, replicaSet.Name)
	if err != nil {
		return err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		podList, err := w.Storage.PodStore.List(labelSelector)
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
	if err := w.KubeClient.ExtensionsV1beta1().ReplicaSets(replicaSet.Namespace).Delete(replicaSet.Name, nil); err != nil {
		return err
	}
	return nil
}

func UpdateReplicaSet(w *controller.Controller, replicaSet *extensions.ReplicaSet) (*extensions.ReplicaSet, error) {
	if _, err := w.KubeClient.ExtensionsV1beta1().ReplicaSets(replicaSet.Namespace).Update(replicaSet); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		nReplicaset, err := w.Storage.ReplicaSetStore.ReplicaSets(replicaSet.Namespace).Get(replicaSet.Name)
		if err != nil {
			return nil, err
		}
		if nReplicaset.Status.ReadyReplicas == nReplicaset.Status.Replicas {
			return nReplicaset, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create ReplicaSet")
		}
		check++
	}
}
