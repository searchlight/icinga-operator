package framework

import (
	"github.com/appscode/go/crypto/rand"
	cs "github.com/appscode/searchlight/client/clientset/versioned"
	mon_cs "github.com/appscode/searchlight/client/clientset/versioned/typed/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/client-go/kubernetes"
)

type Framework struct {
	kubeClient       kubernetes.Interface
	apiExtKubeClient crd_cs.ApiextensionsV1beta1Interface
	extClient        cs.Interface
	icingaClient     *icinga.Client
	namespace        string
	name             string
	Provider         string
	storageClass     string
}

func New(kubeClient kubernetes.Interface, apiExtKubeClient crd_cs.ApiextensionsV1beta1Interface, extClient cs.Interface, icingaClient *icinga.Client, provider, storageClass string) *Framework {
	return &Framework{
		kubeClient:       kubeClient,
		apiExtKubeClient: apiExtKubeClient,
		extClient:        extClient,
		icingaClient:     icingaClient,
		name:             "searchlight-operator",
		namespace:        rand.WithUniqSuffix("searchlight"), // "searchlight-42e4fy",
		Provider:         provider,
		storageClass:     storageClass,
	}
}

func (f *Framework) SetIcingaClient(icingaClient *icinga.Client) *Framework {
	f.icingaClient = icingaClient
	return f
}

func (f *Framework) Invoke() *Invocation {
	return &Invocation{
		Framework: f,
		app:       rand.WithUniqSuffix("searchlight-e2e"),
	}
}

func (f *Framework) KubeClient() kubernetes.Interface {
	return f.kubeClient
}

func (f *Framework) MonitoringClient() mon_cs.MonitoringV1alpha1Interface {
	return f.extClient.MonitoringV1alpha1()
}

type Invocation struct {
	*Framework
	app string
}
