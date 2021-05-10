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
	"errors"
	"fmt"
	"reflect"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"kmodules.xyz/client-go/meta"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return api.SchemeGroupVersion.WithKind(meta.GetKind(v))
}

func AssignTypeKind(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%v must be a pointer", v)
	}

	switch u := v.(type) {
	case *api.ClusterAlert:
		u.APIVersion = api.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *api.NodeAlert:
		u.APIVersion = api.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	case *api.PodAlert:
		u.APIVersion = api.SchemeGroupVersion.String()
		u.Kind = meta.GetKind(v)
		return nil
	}
	return errors.New("unknown api object type")
}
