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

package operator

import (
	"time"

	cs "go.searchlight.dev/icinga-operator/client/clientset/versioned"
	mon_informers "go.searchlight.dev/icinga-operator/client/informers/externalversions"
	"go.searchlight.dev/icinga-operator/pkg/eventer"
	"go.searchlight.dev/icinga-operator/pkg/icinga"

	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	reg_util "kmodules.xyz/client-go/admissionregistration/v1beta1"
	"kmodules.xyz/client-go/discovery"
	hooks "kmodules.xyz/webhook-runtime/admission/v1beta1"
)

const (
	validatingWebhook = "admission.monitoring.appscode.com"
)

type Config struct {
	ConfigRoot       string
	ConfigSecretName string
	ResyncPeriod     time.Duration
	MaxNumRequeues   int
	NumThreads       int
	IncidentTTL      time.Duration
	// V logging level, the value of the -v flag
	Verbosity string
}

type OperatorConfig struct {
	Config

	ClientConfig   *rest.Config
	KubeClient     kubernetes.Interface
	ExtClient      cs.Interface
	CRDClient      crd_cs.Interface
	IcingaClient   *icinga.Client // TODO: init
	AdmissionHooks []hooks.AdmissionHook
}

func NewOperatorConfig(clientConfig *rest.Config) *OperatorConfig {
	return &OperatorConfig{
		ClientConfig: clientConfig,
	}
}

func (c *OperatorConfig) New() (*Operator, error) {
	if err := discovery.IsDefaultSupportedVersion(c.KubeClient); err != nil {
		return nil, err
	}

	op := &Operator{
		Config:              c.Config,
		clientConfig:        c.ClientConfig,
		kubeClient:          c.KubeClient,
		kubeInformerFactory: informers.NewSharedInformerFactory(c.KubeClient, c.ResyncPeriod),
		crdClient:           c.CRDClient,
		extClient:           c.ExtClient,
		monInformerFactory:  mon_informers.NewSharedInformerFactory(c.ExtClient, c.ResyncPeriod),
		icingaClient:        c.IcingaClient,
		clusterHost:         icinga.NewClusterHost(c.IcingaClient, c.Verbosity),
		nodeHost:            icinga.NewNodeHost(c.IcingaClient, c.Verbosity),
		podHost:             icinga.NewPodHost(c.IcingaClient, c.Verbosity),
		recorder:            eventer.NewEventRecorder(c.KubeClient, "Searchlight operator"),
	}

	if err := op.ensureCustomResourceDefinitions(); err != nil {
		return nil, err
	}
	if err := reg_util.UpdateValidatingWebhookCABundle(op.clientConfig, validatingWebhook); err != nil {
		return nil, err
	}
	op.initNamespaceWatcher()
	op.initNodeWatcher()
	op.initPodWatcher()
	op.initClusterAlertWatcher()
	op.initNodeAlertWatcher()
	op.initPodAlertWatcher()
	op.initPluginWatcher()
	return op, nil
}
