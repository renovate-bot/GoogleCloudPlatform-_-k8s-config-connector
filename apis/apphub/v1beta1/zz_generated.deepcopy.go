//go:build !ignore_autogenerated

// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	refsv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/apis/refs/v1beta1"
	k8sv1beta1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/k8s/v1beta1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppHubApplication) DeepCopyInto(out *AppHubApplication) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppHubApplication.
func (in *AppHubApplication) DeepCopy() *AppHubApplication {
	if in == nil {
		return nil
	}
	out := new(AppHubApplication)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppHubApplication) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppHubApplicationList) DeepCopyInto(out *AppHubApplicationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]AppHubApplication, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppHubApplicationList.
func (in *AppHubApplicationList) DeepCopy() *AppHubApplicationList {
	if in == nil {
		return nil
	}
	out := new(AppHubApplicationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AppHubApplicationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppHubApplicationObservedState) DeepCopyInto(out *AppHubApplicationObservedState) {
	*out = *in
	if in.CreateTime != nil {
		in, out := &in.CreateTime, &out.CreateTime
		*out = new(string)
		**out = **in
	}
	if in.UpdateTime != nil {
		in, out := &in.UpdateTime, &out.UpdateTime
		*out = new(string)
		**out = **in
	}
	if in.Uid != nil {
		in, out := &in.Uid, &out.Uid
		*out = new(string)
		**out = **in
	}
	if in.State != nil {
		in, out := &in.State, &out.State
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppHubApplicationObservedState.
func (in *AppHubApplicationObservedState) DeepCopy() *AppHubApplicationObservedState {
	if in == nil {
		return nil
	}
	out := new(AppHubApplicationObservedState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppHubApplicationSpec) DeepCopyInto(out *AppHubApplicationSpec) {
	*out = *in
	if in.DisplayName != nil {
		in, out := &in.DisplayName, &out.DisplayName
		*out = new(string)
		**out = **in
	}
	if in.Description != nil {
		in, out := &in.Description, &out.Description
		*out = new(string)
		**out = **in
	}
	if in.Attributes != nil {
		in, out := &in.Attributes, &out.Attributes
		*out = new(Attributes)
		(*in).DeepCopyInto(*out)
	}
	if in.Scope != nil {
		in, out := &in.Scope, &out.Scope
		*out = new(Scope)
		(*in).DeepCopyInto(*out)
	}
	if in.Parent != nil {
		in, out := &in.Parent, &out.Parent
		*out = new(Parent)
		(*in).DeepCopyInto(*out)
	}
	if in.ResourceID != nil {
		in, out := &in.ResourceID, &out.ResourceID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppHubApplicationSpec.
func (in *AppHubApplicationSpec) DeepCopy() *AppHubApplicationSpec {
	if in == nil {
		return nil
	}
	out := new(AppHubApplicationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AppHubApplicationStatus) DeepCopyInto(out *AppHubApplicationStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]k8sv1beta1.Condition, len(*in))
		copy(*out, *in)
	}
	if in.ObservedGeneration != nil {
		in, out := &in.ObservedGeneration, &out.ObservedGeneration
		*out = new(int64)
		**out = **in
	}
	if in.ExternalRef != nil {
		in, out := &in.ExternalRef, &out.ExternalRef
		*out = new(string)
		**out = **in
	}
	if in.ObservedState != nil {
		in, out := &in.ObservedState, &out.ObservedState
		*out = new(AppHubApplicationObservedState)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AppHubApplicationStatus.
func (in *AppHubApplicationStatus) DeepCopy() *AppHubApplicationStatus {
	if in == nil {
		return nil
	}
	out := new(AppHubApplicationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationIdentity) DeepCopyInto(out *ApplicationIdentity) {
	*out = *in
	if in.parent != nil {
		in, out := &in.parent, &out.parent
		*out = new(ApplicationParent)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationIdentity.
func (in *ApplicationIdentity) DeepCopy() *ApplicationIdentity {
	if in == nil {
		return nil
	}
	out := new(ApplicationIdentity)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationParent) DeepCopyInto(out *ApplicationParent) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationParent.
func (in *ApplicationParent) DeepCopy() *ApplicationParent {
	if in == nil {
		return nil
	}
	out := new(ApplicationParent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApplicationRef) DeepCopyInto(out *ApplicationRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApplicationRef.
func (in *ApplicationRef) DeepCopy() *ApplicationRef {
	if in == nil {
		return nil
	}
	out := new(ApplicationRef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Attributes) DeepCopyInto(out *Attributes) {
	*out = *in
	if in.Criticality != nil {
		in, out := &in.Criticality, &out.Criticality
		*out = new(Criticality)
		(*in).DeepCopyInto(*out)
	}
	if in.Environment != nil {
		in, out := &in.Environment, &out.Environment
		*out = new(Environment)
		(*in).DeepCopyInto(*out)
	}
	if in.DeveloperOwners != nil {
		in, out := &in.DeveloperOwners, &out.DeveloperOwners
		*out = make([]ContactInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.OperatorOwners != nil {
		in, out := &in.OperatorOwners, &out.OperatorOwners
		*out = make([]ContactInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.BusinessOwners != nil {
		in, out := &in.BusinessOwners, &out.BusinessOwners
		*out = make([]ContactInfo, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Attributes.
func (in *Attributes) DeepCopy() *Attributes {
	if in == nil {
		return nil
	}
	out := new(Attributes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContactInfo) DeepCopyInto(out *ContactInfo) {
	*out = *in
	if in.DisplayName != nil {
		in, out := &in.DisplayName, &out.DisplayName
		*out = new(string)
		**out = **in
	}
	if in.Email != nil {
		in, out := &in.Email, &out.Email
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContactInfo.
func (in *ContactInfo) DeepCopy() *ContactInfo {
	if in == nil {
		return nil
	}
	out := new(ContactInfo)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Criticality) DeepCopyInto(out *Criticality) {
	*out = *in
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Criticality.
func (in *Criticality) DeepCopy() *Criticality {
	if in == nil {
		return nil
	}
	out := new(Criticality)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Environment) DeepCopyInto(out *Environment) {
	*out = *in
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Environment.
func (in *Environment) DeepCopy() *Environment {
	if in == nil {
		return nil
	}
	out := new(Environment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Parent) DeepCopyInto(out *Parent) {
	*out = *in
	if in.ProjectRef != nil {
		in, out := &in.ProjectRef, &out.ProjectRef
		*out = new(refsv1beta1.ProjectRef)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Parent.
func (in *Parent) DeepCopy() *Parent {
	if in == nil {
		return nil
	}
	out := new(Parent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Scope) DeepCopyInto(out *Scope) {
	*out = *in
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scope.
func (in *Scope) DeepCopy() *Scope {
	if in == nil {
		return nil
	}
	out := new(Scope)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadProperties) DeepCopyInto(out *WorkloadProperties) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadProperties.
func (in *WorkloadProperties) DeepCopy() *WorkloadProperties {
	if in == nil {
		return nil
	}
	out := new(WorkloadProperties)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadPropertiesObservedState) DeepCopyInto(out *WorkloadPropertiesObservedState) {
	*out = *in
	if in.GcpProject != nil {
		in, out := &in.GcpProject, &out.GcpProject
		*out = new(string)
		**out = **in
	}
	if in.Location != nil {
		in, out := &in.Location, &out.Location
		*out = new(string)
		**out = **in
	}
	if in.Zone != nil {
		in, out := &in.Zone, &out.Zone
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadPropertiesObservedState.
func (in *WorkloadPropertiesObservedState) DeepCopy() *WorkloadPropertiesObservedState {
	if in == nil {
		return nil
	}
	out := new(WorkloadPropertiesObservedState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadReference) DeepCopyInto(out *WorkloadReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadReference.
func (in *WorkloadReference) DeepCopy() *WorkloadReference {
	if in == nil {
		return nil
	}
	out := new(WorkloadReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *WorkloadReferenceObservedState) DeepCopyInto(out *WorkloadReferenceObservedState) {
	*out = *in
	if in.URI != nil {
		in, out := &in.URI, &out.URI
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new WorkloadReferenceObservedState.
func (in *WorkloadReferenceObservedState) DeepCopy() *WorkloadReferenceObservedState {
	if in == nil {
		return nil
	}
	out := new(WorkloadReferenceObservedState)
	in.DeepCopyInto(out)
	return out
}
