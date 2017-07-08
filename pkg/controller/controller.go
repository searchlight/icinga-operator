package controller

import (
	"time"

	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	acs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/data"
	_ "github.com/appscode/searchlight/pkg/controller/host/localhost"
	_ "github.com/appscode/searchlight/pkg/controller/host/node"
	_ "github.com/appscode/searchlight/pkg/controller/host/pod"
	"github.com/appscode/searchlight/pkg/controller/types"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type Controller struct {
	KubeClient clientset.Interface
	ExtClient  acs.ExtensionInterface

	opt *types.Option

	Commands   map[string]*types.IcingaCommands
	syncPeriod time.Duration
}

func New(kubeClient clientset.Interface, extClient acs.ExtensionInterface) (*Controller, error) {
	ctrl := &Controller{
		KubeClient: kubeClient,
		ExtClient:  extClient,
		syncPeriod: 5 * time.Minute,
	}
	cmdList, err := data.LoadIcingaData()
	if err != nil {
		return nil, err
	}
	ctrl.Commands = make(map[string]*types.IcingaCommands)
	for _, cmd := range cmdList.Command {
		vars := make(map[string]data.CommandVar)
		for _, v := range cmd.Vars {
			vars[v.Name] = v
		}
		ctrl.Commands[cmd.Name] = &types.IcingaCommands{
			Name:   cmd.Name,
			Vars:   vars,
			States: cmd.States,
		}
	}
	return ctrl, nil
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
