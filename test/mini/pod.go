package mini

import (
	"github.com/appscode/searchlight/cmd/searchlight/app"
	kapi "k8s.io/kubernetes/pkg/api"
)

func DeletePod(watcher *app.Watcher, pod *kapi.Pod) error {
	// Delete Pod
	if err := watcher.Client.Core().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
		return err
	}
	return nil
}
