package operator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/appscode/log"
	"github.com/appscode/pat"
	tapi "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	tapi_v1alpha1 "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	tcs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
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
	KubeClient       clientset.Interface
	ApiExtKubeClient apiextensionsclient.Interface
	ExtClient        tcs.MonitoringV1alpha1Interface
	IcingaClient     *icinga.Client // TODO: init

	Opt         Options
	clusterHost *icinga.ClusterHost
	nodeHost    *icinga.NodeHost
	podHost     *icinga.PodHost
	recorder    record.EventRecorder
}

func New(kubeClient clientset.Interface, apiExtKubeClient apiextensionsclient.Interface, extClient tcs.MonitoringV1alpha1Interface, icingaClient *icinga.Client, opt Options) *Operator {
	return &Operator{
		KubeClient:       kubeClient,
		ApiExtKubeClient: apiExtKubeClient,
		ExtClient:        extClient,
		IcingaClient:     icingaClient,
		Opt:              opt,
		clusterHost:      icinga.NewClusterHost(kubeClient, extClient, icingaClient),
		nodeHost:         icinga.NewNodeHost(kubeClient, extClient, icingaClient),
		podHost:          icinga.NewPodHost(kubeClient, extClient, icingaClient),
		recorder:         eventer.NewEventRecorder(kubeClient, "Searchlight operator"),
	}
}

func (op *Operator) Setup() error {
	log.Infoln("Ensuring CustomResourceDefinition")

	if err := op.ensureCustomResourceDefinition(tapi.ResourceKindClusterAlert, tapi.ResourceTypeClusterAlert); err != nil {
		return err
	}
	if err := op.ensureCustomResourceDefinition(tapi.ResourceKindNodeAlert, tapi.ResourceTypeNodeAlert); err != nil {
		return err
	}
	if err := op.ensureCustomResourceDefinition(tapi.ResourceKindPodAlert, tapi.ResourceTypePodAlert); err != nil {
		return err
	}
	return nil
}

func (op *Operator) ensureCustomResourceDefinition(resourceKind, resourceType string) error {
	name := resourceType + "." + tapi_v1alpha1.SchemeGroupVersion.Group
	_, err := op.ApiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions().Get(name, metav1.GetOptions{})
	if !kerr.IsNotFound(err) {
		return err
	}

	crd := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   tapi_v1alpha1.SchemeGroupVersion.Group,
			Version: tapi_v1alpha1.SchemeGroupVersion.Version,
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: resourceType,
				Kind:   resourceKind,
			},
		},
	}

	_, err = op.ApiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	return err
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
