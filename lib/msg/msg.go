package msg

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Custom call for custom cases
func Custom(e error, t string) {

	m := fmt.Sprintf("%s", e)

	callLogrus(m, t)

}

// OverwrittingMeta call when trying to execute external binaries
func OverwrittingMeta(e error, t string) {

	m := fmt.Sprintf("A meta data folder was found and will not be overwritten. %s", e)

	callLogrus(m, t)

}

// MarshalFile call for a failed Marshaling event
func MarshalFile(e error, t string) {

	m := fmt.Sprintf("Cannot marshal file. %s", e)

	callLogrus(m, t)
}

// SerializeFile call for a failed serialization attempt
func SerializeFile(e error, t string) {

	m := fmt.Sprintf("Cannot serialize file. %s", e)

	callLogrus(m, t)
}

// CopyingFile call when trying to copy files to another location
func CopyingFile(e error, t string) {

	m := fmt.Sprintf("Cannot copy or move file. %s", e)

	callLogrus(m, t)

}

// CastFloatToString call when trying to cast a float number to string
func CastFloatToString(e error, t string) {

	m := fmt.Sprintf("Cannot cast float information to string. %s", e)

	callLogrus(m, t)

}

// Plotter call for faled plotter instantiation
func Plotter(e error, t string) {

	m := fmt.Sprintf("Could not instantiate plotter. %s", e)

	callLogrus(m, t)

}

// ReadFile call for file not found
func ReadFile(e error, t string) {

	m := fmt.Sprintf("Cannot read file. %s", e)

	callLogrus(m, t)

}

// ReadingMzMLZlib call when trying to erad mzML zlibed spectra
func ReadingMzMLZlib(e error, t string) {

	m := fmt.Sprintf("Error trying to read mzML zlib data. %s", e)

	callLogrus(m, t)

}

// WriteFile call for failed file writing event
func WriteFile(e error, t string) {

	m := fmt.Sprintf("Cannot write file. %s", e)

	callLogrus(m, t)

}

// WriteToFile call for failed file writing event
func WriteToFile(e error, t string) {

	m := fmt.Sprintf("Cannot write to file. %s", e)

	callLogrus(m, t)

}

// DeployAsset call for failed asset deployment
func DeployAsset(e error, t string) {

	m := fmt.Sprintf("Cannot deploy asset. %s", e)

	callLogrus(m, t)

}

// DecodeMsgPck call for failed msgpack decoding
func DecodeMsgPck(e error, t string) {

	m := fmt.Sprintf("Cannot decode packed binary. %s", e)

	callLogrus(m, t)

}

// InputNotFound call empty parameters list
func InputNotFound(e error, t string) {

	m := fmt.Sprintf("Missing input file. %s", e)

	callLogrus(m, t)

}

// NoParametersFound call empty parameters list
func NoParametersFound(e error, t string) {

	m := fmt.Sprintf("Missing input parameters. %s", e)

	callLogrus(m, t)

}

// DatabaseNotFound call for a missing database file
func DatabaseNotFound(e error, t string) {

	m := fmt.Sprintf("Database not found. %s", e)

	callLogrus(m, t)

}

// NoSpectraFound call empty Spectra structs
func NoSpectraFound(e error, t string) {

	m := fmt.Sprintf("No Spectra was found in data set. %s", e)

	callLogrus(m, t)

}

// NoPSMFound call empty PSM structs
func NoPSMFound(e error, t string) {

	m := fmt.Sprintf("No PSM was found in data set. %s", e)

	callLogrus(m, t)

}

// QuantifyingData call when trying to do quantification on a data set with problems
func QuantifyingData(e error, t string) {

	m := fmt.Sprintf("Cannot quantify data set. %s", e)

	callLogrus(m, t)
}

// NoProteinFound call empty Protein structs
func NoProteinFound(e error, t string) {

	m := fmt.Sprintf("No Protein was found in data set. %s", e)

	callLogrus(m, t)

}

// Comet call when running the Comet search engine
func Comet(e error, t string) {

	m := fmt.Sprintf("Missing parameter file or data file for analysis. %s", e)

	callLogrus(m, t)

}

// UnsupportedDistribution call for error trying to determine OS distribution
func UnsupportedDistribution(e error, t string) {

	m := fmt.Sprintf("Cannot determine OS distribtion for binary version deployment. %s", e)

	callLogrus(m, t)

}

// ExecutingBinary call when trying to execute external binaries
func ExecutingBinary(e error, t string) {

	m := fmt.Sprintf("Cannot execute program. %s", e)

	callLogrus(m, t)

}

// WorkspaceNotFound call when trying to locate a workspace
func WorkspaceNotFound(e error, t string) {

	m := fmt.Sprintf("Workspace not found. %s", e)

	callLogrus(m, t)

}

// GettingLocalDir call when trying to pinpoint current directory
func GettingLocalDir(e error, t string) {

	m := fmt.Sprintf("Cannot verify local directory path. %s", e)

	callLogrus(m, t)

}

// CreatingMetaDirectory call when trying to create a meta directory
func CreatingMetaDirectory(e error, t string) {

	m := fmt.Sprintf("Cannot create meta directory; check folder permissions. %s", e)

	callLogrus(m, t)

}

// LocatingTemDirecotry call when trying to locate the Temp directory
func LocatingTemDirecotry(e error, t string) {

	m := fmt.Sprintf("Cannot locate temporary directory. %s", e)

	callLogrus(m, t)

}

// LocatingMetaDirecotry call when trying to locate the Meta directory
func LocatingMetaDirecotry(e error, t string) {

	m := fmt.Sprintf("Cannot locate meta directory. %s", e)

	callLogrus(m, t)

}

// ArchivingMetaDirecotry call when trying to archive the Meta directory
func ArchivingMetaDirecotry(e error, t string) {

	m := fmt.Sprintf("Cannot archive meta directory, chekc your zip libraries. %s", e)

	callLogrus(m, t)

}

// DeletingMetaDirecotry call when trying to delete the Meta directory
func DeletingMetaDirecotry(e error, t string) {

	m := fmt.Sprintf("Cannot delete meta directory, check your permissions. %s", e)

	callLogrus(m, t)

}

// ParsingFASTA call when trying parse a protein FASTA database
func ParsingFASTA(e error, t string) {

	m := fmt.Sprintf("Cannot parse the FASTA file, check for formatting errors or malformed headers. %s", e)

	callLogrus(m, t)

}

// ParsingFASTAHeader call when trying parse a protein FASTA database
func ParsingFASTAHeader(e error, t string) {

	m := fmt.Sprintf("Malformed FASTA header. %s", e)

	callLogrus(m, t)

}

// Done call when a process is ready
func Done() {

	m := "Done"

	callLogrus(m, "info")

}

// Executing declares the command or program and the version
func Executing(s, v string) {

	m := fmt.Sprintf("Executing %s %s", s, v)

	callLogrus(m, "info")

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
	case "fatal":
		logrus.Error(m)
		panic(m)
	case "error":
		logrus.Error(m)
		os.Exit(1)
	default:
		logrus.Error(m)
		os.Exit(1)
	}

}
