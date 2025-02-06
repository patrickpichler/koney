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

package utils

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// GetPodCondition looks for the conditionType in the pod conditions and returns its status.
// If no condition of this type is present, we return Unknown.
func GetPodCondition(conditions *[]corev1.PodCondition, conditionType corev1.PodConditionType) corev1.ConditionStatus {
	for _, condition := range *conditions {
		if condition.Type == conditionType {
			return condition.Status
		}
	}
	return corev1.ConditionUnknown
}

// GetDeploymentCondition looks for the conditionType in the deployment conditions and returns its status.
// If no condition of this type is present, we return Unknown.
func GetDeploymentCondition(conditions *[]appsv1.DeploymentCondition, conditionType appsv1.DeploymentConditionType) corev1.ConditionStatus {
	for _, condition := range *conditions {
		if condition.Type == conditionType {
			return condition.Status
		}
	}
	return corev1.ConditionUnknown
}
