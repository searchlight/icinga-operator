package mini

import (
	"errors"
	"time"

	"github.com/appscode/go/types"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func CreateDeployment(w *controller.Controller, namespace string) (*extensions.Deployment, error) {
	deployment := &extensions.Deployment{}
	deployment.Namespace = namespace
	if err := CreateKubernetesObject(w.KubeClient, deployment); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		nDeployment, err := w.Storage.DeploymentStore.Deployments(deployment.Namespace).Get(deployment.Name)
		if err != nil {
			return nil, err
		}

		if *deployment.Spec.Replicas == nDeployment.Status.AvailableReplicas {
			return nDeployment, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create Deployment")
		}
		check++
	}
}

func DeleteDeployment(w *controller.Controller, deployment *extensions.Deployment) error {
	deployment, err := w.KubeClient.ExtensionsV1beta1().Deployments(deployment.Namespace).Get(deployment.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// Update Deployment
	deployment.Spec.Replicas = types.Int32P(0)
	if _, err := w.KubeClient.ExtensionsV1beta1().Deployments(deployment.Namespace).Update(deployment); err != nil {
		return err
	}

	labelSelector, err := icinga.GetLabels(w.KubeClient, deployment.Namespace, icinga.TypeDeployments, deployment.Name)
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
			return errors.New("Fail to delete Deployment Pods")
		}
		check++
	}

	// Delete Deployment
	if err := w.KubeClient.ExtensionsV1beta1().Deployments(deployment.Namespace).Delete(deployment.Name, nil); err != nil {
		return err
	}
	return nil
}
