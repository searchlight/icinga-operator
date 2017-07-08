package controller

import (
	"time"

	"github.com/appscode/errors"
	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	acs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/client/icinga"
	_ "github.com/appscode/searchlight/pkg/controller/host/localhost"
	_ "github.com/appscode/searchlight/pkg/controller/host/node"
	_ "github.com/appscode/searchlight/pkg/controller/host/pod"
	"github.com/appscode/searchlight/pkg/controller/types"
	"github.com/appscode/searchlight/pkg/stash"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type Controller struct {
	KubeClient clientset.Interface
	ExtClient  acs.ExtensionInterface

	opt *types.Option

	CommandList map[string]*types.IcingaCommands
	syncPeriod  time.Duration
}

func New2(kubeClient clientset.Interface,
	icingaClient *icinga.IcingaClient,
	extClient acs.ExtensionInterface,
	storage stash.Storage) *Controller {
	data, err := getIcingaDataMap()
	if err != nil {
		log.Errorln("Icinga data not found")
	}
	ctx := &types.Option{
		KubeClient:   kubeClient,
		ExtClient:    extClient,
		IcingaData:   data,
		IcingaClient: icingaClient,
		Storage:      storage,
	}
	return &Controller{opt: ctx}
}

func New(kubeClient clientset.Interface, extClient acs.ExtensionInterface) (*Controller, error) {
	ctrl := &Controller{
		KubeClient: kubeClient,
		ExtClient:  extClient,
		syncPeriod: 5 * time.Minute,
	}
	data, err := data.LoadIcingaData()
	if err != nil {
		return nil, err
	}
	ctrl.CommandList = make(map[string]*types.IcingaCommands)
	for _, command := range data.Command {
		varsMap := make(map[string]data.CommandVar)
		for _, v := range command.Vars {
			varsMap[v.Name] = v
		}

		ctrl.CommandList[command.Name] = &types.IcingaCommands{
			HostType: command.ObjectToHost,
			VarInfo:  varsMap,
		}
	}
	return ctrl, nil
}

func getIcingaDataMap() (map[string]*types.IcingaCommands, error) {
	data, err := data.LoadIcingaData()
	if err != nil {
		return nil, errors.FromErr(err).Err()
	}

	m := make(map[string]*types.IcingaCommands)
	for _, command := range data.Command {
		varsMap := make(map[string]data.CommandVar)
		for _, v := range command.Vars {
			varsMap[v.Name] = v
		}

		m[command.Name] = &types.IcingaCommands{
			HostType: command.ObjectToHost,
			VarInfo:  varsMap,
		}
	}
	return m, nil
}

func (c *Controller) Setup() {
	log.Infoln("Ensuring ThirdPartyResource")

	if err := c.ensureThirdPartyResource(aci.ResourceNamePodAlert + "." + aci.V1alpha1SchemeGroupVersion.Group); err != nil {
		log.Fatalln(err)
	}
	if err := c.ensureThirdPartyResource(aci.ResourceNameNodeAlert + "." + aci.V1alpha1SchemeGroupVersion.Group); err != nil {
		log.Fatalln(err)
	}
	if err := c.ensureThirdPartyResource(aci.ResourceNameClusterAlert + "." + aci.V1alpha1SchemeGroupVersion.Group); err != nil {
		log.Fatalln(err)
	}
}

func (c *Controller) ensureThirdPartyResource(resourceName string) error {
	_, err := c.KubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(resourceName, metav1.GetOptions{})
	if !kerr.IsNotFound(err) {
		return err
	}

	thirdPartyResource := &extensions.ThirdPartyResource{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "ThirdPartyResource",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Description: "Searchlight by AppsCode - Alerts for Kubernetes",
		Versions: []extensions.APIVersion{
			{
				Name: aci.V1alpha1SchemeGroupVersion.Version,
			},
		},
	}

	_, err = c.KubeClient.ExtensionsV1beta1().ThirdPartyResources().Create(thirdPartyResource)
	return err
}

func (c *Controller) Run() {
	go c.WatchNamespaces()
	go c.WatchPods()
	go c.WatchNodes()
	go c.WatchNamespaces()
	go c.WatchPodAlerts()
	go c.WatchNodeAlerts()
	go c.WatchClusterAlerts()
}
