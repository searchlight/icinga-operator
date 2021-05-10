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

package framework

import (
	cs "go.searchlight.dev/icinga-operator/client/clientset/versioned"
	mon_cs "go.searchlight.dev/icinga-operator/client/clientset/versioned/typed/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/pkg/icinga"

	"gomodules.xyz/x/crypto/rand"
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
