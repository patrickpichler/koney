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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	testNamespace = "test-namespace"
	testCrdName   = "test-crd-name"

	fooType    = "FooType"
	fooReason  = "FooReason"
	fooMessage = "FooMessage"

	barType      = "BarType"
	barReasonOne = "BarReason_1"
	barReasonTwo = "BarReason_2"
	barMessage   = "BarMessage"
)

var (
	deceptionPolicy *DeceptionPolicy
	fooCondition    DeceptionPolicyCondition
	barCondition    DeceptionPolicyCondition
)

func initializeTestConditions() {
	fooCondition = DeceptionPolicyCondition{
		Type:    fooType,
		Status:  metav1.ConditionTrue,
		Reason:  fooReason,
		Message: fooMessage,
	}
	barCondition = DeceptionPolicyCondition{
		Type:    barType,
		Status:  metav1.ConditionTrue,
		Reason:  barReasonOne,
		Message: barMessage,
	}
}

func resetDeceptionPolicy() {
	deceptionPolicy = &DeceptionPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testCrdName,
			Namespace: testNamespace,
		},
	}
}

var _ = Describe("ContainsCondition", func() {
	BeforeEach(func() {
		resetDeceptionPolicy()
	})

	Context("when no conditions are set", func() {
		It("should return false", func() {
			Expect(deceptionPolicy.Status.ContainsCondition(fooType)).To(BeFalse())
		})
	})

	Context("when the condition is not set", func() {
		It("should return false", func() {
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, fooCondition)

			Expect(deceptionPolicy.Status.ContainsCondition(barType)).To(BeFalse())
		})
	})

	Context("when the condition is set", func() {
		It("should return true", func() {
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, fooCondition)

			Expect(deceptionPolicy.Status.ContainsCondition(fooType)).To(BeTrue())
		})
	})

	Context("when multiple conditions are set", func() {
		It("should return true", func() {
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, fooCondition)
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, barCondition)

			Expect(deceptionPolicy.Status.ContainsCondition(fooType)).To(BeTrue())
			Expect(deceptionPolicy.Status.ContainsCondition(barType)).To(BeTrue())
		})
	})
})

var _ = Describe("PutCondition", func() {
	BeforeEach(func() {
		resetDeceptionPolicy()
	})

	Context("when the condition is not set", func() {
		It("should create the condition", func() {
			dirty := deceptionPolicy.Status.PutCondition(fooType, metav1.ConditionTrue, fooReason, fooMessage)

			Expect(dirty).To(BeTrue())
			Expect(deceptionPolicy.Status.ContainsCondition(fooType)).To(BeTrue())
			Expect(deceptionPolicy.Status.Conditions[0].Status).To(Equal(metav1.ConditionTrue))
			Expect(deceptionPolicy.Status.Conditions[0].Reason).To(Equal(fooReason))
			Expect(deceptionPolicy.Status.Conditions[0].Message).To(Equal(fooMessage))
		})
	})

	Context("when the condition is set", func() {
		It("should update the condition", func() {
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, fooCondition)

			dirty := deceptionPolicy.Status.PutCondition(fooType, metav1.ConditionFalse, barReasonOne, barMessage)

			Expect(dirty).To(BeTrue())
			Expect(deceptionPolicy.Status.ContainsCondition(fooType)).To(BeTrue())
			Expect(deceptionPolicy.Status.Conditions[0].Status).To(Equal(metav1.ConditionFalse))
			Expect(deceptionPolicy.Status.Conditions[0].Reason).To(Equal(barReasonOne))
			Expect(deceptionPolicy.Status.Conditions[0].Message).To(Equal(barMessage))
		})
	})

	Context("when the condition is already set as desired", func() {
		It("should not update the condition", func() {
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, fooCondition)

			dirty := deceptionPolicy.Status.PutCondition(fooType, metav1.ConditionTrue, fooReason, fooMessage)

			Expect(dirty).To(BeFalse())
			Expect(deceptionPolicy.Status.ContainsCondition(fooType)).To(BeTrue())
			Expect(deceptionPolicy.Status.Conditions[0].Status).To(Equal(metav1.ConditionTrue))
			Expect(deceptionPolicy.Status.Conditions[0].Reason).To(Equal(fooReason))
			Expect(deceptionPolicy.Status.Conditions[0].Message).To(Equal(fooMessage))
		})
	})

	Context("when multiple conditions are set", func() {
		It("should update the correct condition", func() {
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, fooCondition)
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, barCondition)

			dirty := deceptionPolicy.Status.PutCondition(barType, metav1.ConditionFalse, barReasonTwo, "0/1 decoys deployed (0 skipped)")

			Expect(dirty).To(BeTrue())
			Expect(deceptionPolicy.Status.ContainsCondition(barType)).To(BeTrue())
			Expect(deceptionPolicy.Status.Conditions[1].Status).To(Equal(metav1.ConditionFalse))
			Expect(deceptionPolicy.Status.Conditions[1].Reason).To(Equal(barReasonTwo))
			Expect(deceptionPolicy.Status.Conditions[1].Message).To(Equal("0/1 decoys deployed (0 skipped)"))

			// other conditions should not be affected
			Expect(deceptionPolicy.Status.ContainsCondition(fooType)).To(BeTrue())
			Expect(deceptionPolicy.Status.Conditions[0].Status).To(Equal(metav1.ConditionTrue))
			Expect(deceptionPolicy.Status.Conditions[0].Reason).To(Equal(fooReason))
			Expect(deceptionPolicy.Status.Conditions[0].Message).To(Equal(fooMessage))
		})
	})
})

var _ = Describe("Equals", func() {
	BeforeEach(func() {
		resetDeceptionPolicy()
	})

	Context("when objects are equal", func() {
		It("should return true", func() {
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, fooCondition)

			Expect(fooCondition.Equals(&fooCondition)).To(BeTrue())
		})
	})

	Context("when objects are not equal", func() {
		It("should return false", func() {
			deceptionPolicy.Status.Conditions = append(deceptionPolicy.Status.Conditions, fooCondition)

			Expect(fooCondition.Equals(&barCondition)).To(BeFalse())
		})
	})
})
