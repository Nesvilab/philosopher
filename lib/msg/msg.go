package msg

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Custom call for custom cases
func Custom(e error, t string) {

	m := fmt.Sprintf("%s", e)

	callLogrus(m, t)

	return
}

// OverwrittingMeta call when trying to execute external binaries
func OverwrittingMeta(e error, t string) {

	m := fmt.Sprintf("A meta data folder was found and will not be overwritten. %s", e)

	callLogrus(m, t)

	return
}

// MarshalFile call for a failed Marshaling event
func MarshalFile(e error, t string) {

	m := fmt.Sprintf("Cannot marshal file. %s", e)

	callLogrus(m, t)

	return
}

// SerializeFile call for a failed serialization attempt
func SerializeFile(e error, t string) {

	m := fmt.Sprintf("Cannot serialize file. %s", e)

	callLogrus(m, t)

	return
}

// CopyingFile call when trying to copy files to another location
func CopyingFile(e error, t string) {

	m := fmt.Sprintf("Cannot copy or move file. %s", e)

	callLogrus(m, t)

	return
}

// CastFloatToString call when trying to cast a float number to string
func CastFloatToString(e error, t string) {

	m := fmt.Sprintf("Cannot cast float information to string. %s", e)

	callLogrus(m, t)

	return
}

// Plotter call for faled plotter instantiation
func Plotter(e error, t string) {

	m := fmt.Sprintf("Could not instantiate plotter. %s", e)

	callLogrus(m, t)

	return
}

// ReadFile call for file not found
func ReadFile(e error, t string) {

	m := fmt.Sprintf("Cannot read file. %s", e)

	callLogrus(m, t)

	return
}

// ReadingMzMLZlib call when trying to erad mzML zlibed spectra
func ReadingMzMLZlib(e error, t string) {

	m := fmt.Sprintf("Error trying to read mzML zlib data. %s", e)

	callLogrus(m, t)

	return
}

// WriteFile call for failed file writing event
func WriteFile(e error, t string) {

	m := fmt.Sprintf("Cannot write file. %s", e)

	callLogrus(m, t)

	return
}

// WriteToFile call for failed file writing event
func WriteToFile(e error, t string) {

	m := fmt.Sprintf("Cannot write to file. %s", e)

	callLogrus(m, t)

	return
}

// DeployAsset call for failed asset deployment
func DeployAsset(e error, t string) {

	m := fmt.Sprintf("Cannot deploy asset. %s", e)

	callLogrus(m, t)

	return
}

// DecodeMsgPck call for failed msgpack decoding
func DecodeMsgPck(e error, t string) {

	m := fmt.Sprintf("Cannot decode packed binary. %s", e)

	callLogrus(m, t)

	return
}

// InputNotFound call empty parameters list
func InputNotFound(e error, t string) {

	m := fmt.Sprintf("Missing input file. %s", e)

	callLogrus(m, t)

	return
}

// NoParametersFound call empty parameters list
func NoParametersFound(e error, t string) {

	m := fmt.Sprintf("Missing input parameters. %s", e)

	callLogrus(m, t)

	return
}

// DatabaseNotFound call for a missing database file
func DatabaseNotFound(e error, t string) {

	m := fmt.Sprintf("Database not found. %s", e)

	callLogrus(m, t)

	return
}

// NoSpectraFound call empty Spectra structs
func NoSpectraFound(e error, t string) {

	m := fmt.Sprintf("No Spectra was found in data set. %s", e)

	callLogrus(m, t)

	return
}

// NoPSMFound call empty PSM structs
func NoPSMFound(e error, t string) {

	m := fmt.Sprintf("No PSM was found in data set. %s", e)

	callLogrus(m, t)

	return
}

// NoProteinFound call empty Protein structs
func NoProteinFound(e error, t string) {

	m := fmt.Sprintf("No Protein was found in data set. %s", e)

	callLogrus(m, t)

	return
}

// Comet call when running the Comet search engine
func Comet(e error, t string) {

	m := fmt.Sprintf("Missing parameter file or data file for analysis. %s", e)

	callLogrus(m, t)

	return
}

// UnsupportedDistribution call for error trying to determine OS distribution
func UnsupportedDistribution(e error, t string) {

	m := fmt.Sprintf("Cannot determine OS distribtion for binary version deployment. %s", e)

	callLogrus(m, t)

	return
}

// ExecutingBinary call when trying to execute external binaries
func ExecutingBinary(e error, t string) {

	m := fmt.Sprintf("Cannot execute program. %s", e)

	callLogrus(m, t)

	return
}

// WorkspaceNotFound call when trying to locate a workspace
func WorkspaceNotFound(e error, t string) {

	m := fmt.Sprintf("Workspace not found. %s", e)

	callLogrus(m, t)

	return
}

// GettingLocalDir call when trying to pinpoint current directory
func GettingLocalDir(e error, t string) {

	m := fmt.Sprintf("Cannot verify local directory path. %s", e)

	callLogrus(m, t)

	return
}

// CreatingMetaDirectory call when trying to create a meta directory
func CreatingMetaDirectory(e error, t string) {

	m := fmt.Sprintf("Cannot create meta directory; check folder permissions. %s", e)

	callLogrus(m, t)

	return
}

// LocatingTemDirecotry call when trying to locate the Temp directory
func LocatingTemDirecotry(e error, t string) {

	m := fmt.Sprintf("Cannot locate temporary directory. %s", e)

	callLogrus(m, t)

	return
}

// LocatingMetaDirecotry call when trying to locate the Meta directory
func LocatingMetaDirecotry(e error, t string) {

	m := fmt.Sprintf("Cannot locate meta directory. %s", e)

	callLogrus(m, t)

	return
}

// ArchivingMetaDirecotry call when trying to archive the Meta directory
func ArchivingMetaDirecotry(e error, t string) {

	m := fmt.Sprintf("Cannot archive meta directory, chekc your zip libraries. %s", e)

	callLogrus(m, t)

	return
}

// DeletingMetaDirecotry call when trying to delete the Meta directory
func DeletingMetaDirecotry(e error, t string) {

	m := fmt.Sprintf("Cannot delete meta directory, check your permissions. %s", e)

	callLogrus(m, t)

	return
}

// ParsingFASTA call when trying parse a protein FASTA database
func ParsingFASTA(e error, t string) {

	m := fmt.Sprintf("Cannot parse the FASTA file, check for formatting errors or malformed headers. %s", e)

	callLogrus(m, t)

	return
}

// Done call when a process is ready
func Done() {

	m := fmt.Sprintf("Done")

	callLogrus(m, "info")

	return
}

// Executing declares the command or program and the version
func Executing(s, v string) {

	m := fmt.Sprintf("Executing %s %s", s, v)

	callLogrus(m, "info")

	return
}

// callLogrus returns the appropriate response for each erro type
func callLogrus(m, t string) {

	switch t {
	case "trace":
		logrus.Trace(m)
	case "debug":
		logrus.Debug(m)
	case "info":
		logrus.Info(m)
	case "warning":
		logrus.Warning(m)
	case "error":
		logrus.Error(m)
	case "fatal":
		logrus.Fatal(m)
	default:
		logrus.Error(m)
	}

	return
}
