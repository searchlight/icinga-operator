package host

import (
	"fmt"
	"regexp"

	"github.com/appscode/errors"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/controller/types"
)

type NodeHost struct {
	*types.Option
}

//-----------------------------------------------------
// set Alert in Icinga LocalHost
func (b *NodeHost) Create(specificObject string) error {
	alertSpec := b.Resource.Spec

	if alertSpec.Check == "" {
		return errors.New("Invalid request").Err()
	}

	// Get Icinga Host Info
	objectList, err := GetObjectList(b.KubeClient, alertSpec.Check, HostTypeNode, b.Resource.Namespace, b.ObjectType, b.ObjectName, specificObject)
	if err != nil {
		return errors.FromErr(err).Err()
	}

	var has bool
	if has, err = CheckIcingaService(b.IcingaClient, b.Resource.Name, objectList); err != nil {
		return errors.FromErr(err).Err()
	}
	if has {
		return nil
	}

	// Create Icinga Host
	if err := CreateIcingaHost(b.IcingaClient, objectList, b.Resource.Namespace); err != nil {
		return errors.FromErr(err).Err()
	}

	if err := b.createIcingaService(objectList); err != nil {
		return errors.FromErr(err).Err()
	}

	if err := CreateIcingaNotification(b.IcingaClient, b.Resource, objectList); err != nil {
		return errors.FromErr(err).Err()
	}

	return nil
}

func setParameterizedVariables(alertSpec tapi.PodAlertSpec, objectName string, commandVars map[string]data.CommandVar, mp map[string]interface{}) (map[string]interface{}, error) {
	for key, val := range alertSpec.Vars {
		if v, found := commandVars[key]; found {
			if !v.Parameterized {
				continue
			}

			reg, err := regexp.Compile("nodename[ ]*=[ ]*'[?]'")
			if err != nil {
				return nil, errors.FromErr(err).Err()
			}
			mp[IVar(key)] = reg.ReplaceAllString(val.(string), fmt.Sprintf("nodename='%s'", objectName))
		} else {
			return nil, errors.Newf("variable %v not found", key).Err()
		}
	}
	return mp, nil
}

func (b *NodeHost) createIcingaService(objectList []*KubeObjectInfo) error {
	alertSpec := b.Resource.Spec

	mp := make(map[string]interface{})
	mp["check_command"] = alertSpec.Check
	if alertSpec.CheckInterval.Seconds() > 0 {
		mp["check_interval"] = alertSpec.CheckInterval.Seconds()
	}

	commandVars := tapi.NodeCommands[alertSpec.Check].Vars
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
		if mp, err = setParameterizedVariables(alertSpec, object.Name, commandVars, mp); err != nil {
			return errors.FromErr(err).Err()
		}

		if err := CreateIcingaService(b.IcingaClient, mp, object, b.Resource.Name); err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}

func (b *NodeHost) Update() error {
	alertSpec := b.Resource.Spec

	// Get Icinga Host Info
	objectList, err := GetObjectList(b.KubeClient, alertSpec.Check, HostTypeNode, b.Resource.Namespace, b.ObjectType, b.ObjectName, "")
	if err != nil {
		return errors.FromErr(err).Err()
	}

	if err := b.updateIcingaService(objectList); err != nil {
		return errors.FromErr(err).Err()
	}

	if err := UpdateIcingaNotification(b.IcingaClient, b.Resource, objectList); err != nil {
		return errors.FromErr(err).Err()
	}
	return nil
}

func (b *NodeHost) updateIcingaService(objectList []*KubeObjectInfo) error {
	alertSpec := b.Resource.Spec

	mp := make(map[string]interface{})
	if alertSpec.CheckInterval.Seconds() > 0 {
		mp["check_interval"] = alertSpec.CheckInterval.Seconds()
	}

	commandVars := tapi.NodeCommands[alertSpec.Check].Vars
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
		if mp, err = setParameterizedVariables(alertSpec, object.Name, commandVars, mp); err != nil {
			return errors.FromErr(err).Err()
		}

		if err := UpdateIcingaService(b.IcingaClient, mp, object, b.Resource.Name); err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}

func (b *NodeHost) Delete(specificObject string) error {
	alertSpec := b.Resource.Spec

	var objectList []*KubeObjectInfo
	if specificObject != "" {
		objectList = append(objectList, &KubeObjectInfo{Name: specificObject + "@" + b.Resource.Namespace})
	} else {
		// Get Icinga Host Info
		var err error
		objectList, err = GetObjectList(b.KubeClient, alertSpec.Check, HostTypeNode,
			b.Resource.Namespace, b.ObjectType, b.ObjectName, specificObject)
		if err != nil {
			return errors.FromErr(err).Err()
		}

	}

	if err := DeleteIcingaService(b.IcingaClient, objectList, b.Resource.Name); err != nil {
		return errors.FromErr(err).Err()
	}

	for _, object := range objectList {
		if err := DeleteIcingaHost(b.IcingaClient, object.Name); err != nil {
			return errors.FromErr(err).Err()
		}
	}
	return nil
}
