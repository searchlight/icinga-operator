package check_volume

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/appscode/go/flags"
	"github.com/appscode/go/net/httpclient"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
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

func getVolumePluginName(volumeSource *kapi.VolumeSource) string {
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

func getPersistentVolumePluginName(volumeSource *kapi.PersistentVolumeSource) string {
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

const (
	hostFactPort = 56977
)

type request struct {
	host      string
	name      string
	warning   float64
	critical  float64
	node_stat bool
	secret    string
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

type authInfo struct {
	ca        []byte
	key       []byte
	crt       []byte
	authToken string
	username  string
	password  string
}

const (
	ca        = "ca.crt"
	key       = "hostfacts.key"
	crt       = "hostfacts.crt"
	authToken = "auth_token"
	username  = "username"
	password  = "password"
)

func getHostfactsSecretData(kubeClient *k8s.KubeClient, secretName string) *authInfo {
	if secretName == "" {
		return nil
	}

	parts := strings.Split(secretName, ".")
	name := parts[0]
	namespace := "default"
	if len(parts) > 1 {
		namespace = parts[1]
	}

	secret, err := kubeClient.Client.Core().Secrets(namespace).Get(name)
	if err != nil {
		return nil
	}

	authData := &authInfo{
		ca:        secret.Data[ca],
		key:       secret.Data[key],
		crt:       secret.Data[crt],
		authToken: string(secret.Data[authToken]),
		username:  string(secret.Data[username]),
		password:  string(secret.Data[password]),
	}

	return authData
}

func getUsage(authInfo *authInfo, hostIP, path string) (*usageStat, error) {
	scheme := "http"
	httpClient := httpclient.Default()
	if authInfo != nil && authInfo.ca != nil {
		scheme = "https"
		httpClient.WithBasicAuth(authInfo.username, authInfo.password).
			WithBearerToken(authInfo.authToken).
			WithTLSConfig(authInfo.ca, authInfo.crt, authInfo.key)
	}

	urlStr := fmt.Sprintf("%v://%v:%v/du?p=%v", scheme, hostIP, hostFactPort, path)
	usages := make([]*usageStat, 1)
	_, err := httpClient.Call(http.MethodGet, urlStr, nil, &usages, true)
	if err != nil {
		return nil, err
	}

	return usages[0], nil
}

func checkResult(field string, warning, critical, result float64) (util.IcingaState, interface{}) {
	if result >= critical {
		return util.Critical, fmt.Sprintf("%v used more than %v%%", field, critical)
	}
	if result >= warning {
		return util.Warning, fmt.Sprintf("%v used more than %v%%", field, warning)
	}
	return util.Ok, "(Disk & Inodes)"
}

func checkDiskStat(kubeClient *k8s.KubeClient, req *request, nodeIP, path string) (util.IcingaState, interface{}) {
	authInfo := getHostfactsSecretData(kubeClient, req.secret)

	usage, err := getUsage(authInfo, nodeIP, path)
	if err != nil {
		return util.Unknown, err
	}

	warning := req.warning
	critical := req.critical
	state, message := checkResult("Disk", warning, critical, usage.UsedPercent)
	if state != util.Ok {
		return state, message
	}
	state, message = checkResult("Inodes", warning, critical, usage.InodesUsedPercent)
	return state, message
}

func checkNodeDiskStat(req *request) (util.IcingaState, interface{}) {
	host := req.host
	parts := strings.Split(host, "@")
	if len(parts) != 2 {
		return util.Unknown, "Invalid icinga host.name"
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return util.Unknown, err
	}

	node_name := parts[0]
	node, err := kubeClient.Client.Core().Nodes().Get(node_name)
	if err != nil {
		return util.Unknown, err
	}

	if node == nil {
		return util.Unknown, "Node not found"
	}

	hostIP := ""
	for _, address := range node.Status.Addresses {
		if address.Type == kapi.NodeInternalIP {
			hostIP = address.Address
		}
	}

	if hostIP == "" {
		return util.Unknown, "Node InternalIP not found"
	}
	return checkDiskStat(kubeClient, req, hostIP, "/")
}

func checkPodVolumeStat(req *request) (util.IcingaState, interface{}) {
	host := req.host
	name := req.name
	parts := strings.Split(host, "@")
	if len(parts) != 2 {
		return util.Unknown, "Invalid icinga host.name"
	}

	kubeClient, err := k8s.NewClient()
	if err != nil {
		return util.Unknown, err
	}

	pod_name := parts[0]
	namespace := parts[1]
	pod, err := kubeClient.Client.Core().Pods(namespace).Get(pod_name)
	if err != nil {
		return util.Unknown, err
	}

	var volumeSourcePluginName = ""
	var volumeSourceName = ""
	for _, volume := range pod.Spec.Volumes {
		if volume.Name == name {
			if volume.PersistentVolumeClaim != nil {
				claim, err := kubeClient.Client.Core().
					PersistentVolumeClaims(namespace).Get(volume.PersistentVolumeClaim.ClaimName)
				if err != nil {
					return util.Unknown, err

				}
				volume, err := kubeClient.Client.Core().PersistentVolumes().Get(claim.Spec.VolumeName)
				if err != nil {
					return util.Unknown, err
				}
				volumeSourcePluginName = getPersistentVolumePluginName(&volume.Spec.PersistentVolumeSource)
				volumeSourceName = volume.Name

			} else {
				volumeSourcePluginName = getVolumePluginName(&volume.VolumeSource)
				volumeSourceName = volume.Name
			}
			break
		}
	}

	if volumeSourcePluginName == "" {
		return util.Unknown, errors.New("Invalid volume source")
	}

	path := fmt.Sprintf("/var/lib/kubelet/pods/%v/volumes/%v/%v", pod.UID, volumeSourcePluginName, volumeSourceName)
	return checkDiskStat(kubeClient, req, pod.Status.HostIP, path)
}

func NewCmd() *cobra.Command {
	var req request

	c := &cobra.Command{
		Use:     "check_volume",
		Short:   "Check kubernetes volume",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host")
			if req.node_stat {
				util.Output(checkNodeDiskStat(&req))
			} else {
				flags.EnsureRequiredFlags(cmd, "name")
				util.Output(checkPodVolumeStat(&req))
			}
		},
	}

	c.Flags().BoolVar(&req.node_stat, "node_stat", false, "Checking Node disk size")
	c.Flags().StringVarP(&req.secret, "secret", "s", "", `Kubernetes secret name`)
	c.Flags().StringVarP(&req.host, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.name, "name", "N", "", "Volume name")
	c.Flags().Float64VarP(&req.warning, "warning", "w", 75.0, "Warning level value (usage percentage)")
	c.Flags().Float64VarP(&req.critical, "critical", "c", 90.0, "Critical level value (usage percentage)")
	return c
}
