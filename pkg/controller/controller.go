package controller

import (
	"time"

	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	acs "github.com/appscode/searchlight/client/clientset"
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
	kubeClient clientset.Interface
	extClient  acs.ExtensionInterface
	syncPeriod time.Duration
	opt        *types.Option
}

func New(kubeClient clientset.Interface,
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

func (c *Controller) Setup() {
	log.Infoln("Ensuring ThirdPartyResource")

	if err := c.ensureThirdPartyResource(aci.ResourceNameAlert + "." + aci.V1alpha1SchemeGroupVersion.Group); err != nil {
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
	_, err := c.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(resourceName, metav1.GetOptions{})
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

	_, err = c.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Create(thirdPartyResource)
	return err
}
