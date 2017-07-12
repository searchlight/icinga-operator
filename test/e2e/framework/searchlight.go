package framework

import (
	"github.com/appscode/go/io"
	"github.com/appscode/go/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"path/filepath"
)

func (f *Invocation) DeploymentSearchlight() *apps.Deployment {
	return &apps.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "searchlight",
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: apps.DeploymentSpec{
			Replicas: types.Int32P(1),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "searchlight",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            "icinga",
							Image:           "appscode/icinga:1.5.9-k8s",
							ImagePullPolicy: apiv1.PullIfNotPresent,
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 5665,
									Name:          "api",
								},
								{
									ContainerPort: 60006,
									Name:          "web",
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/var/pv",
								},
								{
									Name:      "script-volume",
									MountPath: "/var/db-script",
								},
								{
									Name:      "icingaconfig",
									MountPath: "/srv",
								},
							},
						},
						{
							Name:            "ido",
							Image:           "appscode/postgres:9.5-v3-db",
							ImagePullPolicy: apiv1.PullIfNotPresent,
							Args: []string{
								"basic",
								"./setup-db.sh",
							},
							Ports: []apiv1.ContainerPort{
								{
									ContainerPort: 5432,
									Name:          "ido",
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "data-volume",
									MountPath: "/var/pv",
								},
								{
									Name:      "script-volume",
									MountPath: "/var/db-script",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "data-volume",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "script-volume",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{},
							},
						},
						{
							Name: "icingaconfig",
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: "searchlight",
									Items: []apiv1.KeyToPath{
										{
											Key: "ca.crt",
											Path: filepath.Join("pki", "ca.crt"),
										},
										{
											Key: "icinga.crt",
											Path: filepath.Join("pki", "icinga.crt"),
										},
										{
											Key: "icinga.key",
											Path: filepath.Join("pki", "icinga.key"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (f *Invocation) SecretSearchlight(path string) (*apiv1.Secret, error) {
	secret := &apiv1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "searchlight",
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
	}

	secret.Data = make(map[string][]byte)

	ca, err := io.ReadFile(filepath.Join(path, "pki/ca.crt"))
	if err != nil {
		return nil, err
	}
	secret.Data["ca.crt"] = []byte(ca)

	crt, err := io.ReadFile(filepath.Join(path, "pki/icinga.crt"))
	if err != nil {
		return nil, err
	}
	secret.Data["icinga.crt"] = []byte(crt)

	key, err := io.ReadFile(filepath.Join(path, "pki/icinga.key"))
	if err != nil {
		return nil, err
	}
	secret.Data["icinga.key"] = []byte(key)

	ini, err := io.ReadFile(filepath.Join(path, "config.ini"))
	if err != nil {
		return nil, err
	}
	secret.Data["config.ini"] = []byte(ini)

	return secret, err
}
