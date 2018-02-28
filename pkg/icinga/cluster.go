package icinga

import (
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/pkg/errors"
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

func (h *ClusterHost) Create(alert *api.ClusterAlert) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert.Namespace)

	if err := h.EnsureIcingaHost(kh); err != nil {
		return errors.WithStack(err)
	}

	has, err := h.CheckIcingaService(alert.Name, kh)
	if err != nil {
		return errors.WithStack(err)
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
			return errors.WithStack(err)
		}
	} else {
		if err := h.UpdateIcingaService(alert.Name, kh, attrs); err != nil {
			return errors.WithStack(err)
		}
	}

	return h.EnsureIcingaNotification(alert, kh)
}

func (h *ClusterHost) Delete(alert *api.ClusterAlert) error {
	kh := h.getHost(alert.Namespace)
	if err := h.DeleteIcingaService(alert.Name, kh); err != nil {
		return errors.WithStack(err)
	}
	return h.DeleteIcingaHost(kh)
}
