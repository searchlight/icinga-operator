package mini

import (
	"errors"
	"time"

	"github.com/appscode/go/types"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func CreateReplicationController(w *controller.Controller, namespace string) (*apiv1.ReplicationController, error) {
	replicationController := &apiv1.ReplicationController{}
	replicationController.Namespace = namespace
	if err := CreateKubernetesObject(w.KubeClient, replicationController); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		nReplicationController, err := w.Storage.RcStore.ReplicationControllers(replicationController.Namespace).Get(replicationController.Name)
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

func DeleteReplicationController(w *controller.Controller, replicationController *apiv1.ReplicationController) error {
	replicationController, err := w.KubeClient.CoreV1().ReplicationControllers(replicationController.Namespace).Get(replicationController.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// Update ReplicationController
	replicationController.Spec.Replicas = types.Int32P(0)
	if _, err := w.KubeClient.CoreV1().ReplicationControllers(replicationController.Namespace).Update(replicationController); err != nil {
		return err
	}

	labelSelector, err := icinga.GetLabels(w.KubeClient, replicationController.Namespace, icinga.TypeReplicationcontrollers, replicationController.Name)
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
			return errors.New("Fail to delete ReplicationController Pods")
		}
		check++
	}

	// Delete ReplicationController
	if err := w.KubeClient.CoreV1().ReplicationControllers(replicationController.Namespace).Delete(replicationController.Name, nil); err != nil {
		return err
	}
	return nil
}
