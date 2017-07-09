package icinga

import (
	"fmt"
	"regexp"

	"github.com/appscode/errors"
	tapi "github.com/appscode/searchlight/api"
	tcs "github.com/appscode/searchlight/client/clientset"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type PodHost struct {
	commonHost

	KubeClient clientset.Interface
	ExtClient  tcs.ExtensionInterface
}

func NewPodHost(kubeClient clientset.Interface, extClient tcs.ExtensionInterface, IcingaClient *Client) *PodHost {
	return &PodHost{
		KubeClient: kubeClient,
		ExtClient:  extClient,
		commonHost: commonHost{
			IcingaClient: IcingaClient,
		},
	}
}

func (h *PodHost) GetObject(alert tapi.PodAlert, pod apiv1.Pod) KHost {
	return KHost{Name: pod.Name + "@" + alert.Namespace, IP: pod.Status.PodIP}
}

func (h *PodHost) expandVars(alertSpec tapi.PodAlertSpec, kh KHost, attrs map[string]interface{}) error {
	commandVars := tapi.PodCommands[alertSpec.Check].Vars
	for key, val := range alertSpec.Vars {
		if v, found := commandVars[key]; found {
			if v.Parameterized {
				reg, err := regexp.Compile("pod_name[ ]*=[ ]*'[?]'")
				if err != nil {
					return err
				}
				attrs[IVar(key)] = reg.ReplaceAllString(val.(string), fmt.Sprintf("pod_name='%s'", kh.Name))
			} else {
				attrs[IVar(key)] = val
			}
		} else {
			return errors.Newf("variable %v not found", key).Err()
		}
	}
	return nil
}

func (h *PodHost) Create(alert tapi.PodAlert, pod apiv1.Pod) error {
	alertSpec := alert.Spec
	kh := h.GetObject(alert, pod)

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

func (h *PodHost) Update(alert tapi.PodAlert, pod apiv1.Pod) error {
	alertSpec := alert.Spec
	kh := h.GetObject(alert, pod)

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

func (h *PodHost) Delete(alert tapi.PodAlert, pod apiv1.Pod) error {
	kh := h.GetObject(alert, pod)

	if err := h.DeleteIcingaService(alert.Name, kh); err != nil {
		return errors.FromErr(err).Err()
	}
	return h.DeleteIcingaHost(kh.Name)
}
