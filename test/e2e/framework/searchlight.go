package framework

import (
	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (f *Invocation) DeploymentAppSearchlight() *apps.Deployment {
	return &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.name,
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: types.Int32P(1),
			Template: f.getSearchlightPodTemplate(),
		},
	}
}

func (f *Invocation) DeploymentExtensionSearchlight() *extensions.Deployment {
	return &extensions.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.name,
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: extensions.DeploymentSpec{
			Replicas: types.Int32P(1),
			Template: f.getSearchlightPodTemplate(),
		},
	}
}

func (f *Invocation) ServiceSearchlight() *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.name,
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "searchlight",
			},
			Type: apiv1.ServiceTypeLoadBalancer,
			Ports: []apiv1.ServicePort{
				{
					Name:       "icinga",
					Port:       5665,
					TargetPort: intstr.Parse("icinga"),
				},
				{
					Name:       "ui",
					Port:       80,
					TargetPort: intstr.Parse("ui"),
				},
			},
		},
	}
}

func (f *Invocation) getSearchlightPodTemplate() apiv1.PodTemplateSpec {
	return apiv1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:            "icinga",
					Image:           "appscode/icinga:4.0.0-k8s",
					ImagePullPolicy: apiv1.PullIfNotPresent,
					Ports: []apiv1.ContainerPort{
						{
							ContainerPort: 5665,
							Name:          "icinga",
						},
						{
							ContainerPort: 60006,
							Name:          "ui",
						},
					},
					LivenessProbe: &apiv1.Probe{
						Handler: apiv1.Handler{
							HTTPGet: &apiv1.HTTPGetAction{
								Scheme: apiv1.URISchemeHTTPS,
								Port:   intstr.FromInt(5665),
								Path:   "/v1/status",
								HTTPHeaders: []apiv1.HTTPHeader{
									{
										Name:  "Authorization",
										Value: "Basic c3RhdHVzdXNlcjpzdGF0dXNwYXNz",
									},
								},
							},
						},
						InitialDelaySeconds: 300,
						PeriodSeconds:       120,
					},
					VolumeMounts: []apiv1.VolumeMount{
						{
							Name:      "data",
							MountPath: "/srv",
						},
					},
				},
				{
					Name:            "ido",
					Image:           "appscode/postgres:9.5-alpine",
					ImagePullPolicy: apiv1.PullIfNotPresent,
					Env: []apiv1.EnvVar{
						{
							Name:  "PGDATA",
							Value: "/var/lib/postgresql/data/pgdata",
						},
					},
					Ports: []apiv1.ContainerPort{
						{
							ContainerPort: 5432,
							Name:          "ido",
						},
					},
					VolumeMounts: []apiv1.VolumeMount{
						{
							Name:      "data",
							MountPath: "/var/lib/postgresql/data",
						},
					},
				},
				{
					Name:  "busybox",
					Image: "busybox",
					Command: []string{
						"/bin/sh",
						"-c",
						"cp -rf /var/searchlight /srv/searchlight && sleep 1d",
					},
					VolumeMounts: []apiv1.VolumeMount{
						{
							Name:      "data",
							MountPath: "/srv",
						},
						{
							Name:      "icingaconfig",
							MountPath: "/var/",
						},
					},
				},
			},
			Volumes: []apiv1.Volume{
				{
					Name: "data",
					VolumeSource: apiv1.VolumeSource{
						EmptyDir: &apiv1.EmptyDirVolumeSource{},
					},
				},
				{
					Name: "icingaconfig",
					VolumeSource: apiv1.VolumeSource{
						GitRepo: &apiv1.GitRepoVolumeSource{
							Repository: "https://github.com/appscode/icinga-testconfig.git",
							Directory:  ".",
						},
					},
				},
			},
		},
	}
}
