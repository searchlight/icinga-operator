package mailgun

import (
	"strings"

	api "github.com/appscode/api/kubernetes/v1beta1"
	"github.com/appscode/go-notify/mailgun"
	"github.com/appscode/searchlight/plugins/notifier/driver"
	"github.com/appscode/searchlight/plugins/notifier/driver/extpoints"
)

type biblio struct{}

func init() {
	extpoints.Drivers.Register(new(biblio), mailgun.UID)
}

func (b *biblio) Notify(req *api.IncidentNotifyRequest) error {
	parts := strings.Split(req.HostName, "@")
	namespace := parts[1]

	alert, err := driver.GetAlertInfo(namespace, req.KubernetesAlertName)
	if err != nil {
		return err
	}

	mailBody, err := driver.RenderMail(alert, req)
	if err != nil {
		return err
	}

	client, err := mailgun.Default()
	if err != nil {
		return err
	}

	subject := "Notification"
	if sub, found := driver.SubjectMap[req.Type]; found {
		subject = sub
	}

	return client.WithSubject(subject).
		WithBody(mailBody).
		SendHtml()
}
