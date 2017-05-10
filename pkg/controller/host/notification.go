package host

import (
	"encoding/json"
	"strings"

	"github.com/appscode/errors"
	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/searchlight/pkg/client/icinga"
)

func CreateIcingaNotification(icingaClient *icinga.IcingaClient, alert *aci.Alert, objectList []*KubeObjectInfo) error {
	alertSpec := alert.Spec
	for _, object := range objectList {
		var obj IcingaObject
		obj.Templates = []string{"icinga2-notifier-template"}
		mp := make(map[string]interface{})
		mp["interval"] = alertSpec.IcingaParam.AlertIntervalSec
		mp["users"] = []string{"appscode_user"}
		obj.Attrs = mp

		jsonStr, err := json.Marshal(obj)
		if err != nil {
			return errors.New().WithCause(err).Err()
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

func UpdateIcingaNotification(icingaClient *icinga.IcingaClient, alert *aci.Alert, objectList []*KubeObjectInfo) error {
	icignaService := alert.Name
	for _, object := range objectList {
		var obj IcingaObject
		mp := make(map[string]interface{})
		mp["interval"] = alert.Spec.IcingaParam.AlertIntervalSec
		obj.Attrs = mp
		jsonStr, err := json.Marshal(obj)
		if err != nil {
			return errors.New().WithCause(err).Err()
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
