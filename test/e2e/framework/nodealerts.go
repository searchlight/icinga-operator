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
	"k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
)

func (f *Invocation) NodeAlert() *tapi.NodeAlert {
	return &tapi.NodeAlert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("nodealert"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: tapi.NodeAlertSpec{
			CheckInterval: metav1.Duration{time.Second * 5},
			Vars:          make(map[string]interface{}),
		},
	}
}

func (f *Framework) CreateNodeAlert(obj *tapi.NodeAlert) error {
	_, err := f.extClient.NodeAlerts(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) GetNodeAlert(meta metav1.ObjectMeta) (*tapi.NodeAlert, error) {
	return f.extClient.NodeAlerts(meta.Namespace).Get(meta.Name)
}

func (f *Framework) patchNodeAlert(cur *tapi.NodeAlert, transform func(*tapi.NodeAlert) *tapi.NodeAlert) (*tapi.NodeAlert, error) {
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
	glog.V(5).Infof("Patching NodeAlert %s@%s with %s.", cur.Name, cur.Namespace, string(pb))
	return f.extClient.NodeAlerts(cur.Namespace).Patch(cur.Name, types.JSONPatchType, pb)
}

func (f *Framework) TryPatchNodeAlert(meta metav1.ObjectMeta, transform func(*tapi.NodeAlert) *tapi.NodeAlert) (*tapi.NodeAlert, error) {
	attempt := 0
	for ; attempt < kutil.MaxAttempts; attempt = attempt + 1 {
		cur, err := f.extClient.NodeAlerts(meta.Namespace).Get(meta.Name)
		if kerr.IsNotFound(err) {
			return cur, err
		} else if err == nil {
			return f.patchNodeAlert(cur, transform)
		}
		glog.Errorf("Attempt %d failed to patch NodeAlert %s@%s due to %s.", attempt, cur.Name, cur.Namespace, err)
		time.Sleep(kutil.RetryInterval)
	}
	return nil, fmt.Errorf("Failed to patch NodeAlert %s@%s after %d attempts.", meta.Name, meta.Namespace, attempt)
}

func (f *Framework) DeleteNodeAlert(meta metav1.ObjectMeta) error {
	return f.extClient.NodeAlerts(meta.Namespace).Delete(meta.Name)
}

func (f *Framework) getNodeAlertObjects(meta metav1.ObjectMeta, nodeAlertSpec tapi.NodeAlertSpec) ([]icinga.IcingaHost, error) {
	names := make([]string, 0)
	sel := labels.SelectorFromSet(nodeAlertSpec.Selector)
	nodeList, err := f.kubeClient.CoreV1().Nodes().List(
		metav1.ListOptions{
			LabelSelector: sel.String(),
		},
	)
	if err != nil {
		return nil, err
	}

	for _, node := range nodeList.Items {
		names = append(names, node.Name)
	}

	objectList := make([]icinga.IcingaHost, 0)
	for _, name := range names {
		objectList = append(objectList,
			icinga.IcingaHost{
				Type:           icinga.TypeNode,
				AlertNamespace: meta.Namespace,
				ObjectName:     name,
			})
	}

	return objectList, nil
}

func (f *Framework) EventuallyNodeAlertIcingaService(meta metav1.ObjectMeta, nodeAlertSpec tapi.NodeAlertSpec) GomegaAsyncAssertion {
	objectList, err := f.getNodeAlertObjects(meta, nodeAlertSpec)
	Expect(err).NotTo(HaveOccurred())

	in := icinga.NewNodeHost(nil, nil, f.icingaClient).
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
