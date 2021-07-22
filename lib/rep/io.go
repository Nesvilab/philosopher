package rep

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"philosopher/lib/sys"

	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

// SerializeGranular converts the whole structure into sevral small gob files
func (evi *Evidence) SerializeGranular() {

	// create PSM Bin
	SerializePSM(&evi.PSM)

	// create Ion Bin
	SerializeIon(&evi.Ions)

	// create Peptides Bin
	SerializePeptides(&evi.Peptides)

	// create Protein Bin
	SerializeProteins(&evi.Proteins)
}

// SerializePSM creates an ev serial with Evidence data
func SerializePSM(evi *PSMEvidenceList) {

	b, e := msgpack.Marshal(&evi)
	if e != nil {
		logrus.Trace("Cannot marshal PSM data:", e)
	}

	e = ioutil.WriteFile(sys.PSMBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize PSM data:", e)
	}

}

// SerializeIon creates an ev serial with Evidence data
func SerializeIon(evi *IonEvidenceList) {

	b, e := msgpack.Marshal(&evi)
	if e != nil {
		logrus.Trace("Cannot marshal Ions data:", e)
	}

	e = ioutil.WriteFile(sys.IonBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize Ions data:", e)
	}
}

// SerializePeptides creates an ev serial with Evidence data
func SerializePeptides(evi *PeptideEvidenceList) {

	b, e := msgpack.Marshal(&evi)
	if e != nil {
		logrus.Trace("Cannot marshal Peptides data:", e)
	}

	e = ioutil.WriteFile(sys.PepBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize Peptides data:", e)
	}

}

// SerializeProteins creates an ev serial with Evidence data
func SerializeProteins(evi *ProteinEvidenceList) {

	b, e := msgpack.Marshal(&evi)
	if e != nil {
		logrus.Trace("Cannot marshal Proteins data:", e)
	}

	e = ioutil.WriteFile(sys.ProBin(), b, sys.FilePermission())
	if e != nil {
		logrus.Trace("Cannot serialize Proteins data:", e)
	}

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

	b, e := ioutil.ReadFile(sys.PSMBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

}

// RestoreIon restores Ion data
func RestoreIon(evi *IonEvidenceList) {

	b, e := ioutil.ReadFile(sys.IonBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

}

// RestorePeptide restores Peptide data
func RestorePeptide(evi *PeptideEvidenceList) {

	b, e := ioutil.ReadFile(sys.PepBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

}

// RestoreProtein restores Protein data
func RestoreProtein(evi *ProteinEvidenceList) {

	b, e := ioutil.ReadFile(sys.ProBin())
	if e != nil {
		logrus.Fatal("Cannot read file:", e)
	}

	e = msgpack.Unmarshal(b, &evi)
	if e != nil {
		logrus.Fatal("Cannot unmarshal file:", e)
	}

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
