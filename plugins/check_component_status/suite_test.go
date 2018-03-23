package check_component_status

import (
	"testing"
	"time"

	"github.com/appscode/searchlight/client/clientset/versioned/scheme"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"
	clientSetScheme "k8s.io/client-go/kubernetes/scheme"
)

var (
	cs *fake.Clientset
)

const (
	TIMEOUT = 2 * time.Minute
)

func TestPlugin_Check(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TIMEOUT)
	RunSpecsWithDefaultAndCustomReporters(t, "check_component_status Suite", []Reporter{})
}

var _ = BeforeSuite(func() {
	scheme.AddToScheme(clientSetScheme.Scheme)
	cs = fake.NewSimpleClientset()
})

var _ = AfterSuite(func() {})
