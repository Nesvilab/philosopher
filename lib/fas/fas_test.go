package fas_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/prvst/philosopher/lib/fas"
)

var _ = Describe("Fas", func() {

	Context("Testing database parsing", func() {

		It("Accessing workspace", func() {
			e := os.Chdir("../../test/")
			Expect(e).NotTo(HaveOccurred())
		})

		It("Parsing FASTA", func() {
			f, e := ParseFile("db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta")
			Expect(e).NotTo(HaveOccurred())
			Expect(len(f)).To(Equal(40896))
		})

		It("Parsing UniProt FASTA", func() {
			f := ParseUniProtDescriptionMap("db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta")
			Expect(len(f)).To(Equal(20448))
		})

		It("Parsing UniProt FASTA", func() {
			f := ParseUniProtSequencenMap("db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta")
			Expect(len(f)).To(Equal(20448))
		})

		It("Parsing UniProt FASTA", func() {
			f, e := ParseFastaDescription("db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta")
			Expect(e).NotTo(HaveOccurred())
			Expect(len(f)).To(Equal(20448))
		})

		It("Parsing UniProt FASTA", func() {
			f, e := ParseFile("db/uniprot/2019-02-05-td-hsa-reviewed-2019-02-04.fasta")
			Expect(e).NotTo(HaveOccurred())

			f = CleanDatabase(f, "rev_", "cont")
			Expect(e).NotTo(HaveOccurred())
			Expect(len(f)).To(Equal(20448))
		})

	})

})
