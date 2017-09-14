package icinga

import (
	"github.com/appscode/errors"
	tapi "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	tcs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	clientset "k8s.io/client-go/kubernetes"
)

type ClusterHost struct {
	commonHost

	KubeClient clientset.Interface
	ExtClient  tcs.MonitoringV1alpha1Interface
}

func NewClusterHost(kubeClient clientset.Interface, extClient tcs.MonitoringV1alpha1Interface, IcingaClient *Client) *ClusterHost {
	return &ClusterHost{
		KubeClient: kubeClient,
		ExtClient:  extClient,
		commonHost: commonHost{
			IcingaClient: IcingaClient,
		},
	}
}

func (h *ClusterHost) getHost(alert tapi.ClusterAlert) IcingaHost {
	return IcingaHost{
		Type:           TypeCluster,
		AlertNamespace: alert.Namespace,
		IP:             "127.0.0.1",
	}
}

func (h *ClusterHost) Create(alert tapi.ClusterAlert) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert)

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
	commandVars := tapi.ClusterCommands[alertSpec.Check].Vars
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

func (h *ClusterHost) Update(alert tapi.ClusterAlert) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert)

	attrs := make(map[string]interface{})
	if alertSpec.CheckInterval.Seconds() > 0 {
		attrs["check_interval"] = alertSpec.CheckInterval.Seconds()
	}
	commandVars := tapi.ClusterCommands[alertSpec.Check].Vars
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

func (h *ClusterHost) Delete(alert tapi.ClusterAlert) error {
	kh := h.getHost(alert)
	if err := h.DeleteIcingaService(alert.Name, kh); err != nil {
		return errors.FromErr(err).Err()
	}
	return h.DeleteIcingaHost(kh)
}
