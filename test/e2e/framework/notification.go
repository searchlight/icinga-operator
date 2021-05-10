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

package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	incident_api "go.searchlight.dev/icinga-operator/apis/incidents/v1alpha1"
	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (f *Framework) ForceCheckClusterAlert(meta metav1.ObjectMeta, hostname string, times int) error {
	mp := make(map[string]interface{})
	mp["type"] = "Service"
	mp["filter"] = fmt.Sprintf(`service.name == "%s" && host.name == "%s"`, meta.Name, hostname)
	mp["force_check"] = true
	checkNow, err := json.Marshal(mp)
	if err != nil {
		return err
	}

	for i := 0; i < times; i++ {
		f.icingaClient.Actions("reschedule-check").Update([]string{}, string(checkNow)).Do()
	}
	return nil
}

func (f *Framework) SendClusterAlertCustomNotification(meta metav1.ObjectMeta, hostname string) error {
	mp := make(map[string]interface{})
	mp["type"] = "Service"
	mp["filter"] = fmt.Sprintf(`service.name == "%s" && host.name == "%s"`, meta.Name, hostname)
	mp["author"] = "e2e"
	mp["comment"] = "test"
	custom, err := json.Marshal(mp)
	if err != nil {
		return err
	}
	return f.icingaClient.Actions("send-custom-notification").Update([]string{}, string(custom)).Do().Err
}

func (f *Framework) AcknowledgeClusterAlertNotification(meta metav1.ObjectMeta, hostname string) error {

	labelMap := map[string]string{
		api.LabelKeyAlert:            meta.Name,
		api.LabelKeyObjectName:       hostname,
		api.LabelKeyProblemRecovered: "false",
	}

	incidentList, err := f.extClient.MonitoringV1alpha1().Incidents(meta.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(labelMap).String(),
	})
	if err != nil {
		return err
	}

	var lastCreationTimestamp time.Time
	var incident *api.Incident
	for _, item := range incidentList.Items {
		if item.CreationTimestamp.After(lastCreationTimestamp) {
			lastCreationTimestamp = item.CreationTimestamp.Time
			incident = &item
		}
	}

	_, err = f.extClient.IncidentsV1alpha1().Acknowledgements(incident.Namespace).Create(context.TODO(), &incident_api.Acknowledgement{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: incident.Namespace,
			Name:      incident.Name,
		},
		Request: incident_api.AcknowledgementRequest{
			Comment: "test",
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}
