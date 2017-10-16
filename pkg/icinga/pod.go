package icinga

import (
	"bytes"
	"text/template"

	"github.com/appscode/go/errors"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type PodHost struct {
	commonHost

	KubeClient kubernetes.Interface
	ExtClient  cs.MonitoringV1alpha1Interface
}

func NewPodHost(kubeClient kubernetes.Interface, extClient cs.MonitoringV1alpha1Interface, IcingaClient *Client) *PodHost {
	return &PodHost{
		KubeClient: kubeClient,
		ExtClient:  extClient,
		commonHost: commonHost{
			IcingaClient: IcingaClient,
		},
	}
}

func (h *PodHost) getHost(alert api.PodAlert, pod apiv1.Pod) IcingaHost {
	return IcingaHost{
		ObjectName:     pod.Name,
		Type:           TypePod,
		AlertNamespace: alert.Namespace,
		IP:             pod.Status.PodIP,
	}
}

func (h *PodHost) expandVars(alertSpec api.PodAlertSpec, kh IcingaHost, attrs map[string]interface{}) error {
	commandVars := api.PodCommands[alertSpec.Check].Vars
	for key, val := range alertSpec.Vars {
		if v, found := commandVars[key]; found {
			if v.Parameterized {
				type Data struct {
					PodName   string
					PodIP     string
					Namespace string
				}
				tmpl, err := template.New("").Parse(val)
				if err != nil {
					return err
				}
				var buf bytes.Buffer
				err = tmpl.Execute(&buf, Data{PodName: kh.ObjectName, Namespace: kh.AlertNamespace, PodIP: kh.IP})
				if err != nil {
					return err
				}
				attrs[IVar(key)] = buf.String()
			} else {
				attrs[IVar(key)] = val
			}
		} else {
			return errors.Newf("variable %v not found", key).Err()
		}
	}
	return nil
}

func (h *PodHost) Create(alert api.PodAlert, pod apiv1.Pod) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert, pod)

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
	if err := h.expandVars(alertSpec, kh, attrs); err != nil {
		return err
	}
	if err := h.CreateIcingaService(alert.Name, kh, attrs); err != nil {
		return errors.FromErr(err).Err()
	}

	return h.CreateIcingaNotification(alert, kh)
}

func (h *PodHost) Update(alert api.PodAlert, pod apiv1.Pod) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert, pod)

	attrs := make(map[string]interface{})
	if alertSpec.CheckInterval.Seconds() > 0 {
		attrs["check_interval"] = alertSpec.CheckInterval.Seconds()
	}
	if err := h.expandVars(alertSpec, kh, attrs); err != nil {
		return err
	}
	if err := h.UpdateIcingaService(alert.Name, kh, attrs); err != nil {
		return errors.FromErr(err).Err()
	}

	return h.UpdateIcingaNotification(alert, kh)
}

func (h *PodHost) Delete(alert api.PodAlert, pod apiv1.Pod) error {
	kh := h.getHost(alert, pod)

	if err := h.DeleteIcingaService(alert.Name, kh); err != nil {
		return errors.FromErr(err).Err()
	}
	return h.DeleteIcingaHost(kh)
}
