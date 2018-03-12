// +build !ignore_autogenerated

/*
Copyright 2018 The Searchlight Authors.

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

package incidents

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Acknowledgement) DeepCopyInto(out *Acknowledgement) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.Request != nil {
		in, out := &in.Request, &out.Request
		if *in == nil {
			*out = nil
		} else {
			*out = new(AcknowledgementRequest)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Response != nil {
		in, out := &in.Response, &out.Response
		if *in == nil {
			*out = nil
		} else {
			*out = new(AcknowledgementResponse)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Acknowledgement.
func (in *Acknowledgement) DeepCopy() *Acknowledgement {
	if in == nil {
		return nil
	}
	out := new(Acknowledgement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Acknowledgement) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AcknowledgementRequest) DeepCopyInto(out *AcknowledgementRequest) {
	*out = *in
	if in.Comment != nil {
		in, out := &in.Comment, &out.Comment
		if *in == nil {
			*out = nil
		} else {
			*out = new(string)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AcknowledgementRequest.
func (in *AcknowledgementRequest) DeepCopy() *AcknowledgementRequest {
	if in == nil {
		return nil
	}
	out := new(AcknowledgementRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AcknowledgementResponse) DeepCopyInto(out *AcknowledgementResponse) {
	*out = *in
	in.Timestamp.DeepCopyInto(&out.Timestamp)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AcknowledgementResponse.
func (in *AcknowledgementResponse) DeepCopy() *AcknowledgementResponse {
	if in == nil {
		return nil
	}
	out := new(AcknowledgementResponse)
	in.DeepCopyInto(out)
	return out
}
