package notifier

import (
	"fmt"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
)

func (n *notifier) RenderSMS(receiver api.Receiver) string {
	opts := n.options
	var msg string

	switch api.AlertType(opts.notificationType) {
	case api.NotificationAcknowledgement:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is acked.", opts.alertName, opts.hostname, receiver.State)
	case api.NotificationRecovery:
		msg = fmt.Sprintf("Service [%s] for [%s] was in \"%s\" state.\nThis issue is recovered.", opts.alertName, opts.hostname, receiver.State)
	case api.NotificationProblem:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nCheck this issue in Icingaweb.", opts.alertName, opts.hostname, receiver.State)
	default:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.", opts.alertName, opts.hostname, receiver.State)
	}
	if opts.comment != "" {
		if opts.author != "" {
			msg = msg + " " + fmt.Sprintf(`%s says "%s".`, opts.author, opts.comment)
		} else {
			msg = msg + " " + fmt.Sprintf(`Comment: "%s".`, opts.comment)
		}
	}

	return msg
}
