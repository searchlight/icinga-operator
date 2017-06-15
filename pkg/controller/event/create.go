package event

import (
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/controller/types"
	"github.com/appscode/searchlight/pkg/events"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func CreateAlertEvent(kubeClient clientset.Interface, alert *aci.Alert, reason types.EventReason, additionalMessage ...string) {
	timestamp := metav1.NewTime(time.Now().UTC())
	event := &apiv1.Event{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("alert"),
			Namespace: alert.Namespace,
		},
		InvolvedObject: apiv1.ObjectReference{
			Kind:      events.ObjectKindAlert.String(),
			Namespace: alert.Namespace,
			Name:      alert.Name,
		},
		Source: apiv1.EventSource{
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
		event.Type = apiv1.EventTypeWarning
	case types.EventReasonFailedToProceed:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to proceed. Reason: %v`, additionalMessage)
		event.Type = apiv1.EventTypeWarning

	case types.EventReasonCreating:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`creating Icinga objects`)
		event.Type = apiv1.EventTypeNormal
	case types.EventReasonFailedToCreate:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to create Icinga objects. Error: %v`, additionalMessage)
		event.Type = apiv1.EventTypeWarning
	case types.EventReasonSuccessfulCreate:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully created Icinga objects`)
		event.Type = apiv1.EventTypeNormal

	case types.EventReasonUpdating:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`updating Icinga objects`)
	case types.EventReasonFailedToUpdate:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to update Icinga objects. Error: %v`, additionalMessage)
		event.Type = apiv1.EventTypeWarning
	case types.EventReasonSuccessfulUpdate:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully updated Icinga objects.`)
		event.Type = apiv1.EventTypeNormal

	case types.EventReasonDeleting:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`deleting Icinga objects`)
		event.Type = apiv1.EventTypeNormal
	case types.EventReasonFailedToDelete:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to delete Icinga objects. Error: %v`, additionalMessage)
		event.Type = apiv1.EventTypeWarning
	case types.EventReasonSuccessfulDelete:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully deleted Icinga objects.`)
		event.Type = apiv1.EventTypeNormal

	case types.EventReasonSync:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`synchronizing alert for %v.`, additionalMessage[0])
		event.Type = apiv1.EventTypeNormal
	case types.EventReasonFailedToSync:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`failed to synchronize alert for %v. Error: %v`, additionalMessage[0], additionalMessage[1])
		event.Type = apiv1.EventTypeWarning
	case types.EventReasonSuccessfulSync:
		event.Reason = reason.String()
		event.Message = fmt.Sprintf(`successfully synchronized alert for %v.`, additionalMessage[0])
		event.Type = apiv1.EventTypeNormal
	}

	if _, err := kubeClient.Core().Events(alert.Namespace).Create(event); err != nil {
		log.Debugln(err)
	}
}
