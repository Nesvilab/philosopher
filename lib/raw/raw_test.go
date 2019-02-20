package raw_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/prvst/philosopher/lib/raw"
)

var _ = Describe("Raw", func() {

	Context("Testing Raw file parsing", func() {

		var spec *raw.Data
		var ms1 raw.MS1
		var ms2 raw.MS2
		var spec1 raw.Ms1Scan
		var spec2 raw.Ms2Scan
		var e error

		It("Accessing workspace", func() {
			e = os.Chdir("../../test/wrksp/")
			Expect(e).NotTo(HaveOccurred())
		})

		It("Indexing mzML 01A", func() {
			var input []string
			input = append(input, "01_CPTAC_TMTS1-NCI7_Z_JHUZ_20170502_LUMOS.mzML")
			e = raw.IndexMz(input)
			Expect(e).NotTo(HaveOccurred())
		})

		It("Reading mzML 01A MS1", func() {
			spec, e = raw.RestoreFromFile(".", "01_CPTAC_TMTS1-NCI7_Z_JHUZ_20170502_LUMOS", "mzML")
			Expect(e).NotTo(HaveOccurred())
		})

		It("Read mzML 01A MS1 spectra", func() {
			ms1 = raw.GetMS1(spec)
			for _, i := range ms1.Ms1Scan {
				if i.Index == "0" && i.Scan == "1"{
					spec1 = i
					break
				}
			}
			Expect(len(ms1.Ms1Scan)).To(Equal(12405))
		})

		It("mzML 01A MS1 spectra Index", func() {
			Expect(spec1.Index).To(Equal("0"))
		})

		It("mzML 01A MS1 spectra Scan", func() {
			Expect(spec1.Scan).To(Equal("1"))
		})

		It("mzML 01A MS1 stream", func() {
			Expect(len(spec1.Spectrum)).To(Equal(865))
		})

		It("mzML 01A MS1 intensities", func() {
			Expect(spec1.Spectrum[0].Intensity).To(Equal(9104.91796875))
		})

		It("mzML 01A MS1 MZ", func() {
			Expect(spec1.Spectrum[0].Mz).To(Equal(350.1635437011719))
		})


		// It("mzML 01A MS1 spectra Scan Start Time", func() {
		// 	ms1 = raw.GetMS1(spec)
		// 	Expect(ms1.Ms1Scan[0].ScanStartTime).To(Equal(152.34552))
		// })

		It("Read mzML 01A MS2 spectra", func() {
			ms2 = raw.GetMS2(spec)
			for _, i := range ms2.Ms2Scan {
				if i.Index == "2" && i.Scan == "3"{
					spec2 = i
					break
				}
			}
			Expect(len(ms2.Ms2Scan)).To(Equal(41952))
		})

		It("mzML 01 MS2 spectra Index", func() {
			Expect(spec2.Index).To(Equal("2"))
		})

		It("mzML 01 MS2 spectra Scan", func() {
			Expect(spec2.Scan).To(Equal("3"))
		})

		It("mzML 01 MS2 stream", func() {
			Expect(len(spec2.Spectrum)).To(Equal(231))
		})

		It("mzML 01 MS1 intensities", func() {
			Expect(spec2.Spectrum[0].Intensity).To(Equal(371635.9375))
		})

		It("mzML 01 MS2 MZ", func() {
			Expect(spec2.Spectrum[0].Mz).To(Equal(110.07147216796875))
		})

		It("mzML 01 MS2 Parent Index", func() {
			Expect(spec2.Precursor.ParentIndex).To(Equal("1"))
		})

		It("mzML 01 MS2 Parent Scan", func() {
			Expect(spec2.Precursor.ParentScan).To(Equal("2"))
		})

		It("mzML 01 MS2 Charge State", func() {
			Expect(spec2.Precursor.ChargeState).To(Equal(2))
		})

		It("mzML 01 MS2 Parent Selected Ion", func() {
			Expect(spec2.Precursor.SelectedIon).To(Equal(391.201019287109))
		})

		It("mzML 01 MS2 Parent Target Ion", func() {
			Expect(spec2.Precursor.TargetIon).To(Equal(391.2))
		})

		It("mzML 01 MS2 Charge Peak Intensity", func() {
			Expect(spec2.Precursor.PeakIntensity).To(Equal(3.58558525e+06))
		})

		It("mzML 01 MS2 Parent Isolatio nWindow Lower Offset", func() {
			Expect(spec2.Precursor.IsolationWindowLowerOffset).To(Equal(0.34999999404))
		})

		It("mzML 01 MS2 Charge Isolation Window Upper Offset", func() {
			Expect(spec2.Precursor.IsolationWindowUpperOffset).To(Equal(0.34999999404))
		})

	})
})
