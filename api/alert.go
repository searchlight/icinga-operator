package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindAlert = "Alert"
	ResourceNameAlert = "alert"
	ResourceTypeAlert = "alerts"
)

// Alert types for appscode.
type Alert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the Alert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec AlertSpec `json:"spec,omitempty"`

	// Status is the current state of the Alert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Status AlertStatus `json:"status,omitempty"`
}

// AlertList is a collection of Alert.
type AlertList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of Alert.
	Items []Alert `json:"items"`
}

// AlertSpec describes the Alert the user wishes to create.
type AlertSpec struct {
	Selector metav1.LabelSelector `json:"selector,omitempty"`

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
