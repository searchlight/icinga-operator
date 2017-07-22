package notifier

import (
	"fmt"
	"strings"

	aci "github.com/appscode/searchlight/api"
)

const (
	EventTypeProblem         = "PROBLEM"
	EventTypeAcknowledgement = "ACKNOWLEDGEMENT"
	EventTypeRecovery        = "RECOVERY"
)

func RenderSMS(alert aci.Alert, req *Request) (string, error) {
	if strings.ToUpper(req.Type) == EventTypeAcknowledgement {
		return fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is acked.", alert.GetName(), req.HostName, req.State), nil
	} else if strings.ToUpper(req.Type) == EventTypeRecovery {
		return fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is recovered.", alert.GetName(), req.HostName, req.State), nil
	} else if strings.ToUpper(req.Type) == EventTypeProblem {
		return fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nCheck this issue in Icingaweb.", alert.GetName(), req.HostName, req.State), nil
	} else {
		return fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.", alert.GetName(), req.HostName, req.State), nil
	}
}
