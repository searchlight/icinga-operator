package fake

import (
	tapi "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"
)

type FakeNodeAlert struct {
	Fake *testing.Fake
	ns   string
}

var resourceNodeAlert = tapi.V1alpha1SchemeGroupVersion.WithResource(tapi.ResourceTypeNodeAlert)
var kindNodeAlert = tapi.V1alpha1SchemeGroupVersion.WithKind(tapi.ResourceKindNodeAlert)

// Get returns the NodeAlert by name.
func (mock *FakeNodeAlert) Get(name string) (*tapi.NodeAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewGetAction(resourceNodeAlert, mock.ns, name), &tapi.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.NodeAlert), err
}

// List returns the a of NodeAlerts.
func (mock *FakeNodeAlert) List(opts metav1.ListOptions) (*tapi.NodeAlertList, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewListAction(resourceNodeAlert, kindNodeAlert, mock.ns, opts), &tapi.NodeAlert{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &tapi.NodeAlertList{}
	for _, item := range obj.(*tapi.NodeAlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Create creates a new NodeAlert.
func (mock *FakeNodeAlert) Create(svc *tapi.NodeAlert) (*tapi.NodeAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewCreateAction(resourceNodeAlert, mock.ns, svc), &tapi.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.NodeAlert), err
}

// Update updates a NodeAlert.
func (mock *FakeNodeAlert) Update(svc *tapi.NodeAlert) (*tapi.NodeAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateAction(resourceNodeAlert, mock.ns, svc), &tapi.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.NodeAlert), err
}

// Delete deletes a NodeAlert by name.
func (mock *FakeNodeAlert) Delete(name string) error {
	_, err := mock.Fake.
		Invokes(testing.NewDeleteAction(resourceNodeAlert, mock.ns, name), &tapi.NodeAlert{})

	return err
}

func (mock *FakeNodeAlert) UpdateStatus(srv *tapi.NodeAlert) (*tapi.NodeAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateSubresourceAction(resourceNodeAlert, "status", mock.ns, srv), &tapi.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.NodeAlert), err
}

func (mock *FakeNodeAlert) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return mock.Fake.
		InvokesWatch(testing.NewWatchAction(resourceNodeAlert, mock.ns, opts))
}

func (mock *FakeNodeAlert) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*tapi.NodeAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewPatchSubresourceAction(resourceNodeAlert, mock.ns, name, data, subresources...), &tapi.NodeAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.NodeAlert), err
}
