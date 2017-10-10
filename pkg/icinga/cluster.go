package icinga

import (
	"github.com/appscode/go/errors"
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

func (h *ClusterHost) Create(alert api.ClusterAlert) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert.Namespace)

	if has, err := h.CheckIcingaService(alert.Name, kh); err != nil || has {
		return err
	}

	if err := h.CreateIcingaHost(kh); err != nil {
		return errors.FromErr(err).Err()
	}

	attrs := make(map[string]interface{})
	attrs["check_command"] = alertSpec.Check
	if alertSpec.CheckInterval.Seconds() > 0 {
		attrs["check_interval"] = alertSpec.CheckInterval.Seconds()
	}
	commandVars := api.ClusterCommands[alertSpec.Check].Vars
	for key, val := range alertSpec.Vars {
		if _, found := commandVars[key]; found {
			attrs[IVar(key)] = val
		}
	}
	if err := h.CreateIcingaService(alert.Name, kh, attrs); err != nil {
		return errors.FromErr(err).Err()
	}
	return h.CreateIcingaNotification(alert, kh)
}

func (h *ClusterHost) Update(alert api.ClusterAlert) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert.Namespace)

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
	if err := h.UpdateIcingaService(alert.Name, kh, attrs); err != nil {
		return errors.FromErr(err).Err()
	}

	return h.UpdateIcingaNotification(alert, kh)
}

func (h *ClusterHost) Delete(namespace, name string) error {
	kh := h.getHost(namespace)
	if err := h.DeleteIcingaService(name, kh); err != nil {
		return errors.FromErr(err).Err()
	}
	return h.DeleteIcingaHost(kh)
}
