package mini

import (
	"errors"
	"time"

	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func CreateDaemonSet(w *controller.Controller, namespace string) (*extensions.DaemonSet, error) {
	daemonSet := &extensions.DaemonSet{}
	daemonSet.Namespace = namespace
	if err := CreateKubernetesObject(w.KubeClient, daemonSet); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		dmn, err := w.Storage.DaemonSetStore.DaemonSets(daemonSet.Namespace).Get(daemonSet.Name)
		if err != nil {
			return nil, err
		}

		if dmn.Status.DesiredNumberScheduled == dmn.Status.CurrentNumberScheduled {
			return dmn, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create DaemonSet")
		}
		check++
	}
}

func DeleteDaemonSet(watcher *controller.Controller, daemonSet *extensions.DaemonSet) error {
	labelSelector, err := icinga.GetLabels(watcher.KubeClient, daemonSet.Namespace, icinga.TypeDaemonsets, daemonSet.Name)
	if err != nil {
		return err
	}

	// Delete DaemonSet
	if err := watcher.KubeClient.ExtensionsV1beta1().DaemonSets(daemonSet.Namespace).Delete(daemonSet.Name, nil); err != nil {
		return err
	}

	podList, err := watcher.Storage.PodStore.List(labelSelector)
	if err != nil {
		return err
	}

	for _, pod := range podList {
		if err := watcher.KubeClient.CoreV1().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
			return err
		}
	}
	return nil
}
