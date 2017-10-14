package operator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/kutil"
	"github.com/appscode/pat"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
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
}

type Operator struct {
	KubeClient   clientset.Interface
	CRDClient    apiextensionsclient.Interface
	ExtClient    cs.MonitoringV1alpha1Interface
	IcingaClient *icinga.Client // TODO: init

	Opt         Options
	clusterHost *icinga.ClusterHost
	nodeHost    *icinga.NodeHost
	podHost     *icinga.PodHost
	recorder    record.EventRecorder
}

func New(kubeClient clientset.Interface, crdClient apiextensionsclient.Interface, extClient cs.MonitoringV1alpha1Interface, icingaClient *icinga.Client, opt Options) *Operator {
	return &Operator{
		KubeClient:   kubeClient,
		CRDClient:    crdClient,
		ExtClient:    extClient,
		IcingaClient: icingaClient,
		Opt:          opt,
		clusterHost:  icinga.NewClusterHost(kubeClient, extClient, icingaClient),
		nodeHost:     icinga.NewNodeHost(kubeClient, extClient, icingaClient),
		podHost:      icinga.NewPodHost(kubeClient, extClient, icingaClient),
		recorder:     eventer.NewEventRecorder(kubeClient, "Searchlight operator"),
	}
}

func (op *Operator) Setup() error {
	return op.ensureCustomResourceDefinitions()
}

func (op *Operator) ensureCustomResourceDefinitions() error {
	crds := []*apiextensions.CustomResourceDefinition{
		api.ClusterAlert{}.CustomResourceDefinition(),
		api.NodeAlert{}.CustomResourceDefinition(),
		api.PodAlert{}.CustomResourceDefinition(),
	}
	for _, crd := range crds {
		_, err := op.CRDClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(crd.Name, metav1.GetOptions{})
		if kerr.IsNotFound(err) {
			_, err = op.CRDClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
			if err != nil {
				return err
			}
		}
	}
	return kutil.WaitForCRDReady(op.KubeClient.CoreV1().RESTClient(), crds)
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

func (op *Operator) Run() {
	go op.WatchNamespaces()
	go op.WatchPods()
	go op.WatchNodes()
	go op.WatchNamespaces()
	go op.WatchPodAlerts()
	go op.WatchNodeAlerts()
	go op.WatchClusterAlerts()
}

func (op *Operator) RunAndHold() {
	op.Run()
	go op.RunAPIServer()

	m := pat.New()
	m.Get("/metrics", promhttp.Handler())
	http.Handle("/", m)
	log.Infoln("Listening on", op.Opt.WebAddress)
	log.Fatal(http.ListenAndServe(op.Opt.WebAddress, nil))
}
