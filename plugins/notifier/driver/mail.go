package driver

import (
	"time"

	api "github.com/appscode/api/kubernetes/v1beta1"
	aci "github.com/appscode/searchlight/api"
	"github.com/flosch/pongo2"
)

const (
	ObjectType = "monitoring.appscode.com/objectType"
	ObjectName = "monitoring.appscode.com/objectName"
)

var SubjectMap = map[string]string{
	"PROBLEM":         "Problem Detected",
	"ACKNOWLEDGEMENT": "Problem Acknowledged",
	"RECOVERY":        "Problem Recovered",
	"CUSTOM":          "Custom Notification",
}

type labelMap map[string]string

func (s labelMap) ObjectType() string {
	v, _ := s[ObjectType]
	return v
}

func (s labelMap) ObjectName() string {
	v, _ := s[ObjectName]
	return v
}
func getObjectInfo(label map[string]string) (objectType string, objectName string) {
	opts := labelMap(label)
	objectType = opts.ObjectType()
	objectName = opts.ObjectName()
	return
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

func RenderMail(alert *aci.Alert, req *api.IncidentNotifyRequest) (string, error) {
	t := time.Unix(req.Time, 0)

	objectType, objectName := getObjectInfo(alert.Labels)

	data := map[string]interface{}{
		"KubernetesCluster":    req.KubernetesCluster,
		"KubernetesNamespace":  alert.Namespace,
		"kubernetesObjectType": objectType,
		"kubernetesObjectName": objectName,
		"IcingaHostName":       req.HostName,
		"IcingaServiceName":    alert.Name,
		"CheckCommand":         alert.Spec.CheckCommand,
		"IcingaType":           req.Type,
		"IcingaState":          req.State,
		"IcingaOutput":         req.Output,
		"IcingaTime":           t,
	}

	pCtx := pongo2.Context(data)
	return render(&pCtx, notificationMailTemplate)
}
