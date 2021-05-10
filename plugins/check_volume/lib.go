/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package check_volume

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	"github.com/spf13/cobra"
	"gomodules.xyz/envconfig"
	"gomodules.xyz/x/flags"
	"gomodules.xyz/x/net/httpclient"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"kmodules.xyz/client-go/tools/clientcmd"
)

const (
	awsElasticBlockStorePluginName = "kubernetes.io~aws-ebs"
	azureDataDiskPluginName        = "kubernetes.io~azure-disk"
	azureFilePluginName            = "kubernetes.io~azure-file"
	cephfsPluginName               = "kubernetes.io~cephfs"
	cinderVolumePluginName         = "kubernetes.io~cinder"
	configMapPluginName            = "kubernetes.io~configmap"
	downwardAPIPluginName          = "kubernetes.io~downward-api"
	emptyDirPluginName             = "kubernetes.io~empty-dir"
	fcPluginName                   = "kubernetes.io~fc"
	flockerPluginName              = "kubernetes.io~flocker"
	gcePersistentDiskPluginName    = "kubernetes.io~gce-pd"
	gitRepoPluginName              = "kubernetes.io~git-repo"
	glusterfsPluginName            = "kubernetes.io~glusterfs"
	hostPathPluginName             = "kubernetes.io~host-path"
	iscsiPluginName                = "kubernetes.io~iscsi"
	nfsPluginName                  = "kubernetes.io~nfs"
	quobytePluginName              = "kubernetes.io~quobyte"
	rbdPluginName                  = "kubernetes.io~rbd"
	secretPluginName               = "kubernetes.io~secret"
	vsphereVolumePluginName        = "kubernetes.io~vsphere-volume"
)

func getVolumePluginName(volumeSource *core.VolumeSource) string {
	if volumeSource.AWSElasticBlockStore != nil {
		return awsElasticBlockStorePluginName
	} else if volumeSource.AzureDisk != nil {
		return azureDataDiskPluginName
	} else if volumeSource.AzureFile != nil {
		return azureFilePluginName
	} else if volumeSource.CephFS != nil {
		return cephfsPluginName
	} else if volumeSource.Cinder != nil {
		return cinderVolumePluginName
	} else if volumeSource.ConfigMap != nil {
		return configMapPluginName
	} else if volumeSource.DownwardAPI != nil {
		return downwardAPIPluginName
	} else if volumeSource.EmptyDir != nil {
		return emptyDirPluginName
	} else if volumeSource.FC != nil {
		return fcPluginName
	} else if volumeSource.Flocker != nil {
		return flockerPluginName
	} else if volumeSource.GCEPersistentDisk != nil {
		return gcePersistentDiskPluginName
	} else if volumeSource.GitRepo != nil {
		return gitRepoPluginName
	} else if volumeSource.Glusterfs != nil {
		return glusterfsPluginName
	} else if volumeSource.HostPath != nil {
		return hostPathPluginName
	} else if volumeSource.ISCSI != nil {
		return iscsiPluginName
	} else if volumeSource.NFS != nil {
		return nfsPluginName
	} else if volumeSource.Quobyte != nil {
		return quobytePluginName
	} else if volumeSource.RBD != nil {
		return rbdPluginName
	} else if volumeSource.Secret != nil {
		return secretPluginName
	} else if volumeSource.VsphereVolume != nil {
		return vsphereVolumePluginName
	}
	return ""
}

func getPersistentVolumePluginName(volumeSource *core.PersistentVolumeSource) string {
	if volumeSource.AWSElasticBlockStore != nil {
		return awsElasticBlockStorePluginName
	} else if volumeSource.AzureDisk != nil {
		return azureDataDiskPluginName
	} else if volumeSource.AzureFile != nil {
		return azureFilePluginName
	} else if volumeSource.CephFS != nil {
		return cephfsPluginName
	} else if volumeSource.Cinder != nil {
		return cinderVolumePluginName
	} else if volumeSource.FC != nil {
		return fcPluginName
	} else if volumeSource.Flocker != nil {
		return flockerPluginName
	} else if volumeSource.GCEPersistentDisk != nil {
		return gcePersistentDiskPluginName
	} else if volumeSource.Glusterfs != nil {
		return glusterfsPluginName
	} else if volumeSource.HostPath != nil {
		return hostPathPluginName
	} else if volumeSource.ISCSI != nil {
		return iscsiPluginName
	} else if volumeSource.NFS != nil {
		return nfsPluginName
	} else if volumeSource.Quobyte != nil {
		return quobytePluginName
	} else if volumeSource.RBD != nil {
		return rbdPluginName
	} else if volumeSource.VsphereVolume != nil {
		return vsphereVolumePluginName
	}
	return ""
}

type plugin struct {
	client  kubernetes.Interface
	options options
}

var _ plugins.PluginInterface = &plugin{}

func newPluginFromConfig(opts options) (*plugin, error) {
	client, err := clientcmd.ClientFromContext(opts.kubeconfigPath, opts.contextName)
	if err != nil {
		return nil, err
	}
	return &plugin{client, opts}, nil
}

type options struct {
	kubeconfigPath string
	contextName    string
	// options
	nodeStat   bool
	secretName string
	volumeName string
	mountPoint string
	warning    float64
	critical   float64
	// IcingaHost
	host *icinga.IcingaHost
}

func (o *options) complete(cmd *cobra.Command) error {
	hostname, err := cmd.Flags().GetString(plugins.FlagHost)
	if err != nil {
		return err
	}
	o.host, err = icinga.ParseHost(hostname)
	if err != nil {
		return errors.New("invalid icinga host.name")
	}

	if o.host.Type == icinga.TypeNode {
		o.nodeStat = true
	}

	o.kubeconfigPath, err = cmd.Flags().GetString(plugins.FlagKubeConfig)
	if err != nil {
		return err
	}
	o.contextName, err = cmd.Flags().GetString(plugins.FlagKubeConfigContext)
	if err != nil {
		return err
	}
	return nil
}

func (o *options) validate() error {
	if o.nodeStat {
		if o.host.Type != icinga.TypeNode {
			return errors.New("invalid icinga host type")
		}
	} else {
		if o.host.Type != icinga.TypePod {
			return errors.New("invalid icinga host type")
		}
	}
	return nil
}

type usageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

type AuthInfo struct {
	Port       int    `envconfig:"PORT" default:"56977"`
	Username   string `envconfig:"USERNAME"`
	Password   string `envconfig:"PASSWORD"`
	Token      string `envconfig:"TOKEN"`
	CACertData string `envconfig:"CA_CERT_DATA"`
}

func (p *plugin) getAuthInfo() (*AuthInfo, error) {
	opts := p.options
	host := opts.host

	if opts.secretName == "" {
		return &AuthInfo{Port: 56977}, nil
	}

	secret, err := p.client.CoreV1().Secrets(host.AlertNamespace).Get(context.TODO(), opts.secretName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	var au AuthInfo
	err = envconfig.Load("hostfacts", &au, func(key string) (string, bool) {
		v, ok := secret.Data[key]
		if !ok {
			return "", false
		}
		return string(v), true
	})
	if err != nil {
		return nil, err
	}
	return &au, nil
}

func (p *plugin) getStats(au *AuthInfo, hostIP, path string) (*usageStat, error) {
	hc := httpclient.Default().
		WithBasicAuth(au.Username, au.Password).
		WithBearerToken(au.Token)
	if au.CACertData != "" {
		hc = hc.WithBaseURL(fmt.Sprintf("https://%s:%d/du?p=%s", hostIP, au.Port, path)).
			WithTLSConfig([]byte(au.CACertData))
	} else {
		hc = hc.WithBaseURL(fmt.Sprintf("http://%s:%d/du?p=%s", hostIP, au.Port, path))
	}
	usages := make([]*usageStat, 1)
	_, err := hc.Call(http.MethodGet, "", nil, &usages, true)
	if err != nil {
		return nil, err
	}
	return usages[0], nil
}

func checkResult(field string, warning, critical, result float64) (icinga.State, interface{}) {
	if result >= critical {
		return icinga.Critical, fmt.Sprintf("%v used more than %v%%", field, critical)
	}
	if result >= warning {
		return icinga.Warning, fmt.Sprintf("%v used more than %v%%", field, warning)
	}
	return icinga.OK, "(Disk & Inodes)"
}

func (p *plugin) checkVolume(ip, path string) (icinga.State, interface{}) {
	authInfo, err := p.getAuthInfo()
	if err != nil {
		return icinga.Unknown, err
	}
	usage, err := p.getStats(authInfo, ip, path)
	if err != nil {
		return icinga.Unknown, err
	}

	warning := p.options.warning
	critical := p.options.critical
	state, message := checkResult("Disk", warning, critical, usage.UsedPercent)
	if state != icinga.OK {
		return state, message
	}
	state, message = checkResult("Inodes", warning, critical, usage.InodesUsedPercent)
	return state, message
}

func (p *plugin) checkNodeVolume() (icinga.State, interface{}) {
	node, err := p.client.CoreV1().Nodes().Get(context.TODO(), p.options.host.ObjectName, metav1.GetOptions{})
	if err != nil {
		return icinga.Unknown, err
	}

	if node == nil {
		return icinga.Unknown, "Node not found"
	}

	hostIP := ""
	for _, address := range node.Status.Addresses {
		if address.Type == core.NodeInternalIP {
			hostIP = address.Address
		}
	}

	if hostIP == "" {
		return icinga.Unknown, "Node InternalIP not found"
	}
	return p.checkVolume(hostIP, p.options.mountPoint)
}

func (p *plugin) checkPodVolume() (icinga.State, interface{}) {
	opts := p.options
	host := opts.host

	pod, err := p.client.CoreV1().Pods(host.AlertNamespace).Get(context.TODO(), host.ObjectName, metav1.GetOptions{})
	if err != nil {
		return icinga.Unknown, err
	}

	for _, volume := range pod.Spec.Volumes {
		if volume.Name == opts.volumeName {
			if volume.PersistentVolumeClaim != nil {
				claim, err := p.client.CoreV1().PersistentVolumeClaims(host.AlertNamespace).Get(context.TODO(), volume.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
				if err != nil {
					return icinga.Unknown, err
				}
				volume, err := p.client.CoreV1().PersistentVolumes().Get(context.TODO(), claim.Spec.VolumeName, metav1.GetOptions{})
				if err != nil {
					return icinga.Unknown, err
				}
				volumePluginName := getPersistentVolumePluginName(&volume.Spec.PersistentVolumeSource)
				if volumePluginName == hostPathPluginName {
					if claim.Spec.StorageClassName != nil {
						class, err := p.client.StorageV1beta1().StorageClasses().Get(context.TODO(), *claim.Spec.StorageClassName, metav1.GetOptions{})
						if err != nil {
							return icinga.Unknown, err
						}
						if class.Provisioner == "k8s.io/minikube-hostpath" {
							path := fmt.Sprintf("/tmp/hostpath-provisioner/%s", volume.Name)
							return p.checkVolume(pod.Status.HostIP, path)
						}
					}
				}
				path := fmt.Sprintf("/var/lib/kubelet/pods/%v/volumes/%v/%v", pod.UID, volumePluginName, volume.Name)
				return p.checkVolume(pod.Status.HostIP, path)
			} else {
				path := fmt.Sprintf("/var/lib/kubelet/pods/%v/volumes/%v/%v", pod.UID, getVolumePluginName(&volume.VolumeSource), volume.Name)
				return p.checkVolume(pod.Status.HostIP, path)
			}
			break
		}
	}
	return icinga.Unknown, errors.New("Invalid volume source")
}

func (p *plugin) Check() (icinga.State, interface{}) {

	if p.options.nodeStat {
		return p.checkNodeVolume()
	} else {
		return p.checkPodVolume()
	}
}

const (
	flagVolumeName = "volumeName"
	flagMountPoint = "mountPoint"
)

func NewCmd() *cobra.Command {
	var opts options

	c := &cobra.Command{
		Use:   "check_volume",
		Short: "Check kubernetes volume",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, plugins.FlagHost)
			flags.EnsureAlterableFlags(cmd, flagMountPoint, flagVolumeName)

			if err := opts.complete(cmd); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			if err := opts.validate(); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			plugin, err := newPluginFromConfig(opts)
			if err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			icinga.Output(plugin.Check())
		},
	}

	c.Flags().StringP(plugins.FlagHost, "H", "", "Icinga host name")
	c.Flags().StringVarP(&opts.secretName, "secretName", "s", "", `Kubernetes secret name`)
	c.Flags().StringVarP(&opts.volumeName, flagVolumeName, "N", "", "Volume name")
	c.Flags().StringVarP(&opts.mountPoint, flagMountPoint, "M", "", "Mount point")
	c.Flags().Float64VarP(&opts.warning, "warning", "w", 80.0, "Warning level value (usage percentage)")
	c.Flags().Float64VarP(&opts.critical, "critical", "c", 95.0, "Critical level value (usage percentage)")
	return c
}
