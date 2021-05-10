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
	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
)

type ClusterHost struct {
	commonHost
}

func NewClusterHost(IcingaClient *Client, verbosity string) *ClusterHost {
	return &ClusterHost{
		commonHost: commonHost{
			IcingaClient: IcingaClient,
			verbosity:    verbosity,
		},
	}
}

func (h *ClusterHost) getHost(namespace string) IcingaHost {
	return IcingaHost{
		Type:           TypeCluster,
		AlertNamespace: namespace,
		IP:             "127.0.0.1",
	}
}

func (h *ClusterHost) Apply(alert *api.ClusterAlert) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert.Namespace)

	if err := h.reconcileIcingaHost(kh); err != nil {
		return err
	}

	has, err := h.checkIcingaService(alert.Name, kh)
	if err != nil {
		return err
	}

	if alertSpec.Paused {
		if has {
			if err := h.deleteIcingaService(alert.Name, kh); err != nil {
				return err
			}
		}
		return nil
	}

	attrs := make(map[string]interface{})
	if alertSpec.CheckInterval.Seconds() > 0 {
		attrs["check_interval"] = alertSpec.CheckInterval.Seconds()
	}
	cmd, _ := api.ClusterCommands.Get(alertSpec.Check)
	commandVars := cmd.Vars.Fields
	for key, val := range alertSpec.Vars {
		if _, found := commandVars[key]; found {
			attrs[IVar(key)] = val
		}
	}

	if !has {
		attrs["check_command"] = alertSpec.Check
		if err := h.createIcingaService(alert.Name, kh, attrs); err != nil {
			return err
		}
	} else {
		if err := h.updateIcingaService(alert.Name, kh, attrs); err != nil {
			return err
		}
	}

	return h.reconcileIcingaNotification(alert, kh)
}

func (h *ClusterHost) Delete(namespace, name string) error {
	kh := h.getHost(namespace)
	if err := h.deleteIcingaService(name, kh); err != nil {
		return err
	}
	return h.deleteIcingaHost(kh)
}

func (h *ClusterHost) DeleteChecks(cmd string) error {
	return h.deleteIcingaServiceForCheckCommand(cmd)
}
