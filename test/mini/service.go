package mini

import (
	"fmt"
	"os"

	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	kapi "k8s.io/kubernetes/pkg/api"
)

func CreateService(watcher *app.Watcher, namespace string, selector map[string]string) *kapi.Service {
	service := &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Namespace: namespace,
		},
		Spec: kapi.ServiceSpec{
			Selector: selector,
		},
	}
	if err := testing.CreateKubernetesObject(watcher.Client, service); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return service
}

func DeleteService(watcher *app.Watcher, service *kapi.Service) {
	// Delete Service
	if err := watcher.Client.Core().Services(service.Namespace).Delete(service.Name, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
