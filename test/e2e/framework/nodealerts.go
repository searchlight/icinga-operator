package framework

import (
	"time"

	"github.com/appscode/go/crypto/rand"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/test/e2e/matcher"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (f *Invocation) NodeAlert() *api.NodeAlert {
	return &api.NodeAlert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      rand.WithUniqSuffix("nodealert"),
			Namespace: f.namespace,
			Labels: map[string]string{
				"app": f.app,
			},
		},
		Spec: api.NodeAlertSpec{
			CheckInterval: metav1.Duration{time.Second * 5},
			Vars:          make(map[string]string),
		},
	}
}

func (f *Framework) CreateNodeAlert(obj *api.NodeAlert) error {
	_, err := f.extClient.MonitoringV1alpha1().NodeAlerts(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) GetNodeAlert(meta metav1.ObjectMeta) (*api.NodeAlert, error) {
	return f.extClient.MonitoringV1alpha1().NodeAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) DeleteNodeAlert(meta metav1.ObjectMeta) error {
	return f.extClient.MonitoringV1alpha1().NodeAlerts(meta.Namespace).Delete(meta.Name, &metav1.DeleteOptions{})
}

func (f *Framework) getNodeAlertObjects(meta metav1.ObjectMeta, nodeAlertSpec api.NodeAlertSpec) ([]icinga.IcingaHost, error) {
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

func (f *Framework) EventuallyNodeAlertIcingaService(meta metav1.ObjectMeta, nodeAlertSpec api.NodeAlertSpec) GomegaAsyncAssertion {
	objectList, err := f.getNodeAlertObjects(meta, nodeAlertSpec)
	Expect(err).NotTo(HaveOccurred())

	in := icinga.NewNodeHost(f.icingaClient).
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
					icingaServiceState.OK++
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

func (f *Framework) CleanNodeAlert() {
	caList, err := f.extClient.MonitoringV1alpha1().NodeAlerts(f.namespace).List(metav1.ListOptions{})
	if err != nil {
		return
	}
	for _, e := range caList.Items {
		f.extClient.MonitoringV1alpha1().NodeAlerts(f.namespace).Delete(e.Name, &metav1.DeleteOptions{})
	}
}
