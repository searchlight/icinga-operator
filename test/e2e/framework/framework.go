package framework

import (
	"github.com/appscode/go/crypto/rand"
	tcs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/icinga"
	clientset "k8s.io/client-go/kubernetes"
)

type Framework struct {
	kubeClient   clientset.Interface
	extClient    tcs.ExtensionInterface
	icingaClient *icinga.Client
	namespace    string
	name         string
	Provider     string
	storageClass string
}

func New(kubeClient clientset.Interface, extClient tcs.ExtensionInterface, icingaClient *icinga.Client, provider, storageClass string) *Framework {
	return &Framework{
		kubeClient:   kubeClient,
		extClient:    extClient,
		icingaClient: icingaClient,
		name:         "searchlight-operator",
		namespace:    rand.WithUniqSuffix("searchlight"), // "searchlight-42e4fy",
		Provider:     provider,
		storageClass: storageClass,
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
