package mini

import (
	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	kapi "k8s.io/kubernetes/pkg/api"
)

func CreateService(watcher *app.Watcher, namespace string, selector map[string]string) (*kapi.Service, error) {
	service := &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Namespace: namespace,
		},
		Spec: kapi.ServiceSpec{
			Selector: selector,
		},
	}
	if err := testing.CreateKubernetesObject(watcher.Client, service); err != nil {
		return nil, err
	}
	return service, nil
}

func DeleteService(watcher *app.Watcher, service *kapi.Service) error {
	// Delete Service
	if err := watcher.Client.Core().Services(service.Namespace).Delete(service.Name, nil); err != nil {
		return err
	}
	return nil
}
