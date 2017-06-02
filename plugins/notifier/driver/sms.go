package driver

import (
	"fmt"
	"strings"

	api "github.com/appscode/api/kubernetes/v1beta1"
	aci "github.com/appscode/searchlight/api"
)

const (
	EventTypeProblem         = "PROBLEM"
	EventTypeAcknowledgement = "ACKNOWLEDGEMENT"
	EventTypeRecovery        = "RECOVERY"
)

func RenderSMS(alert *aci.Alert, req *api.IncidentNotifyRequest) (string, error) {
	clusterInfo := ""
	if req.KubernetesCluster != "" {
		clusterInfo = fmt.Sprintf(`Cluster: %s.\n`, req.KubernetesCluster)
	}
	if strings.ToUpper(req.Type) == EventTypeAcknowledgement {
		return clusterInfo + fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is acked.", alert.Name, req.HostName, req.State), nil
	} else if strings.ToUpper(req.Type) == EventTypeRecovery {
		return clusterInfo + fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is recovered.", alert.Name, req.HostName, req.State), nil
	} else if strings.ToUpper(req.Type) == EventTypeProblem {
		return clusterInfo + fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nCheck this issue in Icingaweb.", alert.Name, req.HostName, req.State), nil
	} else {
		return clusterInfo + fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.", alert.Name, req.HostName, req.State), nil
	}
}
