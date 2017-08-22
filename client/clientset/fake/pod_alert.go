package fake

import (
	tapi "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/testing"
)

type FakePodAlert struct {
	Fake *testing.Fake
	ns   string
}

var resourcePodAlert = tapi.V1alpha1SchemeGroupVersion.WithResource(tapi.ResourceTypePodAlert)

// Get returns the PodAlert by name.
func (mock *FakePodAlert) Get(name string) (*tapi.PodAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewGetAction(resourcePodAlert, mock.ns, name), &tapi.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.PodAlert), err
}

// List returns the a of PodAlerts.
func (mock *FakePodAlert) List(opts metav1.ListOptions) (*tapi.PodAlertList, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewListAction(resourcePodAlert, mock.ns, opts), &tapi.PodAlert{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &tapi.PodAlertList{}
	for _, item := range obj.(*tapi.PodAlertList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Create creates a new PodAlert.
func (mock *FakePodAlert) Create(svc *tapi.PodAlert) (*tapi.PodAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewCreateAction(resourcePodAlert, mock.ns, svc), &tapi.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.PodAlert), err
}

// Update updates a PodAlert.
func (mock *FakePodAlert) Update(svc *tapi.PodAlert) (*tapi.PodAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateAction(resourcePodAlert, mock.ns, svc), &tapi.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.PodAlert), err
}

// Delete deletes a PodAlert by name.
func (mock *FakePodAlert) Delete(name string) error {
	_, err := mock.Fake.
		Invokes(testing.NewDeleteAction(resourcePodAlert, mock.ns, name), &tapi.PodAlert{})

	return err
}

func (mock *FakePodAlert) UpdateStatus(srv *tapi.PodAlert) (*tapi.PodAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewUpdateSubresourceAction(resourcePodAlert, "status", mock.ns, srv), &tapi.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.PodAlert), err
}

func (mock *FakePodAlert) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	return mock.Fake.
		InvokesWatch(testing.NewWatchAction(resourcePodAlert, mock.ns, opts))
}

func (mock *FakePodAlert) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (*tapi.PodAlert, error) {
	obj, err := mock.Fake.
		Invokes(testing.NewPatchSubresourceAction(resourcePodAlert, mock.ns, name, data, subresources...), &tapi.PodAlert{})

	if obj == nil {
		return nil, err
	}
	return obj.(*tapi.PodAlert), err
}
