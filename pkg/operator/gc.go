/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

import (
	"context"
	"time"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func (op *Operator) gcIncidents() {
	if op.IncidentTTL <= 0 {
		klog.Warningln("skipping garbage collection of incidents")
		return
	}

	ticker := time.NewTicker(op.IncidentTTL)
	go func() {
		for t := range ticker.C {
			klog.Infoln("Incident GC run at", t)

			objects, err := op.extClient.MonitoringV1alpha1().Incidents(core.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				klog.Errorln(err)
				continue
			}

			for _, item := range objects.Items {
				if item.Status.LastNotificationType == api.NotificationRecovery &&
					t.Sub(item.CreationTimestamp.Time) > op.IncidentTTL {
					op.extClient.MonitoringV1alpha1().Incidents(item.Namespace).Delete(context.TODO(), item.Name, metav1.DeleteOptions{})
				}
			}
		}
	}()
}
