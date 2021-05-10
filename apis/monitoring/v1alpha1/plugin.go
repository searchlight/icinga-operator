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

package v1alpha1

import (
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/meta"
)

const (
	ResourceKindSearchlightPlugin     = "SearchlightPlugin"
	ResourcePluralSearchlightPlugin   = "searchlightplugins"
	ResourceSingularSearchlightPlugin = "searchlightplugin"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=searchlightplugins,singular=searchlightplugin,scope=Cluster,categories={monitoring,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Command",type="string",JSONPath=".spec.command"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type SearchlightPlugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec is the desired state of the SearchlightPlugin.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec SearchlightPluginSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// SearchlightPluginSpec describes the SearchlightPlugin the user wishes to create.
type SearchlightPluginSpec struct {
	// Check Command
	Command string `json:"command,omitempty" protobuf:"bytes,1,opt,name=command"`

	// Webhook provides a reference to the service for this SearchlightPlugin.
	// It must communicate on port 80
	Webhook *WebhookServiceSpec `json:"webhook,omitempty" protobuf:"bytes,2,opt,name=webhook"`

	// AlertKinds refers to supports Alert kinds for this plugin
	AlertKinds []string `json:"alertKinds" protobuf:"bytes,3,rep,name=alertKinds"`
	// Supported arguments for SearchlightPlugin
	Arguments PluginArguments `json:"arguments,omitempty" protobuf:"bytes,4,opt,name=arguments"`
	// Supported Icinga Service State
	States []string `json:"states" protobuf:"bytes,5,rep,name=states"`
}

type WebhookServiceSpec struct {
	// Namespace is the namespace of the service
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,1,opt,name=namespace"`
	// Name is the name of the service
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`
}

type VarType string

const (
	VarTypeInteger  VarType = "integer"
	VarTypeNumber   VarType = "number"
	VarTypeBoolean  VarType = "boolean"
	VarTypeString   VarType = "string"
	VarTypeDuration VarType = "duration"
)

type PluginVarField struct {
	Description string  `json:"description,omitempty" protobuf:"bytes,1,opt,name=description"`
	Type        VarType `json:"type" protobuf:"bytes,2,opt,name=type,casttype=VarType"`
}

type PluginVars struct {
	Fields   map[string]PluginVarField `json:"fields" protobuf:"bytes,1,rep,name=fields"`
	Required []string                  `json:"required,omitempty" protobuf:"bytes,2,rep,name=required"`
}

type PluginArguments struct {
	Vars *PluginVars       `json:"vars,omitempty" protobuf:"bytes,1,opt,name=vars"`
	Host map[string]string `json:"host,omitempty" protobuf:"bytes,2,rep,name=host"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SearchlightPluginList is a collection of SearchlightPlugin.
type SearchlightPluginList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of SearchlightPlugin.
	Items []SearchlightPlugin `json:"items" protobuf:"bytes,2,rep,name=items"`
}

var (
	validateVarValue = map[VarType]meta.ParserFunc{}
)

func registerVarValueParser(key VarType, fn meta.ParserFunc) {
	validateVarValue[key] = fn
}

func init() {
	registerVarValueParser(VarTypeInteger, meta.GetInt)
	registerVarValueParser(VarTypeNumber, meta.GetFloat)
	registerVarValueParser(VarTypeBoolean, meta.GetBool)
	registerVarValueParser(VarTypeString, meta.GetString)
	registerVarValueParser(VarTypeDuration, meta.GetDuration)
}

func validateVariables(pluginVars *PluginVars, vars map[string]string) error {
	if pluginVars == nil {
		return nil
	}
	// Check if any invalid variable is provided
	var err error
	for k := range vars {
		p, found := pluginVars.Fields[k]
		if !found {
			return fmt.Errorf("var '%s' is unsupported", k)
		}

		fn, found := validateVarValue[p.Type]
		if !found {
			return errors.Errorf(`type "%v" is not registered`, p.Type)
		}
		if _, err = fn(vars, k); err != nil {
			return errors.Wrapf(err, `validation failure: variable "%s" must be of type %v`, k, p.Type)
		}
	}
	for _, k := range pluginVars.Required {
		if _, ok := vars[k]; !ok {
			return fmt.Errorf("plugin variable '%s' is required", k)
		}
	}

	return nil
}
