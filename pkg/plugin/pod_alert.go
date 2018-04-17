package plugin

import (
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetPodStatusPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod-status",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_pod_status",
			AlertKinds: []string{api.ResourceKindPodAlert},
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

func GetPodVolumePlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod-volume",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_volume",
			AlertKinds: []string{api.ResourceKindPodAlert},
			Arguments: api.PluginArguments{
				Vars: []string{
					"volumeName",
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

func GetPodExecPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod-exec",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_pod_exec",
			AlertKinds: []string{api.ResourceKindPodAlert},
			Arguments: api.PluginArguments{
				Vars: []string{
					"container",
					"cmd",
					"argv",
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
