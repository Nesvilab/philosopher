package rep

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"philosopher/lib/sys"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

// SerializeGranular converts the whole structure into sevral small gob files
func (evi *Evidence) SerializeGranular() {
	wg := sync.WaitGroup{}
	wg.Add(4)
	// create PSM Bin
	go func() { defer wg.Done(); SerializePSM(&evi.PSM) }()
	// create Ion Bin
	go func() { defer wg.Done(); SerializeIon(&evi.Ions) }()
	// create Peptides Bin
	go func() { defer wg.Done(); SerializePeptides(&evi.Peptides) }()
	// create Protein Bin
	go func() { defer wg.Done(); SerializeProteins(&evi.Proteins) }()
	wg.Wait()
}

// SerializePSM creates an ev serial with Evidence data
func SerializePSM(evi *PSMEvidenceList) {
	sys.Serialize(evi, sys.PSMBin())
}

// SerializeIon creates an ev serial with Evidence data
func SerializeIon(evi *IonEvidenceList) {
	sys.Serialize(evi, sys.IonBin())
}

// SerializePeptides creates an ev serial with Evidence data
func SerializePeptides(evi *PeptideEvidenceList) {
	sys.Serialize(evi, sys.PepBin())
}

// SerializeProteins creates an ev serial with Evidence data
func SerializeProteins(evi *ProteinEvidenceList) {
	sys.Serialize(evi, sys.ProBin())
}

// RestoreGranular reads philosopher results files and restore the data sctructure
func (evi *Evidence) RestoreGranular() {

	// PSM
	RestorePSM(&evi.PSM)

	// Ion
	RestoreIon(&evi.Ions)

	// Peptide
	RestorePeptide(&evi.Peptides)

	// Protein
	RestoreProtein(&evi.Proteins)
}

// RestorePSM restores PSM data
func RestorePSM(evi *PSMEvidenceList) {
	sys.Restore(evi, sys.PSMBin(), false)
}

// RestoreIon restores Ion data
func RestoreIon(evi *IonEvidenceList) {
	sys.Restore(evi, sys.IonBin(), false)
}

// RestorePeptide restores Peptide data
func RestorePeptide(evi *PeptideEvidenceList) {
	sys.Restore(evi, sys.PepBin(), false)
}

// RestoreProtein restores Protein data
func RestoreProtein(evi *ProteinEvidenceList) {
	sys.Restore(evi, sys.ProBin(), false)
}

// RestoreGranularWithPath reads philosopher results files and restore the data sctructure
func (evi *Evidence) RestoreGranularWithPath(p string) {

	// PSM
	RestorePSMWithPath(&evi.PSM, p)

	// Ion
	RestoreIonWithPath(&evi.Ions, p)

	// Peptide
	RestorePeptideWithPath(&evi.Peptides, p)

	// Protein
	RestoreProteinWithPath(&evi.Proteins, p)
}

// RestorePSMWithPath restores PSM data
func RestorePSMWithPath(evi *PSMEvidenceList, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.PSMBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

}

// RestoreIonWithPath restores Ion data
func RestoreIonWithPath(evi *IonEvidenceList, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.IonBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

}

// RestorePeptideWithPath restores Ion data
func RestorePeptideWithPath(evi *PeptideEvidenceList, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.PepBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

}

// RestoreProteinWithPath restores Protein data
func RestoreProteinWithPath(evi *ProteinEvidenceList, p string) {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.ProBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

}
