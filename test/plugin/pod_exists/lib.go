package pod_exists

import (
	"fmt"
	"os"

	"github.com/appscode/searchlight/cmd/searchlight/app"
	"k8s.io/kubernetes/pkg/labels"
)

func GetPodCount(watcher *app.Watcher, namespace string) int {
	podList, err := watcher.Storage.PodStore.Pods(namespace).List(labels.Everything())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return len(podList)
}
