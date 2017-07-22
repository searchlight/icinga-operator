package framework

import (
	"errors"
	"time"

	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (f *Framework) EventuallyTPR() GomegaAsyncAssertion {
	label := map[string]string{
		"app": "searchlight",
	}
	return Eventually(
		func() error {
			tprList, err := f.kubeClient.ExtensionsV1beta1().ThirdPartyResources().List(
				metav1.ListOptions{
					LabelSelector: labels.SelectorFromSet(label).String(),
				},
			)
			if err != nil {
				return err
			}

			if len(tprList.Items) != 3 {
				return errors.New("All ThirdPartyResources are not ready")
			}
			return nil
		},
		time.Minute*2,
		time.Second*2,
	)
}
