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

func CreateDaemonSet(watcher *app.Watcher, namespace string) (*extensions.DaemonSet, error) {
	daemonSet := &extensions.DaemonSet{}
	daemonSet.Namespace = namespace
	if err := testing.CreateKubernetesObject(watcher.Client, daemonSet); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		nDaemonSet, exists, err := watcher.Storage.DaemonSetStore.Get(daemonSet)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("DaemonSet not found")
		}

		if nDaemonSet.(*extensions.DaemonSet).Status.DesiredNumberScheduled == nDaemonSet.(*extensions.DaemonSet).Status.CurrentNumberScheduled {
			return nDaemonSet.(*extensions.DaemonSet), nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create DaemonSet")
		}
		check++
	}
}

func DeleteDaemonSet(watcher *app.Watcher, daemonSet *extensions.DaemonSet) error {
	labelSelector, err := util.GetLabels(watcher.Client, daemonSet.Namespace, host.TypeDaemonsets, daemonSet.Name)
	if err != nil {
		return err
	}

	// Delete DaemonSet
	if err := watcher.Client.Extensions().DaemonSets(daemonSet.Namespace).Delete(daemonSet.Name, nil); err != nil {
		return err
	}

	podList, err := watcher.Storage.PodStore.List(labelSelector)
	if err != nil {
		return err
	}

	for _, pod := range podList {
		if err := watcher.Client.Core().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
			return err
		}
	}
	return nil
}
