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

package util

import (
	"context"
	"encoding/json"
	"fmt"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	cs "go.searchlight.dev/icinga-operator/client/clientset/versioned/typed/monitoring/v1alpha1"

	jsonpatch "github.com/evanphx/json-patch"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/v2"
	kutil "kmodules.xyz/client-go"
)

func CreateOrPatchPodAlert(ctx context.Context, c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(in *api.PodAlert) *api.PodAlert, opts metav1.PatchOptions) (*api.PodAlert, kutil.VerbType, error) {
	cur, err := c.PodAlerts(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		klog.V(3).Infof("Creating PodAlert %s/%s.", meta.Namespace, meta.Name)
		out, err := c.PodAlerts(meta.Namespace).Create(ctx, transform(&api.PodAlert{
			TypeMeta: metav1.TypeMeta{
				Kind:       api.ResourceKindPodAlert,
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}), metav1.CreateOptions{
			DryRun:       opts.DryRun,
			FieldManager: opts.FieldManager,
		})
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchPodAlert(ctx, c, cur, transform, opts)
}

func PatchPodAlert(ctx context.Context, c cs.MonitoringV1alpha1Interface, cur *api.PodAlert, transform func(*api.PodAlert) *api.PodAlert, opts metav1.PatchOptions) (*api.PodAlert, kutil.VerbType, error) {
	return PatchPodAlertObject(ctx, c, cur, transform(cur.DeepCopy()), opts)
}

func PatchPodAlertObject(ctx context.Context, c cs.MonitoringV1alpha1Interface, cur, mod *api.PodAlert, opts metav1.PatchOptions) (*api.PodAlert, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonpatch.CreateMergePatch(curJson, modJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	klog.V(3).Infof("Patching PodAlert %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.PodAlerts(cur.Namespace).Patch(ctx, cur.Name, types.MergePatchType, patch, opts)
	return out, kutil.VerbPatched, err
}

func TryUpdatePodAlert(ctx context.Context, c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.PodAlert) *api.PodAlert, opts metav1.UpdateOptions) (result *api.PodAlert, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.PodAlerts(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.PodAlerts(cur.Namespace).Update(ctx, transform(cur.DeepCopy()), opts)
			return e2 == nil, nil
		}
		klog.Errorf("Attempt %d failed to update PodAlert %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update PodAlert %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func UpdatePodAlertStatus(
	ctx context.Context,
	c cs.MonitoringV1alpha1Interface,
	meta metav1.ObjectMeta,
	transform func(*api.PodAlertStatus) (types.UID, *api.PodAlertStatus),
	opts metav1.UpdateOptions,
) (result *api.PodAlert, err error) {
	apply := func(x *api.PodAlert) *api.PodAlert {
		uid, updatedStatus := transform(x.Status.DeepCopy())
		// Ignore status update when uid does not match
		if uid != "" && uid != x.UID {
			return x
		}
		out := &api.PodAlert{
			TypeMeta:   x.TypeMeta,
			ObjectMeta: x.ObjectMeta,
			Spec:       x.Spec,
			Status:     *updatedStatus,
		}
		return out
	}

	attempt := 0
	cur, err := c.PodAlerts(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		var e2 error
		result, e2 = c.PodAlerts(meta.Namespace).UpdateStatus(ctx, apply(cur), opts)
		if kerr.IsConflict(e2) {
			latest, e3 := c.PodAlerts(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
			switch {
			case e3 == nil:
				cur = latest
				return false, nil
			case kutil.IsRequestRetryable(e3):
				return false, nil
			default:
				return false, e3
			}
		} else if err != nil && !kutil.IsRequestRetryable(e2) {
			return false, e2
		}
		return e2 == nil, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update status of PodAlert %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
