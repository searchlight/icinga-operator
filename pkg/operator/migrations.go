/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

import (
	"context"
	"strings"

	"go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/client/clientset/versioned/typed/monitoring/v1alpha1/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
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
	ca, err := op.extClient.MonitoringV1alpha1().ClusterAlerts(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	var errs []error
	for i := range ca.Items {
		_, _, err := util.PatchClusterAlert(context.TODO(), op.extClient.MonitoringV1alpha1(), &ca.Items[i], func(alert *v1alpha1.ClusterAlert) *v1alpha1.ClusterAlert {
			check := strings.Replace(alert.Spec.Check, "_", "-", -1)
			alert.Spec.Check = check
			return alert
		}, metav1.PatchOptions{})
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
	poa, err := op.extClient.MonitoringV1alpha1().PodAlerts(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	var errs []error
	for i := range poa.Items {
		_, _, err := util.PatchPodAlert(context.TODO(), op.extClient.MonitoringV1alpha1(), &poa.Items[i], func(alert *v1alpha1.PodAlert) *v1alpha1.PodAlert {
			check := strings.Replace(alert.Spec.Check, "_", "-", -1)
			alert.Spec.Check = check
			return alert
		}, metav1.PatchOptions{})
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
	noa, err := op.extClient.MonitoringV1alpha1().NodeAlerts(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	var errs []error
	for i := range noa.Items {
		_, _, err := util.PatchNodeAlert(context.TODO(), op.extClient.MonitoringV1alpha1(), &noa.Items[i], func(alert *v1alpha1.NodeAlert) *v1alpha1.NodeAlert {
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
		}, metav1.PatchOptions{})
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return utilerrors.NewAggregate(errs)
	}

	return nil
}
