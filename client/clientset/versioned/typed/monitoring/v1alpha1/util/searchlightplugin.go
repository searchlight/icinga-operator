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

func CreateOrPatchSearchlightPlugin(ctx context.Context, c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(in *api.SearchlightPlugin) *api.SearchlightPlugin, opts metav1.PatchOptions) (*api.SearchlightPlugin, kutil.VerbType, error) {
	cur, err := c.SearchlightPlugins().Get(ctx, meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		klog.V(3).Infof("Creating SearchlightPlugin %s/%s.", meta.Namespace, meta.Name)
		out, err := c.SearchlightPlugins().Create(ctx, transform(&api.SearchlightPlugin{
			TypeMeta: metav1.TypeMeta{
				Kind:       "SearchlightPlugin",
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
	return PatchSearchlightPlugin(ctx, c, cur, transform, opts)
}

func PatchSearchlightPlugin(ctx context.Context, c cs.MonitoringV1alpha1Interface, cur *api.SearchlightPlugin, transform func(*api.SearchlightPlugin) *api.SearchlightPlugin, opts metav1.PatchOptions) (*api.SearchlightPlugin, kutil.VerbType, error) {
	return PatchSearchlightPluginObject(ctx, c, cur, transform(cur.DeepCopy()), opts)
}

func PatchSearchlightPluginObject(ctx context.Context, c cs.MonitoringV1alpha1Interface, cur, mod *api.SearchlightPlugin, opts metav1.PatchOptions) (*api.SearchlightPlugin, kutil.VerbType, error) {
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
	klog.V(3).Infof("Patching SearchlightPlugin %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.SearchlightPlugins().Patch(ctx, cur.Name, types.MergePatchType, patch, opts)
	return out, kutil.VerbPatched, err
}

func TryUpdateSearchlightPlugin(ctx context.Context, c cs.MonitoringV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.SearchlightPlugin) *api.SearchlightPlugin, opts metav1.UpdateOptions) (result *api.SearchlightPlugin, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.SearchlightPlugins().Get(ctx, meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.SearchlightPlugins().Update(ctx, transform(cur.DeepCopy()), opts)
			return e2 == nil, nil
		}
		klog.Errorf("Attempt %d failed to update SearchlightPlugin %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update SearchlightPlugin %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
