package icinga

import (
	"encoding/json"
	"fmt"
	"strings"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/pkg/errors"
)

type commonHost struct {
	IcingaClient *Client
}

func (h *commonHost) EnsureIcingaHost(kh IcingaHost) error {
	host, err := kh.Name()
	if err != nil {
		return errors.WithStack(err)
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
		return errors.Wrap(err, "Failed to Marshal IcingaObject")
	}

	resp = h.IcingaClient.Objects().Hosts(host).Create([]string{}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, string(resp.ResponseBody))
	}
	if resp.Status != 200 {
		return errors.Errorf("can't create Icinga host. Status: %d", resp.Status)
	}
	return nil
}

func (h *commonHost) DeleteIcingaHost(kh IcingaHost) error {
	param := map[string]string{
		"cascade": "1",
	}
	host, err := kh.Name()
	if err != nil {
		return errors.WithStack(err)
	}

	in := fmt.Sprintf(`{"filter": "match(\"%s\",host.name)"}`, host)
	var respService ResponseObject
	if _, err := h.IcingaClient.Objects().Service("").Update([]string{}, in).Do().Into(&respService); err != nil {
		return errors.Wrap(err, "can't get Icinga service")
	}

	if len(respService.Results) == 0 {
		resp := h.IcingaClient.Objects().Hosts("").Delete([]string{}, in).Params(param).Do()
		if resp.Err != nil {
			return errors.Wrap(err, "can't delete Icinga host")
		}
	}
	return nil
}

func (h *commonHost) ForceDeleteIcingaHost(kh IcingaHost) error {
	param := map[string]string{
		"cascade": "1",
	}
	host, err := kh.Name()
	if err != nil {
		return errors.WithStack(err)
	}

	in := fmt.Sprintf(`{"filter": "match(\"%s\",host.name)"}`, host)
	resp := h.IcingaClient.Objects().Hosts("").Delete([]string{}, in).Params(param).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to delete IcingaHost")
	}
	if resp.Status == 200 {
		return nil
	}
	return errors.New("can't delete Icinga host")
}

// createIcingaServiceForCluster
func (h *commonHost) CreateIcingaService(svc string, kh IcingaHost, attrs map[string]interface{}) error {
	obj := IcingaObject{
		Templates: []string{"generic-service"},
		Attrs:     attrs,
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "Failed to Marshal IcingaObject")
	}
	host, err := kh.Name()
	if err != nil {
		return errors.WithStack(err)
	}
	resp := h.IcingaClient.Objects().Service(host).Create([]string{svc}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to create Icinga Service")
	}
	if resp.Status == 200 {
		return nil
	}
	if strings.Contains(string(resp.ResponseBody), "already exists") {
		return nil
	}

	return errors.Errorf("can't create Icinga service. Status: %d", resp.Status)
}

func (h *commonHost) UpdateIcingaService(svc string, kh IcingaHost, attrs map[string]interface{}) error {
	obj := IcingaObject{
		Templates: []string{"generic-service"},
		Attrs:     attrs,
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "Failed to Marshal IcingaObject")
	}
	host, err := kh.Name()
	if err != nil {
		return errors.WithStack(err)
	}
	resp := h.IcingaClient.Objects().Service(host).Update([]string{svc}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to update Icinga Service")
	}
	if resp.Status != 200 {
		return errors.Errorf("can't update Icinga service; %d", resp.Status)
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
		return errors.Wrap(resp.Err, "Failed to delete Icinga Service")
	}
	if resp.Status == 200 || resp.Status == 404 {
		return nil
	}

	return errors.Errorf("Fail to delete service. Status: %d", resp.Status)
}

func (h *commonHost) CheckIcingaService(svc string, kh IcingaHost) (bool, error) {
	in := h.IcingaServiceSearchQuery(svc, kh)
	var respService ResponseObject

	if _, err := h.IcingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
		return true, errors.Wrap(err, "can't check icinga service")
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

func (h *commonHost) CheckIcingaNotification(svc string, kh IcingaHost) (bool, error) {
	in := h.IcingaServiceSearchQuery(svc, kh)
	var respService ResponseObject

	if _, err := h.IcingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
		return true, errors.Wrap(err, "can't check icinga service")
	}
	return len(respService.Results) > 0, nil
}

func (h *commonHost) ReconcileIcingaNotification(alert api.Alert, kh IcingaHost) error {
	obj := IcingaObject{
		Templates: []string{"icinga2-notifier-template"},
		Attrs: map[string]interface{}{
			"interval": int(alert.GetAlertInterval().Seconds()),
			"users":    []string{"searchlight_user"},
		},
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "Failed to Marshal Icinga Object")
	}
	host, err := kh.Name()
	if err != nil {
		return errors.WithStack(err)
	}

	var resp *APIResponse

	resp = h.IcingaClient.Objects().Notifications(host).Create([]string{alert.GetName(), alert.GetName()}, string(jsonStr)).Do()

	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to create Icinga Notification")
	}
	if resp.Status == 200 {
		return nil
	}

	if !strings.Contains(string(resp.ResponseBody), "already exists") {
		return errors.Errorf("Failed to create Icinga notification. Status: %d", resp.Status)
	}

	resp = h.IcingaClient.Objects().Notifications(host).Update([]string{alert.GetName(), alert.GetName()}, string(jsonStr)).Do()

	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to update Icinga Notification")
	}

	if resp.Status != 200 {
		return errors.Errorf("can't update Icinga notification. Status: %d", resp.Status)
	}

	return nil
}
