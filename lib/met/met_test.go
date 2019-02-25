package met_test

import (
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/prvst/philosopher/lib/met"
	"github.com/prvst/philosopher/lib/sys"
)

var _ = Describe("Met", func() {

	Context("Testing the meta data structure", func() {

		var dir string
		var d Data
		var e error

		It("New", func() {
			dir, e = os.Getwd()
			Expect(e).NotTo(HaveOccurred())
			d = New(dir)
			Expect(e).NotTo(HaveOccurred())
		})

		It("UUID", func() {
			Expect(len(d.UUID)).NotTo(Equal(0))
		})

		It("OS", func() {
			Expect(d.OS).To(Equal(runtime.GOOS))
		})

		It("Arch", func() {
			Expect(d.Arch).To(Equal(runtime.GOARCH))
		})

		// It("Distro", func() {
		// 	Expect(d.Distro).To(Equal("Debian"))
		// })

		It("Home", func() {
			//home := fmt.Sprintf("%s%swrksp", dir, string(filepath.Separator))
			Expect(d.Home).To(Equal(dir))
		})

		It("Project Name", func() {
			Expect(d.ProjectName).To(Equal(string(filepath.Base(dir))))
		})

		It("Meta File", func() {
			Expect(d.MetaFile).To(Equal(d.Home + string(filepath.Separator) + sys.Meta()))
		})

		It("Meta Directory", func() {
			Expect(d.MetaDir).To(Equal(d.Home + string(filepath.Separator) + sys.MetaDir()))
		})

		It("Database", func() {
			Expect(d.DB).To(Equal(d.Home + string(filepath.Separator) + sys.DBBin()))
		})

		It("Temp Directory", func() {
			temp, e := sys.GetTemp()
			Expect(e).NotTo(HaveOccurred())
			temp += string(filepath.Separator) + d.UUID
			Expect(d.Temp).To(Equal(temp))
		})

	})

})
