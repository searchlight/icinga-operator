package plugin

import (
	"encoding/json"

	hookapi "github.com/appscode/kutil/admission/api"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	admission "k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"
)

type AdmissionHook struct {
}

func (a *AdmissionHook) Resource() (plural schema.GroupVersionResource, singular string) {
	return schema.GroupVersionResource{
			Group:    "admission.monitoring.appscode.com",
			Version:  "v1alpha1",
			Resource: "admissionreviews",
		},
		"admissionreview"
}

func (a *AdmissionHook) Admit(req *admission.AdmissionRequest) *admission.AdmissionResponse {
	status := &admission.AdmissionResponse{}
	supportedKinds := sets.NewString(api.ResourceKindClusterAlert, api.ResourceKindNodeAlert, api.ResourceKindPodAlert)

	if (req.Operation != admission.Create && req.Operation != admission.Update) ||
		len(req.SubResource) != 0 ||
		req.Kind.Group != api.SchemeGroupVersion.Group ||
		!supportedKinds.Has(req.Kind.Kind) {
		status.Allowed = true
		return status
	}

	switch req.Kind.Kind {
	case api.ResourceKindClusterAlert:
		obj := &api.ClusterAlert{}
		err := json.Unmarshal(req.Object.Raw, obj)
		if err != nil {
			return hookapi.StatusBadRequest(err)
		}
		_, err = obj.IsValid()
		if err != nil {
			return hookapi.StatusForbidden(err)
		}
	case api.ResourceKindNodeAlert:
		obj := &api.NodeAlert{}
		err := json.Unmarshal(req.Object.Raw, obj)
		if err != nil {
			return hookapi.StatusBadRequest(err)
		}
		_, err = obj.IsValid()
		if err != nil {
			return hookapi.StatusForbidden(err)
		}
	case api.ResourceKindPodAlert:
		obj := &api.PodAlert{}
		err := json.Unmarshal(req.Object.Raw, obj)
		if err != nil {
			return hookapi.StatusBadRequest(err)
		}
		_, err = obj.IsValid()
		if err != nil {
			return hookapi.StatusForbidden(err)
		}
	}

	status.Allowed = true
	return status
}

func (a *AdmissionHook) Initialize(cfg *rest.Config, stopCh <-chan struct{}) error {
	return nil
}
