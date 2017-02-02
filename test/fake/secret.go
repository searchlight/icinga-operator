package fake

import (
	"fmt"
	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/client/k8s"
	kapi "k8s.io/kubernetes/pkg/api"
	"os"
	"sync"
)

type fakeIcingaSecret struct {
	name      string
	namespace string
	// To create secret for Icinga2 only once
	once sync.Once
}

var fakeSecret = fakeIcingaSecret{}

func CreateFakeIcingaSecret(kubeClient *k8s.KubeClient, namespace string, secretMap map[string]string) string {
	fakeSecret.once.Do(
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
			if _, err := kubeClient.Client.Core().Secrets(secret.Namespace).Create(secret); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fakeSecret.name = secret.Name
			fakeSecret.namespace = secret.Namespace
		},
	)
	// Returns Secret name with namespace extension
	return fakeSecret.name + "." + fakeSecret.namespace
}
