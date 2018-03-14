package monitoring

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CheckPod string

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PodAlert struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta

	// Spec is the desired state of the PodAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec PodAlertSpec
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodAlertList is a collection of PodAlert.
type PodAlertList struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta

	// Items is the list of PodAlert.
	Items []PodAlert
}

// PodAlertSpec describes the PodAlert the user wishes to create.
type PodAlertSpec struct {
	Selector *metav1.LabelSelector

	PodName *string

	// Icinga CheckCommand name
	Check CheckPod

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
