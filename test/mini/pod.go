package mini

import (
	"errors"
	"time"

	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/testing"
	kapi "k8s.io/kubernetes/pkg/api"
)

func checkPod(watcher *app.Watcher, pod *kapi.Pod) (*kapi.Pod, error) {
	check := 0
	for {
		time.Sleep(time.Second * 30)
		nPod, err := watcher.Storage.PodStore.Pods(pod.Namespace).Get(pod.Name)
		if err != nil {
			return nil, err
		}
		if nPod.Status.Phase == kapi.PodRunning {
			return nPod, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create Pod")
		}
		check++
	}
}

func CreatePod(watcher *app.Watcher, namespace string) (*kapi.Pod, error) {
	pod := &kapi.Pod{}
	pod.Namespace = namespace
	if err := testing.CreateKubernetesObject(watcher.Client, pod); err != nil {
		return nil, err
	}

	return checkPod(watcher, pod)
}

func ReCreatePod(watcher *app.Watcher, pod *kapi.Pod) (*kapi.Pod, error) {
	newPod := &kapi.Pod{
		ObjectMeta: kapi.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
	}
	if err := testing.CreateKubernetesObject(watcher.Client, newPod); err != nil {
		return nil, err
	}

	return checkPod(watcher, newPod)
}

func DeletePod(watcher *app.Watcher, pod *kapi.Pod) error {
	// Delete Pod
	if err := watcher.Client.Core().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
		return err
	}
	return nil
}
