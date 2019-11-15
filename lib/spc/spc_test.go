package spc_test

import (
	"os"

	"philosopher/lib/id"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Spc", func() {

	Context("Testing pepxml parsing", func() {

		var p id.PepXML
		var list id.PepIDList
		var e error

		It("Accessing workspace", func() {
			e = os.Chdir("../../test/wrksp/")
			Expect(e).NotTo(HaveOccurred())
		})

		It("Reading interact.pep.xml", func() {
			p.Read("interact.pep.xml")
			Expect(e).NotTo(HaveOccurred())
			list = append(list, p.PeptideIdentification...)
		})

		It("Checking data structure length", func() {
			Expect(len(p.PeptideIdentification)).NotTo(Equal(0))
		})

		It("Checking the list size", func() {
			Expect(len(list)).To(Equal(64406))
		})

		It("Checking peptide sequence 1", func() {
			Expect(list[0].Peptide).To(Equal("KGPGRPTGSK"))
		})

		It("Checking Spectrum from search hit 1", func() {
			Expect(list[0].Spectrum).To(Equal("b1906_293T_proteinID_01A_QE3_122212.01882.01882.3"))
		})

		It("Checking AssumedCharge from search hit 1", func() {
			Expect(list[0].AssumedCharge).To(Equal(uint8(3)))
		})

		It("Checking CalcNeutralPepMass from search hit 1", func() {
			Expect(list[0].CalcNeutralPepMass).To(Equal(983.5512))
		})

		It("Checking DiscriminantValue from search hit 1", func() {
			Expect(list[0].DiscriminantValue).To(Equal(0.0))
		})

		It("Checking HitRank from search hit 1", func() {
			Expect(list[0].HitRank).To(Equal(uint8(1)))
		})

		It("Checking Hyperscore from search hit 1", func() {
			Expect(list[0].Hyperscore).To(Equal(21.783))
		})

		It("Checking index from search hit 1", func() {
			Expect(list[0].Index).To(Equal(uint32(1)))
		})

		It("Checking IsRejected from search hit 1", func() {
			Expect(list[0].IsRejected).To(Equal(uint8(0)))
		})

		It("Checking IsoMassD from search hit 1", func() {
			Expect(list[0].IsoMassD).To(Equal(0))
		})

		// It("Checking Massdiff from search hit 1", func() {
		// 	Expect(list[0].Massdiff).To(Equal(-0.0042))
		// 	fmt.Println(-0.0042)
		// })

		It("Checking MissedCleavages from search hit 1", func() {
			Expect(list[0].MissedCleavages).To(Equal(uint8(0)))
		})

		It("Checking NextAA from search hit 1", func() {
			Expect(list[0].NextAA).To(Equal("K"))
		})

		It("Checking Nextscore from search hit 1", func() {
			Expect(list[0].Nextscore).To(Equal(16.169))
		})

		It("Checking NumberMatchedIons from search hit 1", func() {
			Expect(list[0].NumberMatchedIons).To(Equal(uint16(11)))
		})

		It("Checking NumberTotalProteins from search hit 1", func() {
			Expect(list[0].NumberTotalProteins).To(Equal(uint16(1)))
		})

		It("Checking PrecursorExpMass from search hit 1", func() {
			Expect(list[0].PrecursorExpMass).To(Equal(0.0))
		})

		It("Checking PrecursorNeutralMass from search hit 1", func() {
			Expect(list[0].PrecursorNeutralMass).To(Equal(983.5470))
		})

		It("Checking PrevAA from search hit 1", func() {
			Expect(list[0].PrevAA).To(Equal("K"))
		})

		It("Checking Probability from search hit 1", func() {
			Expect(list[0].Probability).To(Equal(0.9986))
		})

		It("Checking Protein from search hit 1", func() {
			Expect(list[0].Protein).To(Equal("sp|P26583|HMGB2_HUMAN"))
		})

		It("Checking RetentionTime from search hit 1", func() {
			Expect(list[0].RetentionTime).To(Equal(1591.055))
		})

		It("Checking Scan from search hit 1", func() {
			Expect(list[0].Scan).To(Equal(1882))
		})

		It("Checking ModifiedPeptide from search hit 6568", func() {
			Expect(list[6568].ModifiedPeptide).To(Equal("RGLKPSCTIIPLM[147]K"))
		})

		It("Checking fixed Modification 1 from search hit 6568", func() {
			mod := list[6568].Modifications.Index["C#7#160.0307"]
			Expect(mod.MonoIsotopicMass).To(Equal(160.0307))
			Expect(mod.Position).To(Equal("7"))
			Expect(mod.MassDiff).To(Equal(57.0215))
			Expect(mod.AminoAcid).To(Equal("C"))
		})

		It("Checking variable Modification 2 from search hit 6568", func() {
			mod := list[6568].Modifications.Index["M#13#147.0354"]
			Expect(mod.MonoIsotopicMass).To(Equal(147.0354))
			Expect(mod.Position).To(Equal("13"))
			Expect(mod.MassDiff).To(Equal(15.9949))
			Expect(mod.AminoAcid).To(Equal("M"))
		})

		// It("Checking the expect score of the first peptide ID", func() {
		// 	Expect(list[0].Expectation).To(Equal("8.496e-03"))
		// })

		It("Checking last peptide sequence", func() {
			Expect(list[64405].Peptide).To(Equal("LAVEALSSLDGDLAGR"))
		})

	})

	Context("Testing protxml parsing", func() {

		var p id.ProtXML
		var groups id.GroupList
		var e error

		It("Accessing workspace", func() {
			e = os.Chdir("../../test/wrksp/")
			Expect(e).NotTo(HaveOccurred())
		})

		It("Reading interact.prot.xml", func() {
			p.Read("interact.prot.xml")
			Expect(e).NotTo(HaveOccurred())
			groups = append(groups, p.Groups...)
		})

		It("Checking the number of groups", func() {
			Expect(len(groups)).To(Equal(7926))
		})

		It("Checking index of group 2", func() {
			Expect(groups[1].GroupNumber).To(Equal(uint32(2)))
		})

		It("Checking the probability of group 2", func() {
			Expect(groups[1].Probability).To(Equal(1.0))
		})

		It("Checking the probability of protein 1 in group 2", func() {
			Expect(groups[1].Proteins[0].Probability).To(Equal(1.0))
		})

		It("Checking HasRazor of protein 1 in group 2", func() {
			Expect(groups[1].Proteins[0].HasRazor).To(Equal(false))
		})

		It("Checking the length of protein 1 in group 2", func() {
			Expect(groups[1].Proteins[0].Length).To(Equal("268"))
		})

		It("Checking the number of peptide ions for protein 1 in group 2", func() {
			Expect(len(groups[1].Proteins[0].PeptideIons)).To(Equal(3))
		})

		It("Checking percent coverage of protein 1 in group 2", func() {
			Expect(groups[1].Proteins[0].PercentCoverage).To(Equal(float32(6.300000190734863)))
		})

		It("Checking name of protein 1 in group 2", func() {
			Expect(groups[1].Proteins[0].ProteinName).To(Equal("sp|A0A0B4J2D5|GAL3B_HUMAN"))
		})

		It("Checking top peptide probability for protein 1 in group 2", func() {
			Expect(groups[1].Proteins[0].TopPepProb).To(Equal(float64(0.9989)))
		})

		It("Checking sequence of peptide 1 in protein 1 in group 2", func() {
			Expect(groups[1].Proteins[0].PeptideIons[0].PeptideSequence).To(Equal("EVVEAHVDQK"))
		})

		It("Checking charge of peptide 1 in protein 1 in group 2", func() {
			Expect(groups[1].Proteins[0].PeptideIons[0].Charge).To(Equal(uint8(2)))
		})

		It("Checking uniqueness of peptide 1 in protein 1 in group 2", func() {
			Expect(groups[1].Proteins[0].PeptideIons[0].IsUnique).To(Equal(true))
		})

		It("Checking ModifiedPeptide for peptide 1 in protein 1 in group 17", func() {
			Expect(groups[16].Proteins[0].PeptideIons[12].ModifiedPeptide).To(Equal("IAFIFNNLSQSNM[147]TQK"))
		})

	})
})
