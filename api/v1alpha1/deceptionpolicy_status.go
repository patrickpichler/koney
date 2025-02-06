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

// DeceptionPolicyStatus defines the observed state of DeceptionPolicy
type DeceptionPolicyStatus struct {
	// Conditions is an array of conditions that the DeceptionPolicy can be in.
	// +listType=map
	// +listMapKey=type
	Conditions []DeceptionPolicyCondition `json:"conditions" yaml:"conditions"`
}

// DeceptionPolicyCondition describes the state of one aspect of a DeceptionPolicy at a certain point.
type DeceptionPolicyCondition struct {
	// Type of deception policy condition.
	// The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$`
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=316
	Type string `json:"type" yaml:"type"`

	// Status of the condition, one of True, False, Unknown.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=True;False;Unknown
	Status metav1.ConditionStatus `json:"status" yaml:"status"`

	// LastTransitionTime is the last time the condition transitioned from one status to another,
	// i.e., when the underlying condition changed.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	LastTransitionTime metav1.Time `json:"lastTransitionTime" yaml:"lastTransitionTime"`

	// Reason indicates the reason for the condition's last transition.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:Pattern=`^[A-Za-z]([A-Za-z0-9_,:]*[A-Za-z0-9_])?$`
	Reason string `json:"reason" yaml:"reason"`

	// Message is a human-readable explanation indicating details about the transition.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=32768
	Message string `json:"message" yaml:"message"`
}

// ContainsCondition returns true if the DeceptionPolicy status contains a condition with the provided type.
func (status *DeceptionPolicyStatus) ContainsCondition(conditionType string) bool {
	return status.GetCondition(conditionType) != nil
}

// GetCondition returns a pointer to the first condition with the provided type, if it exists.
func (status *DeceptionPolicyStatus) GetCondition(conditionType string) *DeceptionPolicyCondition {
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return &status.Conditions[i]
		}
	}

	return nil
}

// PutCondition adds a new condition to the DeceptionPolicy status, or updates the first existing condition of the same type, if it exists.
// The function returns true if the conditions were modified as a result of the operation.
func (status *DeceptionPolicyStatus) PutCondition(conditionType string, conditionStatus metav1.ConditionStatus, conditionReason, conditionMessage string) bool {
	return status.PutConditionStruct(DeceptionPolicyCondition{
		Type:               conditionType,
		Status:             conditionStatus,
		LastTransitionTime: metav1.Now(),
		Reason:             conditionReason,
		Message:            conditionMessage,
	})
}

// PutConditionStruct adds a new condition to the DeceptionPolicy status, or updates the first existing condition of the same type, if it exists.
// The function returns true if the conditions were modified as a result of the operation.
func (status *DeceptionPolicyStatus) PutConditionStruct(condition DeceptionPolicyCondition) bool {
	conditionsModified := false

	if existingCondition := status.GetCondition(condition.Type); existingCondition == nil {
		status.Conditions = append(status.Conditions, condition)
		conditionsModified = true
	} else if !condition.Equals(existingCondition) {
		existingCondition.Status = condition.Status
		existingCondition.LastTransitionTime = condition.LastTransitionTime
		existingCondition.Reason = condition.Reason
		existingCondition.Message = condition.Message

		conditionsModified = true
	}

	return conditionsModified
}

// Equals returns true if the conditions are equal (excluding LastTransitionTime).
func (condition *DeceptionPolicyCondition) Equals(other *DeceptionPolicyCondition) bool {
	if condition == other {
		return true
	}
	if condition.Type != other.Type {
		return false
	}
	if condition.Status != other.Status {
		return false
	}
	if condition.Reason != other.Reason {
		return false
	}
	if condition.Message != other.Message {
		return false
	}

	return true
}
