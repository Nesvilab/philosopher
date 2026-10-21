package rep

import (
	"os"
	"testing"

	"github.com/Nesvilab/philosopher/lib/dat"
	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/sys"
)

// A decoy representative promoted to a target among its indistinguishable proteins is reported as that
// target, with the target's database record, and keeps its PSMs. The database is annotated into a fresh
// workspace, so the test runs inside one.
func TestAssembleProteinReport_PromotedProtein(t *testing.T) {

	// not t.TempDir: Restore leaves the database parts open, and Windows cannot remove an open file
	dir, e := os.MkdirTemp("", "philosopher")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	t.Chdir(dir)
	os.Mkdir(sys.MetaDir(), sys.FilePermission())

	fasta := ">rev_sp|P00001|X_HUMAN Protein X OS=Homo sapiens OX=9606 GN=XXX PE=1 SV=1\nKEDITPEPYLNO\n" +
		">tr|Q00002|Y_HUMAN Protein Y OS=Homo sapiens OX=9606 GN=YYY PE=1 SV=1\nMONLYHEREPEPKYY\n" +
		">sp|P00003|Z_HUMAN Protein Z OS=Homo sapiens OX=9606 GN=ZZZ PE=1 SV=1\nMONLYHEREPEPKZZZZZZZZZZZZ\n"
	if e := os.WriteFile("db.fas", []byte(fasta), sys.FilePermission()); e != nil {
		t.Fatal(e)
	}

	db := dat.New()
	db.ProcessDBAndSerialize("db.fas", "rev_", false)
	db.Serialize()

	p := id.ProtXML{
		DecoyTag: "rev_",
		Groups: id.GroupList{
			{
				Proteins: id.ProtIDList{
					{
						OriginalHeader:           "rev_sp|P00001|X_HUMAN Protein X OS=Homo sapiens OX=9606 GN=XXX PE=1 SV=1",
						ProteinName:              "rev_sp|P00001|X_HUMAN",
						Description:              "Protein X OS=Homo sapiens OX=9606 GN=XXX PE=1 SV=1",
						IndistinguishableProtein: []string{"tr|Q00002|Y_HUMAN", "sp|P00003|Z_HUMAN"},
						Length:                   12,
					},
				},
			},
		},
	}
	p.PromoteProteinIDs()

	evi := Evidence{PSM: PSMEvidenceList{{
		Spectrum:     "a.00001.00001.2",
		SpectrumFile: "a.pepXML",
		Peptide:      "ONLYHEREPEPK",
		Protein:      "sp|P00003|Z_HUMAN",
	}}}

	evi.AssembleProteinReport(p.Groups[0].Proteins, 0.5, "", "rev_")

	pe := evi.Proteins[0]
	if pe.PartHeader != "sp|P00003|Z_HUMAN" || pe.ProteinID != "P00003" || pe.GeneNames != "ZZZ" || pe.Length != 25 {
		t.Errorf("PartHeader, ProteinID, GeneNames, Length = %q, %q, %q, %d, want %q, %q, %q, 25",
			pe.PartHeader, pe.ProteinID, pe.GeneNames, pe.Length, "sp|P00003|Z_HUMAN", "P00003", "ZZZ")
	}
	if _, ok := pe.IndiProtein["sp|P00003|Z_HUMAN"]; ok {
		t.Errorf("sp|P00003|Z_HUMAN is listed as indistinguishable from itself")
	}

	evi.SyncPSMToProteins("rev_")

	if len(evi.Proteins) != 1 || len(evi.PSM) != 1 {
		t.Errorf("%d proteins and %d PSMs after the synchronization, want 1 and 1", len(evi.Proteins), len(evi.PSM))
	}
}
