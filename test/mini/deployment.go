package mini

import (
	"errors"
	"time"

	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/util"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

func CreateDeployment(watcher *app.Watcher, namespace string) (*extensions.Deployment, error) {
	deployment := &extensions.Deployment{}
	deployment.Namespace = namespace
	if err := testing.CreateKubernetesObject(watcher.Client, deployment); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		nDeployment, err := watcher.Storage.DeploymentStore.Deployments(deployment.Namespace).Get(deployment.Name)
		if err != nil {
			return nil, err
		}

		if deployment.Spec.Replicas == nDeployment.Status.AvailableReplicas {
			return nDeployment, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create Deployment")
		}
		check++
	}
}

func DeleteDeployment(watcher *app.Watcher, deployment *extensions.Deployment) error {
	// Update Deployment
	deployment.Spec.Replicas = 0
	if _, err := watcher.Client.Extensions().Deployments(deployment.Namespace).Update(deployment); err != nil {
		return err
	}

	labelSelector, err := util.GetLabels(watcher.Client, deployment.Namespace, host.TypeDeployments, deployment.Name)
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
			return errors.New("Fail to delete Deployment Pods")
		}
		check++
	}

	// Delete Deployment
	if err := watcher.Client.Extensions().Deployments(deployment.Namespace).Delete(deployment.Name, nil); err != nil {
		return err
	}
	return nil
}
