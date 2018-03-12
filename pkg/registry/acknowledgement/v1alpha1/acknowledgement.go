package v1beta1

import (
	incident "github.com/appscode/searchlight/apis/incidents/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

// Adapted from https://github.com/openshift/generic-admission-server/blob/master/pkg/registry/admissionreview/admission_review.go

type AdmissionHookFunc func(admissionSpec *incident.Acknowledgement) *incident.Acknowledgement

type REST struct {
	hookFn AdmissionHookFunc
}

var _ rest.Creater = &REST{}
var _ rest.GroupVersionKindProvider = &REST{}

func NewREST(hookFn AdmissionHookFunc) *REST {
	return &REST{
		hookFn: hookFn,
	}
}

func (r *REST) New() runtime.Object {
	return &incident.Acknowledgement{}
}

func (r *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return incident.SchemeGroupVersion.WithKind("Acknowledgement")
}

func (r *REST) Create(ctx apirequest.Context, obj runtime.Object, _ rest.ValidateObjectFunc, _ bool) (runtime.Object, error) {
	admissionReview := obj.(*incident.Acknowledgement)
	return admissionReview, nil
}
