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
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/apis/k8s/v1alpha1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirestoreDatabase) DeepCopyInto(out *FirestoreDatabase) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirestoreDatabase.
func (in *FirestoreDatabase) DeepCopy() *FirestoreDatabase {
	if in == nil {
		return nil
	}
	out := new(FirestoreDatabase)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FirestoreDatabase) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirestoreDatabaseList) DeepCopyInto(out *FirestoreDatabaseList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FirestoreDatabase, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirestoreDatabaseList.
func (in *FirestoreDatabaseList) DeepCopy() *FirestoreDatabaseList {
	if in == nil {
		return nil
	}
	out := new(FirestoreDatabaseList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FirestoreDatabaseList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirestoreDatabaseObservedState) DeepCopyInto(out *FirestoreDatabaseObservedState) {
	*out = *in
	if in.Uid != nil {
		in, out := &in.Uid, &out.Uid
		*out = new(string)
		**out = **in
	}
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
	if in.VersionRetentionPeriod != nil {
		in, out := &in.VersionRetentionPeriod, &out.VersionRetentionPeriod
		*out = new(string)
		**out = **in
	}
	if in.EarliestVersionTime != nil {
		in, out := &in.EarliestVersionTime, &out.EarliestVersionTime
		*out = new(string)
		**out = **in
	}
	if in.KeyPrefix != nil {
		in, out := &in.KeyPrefix, &out.KeyPrefix
		*out = new(string)
		**out = **in
	}
	if in.Etag != nil {
		in, out := &in.Etag, &out.Etag
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirestoreDatabaseObservedState.
func (in *FirestoreDatabaseObservedState) DeepCopy() *FirestoreDatabaseObservedState {
	if in == nil {
		return nil
	}
	out := new(FirestoreDatabaseObservedState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirestoreDatabaseSpec) DeepCopyInto(out *FirestoreDatabaseSpec) {
	*out = *in
	out.ProjectRef = in.ProjectRef
	if in.ResourceID != nil {
		in, out := &in.ResourceID, &out.ResourceID
		*out = new(string)
		**out = **in
	}
	if in.LocationID != nil {
		in, out := &in.LocationID, &out.LocationID
		*out = new(string)
		**out = **in
	}
	if in.ConcurrencyMode != nil {
		in, out := &in.ConcurrencyMode, &out.ConcurrencyMode
		*out = new(string)
		**out = **in
	}
	if in.PointInTimeRecoveryEnablement != nil {
		in, out := &in.PointInTimeRecoveryEnablement, &out.PointInTimeRecoveryEnablement
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirestoreDatabaseSpec.
func (in *FirestoreDatabaseSpec) DeepCopy() *FirestoreDatabaseSpec {
	if in == nil {
		return nil
	}
	out := new(FirestoreDatabaseSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FirestoreDatabaseStatus) DeepCopyInto(out *FirestoreDatabaseStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1alpha1.Condition, len(*in))
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
		*out = new(FirestoreDatabaseObservedState)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FirestoreDatabaseStatus.
func (in *FirestoreDatabaseStatus) DeepCopy() *FirestoreDatabaseStatus {
	if in == nil {
		return nil
	}
	out := new(FirestoreDatabaseStatus)
	in.DeepCopyInto(out)
	return out
}
