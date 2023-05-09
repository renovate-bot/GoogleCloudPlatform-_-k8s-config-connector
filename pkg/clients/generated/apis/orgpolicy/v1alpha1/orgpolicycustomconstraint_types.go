// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Config Connector and manual
//     changes will be clobbered when the file is regenerated.
//
// ----------------------------------------------------------------------------

// *** DISCLAIMER ***
// Config Connector's go-client for CRDs is currently in ALPHA, which means
// that future versions of the go-client may include breaking changes.
// Please try it out and give us feedback!

package v1alpha1

import (
	"github.com/GoogleCloudPlatform/k8s-config-connector/pkg/clients/generated/apis/k8s/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OrgPolicyCustomConstraintSpec struct {
	/* The action to take if the condition is met. Possible values: ["ALLOW", "DENY"]. */
	ActionType string `json:"actionType"`

	/* A CEL condition that refers to a supported service resource, for example 'resource.management.autoUpgrade == false'. For details about CEL usage, see [Common Expression Language](https://cloud.google.com/resource-manager/docs/organization-policy/creating-managing-custom-constraints#common_expression_language). */
	Condition string `json:"condition"`

	/* A human-friendly description of the constraint to display as an error message when the policy is violated. */
	// +optional
	Description *string `json:"description,omitempty"`

	/* A human-friendly name for the constraint. */
	// +optional
	DisplayName *string `json:"displayName,omitempty"`

	/* A list of RESTful methods for which to enforce the constraint. Can be 'CREATE', 'UPDATE', or both. Not all Google Cloud services support both methods. To see supported methods for each service, find the service in [Supported services](https://cloud.google.com/resource-manager/docs/organization-policy/custom-constraint-supported-services). */
	MethodTypes []string `json:"methodTypes"`

	/* Immutable. The parent of the resource, an organization. Format should be 'organizations/{organization_id}'. */
	Parent string `json:"parent"`

	/* Immutable. Optional. The name of the resource. Used for creation and acquisition. When unset, the value of `metadata.name` is used as the default. */
	// +optional
	ResourceID *string `json:"resourceID,omitempty"`

	/* Immutable. Immutable. The fully qualified name of the Google Cloud REST resource containing the object and field you want to restrict. For example, 'container.googleapis.com/NodePool'. */
	ResourceTypes []string `json:"resourceTypes"`
}

type OrgPolicyCustomConstraintStatus struct {
	/* Conditions represent the latest available observations of the
	   OrgPolicyCustomConstraint's current state. */
	Conditions []v1alpha1.Condition `json:"conditions,omitempty"`
	/* ObservedGeneration is the generation of the resource that was most recently observed by the Config Connector controller. If this is equal to metadata.generation, then that means that the current reported status reflects the most recent desired state of the resource. */
	// +optional
	ObservedGeneration *int `json:"observedGeneration,omitempty"`

	/* Output only. The timestamp representing when the constraint was last updated. */
	// +optional
	UpdateTime *string `json:"updateTime,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OrgPolicyCustomConstraint is the Schema for the orgpolicy API
// +k8s:openapi-gen=true
type OrgPolicyCustomConstraint struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OrgPolicyCustomConstraintSpec   `json:"spec,omitempty"`
	Status OrgPolicyCustomConstraintStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OrgPolicyCustomConstraintList contains a list of OrgPolicyCustomConstraint
type OrgPolicyCustomConstraintList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OrgPolicyCustomConstraint `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OrgPolicyCustomConstraint{}, &OrgPolicyCustomConstraintList{})
}