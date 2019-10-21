package mzn_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/nesvilab/philosopher/lib/mzn"
)

var _ = Describe("Mzn", func() {

	Context("Testing Raw file parsing", func() {

		var msd MsData
		var spec Spectrum
		var e error

		It("Accessing workspace", func() {
			e = os.Chdir("../../test/wrksp/")
			Expect(e).NotTo(HaveOccurred())
		})

		It("Reading mzML 01A MS1", func() {
			msd.Read("01_CPTAC_TMTS1-NCI7_Z_JHUZ_20170502_LUMOS.mzML", false, false, false)
			Expect(e).NotTo(HaveOccurred())
		})

		It("Read mzML 01A MS1 spectra", func() {
			for _, i := range msd.Spectra {
				if i.Index == "0" && i.Scan == "1" {
					spec = i
					spec.Decode()
					break
				}
			}
			Expect(len(msd.Spectra)).To(Equal(54357))
		})

		It("mzML 01A MS1 spectra Index", func() {
			Expect(spec.Index).To(Equal("0"))
		})

		It("mzML 01A MS1 spectra Scan", func() {
			Expect(spec.Scan).To(Equal("1"))
		})

		// It("mzML 01A MS1 stream", func() {
		// 	Expect(len(spec)).To(Equal(865))
		// })

		It("mzML 01A MS1 intensities", func() {
			Expect(spec.Intensity.DecodedStream[0]).To(Equal(9104.91796875))
		})

		It("mzML 01A MS1 MZ", func() {
			Expect(spec.Mz.DecodedStream[0]).To(Equal(350.1635437011719))
		})

		// It("mzML 01A MS1 spectra Scan Start Time", func() {
		// 	ms1 = raw.GetMS1(spec)
		// 	Expect(ms1.Ms1Scan[0].ScanStartTime).To(Equal(152.34552))
		// })

		It("Read mzML 01A MS2 spectra", func() {
			for _, i := range msd.Spectra {
				if i.Index == "2" && i.Scan == "3" {
					spec = i
					spec.Decode()
					break
				}
			}
			Expect(len(spec.Mz.DecodedStream)).To(Equal(231))
		})

		It("mzML 01 MS2 spectra Index", func() {
			Expect(spec.Index).To(Equal("2"))
		})

		It("mzML 01 MS2 spectra Scan", func() {
			Expect(spec.Scan).To(Equal("3"))
		})

		// It("mzML 01 MS2 stream", func() {
		// 	Expect(len(spec2.Spectrum)).To(Equal(231))
		// })

		It("mzML 01 MS1 intensities", func() {
			Expect(spec.Intensity.DecodedStream[0]).To(Equal(371635.9375))
		})

		It("mzML 01 MS2 MZ", func() {
			Expect(spec.Mz.DecodedStream[0]).To(Equal(110.07147216796875))
		})

		It("mzML 01 MS2 Parent Index", func() {
			Expect(spec.Precursor.ParentIndex).To(Equal("1"))
		})

		It("mzML 01 MS2 Parent Scan", func() {
			Expect(spec.Precursor.ParentScan).To(Equal("2"))
		})

		It("mzML 01 MS2 Charge State", func() {
			Expect(spec.Precursor.ChargeState).To(Equal(2))
		})

		It("mzML 01 MS2 Parent Selected Ion", func() {
			Expect(spec.Precursor.SelectedIon).To(Equal(391.201019287109))
		})

		It("mzML 01 MS2 Parent Target Ion", func() {
			Expect(spec.Precursor.TargetIon).To(Equal(391.2))
		})

		It("mzML 01 MS2 Charge Peak Intensity", func() {
			Expect(spec.Precursor.PeakIntensity).To(Equal(3.58558525e+06))
		})

		It("mzML 01 MS2 Parent Isolatio nWindow Lower Offset", func() {
			Expect(spec.Precursor.IsolationWindowLowerOffset).To(Equal(0.34999999404))
		})

		It("mzML 01 MS2 Charge Isolation Window Upper Offset", func() {
			Expect(spec.Precursor.IsolationWindowUpperOffset).To(Equal(0.34999999404))
		})

	})

})
