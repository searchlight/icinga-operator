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
	"fmt"

	"go.searchlight.dev/icinga-operator/pkg/icinga"

	"gomodules.xyz/x/crypto/rand"
	core_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Invocation) GetWebHookSecret() *core_v1.Secret {
	return &core_v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("notifier"),
			Namespace: f.namespace,
		},
		StringData: map[string]string{},
	}
}

func (f *Framework) CreateWebHookSecret(obj *core_v1.Secret) error {
	_, err := f.kubeClient.CoreV1().Secrets(obj.Namespace).Create(context.TODO(), obj, metav1.CreateOptions{})
	return err
}

func (f *Invocation) GetIcingaApiPassword(objectMeta metav1.ObjectMeta) (string, error) {
	secret, err := f.kubeClient.CoreV1().Secrets(objectMeta.Namespace).Get(context.TODO(), objectMeta.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	pass, found := secret.Data[icinga.ICINGA_API_PASSWORD]
	if !found {
		return "", fmt.Errorf(`key "%s" is not found in Secret "%s/%s"`, icinga.ICINGA_API_PASSWORD, objectMeta.Namespace, objectMeta.Name)
	}

	return string(pass), nil
}
