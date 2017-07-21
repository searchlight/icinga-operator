package framework

import (
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/log"
	tapi "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (f *Framework) UpdateClusterAlert(meta metav1.ObjectMeta, transformer func(tapi.ClusterAlert) tapi.ClusterAlert) (*tapi.ClusterAlert, error) {
	attempt := 0
	for ; attempt < maxAttempts; attempt = attempt + 1 {
		cur, err := f.extClient.ClusterAlerts(meta.Namespace).Get(meta.Name)
		if err != nil {
			return nil, err
		}

		modified := transformer(*cur)
		updated, err := f.extClient.ClusterAlerts(cur.Namespace).Update(&modified)
		if err == nil {
			return updated, nil
		}

		log.Errorf("Attempt %d failed to update ClusterAlert %s@%s due to %s.", attempt, cur.Name, cur.Namespace, err)
		time.Sleep(updateRetryInterval)
	}

	return nil, fmt.Errorf("Failed to update ClusterAlert %s@%s after %d attempts.", meta.Name, meta.Namespace, attempt)
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
