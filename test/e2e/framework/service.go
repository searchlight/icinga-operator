package framework

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/appscode/searchlight/test/e2e"
	shell "github.com/codeskyblue/go-sh"
	. "github.com/onsi/gomega"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	TEST_HEADLESS_SERVICE = "headless"
)

func (f *Framework) CreateService(obj *apiv1.Service) error {
	_, err := f.kubeClient.CoreV1().Services(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) DeleteService(meta metav1.ObjectMeta) error {
	return f.kubeClient.CoreV1().Services(meta.Namespace).Delete(meta.Name, deleteInForeground())
}

func (f *Framework) GetServiceEndpoint(meta metav1.ObjectMeta, portName string) (string, error) {
	service, err := f.kubeClient.CoreV1().Services(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	var port int32
	for _, p := range service.Spec.Ports {
		if p.Name == portName {
			if strings.ToLower(f.Provider) == "minikube" {
				port = p.NodePort
			} else {
				port = p.Port
			}
		}
	}

	var host interface{} = nil

	if strings.ToLower(f.Provider) == "minikube" {
		if port != 0 {
			ip, err := shell.Command("minikube", "ip").Output()
			if err != nil {
				return "", err
			}
			host = strings.TrimRight(string(ip), "\n")
		}
	}

	if strings.ToLower(f.Provider) == "aws" {
		for _, ingress := range service.Status.LoadBalancer.Ingress {
			if ingress.Hostname != "" {
				host = ingress.Hostname
				break
			} else if ingress.IP != "" {
				host = ingress.IP
			}
		}
	}

	if host != nil {
		return fmt.Sprintf("%v:%v", host, port), nil
	}

	return "", errors.New("API Endpoint not found")
}

func (f *Framework) EventuallyServiceLoadBalancer(meta metav1.ObjectMeta, portName string) GomegaAsyncAssertion {
	return Eventually(
		func() bool {
			endpoint, _ := f.GetServiceEndpoint(meta, portName)
			if endpoint == "" {
				fmt.Println("Waiting for LoadBalancer")
				return false
			}
			e2e.PrintSeparately("LoadBalancer is ready")
			return true
		},
		time.Minute*5,
		time.Second*5,
	)
}
