package framework

import (
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/log"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	sutil "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1/util"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Invocation) PodAlert() *api.PodAlert {
	return &api.PodAlert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("podalert"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: api.PodAlertSpec{
			CheckInterval: metav1.Duration{time.Second * 5},
			Vars:          make(map[string]string),
		},
	}
}

func (f *Framework) CreatePodAlert(obj *api.PodAlert) error {
	_, err := f.extClient.PodAlerts(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) GetPodAlert(meta metav1.ObjectMeta) (*api.PodAlert, error) {
	return f.extClient.PodAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) UpdatePodAlert(meta metav1.ObjectMeta, transformer func(api.PodAlert) api.PodAlert) (*api.PodAlert, error) {
	attempt := 0
	for ; attempt < maxAttempts; attempt = attempt + 1 {
		cur, err := f.extClient.PodAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		modified := transformer(*cur)
		updated, err := f.extClient.PodAlerts(cur.Namespace).Update(&modified)
		if err == nil {
			return updated, nil
		}

		log.Errorf("Attempt %d failed to update PodAlert %s@%s due to %s.", attempt, cur.Name, cur.Namespace, err)
		time.Sleep(updateRetryInterval)
	}

	return nil, fmt.Errorf("Failed to update PodAlert %s@%s after %d attempts.", meta.Name, meta.Namespace, attempt)
}

func (f *Framework) TryPatchPodAlert(meta metav1.ObjectMeta, transform func(alert *api.PodAlert) *api.PodAlert) (*api.PodAlert, error) {
	return sutil.TryPatchPodAlert(f.extClient, meta, transform)
}

func (f *Framework) DeletePodAlert(meta metav1.ObjectMeta) error {
	return f.extClient.PodAlerts(meta.Namespace).Delete(meta.Name, &metav1.DeleteOptions{})
}

func (f *Framework) getPodAlertObjects(meta metav1.ObjectMeta, podAlertSpec api.PodAlertSpec) ([]icinga.IcingaHost, error) {
	names := make([]string, 0)

	if podAlertSpec.PodName != "" {
		pod, err := f.kubeClient.CoreV1().Pods(meta.Namespace).Get(podAlertSpec.PodName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		names = append(names, pod.Name)
	} else {
		sel, err := metav1.LabelSelectorAsSelector(&podAlertSpec.Selector)
		if err != nil {
			return nil, err
		}
		podList, err := f.kubeClient.CoreV1().Pods(meta.Namespace).List(
			metav1.ListOptions{
				LabelSelector: sel.String(),
			},
		)
		if err != nil {
			return nil, err
		}

		for _, pod := range podList.Items {
			names = append(names, pod.Name)
		}

	}

	objectList := make([]icinga.IcingaHost, 0)
	for _, name := range names {
		objectList = append(objectList,
			icinga.IcingaHost{
				Type:           icinga.TypePod,
				AlertNamespace: meta.Namespace,
				ObjectName:     name,
			})
	}

	return objectList, nil
}

func (f *Framework) EventuallyPodAlertIcingaService(meta metav1.ObjectMeta, podAlertSpec api.PodAlertSpec) GomegaAsyncAssertion {
	objectList, err := f.getPodAlertObjects(meta, podAlertSpec)
	Expect(err).NotTo(HaveOccurred())

	in := icinga.NewPodHost(nil, nil, f.icingaClient).
		IcingaServiceSearchQuery(meta.Name, objectList...)

	return Eventually(
		func() matcher.IcingaServiceState {
			var respService icinga.ResponseObject
			status, err := f.icingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService)
			if status == 0 {
				return matcher.IcingaServiceState{Unknown: 1.0}
			}
			Expect(err).NotTo(HaveOccurred())

			var icingaServiceState matcher.IcingaServiceState
			for _, service := range respService.Results {
				if service.Attrs.LastState == 0.0 {
					icingaServiceState.Ok++
				}
				if service.Attrs.LastState == 1.0 {
					icingaServiceState.Warning++
				}
				if service.Attrs.LastState == 2.0 {
					icingaServiceState.Critical++
				}
				if service.Attrs.LastState == 3.0 {
					icingaServiceState.Unknown++
				}
			}
			return icingaServiceState
		},
		time.Minute*5,
		time.Second*5,
	)
}
