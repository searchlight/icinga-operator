/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package notifier

import (
	"context"
	"fmt"
	"time"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/client/clientset/versioned/typed/monitoring/v1alpha1/util"
	"go.searchlight.dev/icinga-operator/pkg/icinga"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

func (n *notifier) appendIncidentNotification(notifications []api.IncidentNotification) []api.IncidentNotification {
	opts := n.options
	notification := api.IncidentNotification{
		Type:           api.AlertType(opts.notificationType),
		CheckOutput:    opts.serviceOutput,
		Author:         &opts.author,
		Comment:        &opts.comment,
		FirstTimestamp: metav1.NewTime(opts.time),
		LastTimestamp:  metav1.NewTime(opts.time),
		LastState:      opts.serviceState,
	}
	notifications = append(notifications, notification)
	return notifications
}

func (n *notifier) updateIncidentNotification(notification api.IncidentNotification) api.IncidentNotification {
	opts := n.options
	notification.CheckOutput = opts.serviceOutput
	notification.Author = &opts.author
	notification.Comment = &opts.comment
	notification.LastTimestamp = metav1.NewTime(opts.time)
	notification.LastState = opts.serviceState
	return notification
}

func (n *notifier) getLabel() map[string]string {
	labelMap := map[string]string{
		api.LabelKeyAlertType:        n.options.host.Type,
		api.LabelKeyAlert:            n.options.alertName,
		api.LabelKeyObjectName:       n.options.host.ObjectName,
		api.LabelKeyProblemRecovered: "false",
	}

	return labelMap
}

func (n *notifier) generateIncidentName() (string, error) {
	host := n.options.host
	t := n.options.time.Format("20060102-1504")

	switch host.Type {
	case icinga.TypePod, icinga.TypeNode:
		return host.Type + "." + host.ObjectName + "." + n.options.alertName + "." + t, nil
	case icinga.TypeCluster:
		return host.Type + "." + n.options.alertName + "." + t, nil
	}

	return "", fmt.Errorf("unknown host type %s", host.Type)
}

func (n *notifier) reconcileIncident() error {
	opts := n.options

	incident, err := n.getIncident()
	if err != nil {
		return err
	}

	if incident != nil {
		notifications := incident.Status.Notifications
		if api.AlertType(opts.notificationType) == api.NotificationCustom {
			notifications = n.appendIncidentNotification(notifications)
		} else {
			updated := false
			for i := len(notifications) - 1; i >= 0; i-- {
				notification := notifications[i]
				if notification.Type == api.NotificationAcknowledgement {
					continue
				}
				if api.AlertType(opts.notificationType) == notification.Type {
					notifications[i] = n.updateIncidentNotification(notification)
					updated = true
					break
				}
			}
			if !updated {
				notifications = n.appendIncidentNotification(notifications)
			}
		}

		incident.Status.LastNotificationType = api.AlertType(opts.notificationType)
		incident.Status.Notifications = notifications

		if api.AlertType(opts.notificationType) == api.NotificationRecovery {
			_, _, err = util.PatchIncident(context.TODO(), n.extClient, incident, func(in *api.Incident) *api.Incident {
				if in.Labels == nil {
					in.Labels = map[string]string{}
				}
				in.Labels[api.LabelKeyProblemRecovered] = "true"
				return in
			}, metav1.PatchOptions{})
			if err != nil {
				return err
			}
		}

		_, err = util.UpdateIncidentStatus(context.TODO(), n.extClient, incident.ObjectMeta, func(in *api.IncidentStatus) (types.UID, *api.IncidentStatus) {
			in.LastNotificationType = api.AlertType(opts.notificationType)
			in.Notifications = notifications
			return incident.UID, in
		}, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	} else {
		name, err := n.generateIncidentName()
		if err != nil {
			return err
		}

		incident := &api.Incident{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: opts.host.AlertNamespace,
				Labels:    n.getLabel(),
			},
			Status: api.IncidentStatus{
				LastNotificationType: api.AlertType(opts.notificationType),
				Notifications:        n.appendIncidentNotification(make([]api.IncidentNotification, 0)),
			},
		}

		if _, err = n.extClient.Incidents(incident.Namespace).Create(context.TODO(), incident, metav1.CreateOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func (n *notifier) getIncident() (*api.Incident, error) {
	incidentList, err := n.extClient.Incidents(n.options.host.AlertNamespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(n.getLabel()).String(),
	})
	if err != nil {
		return nil, err
	}

	var lastCreationTimestamp time.Time
	var incident *api.Incident

	for _, item := range incidentList.Items {
		if item.CreationTimestamp.After(lastCreationTimestamp) {
			lastCreationTimestamp = item.CreationTimestamp.Time
			incident = &item
		}
	}
	return incident, nil
}

func (n *notifier) getLastNonOKState(incident *api.Incident) string {
	var lastTimestamp time.Time
	var lastNonOKState string

	for _, item := range incident.Status.Notifications {
		if item.LastTimestamp.After(lastTimestamp) {
			lastTimestamp = item.LastTimestamp.Time
			if item.LastState == stateCritical || item.LastState == stateWarning {
				lastNonOKState = item.LastState
			}
		}
	}
	return lastNonOKState
}
