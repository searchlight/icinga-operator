package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindPodAlert = "PodAlert"
	ResourceNamePodAlert = "podalert"
	ResourceTypePodAlert = "podalerts"
)

// PodAlert types for appscode.
type PodAlert struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the desired state of the PodAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Spec PodAlertSpec `json:"spec,omitempty"`

	// Status is the current state of the PodAlert.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#spec-and-status
	Status AlertStatus `json:"status,omitempty"`
}

// PodAlertList is a collection of PodAlert.
type PodAlertList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: http://releases.k8s.io/release-1.2/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is the list of PodAlert.
	Items []PodAlert `json:"items"`
}

type CheckPod string

const (
	CheckPodInfluxQuery       CheckPod     = "influx_query"
	CheckPodStatus            CheckPod     = "pod_status"
	CheckPodPrometheusMetric  CheckPod     = "prometheus_metric"
	CheckVolume               CheckPod     = "volume"
	CheckPodExec              CheckPod     = "kube_exec"
	CheckNodeInfluxQuery      CheckNode    = "influx_query"
	CheckNodeDisk             CheckNode    = "node_disk"
	CheckNodeStatus           CheckNode    = "node_status"
	CheckNodePrometheusMetric CheckNode    = "prometheus_metric"
	CheckHttp                 CheckCluster = "any_http"
	CheckComponentStatus      CheckCluster = "component_status"
	CheckJsonPath             CheckCluster = "json_path"
	CheckNodeCount            CheckCluster = "node_count"
	CheckPodExists            CheckCluster = "pod_exists"
	CheckPrometheusMetric     CheckCluster = "prometheus_metric"
	CheckClusterEvent         CheckCluster = "kube_event"
	CheckHelloIcinga          CheckCluster = "hello_icinga"
	CheckDIG                  CheckCluster = "dig"
	CheckDNS                  CheckCluster = "dns"
	CheckDummy                CheckCluster = "dummy"
	CheckICMP                 CheckCluster = "icmp"
)

// PodAlertSpec describes the PodAlert the user wishes to create.
type PodAlertSpec struct {
	Selector metav1.LabelSelector `json:"selector,omitempty"`

	// Icinga CheckCommand name
	Check string `json:"check,omitempty"`

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
