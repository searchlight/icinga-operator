package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindPlugin     = "Plugin"
	ResourcePluralPlugin   = "plugins"
	ResourceSingularPlugin = "plugin"
)

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Plugin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the Plugin.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec PluginSpec `json:"spec,omitempty"`
}

// PluginSpec describes the Plugin the user wishes to create.
type PluginSpec struct {
	// Check Command
	Command string `json:"command,omitempty"`

	// Webhook provides a reference to the service for this Plugin.
	// It must communicate on port 80
	Webhook *WebhookServiceSpec `json:"webhook,omitempty"`

	// AlertKinds refers to supports Alert kinds for this plugin
	AlertKinds []string `json:"alertKinds"`
	// Supported arguments for Plugin
	Arguments PluginArguments `json:"arguments,omitempty"`
	// Supported Icinga Service State
	State []string `json:"state"`
}

type WebhookServiceSpec struct {
	// Namespace is the namespace of the service
	Namespace string `json:"namespace,omitempty"`
	// Name is the name of the service
	Name string `json:"name"`
	// InsecureSkipTLSVerify disables TLS certificate verification when communicating with this webhook.
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`
	// CABundle is a PEM encoded CA bundle which will be used to validate an webhook's serving certificate.
	CABundle []byte `json:"caBundle,omitempty"`
}

type PluginArguments struct {
	Vars    []string          `json:"vars,omitempty"`
	Host    map[string]string `json:"host,omitempty"`
	Service map[string]string `json:"service,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PluginList is a collection of Plugin.
type PluginList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of Plugin.
	Items []Plugin `json:"items"`
}