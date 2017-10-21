package framework

import (
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	. "github.com/onsi/gomega"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	TestSourceDataVolumeName = "source-data"
	TestSourceDataMountPath  = "/source/data"
)

func (f *Invocation) Pod() *core.Pod {
	return &core.Pod{
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

func (f *Framework) CreatePod(obj *core.Pod) (*core.Pod, error) {
	return f.kubeClient.CoreV1().Pods(obj.Namespace).Create(obj)
}

func (f *Framework) DeletePod(meta metav1.ObjectMeta) error {
	return f.kubeClient.CoreV1().Pods(meta.Namespace).Delete(meta.Name, deleteInForeground())
}

func (f *Framework) EventuallyPodRunning(meta metav1.ObjectMeta) GomegaAsyncAssertion {
	return Eventually(
		func() *core.PodList {
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

func (f *Invocation) PodTemplate() core.PodTemplateSpec {
	return core.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: getPodSpec(),
	}
}

func getPodSpec() core.PodSpec {
	return core.PodSpec{
		Containers: []core.Container{
			{
				Name:            "busybox",
				Image:           "busybox",
				ImagePullPolicy: core.PullIfNotPresent,
				Command: []string{
					"sleep",
					"1d",
				},
				VolumeMounts: []core.VolumeMount{
					{
						Name:      TestSourceDataVolumeName,
						MountPath: TestSourceDataMountPath,
					},
				},
			},
		},
		Volumes: []core.Volume{
			{
				Name: TestSourceDataVolumeName,
				VolumeSource: core.VolumeSource{
					EmptyDir: &core.EmptyDirVolumeSource{},
				},
			},
		},
	}
}

func (f *Framework) GetPodList(actual interface{}) (*core.PodList, error) {
	switch obj := actual.(type) {
	case *core.Pod:
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

func (f *Framework) listPods(namespace string, label map[string]string) (*core.PodList, error) {
	return f.kubeClient.CoreV1().Pods(namespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(label).String(),
	})
}
