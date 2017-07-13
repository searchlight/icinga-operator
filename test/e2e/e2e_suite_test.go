package e2e_test

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	tcs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
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

	kubeconfigPath := filepath.Join(userHome, ".kube/config")
	By("Using kubeconfig from " + kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())

	kubeClient := clientset.NewForConfigOrDie(config)
	extClient := tcs.NewForConfigOrDie(config)
	root = framework.New(kubeClient, extClient, nil)

	fmt.Println("Using namespace ", root.Namespace())
	By("Using namespace " + root.Namespace())

	err = root.CreateNamespace()
	Expect(err).NotTo(HaveOccurred())

	configure := &icinga.Configurator{
		NotifierSecretName: "searchlight",
		ConfigRoot:         filepath.Join(userHome),
		Expiry:             10 * 365 * 24 * time.Hour,
	}
	_, err = configure.LoadIcingaConfig()
	Expect(err).NotTo(HaveOccurred())

	icingaSecret, err := root.Invoke().SecretSearchlight(filepath.Join(configure.ConfigRoot, "icinga2"))
	Expect(err).NotTo(HaveOccurred())
	err = root.CreateSecret(icingaSecret)
	Expect(err).NotTo(HaveOccurred())

	icingaDeployment := root.Invoke().DeploymentAppSearchlight()
	err = root.CreateDeploymentApp(icingaDeployment)
	Expect(err).NotTo(HaveOccurred())

})

var _ = AfterSuite(func() {
	fmt.Println("Delete Icinga Deployments")
})
