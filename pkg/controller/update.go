package controller

import (
	"github.com/appscode/errors"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/icinga/host/extpoints"
)

func (c *Controller) Update() error {
	if !c.checkIcingaAvailability() {
		return errors.New("Icinga is down").Err()
	}

	log.Debugln("Starting updating alert", c.opt.Resource.ObjectMeta)

	alertSpec := c.opt.Resource.Spec
	command, found := c.opt.IcingaData[alertSpec.Check]
	if !found {
		return errors.Newf("check_command [%s] not found", alertSpec.Check).Err()
	}
	hostType, found := command.HostType[c.opt.ObjectType]
	if !found {
		return errors.Newf("check_command [%s] is not applicable to %s", alertSpec.Check, c.opt.ObjectType).Err()
	}
	p := extpoints.IcingaHostTypes.Lookup(hostType)
	if p == nil {
		return errors.Newf("IcingaHostType %v is unknown", hostType).Err()
	}
	return p.UpdateAlert(c.opt)
}
