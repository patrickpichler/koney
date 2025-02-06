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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Contains", func() {
	var slice []string

	Context("With a valid slice", func() {
		BeforeEach(func() {
			slice = []string{"foo", "bar", "baz"}
		})

		It("should find the element", func() {
			Expect(Contains(slice, "foo")).To(BeTrue())
		})

		It("should not find non-existing element", func() {
			Expect(Contains(slice, "other")).To(BeFalse())
		})
	})

	Context("With an empty slice", func() {
		BeforeEach(func() {
			slice = []string{}
		})

		It("should not find anything", func() {
			Expect(Contains(slice, "foo")).To(BeFalse())
		})
	})

})
