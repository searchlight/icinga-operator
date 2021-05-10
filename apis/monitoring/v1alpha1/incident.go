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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindIncident     = "Incident"
	ResourcePluralIncident   = "incidents"
	ResourceSingularIncident = "incident"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=incidents,singular=incident,categories={monitoring,appscode,all}
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="LastNotification",type="string",JSONPath=".status.lastNotificationType"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type Incident struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Derived information about the incident.
	// +optional
	Status IncidentStatus `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
}

type IncidentStatus struct {
	// Type of last notification, such as problem, acknowledgement, recovery or custom
	LastNotificationType IncidentNotificationType `json:"lastNotificationType" protobuf:"bytes,1,opt,name=lastNotificationType,casttype=IncidentNotificationType"`

	// Notifications for the incident, such as problem or acknowledgement.
	// +optional
	Notifications []IncidentNotification `json:"notifications,omitempty" protobuf:"bytes,2,rep,name=notifications"`
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
	Type IncidentNotificationType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=IncidentNotificationType"`
	// brief output of check command for the incident
	// +optional
	CheckOutput string `json:"checkOutput" protobuf:"bytes,2,opt,name=checkOutput"`
	// name of user making comment
	// +optional
	Author *string `json:"author,omitempty" protobuf:"bytes,3,opt,name=author"`
	// comment made by user
	// +optional
	Comment *string `json:"comment,omitempty" protobuf:"bytes,4,opt,name=comment"`
	// The time at which this notification was first recorded. (Time of server receipt is in TypeMeta.)
	// +optional
	FirstTimestamp metav1.Time `json:"firstTimestamp,omitempty" protobuf:"bytes,5,opt,name=firstTimestamp"`
	// The time at which the most recent occurrence of this notification was recorded.
	// +optional
	LastTimestamp metav1.Time `json:"lastTimestamp,omitempty" protobuf:"bytes,6,opt,name=lastTimestamp"`
	// state of incident, such as Critical, Warning, OK, Unknown
	LastState string `json:"state" protobuf:"bytes,7,opt,name=state"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IncidentList is a collection of Incident.
type IncidentList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of Incident.
	Items []Incident `json:"items" protobuf:"bytes,2,rep,name=items"`
}
