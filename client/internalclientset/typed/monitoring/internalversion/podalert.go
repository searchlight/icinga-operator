/*
Copyright 2017 The Searchlight Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package internalversion

import (
	monitoring "github.com/appscode/searchlight/apis/monitoring"
	scheme "github.com/appscode/searchlight/client/internalclientset/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PodAlertsGetter has a method to return a PodAlertInterface.
// A group's client should implement this interface.
type PodAlertsGetter interface {
	PodAlerts(namespace string) PodAlertInterface
}

// PodAlertInterface has methods to work with PodAlert resources.
type PodAlertInterface interface {
	Create(*monitoring.PodAlert) (*monitoring.PodAlert, error)
	Update(*monitoring.PodAlert) (*monitoring.PodAlert, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*monitoring.PodAlert, error)
	List(opts v1.ListOptions) (*monitoring.PodAlertList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *monitoring.PodAlert, err error)
	PodAlertExpansion
}

// podAlerts implements PodAlertInterface
type podAlerts struct {
	client rest.Interface
	ns     string
}

// newPodAlerts returns a PodAlerts
func newPodAlerts(c *MonitoringClient, namespace string) *podAlerts {
	return &podAlerts{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the podAlert, and returns the corresponding podAlert object, and an error if there is any.
func (c *podAlerts) Get(name string, options v1.GetOptions) (result *monitoring.PodAlert, err error) {
	result = &monitoring.PodAlert{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podalerts").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PodAlerts that match those selectors.
func (c *podAlerts) List(opts v1.ListOptions) (result *monitoring.PodAlertList, err error) {
	result = &monitoring.PodAlertList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("podalerts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested podAlerts.
func (c *podAlerts) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("podalerts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a podAlert and creates it.  Returns the server's representation of the podAlert, and an error, if there is any.
func (c *podAlerts) Create(podAlert *monitoring.PodAlert) (result *monitoring.PodAlert, err error) {
	result = &monitoring.PodAlert{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("podalerts").
		Body(podAlert).
		Do().
		Into(result)
	return
}

// Update takes the representation of a podAlert and updates it. Returns the server's representation of the podAlert, and an error, if there is any.
func (c *podAlerts) Update(podAlert *monitoring.PodAlert) (result *monitoring.PodAlert, err error) {
	result = &monitoring.PodAlert{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("podalerts").
		Name(podAlert.Name).
		Body(podAlert).
		Do().
		Into(result)
	return
}

// Delete takes name of the podAlert and deletes it. Returns an error if one occurs.
func (c *podAlerts) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podalerts").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *podAlerts) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("podalerts").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched podAlert.
func (c *podAlerts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *monitoring.PodAlert, err error) {
	result = &monitoring.PodAlert{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("podalerts").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
