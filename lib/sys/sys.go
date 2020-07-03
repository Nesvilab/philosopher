package sys

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"philosopher/lib/msg"
)

// GetHome returns the user home directory name
func GetHome() string {

	var home string

	if runtime.GOOS == Windows() {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else if runtime.GOOS == Linux() {
		home = os.Getenv("HOME")
	} else {
		msg.Custom(errors.New("Cannot define your operating system"), "fatal")
	}

	return home
}

// GetTemp retirves the temporary directory name
func GetTemp() string {
	var tmp string

	if runtime.GOOS == Windows() {
		tmp = os.Getenv("Temp")
	} else if runtime.GOOS == Linux() {
		tmp = "/tmp"
	} else {
		msg.Custom(errors.New("Cannot define your operating system"), "fatal")
	}

	return tmp
}

// VerifyTemp allows the definition of a custom folder to be used for deplyments and file creations
func VerifyTemp(f string) {

	if _, err := os.Stat(f); os.IsNotExist(err) {
		msg.Custom(errors.New("Cannot find the custom temporary folder"), "fatal")
	}

	return
}

// GetLinuxFlavor returns the Linux flavor by looking into the lsb_release
func GetLinuxFlavor() string {

	var flavor string

	if runtime.GOOS == Linux() {

		cmd := exec.Command("lsb_release", "-a")
		output, _ := cmd.CombinedOutput()

		if strings.Contains(string(output), Ubuntu()) || strings.Contains(string(output), Mint()) || strings.Contains(string(output), Debian()) {
			flavor = Debian()
		} else if strings.Contains(string(output), Redhat()) {
			flavor = Redhat()
		} else if strings.Contains(string(output), Centos()) {
			flavor = Redhat()
		} else {
			flavor = Redhat()
		}
	} else {
		flavor = Windows()
	}

	return flavor
}

// CopyFile emulates a system copy function. The function needs
// the full qualified names for both origin and destination
func CopyFile(from, to string) {

	// Open original file
	originalFile, e := os.Open(from)
	if e != nil {
		msg.ReadFile(e, "fatal")
	}
	defer originalFile.Close()

	// Create new file
	newFile, e := os.Create(to)
	if e != nil {
		msg.WriteFile(e, "fatal")
	}
	defer newFile.Close()

	// Copy the bytes to destination from source
	_, e = io.Copy(newFile, originalFile)
	if e != nil {
		msg.CopyingFile(e, "fatal")
	}

	// Commit the file contents
	e = newFile.Sync()
	if e != nil {
		msg.Custom(e, "fatal")
	}

	return
}

// Meta file
func Meta() string {
	p := fmt.Sprintf("%s%smeta.bin", MetaDir(), string(filepath.Separator))
	return p
}

// RawBin file
func RawBin() string {
	p := fmt.Sprintf("%s%sraw.bin", MetaDir(), string(filepath.Separator))
	return p
}

// PepxmlBin file
func PepxmlBin() string {
	p := fmt.Sprintf("%s%spepxml.bin", MetaDir(), string(filepath.Separator))
	return p
}

// ProtxmlBin file
func ProtxmlBin() string {
	p := fmt.Sprintf("%s%sprotxml.bin", MetaDir(), string(filepath.Separator))
	return p
}

// PsmBin file
func PsmBin() string {
	p := fmt.Sprintf("%s%spsm.bin", MetaDir(), string(filepath.Separator))
	return p
}

// PepBin file
func PepBin() string {
	p := fmt.Sprintf("%s%spep.bin", MetaDir(), string(filepath.Separator))
	return p
}

// IonBin file
func IonBin() string {
	p := fmt.Sprintf("%s%sion.bin", MetaDir(), string(filepath.Separator))
	return p
}

// ProBin file
func ProBin() string {
	p := fmt.Sprintf("%s%spro.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvBin file
func EvBin() string {
	p := fmt.Sprintf("%s%sev.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvParameterBin file
func EvParameterBin() string {
	p := fmt.Sprintf("%s%sev.param.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvMetaBin file
func EvMetaBin() string {
	p := fmt.Sprintf("%s%sev.meta.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvPSMBin file
func EvPSMBin() string {
	p := fmt.Sprintf("%s%sev.psm.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvPeptideBin file
func EvPeptideBin() string {
	p := fmt.Sprintf("%s%sev.pep.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvProteinBin file
func EvProteinBin() string {
	p := fmt.Sprintf("%s%sev.pro.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvModificationsBin file
func EvModificationsBin() string {
	p := fmt.Sprintf("%s%sev.mod.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvModificationsEvBin file
func EvModificationsEvBin() string {
	p := fmt.Sprintf("%s%sev.mev.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvCombinedBin file
func EvCombinedBin() string {
	p := fmt.Sprintf("%s%sev.com.bin", MetaDir(), string(filepath.Separator))
	return p
}

// EvIonBin file
func EvIonBin() string {
	p := fmt.Sprintf("%s%sev.ion.bin", MetaDir(), string(filepath.Separator))
	return p
}

// DBBin file
func DBBin() string {
	p := fmt.Sprintf("%s%sdb.bin", MetaDir(), string(filepath.Separator))
	return p
}

// LFQBin file
func LFQBin() string {
	p := fmt.Sprintf("%s%slfq.bin", MetaDir(), string(filepath.Separator))
	return p
}

// IsoBin file
func IsoBin() string {
	p := fmt.Sprintf("%s%siso.bin", MetaDir(), string(filepath.Separator))
	return p
}

// MODBin file
func MODBin() string {
	p := fmt.Sprintf("%s%smod.bin", MetaDir(), string(filepath.Separator))
	return p
}

// MetaDir dir
func MetaDir() string {
	return ".meta"
}

// Linux OS
func Linux() string {
	return "linux"
}

// Windows OS
func Windows() string {
	return "windows"
}

// Darwin OS
func Darwin() string {
	return "darwin"
}

// Redhat OS
func Redhat() string {
	return "RedHat"
}

// Ubuntu OS
func Ubuntu() string {
	return "Ubuntu"
}

// Mint OS
func Mint() string {
	return "Mint"
}

// Debian OS
func Debian() string {
	return "Debian"
}

// Centos OS
func Centos() string {
	return "CentOS"
}

// Arch386 arch
func Arch386() string {
	return "386"
}

// FilePermission sets the default permission for every file written to disk
func FilePermission() os.FileMode {
	//return 0644
	return 0755
}
