package eventer

import (
	"github.com/appscode/log"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/record"
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


func NewEventRecorder(client clientset.Interface, component string) record.EventRecorder {
	// Event Broadcaster
	broadcaster := record.NewBroadcaster()
	broadcaster.StartEventWatcher(
		func(event *apiv1.Event) {
			if _, err := client.CoreV1().Events(event.Namespace).Create(event); err != nil {
				log.Errorln(err)
			}
		},
	)
	// Event Recorder
	return broadcaster.NewRecorder(api.Scheme, apiv1.EventSource{Component: component})
}
