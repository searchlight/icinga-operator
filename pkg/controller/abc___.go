package controller

import (
	tapi "github.com/appscode/searchlight/api"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func (c *Controller) EnsureIcingaSS(resource *apiv1.Pod, old, new *tapi.PodAlert) (err error) {
	return nil
}

func (c *Controller) EnsureIcingaSSDeleted(resource *apiv1.Pod, restic *tapi.PodAlert) (err error) {
	return nil
}
