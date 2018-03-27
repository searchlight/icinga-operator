package notifier

import (
	"fmt"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
)

type SMS struct {
	AlertName        string
	NotificationType string
	ServiceState     string
	Author           string
	Comment          string
	Hostname         string
}

func (n *notifier) RenderSMS(receiver api.Receiver) string {
	opts := n.options
	m := &SMS{
		AlertName:        opts.alertName,
		NotificationType: opts.notificationType,
		ServiceState:     receiver.State,
		Author:           opts.author,
		Comment:          opts.comment,
		Hostname:         opts.hostname,
	}

	return m.Render()
}

func (m *SMS) Render() string {
	var msg string
	switch api.AlertType(m.NotificationType) {
	case api.NotificationAcknowledgement:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nThis issue is acked.", m.AlertName, m.Hostname, m.ServiceState)
	case api.NotificationRecovery:
		msg = fmt.Sprintf("Service [%s] for [%s] was in \"%s\" state.\nThis issue is recovered.", m.AlertName, m.Hostname, m.ServiceState)
	case api.NotificationProblem:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.\nCheck this issue in Icingaweb.", m.AlertName, m.Hostname, m.ServiceState)
	default:
		msg = fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.", m.AlertName, m.Hostname, m.ServiceState)
	}
	if m.Comment != "" {
		if m.Author != "" {
			msg = msg + " " + fmt.Sprintf(`%s says "%s".`, m.Author, m.Comment)
		} else {
			msg = msg + " " + fmt.Sprintf(`Comment: "%s".`, m.Comment)
		}
	}
	return msg
}
