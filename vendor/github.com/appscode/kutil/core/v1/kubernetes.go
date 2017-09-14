package v1

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/appscode/kutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func GetGroupVersionKind(v interface{}) schema.GroupVersionKind {
	return apiv1.SchemeGroupVersion.WithKind(kutil.GetKind(v))
}

func AssignTypeKind(v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return fmt.Errorf("%v must be a pointer", v)
	}

	switch u := v.(type) {
	case *apiv1.Pod:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.ReplicationController:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.ConfigMap:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.Secret:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.Service:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.PersistentVolumeClaim:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.PersistentVolume:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.Node:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.ServiceAccount:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.Namespace:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.Endpoints:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.ComponentStatus:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.LimitRange:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	case *apiv1.Event:
		u.APIVersion = apiv1.SchemeGroupVersion.String()
		u.Kind = kutil.GetKind(v)
		return nil
	}
	return errors.New("unknown api object type")
}
