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

package v1alpha1

type Receiver struct {
	// For which state notification will be sent
	State string `json:"state,omitempty" protobuf:"bytes,1,opt,name=state"`

	// To whom notification will be sent
	To []string `json:"to,omitempty" protobuf:"bytes,2,rep,name=to"`

	// How this notification will be sent
	Notifier string `json:"notifier,omitempty" protobuf:"bytes,3,opt,name=notifier"`
}
