package controller

import (
	"time"

	"github.com/appscode/log"
	tapi "github.com/appscode/searchlight/api"
	tcs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/controller/types"
	"github.com/appscode/searchlight/pkg/eventer"
	icinga "github.com/appscode/searchlight/pkg/icinga/client"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/record"
)

type Controller struct {
	KubeClient   clientset.Interface
	ExtClient    tcs.ExtensionInterface
	IcingaClient *icinga.IcingaClient // TODO: init

	// Event Recorder
	recorder   record.EventRecorder
	opt        *types.Option
	syncPeriod time.Duration
}

func New(kubeClient clientset.Interface, extClient tcs.ExtensionInterface) *Controller {
	return &Controller{
		KubeClient: kubeClient,
		ExtClient:  extClient,
		recorder:   eventer.NewEventRecorder(kubeClient, "Searchlight operator"),
		syncPeriod: 5 * time.Minute,
	}
}

func (c *Controller) Setup() {
	log.Infoln("Ensuring ThirdPartyResource")

	if err := c.ensureThirdPartyResource(tapi.ResourceNamePodAlert + "." + tapi.V1alpha1SchemeGroupVersion.Group); err != nil {
		log.Fatalln(err)
	}
	if err := c.ensureThirdPartyResource(tapi.ResourceNameNodeAlert + "." + tapi.V1alpha1SchemeGroupVersion.Group); err != nil {
		log.Fatalln(err)
	}
	if err := c.ensureThirdPartyResource(tapi.ResourceNameClusterAlert + "." + tapi.V1alpha1SchemeGroupVersion.Group); err != nil {
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
				Name: tapi.V1alpha1SchemeGroupVersion.Version,
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
