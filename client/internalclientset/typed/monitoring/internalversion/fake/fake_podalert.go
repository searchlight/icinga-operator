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

package fake

import (
	monitoring "github.com/appscode/searchlight/apis/monitoring"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePodAlerts implements PodAlertInterface
type FakePodAlerts struct {
	Fake *FakeMonitoring
	ns   string
}

var podalertsResource = schema.GroupVersionResource{Group: "monitoring.appscode.com", Version: "", Resource: "podalerts"}

var podalertsKind = schema.GroupVersionKind{Group: "monitoring.appscode.com", Version: "", Kind: "PodAlert"}

// Get takes name of the podAlert, and returns the corresponding podAlert object, and an error if there is any.
func (c *FakePodAlerts) Get(name string, options v1.GetOptions) (result *monitoring.PodAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(podalertsResource, c.ns, name), &monitoring.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*monitoring.PodAlert), err
}

// List takes label and field selectors, and returns the list of PodAlerts that match those selectors.
func (c *FakePodAlerts) List(opts v1.ListOptions) (result *monitoring.PodAlertList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(podalertsResource, podalertsKind, c.ns, opts), &monitoring.PodAlertList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &monitoring.PodAlertList{}
	for _, item := range obj.(*monitoring.PodAlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested podAlerts.
func (c *FakePodAlerts) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(podalertsResource, c.ns, opts))

}

// Create takes the representation of a podAlert and creates it.  Returns the server's representation of the podAlert, and an error, if there is any.
func (c *FakePodAlerts) Create(podAlert *monitoring.PodAlert) (result *monitoring.PodAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(podalertsResource, c.ns, podAlert), &monitoring.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*monitoring.PodAlert), err
}

// Update takes the representation of a podAlert and updates it. Returns the server's representation of the podAlert, and an error, if there is any.
func (c *FakePodAlerts) Update(podAlert *monitoring.PodAlert) (result *monitoring.PodAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(podalertsResource, c.ns, podAlert), &monitoring.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*monitoring.PodAlert), err
}

// Delete takes name of the podAlert and deletes it. Returns an error if one occurs.
func (c *FakePodAlerts) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(podalertsResource, c.ns, name), &monitoring.PodAlert{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePodAlerts) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(podalertsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &monitoring.PodAlertList{})
	return err
}

// Patch applies the patch and returns the patched podAlert.
func (c *FakePodAlerts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *monitoring.PodAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(podalertsResource, c.ns, name, data, subresources...), &monitoring.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*monitoring.PodAlert), err
}
