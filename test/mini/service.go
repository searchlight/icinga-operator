package mini

import (
	"github.com/appscode/searchlight/pkg/testing"
	"github.com/appscode/searchlight/pkg/watcher"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func CreateService(w *watcher.Watcher, namespace string, selector map[string]string) (*apiv1.Service, error) {
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
		},
		Spec: apiv1.ServiceSpec{
			Selector: selector,
		},
	}
	if err := testing.CreateKubernetesObject(w.KubeClient, service); err != nil {
		return nil, err
	}
	return service, nil
}

func DeleteService(w *watcher.Watcher, service *apiv1.Service) error {
	// Delete Service
	if err := w.KubeClient.CoreV1().Services(service.Namespace).Delete(service.Name, nil); err != nil {
		return err
	}
	return nil
}
