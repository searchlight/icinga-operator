package mini

import (
	"fmt"
	"os"
	"time"

	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/apps"
)

func CreateStatefulSet(watcher *app.Watcher, namespace string) *apps.StatefulSet {
	// Create Service
	service := CreateService(watcher, namespace, nil)

	statefulSet := &apps.StatefulSet{
		ObjectMeta: kapi.ObjectMeta{
			Namespace: namespace,
		},
		Spec: apps.StatefulSetSpec{
			ServiceName: service.Name,
		},
	}

	if err := testing.CreateKubernetesObject(watcher.Client, statefulSet); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		nStatefulSet, exists, err := watcher.Storage.StatefulSetStore.Get(statefulSet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !exists {
			fmt.Println("StatefulSet not found")
			os.Exit(1)
		}

		if nStatefulSet.(*apps.StatefulSet).Status.Replicas == statefulSet.Spec.Replicas {
			return nStatefulSet.(*apps.StatefulSet)
		}

		if check > 6 {
			fmt.Println("Fail to create StatefulSet")
			os.Exit(1)
		}
		check++
	}
}

func DeleteStatefulSet(watcher *app.Watcher, statefulSet *apps.StatefulSet) {
	// Update StatefulSet
	statefulSet.Spec.Replicas = 0
	if _, err := watcher.Client.Apps().StatefulSets(statefulSet.Namespace).Update(statefulSet); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	labelSelector, err := util.GetLabels(watcher.Client, statefulSet.Namespace, host.TypeStatefulSet, statefulSet.Name)
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
			fmt.Println("Fail to delete StatefulSet Pods")
			os.Exit(1)
		}
		check++
	}

	// Delete StatefulSet
	if err := watcher.Client.Apps().StatefulSets(statefulSet.Namespace).Delete(statefulSet.Name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
