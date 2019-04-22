package rep

import (
	"fmt"
	"io/ioutil"
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

	// create EV Parameters
	er := SerializeEVParameters(e)
	if er != nil {
		return er
	}

	// create EV PSM
	er = SerializeEVPSM(e)
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

// SerializeEVParameters creates an ev serial with Parameter data
func SerializeEVParameters(e *Evidence) *err.Error {

	b, er := msgpack.Marshal(&e.Parameters)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvParameterBin(), b, sys.FilePermission())
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// SerializeEVPSM creates an ev serial with Evidence data
func SerializeEVPSM(e *Evidence) *err.Error {

	b, er := msgpack.Marshal(&e.PSM)
	if er != nil {
		return &err.Error{Type: err.CannotCreateOutputFile, Class: err.FATA, Argument: er.Error()}
	}

	er = ioutil.WriteFile(sys.EvPSMBin(), b, sys.FilePermission())
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

	er = ioutil.WriteFile(sys.EvIonBin(), b, sys.FilePermission())
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

	er = ioutil.WriteFile(sys.EvPeptideBin(), b, sys.FilePermission())
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

	er = ioutil.WriteFile(sys.EvProteinBin(), b, sys.FilePermission())
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

	er = ioutil.WriteFile(sys.EvModificationsBin(), b, sys.FilePermission())
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

	er = ioutil.WriteFile(sys.EvModificationsEvBin(), b, sys.FilePermission())
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

	er = ioutil.WriteFile(sys.EvCombinedBin(), b, sys.FilePermission())
	if er != nil {
		return &err.Error{Type: err.CannotSerializeData, Class: err.FATA, Argument: er.Error()}
	}

	return nil
}

// Restore reads philosopher results files and restore the data sctructure
func (e *Evidence) Restore() *err.Error {

	b, er := ioutil.ReadFile(sys.EvBin())
	if er != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: ": database data may be corrupted"}
	}

	er = msgpack.Unmarshal(b, &e)
	if er != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: ": database data may be corrupted"}
	}

	return nil
}

// RestoreGranular reads philosopher results files and restore the data sctructure
func (e *Evidence) RestoreGranular() *err.Error {

	// Parameters
	err := RestoreEVParameters(e)
	if err != nil {
		return err
	}

	// PSM
	err = RestoreEVPSM(e)
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

// RestoreEVParameters restores Ev PSM data
func RestoreEVParameters(d *Evidence) *err.Error {

	b, e := ioutil.ReadFile(sys.EvParameterBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Parameters)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}
	return nil
}

// RestoreEVPSM restores Ev PSM data
func RestoreEVPSM(d *Evidence) *err.Error {

	b, e := ioutil.ReadFile(sys.EvPSMBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.PSM)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}
	return nil
}

// RestoreEVIon restores Ev Ion data
func RestoreEVIon(d *Evidence) *err.Error {

	b, e := ioutil.ReadFile(sys.EvIonBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Ions)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVPeptide restores Ev Ion data
func RestoreEVPeptide(d *Evidence) *err.Error {

	b, e := ioutil.ReadFile(sys.EvPeptideBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Peptides)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVProtein restores Ev Protein data
func RestoreEVProtein(d *Evidence) *err.Error {

	b, e := ioutil.ReadFile(sys.EvProteinBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Proteins)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVMods restores Ev Mods data
func RestoreEVMods(d *Evidence) *err.Error {

	b, e := ioutil.ReadFile(sys.EvModificationsBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Mods)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVModifications restores Ev Mods data
func RestoreEVModifications(d *Evidence) *err.Error {

	b, e := ioutil.ReadFile(sys.EvModificationsEvBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Modifications)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVCombined restores Ev Mods data
func RestoreEVCombined(d *Evidence) *err.Error {

	b, e := ioutil.ReadFile(sys.EvCombinedBin())
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.CombinedProtein)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreGranularWithPath reads philosopher results files and restore the data sctructure
func (e *Evidence) RestoreGranularWithPath(p string) *err.Error {

	// Parameters
	err := RestoreEVParametersWithPath(e, p)
	if err != nil {
		return err
	}

	// PSM
	err = RestoreEVPSMWithPath(e, p)
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

// RestoreEVParametersWithPath restores Ev PSM data
func RestoreEVParametersWithPath(d *Evidence, p string) *err.Error {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvParameterBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Parameters)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVPSMWithPath restores Ev PSM data
func RestoreEVPSMWithPath(d *Evidence, p string) *err.Error {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvPSMBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.PSM)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVIonWithPath restores Ev Ion data
func RestoreEVIonWithPath(d *Evidence, p string) *err.Error {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvIonBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Ions)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVPeptideWithPath restores Ev Ion data
func RestoreEVPeptideWithPath(d *Evidence, p string) *err.Error {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvPeptideBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Peptides)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVProteinWithPath restores Ev Protein data
func RestoreEVProteinWithPath(d *Evidence, p string) *err.Error {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvProteinBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Proteins)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVModsWithPath restores Ev Mods data
func RestoreEVModsWithPath(d *Evidence, p string) *err.Error {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvModificationsBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Mods)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVModificationsWithPath restores Ev Mods data
func RestoreEVModificationsWithPath(d *Evidence, p string) *err.Error {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvModificationsEvBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.Modifications)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}

// RestoreEVCombinedWithPath restores Ev Mods data
func RestoreEVCombinedWithPath(d *Evidence, p string) *err.Error {

	path := fmt.Sprintf("%s%s%s", p, string(filepath.Separator), sys.EvCombinedBin())

	b, e := ioutil.ReadFile(path)
	if e != nil {
		return &err.Error{Type: err.CannotOpenFile, Class: err.FATA, Argument: e.Error()}
	}

	e = msgpack.Unmarshal(b, &d.CombinedProtein)
	if e != nil {
		return &err.Error{Type: err.CannotRestoreGob, Class: err.FATA, Argument: e.Error()}
	}

	return nil
}
