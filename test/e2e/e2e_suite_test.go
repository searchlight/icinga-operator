package e2e_test

import (
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	cs "github.com/appscode/searchlight/client"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/operator"
	"github.com/appscode/searchlight/test/e2e"
	"github.com/appscode/searchlight/test/e2e/framework"
	. "github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var provider string
var storageClass string

func init() {
	flag.StringVar(&provider, "provider", "minikube", "Kubernetes cloud provider")
	flag.StringVar(&storageClass, "storageclass", "", "Kubernetes StorageClass name")
}

const (
	TIMEOUT = 20 * time.Minute
)

var (
	op   *operator.Operator
	root *framework.Framework
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TIMEOUT)

	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "e2e Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {
	// Kubernetes config
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube/config")
	By("Using kubeconfig from " + kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())
	// Clients
	kubeClient := kubernetes.NewForConfigOrDie(config)
	apiExtKubeClient := crd_cs.NewForConfigOrDie(config)
	extClient := cs.NewForConfigOrDie(config)
	// Framework
	root = framework.New(kubeClient, apiExtKubeClient, extClient, nil, provider, storageClass)

	e2e.PrintSeparately("Using namespace " + root.Namespace())

	// Create namespace
	err = root.CreateNamespace()
	Expect(err).NotTo(HaveOccurred())

	// Create Searchlight deployment
	slDeployment := root.Invoke().DeploymentSearchlight()
	err = root.CreateDeployment(slDeployment)
	Expect(err).NotTo(HaveOccurred())
	By("Waiting for Running pods")
	root.EventuallyDeployment(slDeployment.ObjectMeta).Should(HaveRunningPods(*slDeployment.Spec.Replicas))

	// Create Searchlight service
	slService := root.Invoke().ServiceSearchlight()
	err = root.CreateService(slService)
	Expect(err).NotTo(HaveOccurred())
	root.EventuallyServiceLoadBalancer(slService.ObjectMeta, "icinga").Should(BeTrue())

	// Get Icinga Ingress Hostname
	endpoint, err := root.GetServiceEndpoint(slService.ObjectMeta, "icinga")
	Expect(err).NotTo(HaveOccurred())

	// Icinga Config
	cfg := &icinga.Config{
		Endpoint: fmt.Sprintf("https://%v/v1", endpoint),
	}

	cfg.BasicAuth.Username = e2e.ICINGA_API_USER
	cfg.BasicAuth.Password = e2e.ICINGA_API_PASSWORD

	// Icinga Client
	icingaClient := icinga.NewClient(*cfg)
	root = root.SetIcingaClient(icingaClient)
	root.EventuallyIcingaAPI().Should(Succeed())

	icingawebEndpoint, err := root.GetServiceEndpoint(slService.ObjectMeta, "ui")
	Expect(err).NotTo(HaveOccurred())
	fmt.Println()
	fmt.Println("Icingaweb2:     ", fmt.Sprintf("http://%v/icingaweb2", icingawebEndpoint))
	fmt.Println("Login password: ", e2e.ICINGA_WEB_UI_PASSWORD)
	fmt.Println()

	// Controller
	op = operator.New(kubeClient, apiExtKubeClient, extClient, icingaClient, operator.Options{})
	err = op.Setup()
	Expect(err).NotTo(HaveOccurred())
	op.Run(nil)
})

var _ = AfterSuite(func() {
	err := root.DeleteNamespace()
	Expect(err).NotTo(HaveOccurred())
	e2e.PrintSeparately("Deleted namespace")
})
