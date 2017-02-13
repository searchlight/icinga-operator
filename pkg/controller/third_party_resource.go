package controller

import (
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
)

func CreateThirdPartyResource(kubeClient clientset.Interface) error {
	_, err := kubeClient.Extensions().ThirdPartyResources().Get("alert.appscode.com")
	if !errors.IsNotFound(err) {
		return err
	}

	thirdPartyResource := &extensions.ThirdPartyResource{
		TypeMeta: unversioned.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "ThirdPartyResource",
		},
		ObjectMeta: kapi.ObjectMeta{
			Name: "alert.appscode.com",
		},
		Versions: []extensions.APIVersion{
			{
				Name: "v1beta1",
			},
		},
	}

	if _, err := kubeClient.Extensions().ThirdPartyResources().Create(thirdPartyResource); err != nil {
		return err
	}
	return nil
}
