package v1alpha1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/appscode/kutil"
	aci "github.com/appscode/searchlight/api"
	tcs "github.com/appscode/searchlight/client/clientset"
	"github.com/golang/glog"
	"github.com/mattbaird/jsonpatch"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func CreateOrPatchPodAlert(c tcs.ExtensionInterface, meta metav1.ObjectMeta, transform func(alert *aci.PodAlert) *aci.PodAlert) (*aci.PodAlert, error) {
	cur, err := c.PodAlerts(meta.Namespace).Get(meta.Name)
	if kerr.IsNotFound(err) {
		return c.PodAlerts(meta.Namespace).Create(transform(&aci.PodAlert{ObjectMeta: meta}))
	} else if err != nil {
		return nil, err
	}
	return PatchPodAlert(c, cur, transform)
}

func PatchPodAlert(c tcs.ExtensionInterface, cur *aci.PodAlert, transform func(*aci.PodAlert) *aci.PodAlert) (*aci.PodAlert, error) {
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
	glog.V(5).Infof("Patching PodAlert %s@%s with %s.", cur.Name, cur.Namespace, string(pb))
	result, err := c.PodAlerts(cur.Namespace).Patch(cur.Name, types.JSONPatchType, pb)
	return result, err
}

func TryPatchPodAlert(c tcs.ExtensionInterface, meta metav1.ObjectMeta, transform func(*aci.PodAlert) *aci.PodAlert) (*aci.PodAlert, error) {
	attempt := 0
	for ; attempt < kutil.MaxAttempts; attempt = attempt + 1 {
		cur, err := c.PodAlerts(meta.Namespace).Get(meta.Name)
		if kerr.IsNotFound(err) {
			return cur, err
		} else if err == nil {
			return PatchPodAlert(c, cur, transform)
		}
		glog.Errorf("Attempt %d failed to patch PodAlert %s@%s due to %s.", attempt, cur.Name, cur.Namespace, err)
		time.Sleep(kutil.RetryInterval)
	}
	return nil, fmt.Errorf("Failed to patch PodAlert %s@%s after %d attempts.", meta.Name, meta.Namespace, attempt)
}

func TryUpdatePodAlert(c tcs.ExtensionInterface, meta metav1.ObjectMeta, transform func(*aci.PodAlert) *aci.PodAlert) (*aci.PodAlert, error) {
	attempt := 0
	for ; attempt < kutil.MaxAttempts; attempt = attempt + 1 {
		cur, err := c.PodAlerts(meta.Namespace).Get(meta.Name)
		if kerr.IsNotFound(err) {
			return cur, err
		} else if err == nil {
			oJson, err := json.Marshal(cur)
			if err != nil {
				return nil, err
			}
			modified := transform(cur)
			mJson, err := json.Marshal(modified)
			if err != nil {
				return nil, err
			}
			if bytes.Equal(oJson, mJson) {
				return cur, err
			}

			result, err := c.PodAlerts(cur.Namespace).Update(transform(cur))
			return result, err
		}
		glog.Errorf("Attempt %d failed to update PodAlert %s@%s due to %s.", attempt, cur.Name, cur.Namespace, err)
		time.Sleep(kutil.RetryInterval)
	}
	return nil, fmt.Errorf("Failed to update PodAlert %s@%s after %d attempts.", meta.Name, meta.Namespace, attempt)
}
