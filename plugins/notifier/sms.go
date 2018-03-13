package notifier

import (
	"fmt"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
)

func RenderSMS(alert api.Alert, req *Request) string {
	var msg string

	switch api.AlertType(req.Type) {
	case api.NotificationAcknowledgement:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is acked.", alert.GetName(), req.HostName, req.State)
	case api.NotificationRecovery:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is recovered.", alert.GetName(), req.HostName, req.State)
	case api.NotificationProblem:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nCheck this issue in Icingaweb.", alert.GetName(), req.HostName, req.State)
	default:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.", alert.GetName(), req.HostName, req.State)
	}
	if req.Comment != "" {
		if req.Author != "" {
			msg = msg + " " + fmt.Sprintf(`%s says "%s".`, req.Author, req.Comment)
		} else {
			msg = msg + " " + fmt.Sprintf(`Comment: "%s".`, req.Comment)
		}
	}

	return msg
}
