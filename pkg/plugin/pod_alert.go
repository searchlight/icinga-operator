/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"

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
			States: []string{stateOK, stateCritical, stateUnknown},
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
				Vars: &api.PluginVars{
					Fields: map[string]api.PluginVarField{
						"volumeName": {
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
					Required: []string{"volumeName"},
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
				Vars: &api.PluginVars{
					Fields: map[string]api.PluginVarField{
						"container": {
							Type: api.VarTypeString,
						},
						"cmd": {
							Type: api.VarTypeString,
						},
						"argv": {
							Type: api.VarTypeString,
						},
					},
					Required: []string{"argv"},
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
