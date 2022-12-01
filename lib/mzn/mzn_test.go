package mzn_test

import (
	"testing"

	"philosopher/lib/mzn"
	"philosopher/lib/tes"
	"philosopher/lib/uti"
)

var msd mzn.MsData
var spec mzn.Spectrum

func TestRawFileParsing(t *testing.T) {

	tes.SetupTestEnv()
	msd.Read("z04397_tc-o238g-setB_MS3.mzML")
	tes.ShutDowTestEnv()

}

func TestMS1Spectra(t *testing.T) {

	for _, i := range msd.Spectra {
		if i.Index == "100" && i.Scan == "101" {
			spec = i
			spec.Decode()
			break
		}
	}

	if len(msd.Spectra) != 80926 {
		t.Errorf("Spectra number is incorrect, got %d, want %d", len(msd.Spectra), 80926)
	}

	if spec.Index != "100" {
		t.Errorf("Spectrum index is incorrect, got %s, want %d", spec.Index, 100)
	}

	if uti.ToFixed(spec.Intensity.DecodedStream[0], 3) != 11188.795 {
		t.Errorf("Spectrum Intensity Stream is incorrect, got %f, want %f", spec.Intensity.DecodedStream[0], 11188.795)
	}

	if uti.ToFixed(spec.Mz.DecodedStream[0], 3) != 400.056 {
		t.Errorf("Spectrum index is incorrect, got %f, want %f", spec.Mz.DecodedStream[0], 400.056)
	}

	if spec.Index != "100" {
		t.Errorf("Spectrum index is incorrect, got %s, want %d", spec.Index, 100)
	}

}

func TestMS2Spectra(t *testing.T) {

	for _, i := range msd.Spectra {
		if i.Index == "101" && i.Scan == "102" {
			spec = i
			spec.Decode()
			break
		}
	}

	if len(spec.Mz.DecodedStream) != 107 {
		t.Errorf("MS2 Spectra number is incorrect, got %d, want %d", len(spec.Mz.DecodedStream), 107)
	}

	if spec.Index != "101" {
		t.Errorf("Spectrum index is incorrect, got %s, want %d", spec.Index, 101)
	}

	if spec.Scan != "102" {
		t.Errorf("Spectrum scan is incorrect, got %s, want %d", spec.Scan, 102)
	}

	if uti.ToFixed(spec.Intensity.DecodedStream[0], 4) != 4.1324 {
		t.Errorf("Spectrum Intensity is incorrect, got %f, want %f", spec.Intensity.DecodedStream[0], 4.1324)
	}

	if uti.ToFixed(spec.Mz.DecodedStream[0], 3) != 134.989 {
		t.Errorf("Spectrum MZ is incorrect, got %f, want %f", spec.Mz.DecodedStream[0], 134.989)
	}

	if spec.Precursor.ParentIndex != "99" {
		t.Errorf("Spectrum parent index is incorrect, got %s, want %d", spec.Precursor.ParentIndex, 99)
	}

	if spec.Precursor.ParentScan != "100" {
		t.Errorf("Spectrum parent scan is incorrect, got %s, want %d", spec.Precursor.ParentScan, 100)
	}

	if spec.Precursor.ChargeState != 2 {
		t.Errorf("Spectrum charge state is incorrect, got %d, want %d", spec.Precursor.ChargeState, 2)
	}

	if uti.ToFixed(spec.Precursor.SelectedIon, 4) != 423.7361 {
		t.Errorf("Spectrum selected ion is incorrect, got %f want %f", spec.Precursor.SelectedIon, 423.7361)
	}

	if uti.ToFixed(spec.Precursor.TargetIon, 4) != 423.7361 {
		t.Errorf("Spectrum target ion is incorrect, got %f, want %f", spec.Precursor.TargetIon, 423.7361)
	}

	if uti.ToFixed(spec.Precursor.SelectedIonIntensity, 4) != 204667.7973 {
		t.Errorf("Spectrum precursor intensity is incorrect, got %f, want %f", spec.Precursor.SelectedIonIntensity, 204667.7973)
	}

	if uti.ToFixed(spec.Precursor.IsolationWindowLowerOffset, 4) != 0.2500 {
		t.Errorf("Spectrum number is incorrect, got %f, want %f", spec.Precursor.IsolationWindowLowerOffset, 0.2500)
	}
}
