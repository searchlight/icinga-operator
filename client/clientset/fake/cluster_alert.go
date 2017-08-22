package fake

import (
	tapi "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"
)

type FakeClusterAlert struct {
	Fake *testing.Fake
	ns   string
}

var resourceClusterAlert = tapi.V1alpha1SchemeGroupVersion.WithResource(tapi.ResourceTypeClusterAlert)

// Get returns the ClusterAlert by name.
func (mock *FakeClusterAlert) Get(name string) (*tapi.ClusterAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewGetAction(resourceClusterAlert, mock.ns, name), &tapi.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.ClusterAlert), err
}

// List returns the a of ClusterAlerts.
func (mock *FakeClusterAlert) List(opts metav1.ListOptions) (*tapi.ClusterAlertList, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewListAction(resourceClusterAlert, mock.ns, opts), &tapi.ClusterAlert{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &tapi.ClusterAlertList{}
	for _, item := range obj.(*tapi.ClusterAlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Create creates a new ClusterAlert.
func (mock *FakeClusterAlert) Create(svc *tapi.ClusterAlert) (*tapi.ClusterAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewCreateAction(resourceClusterAlert, mock.ns, svc), &tapi.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.ClusterAlert), err
}

// Update updates a ClusterAlert.
func (mock *FakeClusterAlert) Update(svc *tapi.ClusterAlert) (*tapi.ClusterAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateAction(resourceClusterAlert, mock.ns, svc), &tapi.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.ClusterAlert), err
}

// Delete deletes a ClusterAlert by name.
func (mock *FakeClusterAlert) Delete(name string) error {
	_, err := mock.Fake.
		Invokes(testing.NewDeleteAction(resourceClusterAlert, mock.ns, name), &tapi.ClusterAlert{})

	return err
}

func (mock *FakeClusterAlert) UpdateStatus(srv *tapi.ClusterAlert) (*tapi.ClusterAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateSubresourceAction(resourceClusterAlert, "status", mock.ns, srv), &tapi.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.ClusterAlert), err
}

func (mock *FakeClusterAlert) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return mock.Fake.
		InvokesWatch(testing.NewWatchAction(resourceClusterAlert, mock.ns, opts))
}

func (mock *FakeClusterAlert) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*tapi.ClusterAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewPatchSubresourceAction(resourceClusterAlert, mock.ns, name, data, subresources...), &tapi.ClusterAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.ClusterAlert), err
}
