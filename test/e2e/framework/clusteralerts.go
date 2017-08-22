package framework

import (
	"fmt"
	"time"

	"encoding/json"
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

func (f *Invocation) ClusterAlert() *tapi.ClusterAlert {
	return &tapi.ClusterAlert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("clusteralert"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: tapi.ClusterAlertSpec{
			CheckInterval: metav1.Duration{time.Second * 5},
			Vars:          make(map[string]interface{}),
		},
	}
}

func (f *Framework) CreateClusterAlert(obj *tapi.ClusterAlert) error {
	_, err := f.extClient.ClusterAlerts(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) GetClusterAlert(meta metav1.ObjectMeta) (*tapi.ClusterAlert, error) {
	return f.extClient.ClusterAlerts(meta.Namespace).Get(meta.Name)
}

func (f *Framework) patchClusterAlert(cur *tapi.ClusterAlert, transform func(*tapi.ClusterAlert) *tapi.ClusterAlert) (*tapi.ClusterAlert, error) {
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
	glog.V(5).Infof("Patching ClusterAlert %s@%s with %s.", cur.Name, cur.Namespace, string(pb))
	return f.extClient.ClusterAlerts(cur.Namespace).Patch(cur.Name, types.JSONPatchType, pb)
}

func (f *Framework) TryPatchClusterAlert(meta metav1.ObjectMeta, transform func(*tapi.ClusterAlert) *tapi.ClusterAlert) (*tapi.ClusterAlert, error) {
	attempt := 0
	for ; attempt < kutil.MaxAttempts; attempt = attempt + 1 {
		cur, err := f.extClient.ClusterAlerts(meta.Namespace).Get(meta.Name)
		if kerr.IsNotFound(err) {
			return cur, err
		} else if err == nil {
			return f.patchClusterAlert(cur, transform)
		}
		glog.Errorf("Attempt %d failed to patch ClusterAlert %s@%s due to %s.", attempt, cur.Name, cur.Namespace, err)
		time.Sleep(kutil.RetryInterval)
	}
	return nil, fmt.Errorf("Failed to patch ClusterAlert %s@%s after %d attempts.", meta.Name, meta.Namespace, attempt)
}

func (f *Framework) DeleteClusterAlert(meta metav1.ObjectMeta) error {
	return f.extClient.ClusterAlerts(meta.Namespace).Delete(meta.Name)
}

func (f *Framework) getClusterAlertObjects(meta metav1.ObjectMeta, clusterAlertSpec tapi.ClusterAlertSpec) ([]icinga.IcingaHost, error) {
	objectList := []icinga.IcingaHost{
		{
			Type:           icinga.TypeCluster,
			AlertNamespace: meta.Namespace,
		},
	}
	return objectList, nil
}

func (f *Framework) EventuallyClusterAlertIcingaService(meta metav1.ObjectMeta, nodeAlertSpec tapi.ClusterAlertSpec) GomegaAsyncAssertion {
	objectList, err := f.getClusterAlertObjects(meta, nodeAlertSpec)
	Expect(err).NotTo(HaveOccurred())

	in := icinga.NewClusterHost(nil, nil, f.icingaClient).
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
