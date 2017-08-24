package fake

import (
	"github.com/appscode/searchlight/client/clientset"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/testing"
)

type FakeExtensionClient struct {
	*testing.Fake
}

var _ clientset.ExtensionInterface = &FakeExtensionClient{}

func NewFakeExtensionClient(objects ...runtime.Object) *FakeExtensionClient {
	o := testing.NewObjectTracker(api.Scheme, api.Codecs.UniversalDecoder())
	for _, obj := range objects {
		if obj.GetObjectKind().GroupVersionKind().Group == "monitoring.appscode.com" {
			if err := o.Add(obj); err != nil {
				panic(err)
			}
		}
	}

	fakePtr := testing.Fake{}
	fakePtr.AddReactor("*", "*", testing.ObjectReaction(o))

	fakePtr.AddWatchReactor("*", testing.DefaultWatchReactor(watch.NewFake(), nil))

	return &FakeExtensionClient{&fakePtr}
}

func (c *FakeExtensionClient) PodAlerts(namespace string) clientset.PodAlertInterface {
	return &FakePodAlert{c.Fake, namespace}
}

func (c *FakeExtensionClient) NodeAlerts(namespace string) clientset.NodeAlertInterface {
	return &FakeNodeAlert{c.Fake, namespace}
}

func (c *FakeExtensionClient) ClusterAlerts(namespace string) clientset.ClusterAlertInterface {
	return &FakeClusterAlert{c.Fake, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeExtensionClient) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
