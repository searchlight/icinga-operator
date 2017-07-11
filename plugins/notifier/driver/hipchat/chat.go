package hipchat

import (
	api "github.com/appscode/api/kubernetes/v1beta1"
	"github.com/appscode/go-notify/hipchat"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins/notifier/driver"
	"github.com/appscode/searchlight/plugins/notifier/driver/extpoints"
)

type biblio struct{}

func init() {
	extpoints.Drivers.Register(new(biblio), hipchat.UID)
}

func (b *biblio) Notify(req *api.IncidentNotifyRequest) error {
	host, err := icinga.ParseHost(req.HostName)
	if err != nil {
		return err
	}

	alert, err := driver.GetAlertInfo(host.AlertNamespace, req.KubernetesAlertName)
	if err != nil {
		return err
	}

	chatBody, err := driver.RenderSMS(alert, req)
	if err != nil {
		return err
	}

	client, err := hipchat.Default()
	if err != nil {
		return err
	}

	return client.WithBody(chatBody).Send()
}
