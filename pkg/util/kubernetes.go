package util

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/appscode/go-notify/unified"
	tapi "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	tcs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func OperatorNamespace() string {
	if ns := os.Getenv("OPERATOR_NAMESPACE"); ns != "" {
		return ns
	}
	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns
		}
	}
	return apiv1.NamespaceDefault
}

func CheckNotifiers(kubeClient clientset.Interface, alert tapi.Alert) error {
	if alert.GetNotifierSecretName() == "" && len(alert.GetReceivers()) == 0 {
		return nil
	}
	secret, err := kubeClient.CoreV1().Secrets(alert.GetNamespace()).Get(alert.GetNotifierSecretName(), metav1.GetOptions{})
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

func FindPodAlert(searchlightClient tcs.MonitoringV1alpha1Interface, obj metav1.ObjectMeta) ([]*tapi.PodAlert, error) {
	alerts, err := searchlightClient.PodAlerts(obj.Namespace).List(metav1.ListOptions{LabelSelector: labels.Everything().String()})
	if kerr.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]*tapi.PodAlert, 0)
	for i, alert := range alerts.Items {
		if ok, _ := alert.IsValid(); !ok {
			continue
		}
		if alert.Spec.PodName != "" && alert.Spec.PodName != obj.Name {
			continue
		}
		if selector, err := metav1.LabelSelectorAsSelector(&alert.Spec.Selector); err == nil {
			if selector.Matches(labels.Set(obj.Labels)) {
				result = append(result, &alerts.Items[i])
			}
		}
	}
	return result, nil
}

func FindNodeAlert(searchlightClient tcs.MonitoringV1alpha1Interface, obj metav1.ObjectMeta) ([]*tapi.NodeAlert, error) {
	alerts, err := searchlightClient.NodeAlerts(obj.Namespace).List(metav1.ListOptions{LabelSelector: labels.Everything().String()})
	if kerr.IsNotFound(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]*tapi.NodeAlert, 0)
	for i, alert := range alerts.Items {
		if ok, _ := alert.IsValid(); !ok {
			continue
		}
		if alert.Spec.NodeName != "" && alert.Spec.NodeName != obj.Name {
			continue
		}
		selector := labels.SelectorFromSet(alert.Spec.Selector)
		if selector.Matches(labels.Set(obj.Labels)) {
			result = append(result, &alerts.Items[i])
		}
	}
	return result, nil
}
