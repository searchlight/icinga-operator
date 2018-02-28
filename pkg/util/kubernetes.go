package util

import (
	"github.com/appscode/go-notify/unified"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	slite_listers "github.com/appscode/searchlight/client/listers/monitoring/v1alpha1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func CheckNotifiers(client kubernetes.Interface, alert api.Alert) error {
	if alert.GetNotifierSecretName() == "" && len(alert.GetReceivers()) == 0 {
		return nil
	}
	secret, err := client.CoreV1().Secrets(alert.GetNamespace()).Get(alert.GetNotifierSecretName(), metav1.GetOptions{})
	if err != nil {
		return err
	}
	for _, r := range alert.GetReceivers() {
		_, err = unified.LoadVia(r.Notifier, func(key string) (value string, found bool) {
			var bytes []byte
			bytes, found = secret.Data[key]
			value = string(bytes)
			return
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func FindPodAlert(pa slite_listers.PodAlertLister, obj metav1.ObjectMeta) ([]*api.PodAlert, error) {
	alerts, err := pa.PodAlerts(obj.Namespace).List(labels.Everything())
	if kerr.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]*api.PodAlert, 0)
	for i, alert := range alerts {
		if ok, _ := alert.IsValid(); !ok {
			continue
		}
		if alert.Spec.PodName != "" && alert.Spec.PodName != obj.Name {
			continue
		}
		if selector, err := metav1.LabelSelectorAsSelector(&alert.Spec.Selector); err == nil {
			if selector.Matches(labels.Set(obj.Labels)) {
				result = append(result, alerts[i])
			}
		}
	}
	return result, nil
}

func FindNodeAlert(na slite_listers.NodeAlertLister, obj metav1.ObjectMeta) ([]*api.NodeAlert, error) {
	alerts, err := na.NodeAlerts(obj.Namespace).List(labels.Everything())
	if kerr.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]*api.NodeAlert, 0)
	for i, alert := range alerts {
		if ok, _ := alert.IsValid(); !ok {
			continue
		}
		if alert.Spec.NodeName != "" && alert.Spec.NodeName != obj.Name {
			continue
		}
		selector := labels.SelectorFromSet(alert.Spec.Selector)
		if selector.Matches(labels.Set(obj.Labels)) {
			result = append(result, alerts[i])
		}
	}
	return result, nil
}
