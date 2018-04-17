package operator

import (
	"strings"

	utilerrors "github.com/appscode/go/util/errors"
	"github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/client/clientset/versioned/typed/monitoring/v1alpha1/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (op *Operator) MigrateAlerts() error {
	var errs []error
	if err := op.MigrateClusterAlerts(); err != nil {
		errs = append(errs, err)
	}
	if err := op.MigratePodAlert(); err != nil {
		errs = append(errs, err)
	}
	if err := op.MigrateNodeAlert(); err != nil {
		errs = append(errs, err)
	}

	return utilerrors.NewAggregate(errs)
}

func (op *Operator) MigrateClusterAlerts() error {
	ca, err := op.extClient.MonitoringV1alpha1().ClusterAlerts(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	var errs []error
	for i := range ca.Items {
		_, _, err := util.PatchClusterAlert(op.extClient.MonitoringV1alpha1(), &ca.Items[i], func(alert *v1alpha1.ClusterAlert) *v1alpha1.ClusterAlert {
			check := strings.Replace(alert.Spec.Check, "_", "-", -1)
			alert.Spec.Check = check
			return alert
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	return nil
}

func (op *Operator) MigratePodAlert() error {
	poa, err := op.extClient.MonitoringV1alpha1().PodAlerts(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	var errs []error
	for i := range poa.Items {
		_, _, err := util.PatchPodAlert(op.extClient.MonitoringV1alpha1(), &poa.Items[i], func(alert *v1alpha1.PodAlert) *v1alpha1.PodAlert {
			check := strings.Replace(alert.Spec.Check, "_", "-", -1)
			alert.Spec.Check = check
			return alert
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	return nil
}

func (op *Operator) MigrateNodeAlert() error {
	noa, err := op.extClient.MonitoringV1alpha1().NodeAlerts(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	var errs []error
	for i := range noa.Items {
		_, _, err := util.PatchNodeAlert(op.extClient.MonitoringV1alpha1(), &noa.Items[i], func(alert *v1alpha1.NodeAlert) *v1alpha1.NodeAlert {
			check := strings.Replace(alert.Spec.Check, "_", "-", -1)
			alert.Spec.Check = check

			if check == api.CheckNodeVolume {
				mp, found := alert.Spec.Vars["mountpoint"]
				if found {
					delete(alert.Spec.Vars, "mountpoint")
					alert.Spec.Vars["mountPoint"] = mp
				}
			}

			return alert
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	return nil
}
