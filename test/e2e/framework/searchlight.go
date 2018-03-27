package framework

import (
	"github.com/appscode/go/types"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (f *Invocation) DeploymentSearchlight() *apps.Deployment {
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

func (f *Invocation) ServiceSearchlight() *core.Service {
	return &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      f.name,
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: core.ServiceSpec{
			Selector: map[string]string{
				"app": "searchlight",
			},
			Type: core.ServiceTypeLoadBalancer,
			Ports: []core.ServicePort{
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

func (f *Invocation) getSearchlightPodTemplate() core.PodTemplateSpec {
	return core.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:            "icinga",
					Image:           "appscode/icinga:6.0.0-rc.0-k8s",
					ImagePullPolicy: core.PullIfNotPresent,
					Ports: []core.ContainerPort{
						{
							ContainerPort: 5665,
							Name:          "icinga",
						},
						{
							ContainerPort: 60006,
							Name:          "ui",
						},
					},
					LivenessProbe: &core.Probe{
						Handler: core.Handler{
							HTTPGet: &core.HTTPGetAction{
								Scheme: core.URISchemeHTTPS,
								Port:   intstr.FromInt(5665),
								Path:   "/v1/status",
								HTTPHeaders: []core.HTTPHeader{
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
					VolumeMounts: []core.VolumeMount{
						{
							Name:      "data",
							MountPath: "/srv",
						},
					},
				},
				{
					Name:            "ido",
					Image:           "appscode/postgres:9.5-alpine",
					ImagePullPolicy: core.PullIfNotPresent,
					Env: []core.EnvVar{
						{
							Name:  "PGDATA",
							Value: "/var/lib/postgresql/data/pgdata",
						},
					},
					Ports: []core.ContainerPort{
						{
							ContainerPort: 5432,
							Name:          "ido",
						},
					},
					VolumeMounts: []core.VolumeMount{
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
					VolumeMounts: []core.VolumeMount{
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
			Volumes: []core.Volume{
				{
					Name: "data",
					VolumeSource: core.VolumeSource{
						EmptyDir: &core.EmptyDirVolumeSource{},
					},
				},
				{
					Name: "icingaconfig",
					VolumeSource: core.VolumeSource{
						GitRepo: &core.GitRepoVolumeSource{
							Repository: "https://github.com/appscode/icinga-testconfig.git",
							Directory:  ".",
						},
					},
				},
			},
		},
	}
}
