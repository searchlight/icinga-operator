package testing

import (
	"errors"

	"github.com/appscode/go/crypto/rand"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const (
	DefaultNamespace string = "default"
	Replica          int32  = 2
	Image            string = "busybox"
)

func fixNamespace(ns string) string {
	if ns == "" {
		return DefaultNamespace
	}
	return ns
}

func fixServiceSpec(serviceSpec apiv1.ServiceSpec) apiv1.ServiceSpec {
	if serviceSpec.Selector == nil {
		serviceSpec.Selector = map[string]string{
			"object/random": rand.Characters(6),
		}
	}
	if len(serviceSpec.Ports) == 0 {
		serviceSpec.Ports = []apiv1.ServicePort{
			{
				Port: 80,
			},
		}
	}
	return serviceSpec
}

func fixPodSpec(podSpec apiv1.PodSpec) apiv1.PodSpec {
	if len(podSpec.Containers) == 0 {
		podSpec.Containers = []apiv1.Container{
			{
				Name:    rand.WithUniqSuffix("container"),
				Image:   Image,
				Command: []string{"sleep", "3600"},
			},
		}
	}
	return podSpec
}

func fixPodTemplateSpec(template apiv1.PodTemplateSpec) apiv1.PodTemplateSpec {
	if template.Labels == nil {
		template.Labels = map[string]string{
			"object/random": rand.Characters(6),
		}
	}

	template.Spec = fixPodSpec(template.Spec)
	return template
}

func fixPodTemplateSpecPtr(template *apiv1.PodTemplateSpec) *apiv1.PodTemplateSpec {
	if template == nil {
		template = &apiv1.PodTemplateSpec{}
	}

	fixedTemplate := fixPodTemplateSpec(*template)
	return &fixedTemplate
}

// CreateKubernetesObject will create kubernetes objects
// Pass kubernetes clientset.Interface and object pointer (Example: CreateKubernetesObject(client, &extensions.DaemonSet{}))
func CreateKubernetesObject(kubeClient clientset.Interface, kubeObject interface{}) (err error) {
	switch kubeObject.(type) {
	case *apiv1.ReplicationController:
		replicationController := kubeObject.(*apiv1.ReplicationController)
		if replicationController.Name == "" {
			replicationController.Name = rand.WithUniqSuffix("e2e-rc")
		}
		replicationController.Spec.Template = fixPodTemplateSpecPtr(replicationController.Spec.Template)
		replicationController.Spec.Selector = replicationController.Spec.Template.Labels
		if replicationController.Spec.Replicas == 0 {
			replicationController.Spec.Replicas = Replica
		}
		replicationController, err = kubeClient.Core().ReplicationControllers(fixNamespace(replicationController.Namespace)).Create(replicationController)
		return
	case *extensions.DaemonSet:
		daemonSet := kubeObject.(*extensions.DaemonSet)
		if daemonSet.Name == "" {
			daemonSet.Name = rand.WithUniqSuffix("e2e-daemonset")
		}
		daemonSet.Spec.Template = fixPodTemplateSpec(daemonSet.Spec.Template)
		daemonSet.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: daemonSet.Spec.Template.Labels,
		}
		daemonSet, err = kubeClient.Extensions().DaemonSets(fixNamespace(daemonSet.Namespace)).Create(daemonSet)
		return
	case *apps.StatefulSet:
		statefulSet := kubeObject.(*apps.StatefulSet)
		if statefulSet.Name == "" {
			statefulSet.Name = rand.WithUniqSuffix("e2e-statefulset")
		}
		statefulSet.Spec.Template = fixPodTemplateSpec(statefulSet.Spec.Template)
		statefulSet.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: statefulSet.Spec.Template.Labels,
		}
		if statefulSet.Spec.Replicas == 0 {
			statefulSet.Spec.Replicas = Replica
		}
		statefulSet, err = kubeClient.Apps().StatefulSets(fixNamespace(statefulSet.Namespace)).Create(statefulSet)
		return
	case *extensions.ReplicaSet:
		replicaSet := kubeObject.(*extensions.ReplicaSet)
		if replicaSet.Name == "" {
			replicaSet.Name = rand.WithUniqSuffix("e2e-replicaset")
		}
		replicaSet.Spec.Template = fixPodTemplateSpec(replicaSet.Spec.Template)
		replicaSet.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: replicaSet.Spec.Template.Labels,
		}
		if replicaSet.Spec.Replicas == 0 {
			replicaSet.Spec.Replicas = Replica
		}
		replicaSet, err = kubeClient.Extensions().ReplicaSets(fixNamespace(replicaSet.Namespace)).Create(replicaSet)
		return
	case *extensions.Deployment:
		deployment := kubeObject.(*extensions.Deployment)
		if deployment.Name == "" {
			deployment.Name = rand.WithUniqSuffix("e2e-deployment")
		}
		deployment.Spec.Template = fixPodTemplateSpec(deployment.Spec.Template)
		deployment.Spec.Selector = &metav1.LabelSelector{
			MatchLabels: deployment.Spec.Template.Labels,
		}
		if deployment.Spec.Replicas == 0 {
			deployment.Spec.Replicas = Replica
		}
		deployment, err = kubeClient.Extensions().Deployments(fixNamespace(deployment.Namespace)).Create(deployment)
		return
	case *apiv1.Pod:
		pod := kubeObject.(*apiv1.Pod)
		if pod.Name == "" {
			pod.Name = rand.WithUniqSuffix("e2e-pod")
		}
		pod.Spec = fixPodSpec(pod.Spec)
		pod, err = kubeClient.Core().Pods(fixNamespace(pod.Namespace)).Create(pod)
		return
	case *apiv1.Service:
		service := kubeObject.(*apiv1.Service)
		if service.Name == "" {
			service.Name = rand.WithUniqSuffix("e2e-svc")
		}
		service.Spec = fixServiceSpec(service.Spec)
		service, err = kubeClient.Core().Services(fixNamespace(service.Namespace)).Create(service)
		return
	default:
		err = errors.New("Unknown objectType")
	}
	return
}
