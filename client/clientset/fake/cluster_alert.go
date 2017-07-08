package fake

import (
	aci "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"
)

type FakeClusterAlert struct {
	Fake *testing.Fake
	ns   string
}

var resourceClusterAlert = schema.GroupVersionResource{Group: "monitoring.appscode.com", Version: "v1alpha1", Resource: "clusteralerts"}

// Get returns the ClusterAlert by name.
func (mock *FakeClusterAlert) Get(name string) (*aci.ClusterAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewGetAction(resourceClusterAlert, mock.ns, name), &aci.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.ClusterAlert), err
}

// List returns the a of ClusterAlerts.
func (mock *FakeClusterAlert) List(opts metav1.ListOptions) (*aci.ClusterAlertList, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewListAction(resourceClusterAlert, mock.ns, opts), &aci.ClusterAlert{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &aci.ClusterAlertList{}
	for _, item := range obj.(*aci.ClusterAlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Create creates a new ClusterAlert.
func (mock *FakeClusterAlert) Create(svc *aci.ClusterAlert) (*aci.ClusterAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewCreateAction(resourceClusterAlert, mock.ns, svc), &aci.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.ClusterAlert), err
}

// Update updates a ClusterAlert.
func (mock *FakeClusterAlert) Update(svc *aci.ClusterAlert) (*aci.ClusterAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateAction(resourceClusterAlert, mock.ns, svc), &aci.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.ClusterAlert), err
}

// Delete deletes a ClusterAlert by name.
func (mock *FakeClusterAlert) Delete(name string) error {
	_, err := mock.Fake.
		Invokes(testing.NewDeleteAction(resourceClusterAlert, mock.ns, name), &aci.ClusterAlert{})

	return err
}

func (mock *FakeClusterAlert) UpdateStatus(srv *aci.ClusterAlert) (*aci.ClusterAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateSubresourceAction(resourceClusterAlert, "status", mock.ns, srv), &aci.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.ClusterAlert), err
}

func (mock *FakeClusterAlert) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return mock.Fake.
		InvokesWatch(testing.NewWatchAction(resourceClusterAlert, mock.ns, opts))
}
