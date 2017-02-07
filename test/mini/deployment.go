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

func CreateDeployment(watcher *app.Watcher, namespace string) *extensions.Deployment {
	deployment := &extensions.Deployment{}
	deployment.Namespace = namespace
	if err := testing.CreateKubernetesObject(watcher.Client, deployment); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		nDeployment, err := watcher.Storage.DeploymentStore.Deployments(deployment.Namespace).Get(deployment.Name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if deployment.Spec.Replicas == nDeployment.Status.AvailableReplicas {
			break
		}

		if check > 6 {
			fmt.Println("Fail to create Deployment")
			os.Exit(1)
		}
		check++
	}
	return deployment
}

func DeleteDeployment(watcher *app.Watcher, deployment *extensions.Deployment) {
	// Create ReplicationController
	deployment.Spec.Replicas = 0
	if _, err := watcher.Client.Extensions().Deployments(deployment.Namespace).Update(deployment); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	labelSelector, err := util.GetLabels(watcher.Client, deployment.Namespace, host.TypeDeployments, deployment.Name)
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
			fmt.Println("Fail to delete Deployment Pods")
			os.Exit(1)
		}
		check++
	}

	// Create ReplicationController
	if err := watcher.Client.Extensions().Deployments(deployment.Namespace).Delete(deployment.Name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
