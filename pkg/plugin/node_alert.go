package plugin

import (
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNodeStatusPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-status",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_node_status",
			AlertKinds: []string{api.ResourceKindNodeAlert},
			Arguments: api.PluginArguments{
				Host: map[string]string{
					"host": "name",
					"v":    "vars.verbosity",
				},
			},
			State: []string{stateOK, stateCritical, stateUnknown},
		},
	}
}

func GetNodeVolumePlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-volume",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_volume",
			AlertKinds: []string{api.ResourceKindNodeAlert},
			Arguments: api.PluginArguments{
				Vars: []string{
					"mountPoint",
					"secretName",
					"warning",
					"critical",
				},
				Host: map[string]string{
					"host": "name",
					"v":    "vars.verbosity",
				},
			},
			State: []string{stateOK, stateCritical, stateUnknown},
		},
	}
}
