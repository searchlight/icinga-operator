package framework

import (
	"fmt"
	"time"

	"github.com/appscode/go/crypto/rand"
	"github.com/appscode/go/log"
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
	_, err := f.extClient.NodeAlerts(obj.Namespace).Create(obj)
	return err
}

func (f *Framework) GetNodeAlert(meta metav1.ObjectMeta) (*api.NodeAlert, error) {
	return f.extClient.NodeAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
}

func (f *Framework) UpdateNodeAlert(meta metav1.ObjectMeta, transformer func(api.NodeAlert) api.NodeAlert) (*api.NodeAlert, error) {
	attempt := 0
	for ; attempt < maxAttempts; attempt = attempt + 1 {
		cur, err := f.extClient.NodeAlerts(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}

		modified := transformer(*cur)
		updated, err := f.extClient.NodeAlerts(cur.Namespace).Update(&modified)
		if err == nil {
			return updated, nil
		}

		log.Errorf("Attempt %d failed to update NodeAlert %s@%s due to %s.", attempt, cur.Name, cur.Namespace, err)
		time.Sleep(updateRetryInterval)
	}

	return nil, fmt.Errorf("Failed to update NodeAlert %s@%s after %d attempts.", meta.Name, meta.Namespace, attempt)
}

func (f *Framework) DeleteNodeAlert(meta metav1.ObjectMeta) error {
	return f.extClient.NodeAlerts(meta.Namespace).Delete(meta.Name, &metav1.DeleteOptions{})
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
