package monitoring

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:skipVerbs=updateStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Incident struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta

	// Derived information about the incident.
	// +optional
	Status IncidentStatus
}

type IncidentStatus struct {
	// state of incident, such as Critical, Warning, OK
	LastNotificationType IncidentNotificationType

	// Notifications for the incident, such as problem or acknowledge.
	// +optional
	Notifications []IncidentNotification
}

type IncidentNotificationType string

type IncidentNotification struct {
	// incident notification type.
	Type IncidentNotificationType
	// brief output of check command for the incident
	// +optional
	CheckOutput string
	// name of user making comment
	// +optional
	Author *string
	// comment made by user
	// +optional
	Comment *string
	// The time at which this notification was first recorded. (Time of server receipt is in TypeMeta.)
	// +optional
	FirstTimestamp metav1.Time
	// The time at which the most recent occurrence of this notification was recorded.
	// +optional
	LastTimestamp metav1.Time
	// state of incident, such as Critical, Warning, OK, Unknown
	LastState string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IncidentList is a collection of Incident.
type IncidentList struct {
	metav1.TypeMeta
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta

	// Items is the list of Incident.
	Items []Incident
}
