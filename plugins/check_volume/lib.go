package check_volume

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/appscode/envconfig"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/net/httpclient"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

func getVolumePluginName(volumeSource *apiv1.VolumeSource) string {
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

func getPersistentVolumePluginName(volumeSource *apiv1.PersistentVolumeSource) string {
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

type Request struct {
	masterURL      string
	kubeconfigPath string

	Host       string
	NodeStat   bool
	SecretName string
	VolumeName string
	Warning    float64
	Critical   float64
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

func getAuthInfo(kubeClient kubernetes.Interface, secretName, secretNamespace string) (*AuthInfo, error) {
	if secretName == "" {
		return &AuthInfo{Port: 56977}, nil
	}

	secret, err := kubeClient.CoreV1().Secrets(secretNamespace).Get(secretName, metav1.GetOptions{})
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

func getStats(au *AuthInfo, hostIP, path string) (*usageStat, error) {
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
		return icinga.CRITICAL, fmt.Sprintf("%v used more than %v%%", field, critical)
	}
	if result >= warning {
		return icinga.WARNING, fmt.Sprintf("%v used more than %v%%", field, warning)
	}
	return icinga.OK, "(Disk & Inodes)"
}

func checkVolume(kubeClient kubernetes.Interface, req *Request, namespace, ip, path string) (icinga.State, interface{}) {
	authInfo, err := getAuthInfo(kubeClient, req.SecretName, namespace)
	if err != nil {
		return icinga.UNKNOWN, err
	}
	usage, err := getStats(authInfo, ip, path)
	if err != nil {
		return icinga.UNKNOWN, err
	}

	warning := req.Warning
	critical := req.Critical
	state, message := checkResult("Disk", warning, critical, usage.UsedPercent)
	if state != icinga.OK {
		return state, message
	}
	state, message = checkResult("Inodes", warning, critical, usage.InodesUsedPercent)
	return state, message
}

func checkNodeVolume(req *Request) (icinga.State, interface{}) {
	host, err := icinga.ParseHost(req.Host)
	if err != nil {
		return icinga.UNKNOWN, "Invalid icinga host.name"
	}
	if host.Type != icinga.TypeNode {
		return icinga.UNKNOWN, "Invalid icinga host type"
	}

	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		return icinga.UNKNOWN, err
	}
	kubeClient := kubernetes.NewForConfigOrDie(config)
	node, err := kubeClient.CoreV1().Nodes().Get(host.ObjectName, metav1.GetOptions{})
	if err != nil {
		return icinga.UNKNOWN, err
	}

	if node == nil {
		return icinga.UNKNOWN, "Node not found"
	}

	hostIP := ""
	for _, address := range node.Status.Addresses {
		if address.Type == apiv1.NodeInternalIP {
			hostIP = address.Address
		}
	}

	if hostIP == "" {
		return icinga.UNKNOWN, "Node InternalIP not found"
	}
	return checkVolume(kubeClient, req, host.AlertNamespace, hostIP, req.VolumeName)
}

func checkPodVolume(req *Request) (icinga.State, interface{}) {
	host, err := icinga.ParseHost(req.Host)
	if err != nil {
		return icinga.UNKNOWN, "Invalid icinga host.name"
	}
	if host.Type != icinga.TypePod {
		return icinga.UNKNOWN, "Invalid icinga host type"
	}

	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		return icinga.UNKNOWN, err
	}
	kubeClient := kubernetes.NewForConfigOrDie(config)

	pod, err := kubeClient.CoreV1().Pods(host.AlertNamespace).Get(host.ObjectName, metav1.GetOptions{})
	if err != nil {
		return icinga.UNKNOWN, err
	}

	for _, volume := range pod.Spec.Volumes {
		if volume.Name == req.VolumeName {
			if volume.PersistentVolumeClaim != nil {
				claim, err := kubeClient.CoreV1().PersistentVolumeClaims(host.AlertNamespace).Get(volume.PersistentVolumeClaim.ClaimName, metav1.GetOptions{})
				if err != nil {
					return icinga.UNKNOWN, err
				}
				volume, err := kubeClient.CoreV1().PersistentVolumes().Get(claim.Spec.VolumeName, metav1.GetOptions{})
				if err != nil {
					return icinga.UNKNOWN, err
				}
				volumePluginName := getPersistentVolumePluginName(&volume.Spec.PersistentVolumeSource)
				if volumePluginName == hostPathPluginName {
					if claim.Spec.StorageClassName != nil {
						class, err := kubeClient.StorageV1beta1().StorageClasses().Get(*claim.Spec.StorageClassName, metav1.GetOptions{})
						if err != nil {
							return icinga.UNKNOWN, err
						}
						if class.Provisioner == "k8s.io/minikube-hostpath" {
							path := fmt.Sprintf("/tmp/hostpath-provisioner/%s", volume.Name)
							return checkVolume(kubeClient, req, host.AlertNamespace, pod.Status.HostIP, path)
						}
					}
				}
				path := fmt.Sprintf("/var/lib/kubelet/pods/%v/volumes/%v/%v", pod.UID, volumePluginName, volume.Name)
				return checkVolume(kubeClient, req, host.AlertNamespace, pod.Status.HostIP, path)
			} else {
				path := fmt.Sprintf("/var/lib/kubelet/pods/%v/volumes/%v/%v", pod.UID, getVolumePluginName(&volume.VolumeSource), volume.Name)
				return checkVolume(kubeClient, req, host.AlertNamespace, pod.Status.HostIP, path)
			}
			break
		}
	}
	return icinga.UNKNOWN, errors.New("Invalid volume source")
}

func NewCmd() *cobra.Command {
	var req Request

	c := &cobra.Command{
		Use:     "check_volume",
		Short:   "Check kubernetes volume",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host")

			if req.NodeStat {
				icinga.Output(checkNodeVolume(&req))
			} else {
				flags.EnsureRequiredFlags(cmd, "volume_name")
				icinga.Output(checkPodVolume(&req))
			}
		},
	}

	c.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	c.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")

	c.Flags().StringVarP(&req.Host, "host", "H", "", "Icinga host name")
	c.Flags().BoolVar(&req.NodeStat, "nodeStat", false, "Checking Node disk size")
	c.Flags().StringVarP(&req.SecretName, "secretName", "s", "", `Kubernetes secret name`)
	c.Flags().StringVarP(&req.VolumeName, "volumeName", "N", "", "Volume name")
	c.Flags().Float64VarP(&req.Warning, "warning", "w", 80.0, "Warning level value (usage percentage)")
	c.Flags().Float64VarP(&req.Critical, "critical", "c", 95.0, "Critical level value (usage percentage)")
	return c
}
