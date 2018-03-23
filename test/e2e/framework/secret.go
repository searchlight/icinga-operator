package framework

import (
	"github.com/appscode/go/crypto/rand"
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
