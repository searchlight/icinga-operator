package framework

import (
	"encoding/json"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	mp := make(map[string]interface{})
	mp["type"] = "Service"
	mp["filter"] = fmt.Sprintf(`service.name == "%s" && host.name == "%s"`, meta.Name, hostname)
	mp["author"] = "e2e"
	mp["comment"] = "test"
	mp["notify"] = true
	ack, err := json.Marshal(mp)
	if err != nil {
		return err
	}
	return f.icingaClient.Actions("acknowledge-problem").Update([]string{}, string(ack)).Do().Err
}
