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

func GetComponentStatusPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "component-status",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_component_status",
			AlertKinds: []string{api.ResourceKindClusterAlert},
			Arguments: api.PluginArguments{
				Vars: &api.PluginVars{
					Fields: map[string]api.PluginVarField{
						"selector": {
							Type: api.VarTypeString,
						},
						"componentName": {
							Type: api.VarTypeString,
						},
					},
				},
				Host: map[string]string{
					"v": "vars.verbosity",
				},
			},
			States: []string{stateOK, stateCritical, stateUnknown},
		},
	}
}

func GetJsonPathPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "json-path",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_json_path",
			AlertKinds: []string{api.ResourceKindClusterAlert},
			Arguments: api.PluginArguments{
				Vars: &api.PluginVars{
					Fields: map[string]api.PluginVarField{
						"url": {
							Type: api.VarTypeString,
						},
						"secretName": {
							Type: api.VarTypeString,
						},
						"warning": {
							Type: api.VarTypeString,
						},
						"critical": {
							Type: api.VarTypeString,
						},
					},
					Required: []string{"url"},
				},
				Host: map[string]string{
					"host": "name",
					"v":    "vars.verbosity",
				},
			},
			States: []string{stateOK, stateWarning, stateCritical, stateUnknown},
		},
	}
}

func GetNodeExistsPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "node-exists",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_node_exists",
			AlertKinds: []string{api.ResourceKindClusterAlert},
			Arguments: api.PluginArguments{
				Vars: &api.PluginVars{
					Fields: map[string]api.PluginVarField{
						"selector": {
							Type: api.VarTypeString,
						},
						"nodeName": {
							Type: api.VarTypeString,
						},
						"count": {
							Type: api.VarTypeInteger,
						},
					},
				},
				Host: map[string]string{
					"v": "vars.verbosity",
				},
			},
			States: []string{stateOK, stateCritical, stateUnknown},
		},
	}
}

func GetPodExistsPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod-exists",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_pod_exists",
			AlertKinds: []string{api.ResourceKindClusterAlert},
			Arguments: api.PluginArguments{
				Vars: &api.PluginVars{
					Fields: map[string]api.PluginVarField{
						"selector": {
							Type: api.VarTypeString,
						},
						"podName": {
							Type: api.VarTypeString,
						},
						"count": {
							Type: api.VarTypeInteger,
						},
					},
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

func GetEventPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "event",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_event",
			AlertKinds: []string{api.ResourceKindClusterAlert},
			Arguments: api.PluginArguments{
				Vars: &api.PluginVars{
					Fields: map[string]api.PluginVarField{
						"clockSkew": {
							Type: api.VarTypeDuration,
						},
						"involvedObjectName": {
							Type: api.VarTypeString,
						},
						"involvedObjectNamespace": {
							Type: api.VarTypeString,
						},
						"involvedObjectKind": {
							Type: api.VarTypeString,
						},
						"involvedObjectUID": {
							Type: api.VarTypeString,
						},
					},
				},
				Host: map[string]string{
					"host": "name",
					"v":    "vars.verbosity",
				},
			},
			States: []string{stateOK, stateWarning, stateUnknown},
		},
	}
}

func GetCACertPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ca-cert",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_ca_cert",
			AlertKinds: []string{api.ResourceKindClusterAlert},
			Arguments: api.PluginArguments{
				Vars: &api.PluginVars{
					Fields: map[string]api.PluginVarField{
						"warning": {
							Type: api.VarTypeDuration,
						},
						"critical": {
							Type: api.VarTypeDuration,
						},
					},
				},
				Host: map[string]string{
					"v": "vars.verbosity",
				},
			},
			States: []string{stateOK, stateWarning, stateCritical, stateUnknown},
		},
	}
}

func GetCertPlugin() *api.SearchlightPlugin {
	return &api.SearchlightPlugin{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cert",
		},
		TypeMeta: PluginTypeMeta,
		Spec: api.SearchlightPluginSpec{
			Command:    "hyperalert check_cert",
			AlertKinds: []string{api.ResourceKindClusterAlert},
			Arguments: api.PluginArguments{
				Vars: &api.PluginVars{
					Fields: map[string]api.PluginVarField{
						"selector": {
							Type: api.VarTypeString,
						},
						"secretName": {
							Type: api.VarTypeString,
						},
						"secretKey": {
							Type: api.VarTypeString,
						},
						"warning": {
							Type: api.VarTypeDuration,
						},
						"critical": {
							Type: api.VarTypeDuration,
						},
					},
				},
				Host: map[string]string{
					"host": "name",
					"v":    "vars.verbosity",
				},
			},
			States: []string{stateOK, stateWarning, stateCritical, stateUnknown},
		},
	}
}
