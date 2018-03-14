package monitoring

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CheckNode string

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeAlert struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta

	// Spec is the desired state of the NodeAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec NodeAlertSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeAlertList is a collection of NodeAlert.
type NodeAlertList struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta

	// Items is the list of NodeAlert.
	Items []NodeAlert
}

// NodeAlertSpec describes the NodeAlert the user wishes to create.
type NodeAlertSpec struct {
	Selector map[string]string

	NodeName *string

	// Icinga CheckCommand name
	Check CheckNode

	// How frequently Icinga Service will be checked
	CheckInterval metav1.Duration

	// How frequently notifications will be send
	AlertInterval metav1.Duration

	// Secret containing notifier credentials
	NotifierSecretName string

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	Receivers []Receiver

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]string
}
