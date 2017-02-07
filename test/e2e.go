package e2e

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/appscode/go/crypto/rand"
	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/k8s-addons/pkg/testing"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/client"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/apps"
	ext "k8s.io/kubernetes/pkg/apis/extensions"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/labels"
)

const (
	GOHOSTOS         string = "linux"
	GOHOSTARCH       string = "amd64"
	DefaultNamespace string = "default"
)

type dataConfig struct {
	ObjectType   string
	CheckCommand string
	Namespace    string
}

func fixNamespace(ns string) string {
	if ns == "" {
		return DefaultNamespace
	}
	return ns
}

func getHostName(objectType, objectName string, namespace ...string) string {
	object := objectName
	if objectType != "" {
		object = fmt.Sprintf("%s|%s", objectType, objectName)
	}

	if len(namespace) == 1 {
		object = fmt.Sprintf("%s@%s", object, namespace[0])
	} else {
		object = fmt.Sprintf("%s@default", object)
	}
	return object
}

func getClusterCheckData(kubeClient clientset.Interface, checkCommand, namespace string) (name string, count int32, err error) {
	var podList *kapi.PodList
	if podList, err = kubeClient.Core().Pods(fixNamespace(namespace)).List(
		kapi.ListOptions{LabelSelector: labels.Everything()}); err != nil {
		return
	}
	count = int32(len(podList.Items))
	name = getHostName("", checkCommand, namespace)
	return
}

func getKubernetesObjectData(kubeClient clientset.Interface, objectType, namespace string) (name string, count int32, err error) {
	switch objectType {
	case host.TypeReplicationcontrollers:
		replicationController := &kapi.ReplicationController{}
		replicationController.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, replicationController); err != nil {
			return
		}
		name = getHostName(host.TypeReplicationcontrollers, replicationController.Name, replicationController.Namespace)
		count = replicationController.Spec.Replicas
	case host.TypeDaemonsets:
		daemonSet := &ext.DaemonSet{}
		daemonSet.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, daemonSet); err != nil {
			return
		}

		if daemonSet, err = kubeClient.Extensions().
			DaemonSets(daemonSet.Namespace).Get(daemonSet.Name); err != nil {
			return
		}
		name = getHostName(host.TypeDaemonsets, daemonSet.Name, daemonSet.Namespace)
		count = daemonSet.Status.DesiredNumberScheduled
	case host.TypeStatefulSet:
		statefulSet := &apps.StatefulSet{}
		statefulSet.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, statefulSet); err != nil {
			return
		}
		name = getHostName(host.TypeStatefulSet, statefulSet.Name, statefulSet.Namespace)
		count = statefulSet.Spec.Replicas
	case host.TypeReplicasets:
		replicaSet := &ext.ReplicaSet{}
		replicaSet.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, replicaSet); err != nil {
			return
		}
		name = getHostName(host.TypeReplicasets, replicaSet.Name, replicaSet.Namespace)
		count = replicaSet.Spec.Replicas
	case host.TypeDeployments:
		deployment := &ext.Deployment{}
		deployment.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, deployment); err != nil {
			return
		}
		name = getHostName(host.TypeDeployments, deployment.Name, deployment.Namespace)
		count = deployment.Spec.Replicas
	case host.TypePods:
		pod := &kapi.Pod{}
		pod.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, pod); err != nil {
			return
		}
		name = getHostName(host.TypePods, pod.Name, pod.Namespace)

	case host.TypeServices:
		replicaSet := &ext.ReplicaSet{}
		replicaSet.Namespace = namespace
		if err = testing.CreateKubernetesObject(kubeClient, replicaSet); err != nil {
			return
		}

		service := &kapi.Service{
			ObjectMeta: kapi.ObjectMeta{
				Namespace: replicaSet.Namespace,
			},
			Spec: kapi.ServiceSpec{
				Selector: replicaSet.Spec.Selector.MatchLabels,
			},
		}
		if err = testing.CreateKubernetesObject(kubeClient, service); err != nil {
			return
		}
		name = getHostName(host.TypeServices, service.Name, service.Namespace)
		count = replicaSet.Spec.Replicas
	default:
		err = errors.New("Unknown objectType")
	}
	return
}

func getTestData(kubeClient *k8s.KubeClient, dataConfig *dataConfig) (name string, count int32) {
	var err error
	if dataConfig.ObjectType == host.TypeCluster {
		name, count, err = getClusterCheckData(kubeClient.Client, dataConfig.CheckCommand, dataConfig.Namespace)
	} else {
		name, count, err = getKubernetesObjectData(kubeClient.Client, dataConfig.ObjectType, dataConfig.Namespace)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}

func createNewNamespace(kubeClient *k8s.KubeClient, name string) {
	ns := &kapi.Namespace{
		ObjectMeta: kapi.ObjectMeta{
			Name: name,
		},
	}
	_, err := kubeClient.Client.Core().Namespaces().Create(ns)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func deleteNewNamespace(kubeClient *k8s.KubeClient, name string) {
	err := kubeClient.Client.Core().Namespaces().Delete(name, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func createAlertObject(kubeClient *k8s.KubeClient, alert *aci.Alert) (err error) {
	if alert.Name == "" {
		alert.Name = rand.WithUniqSuffix("e2e-alert")
	}

	alert, err = kubeClient.AppscodeExtensionClient.Alert(fixNamespace(alert.Namespace)).Create(alert)
	return
}

func deleteAlertObject(kubeClient *k8s.KubeClient, alert *aci.Alert) (err error) {
	// delete alert
	err = kubeClient.AppscodeExtensionClient.Alert(alert.Namespace).Delete(alert.Name, nil)
	return
}

func runKubeD(context *client.Context) *app.Watcher {
	fmt.Println("-- TestE2E: Waiting for kubed")
	w := &app.Watcher{
		Watcher: acw.Watcher{
			Client:                  context.KubeClient.Client,
			AppsCodeExtensionClient: context.KubeClient.AppscodeExtensionClient,
			SyncPeriod:              time.Minute * 2,
		},
		IcingaClient: context.IcingaClient,
	}

	w.Watcher.Dispatch = w.Dispatch
	go w.Run()
	time.Sleep(time.Second * 10)
	return w
}
