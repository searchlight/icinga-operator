package operator

import (
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/eventer"
	"github.com/appscode/searchlight/pkg/util"
	core "k8s.io/api/core/v1"
)

func (op *Operator) validateAlert(alert api.Alert) bool {
	// Validate IcingaCommand & it's variables.
	// And also check supported IcingaState
	if ok, err := alert.IsValid(); !ok {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonAlertInvalid,
			`Reason: %v`,
			err,
		)
		return false
	}

	// Validate Notifiers configurations
	if err := util.CheckNotifiers(op.KubeClient, alert); err != nil {
		op.recorder.Eventf(
			alert.ObjectReference(),
			core.EventTypeWarning,
			eventer.EventReasonAlertInvalid,
			`Bad notifier config for NodeAlert: "%s@%s". Reason: %v`,
			alert.GetName(), alert.GetNamespace(),
			err,
		)
		return false
	}

	return true
}
