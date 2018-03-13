package notifier

import (
	"fmt"
	"strings"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
)

func RenderSMS(alert api.Alert, req *Request) string {
	var message string
	if strings.ToUpper(req.Type) == EventTypeAcknowledgement {
		message = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is acked.", alert.GetName(), req.HostName, req.State)
	} else if strings.ToUpper(req.Type) == EventTypeRecovery {
		message = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is recovered.", alert.GetName(), req.HostName, req.State)
	} else if strings.ToUpper(req.Type) == EventTypeProblem {
		message = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nCheck this issue in Icingaweb.", alert.GetName(), req.HostName, req.State)
	} else {
		message = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.", alert.GetName(), req.HostName, req.State)
	}
	if req.Comment != "" {
		if req.Author != "" {
			message = message + " " + fmt.Sprintf(`%s says "%s".`, req.Author, req.Comment)
		} else {
			message = message + " " + fmt.Sprintf(`Comment: "%s".`, req.Comment)
		}
	}

	return message
}
