package wrk_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/prvst/philosopher/lib/err"
	. "github.com/prvst/philosopher/lib/wrk"
)

var _ = Describe("Wrk", func() {

	Context("Testing workspace management", func() {

		It("Init", func() {
			Init("0000", "0000")
		})

		It("Clean", func() {
			e := Clean()
			Expect(e).NotTo(HaveOccurred())
			if e != nil {
				Expect(e.Type).To(Equal(err.CannotDeleteMetaDirectory))
				Expect(e.Class).To(Equal(err.FATA))
			}
		})

	})

})
