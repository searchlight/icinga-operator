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

package icinga

import (
	"encoding/json"
	"fmt"
	"strings"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"

	"github.com/pkg/errors"
)

type commonHost struct {
	IcingaClient *Client
	// V logging level, the value of the -v flag
	verbosity string
}

func (h *commonHost) Complete(v string) {
	h.verbosity = v
}

func (h *commonHost) reconcileIcingaHost(kh IcingaHost) error {
	host, err := kh.Name()
	if err != nil {
		return errors.WithStack(err)
	}

	obj := IcingaObject{
		Templates: []string{"generic-host"},
		Attrs: map[string]interface{}{
			"address":         kh.IP,
			IVar("verbosity"): h.verbosity,
		},
	}
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "Failed to Marshal IcingaObject")
	}

	resp := h.IcingaClient.Hosts(host).Create([]string{}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, string(resp.ResponseBody))
	}
	if resp.Status == 200 {
		return nil
	}

	if !strings.Contains(string(resp.ResponseBody), "already exists") {
		return errors.Errorf("Failed to create Icinga Host. Status: %d", resp.Status)
	}

	resp = h.IcingaClient.Hosts(host).Update([]string{}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to update Icinga Host")
	}
	if resp.Status != 200 {
		return errors.Errorf("can't update Icinga Host. Status: %d", resp.Status)
	}

	return nil
}

func (h *commonHost) deleteIcingaHost(kh IcingaHost) error {
	param := map[string]string{
		"cascade": "1",
	}
	host, err := kh.Name()
	if err != nil {
		return errors.WithStack(err)
	}

	in := fmt.Sprintf(`{"filter": "match(\"%s\",host.name)"}`, host)
	var respService ResponseObject
	if _, err := h.IcingaClient.Service("").Update([]string{}, in).Do().Into(&respService); err != nil {
		return errors.Wrap(err, "can't get Icinga service")
	}

	if len(respService.Results) == 0 {
		resp := h.IcingaClient.Hosts("").Delete([]string{}, in).Params(param).Do()
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
	resp := h.IcingaClient.Hosts("").Delete([]string{}, in).Params(param).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to delete IcingaHost")
	}
	if resp.Status == 200 {
		return nil
	}
	return errors.New("can't delete Icinga host")
}

// createIcingaServiceForCluster
func (h *commonHost) createIcingaService(svc string, kh IcingaHost, attrs map[string]interface{}) error {
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
	resp := h.IcingaClient.Service(host).Create([]string{svc}, string(jsonStr)).Do()
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

func (h *commonHost) updateIcingaService(svc string, kh IcingaHost, attrs map[string]interface{}) error {
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
	resp := h.IcingaClient.Service(host).Update([]string{svc}, string(jsonStr)).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to update Icinga Service")
	}
	if resp.Status != 200 {
		return errors.Errorf("can't update Icinga service; %d", resp.Status)
	}
	return nil
}

func (h *commonHost) deleteIcingaService(svc string, kh IcingaHost) error {
	param := map[string]string{
		"cascade": "1",
	}
	in := h.IcingaServiceSearchQuery(svc, kh)

	resp := h.IcingaClient.Service("").Delete([]string{}, in).Params(param).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to delete Icinga Service")
	}
	if resp.Status == 200 || resp.Status == 404 {
		return nil
	}

	return errors.Errorf("Fail to delete service. Status: %d", resp.Status)
}

func (h *commonHost) deleteIcingaServiceForCheckCommand(name string) error {
	param := map[string]string{
		"cascade": "1",
	}
	in := fmt.Sprintf(`{"filter": "match(\"%s\",service.check_command)"}`, name)

	resp := h.IcingaClient.Service("").Delete([]string{}, in).Params(param).Do()
	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to delete Icinga Service")
	}
	if resp.Status == 200 || resp.Status == 404 {
		return nil
	}

	return errors.Errorf("Fail to delete service. Status: %d", resp.Status)
}

func (h *commonHost) checkIcingaService(svc string, kh IcingaHost) (bool, error) {
	in := h.IcingaServiceSearchQuery(svc, kh)
	var respService ResponseObject

	if _, err := h.IcingaClient.Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
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

func (h *commonHost) reconcileIcingaNotification(alert api.Alert, kh IcingaHost) error {
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

	resp = h.IcingaClient.Notifications(host).Create([]string{alert.GetName(), alert.GetName()}, string(jsonStr)).Do()

	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to create Icinga Notification")
	}
	if resp.Status == 200 {
		return nil
	}

	if !strings.Contains(string(resp.ResponseBody), "already exists") {
		return errors.Errorf("Failed to create Icinga notification. Status: %d", resp.Status)
	}

	resp = h.IcingaClient.Notifications(host).Update([]string{alert.GetName(), alert.GetName()}, string(jsonStr)).Do()

	if resp.Err != nil {
		return errors.Wrap(resp.Err, "Failed to update Icinga Notification")
	}

	if resp.Status != 200 {
		return errors.Errorf("can't update Icinga notification. Status: %d", resp.Status)
	}

	return nil
}
