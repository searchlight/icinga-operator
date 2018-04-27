/*
Copyright 2018 The Searchlight Authors.

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

// FakeSearchlightPlugins implements SearchlightPluginInterface
type FakeSearchlightPlugins struct {
	Fake *FakeMonitoringV1alpha1
}

var searchlightpluginsResource = schema.GroupVersionResource{Group: "monitoring.appscode.com", Version: "v1alpha1", Resource: "searchlightplugins"}

var searchlightpluginsKind = schema.GroupVersionKind{Group: "monitoring.appscode.com", Version: "v1alpha1", Kind: "SearchlightPlugin"}

// Get takes name of the searchlightPlugin, and returns the corresponding searchlightPlugin object, and an error if there is any.
func (c *FakeSearchlightPlugins) Get(name string, options v1.GetOptions) (result *v1alpha1.SearchlightPlugin, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(searchlightpluginsResource, name), &v1alpha1.SearchlightPlugin{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SearchlightPlugin), err
}

// List takes label and field selectors, and returns the list of SearchlightPlugins that match those selectors.
func (c *FakeSearchlightPlugins) List(opts v1.ListOptions) (result *v1alpha1.SearchlightPluginList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(searchlightpluginsResource, searchlightpluginsKind, opts), &v1alpha1.SearchlightPluginList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SearchlightPluginList{}
	for _, item := range obj.(*v1alpha1.SearchlightPluginList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested searchlightPlugins.
func (c *FakeSearchlightPlugins) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(searchlightpluginsResource, opts))
}

// Create takes the representation of a searchlightPlugin and creates it.  Returns the server's representation of the searchlightPlugin, and an error, if there is any.
func (c *FakeSearchlightPlugins) Create(searchlightPlugin *v1alpha1.SearchlightPlugin) (result *v1alpha1.SearchlightPlugin, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(searchlightpluginsResource, searchlightPlugin), &v1alpha1.SearchlightPlugin{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SearchlightPlugin), err
}

// Update takes the representation of a searchlightPlugin and updates it. Returns the server's representation of the searchlightPlugin, and an error, if there is any.
func (c *FakeSearchlightPlugins) Update(searchlightPlugin *v1alpha1.SearchlightPlugin) (result *v1alpha1.SearchlightPlugin, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(searchlightpluginsResource, searchlightPlugin), &v1alpha1.SearchlightPlugin{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SearchlightPlugin), err
}

// Delete takes name of the searchlightPlugin and deletes it. Returns an error if one occurs.
func (c *FakeSearchlightPlugins) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(searchlightpluginsResource, name), &v1alpha1.SearchlightPlugin{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSearchlightPlugins) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(searchlightpluginsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.SearchlightPluginList{})
	return err
}

// Patch applies the patch and returns the patched searchlightPlugin.
func (c *FakeSearchlightPlugins) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SearchlightPlugin, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(searchlightpluginsResource, name, data, subresources...), &v1alpha1.SearchlightPlugin{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SearchlightPlugin), err
}
