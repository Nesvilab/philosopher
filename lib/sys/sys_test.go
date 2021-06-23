package sys_test

import (
	"fmt"
	"path/filepath"
	"philosopher/lib/sys"
	"testing"
)

func TestSysMeta(t *testing.T) {

	p := fmt.Sprintf("%s%smeta.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.Meta() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.Meta())
	}

	p = fmt.Sprintf("%s%sraw.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.RawBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.RawBin())
	}

	p = fmt.Sprintf("%s%spepxml.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.PepxmlBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.PepxmlBin())
	}

	p = fmt.Sprintf("%s%sprotxml.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.ProtxmlBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.ProtxmlBin())
	}

	p = fmt.Sprintf("%s%spsm.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.PsmBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.PsmBin())
	}

	p = fmt.Sprintf("%s%spep.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.PepBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.PepBin())
	}

	p = fmt.Sprintf("%s%sion.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.IonBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.IonBin())
	}

	p = fmt.Sprintf("%s%spro.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.ProBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.ProBin())
	}

	p = fmt.Sprintf("%s%sev.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.EvBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.EvBin())
	}

	p = fmt.Sprintf("%s%sev.meta.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.EvMetaBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.EvMetaBin())
	}

	p = fmt.Sprintf("%s%sev.pep.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.EvPeptideBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.EvPeptideBin())
	}

	p = fmt.Sprintf("%s%sev.pro.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.EvProteinBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.EvProteinBin())
	}

	// p = fmt.Sprintf("%s%sev.mod.bin", sys.MetaDir(), string(filepath.Separator))
	// if p != sys.EvModificationsBin() {
	// 	t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.EvModificationsBin())
	// }

	// p = fmt.Sprintf("%s%sev.mev.bin", sys.MetaDir(), string(filepath.Separator))
	// if p != sys.EvModificationsEvBin() {
	// 	t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.EvModificationsEvBin())
	// }

	// p = fmt.Sprintf("%s%sev.com.bin", sys.MetaDir(), string(filepath.Separator))
	// if p != sys.EvCombinedBin() {
	// 	t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.EvCombinedBin())
	// }

	p = fmt.Sprintf("%s%sev.ion.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.EvIonBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.EvIonBin())
	}

	p = fmt.Sprintf("%s%sdb.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.DBBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.DBBin())
	}

	p = fmt.Sprintf("%s%smod.bin", sys.MetaDir(), string(filepath.Separator))
	if p != sys.MODBin() {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", p, sys.MODBin())
	}
}

func TestSysNames(t *testing.T) {

	if sys.MetaDir() != ".meta" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", ".meta", sys.MetaDir())
	}
	if sys.Linux() != "linux" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", "linux", sys.Linux())
	}
	if sys.Windows() != "windows" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", ".windows", sys.Windows())
	}
	if sys.Darwin() != "darwin" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", "darwin", sys.Darwin())
	}
	if sys.Redhat() != "RedHat" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", "redhat", sys.Redhat())
	}
	if sys.Ubuntu() != "Ubuntu" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", "ubuntu", sys.Ubuntu())
	}
	if sys.Mint() != "Mint" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", "mint", sys.Mint())
	}
	if sys.Debian() != "Debian" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", "debian", sys.Debian())
	}
	if sys.Centos() != "CentOS" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", "centos", sys.Centos())
	}
	if sys.Arch386() != "386" {
		t.Errorf("Meta path or name is incorrect, got %s, want %s", "386", sys.Arch386())
	}
	if sys.FilePermission() != 0755 {
		t.Errorf("Meta path or name is incorrect, got %d, want %s", 0755, sys.FilePermission())
	}

}
