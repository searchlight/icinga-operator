package framework

import (
	"encoding/json"
	"fmt"
	"time"

	incident_api "github.com/appscode/searchlight/apis/incidents/v1alpha1"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
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

	incidentList, err := f.extClient.MonitoringV1alpha1().Incidents(meta.Namespace).List(metav1.ListOptions{
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

	_, err = f.extClient.IncidentsV1alpha1().Acknowledgements(incident.Namespace).Create(&incident_api.Acknowledgement{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: incident.Namespace,
			Name:      incident.Name,
		},
		Request: incident_api.AcknowledgementRequest{
			Comment: "test",
		},
	})
	if err != nil {
		return err
	}

	return nil
}
