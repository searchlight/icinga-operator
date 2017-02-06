package pod_status

import (
	"fmt"
	"os"

	config "github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/test/plugin"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
)

func GetStatusCodeForPodStatus(kubeClient *config.KubeClient, hostname string) util.IcingaState {
	objectType, objectName, namespace := plugin.GetKubeObjectInfo(hostname)

	var err error
	if objectType == host.TypePods {
		pod, err := kubeClient.Client.Core().Pods(namespace).Get(objectName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if !(pod.Status.Phase == kapi.PodSucceeded || pod.Status.Phase == kapi.PodRunning) {
			return util.CRITICAL
		}

	} else {
		labelSelector := labels.Everything()
		if objectType != "" {
			labelSelector, err = util.GetLabels(kubeClient, namespace, objectType, objectName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		var podList *kapi.PodList
		podList, err = kubeClient.Client.Core().Pods(namespace).List(kapi.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, pod := range podList.Items {
			if !(pod.Status.Phase == kapi.PodSucceeded || pod.Status.Phase == kapi.PodRunning) {
				return util.CRITICAL
			}
		}
	}
	return util.OK
}
