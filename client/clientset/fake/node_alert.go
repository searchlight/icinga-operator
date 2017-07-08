package fake

import (
	aci "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"
)

type FakeNodeAlert struct {
	Fake *testing.Fake
	ns   string
}

var resourceNodeAlert = schema.GroupVersionResource{Group: "monitoring.appscode.com", Version: "v1alpha1", Resource: "nodealerts"}

// Get returns the NodeAlert by name.
func (mock *FakeNodeAlert) Get(name string) (*aci.NodeAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewGetAction(resourceNodeAlert, mock.ns, name), &aci.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.NodeAlert), err
}

// List returns the a of NodeAlerts.
func (mock *FakeNodeAlert) List(opts metav1.ListOptions) (*aci.NodeAlertList, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewListAction(resourceNodeAlert, mock.ns, opts), &aci.NodeAlert{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &aci.NodeAlertList{}
	for _, item := range obj.(*aci.NodeAlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Create creates a new NodeAlert.
func (mock *FakeNodeAlert) Create(svc *aci.NodeAlert) (*aci.NodeAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewCreateAction(resourceNodeAlert, mock.ns, svc), &aci.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.NodeAlert), err
}

// Update updates a NodeAlert.
func (mock *FakeNodeAlert) Update(svc *aci.NodeAlert) (*aci.NodeAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateAction(resourceNodeAlert, mock.ns, svc), &aci.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.NodeAlert), err
}

// Delete deletes a NodeAlert by name.
func (mock *FakeNodeAlert) Delete(name string) error {
	_, err := mock.Fake.
		Invokes(testing.NewDeleteAction(resourceNodeAlert, mock.ns, name), &aci.NodeAlert{})

	return err
}

func (mock *FakeNodeAlert) UpdateStatus(srv *aci.NodeAlert) (*aci.NodeAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateSubresourceAction(resourceNodeAlert, "status", mock.ns, srv), &aci.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.NodeAlert), err
}

func (mock *FakeNodeAlert) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return mock.Fake.
		InvokesWatch(testing.NewWatchAction(resourceNodeAlert, mock.ns, opts))
}
