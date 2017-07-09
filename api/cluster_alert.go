package api

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindClusterAlert = "ClusterAlert"
	ResourceNameClusterAlert = "clusteralert"
	ResourceTypeClusterAlert = "clusteralerts"
)

// ClusterAlert types for appscode.
type ClusterAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the ClusterAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec ClusterAlertSpec `json:"spec,omitempty"`

	// Status is the current state of the ClusterAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Status AlertStatus `json:"status,omitempty"`
}

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

	// NotifierParams contains information to send notifications for Incident
	// State, UserUid, Method
	Receivers []Receiver `json:"receivers,omitempty"`

	// Vars contains Icinga Service variables to be used in CheckCommand
	Vars map[string]interface{} `json:"vars,omitempty"`
}

func (s ClusterAlertSpec) IsValid() (bool, error) {
	cmd, ok := ClusterCommands[s.Check]
	if !ok {
		return false, fmt.Errorf("%s is not a valid cluster check command.", s.Check)
	}
	for k := range s.Vars {
		if _, ok := cmd.Vars[k]; !ok {
			return false, fmt.Errorf("Var %s is unsupported for check command %s.", k, s.Check)
		}
	}
	for _, rcv := range s.Receivers {
		found := false
		for _, state := range cmd.States {
			if state == rcv.State {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Errorf("State %s is unsupported for check command %s.", rcv.State, s.Check)
		}
	}
	return true, nil
}
