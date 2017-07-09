package host

import (
	"github.com/appscode/errors"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/controller/types"
)

type ClusterHost struct {
	*types.Option
}

//-----------------------------------------------------

// set Alert in Icinga LocalHost
func (h *ClusterHost) Create() error {
	alertSpec := h.Resource.Spec
	if alertSpec.Check == "" {
		return errors.New("Invalid request").Err()
	}

	// Get Icinga Host Info
	objectList, err := GetObjectList(h.KubeClient, alertSpec.Check, HostTypeLocalhost, h.Resource.Namespace, h.ObjectType, h.ObjectName, "")
	if err != nil {
		return errors.FromErr(err).Err()
	}

	var has bool
	if has, err = CheckIcingaService(h.IcingaClient, h.Resource.Name, objectList); err != nil {
		return errors.FromErr(err).Err()
	}
	if has {
		return nil
	}

	// Create Icinga Host
	if err := CreateIcingaHost(h.IcingaClient, objectList, h.Resource.Namespace); err != nil {
		return errors.FromErr(err).Err()
	}

	if err := h.createIcingaService(objectList); err != nil {
		return errors.FromErr(err).Err()
	}

	if err := CreateIcingaNotification(h.IcingaClient, h.Resource, objectList); err != nil {
		return errors.FromErr(err).Err()
	}

	return nil
}

func (h *ClusterHost) createIcingaService(objectList []*KubeObjectInfo) error {
	alertSpec := h.Resource.Spec

	mp := make(map[string]interface{})
	mp["check_command"] = alertSpec.Check
	if alertSpec.CheckInterval.Seconds() > 0 {
		mp["check_interval"] = alertSpec.CheckInterval.Seconds()
	}

	commandVars := tapi.ClusterCommands[alertSpec.Check].Vars

	for key, val := range alertSpec.Vars {
		if _, found := commandVars[key]; found {
			mp[IVar(key)] = val
		}
	}

	return CreateIcingaService(h.IcingaClient, mp, objectList[0], h.Resource.Name)
}

func (h *ClusterHost) Update() error {
	alertSpec := h.Resource.Spec

	// Get Icinga Host Info
	objectList, err := GetObjectList(h.KubeClient, alertSpec.Check, HostTypeLocalhost, h.Resource.Namespace, h.ObjectType, h.ObjectName, "")
	if err != nil {
		return errors.FromErr(err).Err()
	}

	if err := h.updateIcingaService(objectList); err != nil {
		return errors.FromErr(err).Err()
	}

	if err := UpdateIcingaNotification(h.IcingaClient, h.Resource, objectList); err != nil {
		return errors.FromErr(err).Err()
	}
	return nil
}

func (h *ClusterHost) updateIcingaService(objectList []*KubeObjectInfo) error {
	alertSpec := h.Resource.Spec

	mp := make(map[string]interface{})
	if alertSpec.CheckInterval.Seconds() > 0 {
		mp["check_interval"] = alertSpec.CheckInterval.Seconds()
	}

	commandVars := tapi.ClusterCommands[alertSpec.Check].Vars
	for key, val := range alertSpec.Vars {
		if _, found := commandVars[key]; found {
			mp[IVar(key)] = val
		}
	}

	for _, object := range objectList {
		if err := UpdateIcingaService(h.IcingaClient, mp, object, h.Resource.Name); err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}

func (h *ClusterHost) Delete() error {
	alertSpec := h.Resource.Spec

	objectList, err := GetObjectList(h.KubeClient, alertSpec.Check, HostTypeLocalhost, h.Resource.Namespace, h.ObjectType, h.ObjectName, "")
	if err != nil {
		return errors.FromErr(err).Err()
	}

	if err := DeleteIcingaService(h.IcingaClient, objectList, h.Resource.Name); err != nil {
		return errors.FromErr(err).Err()
	}

	for _, object := range objectList {
		if err := DeleteIcingaHost(h.IcingaClient, object.Name); err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}
