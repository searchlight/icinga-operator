package localhost

import (
	"github.com/appscode/errors"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/pkg/controller/host/extpoints"
	"github.com/appscode/searchlight/pkg/controller/types"
)

func init() {
	extpoints.IcingaHostTypes.Register(new(icingaHost), host.HostTypeLocalhost)
}

type icingaHost struct {
}

type biblio struct {
	*types.Option
}

func (p *icingaHost) CreateAlert(ctx *types.Option, specificObject string) error {
	if specificObject != "" {
		return nil
	}
	return (&biblio{ctx}).create()
}

func (p *icingaHost) UpdateAlert(ctx *types.Option) error {
	return (&biblio{ctx}).update()
}

func (p *icingaHost) DeleteAlert(ctx *types.Option, specificObject string) error {
	return (&biblio{ctx}).delete()
}

//-----------------------------------------------------

// set Alert in Icinga LocalHost
func (b *biblio) create() error {
	alertSpec := b.Resource.Spec
	if alertSpec.Check == "" {
		return errors.New("Invalid request").Err()
	}

	// Get Icinga Host Info
	objectList, err := host.GetObjectList(b.KubeClient, alertSpec.Check, host.HostTypeLocalhost, b.Resource.Namespace, b.ObjectType, b.ObjectName, "")
	if err != nil {
		return errors.New().WithCause(err).Err()
	}

	var has bool
	if has, err = host.CheckIcingaService(b.IcingaClient, b.Resource.Name, objectList); err != nil {
		return errors.New().WithCause(err).Err()
	}
	if has {
		return nil
	}

	// Create Icinga Host
	if err := host.CreateIcingaHost(b.IcingaClient, objectList, b.Resource.Namespace); err != nil {
		return errors.New().WithCause(err).Err()
	}

	if err := b.createIcingaService(objectList); err != nil {
		return errors.New().WithCause(err).Err()
	}

	if err := host.CreateIcingaNotification(b.IcingaClient, b.Resource, objectList); err != nil {
		return errors.New().WithCause(err).Err()
	}

	return nil
}

func (b *biblio) createIcingaService(objectList []*host.KubeObjectInfo) error {
	alertSpec := b.Resource.Spec

	mp := make(map[string]interface{})
	mp["check_command"] = alertSpec.Check
	if alertSpec.CheckInterval.Seconds() > 0 {
		mp["check_interval"] = alertSpec.CheckInterval.Seconds()
	}

	commandVars := b.IcingaData[alertSpec.Check].VarInfo

	for key, val := range alertSpec.Vars {
		if _, found := commandVars[key]; found {
			mp[host.IVar(key)] = val
		}
	}

	return host.CreateIcingaService(b.IcingaClient, mp, objectList[0], b.Resource.Name)
}

func (b *biblio) update() error {
	alertSpec := b.Resource.Spec

	// Get Icinga Host Info
	objectList, err := host.GetObjectList(b.KubeClient, alertSpec.Check, host.HostTypeLocalhost, b.Resource.Namespace, b.ObjectType, b.ObjectName, "")
	if err != nil {
		return errors.New().WithCause(err).Err()
	}

	if err := b.updateIcingaService(objectList); err != nil {
		return errors.New().WithCause(err).Err()
	}

	if err := host.UpdateIcingaNotification(b.IcingaClient, b.Resource, objectList); err != nil {
		return errors.New().WithCause(err).Err()
	}
	return nil
}

func (b *biblio) updateIcingaService(objectList []*host.KubeObjectInfo) error {
	alertSpec := b.Resource.Spec

	mp := make(map[string]interface{})
	if alertSpec.CheckInterval.Seconds() > 0 {
		mp["check_interval"] = alertSpec.CheckInterval.Seconds()
	}

	commandVars := b.IcingaData[alertSpec.Check].VarInfo
	for key, val := range alertSpec.Vars {
		if _, found := commandVars[key]; found {
			mp[host.IVar(key)] = val
		}
	}

	for _, object := range objectList {
		if err := host.UpdateIcingaService(b.IcingaClient, mp, object, b.Resource.Name); err != nil {
			return errors.New().WithCause(err).Err()
		}
	}
	return nil
}

func (b *biblio) delete() error {
	alertSpec := b.Resource.Spec

	objectList, err := host.GetObjectList(b.KubeClient, alertSpec.Check, host.HostTypeLocalhost, b.Resource.Namespace, b.ObjectType, b.ObjectName, "")
	if err != nil {
		return errors.New().WithCause(err).Err()
	}

	if err := host.DeleteIcingaService(b.IcingaClient, objectList, b.Resource.Name); err != nil {
		return errors.New().WithCause(err).Err()
	}

	for _, object := range objectList {
		if err := host.DeleteIcingaHost(b.IcingaClient, object.Name); err != nil {
			return errors.New().WithCause(err).Err()
		}
	}
	return nil
}
