/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package framework

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	shell "github.com/codeskyblue/go-sh"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	TEST_HEADLESS_SERVICE = "headless"
)

func (f *Framework) CreateService(obj *core.Service) error {
	_, err := f.kubeClient.CoreV1().Services(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

func (f *Framework) DeleteService(meta metav1.ObjectMeta) error {
	return f.kubeClient.CoreV1().Services(meta.Namespace).Delete(context.TODO(), meta.Name, *deleteInForeground())
}

func (f *Framework) GetServiceEndpoint(meta metav1.ObjectMeta, portName string) (string, error) {
	service, err := f.kubeClient.CoreV1().Services(meta.Namespace).Get(context.TODO(), meta.Name, metav1.GetOptions{})
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
	} else {
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
			PrintSeparately("LoadBalancer is ready")
			return true
		},
		time.Minute*5,
		time.Second*5,
	)
}
