package framework

import (
	"fmt"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/icinga"
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
	_, err := f.kubeClient.CoreV1().Secrets(obj.Namespace).Create(obj)
	return err
}

func (f *Invocation) GetIcingaApiPassword(objectMeta metav1.ObjectMeta) (string, error) {
	secret, err := f.kubeClient.CoreV1().Secrets(objectMeta.Namespace).Get(objectMeta.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	pass, found := secret.Data[icinga.ICINGA_API_PASSWORD]
	if !found {
		return "", fmt.Errorf(`key "%s" is not found in Secret "%s/%s"`, icinga.ICINGA_API_PASSWORD, objectMeta.Namespace, objectMeta.Name)
	}

	return string(pass), nil
}
