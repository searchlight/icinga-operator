package e2e

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	cs "github.com/appscode/searchlight/client/clientset/versioned"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/operator"
	"github.com/appscode/searchlight/test/e2e/framework"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"kmodules.xyz/client-go/logs"
)

var (
	provider           string
	storageClass       string
	searchlightService string
)

func init() {
	flag.StringVar(&provider, "provider", "minikube", "Kubernetes cloud provider")
	flag.StringVar(&storageClass, "storageclass", "", "Kubernetes StorageClass name")
	flag.StringVar(&searchlightService, "searchlight-service", "", "Running searchlight reference")
}

const (
	TIMEOUT = 20 * time.Minute
)

var (
	op   *operator.Operator
	root *framework.Framework
)

func TestE2e(t *testing.T) {
	logs.InitLogs()
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TIMEOUT)

	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "e2e Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {

	Expect(searchlightService).ShouldNot(BeEmpty())

	// Kubernetes config
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube/config")
	By("Using kubeconfig from " + kubeconfigPath)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	Expect(err).NotTo(HaveOccurred())
	// Clients
	Expect(err).NotTo(HaveOccurred())
	kubeClient := kubernetes.NewForConfigOrDie(config)
	apiExtKubeClient := crd_cs.NewForConfigOrDie(config)
	extClient := cs.NewForConfigOrDie(config)
	// Framework
	root = framework.New(kubeClient, apiExtKubeClient, extClient, nil, provider, storageClass)

	framework.PrintSeparately("Using namespace " + root.Namespace())

	// Create namespace
	err = root.CreateNamespace()
	Expect(err).NotTo(HaveOccurred())

	parts := strings.Split(searchlightService, "@")
	Expect(len(parts)).Should(BeIdenticalTo(2))
	om := metav1.ObjectMeta{
		Name:      parts[0],
		Namespace: parts[1],
	}
	slService := &core.Service{ObjectMeta: om}

	// Get Icinga Ingress Hostname
	endpoint, err := root.GetServiceEndpoint(slService.ObjectMeta, "icinga")
	Expect(err).NotTo(HaveOccurred())

	// Icinga Config
	cfg := &icinga.Config{
		Endpoint: fmt.Sprintf("https://%v/v1", endpoint),
		CACert:   nil,
	}

	cfg.BasicAuth.Username = ICINGA_API_USER
	cfg.BasicAuth.Password, err = root.Invoke().GetIcingaApiPassword(om)
	Expect(err).NotTo(HaveOccurred())

	// Icinga Client
	icingaClient := icinga.NewClient(*cfg)
	root = root.SetIcingaClient(icingaClient)
	root.EventuallyIcingaAPI().Should(Succeed())

	icingawebEndpoint, err := root.GetServiceEndpoint(slService.ObjectMeta, "ui")
	Expect(err).NotTo(HaveOccurred())
	fmt.Println()
	fmt.Println("Icingaweb2:     ", fmt.Sprintf("http://%v/", icingawebEndpoint))
	fmt.Println("Login password: ", ICINGA_WEB_UI_PASSWORD)
	fmt.Println()

})

var _ = AfterSuite(func() {
	root.CleanPodAlert()
	root.CleanNodeAlert()
	root.CleanClusterAlert()
	err := root.DeleteNamespace()
	Expect(err).NotTo(HaveOccurred())
	framework.PrintSeparately("Deleted namespace")
})
