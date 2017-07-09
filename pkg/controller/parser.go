package controller

import (
	"github.com/appscode/log"
)

func (c *Controller) parseAlertOptions() {
	if c.opt.Resource == nil {
		log.Infoln("Config is nil, nothing to parse")
		return
	}
	log.Infoln("Parsing labels.")
	c.opt.ObjectType, c.opt.ObjectName = GetObjectInfo(c.opt.Resource.ObjectMeta.Labels)
}
