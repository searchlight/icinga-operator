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
	v1alpha1 "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterAlerts implements ClusterAlertInterface
type FakeClusterAlerts struct {
	Fake *FakeMonitoringV1alpha1
	ns   string
}

var clusteralertsResource = schema.GroupVersionResource{Group: "monitoring.appscode.com", Version: "v1alpha1", Resource: "clusteralerts"}

var clusteralertsKind = schema.GroupVersionKind{Group: "monitoring.appscode.com", Version: "v1alpha1", Kind: "ClusterAlert"}

// Get takes name of the clusterAlert, and returns the corresponding clusterAlert object, and an error if there is any.
func (c *FakeClusterAlerts) Get(name string, options v1.GetOptions) (result *v1alpha1.ClusterAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(clusteralertsResource, c.ns, name), &v1alpha1.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterAlert), err
}

// List takes label and field selectors, and returns the list of ClusterAlerts that match those selectors.
func (c *FakeClusterAlerts) List(opts v1.ListOptions) (result *v1alpha1.ClusterAlertList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(clusteralertsResource, clusteralertsKind, c.ns, opts), &v1alpha1.ClusterAlertList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ClusterAlertList{}
	for _, item := range obj.(*v1alpha1.ClusterAlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterAlerts.
func (c *FakeClusterAlerts) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(clusteralertsResource, c.ns, opts))

}

// Create takes the representation of a clusterAlert and creates it.  Returns the server's representation of the clusterAlert, and an error, if there is any.
func (c *FakeClusterAlerts) Create(clusterAlert *v1alpha1.ClusterAlert) (result *v1alpha1.ClusterAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(clusteralertsResource, c.ns, clusterAlert), &v1alpha1.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterAlert), err
}

// Update takes the representation of a clusterAlert and updates it. Returns the server's representation of the clusterAlert, and an error, if there is any.
func (c *FakeClusterAlerts) Update(clusterAlert *v1alpha1.ClusterAlert) (result *v1alpha1.ClusterAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(clusteralertsResource, c.ns, clusterAlert), &v1alpha1.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterAlert), err
}

// Delete takes name of the clusterAlert and deletes it. Returns an error if one occurs.
func (c *FakeClusterAlerts) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(clusteralertsResource, c.ns, name), &v1alpha1.ClusterAlert{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterAlerts) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(clusteralertsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ClusterAlertList{})
	return err
}

// Patch applies the patch and returns the patched clusterAlert.
func (c *FakeClusterAlerts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ClusterAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(clusteralertsResource, c.ns, name, data, subresources...), &v1alpha1.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterAlert), err
}
