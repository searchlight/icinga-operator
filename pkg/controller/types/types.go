package types

import (
	"sync"

	aci "github.com/appscode/searchlight/api"
	acs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/stash"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

type EventReason string

const (
	EventReasonNotFound        EventReason = "NotFound"
	EventReasonFailedToProceed EventReason = "FailedToProceed"

	// Icinga objects create event list
	EventReasonCreating         EventReason = "Creating"
	EventReasonFailedToCreate   EventReason = "FailedToCreate"
	EventReasonSuccessfulCreate EventReason = "SuccessfulCreate"

	// Icinga objects update event list
	EventReasonUpdating         EventReason = "Updating"
	EventReasonFailedToUpdate   EventReason = "FailedToUpdate"
	EventReasonSuccessfulUpdate EventReason = "SuccessfulUpdate"

	// Icinga objects delete event list
	EventReasonDeleting         EventReason = "Deleting"
	EventReasonFailedToDelete   EventReason = "FailedToDelete"
	EventReasonSuccessfulDelete EventReason = "SuccessfulDelete"

	// Icinga objects sync event list
	EventReasonSync           EventReason = "Sync"
	EventReasonFailedToSync   EventReason = "FailedToSync"
	EventReasonSuccessfulSync EventReason = "SuccessfulSync"
)

func (r EventReason) String() string {
	return string(r)
}

const (
	AcknowledgeTimestamp string = "acknowledgement_timestamp"
)

type IcingaData struct {
	HostType map[string]string
	VarInfo  map[string]data.CommandVar
}

type Context struct {
	// kubernetes client
	KubeClient clientset.Interface
	ExtClient  acs.ExtensionInterface

	IcingaClient *icinga.IcingaClient
	IcingaData   map[string]*IcingaData

	Resource   *aci.Alert
	ObjectType string
	ObjectName string

	Storage *stash.Storage
	sync.Mutex
}

type KubeOptions struct {
	ObjectType string
	ObjectName string
}

type Ancestors struct {
	Type  string   `json:"type,omitempty"`
	Names []string `json:"names,omitempty"`
}

type AlertEventMessage struct {
	IncidentEventId int64  `json:"incident_event_id,omitempty"`
	Comment         string `json:"comment,omitempty"`
	UserName        string `json:"username,omitempty"`
}
