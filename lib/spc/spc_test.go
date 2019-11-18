package spc_test

import (
	"os"
	"philosopher/lib/id"
	"testing"
)

func TestSpCPepXML(t *testing.T) {

	var list id.PepIDList
	var p id.PepXML

	os.Chdir("../../test/wrksp/")

	p.Read("interact.pep.xml")
	list = append(list, p.PeptideIdentification...)

	if len(list) != 64406 {
		t.Errorf("Spectra number is incorrect, got %d, want %d", 0, 64406)
	}

	if list[0].Peptide != "KGPGRPTGSK" {
		t.Errorf("Peptide sequence is incorrect, got %s, want %s", list[0].Peptide, "KGPGRPTGSK")
	}

	if list[0].Spectrum != "b1906_293T_proteinID_01A_QE3_122212.01882.01882.3" {
		t.Errorf("Spectrum is incorrect, got %s, want %s", list[0].Spectrum, "b1906_293T_proteinID_01A_QE3_122212.01882.01882.3")
	}

	if list[0].AssumedCharge != uint8(3) {
		t.Errorf("AssumedCharge is incorrect, got %d, want %d", list[0].AssumedCharge, uint8(3))
	}

	if list[0].CalcNeutralPepMass != 983.5512 {
		t.Errorf("CalcNeutralPepMass is incorrect, got %f, want %f", list[0].CalcNeutralPepMass, 983.5512)
	}

	if list[0].DiscriminantValue != 0.0 {
		t.Errorf("Discriminant is incorrect, got %f, want %f", list[0].DiscriminantValue, 0.0)
	}

	if list[0].HitRank != uint8(1) {
		t.Errorf("Hit Rank is incorrect, got %d, want %d", list[0].HitRank, uint8(1))
	}

	if list[0].Hyperscore != 21.783 {
		t.Errorf("Hyperscore is incorrect, got %f, want %f", list[0].Hyperscore, 21.783)
	}

	if list[0].Index != uint32(1) {
		t.Errorf("Index is incorrect, got %d, want %d", list[0].Index, uint32(1))
	}

	if list[0].IsRejected != uint8(0) {
		t.Errorf("IsRejected is incorrect, got %d, want %d", list[0].IsRejected, uint8(0))
	}

	if list[0].IsoMassD != 0 {
		t.Errorf("IsoMassD is incorrect, got %d, want %d", list[0].IsoMassD, 0)
	}

	if list[0].MissedCleavages != uint8(0) {
		t.Errorf("Missed Cleavages is incorrect, got %d, want %d", list[0].MissedCleavages, uint8(0))
	}

	if list[0].NextAA != "K" {
		t.Errorf("NextAA is incorrect, got %s, want %s", list[0].NextAA, "K")
	}

	if list[0].Nextscore != 16.169 {
		t.Errorf("Nextscore is incorrect, got %f, want %f", list[0].Nextscore, 16.169)
	}

	if list[0].NumberMatchedIons != uint16(11) {
		t.Errorf("Number Matched Ions is incorrect, got %d, want %d", list[0].NumberMatchedIons, uint16(11))
	}

	if list[0].NumberTotalProteins != 1 {
		t.Errorf("Number Total Proteins is incorrect, got %d, want %d", list[0].NumberTotalProteins, 1)
	}

	if list[0].PrecursorExpMass != 0.0 {
		t.Errorf("PrecursorExpMass is incorrect, got %f, want %f", list[0].PrecursorExpMass, 0.0)
	}

	if list[0].PrecursorNeutralMass != 983.5470 {
		t.Errorf("PrecursorNeutralMass is incorrect, got %f, want %f", list[0].PrecursorNeutralMass, 983.5470)
	}

	if list[0].PrevAA != "K" {
		t.Errorf("PrevAA is incorrect, got %s, want %s", list[0].PrevAA, "K")
	}

	if list[0].Probability != 0.9986 {
		t.Errorf("Probability is incorrect, got %f, want %f", list[0].Probability, 0.9986)
	}

	if list[0].Protein != "sp|P26583|HMGB2_HUMAN" {
		t.Errorf("Protein is incorrect, got %s, want %s", list[0].Protein, "sp|P26583|HMGB2_HUMAN")
	}

	if list[0].RetentionTime != 1591.055 {
		t.Errorf("RetentionTime is incorrect, got %f, want %f", list[0].RetentionTime, 1591.055)
	}

	if list[0].Scan != 1882 {
		t.Errorf("Scan is incorrect, got %d, want %d", list[0].Scan, 1882)
	}

	mod1 := list[6568].Modifications.Index["C#7#160.0307"]
	if mod1.MonoIsotopicMass != 160.0307 {
		t.Errorf("MonoIsotopic Mass 1 is incorrect, got %f, want %f", mod1.MonoIsotopicMass, 160.0307)
	}

	mod2 := list[6568].Modifications.Index["M#13#147.0354"]
	if mod2.MonoIsotopicMass != 147.0354 {
		t.Errorf("MonoIsotopic Mass 2 is incorrect, got %f, want %f", mod2.MonoIsotopicMass, 147.0354)
	}

	if list[0].Peptide != "KGPGRPTGSK" {
		t.Errorf("Peptide is incorrect, got %s, want %s", list[0].Peptide, "KGPGRPTGSK")
	}

	if list[0].Expectation != 8.496e-03 {
		t.Errorf("Expectation is incorrect, got %f, want %f", list[0].Expectation, 8.496e-03)
	}

}

// 		// It("Checking Massdiff from search hit 1", func() {
// 		// 	Expect(list[0].Massdiff).To(Equal(-0.0042))
// 		// 	fmt.Println(-0.0042)
// 		// })

// 		It("Checking fixed Modification 1 from search hit 6568", func() {
// 			mod := list[6568].Modifications.Index["C#7#160.0307"]
// 			Expect(mod.MonoIsotopicMass).To(Equal(160.0307))
// 			Expect(mod.Position).To(Equal("7"))
// 			Expect(mod.MassDiff).To(Equal(57.0215))
// 			Expect(mod.AminoAcid).To(Equal("C"))
// 		})

// 		It("Checking variable Modification 2 from search hit 6568", func() {
// 			mod := list[6568].Modifications.Index["M#13#147.0354"]
// 			Expect(mod.MonoIsotopicMass).To(Equal(147.0354))
// 			Expect(mod.Position).To(Equal("13"))
// 			Expect(mod.MassDiff).To(Equal(15.9949))
// 			Expect(mod.AminoAcid).To(Equal("M"))
// 		})

func TestSpCProtXML(t *testing.T) {

	var p id.ProtXML
	var groups id.GroupList

	os.Chdir("../../test/wrksp/")

	p.Read("interact.prot.xml")
	groups = append(groups, p.Groups...)

}

// 	Context("Testing protxml parsing", func() {

// 		var p id.ProtXML
// 		var groups id.GroupList
// 		var e error

// 		It("Accessing workspace", func() {
// 			e = os.Chdir("../../test/wrksp/")
// 			Expect(e).NotTo(HaveOccurred())
// 		})

// 		It("Reading interact.prot.xml", func() {
// 			p.Read("interact.prot.xml")
// 			Expect(e).NotTo(HaveOccurred())
// 			groups = append(groups, p.Groups...)
// 		})

// 		It("Checking the number of groups", func() {
// 			Expect(len(groups)).To(Equal(7926))
// 		})

// 		It("Checking index of group 2", func() {
// 			Expect(groups[1].GroupNumber).To(Equal(uint32(2)))
// 		})

// 		It("Checking the probability of group 2", func() {
// 			Expect(groups[1].Probability).To(Equal(1.0))
// 		})

// 		It("Checking the probability of protein 1 in group 2", func() {
// 			Expect(groups[1].Proteins[0].Probability).To(Equal(1.0))
// 		})

// 		It("Checking HasRazor of protein 1 in group 2", func() {
// 			Expect(groups[1].Proteins[0].HasRazor).To(Equal(false))
// 		})

// 		It("Checking the length of protein 1 in group 2", func() {
// 			Expect(groups[1].Proteins[0].Length).To(Equal("268"))
// 		})

// 		It("Checking the number of peptide ions for protein 1 in group 2", func() {
// 			Expect(len(groups[1].Proteins[0].PeptideIons)).To(Equal(3))
// 		})

// 		It("Checking percent coverage of protein 1 in group 2", func() {
// 			Expect(groups[1].Proteins[0].PercentCoverage).To(Equal(float32(6.300000190734863)))
// 		})

// 		It("Checking name of protein 1 in group 2", func() {
// 			Expect(groups[1].Proteins[0].ProteinName).To(Equal("sp|A0A0B4J2D5|GAL3B_HUMAN"))
// 		})

// 		It("Checking top peptide probability for protein 1 in group 2", func() {
// 			Expect(groups[1].Proteins[0].TopPepProb).To(Equal(float64(0.9989)))
// 		})

// 		It("Checking sequence of peptide 1 in protein 1 in group 2", func() {
// 			Expect(groups[1].Proteins[0].PeptideIons[0].PeptideSequence).To(Equal("EVVEAHVDQK"))
// 		})

// 		It("Checking charge of peptide 1 in protein 1 in group 2", func() {
// 			Expect(groups[1].Proteins[0].PeptideIons[0].Charge).To(Equal(uint8(2)))
// 		})

// 		It("Checking uniqueness of peptide 1 in protein 1 in group 2", func() {
// 			Expect(groups[1].Proteins[0].PeptideIons[0].IsUnique).To(Equal(true))
// 		})

// 		It("Checking ModifiedPeptide for peptide 1 in protein 1 in group 17", func() {
// 			Expect(groups[16].Proteins[0].PeptideIons[12].ModifiedPeptide).To(Equal("IAFIFNNLSQSNM[147]TQK"))
// 		})

// 	})
// })
