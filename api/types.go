package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AlertPhase string

const (
	// used for Alert that are currently creating
	AlertPhaseCreating AlertPhase = "Creating"
	// used for Alert that are created
	AlertPhaseCreated AlertPhase = "Created"
	// used for Alert that are currently deleting
	AlertPhaseDeleting AlertPhase = "Deleting"
	// used for Alert that are Failed
	AlertPhaseFailed AlertPhase = "Failed"
)

type AlertStatus struct {
	CreationTime *metav1.Time `json:"creationTime,omitempty"`
	UpdateTime   *metav1.Time `json:"updateTime,omitempty"`
	Phase        AlertPhase   `json:"phase,omitempty"`
	Reason       string       `json:"reason,omitempty"`
}

type IcingaParam struct {
	// How frequently Icinga Service will be checked
	CheckIntervalSec int64 `json:"checkIntervalSec,omitempty"`

	// How frequently notifications will be send
	AlertIntervalSec int64 `json:"alertIntervalSec,omitempty"`
}

type NotifierParam struct {
	// For which state notification will be sent
	State string `json:"state,omitempty"`

	// To whom notification will be sent
	UserUid string `json:"userUid,omitempty"`

	// How this notification will be sent
	Method string `json:"method,omitempty"`
}
