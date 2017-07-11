package icinga

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/appscode/errors"
	tapi "github.com/appscode/searchlight/api"
)

type commonHost struct {
	IcingaClient *Client
}

func (h *commonHost) CreateIcingaHost(kh IcingaHost) error {
	host, err := kh.Name()
	if err != nil {
		return err
	}
	resp := h.IcingaClient.Objects().Hosts(host).Get([]string{}).Do()
	if resp.Status == 200 {
		return nil
	}
	obj := IcingaObject{
		Templates: []string{"generic-host"},
		Attrs: map[string]interface{}{
			"address": kh.IP,
		},
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.FromErr(err).Err()
	}

	resp = h.IcingaClient.Objects().Hosts(host).Create([]string{}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.FromErr(resp.Err).Err()
	}
	if resp.Status != 200 {
		return errors.New("Can't create Icinga host").Err()
	}
	return nil
}

func (h *commonHost) DeleteIcingaHost(kh IcingaHost) error {
	param := map[string]string{
		"cascade": "1",
	}
	host, err := kh.Name()
	if err != nil {
		return err
	}

	in := fmt.Sprintf(`{"filter": "match(\"%s\",host.name)"}`, host)
	var respService ResponseObject
	if _, err := h.IcingaClient.Objects().Service("").Update([]string{}, in).Do().Into(&respService); err != nil {
		return errors.New("Can't get Icinga service").Err()
	}

	if len(respService.Results) <= 1 {
		resp := h.IcingaClient.Objects().Hosts("").Delete([]string{}, in).Params(param).Do()
		if resp.Err != nil {
			return errors.New("Can't delete Icinga host").Err()
		}
	}
	return nil
}

// createIcingaServiceForCluster
func (h *commonHost) CreateIcingaService(svc string, kh IcingaHost, attrs map[string]interface{}) error {
	obj := IcingaObject{
		Templates: []string{"generic-service"},
		Attrs:     attrs,
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.FromErr(err).Err()
	}
	host, err := kh.Name()
	if err != nil {
		return err
	}
	resp := h.IcingaClient.Objects().Service(host).Create([]string{svc}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.FromErr(resp.Err).Err()
	}
	if resp.Status == 200 {
		return nil
	}
	if strings.Contains(string(resp.ResponseBody), "already exists") {
		return nil
	}
	return errors.New("Can't create Icinga service").Err()
}

func (h *commonHost) UpdateIcingaService(svc string, kh IcingaHost, attrs map[string]interface{}) error {
	obj := IcingaObject{
		Templates: []string{"generic-service"},
		Attrs:     attrs,
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.FromErr(err).Err()
	}
	host, err := kh.Name()
	if err != nil {
		return err
	}
	resp := h.IcingaClient.Objects().Service(host).Update([]string{svc}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.FromErr(resp.Err).Err()
	}
	if resp.Status != 200 {
		return errors.New("Can't update Icinga service").Err()
	}
	return nil
}

func (h *commonHost) DeleteIcingaService(svc string, kh IcingaHost) error {
	param := map[string]string{
		"cascade": "1",
	}
	in := h.IcingaServiceSearchQuery(svc, kh)
	resp := h.IcingaClient.Objects().Service("").Delete([]string{}, in).Params(param).Do()

	if resp.Err != nil {
		return errors.FromErr(resp.Err).Err()
	}
	if resp.Status == 200 {
		return nil
	}
	return errors.New("Fail to delete service").Err()
}

func (h *commonHost) CheckIcingaService(svc string, kh IcingaHost) (bool, error) {
	in := h.IcingaServiceSearchQuery(svc, kh)
	var respService ResponseObject

	if _, err := h.IcingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
		return true, errors.New("Can't check icinga service").Err()
	}
	return len(respService.Results) > 0, nil
}

func (h *commonHost) IcingaServiceSearchQuery(svc string, kids ...IcingaHost) string {
	matchHost := ""
	for i, kh := range kids {
		if i > 0 {
			matchHost = matchHost + "||"
		}
		host, _ := kh.Name()
		matchHost = matchHost + fmt.Sprintf(`match(\"%s\",host.name)`, host)
	}
	return fmt.Sprintf(`{"filter": "(%s)&&match(\"%s\",service.name)"}`, matchHost, svc)
}

func (h *commonHost) CreateIcingaNotification(alert tapi.Alert, kh IcingaHost) error {
	obj := IcingaObject{
		Templates: []string{"icinga2-notifier-template"},
		Attrs: map[string]interface{}{
			"interval": int(alert.GetAlertInterval().Seconds()),
			"users":    []string{"appscode_user"},
		},
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.FromErr(err).Err()
	}
	host, err := kh.Name()
	if err != nil {
		return err
	}
	resp := h.IcingaClient.Objects().
		Notifications(host).
		Create([]string{alert.GetName(), alert.GetName()}, string(jsonStr)).
		Do()
	if resp.Err != nil {
		return errors.FromErr(resp.Err).Err()
	}
	if resp.Status == 200 || strings.Contains(string(resp.ResponseBody), "already exists") {
		return nil
	}
	return errors.New("Can't create Icinga notification").Err()
}

func (h *commonHost) UpdateIcingaNotification(alert tapi.Alert, kh IcingaHost) error {
	obj := IcingaObject{
		Attrs: map[string]interface{}{
			"interval": int(alert.GetAlertInterval().Seconds()),
		},
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.FromErr(err).Err()
	}
	host, err := kh.Name()
	if err != nil {
		return err
	}
	resp := h.IcingaClient.Objects().
		Notifications(host).
		Update([]string{alert.GetName(), alert.GetName()}, string(jsonStr)).
		Do()
	if resp.Err != nil {
		return errors.FromErr(resp.Err).Err()
	}
	if resp.Status != 200 {
		return errors.New("Can't update Icinga notification").Err()
	}
	return nil
}
