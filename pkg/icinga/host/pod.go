package host

import (
	"fmt"
	"regexp"

	"github.com/appscode/errors"
	tapi "github.com/appscode/searchlight/api"
	tcs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/data"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type PodHost struct {
	commonHost

	KubeClient clientset.Interface
	ExtClient  tcs.ExtensionInterface
}

func (h *PodHost) GetObject(alert tapi.PodAlert, pod apiv1.Pod) KHost {
	return KHost{Name: pod.Name + "@" + alert.Namespace, IP: pod.Status.PodIP}
}

// set Alert in Icinga LocalHost
func (h *PodHost) Create(alert tapi.PodAlert, pod apiv1.Pod, specificObject string) error {
	alertSpec := alert.Spec

	if alertSpec.Check == "" {
		return errors.New("Invalid request").Err()
	}

	// Get Icinga Host Info
	objectList, err := GetObjectList(h.KubeClient, alertSpec.Check, HostTypePod, alert.Namespace, h.ObjectType, h.ObjectName, specificObject)
	if err != nil {
		return errors.FromErr(err).Err()
	}

	var has bool
	if has, err = h.CheckIcingaService(alert.Name, objectList); err != nil {
		return errors.FromErr(err).Err()
	}
	if has {
		return nil
	}

	// Create Icinga Host
	if err := h.CreateIcingaHost(objectList, alert.Namespace); err != nil {
		return errors.FromErr(err).Err()
	}

	if err := h.createIcingaService(objectList); err != nil {
		return errors.FromErr(err).Err()
	}

	if err := h.CreateIcingaNotification(alert, objectList); err != nil {
		return errors.FromErr(err).Err()
	}

	return nil
}

func (h *PodHost) setParameterizedVariables(alert tapi.PodAlert, pod apiv1.Pod, alertSpec tapi.PodAlertSpec, objectName string, commandVars map[string]data.CommandVar, mp map[string]interface{}) (map[string]interface{}, error) {
	for key, val := range alertSpec.Vars {
		if v, found := commandVars[key]; found {
			if !v.Parameterized {
				continue
			}
			reg, err := regexp.Compile("pod_name[ ]*=[ ]*'[?]'")
			if err != nil {
				return nil, errors.FromErr(err).Err()
			}
			mp[IVar(key)] = reg.ReplaceAllString(val.(string), fmt.Sprintf("pod_name='%s'", objectName))
		} else {
			return nil, errors.Newf("variable %v not found", key).Err()
		}
	}
	return mp, nil
}

func (h *PodHost) createIcingaService(alert tapi.PodAlert, pod apiv1.Pod, objectList []*KHost) error {
	alertSpec := alert.Spec

	mp := make(map[string]interface{})
	mp["check_command"] = alertSpec.Check
	if alertSpec.CheckInterval.Seconds() > 0 {
		mp["check_interval"] = alertSpec.CheckInterval.Seconds()
	}

	commandVars := tapi.PodCommands[alertSpec.Check].Vars
	for key, val := range alertSpec.Vars {
		if v, found := commandVars[key]; found {
			if v.Parameterized {
				continue
			}
			mp[IVar(key)] = val
		}
	}

	for _, object := range objectList {
		var err error
		if mp, err = h.setParameterizedVariables(alertSpec, object.Name, commandVars, mp); err != nil {
			return errors.FromErr(err).Err()
		}

		if err := h.CreateIcingaService(mp, object, alert.Name); err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}

func (h *PodHost) Update(alert tapi.PodAlert, pod apiv1.Pod) error {
	alertSpec := alert.Spec

	// Get Icinga Host Info
	objectList, err := GetObjectList(h.KubeClient, alertSpec.Check, HostTypePod, alert.Namespace, h.ObjectType, h.ObjectName, "")
	if err != nil {
		return errors.FromErr(err).Err()
	}

	if err := h.updateIcingaService(objectList); err != nil {
		return errors.FromErr(err).Err()
	}

	if err := h.UpdateIcingaNotification(alert, objectList); err != nil {
		return errors.FromErr(err).Err()
	}
	return nil
}

func (h *PodHost) updateIcingaService(alert tapi.PodAlert, pod apiv1.Pod, objectList []*KHost) error {
	alertSpec := alert.Spec

	mp := make(map[string]interface{})
	if alertSpec.CheckInterval.Seconds() > 0 {
		mp["check_interval"] = alertSpec.CheckInterval.Seconds()
	}

	commandVars := tapi.PodCommands[alertSpec.Check].Vars
	for key, val := range alertSpec.Vars {
		if v, found := commandVars[key]; found {
			if v.Parameterized {
				continue
			}
			mp[IVar(key)] = val
		}
	}

	for _, object := range objectList {
		var err error
		if mp, err = h.setParameterizedVariables(alertSpec, object.Name, commandVars, mp); err != nil {
			return errors.FromErr(err).Err()
		}

		if err := h.UpdateIcingaService(mp, object, alert.Name); err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}

func (h *PodHost) Delete(alert tapi.PodAlert, pod apiv1.Pod, specificObject string) error {
	alertSpec := alert.Spec
	var objectList []*KHost
	if specificObject != "" {
		objectList = append(objectList, &KHost{Name: specificObject + "@" + alert.Namespace})
	} else {
		// Get Icinga Host Info
		var err error
		objectList, err = GetObjectList(h.KubeClient, alertSpec.Check, HostTypePod, alert.Namespace, h.ObjectType, h.ObjectName, specificObject)
		if err != nil {
			return errors.FromErr(err).Err()
		}
	}

	if err := h.DeleteIcingaService(objectList, alert.Name); err != nil {
		return errors.FromErr(err).Err()
	}

	for _, object := range objectList {
		if err := h.DeleteIcingaHost(object.Name); err != nil {
			return errors.FromErr(err).Err()
		}
	}

	return nil
}
