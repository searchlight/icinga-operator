package icinga

import (
	"bytes"
	"text/template"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
)

type PodHost struct {
	commonHost
}

func NewPodHost(IcingaClient *Client, verbosity string) *PodHost {
	return &PodHost{
		commonHost: commonHost{
			IcingaClient: IcingaClient,
			verbosity:    verbosity,
		},
	}
}

func (h *PodHost) getHost(namespace string, pod *core.Pod) IcingaHost {
	return IcingaHost{
		ObjectName:     pod.Name,
		Type:           TypePod,
		AlertNamespace: namespace,
		IP:             pod.Status.PodIP,
	}
}

func (h *PodHost) expandVars(alertSpec api.PodAlertSpec, kh IcingaHost, attrs map[string]interface{}) error {
	cmd, _ := api.PodCommands.Get(alertSpec.Check)
	commandVars := cmd.Vars
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
			return errors.Errorf("variable %v not found", key)
		}
	}
	return nil
}

func (h *PodHost) Apply(alert *api.PodAlert, pod *core.Pod) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert.Namespace, pod)

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
	if err := h.expandVars(alertSpec, kh, attrs); err != nil {
		return err
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

func (h *PodHost) Delete(alertNamespace, alertName string, pod *core.Pod) error {
	kh := h.getHost(alertNamespace, pod)

	if err := h.deleteIcingaService(alertName, kh); err != nil {
		return err
	}
	return h.deleteIcingaHost(kh)
}

func (h *PodHost) DeleteChecks(cmd string) error {
	return h.deleteIcingaServiceForCheckCommand(cmd)
}
