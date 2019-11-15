package sys_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "philosopher/lib/sys"
)

var _ = Describe("Sys", func() {

	Context("meta locators and directories getters", func() {

		It("Meta", func() {
			p := fmt.Sprintf("%s%smeta.bin", MetaDir(), string(filepath.Separator))
			Expect(Meta()).To(Equal(p))
		})

		It("RawBin", func() {
			p := fmt.Sprintf("%s%sraw.bin", MetaDir(), string(filepath.Separator))
			Expect(RawBin()).To(Equal(p))
		})

		It("PepXML", func() {
			p := fmt.Sprintf("%s%spepxml.bin", MetaDir(), string(filepath.Separator))
			Expect(PepxmlBin()).To(Equal(p))
		})

		It("ProtXML", func() {
			p := fmt.Sprintf("%s%sprotxml.bin", MetaDir(), string(filepath.Separator))
			Expect(ProtxmlBin()).To(Equal(p))
		})

		It("PSMBin", func() {
			p := fmt.Sprintf("%s%spsm.bin", MetaDir(), string(filepath.Separator))
			Expect(PsmBin()).To(Equal(p))
		})

		It("PepBin", func() {
			p := fmt.Sprintf("%s%spep.bin", MetaDir(), string(filepath.Separator))
			Expect(PepBin()).To(Equal(p))
		})

		It("IonBin", func() {
			p := fmt.Sprintf("%s%sion.bin", MetaDir(), string(filepath.Separator))
			Expect(IonBin()).To(Equal(p))
		})

		It("ProBin", func() {
			p := fmt.Sprintf("%s%spro.bin", MetaDir(), string(filepath.Separator))
			Expect(ProBin()).To(Equal(p))
		})

		It("EvBin", func() {
			p := fmt.Sprintf("%s%sev.bin", MetaDir(), string(filepath.Separator))
			Expect(EvBin()).To(Equal(p))
		})

		It("EvBin", func() {
			p := fmt.Sprintf("%s%sev.bin", MetaDir(), string(filepath.Separator))
			Expect(EvBin()).To(Equal(p))
		})

		It("EvMetaBin", func() {
			p := fmt.Sprintf("%s%sev.meta.bin", MetaDir(), string(filepath.Separator))
			Expect(EvMetaBin()).To(Equal(p))
		})

		It("EvPSMBin", func() {
			p := fmt.Sprintf("%s%sev.psm.bin", MetaDir(), string(filepath.Separator))
			Expect(EvPSMBin()).To(Equal(p))
		})

		It("EvPeptideBin", func() {
			p := fmt.Sprintf("%s%sev.pep.bin", MetaDir(), string(filepath.Separator))
			Expect(EvPeptideBin()).To(Equal(p))
		})

		It("EvProteinBin", func() {
			p := fmt.Sprintf("%s%sev.pro.bin", MetaDir(), string(filepath.Separator))
			Expect(EvProteinBin()).To(Equal(p))
		})

		It("EvModificationsBin", func() {
			p := fmt.Sprintf("%s%sev.mod.bin", MetaDir(), string(filepath.Separator))
			Expect(EvModificationsBin()).To(Equal(p))
		})

		It("EvModificationsEvBin", func() {
			p := fmt.Sprintf("%s%sev.mev.bin", MetaDir(), string(filepath.Separator))
			Expect(EvModificationsEvBin()).To(Equal(p))
		})

		It("EvCombinedBin", func() {
			p := fmt.Sprintf("%s%sev.com.bin", MetaDir(), string(filepath.Separator))
			Expect(EvCombinedBin()).To(Equal(p))
		})

		It("EvIonBin", func() {
			p := fmt.Sprintf("%s%sev.ion.bin", MetaDir(), string(filepath.Separator))
			Expect(EvIonBin()).To(Equal(p))
		})

		It("DBBin", func() {
			p := fmt.Sprintf("%s%sdb.bin", MetaDir(), string(filepath.Separator))
			Expect(DBBin()).To(Equal(p))
		})

		It("MODBin", func() {
			p := fmt.Sprintf("%s%smod.bin", MetaDir(), string(filepath.Separator))
			Expect(MODBin()).To(Equal(p))
		})

	})

	Context("system and meta file names", func() {

		It("MetaDir", func() {
			Expect(MetaDir()).To(Equal(".meta"))
		})

		It("Linux", func() {
			Expect(Linux()).To(Equal("linux"))
		})

		It("Windows", func() {
			Expect(Windows()).To(Equal("windows"))
		})

		It("Darwin", func() {
			Expect(Darwin()).To(Equal("darwin"))
		})

		It("Redhat", func() {
			Expect(Redhat()).To(Equal("RedHat"))
		})

		It("Ubuntu", func() {
			Expect(Ubuntu()).To(Equal("Ubuntu"))
		})

		It("Mint", func() {
			Expect(Mint()).To(Equal("Mint"))
		})

		It("Debian", func() {
			Expect(Debian()).To(Equal("Debian"))
		})

		It("Centos", func() {
			Expect(Centos()).To(Equal("CentOS"))
		})

		It("Arch386", func() {
			Expect(Arch386()).To(Equal("386"))
		})

		It("FilePermission", func() {
			Expect(FilePermission()).To(Equal(os.FileMode(0755)))
		})

	})

})
