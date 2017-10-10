package operator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/appscode/go/log"
	apiext_util "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/pat"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	slite_listers "github.com/appscode/searchlight/listers/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/golang/glog"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	ecs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	core_listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
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
}

type Operator struct {
	KubeClient   kubernetes.Interface
	CRDClient    ecs.ApiextensionsV1beta1Interface
	ExtClient    cs.MonitoringV1alpha1Interface
	IcingaClient *icinga.Client // TODO: init

	Opt         Options
	clusterHost *icinga.ClusterHost
	nodeHost    *icinga.NodeHost
	podHost     *icinga.PodHost
	recorder    record.EventRecorder

	// Namespace
	nsIndexer  cache.Indexer
	nsInformer cache.Controller

	// Node
	nQueue    workqueue.RateLimitingInterface
	nIndexer  cache.Indexer
	nInformer cache.Controller
	nLister   core_listers.NodeLister

	// Pod
	pQueue    workqueue.RateLimitingInterface
	pIndexer  cache.Indexer
	pInformer cache.Controller
	pLister   core_listers.PodLister

	// ClusterAlert
	caQueue    workqueue.RateLimitingInterface
	caIndexer  cache.Indexer
	caInformer cache.Controller
	caLister   slite_listers.ClusterAlertLister

	// NodeAlert
	naQueue    workqueue.RateLimitingInterface
	naIndexer  cache.Indexer
	naInformer cache.Controller
	naLister   slite_listers.NodeAlertLister

	// PodAlert
	paQueue    workqueue.RateLimitingInterface
	paIndexer  cache.Indexer
	paInformer cache.Controller
	paLister   slite_listers.PodAlertLister
}

func New(kubeClient kubernetes.Interface, crdClient ecs.ApiextensionsV1beta1Interface, extClient cs.MonitoringV1alpha1Interface, icingaClient *icinga.Client, opt Options) *Operator {
	return &Operator{
		KubeClient:   kubeClient,
		CRDClient:    crdClient,
		ExtClient:    extClient,
		IcingaClient: icingaClient,
		Opt:          opt,
		clusterHost:  icinga.NewClusterHost(icingaClient),
		nodeHost:     icinga.NewNodeHost(icingaClient),
		podHost:      icinga.NewPodHost(icingaClient),
		recorder:     eventer.NewEventRecorder(kubeClient, "Searchlight operator"),
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

	log.Infoln("Listening on", op.Opt.APIAddress)
	log.Fatal(http.ListenAndServe(op.Opt.APIAddress, router))
}

func (op *Operator) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer op.nQueue.ShutDown()
	defer op.pQueue.ShutDown()
	defer op.caQueue.ShutDown()
	defer op.naQueue.ShutDown()
	defer op.paQueue.ShutDown()
	glog.Info("Starting Searchlight controller")

	go op.nsInformer.Run(stopCh)
	go op.nInformer.Run(stopCh)
	go op.pInformer.Run(stopCh)
	go op.caInformer.Run(stopCh)
	go op.naInformer.Run(stopCh)
	go op.paInformer.Run(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	if !cache.WaitForCacheSync(stopCh, op.nsInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	if !cache.WaitForCacheSync(stopCh, op.nInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	if !cache.WaitForCacheSync(stopCh, op.pInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	if !cache.WaitForCacheSync(stopCh, op.caInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	if !cache.WaitForCacheSync(stopCh, op.naInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}
	if !cache.WaitForCacheSync(stopCh, op.paInformer.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(op.runNodeWatcher, time.Second, stopCh)
		go wait.Until(op.runPodWatcher, time.Second, stopCh)
		go wait.Until(op.runClusterAlertWatcher, time.Second, stopCh)
		go wait.Until(op.runNodeAlertWatcher, time.Second, stopCh)
		go wait.Until(op.runPodAlertWatcher, time.Second, stopCh)
	}

	<-stopCh
	glog.Info("Stopping Searchlight controller")
}
