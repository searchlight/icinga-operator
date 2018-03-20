package notifier

import (
	"fmt"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
)

func (p *plugin) RenderSMS(alert api.Alert) string {
	opts := p.options
	var msg string

	switch api.AlertType(opts.notificationType) {
	case api.NotificationAcknowledgement:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is acked.", alert.GetName(), opts.hostname, opts.serviceState)
	case api.NotificationRecovery:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is recovered.", alert.GetName(), opts.hostname, opts.serviceState)
	case api.NotificationProblem:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nCheck this issue in Icingaweb.", alert.GetName(), opts.hostname, opts.serviceState)
	default:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.", alert.GetName(), opts.hostname, opts.serviceState)
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
