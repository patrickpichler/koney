// Copyright (c) 2025 Dynatrace LLC
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// DeceptionPolicy is the Schema for the deceptionpolicies API
type DeceptionPolicy struct {
	metav1.TypeMeta `json:",inline" yaml:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`

	// Spec is the specification of the DeceptionPolicy.
	Spec DeceptionPolicySpec `json:"spec,omitempty" yaml:"spec,omitempty"`

	// Status is the status of the DeceptionPolicy.
	Status DeceptionPolicyStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DeceptionPolicyList contains a list of DeceptionPolicy
type DeceptionPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeceptionPolicy `json:"items"`
}

// DeceptionPolicySpec defines the desired state of DeceptionPolicy
type DeceptionPolicySpec struct {
	// Traps is a list of traps to be deployed by the deception policy.
	// Each trap represents a cyber deception technique.
	Traps []Trap `json:"traps,omitempty" yaml:"traps,omitempty"`

	// StrictValidation is a flag that indicates whether the policy should be strictly validated.
	// If set to true, the traps will be deployed only if all the traps in the policy are valid.
	// If set to false, the valid traps will be deployed even if some of the traps are invalid.
	// By default, it is set to true.
	// +optional
	// +kubebuilder:default:=true
	StrictValidation *bool `json:"strictValidation,omitempty" yaml:"strictValidation,omitempty"`

	// MutateExisting is a flag to also allow adding traps to existing resources.
	// Typically, that means that existing resource definitions will be updated to include the traps.
	// Depending on the decoy and captor deployment strategies, this may require restarting the pods.
	// +optional
	// +kubebuilder:default=true
	MutateExisting *bool `json:"mutateExisting,omitempty" yaml:"mutateExisting,omitempty"`
}

func init() {
	SchemeBuilder.Register(&DeceptionPolicy{}, &DeceptionPolicyList{})
}
