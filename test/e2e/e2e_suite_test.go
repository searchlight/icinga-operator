package e2e_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	tcs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/test/e2e"
	"github.com/appscode/searchlight/test/e2e/framework"
	"github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	TIMEOUT = 20 * time.Minute
)

var (
	ctrl *controller.Controller
	root *framework.Framework
)

func TestE2e(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TIMEOUT)

	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "e2e Suite", []Reporter{junitReporter})
}

const KUBE = "minikube"

var _ = BeforeSuite(func() {
	userHome, err := homedir.Dir()
	Expect(err).NotTo(HaveOccurred())

	// Kubernetes config
	kubeconfigPath := filepath.Join(userHome, ".kube/config")
	By("Using kubeconfig from " + kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())

	// Clients
	kubeClient := clientset.NewForConfigOrDie(config)
	extClient := tcs.NewForConfigOrDie(config)
	// Framework
	root = framework.New(kubeClient, extClient, nil)

	e2e.PrintSeparately("Using namespace: ", root.Namespace())
	By("Using namespace " + root.Namespace())

	// Create namespace
	err = root.CreateNamespace()
	Expect(err).NotTo(HaveOccurred())

	// Create Icinga secret
	configure := &icinga.Configurator{
		ConfigRoot: filepath.Join(userHome),
		Expiry:     10 * 365 * 24 * time.Hour,
	}
	cfg, err := configure.LoadIcingaConfig()
	Expect(err).NotTo(HaveOccurred())
	icingaSecret, err := root.Invoke().SecretSearchlight(filepath.Join(configure.ConfigRoot, "icinga2"))
	Expect(err).NotTo(HaveOccurred())
	err = root.CreateSecret(icingaSecret)
	Expect(err).NotTo(HaveOccurred())

	// Create Searchlight deployment
	searchlightDeployment := root.Invoke().DeploymentExtensionSearchlight()
	err = root.CreateDeploymentExtension(searchlightDeployment)
	Expect(err).NotTo(HaveOccurred())

	//
	for {
		deployment, err := root.GetDeploymentExtension(searchlightDeployment.ObjectMeta)
		Expect(err).NotTo(HaveOccurred())

		if deployment.Status.AvailableReplicas != 0 {
			e2e.PrintSeparately("Searchlight deployment is available")
			break
		}
		fmt.Println("Waiting for Searchlight deployment to be available")
		time.Sleep(5 * time.Second)
	}

	// Create Searchlight svc
	searchlightService := root.Invoke().ServiceSearchlight()
	err = root.CreateService(searchlightService)
	Expect(err).NotTo(HaveOccurred())

	// Get Icinga Ingress Hostname
	icingaHost, err := root.GetServiceIngressHost(searchlightService.ObjectMeta)
	Expect(err).NotTo(HaveOccurred())

	icingaClient := icinga.NewClient(*cfg)
	icingaClient.SetEndpoint(fmt.Sprintf("https://%s:5665/v1", icingaHost))
	for {
		if icingaClient.Check().Get(nil).Do().Status == 200 {
			e2e.PrintSeparately("Connected to icinga api")
			break
		}
		fmt.Println("Waiting for icinga to start")
		time.Sleep(5 * time.Second)
	}
})

var _ = AfterSuite(func() {
	err := root.DeleteNamespace()
	Expect(err).NotTo(HaveOccurred())
	e2e.PrintSeparately("Deleted namespace")
})
