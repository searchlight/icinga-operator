package operator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/appscode/go/log"
	apiext_util "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/tools/queue"
	"github.com/appscode/pat"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/clientset/versioned"
	searchlightinformers "github.com/appscode/searchlight/client/informers/externalversions"
	slite_listers "github.com/appscode/searchlight/client/listers/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/golang/glog"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	ecs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	core_listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

type Options struct {
	Master     string
	KubeConfig string

	ConfigRoot       string
	ConfigSecretName string
	APIAddress       string
	WebAddress       string
	ResyncPeriod     time.Duration
	MaxNumRequeues   int
	NumThreads       int
}

type Operator struct {
	KubeClient   kubernetes.Interface
	CRDClient    ecs.ApiextensionsV1beta1Interface
	ExtClient    cs.Interface
	IcingaClient *icinga.Client // TODO: init

	kubeInformerFactory        informers.SharedInformerFactory
	searchlightInformerFactory searchlightinformers.SharedInformerFactory

	options     Options
	clusterHost *icinga.ClusterHost
	nodeHost    *icinga.NodeHost
	podHost     *icinga.PodHost
	recorder    record.EventRecorder

	// Namespace
	nsInformer cache.SharedIndexInformer

	// Node
	nQueue    *queue.Worker
	nInformer cache.SharedIndexInformer
	nLister   core_listers.NodeLister

	// Pod
	pQueue    *queue.Worker
	pInformer cache.SharedIndexInformer
	pLister   core_listers.PodLister

	// ClusterAlert
	caQueue    *queue.Worker
	caInformer cache.SharedIndexInformer
	caLister   slite_listers.ClusterAlertLister

	// NodeAlert
	naQueue    *queue.Worker
	naInformer cache.SharedIndexInformer
	naLister   slite_listers.NodeAlertLister

	// PodAlert
	paQueue    *queue.Worker
	paInformer cache.SharedIndexInformer
	paLister   slite_listers.PodAlertLister
}

const annotationAlertsName = "monitoring.appscode.com/alerts"

func New(kubeClient kubernetes.Interface, crdClient ecs.ApiextensionsV1beta1Interface, extClient cs.Interface, icingaClient *icinga.Client, opt Options) *Operator {
	return &Operator{
		KubeClient:                 kubeClient,
		kubeInformerFactory:        informers.NewSharedInformerFactory(kubeClient, opt.ResyncPeriod),
		CRDClient:                  crdClient,
		ExtClient:                  extClient,
		searchlightInformerFactory: searchlightinformers.NewSharedInformerFactory(extClient, opt.ResyncPeriod),
		IcingaClient:               icingaClient,
		options:                    opt,
		clusterHost:                icinga.NewClusterHost(icingaClient),
		nodeHost:                   icinga.NewNodeHost(icingaClient),
		podHost:                    icinga.NewPodHost(icingaClient),
		recorder:                   eventer.NewEventRecorder(kubeClient, "Searchlight operator"),
	}
}

func (op *Operator) Setup() error {
	if err := op.ensureCustomResourceDefinitions(); err != nil {
		return err
	}
	op.initNamespaceWatcher()
	op.initNodeWatcher()
	op.initPodWatcher()
	op.initClusterAlertWatcher()
	op.initNodeAlertWatcher()
	op.initPodAlertWatcher()
	return nil
}

func (op *Operator) ensureCustomResourceDefinitions() error {
	crds := []*crd_api.CustomResourceDefinition{
		api.ClusterAlert{}.CustomResourceDefinition(),
		api.NodeAlert{}.CustomResourceDefinition(),
		api.PodAlert{}.CustomResourceDefinition(),
	}
	return apiext_util.RegisterCRDs(op.CRDClient, crds)
}

func (op *Operator) RunAPIServer() {
	router := pat.New()

	// For notification acknowledgement
	ackPattern := fmt.Sprintf("/monitoring.appscode.com/v1alpha1/namespaces/%s/%s/%s", PathParamNamespace, PathParamType, PathParamName)
	ackHandler := func(w http.ResponseWriter, r *http.Request) {
		Acknowledge(op.IcingaClient, w, r)
	}
	router.Post(ackPattern, http.HandlerFunc(ackHandler))

	router.Get("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))

	log.Infoln("Listening on", op.options.APIAddress)
	log.Fatal(http.ListenAndServe(op.options.APIAddress, router))
}

func (op *Operator) Run(stopCh chan struct{}) {
	defer runtime.HandleCrash()

	glog.Info("Starting Searchlight controller")

	go op.kubeInformerFactory.Start(stopCh)
	go op.searchlightInformerFactory.Start(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	for _, v := range op.kubeInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}
	for _, v := range op.searchlightInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	op.nQueue.Run(stopCh)
	op.pQueue.Run(stopCh)
	op.caQueue.Run(stopCh)
	op.naQueue.Run(stopCh)
	op.paQueue.Run(stopCh)

	<-stopCh
	glog.Info("Stopping Searchlight controller")
}
