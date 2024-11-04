/*
Copyright 2021 Nho Luong DevOps All rights reserved.

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

package kubeconfig

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/minikube/pkg/version"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// implementing the runtime.Object internally so we can write extensions to kubeconfig

// Extension represents information to identify clusters and contexts
type Extension struct {
	runtime.TypeMeta `json:",inline"`
	Version          string `json:"version"`
	Provider         string `json:"provider"`
	LastUpdate       string `json:"last-update"`
}

// NewExtension returns a minikube formatted kubeconfig's extension block to idenity clusters and contexts
func NewExtension() *Extension {
	return &Extension{
		Provider: "minikube.sigs.k8s.io",
		Version:  version.GetVersion(),
		// time format matching other RFC in notify.go
		LastUpdate: time.Now().Format(time.RFC1123)}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Extension.
func (in *Extension) DeepCopy() *Extension {
	if in == nil {
		return nil
	}
	out := new(Extension)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Extension) DeepCopyObject() runtime.Object {
	return in.DeepCopy()
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Extension) DeepCopyInto(out *Extension) {
	*out = *in
	out.TypeMeta = in.TypeMeta
}
