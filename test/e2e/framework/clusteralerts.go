package framework

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	kutil "github.com/appscode/kutil/searchlight/v1alpha1"
	tapi "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
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
	return f.extClient.ClusterAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) TryPatchClusterAlert(meta metav1.ObjectMeta, transform func(*tapi.ClusterAlert) *tapi.ClusterAlert) (*tapi.ClusterAlert, error) {
	return kutil.TryPatchClusterAlert(f.extClient, meta, transform)
}

func (f *Framework) DeleteClusterAlert(meta metav1.ObjectMeta) error {
	return f.extClient.ClusterAlerts(meta.Namespace).Delete(meta.Name, &metav1.DeleteOptions{})
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
