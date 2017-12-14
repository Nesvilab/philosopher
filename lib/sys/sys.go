package sys

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// GetHome returns the user home directory name
func GetHome() (string, error) {

	var home string

	if runtime.GOOS == Windows() {
		home = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
	} else if runtime.GOOS == Linux() {
		home = os.Getenv("HOME")
	} else {
		return "", errors.New("Cannot define your operating system")
	}

	return home, nil
}

// GetTemp retirves the temporary directory name
func GetTemp() (string, error) {

	var tmp string

	if runtime.GOOS == Windows() {
		tmp = os.Getenv("Temp")
	} else if runtime.GOOS == Linux() {
		tmp = "/tmp"
	} else {
		return "", errors.New("Cannot define your operating system")
	}

	return tmp, nil
}

// GetLinuxFlavor returns the Linux flavor by looking into the lsb_release
func GetLinuxFlavor() (string, error) {

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

	return flavor, nil
}

// CopyFile emulates a system copy function. The function needs
// the full qualified names for both origin and destination
func CopyFile(from, to string) error {

	// Open original file
	originalFile, err := os.Open(from)
	if err != nil {
		log.Fatal(err)
	}
	defer originalFile.Close()

	// Create new file
	newFile, err := os.Create(to)
	if err != nil {
		log.Fatal(err)
	}
	defer newFile.Close()

	// Copy the bytes to destination from source
	_, err = io.Copy(newFile, originalFile)
	if err != nil {
		return errors.New("Error copying file")
	}

	// Commit the file contents
	err = newFile.Sync()
	if err != nil {
		return err
	}

	return nil
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
