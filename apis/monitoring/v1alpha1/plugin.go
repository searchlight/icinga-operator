package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindSearchlightPlugin     = "SearchlightPlugin"
	ResourcePluralSearchlightPlugin   = "searchlightplugins"
	ResourceSingularSearchlightPlugin = "searchlightplugin"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SearchlightPlugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the SearchlightPlugin.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec SearchlightPluginSpec `json:"spec,omitempty"`
}

// SearchlightPluginSpec describes the SearchlightPlugin the user wishes to create.
type SearchlightPluginSpec struct {
	// Check Command
	Command string `json:"command,omitempty"`

	// Webhook provides a reference to the service for this SearchlightPlugin.
	// It must communicate on port 80
	Webhook *WebhookServiceSpec `json:"webhook,omitempty"`

	// AlertKinds refers to supports Alert kinds for this plugin
	AlertKinds []string `json:"alertKinds"`
	// Supported arguments for SearchlightPlugin
	Arguments PluginArguments `json:"arguments,omitempty"`
	// Supported Icinga Service State
	State []string `json:"state"`
}

type WebhookServiceSpec struct {
	// Namespace is the namespace of the service
	Namespace string `json:"namespace,omitempty"`
	// Name is the name of the service
	Name string `json:"name"`
}

type PluginArguments struct {
	Vars []string          `json:"vars,omitempty"`
	Host map[string]string `json:"host,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SearchlightPluginList is a collection of SearchlightPlugin.
type SearchlightPluginList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of SearchlightPlugin.
	Items []SearchlightPlugin `json:"items"`
}
