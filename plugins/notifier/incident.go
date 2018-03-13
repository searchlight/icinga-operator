package notifier

import (
	"fmt"
	"time"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/clientset/versioned/typed/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func appendIncidentNotification(notifications []api.IncidentNotification, req *Request) []api.IncidentNotification {
	notification := api.IncidentNotification{
		Type:           api.AlertType(req.Type),
		CheckOutput:    req.Output,
		Author:         &req.Author,
		Comment:        &req.Comment,
		FirstTimestamp: metav1.NewTime(req.Time),
		LastTimestamp:  metav1.NewTime(req.Time),
		LastState:      req.State,
	}
	notifications = append(notifications, notification)
	return notifications
}

func updateIncidentNotification(notification api.IncidentNotification, req *Request) api.IncidentNotification {
	notification.CheckOutput = req.Output
	notification.Author = &req.Author
	notification.Comment = &req.Comment
	notification.LastTimestamp = metav1.NewTime(req.Time)
	notification.LastState = req.State
	return notification
}

func getLabel(req *Request, icingaHost *icinga.IcingaHost) map[string]string {
	labelMap := map[string]string{
		api.LabelKeyAlertType:        icingaHost.Type,
		api.LabelKeyAlert:            req.AlertName,
		api.LabelKeyObjectName:       icingaHost.ObjectName,
		api.LabelKeyProblemRecovered: "false",
	}

	return labelMap
}

func generateIncidentName(req *Request) (string, error) {
	host, err := icinga.ParseHost(req.HostName)
	if err != nil {
		return "", err
	}

	t := req.Time.Format("20060102-1504")

	switch host.Type {
	case icinga.TypePod, icinga.TypeNode:
		return host.Type + "." + host.ObjectName + "." + req.AlertName + "." + t, nil
	case icinga.TypeCluster:
		return host.Type + "." + req.AlertName + "." + t, nil
	}

	return "", fmt.Errorf("unknown host type %s", host.Type)
}

func reconcileIncident(client *cs.MonitoringV1alpha1Client, req *Request) error {
	host, err := icinga.ParseHost(req.HostName)
	if err != nil {
		return err
	}

	incidentList, err := client.Incidents(host.AlertNamespace).List(metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(getLabel(req, host)).String(),
	})
	if err != nil {
		return err
	}

	var lastCreationTimestamp time.Time
	var incident *api.Incident

	for _, item := range incidentList.Items {
		if item.CreationTimestamp.After(lastCreationTimestamp) {
			lastCreationTimestamp = item.CreationTimestamp.Time
			incident = &item
		}
	}

	if incident != nil {
		notifications := incident.Status.Notifications
		if api.AlertType(req.Type) == api.NotificationCustom {
			notifications = appendIncidentNotification(notifications, req)
		} else {
			updated := false
			for i := len(notifications) - 1; i >= 0; i-- {
				notification := notifications[i]
				if notification.Type == api.NotificationAcknowledgement {
					continue
				}
				if api.AlertType(req.Type) == notification.Type {
					notifications[i] = updateIncidentNotification(notification, req)
					updated = true
					break
				}
			}
			if !updated {
				notifications = appendIncidentNotification(notifications, req)
			}
		}

		incident.Status.LastNotificationType = api.AlertType(req.Type)
		incident.Status.Notifications = notifications

		if api.AlertType(req.Type) == api.NotificationRecovery {
			incident.Labels[api.LabelKeyProblemRecovered] = "true"
		}

		if _, err := client.Incidents(incident.Namespace).Update(incident); err != nil {
			return err
		}
	} else {
		name, err := generateIncidentName(req)
		if err != nil {
			return err
		}

		incident := &api.Incident{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: host.AlertNamespace,
				Labels:    getLabel(req, host),
			},
			Status: api.IncidentStatus{
				LastNotificationType: api.AlertType(req.Type),
				Notifications:        appendIncidentNotification(make([]api.IncidentNotification, 0), req),
			},
		}

		if _, err = client.Incidents(incident.Namespace).Create(incident); err != nil {
			return err
		}
	}

	return nil
}
