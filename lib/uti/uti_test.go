package uti_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/prvst/philosopher/lib/uti"
)

var _ = Describe("Uti", func() {

	Context("Testing utils functions", func() {

		It("Roud", func() {
			x := Round(5.3557876867, 5, 2)
			Expect(x).To(Equal(5.35))
		})

		It("ToFixed", func() {
			x := ToFixed(5.3557876867, 3)
			Expect(x).To(Equal(5.355))
		})

	})

})
