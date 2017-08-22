package framework

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/kutil"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/test/e2e/matcher"
	"github.com/golang/glog"
	"github.com/mattbaird/jsonpatch"
	. "github.com/onsi/gomega"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
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
		Spec: tapi.PodAlertSpec{
			CheckInterval: metav1.Duration{time.Second * 5},
			Vars:          make(map[string]interface{}),
		},
	}
}

func (f *Framework) CreatePodAlert(obj *tapi.PodAlert) error {
	_, err := f.extClient.PodAlerts(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) GetPodAlert(meta metav1.ObjectMeta) (*tapi.PodAlert, error) {
	return f.extClient.PodAlerts(meta.Namespace).Get(meta.Name)
}

func (f *Framework) patchPodAlert(cur *tapi.PodAlert, transform func(*tapi.PodAlert) *tapi.PodAlert) (*tapi.PodAlert, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, err
	}

	modJson, err := json.Marshal(transform(cur))
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.CreatePatch(curJson, modJson)
	if err != nil {
		return nil, err
	}
	if len(patch) == 0 {
		return cur, nil
	}
	pb, err := json.MarshalIndent(patch, "", "  ")
	if err != nil {
		return nil, err
	}
	glog.V(5).Infof("Patching PodAlert %s@%s with %s.", cur.Name, cur.Namespace, string(pb))
	return f.extClient.PodAlerts(cur.Namespace).Patch(cur.Name, types.JSONPatchType, pb)
}

func (f *Framework) TryPatchPodAlert(meta metav1.ObjectMeta, transform func(*tapi.PodAlert) *tapi.PodAlert) (*tapi.PodAlert, error) {
	attempt := 0
	for ; attempt < kutil.MaxAttempts; attempt = attempt + 1 {
		cur, err := f.extClient.PodAlerts(meta.Namespace).Get(meta.Name)
		if kerr.IsNotFound(err) {
			return cur, err
		} else if err == nil {
			return f.patchPodAlert(cur, transform)
		}
		glog.Errorf("Attempt %d failed to patch PodAlert %s@%s due to %s.", attempt, cur.Name, cur.Namespace, err)
		time.Sleep(kutil.RetryInterval)
	}
	return nil, fmt.Errorf("Failed to patch PodAlert %s@%s after %d attempts.", meta.Name, meta.Namespace, attempt)
}

func (f *Framework) DeletePodAlert(meta metav1.ObjectMeta) error {
	return f.extClient.PodAlerts(meta.Namespace).Delete(meta.Name)
}

func (f *Framework) getPodAlertObjects(meta metav1.ObjectMeta, podAlertSpec tapi.PodAlertSpec) ([]icinga.IcingaHost, error) {
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

func (f *Framework) EventuallyPodAlertIcingaService(meta metav1.ObjectMeta, podAlertSpec tapi.PodAlertSpec) GomegaAsyncAssertion {
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
