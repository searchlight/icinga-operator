package controller

import (
	"github.com/appscode/errors"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/controller/host/extpoints"
)

func (b *IcingaController) Delete(specificObject ...string) error {
	if !b.checkIcingaAvailability() {
		return errors.New("Icinga is down").Err()
	}

	log.Debugln("Starting deleting alert", b.ctx.Resource.ObjectMeta)

	object := ""
	if len(specificObject) > 0 {
		object = specificObject[0]
	}

	alertSpec := b.ctx.Resource.Spec
	command, found := b.ctx.IcingaData[alertSpec.CheckCommand]
	if !found {
		return errors.Newf("check_command [%s] not found", alertSpec.CheckCommand).Err()
	}
	hostType, found := command.HostType[b.ctx.ObjectType]
	if !found {
		return errors.Newf("check_command [%s] is not applicable to %s", alertSpec.CheckCommand, b.ctx.ObjectType).Err()
	}
	p := extpoints.IcingaHostTypes.Lookup(hostType)
	if p == nil {
		return errors.Newf("IcingaHostType %v is unknown", hostType).Err()
	}
	return p.DeleteAlert(b.ctx, object)
}
