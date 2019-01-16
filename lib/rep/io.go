package rep

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/prvst/philosopher/lib/err"
	"github.com/prvst/philosopher/lib/sys"
	"github.com/vmihailenco/msgpack"
)

// Serialize converts the whole structure to a gob file
func (e *Evidence) Serialize() *err.Error {

	b, er := msgpack.Marshal(&e)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvBin(), b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// SerializeGranular converts the whole structure into sevral small gob files
func (e *Evidence) SerializeGranular() *err.Error {

	// create EV PSM
	er := SerializeEVPSM(e)
	if er != nil {
		return er
	}

	// create EV Ion
	er = SerializeEVIon(e)
	if er != nil {
		return er
	}

	// create EV Peptides
	er = SerializeEVPeptides(e)
	if er != nil {
		return er
	}

	// create EV Ion
	er = SerializeEVProteins(e)
	if er != nil {
		return er
	}

	// create EV Mods
	er = SerializeEVMods(e)
	if er != nil {
		return er
	}

	// create EV Modifications
	er = SerializeEVModifications(e)
	if er != nil {
		return er
	}

	// create EV Combined
	er = SerializeEVCombined(e)
	if er != nil {
		return er
	}

	return nil
}

// SerializeEVPSM creates an ev serial with Evidence data
func SerializeEVPSM(e *Evidence) *err.Error {

	b, er := msgpack.Marshal(&e.PSM)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvPSMBin(), b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// SerializeEVIon creates an ev serial with Evidence data
func SerializeEVIon(e *Evidence) *err.Error {

	b, er := msgpack.Marshal(&e.Ions)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvIonBin(), b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// SerializeEVPeptides creates an ev serial with Evidence data
func SerializeEVPeptides(e *Evidence) *err.Error {

	b, er := msgpack.Marshal(&e.Peptides)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvPeptideBin(), b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// SerializeEVProteins creates an ev serial with Evidence data
func SerializeEVProteins(e *Evidence) *err.Error {

	b, er := msgpack.Marshal(&e.Proteins)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvProteinBin(), b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// SerializeEVMods creates an ev serial with Evidence data
func SerializeEVMods(e *Evidence) *err.Error {

	b, er := msgpack.Marshal(&e.Mods)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvModificationsBin(), b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// SerializeEVModifications creates an ev serial with Evidence data
func SerializeEVModifications(e *Evidence) *err.Error {

	b, er := msgpack.Marshal(&e.Modifications)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvModificationsEvBin(), b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// SerializeEVCombined creates an ev serial with Evidence data
func SerializeEVCombined(e *Evidence) *err.Error {

	b, er := msgpack.Marshal(&e.CombinedProtein)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvCombinedBin(), b, 0644)
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (e *Evidence) Restore() error {

	file, _ := os.Open(sys.EvBin())

	dec := msgpack.NewDecoder(file)
	err := dec.Decode(&e)
	if err != nil {
		return errors.New("Could not restore Philosopher result. Please check file path")
	}

	return nil
}

// RestoreGranular reads philosopher results files and restore the data sctructure
func (e *Evidence) RestoreGranular() *err.Error {

	// PSM
	err := RestoreEVPSM(e)
	if err != nil {
		return err
	}

	// Ion
	err = RestoreEVIon(e)
	if err != nil {
		return err
	}

	// Peptide
	err = RestoreEVPeptide(e)
	if err != nil {
		return err
	}

	// Protein
	err = RestoreEVProtein(e)
	if err != nil {
		return err
	}

	// Mods
	err = RestoreEVMods(e)
	if err != nil {
		return err
	}

	// Modifications
	err = RestoreEVModifications(e)
	if err != nil {
		return err
	}

	// Combined
	err = RestoreEVCombined(e)
	if err != nil {
		return err
	}

	return nil
}

// RestoreEVPSM restores Ev PSM data
func RestoreEVPSM(e *Evidence) *err.Error {
	f, _ := os.Open(sys.EvPSMBin())
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.PSM)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVIon restores Ev Ion data
func RestoreEVIon(e *Evidence) *err.Error {
	f, _ := os.Open(sys.EvIonBin())
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Ions)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVPeptide restores Ev Ion data
func RestoreEVPeptide(e *Evidence) *err.Error {
	f, _ := os.Open(sys.EvPeptideBin())
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Peptides)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVProtein restores Ev Protein data
func RestoreEVProtein(e *Evidence) *err.Error {
	f, _ := os.Open(sys.EvProteinBin())
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Proteins)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVMods restores Ev Mods data
func RestoreEVMods(e *Evidence) *err.Error {
	f, _ := os.Open(sys.EvModificationsBin())
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Mods)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVModifications restores Ev Mods data
func RestoreEVModifications(e *Evidence) *err.Error {
	f, _ := os.Open(sys.EvModificationsEvBin())
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Modifications)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVCombined restores Ev Mods data
func RestoreEVCombined(e *Evidence) *err.Error {
	f, _ := os.Open(sys.EvCombinedBin())
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.CombinedProtein)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreGranularWithPath reads philosopher results files and restore the data sctructure
func (e *Evidence) RestoreGranularWithPath(p string) *err.Error {

	// PSM
	err := RestoreEVPSMWithPath(e, p)
	if err != nil {
		return err
	}

	// Ion
	err = RestoreEVIonWithPath(e, p)
	if err != nil {
		return err
	}

	// Peptide
	err = RestoreEVPeptideWithPath(e, p)
	if err != nil {
		return err
	}

	// Protein
	err = RestoreEVProteinWithPath(e, p)
	if err != nil {
		return err
	}

	// Mods
	err = RestoreEVModsWithPath(e, p)
	if err != nil {
		return err
	}

	// Modifications
	err = RestoreEVModificationsWithPath(e, p)
	if err != nil {
		return err
	}

	// Combined
	err = RestoreEVCombinedWithPath(e, p)
	if err != nil {
		return err
	}

	return nil
}

// RestoreEVPSMWithPath restores Ev PSM data
func RestoreEVPSMWithPath(e *Evidence, p string) *err.Error {

	//path := sys.EvPSMBin()
	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvPSMBin())

	// if strings.Contains(p, string(filepath.Separator)) {
	// 	path = fmt.Sprintf("%s%s", p, sys.EvPSMBin())
	// } else {
	// 	path = fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvPSMBin())
	// }

	f, _ := os.Open(path)
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.PSM)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVIonWithPath restores Ev Ion data
func RestoreEVIonWithPath(e *Evidence, p string) *err.Error {

	// path := sys.EvIonBin()
	//
	// if strings.Contains(p, string(filepath.Separator)) {
	// 	path = fmt.Sprintf("%s%s", p, sys.EvIonBin())
	// } else {
	// 	path = fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvIonBin())
	// }

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvIonBin())

	f, _ := os.Open(path)
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Ions)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVPeptideWithPath restores Ev Ion data
func RestoreEVPeptideWithPath(e *Evidence, p string) *err.Error {

	// path := sys.EvPeptideBin()
	//
	// if strings.Contains(p, string(filepath.Separator)) {
	// 	path = fmt.Sprintf("%s%s", p, sys.EvPeptideBin())
	// } else {
	// 	path = fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvPeptideBin())
	// }

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvPeptideBin())

	f, _ := os.Open(path)
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Peptides)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVProteinWithPath restores Ev Protein data
func RestoreEVProteinWithPath(e *Evidence, p string) *err.Error {

	// path := sys.EvProteinBin()
	//
	// if strings.Contains(p, string(filepath.Separator)) {
	// 	path = fmt.Sprintf("%s%s", p, sys.EvProteinBin())
	// } else {
	// 	path = fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvProteinBin())
	// }

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvProteinBin())

	f, _ := os.Open(path)
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Proteins)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVModsWithPath restores Ev Mods data
func RestoreEVModsWithPath(e *Evidence, p string) *err.Error {

	// path := sys.EvModificationsBin()
	//
	// if strings.Contains(p, string(filepath.Separator)) {
	// 	path = fmt.Sprintf("%s%s", p, sys.EvModificationsBin())
	// } else {
	// 	path = fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvModificationsBin())
	// }

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvModificationsBin())

	f, _ := os.Open(path)
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Mods)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVModificationsWithPath restores Ev Mods data
func RestoreEVModificationsWithPath(e *Evidence, p string) *err.Error {

	// path := sys.EvModificationsEvBin()
	//
	// if strings.Contains(p, string(filepath.Separator)) {
	// 	path = fmt.Sprintf("%s%s", p, sys.EvModificationsEvBin())
	// } else {
	// 	path = fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvModificationsEvBin())
	// }

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvModificationsEvBin())

	f, _ := os.Open(path)
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.Modifications)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}

// RestoreEVCombinedWithPath restores Ev Mods data
func RestoreEVCombinedWithPath(e *Evidence, p string) *err.Error {

	// path := sys.EvCombinedBin()
	//
	// if strings.Contains(p, string(filepath.Separator)) {
	// 	path = fmt.Sprintf("%s%s", p, sys.EvCombinedBin())
	// } else {
	// 	path = fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvCombinedBin())
	// }

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvCombinedBin())

	f, _ := os.Open(path)
	d := msgpack.NewDecoder(f)
	er := d.Decode(&e.CombinedProtein)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: er.Error()}
	}
	return nil
}
