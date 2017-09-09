package install

import (
	tapi "github.com/appscode/searchlight/apis/monitoring"
	"github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	"k8s.io/apimachinery/pkg/runtime"
)

// Install registers the API group and adds types to a scheme
func Install(groupFactoryRegistry announced.APIGroupFactoryRegistry, registry *registered.APIRegistrationManager, scheme *runtime.Scheme) {
	if err := announced.NewGroupMetaFactory(
		&announced.GroupMetaFactoryArgs{
			GroupName:                  tapi.GroupName,
			VersionPreferenceOrder:     []string{v1alpha1.SchemeGroupVersion.Version},
			ImportPrefix:               "github.com/appscode/searchlight/apis/monitoring",
			AddInternalObjectsToScheme: tapi.AddToScheme,
		},
		announced.VersionToSchemeFunc{
			v1alpha1.SchemeGroupVersion.Version: v1alpha1.AddToScheme,
		},
	).Announce(groupFactoryRegistry).RegisterAndEnable(registry, scheme); err != nil {
		panic(err)
	}
}
