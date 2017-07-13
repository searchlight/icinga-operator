package framework

import (
	"errors"
	"fmt"
	"time"

	"github.com/appscode/searchlight/test/e2e"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func (f *Framework) CreateService(obj *apiv1.Service) error {
	_, err := f.kubeClient.CoreV1().Services(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) DeleteService(meta metav1.ObjectMeta) error {
	return f.kubeClient.CoreV1().Services(meta.Namespace).Delete(meta.Name, deleteInForeground())
}

func (f *Framework) GetServiceIngressHost(meta metav1.ObjectMeta) (string, error) {

	then := time.Now()
	now := time.Now()

	for now.Sub(then) < time.Minute*5 {
		service, err := f.kubeClient.CoreV1().Services(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if err != nil {
			return "", err
		}

		for _, ingress := range service.Status.LoadBalancer.Ingress {
			if ingress.Hostname != "" {
				e2e.PrintSeparately("Ingress Hostname is available")
				return ingress.Hostname, nil
			}
		}

		fmt.Println("Waiting for Ingress Hostname in service")
		time.Sleep(time.Second * 10)
		now = time.Now()
	}

	return "", errors.New("Failed to get Ingress Hostname in service")
}
