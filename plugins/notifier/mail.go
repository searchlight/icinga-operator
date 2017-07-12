package notifier

import (
	"time"

	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/flosch/pongo2"
)

var subjectMap = map[string]string{
	"PROBLEM":         "Problem Detected",
	"ACKNOWLEDGEMENT": "Problem Acknowledged",
	"RECOVERY":        "Problem Recovered",
	"CUSTOM":          "Custom Notification",
}

func render(ctx *pongo2.Context, template string) (string, error) {
	tpl, err := pongo2.FromString(template)
	if err != nil {
		return "", err
	}

	body, err := tpl.Execute(*ctx)
	if err != nil {
		return "", err
	}
	return body, nil
}

func RenderMail(alert aci.Alert, req *Request) (string, error) {
	t := time.Unix(req.Time, 0)

	host, err := icinga.ParseHost(req.HostName)
	if err != nil {
		return "", err
	}

	data := map[string]interface{}{
		"KubernetesNamespace":  host.AlertNamespace,
		"kubernetesAlertType":  host.Type,
		"kubernetesAlertName":  alert.GetName(),
		"kubernetesObjectName": host.ObjectName,
		"IcingaHostName":       req.HostName,
		"IcingaServiceName":    alert.GetName(),
		"CheckCommand":         alert.Command(),
		"IcingaType":           req.Type,
		"IcingaState":          req.State,
		"IcingaOutput":         req.Output,
		"IcingaTime":           t,
	}

	pCtx := pongo2.Context(data)
	return render(&pCtx, notificationMailTemplate)
}
