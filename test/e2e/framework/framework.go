package framework

import (
	"github.com/appscode/go/crypto/rand"
	tcs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	clientset "k8s.io/client-go/kubernetes"
)

type Framework struct {
	kubeClient       clientset.Interface
	apiExtKubeClient apiextensionsclient.Interface
	extClient        tcs.MonitoringV1alpha1Interface
	icingaClient     *icinga.Client
	namespace        string
	name             string
	Provider         string
	storageClass     string
}

func New(kubeClient clientset.Interface, apiExtKubeClient apiextensionsclient.Interface, extClient tcs.MonitoringV1alpha1Interface, icingaClient *icinga.Client, provider, storageClass string) *Framework {
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

type Invocation struct {
	*Framework
	app string
}
