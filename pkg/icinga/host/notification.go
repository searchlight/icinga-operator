package host

import (
	"encoding/json"
	"strings"

	"github.com/appscode/errors"
	aci "github.com/appscode/searchlight/api"
	icinga "github.com/appscode/searchlight/pkg/icinga/client"
)

func CreateIcingaNotification(icingaClient *icinga.IcingaClient, alert *aci.PodAlert, objectList []*KubeObjectInfo) error {
	alertSpec := alert.Spec
	for _, object := range objectList {
		var obj IcingaObject
		obj.Templates = []string{"icinga2-notifier-template"}
		mp := make(map[string]interface{})
		mp["interval"] = alertSpec.AlertInterval
		mp["users"] = []string{"appscode_user"}
		obj.Attrs = mp

		jsonStr, err := json.Marshal(obj)
		if err != nil {
			return errors.FromErr(err).Err()
		}

		resp := icingaClient.Objects().Notifications(object.Name).Create([]string{alert.Name, alert.Name}, string(jsonStr)).Do()
		if resp.Err != nil {
			return errors.New().WithCause(resp.Err).Err()
		}
		if resp.Status == 200 {
			continue
		}
		if strings.Contains(string(resp.ResponseBody), "already exists") {
			continue
		}

		return errors.New("Can't create Icinga notification").Err()
	}
	return nil
}

func UpdateIcingaNotification(icingaClient *icinga.IcingaClient, alert *aci.PodAlert, objectList []*KubeObjectInfo) error {
	icignaService := alert.Name
	for _, object := range objectList {
		var obj IcingaObject
		mp := make(map[string]interface{})
		mp["interval"] = alert.Spec.AlertInterval
		obj.Attrs = mp
		jsonStr, err := json.Marshal(obj)
		if err != nil {
			return errors.FromErr(err).Err()
		}
		resp := icingaClient.Objects().Notifications(object.Name).Update([]string{icignaService, icignaService}, string(jsonStr)).Do()

		if resp.Err != nil {
			return errors.New().WithCause(resp.Err).Err()
		}
		if resp.Status != 200 {
			return errors.New("Can't update Icinga notification").Err()
		}
	}
	return nil
}
