package openapi

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
)

type ImmutableStorage struct {
	cfg ResourceInfo
}

var _ rest.GroupVersionKindProvider = &ImmutableStorage{}
var _ rest.Creater = &ImmutableStorage{}
var _ rest.GracefulDeleter = &ImmutableStorage{}

func NewImmutableStorage(cfg ResourceInfo) *ImmutableStorage {
	return &ImmutableStorage{cfg}
}

func (r *ImmutableStorage) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return r.cfg.gvk
}

// Getter
func (r *ImmutableStorage) New() runtime.Object {
	return r.cfg.obj
}

func (r *ImmutableStorage) Create(ctx apirequest.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, includeUninitialized bool) (runtime.Object, error) {
	return r.New(), nil
}

// GracefulDeleter
func (r *ImmutableStorage) Delete(ctx apirequest.Context, name string, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	return r.New(), true, nil
}
