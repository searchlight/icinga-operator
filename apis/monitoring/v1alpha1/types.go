package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type ClusterAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the ClusterAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec ClusterAlertSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterAlertList is a collection of ClusterAlert.
type ClusterAlertList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of ClusterAlert.
	Items []ClusterAlert `json:"items"`
}

// ClusterAlertSpec describes the ClusterAlert the user wishes to create.
type ClusterAlertSpec struct {
	// Icinga CheckCommand name
	Check CheckCluster `json:"check,omitempty"`

	// How frequently Icinga Service will be checked
	CheckInterval metav1.Duration `json:"checkInterval,omitempty"`

	// How frequently notifications will be send
	AlertInterval metav1.Duration `json:"alertInterval,omitempty"`

	// Secret containing notifier credentials
	NotifierSecretName string `json:"notifierSecretName,omitempty"`

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	Receivers []Receiver `json:"receivers,omitempty"`

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]string `json:"vars,omitempty"`
}

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type NodeAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the NodeAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec NodeAlertSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

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

	NodeName *string `json:"nodeName,omitempty"`

	// Icinga CheckCommand name
	Check CheckNode `json:"check,omitempty"`

	// How frequently Icinga Service will be checked
	CheckInterval metav1.Duration `json:"checkInterval,omitempty"`

	// How frequently notifications will be send
	AlertInterval metav1.Duration `json:"alertInterval,omitempty"`

	// Secret containing notifier credentials
	NotifierSecretName string `json:"notifierSecretName,omitempty"`

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	Receivers []Receiver `json:"receivers,omitempty"`

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]string `json:"vars,omitempty"`
}

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PodAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the PodAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec PodAlertSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodAlertList is a collection of PodAlert.
type PodAlertList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of PodAlert.
	Items []PodAlert `json:"items"`
}

// PodAlertSpec describes the PodAlert the user wishes to create.
type PodAlertSpec struct {
	Selector *metav1.LabelSelector `json:"selector,omitempty"`

	PodName *string `json:"podName,omitempty"`

	// Icinga CheckCommand name
	Check CheckPod `json:"check,omitempty"`

	// How frequently Icinga Service will be checked
	CheckInterval metav1.Duration `json:"checkInterval,omitempty"`

	// How frequently notifications will be send
	AlertInterval metav1.Duration `json:"alertInterval,omitempty"`

	// Secret containing notifier credentials
	NotifierSecretName string `json:"notifierSecretName,omitempty"`

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	Receivers []Receiver `json:"receivers,omitempty"`

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]string `json:"vars,omitempty"`
}

type Receiver struct {
	// For which state notification will be sent
	State string `json:"state,omitempty"`

	// To whom notification will be sent
	To []string `json:"to,omitempty"`

	// How this notification will be sent
	Notifier string `json:"notifier,omitempty"`
}

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Incident struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Derived information about the incident.
	// +optional
	Status IncidentStatus `json:"status,omitempty"`
}

type IncidentStatus struct {
	// state of incident, such as CRITICAL, WARNING, OK
	LastNotificationType string `json:"lastNotificationType"`

	// Notifications for the incident, such as problem or acknowledge.
	// +optional
	Notifications []IncidentNotification `json:"notifications,omitempty"`
}

type IncidentNotificationType string

// These are the possible notifications for an incident.
const (
	NotificationProblem         IncidentNotificationType = "Problem"
	NotificationAcknowledgement IncidentNotificationType = "Acknowledgement"
	NotificationRecovery        IncidentNotificationType = "Recovery"
	NotificationCustom          IncidentNotificationType = "Custom"
)

type IncidentNotification struct {
	// incident notification type.
	Type IncidentNotificationType `json:"type"`
	// brief output of check command for the incident
	// +optional
	CheckOutput string `json:"checkOutput"`
	// name of user making comment
	// +optional
	Author *string `json:"comment,omitempty"`
	// comment made by user
	// +optional
	Comment *string `json:"comment,omitempty"`
	// The time at which this notification was first recorded. (Time of server receipt is in TypeMeta.)
	// +optional
	FirstTimestamp metav1.Time `json:"firstTimestamp,omitempty"`
	// The time at which the most recent occurrence of this notification was recorded.
	// +optional
	LastTimestamp metav1.Time `json:"lastTimestamp,omitempty"`
	// state of incident, such as CRITICAL, WARNING, OK, UNKNOWN
	LastState string `json:"state"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IncidentList is a collection of Incident.
type IncidentList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of Incident.
	Items []Incident `json:"items"`
}
