package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindNodeAlert = "NodeAlert"
	ResourceNameNodeAlert = "nodealert"
	ResourceTypeNodeAlert = "nodealerts"
)

// NodeAlert types for appscode.
type NodeAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the NodeAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec NodeAlertSpec `json:"spec,omitempty"`

	// Status is the current state of the NodeAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Status AlertStatus `json:"status,omitempty"`
}

// NodeAlertList is a collection of NodeAlert.
type NodeAlertList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of NodeAlert.
	Items []NodeAlert `json:"items"`
}

// NodeAlertSpec describes the NodeAlert the user wishes to create.
type NodeAlertSpec struct {
	Selector map[string]string `json:"selector,omitempty"`

	// IcingaParam contains parameters for Icinga config
	IcingaParam *IcingaParam `json:"icingaParam,omitempty"`

	// Icinga CheckCommand name
	CheckCommand string `json:"checkCommand,omitempty"`

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	NotifierParams []NotifierParam `json:"notifierParams,omitempty"`

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]interface{} `json:"vars,omitempty"`
}
