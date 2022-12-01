package spc_test

import (
	. "philosopher/lib/spc"
	"philosopher/lib/tes"
	"testing"

	_ "github.com/rogpeppe/go-charset/data"
)

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
// 			Expect(le0
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
// }

func TestPepXML_Parse(t *testing.T) {

	tes.SetupTestEnv()

	type fields struct {
		Name                 string
		MsmsPipelineAnalysis MsmsPipelineAnalysis
	}
	type args struct {
		f string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Testing pepXML parsing",
			args: args{f: "interact.pep.xml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PepXML{
				Name:                 tt.fields.Name,
				MsmsPipelineAnalysis: tt.fields.MsmsPipelineAnalysis,
			}

			p.Parse(tt.args.f)

			if len(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery) != 12615 {
				t.Errorf("Spectra number is incorrect, got %d, want %d", len(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery), 12615)
			}

			if string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].Peptide) != "SYPPEFINR" {
				t.Errorf("Peptide sequence is incorrect, got %s, want %s", string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].Peptide), "SYPPEFINR")
			}

			if string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].Spectrum) != "z04397_tc-o238g-setB_MS3.00663.00663.3" {
				t.Errorf("Spectrum is incorrect, got %s, want %s", string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].Spectrum), "z04397_tc-o238g-setB_MS3.00663.00663.3")
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].AssumedCharge != uint8(3) {
				t.Errorf("AssumedCharge is incorrect, got %d, want %d", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].AssumedCharge, uint8(3))
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].CalcNeutralPepMass != 1350.713400 {
				t.Errorf("CalcNeutralPepMass is incorrect, got %f, want %f", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].CalcNeutralPepMass, 1350.713400)
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].HitRank != uint8(1) {
				t.Errorf("Hit Rank is incorrect, got %d, want %d", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].HitRank, uint8(1))
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].Index != 28 {
				t.Errorf("Index is incorrect, got %d, want %d", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].Index, 28)
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].IsRejected != uint8(0) {
				t.Errorf("IsRejected is incorrect, got %d, want %d", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].IsRejected, uint8(0))
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].MissedCleavages != uint8(0) {
				t.Errorf("Missed Cleavages is incorrect, got %d, want %d", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].MissedCleavages, uint8(0))
			}

			if string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].NextAA) != "I" {
				t.Errorf("NextAA is incorrect, got %s, want %s", string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].NextAA), "I")
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].MatchedIons != uint16(8) {
				t.Errorf("Number Matched Ions is incorrect, got %d, want %d", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].MatchedIons, uint16(8))
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].TotalProteins != 1 {
				t.Errorf("Number Total Proteins is incorrect, got %d, want %d", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].TotalProteins, 1)
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].PrecursorNeutralMass != 1350.712400 {
				t.Errorf("PrecursorNeutralMass is incorrect, got %f, want %f", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].PrecursorNeutralMass, 1350.712400)
			}

			if string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].PrevAA) != "R" {
				t.Errorf("PrevAA is incorrect, got %s, want %s", string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].PrevAA), "R")
			}

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].AnalysisResult[0].PeptideProphetResult.Probability != 0.899159 {
				t.Errorf("Probability is incorrect, got %f, want %f", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].AnalysisResult[0].PeptideProphetResult.Probability, 0.899159)
			}

			if string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].Protein) != "sp|P33416|HSP78_YEAST" {
				t.Errorf("Protein is incorrect, got %s, want %s", string(p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].SearchResult.SearchHit[0].Protein), "sp|P33416|HSP78_YEAST")
			}

			// if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].RetentionTimeSec != 336.303328 {
			// 	t.Errorf("RetentionTime is incorrect, got %f, want %f", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].RetentionTimeSec, 336.303328)
			// }

			if p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].StartScan != 663 {
				t.Errorf("Scan is incorrect, got %d, want %d", p.MsmsPipelineAnalysis.MsmsRunSummary.SpectrumQuery[0].StartScan, 663)
			}

			// mod1 := list[6568].Modifications.Index["C#7#160.0307"]
			// if mod1.MonoIsotopicMass != 160.0307 {
			// 	t.Errorf("MonoIsotopic Mass 1 is incorrect, got %f, want %f", mod1.MonoIsotopicMass, 160.0307)
			// }

			// mod2 := list[6568].Modifications.Index["M#13#147.0354"]
			// if mod2.MonoIsotopicMass != 147.0354 {
			// 	t.Errorf("MonoIsotopic Mass 2 is incorrect, got %f, want %f", mod2.MonoIsotopicMass, 147.0354)
			// }
		})
	}
}

func TestProtXML_Parse(t *testing.T) {

	tes.SetupTestEnv()

	type fields struct {
		Name           string
		ProteinSummary ProteinSummary
	}
	type args struct {
		f string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "testing protXML parsing",
			args: args{"interact.prot.xml"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProtXML{
				Name:           tt.fields.Name,
				ProteinSummary: tt.fields.ProteinSummary,
			}

			p.Parse(tt.args.f)

			if len(p.ProteinSummary.ProteinGroup) != 2331 {
				t.Errorf("Number of protein groups is incorrect, got %d, want %d", len(p.ProteinSummary.ProteinGroup), 2331)
			}

			if string(p.ProteinSummary.ProteinGroup[0].Protein[0].ProteinName) != "contam_sp|O77727|K1C15_SHEEP" {
				t.Errorf("Protein group 1 name is incorrect, got %s, want %s", p.ProteinSummary.ProteinGroup[0].Protein[0].ProteinName, "contam_sp|O77727|K1C15_SHEEP")
			}

			if p.ProteinSummary.ProteinGroup[0].Protein[0].TotalNumberPeptides != 3 {
				t.Errorf("Total peptides for protein group 1 is incorrect, got %d, want %d", p.ProteinSummary.ProteinGroup[0].Protein[0].TotalNumberPeptides, 3)
			}

			if p.ProteinSummary.ProteinGroup[5].Protein[0].TotalNumberPeptides != 5 {
				t.Errorf("Total peptides for protein group 6 is incorrect, got %d, want %d", p.ProteinSummary.ProteinGroup[5].Protein[0].TotalNumberPeptides, 5)
			}

			// if string(p.ProteinSummary.ProteinGroup[5].Protein[0].ProteinName) != "sp|A0FGR8|ESYT2_HUMAN" {
			// 	t.Errorf("Protein 1 name in protein group 6 is incorrect, got %s, want %s", fmt.Sprint(p.ProteinSummary.ProteinGroup[5].Protein[0].TotalNumberPeptides), "sp|A0FGR8|ESYT2_HUMAN")
			// }

			if string(p.ProteinSummary.ProteinGroup[5].Protein[0].Peptide[0].PeptideSequence) != "ALNEINQFYQK" {
				t.Errorf("Peptide sequence 1 in protein 1, group 6 is incorrect, got %s, want %s", string(p.ProteinSummary.ProteinGroup[5].Protein[0].Peptide[0].PeptideSequence), "ALNEINQFYQK")
			}

			if p.ProteinSummary.ProteinGroup[5].Protein[0].Peptide[0].Charge != 2 {
				t.Errorf("Charge of peptide 1, protein 1, group 6 is incorrect, got %d, want %d", p.ProteinSummary.ProteinGroup[5].Protein[0].Peptide[0].Charge, 2)
			}

		})
	}

}
