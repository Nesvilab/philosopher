package mzn_test

import (
	"testing"

	"philosopher/lib/mzn"
	"philosopher/lib/tes"
)

var msd mzn.MsData
var spec mzn.Spectrum

func TestRawFileParsing(t *testing.T) {

	tes.SetupTestEnv()
	msd.Read("01_CPTAC_TMTS1-NCI7_Z_JHUZ_20170502_LUMOS.mzML", false, false, false)
	tes.ShutDowTestEnv()

}

func TestMS1Spectra(t *testing.T) {

	for _, i := range msd.Spectra {
		if i.Index == "0" && i.Scan == "1" {
			spec = i
			spec.Decode()
			break
		}
	}

	if len(msd.Spectra) != 54357 {
		t.Errorf("Spectra number is incorrect, got %d, want %d", len(msd.Spectra), 54357)
	}

	if spec.Index != "0" {
		t.Errorf("Spectrum index is incorrect, got %s, want %d", spec.Index, 54357)
	}

	if spec.Intensity.DecodedStream[0] != 9104.91796875 {
		t.Errorf("Spectrum Intensity Stream is incorrect, got %f, want %f", spec.Intensity.DecodedStream[0], 9104.91796875)
	}

	if spec.Mz.DecodedStream[0] != 350.1635437011719 {
		t.Errorf("Spectrum index is incorrect, got %f, want %f", spec.Mz.DecodedStream[0], 350.1635437011719)
	}

	if spec.Index != "0" {
		t.Errorf("Spectrum index is incorrect, got %s, want %d", spec.Index, 54357)
	}

}

func TestMS2Spectra(t *testing.T) {

	for _, i := range msd.Spectra {
		if i.Index == "2" && i.Scan == "3" {
			spec = i
			spec.Decode()
			break
		}
	}

	if len(spec.Mz.DecodedStream) != 231 {
		t.Errorf("MS2 Spectra number is incorrect, got %d, want %d", len(spec.Mz.DecodedStream), 231)
	}

	if spec.Index != "2" {
		t.Errorf("Spectrum index is incorrect, got %s, want %d", spec.Index, 2)
	}

	if spec.Scan != "3" {
		t.Errorf("Spectrum scan is incorrect, got %s, want %d", spec.Scan, 3)
	}

	if spec.Intensity.DecodedStream[0] != 371635.9375 {
		t.Errorf("Spectrum Intensity is incorrect, got %f, want %f", spec.Intensity.DecodedStream[0], 371635.9375)
	}

	if spec.Mz.DecodedStream[0] != 110.07147216796875 {
		t.Errorf("Spectrum MZ is incorrect, got %f, want %f", spec.Mz.DecodedStream[0], 110.07147216796875)
	}

	if spec.Precursor.ParentIndex != "1" {
		t.Errorf("Spectrum parent index is incorrect, got %s, want %d", spec.Precursor.ParentIndex, 1)
	}

	if spec.Precursor.ParentScan != "2" {
		t.Errorf("Spectrum parent scan is incorrect, got %s, want %d", spec.Precursor.ParentScan, 2)
	}

	if spec.Precursor.ChargeState != 2 {
		t.Errorf("Spectrum charge state is incorrect, got %d, want %d", spec.Precursor.ChargeState, 2)
	}

	if spec.Precursor.SelectedIon != 391.201019287109 {
		t.Errorf("Spectrum selected ion is incorrect, got %f want %f", spec.Precursor.SelectedIon, 391.201019287109)
	}

	if spec.Precursor.TargetIon != 391.2 {
		t.Errorf("Spectrum target ion is incorrect, got %f, want %f", spec.Precursor.TargetIon, 391.2)
	}

	if spec.Precursor.PeakIntensity != 3.58558525e+06 {
		t.Errorf("Spectrum precursor intensity is incorrect, got %f, want %f", spec.Precursor.PeakIntensity, 3.58558525e+06)
	}

	if spec.Precursor.IsolationWindowLowerOffset != 0.34999999404 {
		t.Errorf("Spectrum number is incorrect, got %f, want %f", spec.Precursor.IsolationWindowLowerOffset, 0.34999999404)
	}
}
