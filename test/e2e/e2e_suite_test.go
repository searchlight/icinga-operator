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
	//. "github.com/appscode/searchlight/test/e2e/matcher"
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

	//// Create namespace
	//err = root.CreateNamespace()
	//Expect(err).NotTo(HaveOccurred())
	//
	//// Create Searchlight deployment
	//slDeployment := root.Invoke().DeploymentExtensionSearchlight()
	//err = root.CreateDeploymentExtension(slDeployment)
	//Expect(err).NotTo(HaveOccurred())
	//By("Waiting for Running pods")
	//root.EventuallyDeploymentExtension(slDeployment.ObjectMeta).Should(HaveRunningPods(*slDeployment.Spec.Replicas))
	//
	//// Create Searchlight service
	//slService := root.Invoke().ServiceSearchlight()
	//err = root.CreateService(slService)
	//Expect(err).NotTo(HaveOccurred())
	//// Get Icinga Ingress Hostname
	//icingaHost, err := root.GetServiceIngressHost(slService.ObjectMeta)
	//Expect(err).NotTo(HaveOccurred())

	icingaHost := "a365ddb756b8311e78b6912f236046fb-578643977.us-east-1.elb.amazonaws.com"

	// Icinga Config
	cfg := &icinga.Config{
		Endpoint: fmt.Sprintf("https://%s:5665/v1", icingaHost),
	}
	cfg.BasicAuth.Username = e2e.ICINGA_API_USER
	cfg.BasicAuth.Password = e2e.ICINGA_API_PASSWORD

	// Icinga Client
	icingaClient := icinga.NewClient(*cfg)
	root = root.SetIcingaClient(icingaClient)
	root.EventuallyIcingaAPI().Should(Succeed())

	fmt.Println()
	fmt.Println("Icingaweb2:     ", fmt.Sprintf("http://%s/icingaweb2", icingaHost))
	fmt.Println("Login password: ", e2e.ICINGA_WEB_ADMIN_PASSWORD)
	fmt.Println()

	// Controller
	ctrl = controller.New(kubeClient, extClient, icingaClient)
	err = ctrl.Setup()
	Expect(err).NotTo(HaveOccurred())
	ctrl.Run()
	root.EventuallyTPR().Should(Succeed())
})

var _ = AfterSuite(func() {
	 //err := root.DeleteNamespace()
	 //Expect(err).NotTo(HaveOccurred())
	 //e2e.PrintSeparately("Deleted namespace")
})
