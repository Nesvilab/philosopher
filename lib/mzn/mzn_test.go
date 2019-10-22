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
			msd.Read("b1906_293T_proteinID_01A_QE3_122212.mzML", false, false, false)
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
			Expect(len(msd.Spectra)).To(Equal(67599))
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
			Expect(spec.Intensity.DecodedStream[0]).To(Equal(930.33349609375))
		})

		It("mzML 01A MS1 MZ", func() {
			Expect(spec.Mz.DecodedStream[0]).To(Equal(301.1378479003906))
		})

		It("Read mzML 01A MS2 spectra", func() {
			for _, i := range msd.Spectra {
				if i.Index == "2135" && i.Scan == "2136" {
					spec = i
					spec.Decode()
					break
				}
			}
			Expect(len(spec.Mz.DecodedStream)).To(Equal(2522))
		})

		It("mzML 01 MS2 spectra Index", func() {
			Expect(spec.Index).To(Equal("2135"))
		})

		It("mzML 01 MS2 spectra Scan", func() {
			Expect(spec.Scan).To(Equal("2136"))
		})

		It("mzML 01 MS1 intensities", func() {
			Expect(spec.Intensity.DecodedStream[0]).To(Equal(11408.326171875))
		})

		It("mzML 01 MS2 MZ", func() {
			Expect(spec.Mz.DecodedStream[0]).To(Equal(300.06072998046875))
		})

		It("mzML 01 MS2 Parent Index", func() {
			Expect(spec.Precursor.ParentIndex).To(Equal(""))
		})

		It("mzML 01 MS2 Parent Scan", func() {
			Expect(spec.Precursor.ParentScan).To(Equal(""))
		})

		It("mzML 01 MS2 Charge State", func() {
			Expect(spec.Precursor.ChargeState).To(Equal(0))
		})

		It("mzML 01 MS2 Parent Selected Ion", func() {
			Expect(spec.Precursor.SelectedIon).To(Equal(0.0))
		})

		It("mzML 01 MS2 Parent Target Ion", func() {
			Expect(spec.Precursor.TargetIon).To(Equal(0.0))
		})

		It("mzML 01 MS2 Charge Peak Intensity", func() {
			Expect(spec.Precursor.PeakIntensity).To(Equal(0.0))
		})

		It("mzML 01 MS2 Parent Isolatio nWindow Lower Offset", func() {
			Expect(spec.Precursor.IsolationWindowLowerOffset).To(Equal(0.0))
		})

		It("mzML 01 MS2 Charge Isolation Window Upper Offset", func() {
			Expect(spec.Precursor.IsolationWindowUpperOffset).To(Equal(0.0))
		})

	})

})
