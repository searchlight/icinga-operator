package mini

import (
	"errors"
	"time"

	"github.com/appscode/searchlight/pkg/controller"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func checkPod(w *controller.Controller, pod *apiv1.Pod) (*apiv1.Pod, error) {
	check := 0
	for {
		time.Sleep(time.Second * 30)
		nPod, err := w.Storage.PodStore.Pods(pod.Namespace).Get(pod.Name)
		if err != nil {
			return nil, err
		}
		if nPod.Status.Phase == apiv1.PodRunning {
			return nPod, nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create Pod")
		}
		check++
	}
}

func CreatePod(w *controller.Controller, namespace string) (*apiv1.Pod, error) {
	pod := &apiv1.Pod{}
	pod.Namespace = namespace
	if err := CreateKubernetesObject(w.KubeClient, pod); err != nil {
		return nil, err
	}

	return checkPod(w, pod)
}

func ReCreatePod(w *controller.Controller, pod *apiv1.Pod) (*apiv1.Pod, error) {
	newPod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pod.Name,
			Namespace: pod.Namespace,
		},
	}
	if err := CreateKubernetesObject(w.KubeClient, newPod); err != nil {
		return nil, err
	}

	return checkPod(w, newPod)
}

func DeletePod(w *controller.Controller, pod *apiv1.Pod) error {
	// Delete Pod
	if err := w.KubeClient.CoreV1().Pods(pod.Namespace).Delete(pod.Name, nil); err != nil {
		return err
	}
	return nil
}
