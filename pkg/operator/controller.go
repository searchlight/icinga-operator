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
	"fmt"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	cs "go.searchlight.dev/icinga-operator/client/clientset/versioned"
	mon_informers "go.searchlight.dev/icinga-operator/client/informers/externalversions"
	mon_listers "go.searchlight.dev/icinga-operator/client/listers/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/pkg/icinga"

	"github.com/golang/glog"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	core_listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	reg_util "kmodules.xyz/client-go/admissionregistration/v1beta1"
	"kmodules.xyz/client-go/apiextensions"
	"kmodules.xyz/client-go/tools/queue"
)

type Operator struct {
	Config

	clientConfig *rest.Config
	kubeClient   kubernetes.Interface
	crdClient    crd_cs.Interface
	extClient    cs.Interface
	icingaClient *icinga.Client // TODO: init

	clusterHost *icinga.ClusterHost
	nodeHost    *icinga.NodeHost
	podHost     *icinga.PodHost
	recorder    record.EventRecorder

	kubeInformerFactory informers.SharedInformerFactory
	monInformerFactory  mon_informers.SharedInformerFactory

	// Namespace
	nsInformer cache.SharedIndexInformer
	nsLister   core_listers.NamespaceLister

	// Node
	nodeQueue    *queue.Worker
	nodeInformer cache.SharedIndexInformer
	nodeLister   core_listers.NodeLister

	// Pod
	podQueue    *queue.Worker
	podInformer cache.SharedIndexInformer
	podLister   core_listers.PodLister

	// ClusterAlert
	caQueue    *queue.Worker
	caInformer cache.SharedIndexInformer
	caLister   mon_listers.ClusterAlertLister

	// NodeAlert
	naQueue    *queue.Worker
	naInformer cache.SharedIndexInformer
	naLister   mon_listers.NodeAlertLister

	// PodAlert
	paQueue    *queue.Worker
	paInformer cache.SharedIndexInformer
	paLister   mon_listers.PodAlertLister

	// SearchlightPlugin
	pluginQueue    *queue.Worker
	pluginInformer cache.SharedIndexInformer
	pluginLister   mon_listers.SearchlightPluginLister
}

func (op *Operator) ensureCustomResourceDefinitions() error {
	crds := []*apiextensions.CustomResourceDefinition{
		api.ClusterAlert{}.CustomResourceDefinition(),
		api.NodeAlert{}.CustomResourceDefinition(),
		api.PodAlert{}.CustomResourceDefinition(),
		api.Incident{}.CustomResourceDefinition(),
		api.SearchlightPlugin{}.CustomResourceDefinition(),
	}
	return apiextensions.RegisterCRDs(op.crdClient, crds)
}

func (op *Operator) RunInformers(stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	glog.Info("Starting Searchlight controller")

	go op.kubeInformerFactory.Start(stopCh)
	go op.monInformerFactory.Start(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	for _, v := range op.kubeInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}
	for _, v := range op.monInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	op.nodeQueue.Run(stopCh)
	op.podQueue.Run(stopCh)
	op.caQueue.Run(stopCh)
	op.naQueue.Run(stopCh)
	op.paQueue.Run(stopCh)
	op.pluginQueue.Run(stopCh)

	<-stopCh
	glog.Info("Stopping Searchlight controller")
}

func (op *Operator) Run(stopCh <-chan struct{}) error {
	err := op.MigrateAlerts()
	if err != nil {
		return err
	}

	op.gcIncidents()

	// Create build-in SearchlighPlugin
	if err := op.createBuiltinSearchlightPlugin(); err != nil {
		return err
	}

	go op.RunInformers(stopCh)

	cancel, _ := reg_util.SyncValidatingWebhookCABundle(op.clientConfig, validatingWebhook)

	<-stopCh

	cancel()
	return nil
}
