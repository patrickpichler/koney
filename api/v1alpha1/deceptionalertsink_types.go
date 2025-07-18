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

// DeceptionAlertSink is the Schema for the deceptionalertsinks API
type DeceptionAlertSink struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the specification of the DeceptionAlertSinkSpec.
	Spec DeceptionAlertSinkSpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// DeceptionAlertSinkList contains a list of DeceptionAlertSink
type DeceptionAlertSinkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeceptionAlertSink `json:"items"`
}

// DeceptionAlertSinkSpec defines the desired state of DeceptionAlertSink
type DeceptionAlertSinkSpec struct {
	// Dynatrace describes how to send alerts to Dynatrace
	Dynatrace DynatraceSinkSpec `json:"dynatrace,omitempty" yaml:"dynatrace,omitempty"`
}

type DynatraceSinkSpec struct {
	// SecretName references the name of a secret holding `apiToken` and `apiUrl` to connect to the Dynatrace environment.
	SecretName string `json:"secretName,omitempty" yaml:"secretName,omitempty"`

	// Severity describes the severity level upong ingest to Dynatrace.
	// +kubebuilder:validation:Enum=CRITICAL;HIGH;MEDIUM;LOW
	// +optional
	// +kubebuilder:default="HIGH"
	Severity string `json:"severity,omitempty" yaml:"severity,omitempty"`
}

func init() {
	SchemeBuilder.Register(&DeceptionAlertSink{}, &DeceptionAlertSinkList{})
}
