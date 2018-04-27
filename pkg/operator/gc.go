package operator

import (
	"time"

	"github.com/appscode/go/log"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (op *Operator) gcIncidents() {
	if op.IncidentTTL <= 0 {
		log.Warningln("skipping garbage collection of incidents")
		return
	}

	ticker := time.NewTicker(op.IncidentTTL)
	go func() {
		for t := range ticker.C {
			log.Infoln("Incident GC run at", t)

			objects, err := op.extClient.MonitoringV1alpha1().Incidents(core.NamespaceAll).List(metav1.ListOptions{})
			if err != nil {
				log.Errorln(err)
				continue
			}

			for _, item := range objects.Items {
				if item.Status.LastNotificationType == api.NotificationRecovery &&
					t.Sub(item.CreationTimestamp.Time) > op.IncidentTTL {
					op.extClient.MonitoringV1alpha1().Incidents(item.Namespace).Delete(item.Name, nil)
				}
			}
		}
	}()
}
