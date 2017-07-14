package framework

import (
	"github.com/appscode/go/crypto/rand"
	tapi "github.com/appscode/searchlight/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Invocation) PodAlert() *tapi.PodAlert {
	return &tapi.PodAlert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("podalert"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: tapi.PodAlertSpec{},
	}
}

func (f *Framework) CreatePodAlert(obj *tapi.PodAlert) error {
	_, err := f.extClient.PodAlerts(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) DeletePodAlert(meta metav1.ObjectMeta) error {
	return f.extClient.PodAlerts(meta.Namespace).Delete(meta.Name)
}
