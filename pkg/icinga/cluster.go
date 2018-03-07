package icinga

import (
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
)

type ClusterHost struct {
	commonHost
}

func NewClusterHost(IcingaClient *Client) *ClusterHost {
	return &ClusterHost{
		commonHost: commonHost{
			IcingaClient: IcingaClient,
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

	if err := h.EnsureIcingaHost(kh); err != nil {
		return err
	}

	has, err := h.CheckIcingaService(alert.Name, kh)
	if err != nil {
		return err
	}

	attrs := make(map[string]interface{})
	if alertSpec.CheckInterval.Seconds() > 0 {
		attrs["check_interval"] = alertSpec.CheckInterval.Seconds()
	}
	commandVars := api.ClusterCommands[alertSpec.Check].Vars
	for key, val := range alertSpec.Vars {
		if _, found := commandVars[key]; found {
			attrs[IVar(key)] = val
		}
	}

	if !has {
		attrs["check_command"] = alertSpec.Check
		if err := h.CreateIcingaService(alert.Name, kh, attrs); err != nil {
			return err
		}
	} else {
		if err := h.UpdateIcingaService(alert.Name, kh, attrs); err != nil {
			return err
		}
	}

	return h.ReconcileIcingaNotification(alert, kh)
}

func (h *ClusterHost) Delete(namespace, name string) error {
	kh := h.getHost(namespace)
	if err := h.DeleteIcingaService(name, kh); err != nil {
		return err
	}
	return h.DeleteIcingaHost(kh)
}
