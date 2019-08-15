package err

import "github.com/sirupsen/logrus"

// WarnCustom : WARN call for custom cases
func WarnCustom(e error) {
	logrus.Warn(e)
}

// OverwrittingMeta : WARN call when trying to execute external binaries
func OverwrittingMeta() {
	logrus.Warn("A meta data folder was found and will not be overwritten")
}

// MarshalFile : TRACE call for a failed Marshaling event
func MarshalFile(e error) {
	logrus.Trace("Cannot read file: ", e)
}

// SerializeFile : TRACE call for a failed serialization attempt
func SerializeFile(e error) {
	logrus.Trace("Cannot serialize file: ", e)
}

// CopyingFile : TRACE call when trying to copy files to another location
func CopyingFile(e error) {
	logrus.Trace("Cannot copy or mvoe file: ", e)
}

// CastFloatToString : TRACE call when trying to cast a float number to string
func CastFloatToString(e error) {
	logrus.Trace("Cannot cast float information to string")
}

// ErrorCustom : ERROR call for custom cases
func ErrorCustom(e error) {
	logrus.Error(e)
}

// Plotter : WARN call for faled plotter instantiation
func Plotter(e error) {
	logrus.Fatal("Could not instantiate plotter: ", e)
}

// ReadFile : FATAL call for file not found
func ReadFile(e error) {
	logrus.Fatal("Cannot read file: ", e)
}

// ReadingMzMLZlib : FATAL call when trying to erad mzML zlibed spectra
func ReadingMzMLZlib(e error) {
	logrus.Fatal("Error trying to read mzML zlib data:", e)
}

// FatalCustom : FATAL call for custom cases
func FatalCustom(e error) {
	logrus.Fatal(e)
}

// WriteFile : FATAL call for failed file writing event
func WriteFile(e error) {
	logrus.Fatal("Cannot write file: ", e)
}

// WriteToFile : FATAL call for failed file writing into a file event
func WriteToFile(e error) {
	logrus.Fatal("Cannot write file: ", e)
}

// DeployAsset : FATAL call for failed asset deployment
func DeployAsset(e error) {
	logrus.Fatal("Cannot deploy asset: ", e)
}

// DecodeMsgPck : FATAL call for failed msgpack decoding
func DecodeMsgPck(e error) {
	logrus.Fatal("Cannot decode packed binary: ", e)
}

// NoParametersFound : FATAL call empty parameters list
func NoParametersFound(e error) {
	logrus.Fatal("Missing input parameters: ", e)
}

// DatabaseNotFound : FATAL call for a missing database file
func DatabaseNotFound(e error) {
	logrus.Fatal("Database not found: ", e)
}

// NoSpectraFound : FATAL call empty Spectra structs
func NoSpectraFound() {
	logrus.Fatal("No Spectra was found in data set")
}

// NoPSMFound : FATAL call empty PSM structs
func NoPSMFound() {
	logrus.Fatal("No PSM was found in data set")
}

// NoProteinFound : FATAL call empty Protein structs
func NoProteinFound() {
	logrus.Fatal("No Protein was found in data set")
}

// Comet : FATAL call when running the Comet search engine
func Comet() {
	logrus.Fatal("Missing parameter file or data file for analysis")
}

// UnsupportedDistribution : FATAL call for error trying to determine OS distribution
func UnsupportedDistribution() {
	logrus.Fatal("Cannot determine OS distribtion for binary version deployment")
}

// ExecutingBinary : FATAL call when trying to execute external binaries
func ExecutingBinary(e error) {
	logrus.Fatal("Cannot execute program: ", e)
}

// WorkspaceNotFound : FATAL call when trying to locate a workspace
func WorkspaceNotFound(e error) {
	logrus.Fatal("Workspace not found: ", e)
}

// GettingLocalDir : FATAL call when trying to pinpoint current directory
func GettingLocalDir(e error) {
	logrus.Fatal("Cannot verify local directory path: ", e)
}

// CreatingMetaDirectory : FATAL call when trying to create a meta directory
func CreatingMetaDirectory(e error) {
	logrus.Fatal("Cannot create meta directory; check folder permissions: ", e)
}

// LocatingTemDirecotry : FATAL call when trying to locate the Temp directory
func LocatingTemDirecotry(e error) {
	logrus.Fatal("Cannot locate temporary directory: ", e)
}

// LocatingMetaDirecotry : FATAL call when trying to locate the Meta directory
func LocatingMetaDirecotry() {
	logrus.Fatal("Cannot locate meta directory")
}

// ArchivingMetaDirecotry : FATAL call when trying to archive the Meta directory
func ArchivingMetaDirecotry(e error) {
	logrus.Fatal("Cannot archive meta directory, chekc your zip libraries: ", e)
}

// DeletingMetaDirecotry : FATAL call when trying to delete the Meta directory
func DeletingMetaDirecotry(e error) {
	logrus.Fatal("Cannot delete meta directory, check your permissions: ", e)
}

// ParsingFASTA : FATAL call when trying parse a protein FASTA database
func ParsingFASTA() {
	logrus.Fatal("Cannot parse the FASTA file, check for formatting errors or malformed headers")
}
