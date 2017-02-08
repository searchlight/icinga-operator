package mini

import (
	"fmt"
	"sync"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/client/k8s"
	kapi "k8s.io/kubernetes/pkg/api"
)

type icingaSecretInfo struct {
	isSecretSet bool
	name        string
	namespace   string
	// To create secret for Icinga2 only once
	once sync.Once
}

var icingaSecret = icingaSecretInfo{}

func CreateIcingaSecret(kubeClient *k8s.KubeClient, namespace string, secretMap map[string]string) (secretName string, err error) {
	if icingaSecret.isSecretSet {
		secretName = icingaSecret.name + "." + icingaSecret.namespace
		return
	}
	icingaSecret.once.Do(
		func() {
			var secretString string
			for key, val := range secretMap {
				secretString = secretString + fmt.Sprintf("%s=%s\n", key, val)
			}

			secret := &kapi.Secret{
				ObjectMeta: kapi.ObjectMeta{
					Name:      rand.WithUniqSuffix("fake-secret"),
					Namespace: namespace,
				},
				Data: map[string][]byte{
					icinga.ENV: []byte(secretString),
				},
				Type: kapi.SecretTypeOpaque,
			}
			// Create Fake Secret
			if _, err = kubeClient.Client.Core().Secrets(secret.Namespace).Create(secret); err != nil {
				return
			}
			icingaSecret.name = secret.Name
			icingaSecret.namespace = secret.Namespace
			secretName = icingaSecret.name + "." + icingaSecret.namespace
		},
	)
	return
}
