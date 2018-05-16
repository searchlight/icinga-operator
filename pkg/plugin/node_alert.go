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
			States: []string{stateOK, stateCritical, stateUnknown},
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
				Vars: &api.PluginVars{
					Items: map[string]api.PluginVarItem{
						"mountPoint": {
							Type: api.VarTypeString,
						},
						"secretName": {
							Type: api.VarTypeString,
						},
						"warning": {
							Type: api.VarTypeNumber,
						},
						"critical": {
							Type: api.VarTypeNumber,
						},
					},
					Required: []string{"mountPoint"},
				},
				Host: map[string]string{
					"host": "name",
					"v":    "vars.verbosity",
				},
			},
			States: []string{stateOK, stateCritical, stateUnknown},
		},
	}
}
