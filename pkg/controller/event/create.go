package event

import (
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/controller/types"
	"github.com/appscode/searchlight/pkg/events"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/unversioned"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

func CreateAlertEvent(kubeClient clientset.Interface, alert *aci.Alert, reason types.EventReason, additionalMessage ...string) {
	timestamp := unversioned.NewTime(time.Now().UTC())
	event := &kapi.Event{
		ObjectMeta: kapi.ObjectMeta{
			Name:      rand.WithUniqSuffix("alert"),
			Namespace: alert.Namespace,
		},
		InvolvedObject: kapi.ObjectReference{
			Kind:      events.ObjectKindAlert.String(),
			Namespace: alert.Namespace,
			Name:      alert.Name,
		},
		Source: kapi.EventSource{
			Component: "searchlight",
		},

		Count:          1,
		FirstTimestamp: timestamp,
		LastTimestamp:  timestamp,
	}

	switch reason {
	case types.EventReasonNotFound:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to set alerts. Reason: %v`, additionalMessage)
		event.Type = kapi.EventTypeWarning
	case types.EventReasonFailedToProceed:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to proceed. Reason: %v`, additionalMessage)
		event.Type = kapi.EventTypeWarning

	case types.EventReasonCreating:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`creating Icinga objects`)
		event.Type = kapi.EventTypeNormal
	case types.EventReasonFailedToCreate:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to create Icinga objects. Error: %v`, additionalMessage)
		event.Type = kapi.EventTypeWarning
	case types.EventReasonSuccessfulCreate:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully created Icinga objects`)
		event.Type = kapi.EventTypeNormal

	case types.EventReasonUpdating:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`updating Icinga objects`)
	case types.EventReasonFailedToUpdate:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to update Icinga objects. Error: %v`, additionalMessage)
		event.Type = kapi.EventTypeWarning
	case types.EventReasonSuccessfulUpdate:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully updated Icinga objects.`)
		event.Type = kapi.EventTypeNormal

	case types.EventReasonDeleting:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`deleting Icinga objects`)
		event.Type = kapi.EventTypeNormal
	case types.EventReasonFailedToDelete:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to delete Icinga objects. Error: %v`, additionalMessage)
		event.Type = kapi.EventTypeWarning
	case types.EventReasonSuccessfulDelete:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully deleted Icinga objects.`)
		event.Type = kapi.EventTypeNormal

	case types.EventReasonSync:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`synchronizing alert for %v.`, additionalMessage[0])
		event.Type = kapi.EventTypeNormal
	case types.EventReasonFailedToSync:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to synchronize alert for %v. Error: %v`, additionalMessage[0], additionalMessage[1])
		event.Type = kapi.EventTypeWarning
	case types.EventReasonSuccessfulSync:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully synchronized alert for %v.`, additionalMessage[0])
		event.Type = kapi.EventTypeNormal
	}

	if _, err := kubeClient.Core().Events(alert.Namespace).Create(event); err != nil {
		log.Debugln(err)
	}
}
