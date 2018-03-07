package util

import (
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/kutil/meta"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	mon_listers "github.com/appscode/searchlight/client/listers/monitoring/v1alpha1"
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

func FindPodAlert(lister mon_listers.PodAlertLister, obj metav1.ObjectMeta) ([]*api.PodAlert, error) {
	alerts, err := lister.PodAlerts(obj.Namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	result := make([]*api.PodAlert, 0)
	for i := range alerts {
		alert := alerts[i]
		if ok, _ := alert.IsValid(); !ok {
			continue
		}

		if alert.Spec.PodName != nil {
			if *alert.Spec.PodName == obj.Name {
				result = append(result, alert)
			}
		} else if alert.Spec.Selector != nil {
			if selector, err := metav1.LabelSelectorAsSelector(alert.Spec.Selector); err == nil {
				if selector.Matches(labels.Set(obj.Labels)) {
					result = append(result, alert)
				}
			}
		}
	}
	return result, nil
}

func FindNodeAlert(lister mon_listers.NodeAlertLister, obj metav1.ObjectMeta) ([]*api.NodeAlert, error) {
	alerts, err := lister.NodeAlerts(obj.Namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	result := make([]*api.NodeAlert, 0)
	for i := range alerts {
		alert := alerts[i]
		if ok, _ := alert.IsValid(); !ok {
			continue
		}

		if alert.Spec.NodeName != nil {
			if *alert.Spec.NodeName == obj.Name {
				result = append(result, alert)
			}
		} else {
			selector := labels.SelectorFromSet(alert.Spec.Selector)
			if selector.Matches(labels.Set(obj.Labels)) {
				result = append(result, alert)
			}
		}
	}
	return result, nil
}

func NodeAlertEqual(old, new *api.NodeAlert) bool {
	var oldSpec, newSpec *api.NodeAlertSpec
	if old != nil {
		oldSpec = &old.Spec
	}
	if new != nil {
		newSpec = &new.Spec
	}
	return meta.Equal(oldSpec, newSpec)
}

func PodAlertEqual(old, new *api.PodAlert) bool {
	var oldSpec, newSpec *api.PodAlertSpec
	if old != nil {
		oldSpec = &old.Spec
	}
	if new != nil {
		newSpec = &new.Spec
	}
	return meta.Equal(oldSpec, newSpec)
}
