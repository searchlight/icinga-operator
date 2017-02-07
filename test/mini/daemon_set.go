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

func CreateDaemonSet(watcher *app.Watcher, namespace string) *extensions.DaemonSet {
	daemonSet := &extensions.DaemonSet{}
	daemonSet.Namespace = namespace
	if err := testing.CreateKubernetesObject(watcher.Client, daemonSet); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		nDaemonSet, exists, err := watcher.Storage.DaemonSetStore.Get(daemonSet)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !exists {
			fmt.Println("DaemonSet not found")
			os.Exit(1)
		}

		if nDaemonSet.(*extensions.DaemonSet).Status.DesiredNumberScheduled == nDaemonSet.(*extensions.DaemonSet).Status.CurrentNumberScheduled {
			return nDaemonSet.(*extensions.DaemonSet)
		}

		if check > 6 {
			fmt.Println("Fail to create DaemonSet")
			os.Exit(1)
		}
		check++
	}
}

func DeleteDaemonSet(watcher *app.Watcher, daemonSet *extensions.DaemonSet) {
	labelSelector, err := util.GetLabels(watcher.Client, daemonSet.Namespace, host.TypeDaemonsets, daemonSet.Name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Delete DaemonSet
	if err := watcher.Client.Extensions().DaemonSets(daemonSet.Namespace).Delete(daemonSet.Name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	podList, err := watcher.Storage.PodStore.List(labelSelector)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, pod := range podList {
		if err := watcher.Client.Core().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
