//go:build !ignore_autogenerated

/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"go.goms.io/fleet/apis/placement/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AfterStageTask) DeepCopyInto(out *AfterStageTask) {
	*out = *in
	out.WaitTime = in.WaitTime
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AfterStageTask.
func (in *AfterStageTask) DeepCopy() *AfterStageTask {
	if in == nil {
		return nil
	}
	out := new(AfterStageTask)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AfterStageTaskStatus) DeepCopyInto(out *AfterStageTaskStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AfterStageTaskStatus.
func (in *AfterStageTaskStatus) DeepCopy() *AfterStageTaskStatus {
	if in == nil {
		return nil
	}
	out := new(AfterStageTaskStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApprovalRequest) DeepCopyInto(out *ApprovalRequest) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApprovalRequest.
func (in *ApprovalRequest) DeepCopy() *ApprovalRequest {
	if in == nil {
		return nil
	}
	out := new(ApprovalRequest)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ApprovalRequest) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApprovalRequestList) DeepCopyInto(out *ApprovalRequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ApprovalRequest, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApprovalRequestList.
func (in *ApprovalRequestList) DeepCopy() *ApprovalRequestList {
	if in == nil {
		return nil
	}
	out := new(ApprovalRequestList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ApprovalRequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApprovalRequestSpec) DeepCopyInto(out *ApprovalRequestSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApprovalRequestSpec.
func (in *ApprovalRequestSpec) DeepCopy() *ApprovalRequestSpec {
	if in == nil {
		return nil
	}
	out := new(ApprovalRequestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ApprovalRequestStatus) DeepCopyInto(out *ApprovalRequestStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ApprovalRequestStatus.
func (in *ApprovalRequestStatus) DeepCopy() *ApprovalRequestStatus {
	if in == nil {
		return nil
	}
	out := new(ApprovalRequestStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceOverride) DeepCopyInto(out *ClusterResourceOverride) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceOverride.
func (in *ClusterResourceOverride) DeepCopy() *ClusterResourceOverride {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceOverride)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourceOverride) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceOverrideList) DeepCopyInto(out *ClusterResourceOverrideList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterResourceOverride, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceOverrideList.
func (in *ClusterResourceOverrideList) DeepCopy() *ClusterResourceOverrideList {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceOverrideList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourceOverrideList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceOverrideSnapshot) DeepCopyInto(out *ClusterResourceOverrideSnapshot) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceOverrideSnapshot.
func (in *ClusterResourceOverrideSnapshot) DeepCopy() *ClusterResourceOverrideSnapshot {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceOverrideSnapshot)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourceOverrideSnapshot) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceOverrideSnapshotList) DeepCopyInto(out *ClusterResourceOverrideSnapshotList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterResourceOverrideSnapshot, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceOverrideSnapshotList.
func (in *ClusterResourceOverrideSnapshotList) DeepCopy() *ClusterResourceOverrideSnapshotList {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceOverrideSnapshotList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourceOverrideSnapshotList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceOverrideSnapshotSpec) DeepCopyInto(out *ClusterResourceOverrideSnapshotSpec) {
	*out = *in
	in.OverrideSpec.DeepCopyInto(&out.OverrideSpec)
	if in.OverrideHash != nil {
		in, out := &in.OverrideHash, &out.OverrideHash
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceOverrideSnapshotSpec.
func (in *ClusterResourceOverrideSnapshotSpec) DeepCopy() *ClusterResourceOverrideSnapshotSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceOverrideSnapshotSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourceOverrideSpec) DeepCopyInto(out *ClusterResourceOverrideSpec) {
	*out = *in
	if in.ClusterResourceSelectors != nil {
		in, out := &in.ClusterResourceSelectors, &out.ClusterResourceSelectors
		*out = make([]v1beta1.ClusterResourceSelector, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Policy != nil {
		in, out := &in.Policy, &out.Policy
		*out = new(OverridePolicy)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourceOverrideSpec.
func (in *ClusterResourceOverrideSpec) DeepCopy() *ClusterResourceOverrideSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterResourceOverrideSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourcePlacementDisruptionBudget) DeepCopyInto(out *ClusterResourcePlacementDisruptionBudget) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourcePlacementDisruptionBudget.
func (in *ClusterResourcePlacementDisruptionBudget) DeepCopy() *ClusterResourcePlacementDisruptionBudget {
	if in == nil {
		return nil
	}
	out := new(ClusterResourcePlacementDisruptionBudget)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourcePlacementDisruptionBudget) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourcePlacementDisruptionBudgetList) DeepCopyInto(out *ClusterResourcePlacementDisruptionBudgetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterResourcePlacementDisruptionBudget, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourcePlacementDisruptionBudgetList.
func (in *ClusterResourcePlacementDisruptionBudgetList) DeepCopy() *ClusterResourcePlacementDisruptionBudgetList {
	if in == nil {
		return nil
	}
	out := new(ClusterResourcePlacementDisruptionBudgetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourcePlacementDisruptionBudgetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourcePlacementEviction) DeepCopyInto(out *ClusterResourcePlacementEviction) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourcePlacementEviction.
func (in *ClusterResourcePlacementEviction) DeepCopy() *ClusterResourcePlacementEviction {
	if in == nil {
		return nil
	}
	out := new(ClusterResourcePlacementEviction)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourcePlacementEviction) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterResourcePlacementEvictionList) DeepCopyInto(out *ClusterResourcePlacementEvictionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ClusterResourcePlacementEviction, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterResourcePlacementEvictionList.
func (in *ClusterResourcePlacementEvictionList) DeepCopy() *ClusterResourcePlacementEvictionList {
	if in == nil {
		return nil
	}
	out := new(ClusterResourcePlacementEvictionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterResourcePlacementEvictionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterUpdatingStatus) DeepCopyInto(out *ClusterUpdatingStatus) {
	*out = *in
	if in.ResourceOverrideSnapshots != nil {
		in, out := &in.ResourceOverrideSnapshots, &out.ResourceOverrideSnapshots
		*out = make([]v1beta1.NamespacedName, len(*in))
		copy(*out, *in)
	}
	if in.ClusterResourceOverrideSnapshots != nil {
		in, out := &in.ClusterResourceOverrideSnapshots, &out.ClusterResourceOverrideSnapshots
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterUpdatingStatus.
func (in *ClusterUpdatingStatus) DeepCopy() *ClusterUpdatingStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterUpdatingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JSONPatchOverride) DeepCopyInto(out *JSONPatchOverride) {
	*out = *in
	in.Value.DeepCopyInto(&out.Value)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JSONPatchOverride.
func (in *JSONPatchOverride) DeepCopy() *JSONPatchOverride {
	if in == nil {
		return nil
	}
	out := new(JSONPatchOverride)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OverridePolicy) DeepCopyInto(out *OverridePolicy) {
	*out = *in
	if in.OverrideRules != nil {
		in, out := &in.OverrideRules, &out.OverrideRules
		*out = make([]OverrideRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OverridePolicy.
func (in *OverridePolicy) DeepCopy() *OverridePolicy {
	if in == nil {
		return nil
	}
	out := new(OverridePolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OverrideRule) DeepCopyInto(out *OverrideRule) {
	*out = *in
	if in.ClusterSelector != nil {
		in, out := &in.ClusterSelector, &out.ClusterSelector
		*out = new(v1beta1.ClusterSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.JSONPatchOverrides != nil {
		in, out := &in.JSONPatchOverrides, &out.JSONPatchOverrides
		*out = make([]JSONPatchOverride, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OverrideRule.
func (in *OverrideRule) DeepCopy() *OverrideRule {
	if in == nil {
		return nil
	}
	out := new(OverrideRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementDisruptionBudgetSpec) DeepCopyInto(out *PlacementDisruptionBudgetSpec) {
	*out = *in
	if in.MaxUnavailable != nil {
		in, out := &in.MaxUnavailable, &out.MaxUnavailable
		*out = new(intstr.IntOrString)
		**out = **in
	}
	if in.MinAvailable != nil {
		in, out := &in.MinAvailable, &out.MinAvailable
		*out = new(intstr.IntOrString)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementDisruptionBudgetSpec.
func (in *PlacementDisruptionBudgetSpec) DeepCopy() *PlacementDisruptionBudgetSpec {
	if in == nil {
		return nil
	}
	out := new(PlacementDisruptionBudgetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementDisruptionBudgetStatus) DeepCopyInto(out *PlacementDisruptionBudgetStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementDisruptionBudgetStatus.
func (in *PlacementDisruptionBudgetStatus) DeepCopy() *PlacementDisruptionBudgetStatus {
	if in == nil {
		return nil
	}
	out := new(PlacementDisruptionBudgetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementEvictionSpec) DeepCopyInto(out *PlacementEvictionSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementEvictionSpec.
func (in *PlacementEvictionSpec) DeepCopy() *PlacementEvictionSpec {
	if in == nil {
		return nil
	}
	out := new(PlacementEvictionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementEvictionStatus) DeepCopyInto(out *PlacementEvictionStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementEvictionStatus.
func (in *PlacementEvictionStatus) DeepCopy() *PlacementEvictionStatus {
	if in == nil {
		return nil
	}
	out := new(PlacementEvictionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlacementReference) DeepCopyInto(out *PlacementReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlacementReference.
func (in *PlacementReference) DeepCopy() *PlacementReference {
	if in == nil {
		return nil
	}
	out := new(PlacementReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceOverride) DeepCopyInto(out *ResourceOverride) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceOverride.
func (in *ResourceOverride) DeepCopy() *ResourceOverride {
	if in == nil {
		return nil
	}
	out := new(ResourceOverride)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceOverride) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceOverrideList) DeepCopyInto(out *ResourceOverrideList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceOverride, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceOverrideList.
func (in *ResourceOverrideList) DeepCopy() *ResourceOverrideList {
	if in == nil {
		return nil
	}
	out := new(ResourceOverrideList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceOverrideList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceOverrideSnapshot) DeepCopyInto(out *ResourceOverrideSnapshot) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceOverrideSnapshot.
func (in *ResourceOverrideSnapshot) DeepCopy() *ResourceOverrideSnapshot {
	if in == nil {
		return nil
	}
	out := new(ResourceOverrideSnapshot)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceOverrideSnapshot) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceOverrideSnapshotList) DeepCopyInto(out *ResourceOverrideSnapshotList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ResourceOverrideSnapshot, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceOverrideSnapshotList.
func (in *ResourceOverrideSnapshotList) DeepCopy() *ResourceOverrideSnapshotList {
	if in == nil {
		return nil
	}
	out := new(ResourceOverrideSnapshotList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ResourceOverrideSnapshotList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceOverrideSnapshotSpec) DeepCopyInto(out *ResourceOverrideSnapshotSpec) {
	*out = *in
	in.OverrideSpec.DeepCopyInto(&out.OverrideSpec)
	if in.OverrideHash != nil {
		in, out := &in.OverrideHash, &out.OverrideHash
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceOverrideSnapshotSpec.
func (in *ResourceOverrideSnapshotSpec) DeepCopy() *ResourceOverrideSnapshotSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceOverrideSnapshotSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceOverrideSpec) DeepCopyInto(out *ResourceOverrideSpec) {
	*out = *in
	if in.ResourceSelectors != nil {
		in, out := &in.ResourceSelectors, &out.ResourceSelectors
		*out = make([]ResourceSelector, len(*in))
		copy(*out, *in)
	}
	if in.Policy != nil {
		in, out := &in.Policy, &out.Policy
		*out = new(OverridePolicy)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceOverrideSpec.
func (in *ResourceOverrideSpec) DeepCopy() *ResourceOverrideSpec {
	if in == nil {
		return nil
	}
	out := new(ResourceOverrideSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceSelector) DeepCopyInto(out *ResourceSelector) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceSelector.
func (in *ResourceSelector) DeepCopy() *ResourceSelector {
	if in == nil {
		return nil
	}
	out := new(ResourceSelector)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageConfig) DeepCopyInto(out *StageConfig) {
	*out = *in
	if in.LabelSelector != nil {
		in, out := &in.LabelSelector, &out.LabelSelector
		*out = new(v1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.SortingLabelKey != nil {
		in, out := &in.SortingLabelKey, &out.SortingLabelKey
		*out = new(string)
		**out = **in
	}
	if in.AfterStageTasks != nil {
		in, out := &in.AfterStageTasks, &out.AfterStageTasks
		*out = make([]AfterStageTask, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageConfig.
func (in *StageConfig) DeepCopy() *StageConfig {
	if in == nil {
		return nil
	}
	out := new(StageConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StageUpdatingStatus) DeepCopyInto(out *StageUpdatingStatus) {
	*out = *in
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make([]ClusterUpdatingStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.AfterStageTaskStatus != nil {
		in, out := &in.AfterStageTaskStatus, &out.AfterStageTaskStatus
		*out = make([]AfterStageTaskStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = (*in).DeepCopy()
	}
	if in.EndTime != nil {
		in, out := &in.EndTime, &out.EndTime
		*out = (*in).DeepCopy()
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StageUpdatingStatus.
func (in *StageUpdatingStatus) DeepCopy() *StageUpdatingStatus {
	if in == nil {
		return nil
	}
	out := new(StageUpdatingStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagedUpdateRun) DeepCopyInto(out *StagedUpdateRun) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagedUpdateRun.
func (in *StagedUpdateRun) DeepCopy() *StagedUpdateRun {
	if in == nil {
		return nil
	}
	out := new(StagedUpdateRun)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StagedUpdateRun) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagedUpdateRunList) DeepCopyInto(out *StagedUpdateRunList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StagedUpdateRun, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagedUpdateRunList.
func (in *StagedUpdateRunList) DeepCopy() *StagedUpdateRunList {
	if in == nil {
		return nil
	}
	out := new(StagedUpdateRunList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StagedUpdateRunList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagedUpdateRunSpec) DeepCopyInto(out *StagedUpdateRunSpec) {
	*out = *in
	out.PlacementRef = in.PlacementRef
	out.StagedUpdateStrategyRef = in.StagedUpdateStrategyRef
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagedUpdateRunSpec.
func (in *StagedUpdateRunSpec) DeepCopy() *StagedUpdateRunSpec {
	if in == nil {
		return nil
	}
	out := new(StagedUpdateRunSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagedUpdateRunStatus) DeepCopyInto(out *StagedUpdateRunStatus) {
	*out = *in
	in.ApplyStrategy.DeepCopyInto(&out.ApplyStrategy)
	in.StagedUpdateStrategySnapshot.DeepCopyInto(&out.StagedUpdateStrategySnapshot)
	if in.StagesStatus != nil {
		in, out := &in.StagesStatus, &out.StagesStatus
		*out = make([]StageUpdatingStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.DeletionStageStatus.DeepCopyInto(&out.DeletionStageStatus)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagedUpdateRunStatus.
func (in *StagedUpdateRunStatus) DeepCopy() *StagedUpdateRunStatus {
	if in == nil {
		return nil
	}
	out := new(StagedUpdateRunStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagedUpdateStrategy) DeepCopyInto(out *StagedUpdateStrategy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagedUpdateStrategy.
func (in *StagedUpdateStrategy) DeepCopy() *StagedUpdateStrategy {
	if in == nil {
		return nil
	}
	out := new(StagedUpdateStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StagedUpdateStrategy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagedUpdateStrategyList) DeepCopyInto(out *StagedUpdateStrategyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StagedUpdateStrategy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagedUpdateStrategyList.
func (in *StagedUpdateStrategyList) DeepCopy() *StagedUpdateStrategyList {
	if in == nil {
		return nil
	}
	out := new(StagedUpdateStrategyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *StagedUpdateStrategyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StagedUpdateStrategySpec) DeepCopyInto(out *StagedUpdateStrategySpec) {
	*out = *in
	if in.Stages != nil {
		in, out := &in.Stages, &out.Stages
		*out = make([]StageConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StagedUpdateStrategySpec.
func (in *StagedUpdateStrategySpec) DeepCopy() *StagedUpdateStrategySpec {
	if in == nil {
		return nil
	}
	out := new(StagedUpdateStrategySpec)
	in.DeepCopyInto(out)
	return out
}
