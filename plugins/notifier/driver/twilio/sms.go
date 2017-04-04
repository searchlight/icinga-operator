package twilio

import (
	"strings"

	api "github.com/appscode/api/kubernetes/v1beta1"
	"github.com/appscode/go-notify/twilio"
	"github.com/appscode/searchlight/plugins/notifier/driver"
	"github.com/appscode/searchlight/plugins/notifier/driver/extpoints"
)

type biblio struct{}

func init() {
	extpoints.Drivers.Register(new(biblio), twilio.UID)
}

func (b *biblio) Notify(req *api.IncidentNotifyRequest) error {
	parts := strings.Split(req.HostName, "@")
	namespace := parts[1]

	alert, err := driver.GetAlertInfo(namespace, req.KubernetesAlertName)
	if err != nil {
		return err
	}

	smsBody, err := driver.RenderSMS(alert, req)
	if err != nil {
		return err
	}

	client, err := twilio.Default()
	if err != nil {
		return err
	}

	return client.WithBody(smsBody).Send()
}
