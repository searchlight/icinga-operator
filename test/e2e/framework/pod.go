package framework

import (
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	. "github.com/onsi/gomega"
	apps "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	TestSourceDataVolumeName = "source-data"
	TestSourceDataMountPath  = "/source/data"
)

func (f *Invocation) Pod() *apiv1.Pod {
	return &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("pod"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: getPodSpec(),
	}
}

func (f *Framework) CreatePod(obj *apiv1.Pod) (*apiv1.Pod, error) {
	return f.kubeClient.CoreV1().Pods(obj.Namespace).Create(obj)
}

func (f *Framework) DeletePod(meta metav1.ObjectMeta) error {
	return f.kubeClient.CoreV1().Pods(meta.Namespace).Delete(meta.Name, deleteInForeground())
}

func (f *Framework) EventuallyPodRunning(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() *apiv1.PodList {
			obj, err := f.kubeClient.CoreV1().Pods(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			podList, err := f.GetPodList(obj)
			Expect(err).NotTo(HaveOccurred())
			return podList
		},
		time.Minute*5,
		time.Second*5,
	)
}

func (f *Invocation) PodTemplate() apiv1.PodTemplateSpec {
	return apiv1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: getPodSpec(),
	}
}

func getPodSpec() apiv1.PodSpec {
	return apiv1.PodSpec{
		Containers: []apiv1.Container{
			{
				Name:            "busybox",
				Image:           "busybox",
				ImagePullPolicy: apiv1.PullIfNotPresent,
				Command: []string{
					"sleep",
					"1d",
				},
				VolumeMounts: []apiv1.VolumeMount{
					{
						Name:      TestSourceDataVolumeName,
						MountPath: TestSourceDataMountPath,
					},
				},
			},
		},
		Volumes: []apiv1.Volume{
			{
				Name: TestSourceDataVolumeName,
				VolumeSource: apiv1.VolumeSource{
					EmptyDir: &apiv1.EmptyDirVolumeSource{},
				},
			},
		},
	}
}

func (f *Framework) GetPodList(actual interface{}) (*apiv1.PodList, error) {
	switch obj := actual.(type) {
	case *apiv1.Pod:
		return f.listPods(obj.Namespace, obj.Labels)
	case *extensions.ReplicaSet:
		return f.listPods(obj.Namespace, obj.Spec.Selector.MatchLabels)
	case *extensions.Deployment:
		return f.listPods(obj.Namespace, obj.Spec.Selector.MatchLabels)
	case *apps.Deployment:
		return f.listPods(obj.Namespace, obj.Spec.Selector.MatchLabels)
	case *apps.StatefulSet:
		return f.listPods(obj.Namespace, obj.Spec.Selector.MatchLabels)
	default:
		return nil, fmt.Errorf("Unknown object type")
	}
}

func (f *Framework) listPods(namespace string, label map[string]string) (*apiv1.PodList, error) {
	return f.kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(label).String(),
	})
}
