package controller

import (
	"time"
	"encoding/json"
	"net/http"
	"path/filepath"
	"sync"
	"time"
	"github.com/appscode/log"
	"github.com/appscode/pat"
	"github.com/appscode/log"
	tapi "github.com/appscode/searchlight/api"
	tcs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/icinga"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/record"
	"github.com/derekparker/delve/pkg/dwarf/op"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"fmt"
)

type Options struct {
	Master     string
	KubeConfig string

	ConfigRoot       string
	IcingaSecretName  string
	APIAddress        string
	WebAddress        string

	EnableAnalytics bool
}

type Controller struct {
	KubeClient   clientset.Interface
	ExtClient    tcs.ExtensionInterface
	IcingaClient *icinga.Client // TODO: init

	clusterHost *icinga.ClusterHost
	nodeHost    *icinga.NodeHost
	podHost     *icinga.PodHost
	recorder    record.EventRecorder
	SyncPeriod  time.Duration
}

func New(kubeClient clientset.Interface, extClient tcs.ExtensionInterface, icingaClient *icinga.Client) *Controller {
	return &Controller{
		KubeClient:   kubeClient,
		ExtClient:    extClient,
		IcingaClient: icingaClient,
		clusterHost:  icinga.NewClusterHost(kubeClient, extClient, icingaClient),
		nodeHost:     icinga.NewNodeHost(kubeClient, extClient, icingaClient),
		podHost:      icinga.NewPodHost(kubeClient, extClient, icingaClient),
		recorder:     eventer.NewEventRecorder(kubeClient, "Searchlight operator"),
		SyncPeriod:   5 * time.Minute,
	}
}

func (c *Controller) Setup() error {
	log.Infoln("Ensuring ThirdPartyResource")

	if err := c.ensureThirdPartyResource(tapi.ResourceNamePodAlert + "." + tapi.V1alpha1SchemeGroupVersion.Group); err != nil {
		return err
	}
	if err := c.ensureThirdPartyResource(tapi.ResourceNameNodeAlert + "." + tapi.V1alpha1SchemeGroupVersion.Group); err != nil {
		return err
	}
	if err := c.ensureThirdPartyResource(tapi.ResourceNameClusterAlert + "." + tapi.V1alpha1SchemeGroupVersion.Group); err != nil {
		return err
	}
	return nil
}

func (c *Controller) ensureThirdPartyResource(resourceName string) error {
	_, err := c.KubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(resourceName, metav1.GetOptions{})
	if !kerr.IsNotFound(err) {
		return err
	}

	thirdPartyResource := &extensions.ThirdPartyResource{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "ThirdPartyResource",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceName,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Description: "Searchlight by AppsCode - Alerts for Kubernetes",
		Versions: []extensions.APIVersion{
			{
				Name: tapi.V1alpha1SchemeGroupVersion.Version,
			},
		},
	}

	_, err = c.KubeClient.ExtensionsV1beta1().ThirdPartyResources().Create(thirdPartyResource)
	return err
}

func (c *Controller) RunAPIServer() {
	m := pat.New()
	// For notification acknowledgement
	ackPattern := fmt.Sprintf("/monitoring.appscode.com/v1alpha1/namespaces/%s/%s/%s",PathParamNamespace, PathParamType, PathParamName)
	ackHandler := func(w http.ResponseWriter, r *http.Request) {
		Acknowledge(c.IcingaClient, w, r)
	}
	m.Post(ackPattern, http.HandlerFunc(ackHandler))

	http.Handle("/", m)
	log.Infoln("Listening on", apiAddress)
	log.Fatal(http.ListenAndServe(apiAddress, nil))


	// router is default HTTP request multiplexer for kubed. It matches the URL of each
	// incoming request against a list of registered patterns with their associated
	// methods and calls the handler for the pattern that most closely matches the
	// URL.
	//
	// Pattern matching attempts each pattern in the order in which they were
	// registered.
	router := pat.New()
	if op.Config.APIServer.EnableSearchIndex {
		op.SearchIndex.RegisterRouters(router)
	}
	// Enable pod -> service, service -> serviceMonitor indexing
	if op.Config.APIServer.EnableReverseIndex {
		router.Get("/api/v1/namespaces/:namespace/:resource/:name/services", http.HandlerFunc(op.ReverseIndex.Service.ServeHTTP))
		if util.IsPreferredAPIResource(op.KubeClient, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRServiceMonitorsKind) {
			// Add Indexer only if Server support this resource
			router.Get("/apis/"+prom.TPRGroup+"/"+prom.TPRVersion+"/namespaces/:namespace/:resource/:name/"+prom.TPRServiceMonitorName, http.HandlerFunc(op.ReverseIndex.ServiceMonitor.ServeHTTP))
		}
		if util.IsPreferredAPIResource(op.KubeClient, prom.TPRGroup+"/"+prom.TPRVersion, prom.TPRPrometheusesKind) {
			// Add Indexer only if Server support this resource
			router.Get("/apis/"+prom.TPRGroup+"/"+prom.TPRVersion+"/namespaces/:namespace/:resource/:name/"+prom.TPRPrometheusName, http.HandlerFunc(op.ReverseIndex.Prometheus.ServeHTTP))
		}
	}

	router.Get("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }))
	router.Get("/metadata", http.HandlerFunc(op.metadataHandler))
	log.Fatalln(http.ListenAndServe(op.Config.APIServer.Address, router))
}


func (c *Controller) Run() {
	go c.WatchNamespaces()
	go c.WatchPods()
	go c.WatchNodes()
	go c.WatchNamespaces()
	go c.WatchPodAlerts()
	go c.WatchNodeAlerts()
	go c.WatchClusterAlerts()
}
