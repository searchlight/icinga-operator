package fake

import (
	aci "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"
)

type FakePodAlert struct {
	Fake *testing.Fake
	ns   string
}

var resourcePodAlert = schema.GroupVersionResource{Group: "monitoring.appscode.com", Version: "v1alpha1", Resource: "podalerts"}

// Get returns the PodAlert by name.
func (mock *FakePodAlert) Get(name string) (*aci.PodAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewGetAction(resourcePodAlert, mock.ns, name), &aci.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.PodAlert), err
}

// List returns the a of PodAlerts.
func (mock *FakePodAlert) List(opts metav1.ListOptions) (*aci.PodAlertList, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewListAction(resourcePodAlert, mock.ns, opts), &aci.PodAlert{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &aci.PodAlertList{}
	for _, item := range obj.(*aci.PodAlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Create creates a new PodAlert.
func (mock *FakePodAlert) Create(svc *aci.PodAlert) (*aci.PodAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewCreateAction(resourcePodAlert, mock.ns, svc), &aci.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.PodAlert), err
}

// Update updates a PodAlert.
func (mock *FakePodAlert) Update(svc *aci.PodAlert) (*aci.PodAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateAction(resourcePodAlert, mock.ns, svc), &aci.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.PodAlert), err
}

// Delete deletes a PodAlert by name.
func (mock *FakePodAlert) Delete(name string) error {
	_, err := mock.Fake.
		Invokes(testing.NewDeleteAction(resourcePodAlert, mock.ns, name), &aci.PodAlert{})

	return err
}

func (mock *FakePodAlert) UpdateStatus(srv *aci.PodAlert) (*aci.PodAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateSubresourceAction(resourcePodAlert, "status", mock.ns, srv), &aci.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aci.PodAlert), err
}

func (mock *FakePodAlert) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return mock.Fake.
		InvokesWatch(testing.NewWatchAction(resourcePodAlert, mock.ns, opts))
}
