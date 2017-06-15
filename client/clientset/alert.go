package clientset

import (
	aci "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type AlertNamespacer interface {
	Alert(namespace string) AlertInterface
}

type AlertInterface interface {
	List(opts metav1.ListOptions) (*aci.AlertList, error)
	Get(name string) (*aci.Alert, error)
	Create(Alert *aci.Alert) (*aci.Alert, error)
	Update(Alert *aci.Alert) (*aci.Alert, error)
	Delete(name string) error
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	UpdateStatus(Alert *aci.Alert) (*aci.Alert, error)
}

type AlertImpl struct {
	r  rest.Interface
	ns string
}

var _ AlertInterface = &AlertImpl{}

func newAlert(c *ExtensionClient, namespace string) *AlertImpl {
	return &AlertImpl{c.restClient, namespace}
}

func (c *AlertImpl) List(opts metav1.ListOptions) (result *aci.AlertList, err error) {
	result = &aci.AlertList{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource("alerts").
		VersionedParams(&opts, ExtendedCodec).
		Do().
		Into(result)
	return
}

func (c *AlertImpl) Get(name string) (result *aci.Alert, err error) {
	result = &aci.Alert{}
	err = c.r.Get().
		Namespace(c.ns).
		Resource("alerts").
		Name(name).
		Do().
		Into(result)
	return
}

func (c *AlertImpl) Create(alert *aci.Alert) (result *aci.Alert, err error) {
	result = &aci.Alert{}
	err = c.r.Post().
		Namespace(c.ns).
		Resource("alerts").
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *AlertImpl) Update(alert *aci.Alert) (result *aci.Alert, err error) {
	result = &aci.Alert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource("alerts").
		Name(alert.Name).
		Body(alert).
		Do().
		Into(result)
	return
}

func (c *AlertImpl) Delete(name string) (err error) {
	return c.r.Delete().
		Namespace(c.ns).
		Resource("alerts").
		Name(name).
		Do().
		Error()
}

func (c *AlertImpl) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return c.r.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("alerts").
		VersionedParams(&opts, ExtendedCodec).
		Watch()
}

func (c *AlertImpl) UpdateStatus(alert *aci.Alert) (result *aci.Alert, err error) {
	result = &aci.Alert{}
	err = c.r.Put().
		Namespace(c.ns).
		Resource("alerts").
		Name(alert.Name).
		SubResource("status").
		Body(alert).
		Do().
		Into(result)
	return
}
