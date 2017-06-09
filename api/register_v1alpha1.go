package api

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// SchemeGroupVersion is group version used to register these objects
var V1alpha1SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}

var (
	v1alpha1SchemeBuilder = runtime.NewSchemeBuilder(v1addKnownTypes)
	V1betaAddToScheme     = v1alpha1SchemeBuilder.AddToScheme
)

// Adds the list of known types to api.Scheme.
func v1addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(V1alpha1SchemeGroupVersion,
		&Alert{},
		&AlertList{},

		&metav1.ListOptions{},
	)
	metav1.AddToGroupVersion(scheme, V1alpha1SchemeGroupVersion)
	return nil
}
