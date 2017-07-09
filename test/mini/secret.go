package mini

import (
	"fmt"
	"sync"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

type icingaSecretInfo struct {
	isSecretSet bool
	name        string
	namespace   string
	// To create secret for Icinga2 only once
	once sync.Once
}

var icingaSecret = icingaSecretInfo{}

func CreateIcingaSecret(kubeClient *util.KubeClient, namespace string, secretMap map[string]string) (secretName string, err error) {
	if icingaSecret.isSecretSet {
		secretName = icingaSecret.name
		return
	}
	icingaSecret.once.Do(
		func() {
			var secretString string
			for key, val := range secretMap {
				secretString = secretString + fmt.Sprintf("%s=%s\n", key, val)
			}

			secret := &apiv1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      rand.WithUniqSuffix("fake-secret"),
					Namespace: namespace,
				},
				Data: map[string][]byte{
					icinga.ENV: []byte(secretString),
				},
				Type: apiv1.SecretTypeOpaque,
			}
			// Create Fake Secret
			if _, err = kubeClient.Client.CoreV1().Secrets(secret.Namespace).Create(secret); err != nil {
				return
			}
			icingaSecret.name = secret.Name
			icingaSecret.namespace = secret.Namespace
			secretName = icingaSecret.name
		},
	)
	return
}
