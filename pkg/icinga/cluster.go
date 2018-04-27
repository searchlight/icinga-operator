package icinga

import (
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
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
	commandVars := cmd.Vars.Items
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
