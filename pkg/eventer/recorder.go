package eventer

import (
	"github.com/appscode/go/log"
	apiv1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
)

const (
	EventReasonNotFound        = "NotFound"
	EventReasonFailedToProceed = "FailedToProceed"

	EventReasonBadNotifier = "BadNotifier"

	// Icinga objects create event list
	EventReasonCreating         = "Creating"
	EventReasonFailedToCreate   = "FailedToCreate"
	EventReasonSuccessfulCreate = "SuccessfulCreate"

	// Icinga objects update event list
	EventReasonUpdating         = "Updating"
	EventReasonFailedToUpdate   = "FailedToUpdate"
	EventReasonSuccessfulUpdate = "SuccessfulUpdate"

	// Icinga objects delete event list
	EventReasonDeleting         = "Deleting"
	EventReasonFailedToDelete   = "FailedToDelete"
	EventReasonSuccessfulDelete = "SuccessfulDelete"

	// Icinga objects sync event list
	EventReasonSync           = "Sync"
	EventReasonFailedToSync   = "FailedToSync"
	EventReasonSuccessfulSync = "SuccessfulSync"
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
	return broadcaster.NewRecorder(scheme.Scheme, apiv1.EventSource{Component: component})
}
