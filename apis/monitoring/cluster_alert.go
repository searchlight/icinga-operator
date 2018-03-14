package monitoring

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CheckCluster string

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ClusterAlert struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta

	// Spec is the desired state of the ClusterAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec ClusterAlertSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAlertList is a collection of ClusterAlert.
type ClusterAlertList struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta

	// Items is the list of ClusterAlert.
	Items []ClusterAlert
}

// ClusterAlertSpec describes the ClusterAlert the user wishes to create.
type ClusterAlertSpec struct {
	// Icinga CheckCommand name
	Check CheckCluster

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
