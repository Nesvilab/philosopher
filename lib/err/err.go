package err

import "fmt"

// INFO represents the error information class
const INFO string = "info"

// WARN represents the error information class
const WARN string = "warning"

// FATA represents the error information class
const FATA string = "fatal"

const (

	// WorkspaceNotFound is used when a command is executed outside the workspace environment
	WorkspaceNotFound string = "workspace not found. Run 'philosopher workspace --init' to create one or execute the command inside a workspace directory"

	// CannotOverwriteMeta indicates when a new meta folder tries to overwrite an existing one
	CannotOverwriteMeta string = "existing workspace detected, will not overwrite"

	// CorruptedOrEmptyMeta indicates a problem with the meta file
	CorruptedOrEmptyMeta string = "corrupted workspace"

	// CannotFindData indicates an empty or corrupted data structure
	CannotFindData string = "empty data structure"

	// CannotFindMetaDirectory indicates a missing Meta folder
	CannotFindMetaDirectory string = "cannot find the meta data"

	// CannotFindPSMData indicates a missing PSM and peptide data
	CannotFindPSMData string = "cannot find peptide identification data"

	// CannotDeleteMetaDirectory indicates a missing Meta folder
	CannotDeleteMetaDirectory string = "cannot remove the meta data"

	// CannotParseDataBase indicates a problem when reading the database file
	CannotParseDataBase string = "cannot read or parse the database file"

	// CannotZipMetaDirectory indicates a problem trying to zip the metadata
	CannotZipMetaDirectory string = "there was a problem zipping the meta data"

	// CannotSerializeData is used when there is a problem trying to serialize processed data into disk
	CannotSerializeData string = "cannot serialize data structures"

	// CannotRestoreGob is used when a function tries to restore a serielized Gob file
	CannotRestoreGob string = "cannot restore serialized data structures"

	// NoValidationFound indicates missing validation from pepXML files
	NoValidationFound string = "no peptide validation found"

	// CannotOpenFile is used when the given file cant be open, regardless if the file exsts or not
	CannotOpenFile string = "cannot open file"

	// CannotCreateOutputFile is used when th function si trying to create a new text based file as output on disk
	CannotCreateOutputFile string = "cannot create output file, check your writting permissions and disk space"

	// CannotCreateDirectory is used when a directory faisl to be created
	CannotCreateDirectory string = "cannot create directory"

	// CannotStatLocalDirectory is used when a the system faisl to stat the local direcotry via Getwd
	CannotStatLocalDirectory string = "cannot stat current directory"

	// CannotFindUniProtAnnotation is used when the function is trying to fetch annotation data from uniprot
	CannotFindUniProtAnnotation string = "Cannot find annotation information from UniProt, check your connection"

	// CannotDownloadUniProtAnnotation is used when the function is trying to fetch annotation data from uniprot
	CannotDownloadUniProtAnnotation string = "Cannot download annotation file from UniProt, check your writting permissions and connection"

	// CannotInstantiateStruct indcates when a constructor fails
	CannotInstantiateStruct string = "struct cannot be instantiated"

	// CannotParseXML is used when the marshal functions cant read or understand the XML file
	CannotParseXML string = "unable to parse XML file"

	// CannotParseFastaFile is used when there is an error trying to read and parse a FASTA file
	CannotParseFastaFile string = "unable to parse the FASTA file"

	// CannotDeployCrapDB is used when there is a problem tryingto deploy the crap database
	CannotDeployCrapDB string = "unable to deploy crap database, file may be corrupted"

	// CannotIdentifyDatabaseType is used when the database commands needs to guess what is the source of the database file
	CannotIdentifyDatabaseType string = "cannot identify the database type. Your database file contains unformatted headers"

	// CannotConvertFloatToString is used when a float value fails to be converted to string
	CannotConvertFloatToString string = "unable to cast float value to string"

	// NoPSMFound indicates when an PSm list is empty
	NoPSMFound string = "no PSMs were found"

	// CannotExtractAsset is used when there is a problem trying to deploy a binary from bindata
	CannotExtractAsset string = "cannot deploy binary file"

	// CannotExecuteBinary is used when there is a problem trying to run a third-party tool or binary file
	CannotExecuteBinary string = "cannot run program"

	// CannotCopyFile is used when there is a problem trying copy files
	CannotCopyFile string = "cannot copy or move file"

	// UnsupportedDistribution is used when there is an incompatibility between deployed binaries and the OS
	UnsupportedDistribution string = "unsupported distribution for the program"

	// CannotRunComet is used when the comet search engine returns an error
	CannotRunComet string = "cannot run Comet search"

	// CannotRunMSFragger is used when the msfragger search engine returns an error
	CannotRunMSFragger string = "cannot run MSFragger search"

	// CannotRunProgram is used when a given binary fails to run
	CannotRunProgram string = "cannot execute program"

	// UnknownMultiplex is used when the TMT label setting does not exists
	UnknownMultiplex string = "unknown multiplex setting"

	// MethodNotImplemented is used when a incomplete or empty method is called
	MethodNotImplemented string = "method is non-existent"

	// CannotGetLinuxFlavor is used when the LInux falvor cannot be determined
	CannotGetLinuxFlavor string = "cannot determine Linux distribution"

	// CannotDeployAsset is used when an asset is not able to be deploied on disk
	CannotDeployAsset string = "cannot deploy asset on disk"
)

// Error is the base struct for all errors
type Error struct {
	Message  string
	Argument string
	Type     string
	Class    string
}

func (e *Error) Error() string {
	e.Message = e.Type

	if len(e.Argument) > 0 {
		return fmt.Sprintf("%s: %s", e.Message, e.Argument)
	}

	return fmt.Sprintf("%s", e.Message)
}
