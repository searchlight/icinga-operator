package operator

import (
	"time"

	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	cs "github.com/appscode/searchlight/client/clientset/versioned"
	mon_informers "github.com/appscode/searchlight/client/informers/externalversions"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/icinga"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Config struct {
	ConfigRoot       string
	ConfigSecretName string
	OpsAddress       string
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
	CRDClient      crd_cs.ApiextensionsV1beta1Interface
	IcingaClient   *icinga.Client // TODO: init
	AdmissionHooks []hooks.AdmissionHook
}

func NewOperatorConfig(clientConfig *rest.Config) *OperatorConfig {
	return &OperatorConfig{
		ClientConfig: clientConfig,
	}
}

func (c *OperatorConfig) New() (*Operator, error) {
	op := &Operator{
		Config:              c.Config,
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
	op.initNamespaceWatcher()
	op.initNodeWatcher()
	op.initPodWatcher()
	op.initClusterAlertWatcher()
	op.initNodeAlertWatcher()
	op.initPodAlertWatcher()

	return op, nil
}
