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

// FakeNodeAlerts implements NodeAlertInterface
type FakeNodeAlerts struct {
	Fake *FakeMonitoringV1alpha1
	ns   string
}

var nodealertsResource = schema.GroupVersionResource{Group: "monitoring.appscode.com", Version: "v1alpha1", Resource: "nodealerts"}

var nodealertsKind = schema.GroupVersionKind{Group: "monitoring.appscode.com", Version: "v1alpha1", Kind: "NodeAlert"}

func (c *FakeNodeAlerts) Create(nodeAlert *v1alpha1.NodeAlert) (result *v1alpha1.NodeAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(nodealertsResource, c.ns, nodeAlert), &v1alpha1.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NodeAlert), err
}

func (c *FakeNodeAlerts) Update(nodeAlert *v1alpha1.NodeAlert) (result *v1alpha1.NodeAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(nodealertsResource, c.ns, nodeAlert), &v1alpha1.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NodeAlert), err
}

func (c *FakeNodeAlerts) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(nodealertsResource, c.ns, name), &v1alpha1.NodeAlert{})

	return err
}

func (c *FakeNodeAlerts) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(nodealertsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.NodeAlertList{})
	return err
}

func (c *FakeNodeAlerts) Get(name string, options v1.GetOptions) (result *v1alpha1.NodeAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(nodealertsResource, c.ns, name), &v1alpha1.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NodeAlert), err
}

func (c *FakeNodeAlerts) List(opts v1.ListOptions) (result *v1alpha1.NodeAlertList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(nodealertsResource, nodealertsKind, c.ns, opts), &v1alpha1.NodeAlertList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.NodeAlertList{}
	for _, item := range obj.(*v1alpha1.NodeAlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested nodeAlerts.
func (c *FakeNodeAlerts) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(nodealertsResource, c.ns, opts))

}

// Patch applies the patch and returns the patched nodeAlert.
func (c *FakeNodeAlerts) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.NodeAlert, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(nodealertsResource, c.ns, name, data, subresources...), &v1alpha1.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.NodeAlert), err
}
