// +build !ignore_autogenerated

/*
Copyright 2017 The Kubernetes Authors.

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

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package apiregistration

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	reflect "reflect"
)

func init() {
	SchemeBuilder.Register(RegisterDeepCopies)
}

// RegisterDeepCopies adds deep-copy functions to the given scheme. Public
// to allow building arbitrary schemes.
func RegisterDeepCopies(scheme *runtime.Scheme) error {
	return scheme.AddGeneratedDeepCopyFuncs(
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_apiregistration_APIService, InType: reflect.TypeOf(&APIService{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_apiregistration_APIServiceList, InType: reflect.TypeOf(&APIServiceList{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_apiregistration_APIServiceSpec, InType: reflect.TypeOf(&APIServiceSpec{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_apiregistration_APIServiceStatus, InType: reflect.TypeOf(&APIServiceStatus{})},
		conversion.GeneratedDeepCopyFunc{Fn: DeepCopy_apiregistration_ServiceReference, InType: reflect.TypeOf(&ServiceReference{})},
	)
}

func DeepCopy_apiregistration_APIService(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*APIService)
		out := out.(*APIService)
		*out = *in
		if newVal, err := c.DeepCopy(&in.ObjectMeta); err != nil {
			return err
		} else {
			out.ObjectMeta = *newVal.(*v1.ObjectMeta)
		}
		if err := DeepCopy_apiregistration_APIServiceSpec(&in.Spec, &out.Spec, c); err != nil {
			return err
		}
		return nil
	}
}

func DeepCopy_apiregistration_APIServiceList(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*APIServiceList)
		out := out.(*APIServiceList)
		*out = *in
		if in.Items != nil {
			in, out := &in.Items, &out.Items
			*out = make([]APIService, len(*in))
			for i := range *in {
				if err := DeepCopy_apiregistration_APIService(&(*in)[i], &(*out)[i], c); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func DeepCopy_apiregistration_APIServiceSpec(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*APIServiceSpec)
		out := out.(*APIServiceSpec)
		*out = *in
		if in.CABundle != nil {
			in, out := &in.CABundle, &out.CABundle
			*out = make([]byte, len(*in))
			copy(*out, *in)
		}
		return nil
	}
}

func DeepCopy_apiregistration_APIServiceStatus(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*APIServiceStatus)
		out := out.(*APIServiceStatus)
		*out = *in
		return nil
	}
}

func DeepCopy_apiregistration_ServiceReference(in interface{}, out interface{}, c *conversion.Cloner) error {
	{
		in := in.(*ServiceReference)
		out := out.(*ServiceReference)
		*out = *in
		return nil
	}
}
