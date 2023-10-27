package qua

import (
	"errors"
	"fmt"
	"github.com/Nesvilab/philosopher/lib/ibt"
	"github.com/Nesvilab/philosopher/lib/xta"
	"math"
	"strings"

	"github.com/Nesvilab/philosopher/lib/id"
	"github.com/Nesvilab/philosopher/lib/iso"
	"github.com/Nesvilab/philosopher/lib/msg"
	"github.com/Nesvilab/philosopher/lib/mzn"
	"github.com/Nesvilab/philosopher/lib/rep"
	"github.com/Nesvilab/philosopher/lib/scl"
	"github.com/Nesvilab/philosopher/lib/tmt"
	"github.com/Nesvilab/philosopher/lib/trq"
	"github.com/Nesvilab/philosopher/lib/xta2"
)

const (
	mzDeltaWindow float64 = 0.5
)

// prepareLabelStructureWithMS2 instantiates the Label objects and maps them against the fragment scans in order to get the channel intensities
func prepareLabelStructureWithMS2(dir, format, brand, plex string, tol float64, mz mzn.MsData) map[string]iso.Labels {

	// get all spectra names from PSMs and create the label list
	var labels = make(map[string]iso.Labels)
	ppmPrecision := tol / math.Pow(10, 6)

	for _, i := range mz.Spectra {
		if i.Level == "2" {

			var labelData iso.Labels
			if brand == "tmt" {
				labelData = tmt.New(plex)
			} else if brand == "itraq" {
				labelData = trq.New(plex)
			} else if brand == "sclip" {
				labelData = scl.New(plex)
			} else if brand == "ibt" {
				labelData = ibt.New(plex)
			} else if brand == "xtag" {
				labelData = xta.New(plex)
			} else if brand == "xtag2" {
				labelData = xta2.New(plex)
			}

			// left-pad the spectrum scan
			paddedScan := fmt.Sprintf("%05s", i.Scan)

			labelData.Index = i.Index
			labelData.Scan = paddedScan
			labelData.ChargeState = i.Precursor.ChargeState

			for j := range i.Mz.DecodedStream {

				if i.Mz.DecodedStream[j] <= (labelData.Channel1.Mz+(ppmPrecision*labelData.Channel1.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel1.Mz-(ppmPrecision*labelData.Channel1.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel1.Intensity {
						labelData.Channel1.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel2.Mz+(ppmPrecision*labelData.Channel2.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel2.Mz-(ppmPrecision*labelData.Channel2.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel2.Intensity {
						labelData.Channel2.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel3.Mz+(ppmPrecision*labelData.Channel3.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel3.Mz-(ppmPrecision*labelData.Channel3.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel3.Intensity {
						labelData.Channel3.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel4.Mz+(ppmPrecision*labelData.Channel4.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel4.Mz-(ppmPrecision*labelData.Channel4.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel4.Intensity {
						labelData.Channel4.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel5.Mz+(ppmPrecision*labelData.Channel5.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel5.Mz-(ppmPrecision*labelData.Channel5.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel5.Intensity {
						labelData.Channel5.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel6.Mz+(ppmPrecision*labelData.Channel6.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel6.Mz-(ppmPrecision*labelData.Channel6.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel6.Intensity {
						labelData.Channel6.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel7.Mz+(ppmPrecision*labelData.Channel7.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel7.Mz-(ppmPrecision*labelData.Channel7.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel7.Intensity {
						labelData.Channel7.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel8.Mz+(ppmPrecision*labelData.Channel8.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel8.Mz-(ppmPrecision*labelData.Channel8.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel8.Intensity {
						labelData.Channel8.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel9.Mz+(ppmPrecision*labelData.Channel9.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel9.Mz-(ppmPrecision*labelData.Channel9.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel9.Intensity {
						labelData.Channel9.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel10.Mz+(ppmPrecision*labelData.Channel10.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel10.Mz-(ppmPrecision*labelData.Channel10.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel10.Intensity {
						labelData.Channel10.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel11.Mz+(ppmPrecision*labelData.Channel11.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel11.Mz-(ppmPrecision*labelData.Channel11.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel11.Intensity {
						labelData.Channel11.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel12.Mz+(ppmPrecision*labelData.Channel12.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel12.Mz-(ppmPrecision*labelData.Channel12.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel12.Intensity {
						labelData.Channel12.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel13.Mz+(ppmPrecision*labelData.Channel13.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel13.Mz-(ppmPrecision*labelData.Channel13.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel13.Intensity {
						labelData.Channel13.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel14.Mz+(ppmPrecision*labelData.Channel14.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel14.Mz-(ppmPrecision*labelData.Channel14.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel14.Intensity {
						labelData.Channel14.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel15.Mz+(ppmPrecision*labelData.Channel15.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel15.Mz-(ppmPrecision*labelData.Channel15.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel15.Intensity {
						labelData.Channel15.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel16.Mz+(ppmPrecision*labelData.Channel16.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel16.Mz-(ppmPrecision*labelData.Channel16.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel16.Intensity {
						labelData.Channel16.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel17.Mz+(ppmPrecision*labelData.Channel17.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel17.Mz-(ppmPrecision*labelData.Channel17.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel17.Intensity {
						labelData.Channel17.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel18.Mz+(ppmPrecision*labelData.Channel18.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel18.Mz-(ppmPrecision*labelData.Channel18.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel18.Intensity {
						labelData.Channel18.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel19.Mz+(ppmPrecision*labelData.Channel19.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel19.Mz-(ppmPrecision*labelData.Channel19.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel19.Intensity {
						labelData.Channel19.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel20.Mz+(ppmPrecision*labelData.Channel20.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel20.Mz-(ppmPrecision*labelData.Channel20.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel20.Intensity {
						labelData.Channel20.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel21.Mz+(ppmPrecision*labelData.Channel21.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel21.Mz-(ppmPrecision*labelData.Channel21.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel21.Intensity {
						labelData.Channel21.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel22.Mz+(ppmPrecision*labelData.Channel22.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel22.Mz-(ppmPrecision*labelData.Channel22.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel22.Intensity {
						labelData.Channel22.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel23.Mz+(ppmPrecision*labelData.Channel23.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel23.Mz-(ppmPrecision*labelData.Channel23.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel23.Intensity {
						labelData.Channel23.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel24.Mz+(ppmPrecision*labelData.Channel24.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel24.Mz-(ppmPrecision*labelData.Channel24.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel24.Intensity {
						labelData.Channel24.Intensity = i.Intensity.DecodedStream[j]
					}
				}
				if i.Mz.DecodedStream[j] <= (labelData.Channel25.Mz+(ppmPrecision*labelData.Channel25.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel25.Mz-(ppmPrecision*labelData.Channel25.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel25.Intensity {
						labelData.Channel25.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel26.Mz+(ppmPrecision*labelData.Channel26.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel26.Mz-(ppmPrecision*labelData.Channel26.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel26.Intensity {
						labelData.Channel26.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel27.Mz+(ppmPrecision*labelData.Channel27.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel27.Mz-(ppmPrecision*labelData.Channel27.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel27.Intensity {
						labelData.Channel27.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel28.Mz+(ppmPrecision*labelData.Channel28.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel28.Mz-(ppmPrecision*labelData.Channel28.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel28.Intensity {
						labelData.Channel28.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel29.Mz+(ppmPrecision*labelData.Channel29.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel29.Mz-(ppmPrecision*labelData.Channel29.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel29.Intensity {
						labelData.Channel29.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel30.Mz+(ppmPrecision*labelData.Channel30.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel30.Mz-(ppmPrecision*labelData.Channel30.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel30.Intensity {
						labelData.Channel30.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel31.Mz+(ppmPrecision*labelData.Channel31.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel31.Mz-(ppmPrecision*labelData.Channel31.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel31.Intensity {
						labelData.Channel31.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel32.Mz+(ppmPrecision*labelData.Channel32.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel32.Mz-(ppmPrecision*labelData.Channel32.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel32.Intensity {
						labelData.Channel32.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if brand != "sclip" && brand != "xtag" && i.Mz.DecodedStream[j] > 137 {
					break
				} else if i.Mz.DecodedStream[j] > 450 {
					break
				}

			}

			labels[paddedScan] = labelData

		}
	}

	return labels
}

// prepareLabelStructureWithMS3 instantiates the Label objects and maps them against the fragment scans in order to get the channel intensities
func prepareLabelStructureWithMS3(dir, format, brand, plex string, tol float64, mz mzn.MsData) map[string]iso.Labels {

	// get all spectra names from PSMs and create the label list
	var labels = make(map[string]iso.Labels)
	ppmPrecision := tol / math.Pow(10, 6)

	for _, i := range mz.Spectra {
		if i.Level == "3" {

			var labelData iso.Labels
			if brand == "tmt" {
				labelData = tmt.New(plex)
			} else if brand == "itraq" {
				labelData = trq.New(plex)
			} else if brand == "sclip" {
				labelData = scl.New(plex)
			} else if brand == "ibt" {
				labelData = ibt.New(plex)
			} else if brand == "xtag" {
				labelData = xta.New(plex)
			} else if brand == "xtag2" {
				labelData = xta2.New(plex)
			}

			// left-pad the spectrum scan
			paddedScan := fmt.Sprintf("%05s", i.Scan)
			precPaddedScan := fmt.Sprintf("%05s", i.Precursor.ParentScan)

			labelData.Index = i.Index
			labelData.Scan = paddedScan
			labelData.ChargeState = i.Precursor.ChargeState

			for j := range i.Mz.DecodedStream {

				if i.Mz.DecodedStream[j] <= (labelData.Channel1.Mz+(ppmPrecision*labelData.Channel1.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel1.Mz-(ppmPrecision*labelData.Channel1.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel1.Intensity {
						labelData.Channel1.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel2.Mz+(ppmPrecision*labelData.Channel2.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel2.Mz-(ppmPrecision*labelData.Channel2.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel2.Intensity {
						labelData.Channel2.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel3.Mz+(ppmPrecision*labelData.Channel3.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel3.Mz-(ppmPrecision*labelData.Channel3.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel3.Intensity {
						labelData.Channel3.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel4.Mz+(ppmPrecision*labelData.Channel4.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel4.Mz-(ppmPrecision*labelData.Channel4.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel4.Intensity {
						labelData.Channel4.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel5.Mz+(ppmPrecision*labelData.Channel5.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel5.Mz-(ppmPrecision*labelData.Channel5.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel5.Intensity {
						labelData.Channel5.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel6.Mz+(ppmPrecision*labelData.Channel6.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel6.Mz-(ppmPrecision*labelData.Channel6.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel6.Intensity {
						labelData.Channel6.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel7.Mz+(ppmPrecision*labelData.Channel7.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel7.Mz-(ppmPrecision*labelData.Channel7.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel7.Intensity {
						labelData.Channel7.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel8.Mz+(ppmPrecision*labelData.Channel8.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel8.Mz-(ppmPrecision*labelData.Channel8.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel8.Intensity {
						labelData.Channel8.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel9.Mz+(ppmPrecision*labelData.Channel9.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel9.Mz-(ppmPrecision*labelData.Channel9.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel9.Intensity {
						labelData.Channel9.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel10.Mz+(ppmPrecision*labelData.Channel10.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel10.Mz-(ppmPrecision*labelData.Channel10.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel10.Intensity {
						labelData.Channel10.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel11.Mz+(ppmPrecision*labelData.Channel11.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel11.Mz-(ppmPrecision*labelData.Channel11.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel11.Intensity {
						labelData.Channel11.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel12.Mz+(ppmPrecision*labelData.Channel12.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel12.Mz-(ppmPrecision*labelData.Channel12.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel12.Intensity {
						labelData.Channel12.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel13.Mz+(ppmPrecision*labelData.Channel13.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel13.Mz-(ppmPrecision*labelData.Channel13.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel13.Intensity {
						labelData.Channel13.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel14.Mz+(ppmPrecision*labelData.Channel14.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel14.Mz-(ppmPrecision*labelData.Channel14.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel14.Intensity {
						labelData.Channel14.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel15.Mz+(ppmPrecision*labelData.Channel15.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel15.Mz-(ppmPrecision*labelData.Channel15.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel15.Intensity {
						labelData.Channel15.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel16.Mz+(ppmPrecision*labelData.Channel16.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel16.Mz-(ppmPrecision*labelData.Channel16.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel16.Intensity {
						labelData.Channel16.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel17.Mz+(ppmPrecision*labelData.Channel17.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel17.Mz-(ppmPrecision*labelData.Channel17.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel17.Intensity {
						labelData.Channel17.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel18.Mz+(ppmPrecision*labelData.Channel18.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel18.Mz-(ppmPrecision*labelData.Channel18.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel18.Intensity {
						labelData.Channel18.Intensity = i.Intensity.DecodedStream[j]
					}
				}
				if i.Mz.DecodedStream[j] <= (labelData.Channel19.Mz+(ppmPrecision*labelData.Channel19.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel19.Mz-(ppmPrecision*labelData.Channel19.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel19.Intensity {
						labelData.Channel19.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel20.Mz+(ppmPrecision*labelData.Channel20.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel20.Mz-(ppmPrecision*labelData.Channel20.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel20.Intensity {
						labelData.Channel20.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel21.Mz+(ppmPrecision*labelData.Channel21.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel21.Mz-(ppmPrecision*labelData.Channel21.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel21.Intensity {
						labelData.Channel21.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel22.Mz+(ppmPrecision*labelData.Channel22.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel22.Mz-(ppmPrecision*labelData.Channel22.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel22.Intensity {
						labelData.Channel22.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel23.Mz+(ppmPrecision*labelData.Channel23.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel23.Mz-(ppmPrecision*labelData.Channel23.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel23.Intensity {
						labelData.Channel23.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel24.Mz+(ppmPrecision*labelData.Channel24.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel24.Mz-(ppmPrecision*labelData.Channel24.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel24.Intensity {
						labelData.Channel24.Intensity = i.Intensity.DecodedStream[j]
					}
				}
				if i.Mz.DecodedStream[j] <= (labelData.Channel25.Mz+(ppmPrecision*labelData.Channel25.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel25.Mz-(ppmPrecision*labelData.Channel25.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel25.Intensity {
						labelData.Channel25.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel26.Mz+(ppmPrecision*labelData.Channel26.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel26.Mz-(ppmPrecision*labelData.Channel26.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel26.Intensity {
						labelData.Channel26.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel27.Mz+(ppmPrecision*labelData.Channel27.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel27.Mz-(ppmPrecision*labelData.Channel27.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel27.Intensity {
						labelData.Channel27.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel28.Mz+(ppmPrecision*labelData.Channel28.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel28.Mz-(ppmPrecision*labelData.Channel28.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel28.Intensity {
						labelData.Channel28.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel29.Mz+(ppmPrecision*labelData.Channel29.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel29.Mz-(ppmPrecision*labelData.Channel29.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel29.Intensity {
						labelData.Channel29.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel30.Mz+(ppmPrecision*labelData.Channel30.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel30.Mz-(ppmPrecision*labelData.Channel30.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel30.Intensity {
						labelData.Channel30.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel31.Mz+(ppmPrecision*labelData.Channel31.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel31.Mz-(ppmPrecision*labelData.Channel31.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel31.Intensity {
						labelData.Channel31.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if i.Mz.DecodedStream[j] <= (labelData.Channel32.Mz+(ppmPrecision*labelData.Channel32.Mz)) && i.Mz.DecodedStream[j] >= (labelData.Channel32.Mz-(ppmPrecision*labelData.Channel32.Mz)) {
					if i.Intensity.DecodedStream[j] > labelData.Channel32.Intensity {
						labelData.Channel32.Intensity = i.Intensity.DecodedStream[j]
					}
				}

				if brand != "sclip" && brand != "xtag" && i.Mz.DecodedStream[j] > 137 {
					break
				} else if i.Mz.DecodedStream[j] > 450 {
					break
				}

			}

			labels[precPaddedScan] = labelData

		}
	}

	return labels
}

// mapLabeledSpectra maps all labeled spectra to PSMs
func mapLabeledSpectra(labels map[string]iso.Labels, purity float64, evi []rep.PSMEvidence) []rep.PSMEvidence {

	for i := range evi {

		split := strings.Split(evi[i].Spectrum, ".")

		// referenced by scan number
		v, ok := labels[split[2]]
		if ok {

			evi[i].Labels.Spectrum = v.Spectrum
			evi[i].Labels.Index = v.Index
			evi[i].Labels.Scan = v.Scan

			evi[i].Labels.Channel1.Intensity = v.Channel1.Intensity
			evi[i].Labels.Channel1.CustomName = v.Channel1.CustomName

			evi[i].Labels.Channel2.Intensity = v.Channel2.Intensity
			evi[i].Labels.Channel2.CustomName = v.Channel2.CustomName

			evi[i].Labels.Channel3.Intensity = v.Channel3.Intensity
			evi[i].Labels.Channel3.CustomName = v.Channel3.CustomName

			evi[i].Labels.Channel4.Intensity = v.Channel4.Intensity
			evi[i].Labels.Channel4.CustomName = v.Channel4.CustomName

			evi[i].Labels.Channel5.Intensity = v.Channel5.Intensity
			evi[i].Labels.Channel5.CustomName = v.Channel5.CustomName

			evi[i].Labels.Channel6.Intensity = v.Channel6.Intensity
			evi[i].Labels.Channel6.CustomName = v.Channel6.CustomName

			evi[i].Labels.Channel7.Intensity = v.Channel7.Intensity
			evi[i].Labels.Channel7.CustomName = v.Channel7.CustomName

			evi[i].Labels.Channel8.Intensity = v.Channel8.Intensity
			evi[i].Labels.Channel8.CustomName = v.Channel8.CustomName

			evi[i].Labels.Channel9.Intensity = v.Channel9.Intensity
			evi[i].Labels.Channel9.CustomName = v.Channel9.CustomName

			evi[i].Labels.Channel10.Intensity = v.Channel10.Intensity
			evi[i].Labels.Channel10.CustomName = v.Channel10.CustomName

			evi[i].Labels.Channel11.Intensity = v.Channel11.Intensity
			evi[i].Labels.Channel11.CustomName = v.Channel11.CustomName

			evi[i].Labels.Channel12.Intensity = v.Channel12.Intensity
			evi[i].Labels.Channel12.CustomName = v.Channel12.CustomName

			evi[i].Labels.Channel13.Intensity = v.Channel13.Intensity
			evi[i].Labels.Channel13.CustomName = v.Channel13.CustomName

			evi[i].Labels.Channel14.Intensity = v.Channel14.Intensity
			evi[i].Labels.Channel14.CustomName = v.Channel14.CustomName

			evi[i].Labels.Channel15.Intensity = v.Channel15.Intensity
			evi[i].Labels.Channel15.CustomName = v.Channel15.CustomName

			evi[i].Labels.Channel16.Intensity = v.Channel16.Intensity
			evi[i].Labels.Channel16.CustomName = v.Channel16.CustomName

			evi[i].Labels.Channel17.Intensity = v.Channel17.Intensity
			evi[i].Labels.Channel17.CustomName = v.Channel17.CustomName

			evi[i].Labels.Channel18.Intensity = v.Channel18.Intensity
			evi[i].Labels.Channel18.CustomName = v.Channel18.CustomName

			evi[i].Labels.Channel19.Intensity = v.Channel19.Intensity
			evi[i].Labels.Channel19.CustomName = v.Channel19.CustomName

			evi[i].Labels.Channel20.Intensity = v.Channel20.Intensity
			evi[i].Labels.Channel20.CustomName = v.Channel20.CustomName

			evi[i].Labels.Channel21.Intensity = v.Channel21.Intensity
			evi[i].Labels.Channel21.CustomName = v.Channel21.CustomName

			evi[i].Labels.Channel22.Intensity = v.Channel22.Intensity
			evi[i].Labels.Channel22.CustomName = v.Channel22.CustomName

			evi[i].Labels.Channel23.Intensity = v.Channel23.Intensity
			evi[i].Labels.Channel23.CustomName = v.Channel23.CustomName

			evi[i].Labels.Channel24.Intensity = v.Channel24.Intensity
			evi[i].Labels.Channel24.CustomName = v.Channel24.CustomName

			evi[i].Labels.Channel25.Intensity = v.Channel25.Intensity
			evi[i].Labels.Channel25.CustomName = v.Channel25.CustomName

			evi[i].Labels.Channel26.Intensity = v.Channel26.Intensity
			evi[i].Labels.Channel26.CustomName = v.Channel26.CustomName

			evi[i].Labels.Channel27.Intensity = v.Channel27.Intensity
			evi[i].Labels.Channel27.CustomName = v.Channel27.CustomName

			evi[i].Labels.Channel28.Intensity = v.Channel28.Intensity
			evi[i].Labels.Channel28.CustomName = v.Channel28.CustomName

			evi[i].Labels.Channel29.Intensity = v.Channel29.Intensity
			evi[i].Labels.Channel29.CustomName = v.Channel29.CustomName

			evi[i].Labels.Channel30.Intensity = v.Channel30.Intensity
			evi[i].Labels.Channel30.CustomName = v.Channel30.CustomName

			evi[i].Labels.Channel31.Intensity = v.Channel31.Intensity
			evi[i].Labels.Channel31.CustomName = v.Channel31.CustomName

			evi[i].Labels.Channel32.Intensity = v.Channel32.Intensity
			evi[i].Labels.Channel32.CustomName = v.Channel32.CustomName

		}
	}

	return evi
}

// the assignment of usage is only done for general PSM, not for phosphoPSMs
func assignUsage(evi rep.Evidence, spectrumMap map[id.SpectrumType]iso.Labels) rep.Evidence {

	for i := range evi.PSM {
		_, ok := spectrumMap[evi.PSM[i].SpectrumFileName()]
		if ok {
			evi.PSM[i].Labels.IsUsed = true
		}
	}

	return evi
}

func correctUnlabelledSpectra(evi rep.Evidence) rep.Evidence {

	var counter = 0
	var rowSum float64

	for i := range evi.PSM {

		var flag = 0

		rowSum = evi.PSM[i].Labels.Channel1.Intensity + evi.PSM[i].Labels.Channel2.Intensity + evi.PSM[i].Labels.Channel3.Intensity + evi.PSM[i].Labels.Channel4.Intensity + evi.PSM[i].Labels.Channel5.Intensity + evi.PSM[i].Labels.Channel6.Intensity + evi.PSM[i].Labels.Channel7.Intensity + evi.PSM[i].Labels.Channel8.Intensity + evi.PSM[i].Labels.Channel9.Intensity + evi.PSM[i].Labels.Channel10.Intensity + evi.PSM[i].Labels.Channel11.Intensity + evi.PSM[i].Labels.Channel12.Intensity + evi.PSM[i].Labels.Channel13.Intensity + evi.PSM[i].Labels.Channel14.Intensity + evi.PSM[i].Labels.Channel15.Intensity + evi.PSM[i].Labels.Channel16.Intensity + evi.PSM[i].Labels.Channel17.Intensity + evi.PSM[i].Labels.Channel18.Intensity + evi.PSM[i].Labels.Channel19.Intensity + evi.PSM[i].Labels.Channel20.Intensity + evi.PSM[i].Labels.Channel21.Intensity + evi.PSM[i].Labels.Channel22.Intensity + evi.PSM[i].Labels.Channel23.Intensity + evi.PSM[i].Labels.Channel24.Intensity + evi.PSM[i].Labels.Channel25.Intensity + evi.PSM[i].Labels.Channel26.Intensity + evi.PSM[i].Labels.Channel27.Intensity + evi.PSM[i].Labels.Channel28.Intensity + evi.PSM[i].Labels.Channel29.Intensity + evi.PSM[i].Labels.Channel30.Intensity + evi.PSM[i].Labels.Channel31.Intensity + evi.PSM[i].Labels.Channel32.Intensity
		if rowSum > 0 {
			counter++
		}

		if len(evi.PSM[i].Modifications.IndexSlice) < 1 {
			evi.PSM[i].Labels.Channel1.Intensity = 0
			evi.PSM[i].Labels.Channel2.Intensity = 0
			evi.PSM[i].Labels.Channel3.Intensity = 0
			evi.PSM[i].Labels.Channel4.Intensity = 0
			evi.PSM[i].Labels.Channel5.Intensity = 0
			evi.PSM[i].Labels.Channel6.Intensity = 0
			evi.PSM[i].Labels.Channel7.Intensity = 0
			evi.PSM[i].Labels.Channel8.Intensity = 0
			evi.PSM[i].Labels.Channel9.Intensity = 0
			evi.PSM[i].Labels.Channel10.Intensity = 0
			evi.PSM[i].Labels.Channel11.Intensity = 0
			evi.PSM[i].Labels.Channel12.Intensity = 0
			evi.PSM[i].Labels.Channel13.Intensity = 0
			evi.PSM[i].Labels.Channel14.Intensity = 0
			evi.PSM[i].Labels.Channel15.Intensity = 0
			evi.PSM[i].Labels.Channel16.Intensity = 0
			evi.PSM[i].Labels.Channel17.Intensity = 0
			evi.PSM[i].Labels.Channel18.Intensity = 0
			evi.PSM[i].Labels.Channel19.Intensity = 0
			evi.PSM[i].Labels.Channel20.Intensity = 0
			evi.PSM[i].Labels.Channel21.Intensity = 0
			evi.PSM[i].Labels.Channel22.Intensity = 0
			evi.PSM[i].Labels.Channel23.Intensity = 0
			evi.PSM[i].Labels.Channel24.Intensity = 0
			evi.PSM[i].Labels.Channel25.Intensity = 0
			evi.PSM[i].Labels.Channel26.Intensity = 0
			evi.PSM[i].Labels.Channel27.Intensity = 0
			evi.PSM[i].Labels.Channel28.Intensity = 0
			evi.PSM[i].Labels.Channel29.Intensity = 0
			evi.PSM[i].Labels.Channel30.Intensity = 0
			evi.PSM[i].Labels.Channel31.Intensity = 0
			evi.PSM[i].Labels.Channel32.Intensity = 0

		} else {
			for _, j := range evi.PSM[i].Modifications.IndexSlice {
				//if j.MassDiff == 144.1020 || j.MassDiff == 229.1629 || j.MassDiff == 304.2072 {
				if j.MassDiff > 144 {
					flag++
				}
			}

			if flag == 0 {
				evi.PSM[i].Labels.Channel1.Intensity = 0
				evi.PSM[i].Labels.Channel2.Intensity = 0
				evi.PSM[i].Labels.Channel3.Intensity = 0
				evi.PSM[i].Labels.Channel4.Intensity = 0
				evi.PSM[i].Labels.Channel5.Intensity = 0
				evi.PSM[i].Labels.Channel6.Intensity = 0
				evi.PSM[i].Labels.Channel7.Intensity = 0
				evi.PSM[i].Labels.Channel8.Intensity = 0
				evi.PSM[i].Labels.Channel9.Intensity = 0
				evi.PSM[i].Labels.Channel10.Intensity = 0
				evi.PSM[i].Labels.Channel11.Intensity = 0
				evi.PSM[i].Labels.Channel12.Intensity = 0
				evi.PSM[i].Labels.Channel13.Intensity = 0
				evi.PSM[i].Labels.Channel14.Intensity = 0
				evi.PSM[i].Labels.Channel15.Intensity = 0
				evi.PSM[i].Labels.Channel16.Intensity = 0
				evi.PSM[i].Labels.Channel17.Intensity = 0
				evi.PSM[i].Labels.Channel18.Intensity = 0
				evi.PSM[i].Labels.Channel19.Intensity = 0
				evi.PSM[i].Labels.Channel20.Intensity = 0
				evi.PSM[i].Labels.Channel21.Intensity = 0
				evi.PSM[i].Labels.Channel22.Intensity = 0
				evi.PSM[i].Labels.Channel23.Intensity = 0
				evi.PSM[i].Labels.Channel24.Intensity = 0
				evi.PSM[i].Labels.Channel25.Intensity = 0
				evi.PSM[i].Labels.Channel26.Intensity = 0
				evi.PSM[i].Labels.Channel27.Intensity = 0
				evi.PSM[i].Labels.Channel28.Intensity = 0
				evi.PSM[i].Labels.Channel29.Intensity = 0
				evi.PSM[i].Labels.Channel30.Intensity = 0
				evi.PSM[i].Labels.Channel31.Intensity = 0
				evi.PSM[i].Labels.Channel32.Intensity = 0

			}
		}
	}

	if counter < 10 {
		msg.QuantifyingData(errors.New("no reporter ions were detected. Review your parameters, and try again"), "warning")
	}

	return evi
}

// rollUpPeptides gathers PSM info and filters them before summing the instensities to the peptide level
func rollUpPeptides(evi rep.Evidence, spectrumMap map[id.SpectrumType]iso.Labels, phosphoSpectrumMap map[id.SpectrumType]iso.Labels) rep.Evidence {

	for j := range evi.Peptides {

		evi.Peptides[j].Labels = &iso.Labels{}

		for k := range evi.Peptides[j].Spectra {

			i, ok := spectrumMap[k]
			if ok {

				//evi.Peptides[j].Labels = &iso.Labels{}
				evi.Peptides[j].Labels.Channel1.Name = i.Channel1.Name
				evi.Peptides[j].Labels.Channel1.CustomName = i.Channel1.CustomName
				evi.Peptides[j].Labels.Channel1.Mz = i.Channel1.Mz
				evi.Peptides[j].Labels.Channel1.Intensity += i.Channel1.Intensity

				evi.Peptides[j].Labels.Channel2.Name = i.Channel2.Name
				evi.Peptides[j].Labels.Channel2.CustomName = i.Channel2.CustomName
				evi.Peptides[j].Labels.Channel2.Mz = i.Channel2.Mz
				evi.Peptides[j].Labels.Channel2.Intensity += i.Channel2.Intensity

				evi.Peptides[j].Labels.Channel3.Name = i.Channel3.Name
				evi.Peptides[j].Labels.Channel3.CustomName = i.Channel3.CustomName
				evi.Peptides[j].Labels.Channel3.Mz = i.Channel3.Mz
				evi.Peptides[j].Labels.Channel3.Intensity += i.Channel3.Intensity

				evi.Peptides[j].Labels.Channel4.Name = i.Channel4.Name
				evi.Peptides[j].Labels.Channel4.CustomName = i.Channel4.CustomName
				evi.Peptides[j].Labels.Channel4.Mz = i.Channel4.Mz
				evi.Peptides[j].Labels.Channel4.Intensity += i.Channel4.Intensity

				evi.Peptides[j].Labels.Channel5.Name = i.Channel5.Name
				evi.Peptides[j].Labels.Channel5.CustomName = i.Channel5.CustomName
				evi.Peptides[j].Labels.Channel5.Mz = i.Channel5.Mz
				evi.Peptides[j].Labels.Channel5.Intensity += i.Channel5.Intensity

				evi.Peptides[j].Labels.Channel6.Name = i.Channel6.Name
				evi.Peptides[j].Labels.Channel6.CustomName = i.Channel6.CustomName
				evi.Peptides[j].Labels.Channel6.Mz = i.Channel6.Mz
				evi.Peptides[j].Labels.Channel6.Intensity += i.Channel6.Intensity

				evi.Peptides[j].Labels.Channel7.Name = i.Channel7.Name
				evi.Peptides[j].Labels.Channel7.CustomName = i.Channel7.CustomName
				evi.Peptides[j].Labels.Channel7.Mz = i.Channel7.Mz
				evi.Peptides[j].Labels.Channel7.Intensity += i.Channel7.Intensity

				evi.Peptides[j].Labels.Channel8.Name = i.Channel8.Name
				evi.Peptides[j].Labels.Channel8.CustomName = i.Channel8.CustomName
				evi.Peptides[j].Labels.Channel8.Mz = i.Channel8.Mz
				evi.Peptides[j].Labels.Channel8.Intensity += i.Channel8.Intensity

				evi.Peptides[j].Labels.Channel9.Name = i.Channel9.Name
				evi.Peptides[j].Labels.Channel9.CustomName = i.Channel9.CustomName
				evi.Peptides[j].Labels.Channel9.Mz = i.Channel9.Mz
				evi.Peptides[j].Labels.Channel9.Intensity += i.Channel9.Intensity

				evi.Peptides[j].Labels.Channel10.Name = i.Channel10.Name
				evi.Peptides[j].Labels.Channel10.CustomName = i.Channel10.CustomName
				evi.Peptides[j].Labels.Channel10.Mz = i.Channel10.Mz
				evi.Peptides[j].Labels.Channel10.Intensity += i.Channel10.Intensity

				evi.Peptides[j].Labels.Channel11.Name = i.Channel11.Name
				evi.Peptides[j].Labels.Channel11.CustomName = i.Channel11.CustomName
				evi.Peptides[j].Labels.Channel11.Mz = i.Channel11.Mz
				evi.Peptides[j].Labels.Channel11.Intensity += i.Channel11.Intensity

				evi.Peptides[j].Labels.Channel12.Name = i.Channel12.Name
				evi.Peptides[j].Labels.Channel12.CustomName = i.Channel12.CustomName
				evi.Peptides[j].Labels.Channel12.Mz = i.Channel12.Mz
				evi.Peptides[j].Labels.Channel12.Intensity += i.Channel12.Intensity

				evi.Peptides[j].Labels.Channel13.Name = i.Channel13.Name
				evi.Peptides[j].Labels.Channel13.CustomName = i.Channel13.CustomName
				evi.Peptides[j].Labels.Channel13.Mz = i.Channel13.Mz
				evi.Peptides[j].Labels.Channel13.Intensity += i.Channel13.Intensity

				evi.Peptides[j].Labels.Channel14.Name = i.Channel14.Name
				evi.Peptides[j].Labels.Channel14.CustomName = i.Channel14.CustomName
				evi.Peptides[j].Labels.Channel14.Mz = i.Channel14.Mz
				evi.Peptides[j].Labels.Channel14.Intensity += i.Channel14.Intensity

				evi.Peptides[j].Labels.Channel15.Name = i.Channel15.Name
				evi.Peptides[j].Labels.Channel15.CustomName = i.Channel15.CustomName
				evi.Peptides[j].Labels.Channel15.Mz = i.Channel15.Mz
				evi.Peptides[j].Labels.Channel15.Intensity += i.Channel15.Intensity

				evi.Peptides[j].Labels.Channel16.Name = i.Channel16.Name
				evi.Peptides[j].Labels.Channel16.CustomName = i.Channel16.CustomName
				evi.Peptides[j].Labels.Channel16.Mz = i.Channel16.Mz
				evi.Peptides[j].Labels.Channel16.Intensity += i.Channel16.Intensity

				evi.Peptides[j].Labels.Channel17.Name = i.Channel17.Name
				evi.Peptides[j].Labels.Channel17.CustomName = i.Channel17.CustomName
				evi.Peptides[j].Labels.Channel17.Mz = i.Channel17.Mz
				evi.Peptides[j].Labels.Channel17.Intensity += i.Channel17.Intensity

				evi.Peptides[j].Labels.Channel18.Name = i.Channel18.Name
				evi.Peptides[j].Labels.Channel18.CustomName = i.Channel18.CustomName
				evi.Peptides[j].Labels.Channel18.Mz = i.Channel18.Mz
				evi.Peptides[j].Labels.Channel18.Intensity += i.Channel18.Intensity

				evi.Peptides[j].Labels.Channel19.Name = i.Channel19.Name
				evi.Peptides[j].Labels.Channel19.CustomName = i.Channel19.CustomName
				evi.Peptides[j].Labels.Channel19.Mz = i.Channel19.Mz
				evi.Peptides[j].Labels.Channel19.Intensity += i.Channel19.Intensity

				evi.Peptides[j].Labels.Channel20.Name = i.Channel20.Name
				evi.Peptides[j].Labels.Channel20.CustomName = i.Channel20.CustomName
				evi.Peptides[j].Labels.Channel20.Mz = i.Channel20.Mz
				evi.Peptides[j].Labels.Channel20.Intensity += i.Channel20.Intensity

				evi.Peptides[j].Labels.Channel21.Name = i.Channel21.Name
				evi.Peptides[j].Labels.Channel21.CustomName = i.Channel21.CustomName
				evi.Peptides[j].Labels.Channel21.Mz = i.Channel21.Mz
				evi.Peptides[j].Labels.Channel21.Intensity += i.Channel21.Intensity

				evi.Peptides[j].Labels.Channel22.Name = i.Channel22.Name
				evi.Peptides[j].Labels.Channel22.CustomName = i.Channel22.CustomName
				evi.Peptides[j].Labels.Channel22.Mz = i.Channel22.Mz
				evi.Peptides[j].Labels.Channel22.Intensity += i.Channel22.Intensity

				evi.Peptides[j].Labels.Channel23.Name = i.Channel23.Name
				evi.Peptides[j].Labels.Channel23.CustomName = i.Channel23.CustomName
				evi.Peptides[j].Labels.Channel23.Mz = i.Channel23.Mz
				evi.Peptides[j].Labels.Channel23.Intensity += i.Channel23.Intensity

				evi.Peptides[j].Labels.Channel24.Name = i.Channel24.Name
				evi.Peptides[j].Labels.Channel24.CustomName = i.Channel24.CustomName
				evi.Peptides[j].Labels.Channel24.Mz = i.Channel24.Mz
				evi.Peptides[j].Labels.Channel24.Intensity += i.Channel24.Intensity

				evi.Peptides[j].Labels.Channel25.Name = i.Channel25.Name
				evi.Peptides[j].Labels.Channel25.CustomName = i.Channel25.CustomName
				evi.Peptides[j].Labels.Channel25.Mz = i.Channel25.Mz
				evi.Peptides[j].Labels.Channel25.Intensity += i.Channel25.Intensity

				evi.Peptides[j].Labels.Channel26.Name = i.Channel26.Name
				evi.Peptides[j].Labels.Channel26.CustomName = i.Channel26.CustomName
				evi.Peptides[j].Labels.Channel26.Mz = i.Channel26.Mz
				evi.Peptides[j].Labels.Channel26.Intensity += i.Channel26.Intensity

				evi.Peptides[j].Labels.Channel27.Name = i.Channel27.Name
				evi.Peptides[j].Labels.Channel27.CustomName = i.Channel27.CustomName
				evi.Peptides[j].Labels.Channel27.Mz = i.Channel27.Mz
				evi.Peptides[j].Labels.Channel27.Intensity += i.Channel27.Intensity

				evi.Peptides[j].Labels.Channel28.Name = i.Channel28.Name
				evi.Peptides[j].Labels.Channel28.CustomName = i.Channel28.CustomName
				evi.Peptides[j].Labels.Channel28.Mz = i.Channel28.Mz
				evi.Peptides[j].Labels.Channel28.Intensity += i.Channel28.Intensity

				evi.Peptides[j].Labels.Channel29.Name = i.Channel29.Name
				evi.Peptides[j].Labels.Channel29.CustomName = i.Channel29.CustomName
				evi.Peptides[j].Labels.Channel29.Mz = i.Channel29.Mz
				evi.Peptides[j].Labels.Channel29.Intensity += i.Channel29.Intensity

				evi.Peptides[j].Labels.Channel30.Name = i.Channel30.Name
				evi.Peptides[j].Labels.Channel30.CustomName = i.Channel30.CustomName
				evi.Peptides[j].Labels.Channel30.Mz = i.Channel30.Mz
				evi.Peptides[j].Labels.Channel30.Intensity += i.Channel30.Intensity

				evi.Peptides[j].Labels.Channel31.Name = i.Channel31.Name
				evi.Peptides[j].Labels.Channel31.CustomName = i.Channel31.CustomName
				evi.Peptides[j].Labels.Channel31.Mz = i.Channel31.Mz
				evi.Peptides[j].Labels.Channel31.Intensity += i.Channel31.Intensity

				evi.Peptides[j].Labels.Channel32.Name = i.Channel32.Name
				evi.Peptides[j].Labels.Channel32.CustomName = i.Channel32.CustomName
				evi.Peptides[j].Labels.Channel32.Mz = i.Channel32.Mz
				evi.Peptides[j].Labels.Channel32.Intensity += i.Channel32.Intensity

			}

			i, ok = phosphoSpectrumMap[k]
			if ok {
				evi.Peptides[j].PhosphoLabels = &iso.Labels{}
				evi.Peptides[j].PhosphoLabels.Channel1.Name = i.Channel1.Name
				evi.Peptides[j].PhosphoLabels.Channel1.CustomName = i.Channel1.CustomName
				evi.Peptides[j].PhosphoLabels.Channel1.Mz = i.Channel1.Mz
				evi.Peptides[j].PhosphoLabels.Channel1.Intensity += i.Channel1.Intensity

				evi.Peptides[j].PhosphoLabels.Channel2.Name = i.Channel2.Name
				evi.Peptides[j].PhosphoLabels.Channel2.CustomName = i.Channel2.CustomName
				evi.Peptides[j].PhosphoLabels.Channel2.Mz = i.Channel2.Mz
				evi.Peptides[j].PhosphoLabels.Channel2.Intensity += i.Channel2.Intensity

				evi.Peptides[j].PhosphoLabels.Channel3.Name = i.Channel3.Name
				evi.Peptides[j].PhosphoLabels.Channel3.CustomName = i.Channel3.CustomName
				evi.Peptides[j].PhosphoLabels.Channel3.Mz = i.Channel3.Mz
				evi.Peptides[j].PhosphoLabels.Channel3.Intensity += i.Channel3.Intensity

				evi.Peptides[j].PhosphoLabels.Channel4.Name = i.Channel4.Name
				evi.Peptides[j].PhosphoLabels.Channel4.CustomName = i.Channel4.CustomName
				evi.Peptides[j].PhosphoLabels.Channel4.Mz = i.Channel4.Mz
				evi.Peptides[j].PhosphoLabels.Channel4.Intensity += i.Channel4.Intensity

				evi.Peptides[j].PhosphoLabels.Channel5.Name = i.Channel5.Name
				evi.Peptides[j].PhosphoLabels.Channel5.CustomName = i.Channel5.CustomName
				evi.Peptides[j].PhosphoLabels.Channel5.Mz = i.Channel5.Mz
				evi.Peptides[j].PhosphoLabels.Channel5.Intensity += i.Channel5.Intensity

				evi.Peptides[j].PhosphoLabels.Channel6.Name = i.Channel6.Name
				evi.Peptides[j].PhosphoLabels.Channel6.CustomName = i.Channel6.CustomName
				evi.Peptides[j].PhosphoLabels.Channel6.Mz = i.Channel6.Mz
				evi.Peptides[j].PhosphoLabels.Channel6.Intensity += i.Channel6.Intensity

				evi.Peptides[j].PhosphoLabels.Channel7.Name = i.Channel7.Name
				evi.Peptides[j].PhosphoLabels.Channel7.CustomName = i.Channel7.CustomName
				evi.Peptides[j].PhosphoLabels.Channel7.Mz = i.Channel7.Mz
				evi.Peptides[j].PhosphoLabels.Channel7.Intensity += i.Channel7.Intensity

				evi.Peptides[j].PhosphoLabels.Channel8.Name = i.Channel8.Name
				evi.Peptides[j].PhosphoLabels.Channel8.CustomName = i.Channel8.CustomName
				evi.Peptides[j].PhosphoLabels.Channel8.Mz = i.Channel8.Mz
				evi.Peptides[j].PhosphoLabels.Channel8.Intensity += i.Channel8.Intensity

				evi.Peptides[j].PhosphoLabels.Channel9.Name = i.Channel9.Name
				evi.Peptides[j].PhosphoLabels.Channel9.CustomName = i.Channel9.CustomName
				evi.Peptides[j].PhosphoLabels.Channel9.Mz = i.Channel9.Mz
				evi.Peptides[j].PhosphoLabels.Channel9.Intensity += i.Channel9.Intensity

				evi.Peptides[j].PhosphoLabels.Channel10.Name = i.Channel10.Name
				evi.Peptides[j].PhosphoLabels.Channel10.CustomName = i.Channel10.CustomName
				evi.Peptides[j].PhosphoLabels.Channel10.Mz = i.Channel10.Mz
				evi.Peptides[j].PhosphoLabels.Channel10.Intensity += i.Channel10.Intensity

				evi.Peptides[j].PhosphoLabels.Channel11.Name = i.Channel11.Name
				evi.Peptides[j].PhosphoLabels.Channel11.CustomName = i.Channel11.CustomName
				evi.Peptides[j].PhosphoLabels.Channel11.Mz = i.Channel11.Mz
				evi.Peptides[j].PhosphoLabels.Channel11.Intensity += i.Channel11.Intensity

				evi.Peptides[j].PhosphoLabels.Channel12.Name = i.Channel12.Name
				evi.Peptides[j].PhosphoLabels.Channel12.CustomName = i.Channel12.CustomName
				evi.Peptides[j].PhosphoLabels.Channel12.Mz = i.Channel12.Mz
				evi.Peptides[j].PhosphoLabels.Channel12.Intensity += i.Channel12.Intensity

				evi.Peptides[j].PhosphoLabels.Channel13.Name = i.Channel13.Name
				evi.Peptides[j].PhosphoLabels.Channel13.CustomName = i.Channel13.CustomName
				evi.Peptides[j].PhosphoLabels.Channel13.Mz = i.Channel13.Mz
				evi.Peptides[j].PhosphoLabels.Channel13.Intensity += i.Channel13.Intensity

				evi.Peptides[j].PhosphoLabels.Channel14.Name = i.Channel14.Name
				evi.Peptides[j].PhosphoLabels.Channel14.CustomName = i.Channel14.CustomName
				evi.Peptides[j].PhosphoLabels.Channel14.Mz = i.Channel14.Mz
				evi.Peptides[j].PhosphoLabels.Channel14.Intensity += i.Channel14.Intensity

				evi.Peptides[j].PhosphoLabels.Channel15.Name = i.Channel15.Name
				evi.Peptides[j].PhosphoLabels.Channel15.CustomName = i.Channel15.CustomName
				evi.Peptides[j].PhosphoLabels.Channel15.Mz = i.Channel15.Mz
				evi.Peptides[j].PhosphoLabels.Channel15.Intensity += i.Channel15.Intensity

				evi.Peptides[j].PhosphoLabels.Channel16.Name = i.Channel16.Name
				evi.Peptides[j].PhosphoLabels.Channel16.CustomName = i.Channel16.CustomName
				evi.Peptides[j].PhosphoLabels.Channel16.Mz = i.Channel16.Mz
				evi.Peptides[j].PhosphoLabels.Channel16.Intensity += i.Channel16.Intensity

				evi.Peptides[j].PhosphoLabels.Channel17.Name = i.Channel17.Name
				evi.Peptides[j].PhosphoLabels.Channel17.CustomName = i.Channel17.CustomName
				evi.Peptides[j].PhosphoLabels.Channel17.Mz = i.Channel17.Mz
				evi.Peptides[j].PhosphoLabels.Channel17.Intensity += i.Channel17.Intensity

				evi.Peptides[j].PhosphoLabels.Channel18.Name = i.Channel18.Name
				evi.Peptides[j].PhosphoLabels.Channel18.CustomName = i.Channel18.CustomName
				evi.Peptides[j].PhosphoLabels.Channel18.Mz = i.Channel18.Mz
				evi.Peptides[j].PhosphoLabels.Channel18.Intensity += i.Channel18.Intensity

				evi.Peptides[j].PhosphoLabels.Channel19.Name = i.Channel19.Name
				evi.Peptides[j].PhosphoLabels.Channel19.CustomName = i.Channel19.CustomName
				evi.Peptides[j].PhosphoLabels.Channel19.Mz = i.Channel19.Mz
				evi.Peptides[j].PhosphoLabels.Channel19.Intensity += i.Channel19.Intensity

				evi.Peptides[j].PhosphoLabels.Channel20.Name = i.Channel20.Name
				evi.Peptides[j].PhosphoLabels.Channel20.CustomName = i.Channel20.CustomName
				evi.Peptides[j].PhosphoLabels.Channel20.Mz = i.Channel20.Mz
				evi.Peptides[j].PhosphoLabels.Channel20.Intensity += i.Channel20.Intensity

				evi.Peptides[j].PhosphoLabels.Channel21.Name = i.Channel21.Name
				evi.Peptides[j].PhosphoLabels.Channel21.CustomName = i.Channel21.CustomName
				evi.Peptides[j].PhosphoLabels.Channel21.Mz = i.Channel21.Mz
				evi.Peptides[j].PhosphoLabels.Channel21.Intensity += i.Channel21.Intensity

				evi.Peptides[j].PhosphoLabels.Channel22.Name = i.Channel22.Name
				evi.Peptides[j].PhosphoLabels.Channel22.CustomName = i.Channel22.CustomName
				evi.Peptides[j].PhosphoLabels.Channel22.Mz = i.Channel22.Mz
				evi.Peptides[j].PhosphoLabels.Channel22.Intensity += i.Channel22.Intensity

				evi.Peptides[j].PhosphoLabels.Channel23.Name = i.Channel23.Name
				evi.Peptides[j].PhosphoLabels.Channel23.CustomName = i.Channel23.CustomName
				evi.Peptides[j].PhosphoLabels.Channel23.Mz = i.Channel23.Mz
				evi.Peptides[j].PhosphoLabels.Channel23.Intensity += i.Channel23.Intensity

				evi.Peptides[j].PhosphoLabels.Channel24.Name = i.Channel24.Name
				evi.Peptides[j].PhosphoLabels.Channel24.CustomName = i.Channel24.CustomName
				evi.Peptides[j].PhosphoLabels.Channel24.Mz = i.Channel24.Mz
				evi.Peptides[j].PhosphoLabels.Channel24.Intensity += i.Channel24.Intensity

				evi.Peptides[j].PhosphoLabels.Channel25.Name = i.Channel25.Name
				evi.Peptides[j].PhosphoLabels.Channel25.CustomName = i.Channel25.CustomName
				evi.Peptides[j].PhosphoLabels.Channel25.Mz = i.Channel25.Mz
				evi.Peptides[j].PhosphoLabels.Channel25.Intensity += i.Channel25.Intensity

				evi.Peptides[j].PhosphoLabels.Channel26.Name = i.Channel26.Name
				evi.Peptides[j].PhosphoLabels.Channel26.CustomName = i.Channel26.CustomName
				evi.Peptides[j].PhosphoLabels.Channel26.Mz = i.Channel26.Mz
				evi.Peptides[j].PhosphoLabels.Channel26.Intensity += i.Channel26.Intensity

				evi.Peptides[j].PhosphoLabels.Channel27.Name = i.Channel27.Name
				evi.Peptides[j].PhosphoLabels.Channel27.CustomName = i.Channel27.CustomName
				evi.Peptides[j].PhosphoLabels.Channel27.Mz = i.Channel27.Mz
				evi.Peptides[j].PhosphoLabels.Channel27.Intensity += i.Channel27.Intensity

				evi.Peptides[j].PhosphoLabels.Channel28.Name = i.Channel28.Name
				evi.Peptides[j].PhosphoLabels.Channel28.CustomName = i.Channel28.CustomName
				evi.Peptides[j].PhosphoLabels.Channel28.Mz = i.Channel28.Mz
				evi.Peptides[j].PhosphoLabels.Channel28.Intensity += i.Channel28.Intensity

				evi.Peptides[j].PhosphoLabels.Channel29.Name = i.Channel29.Name
				evi.Peptides[j].PhosphoLabels.Channel29.CustomName = i.Channel29.CustomName
				evi.Peptides[j].PhosphoLabels.Channel29.Mz = i.Channel29.Mz
				evi.Peptides[j].PhosphoLabels.Channel29.Intensity += i.Channel29.Intensity

				evi.Peptides[j].PhosphoLabels.Channel30.Name = i.Channel30.Name
				evi.Peptides[j].PhosphoLabels.Channel30.CustomName = i.Channel30.CustomName
				evi.Peptides[j].PhosphoLabels.Channel30.Mz = i.Channel30.Mz
				evi.Peptides[j].PhosphoLabels.Channel30.Intensity += i.Channel30.Intensity

				evi.Peptides[j].PhosphoLabels.Channel31.Name = i.Channel31.Name
				evi.Peptides[j].PhosphoLabels.Channel31.CustomName = i.Channel31.CustomName
				evi.Peptides[j].PhosphoLabels.Channel31.Mz = i.Channel31.Mz
				evi.Peptides[j].PhosphoLabels.Channel31.Intensity += i.Channel31.Intensity

				evi.Peptides[j].PhosphoLabels.Channel32.Name = i.Channel32.Name
				evi.Peptides[j].PhosphoLabels.Channel32.CustomName = i.Channel32.CustomName
				evi.Peptides[j].PhosphoLabels.Channel32.Mz = i.Channel32.Mz
				evi.Peptides[j].PhosphoLabels.Channel32.Intensity += i.Channel32.Intensity
			}

		}
	}

	return evi
}

// rollUpPeptideIons gathers PSM info and filters them before summing the instensities to the peptide ION level
func rollUpPeptideIons(evi rep.Evidence, spectrumMap map[id.SpectrumType]iso.Labels, phosphoSpectrumMap map[id.SpectrumType]iso.Labels) rep.Evidence {

	for j := range evi.Ions {

		evi.Ions[j].Labels = &iso.Labels{}

		for k := range evi.Ions[j].Spectra {

			i, ok := spectrumMap[k]
			if ok {
				//evi.Ions[j].Labels = &iso.Labels{}
				evi.Ions[j].Labels.Channel1.Name = i.Channel1.Name
				evi.Ions[j].Labels.Channel1.CustomName = i.Channel1.CustomName
				evi.Ions[j].Labels.Channel1.Mz = i.Channel1.Mz
				evi.Ions[j].Labels.Channel1.Intensity += i.Channel1.Intensity

				evi.Ions[j].Labels.Channel2.Name = i.Channel2.Name
				evi.Ions[j].Labels.Channel2.CustomName = i.Channel2.CustomName
				evi.Ions[j].Labels.Channel2.Mz = i.Channel2.Mz
				evi.Ions[j].Labels.Channel2.Intensity += i.Channel2.Intensity

				evi.Ions[j].Labels.Channel3.Name = i.Channel3.Name
				evi.Ions[j].Labels.Channel3.CustomName = i.Channel3.CustomName
				evi.Ions[j].Labels.Channel3.Mz = i.Channel3.Mz
				evi.Ions[j].Labels.Channel3.Intensity += i.Channel3.Intensity

				evi.Ions[j].Labels.Channel4.Name = i.Channel4.Name
				evi.Ions[j].Labels.Channel4.CustomName = i.Channel4.CustomName
				evi.Ions[j].Labels.Channel4.Mz = i.Channel4.Mz
				evi.Ions[j].Labels.Channel4.Intensity += i.Channel4.Intensity

				evi.Ions[j].Labels.Channel5.Name = i.Channel5.Name
				evi.Ions[j].Labels.Channel5.CustomName = i.Channel5.CustomName
				evi.Ions[j].Labels.Channel5.Mz = i.Channel5.Mz
				evi.Ions[j].Labels.Channel5.Intensity += i.Channel5.Intensity

				evi.Ions[j].Labels.Channel6.Name = i.Channel6.Name
				evi.Ions[j].Labels.Channel6.CustomName = i.Channel6.CustomName
				evi.Ions[j].Labels.Channel6.Mz = i.Channel6.Mz
				evi.Ions[j].Labels.Channel6.Intensity += i.Channel6.Intensity

				evi.Ions[j].Labels.Channel7.Name = i.Channel7.Name
				evi.Ions[j].Labels.Channel7.CustomName = i.Channel7.CustomName
				evi.Ions[j].Labels.Channel7.Mz = i.Channel7.Mz
				evi.Ions[j].Labels.Channel7.Intensity += i.Channel7.Intensity

				evi.Ions[j].Labels.Channel8.Name = i.Channel8.Name
				evi.Ions[j].Labels.Channel8.CustomName = i.Channel8.CustomName
				evi.Ions[j].Labels.Channel8.Mz = i.Channel8.Mz
				evi.Ions[j].Labels.Channel8.Intensity += i.Channel8.Intensity

				evi.Ions[j].Labels.Channel9.Name = i.Channel9.Name
				evi.Ions[j].Labels.Channel9.CustomName = i.Channel9.CustomName
				evi.Ions[j].Labels.Channel9.Mz = i.Channel9.Mz
				evi.Ions[j].Labels.Channel9.Intensity += i.Channel9.Intensity

				evi.Ions[j].Labels.Channel10.Name = i.Channel10.Name
				evi.Ions[j].Labels.Channel10.CustomName = i.Channel10.CustomName
				evi.Ions[j].Labels.Channel10.Mz = i.Channel10.Mz
				evi.Ions[j].Labels.Channel10.Intensity += i.Channel10.Intensity

				evi.Ions[j].Labels.Channel11.Name = i.Channel11.Name
				evi.Ions[j].Labels.Channel11.CustomName = i.Channel11.CustomName
				evi.Ions[j].Labels.Channel11.Mz = i.Channel11.Mz
				evi.Ions[j].Labels.Channel11.Intensity += i.Channel11.Intensity

				evi.Ions[j].Labels.Channel12.Name = i.Channel12.Name
				evi.Ions[j].Labels.Channel12.CustomName = i.Channel12.CustomName
				evi.Ions[j].Labels.Channel12.Mz = i.Channel12.Mz
				evi.Ions[j].Labels.Channel12.Intensity += i.Channel12.Intensity

				evi.Ions[j].Labels.Channel13.Name = i.Channel13.Name
				evi.Ions[j].Labels.Channel13.CustomName = i.Channel13.CustomName
				evi.Ions[j].Labels.Channel13.Mz = i.Channel13.Mz
				evi.Ions[j].Labels.Channel13.Intensity += i.Channel13.Intensity

				evi.Ions[j].Labels.Channel14.Name = i.Channel14.Name
				evi.Ions[j].Labels.Channel14.CustomName = i.Channel14.CustomName
				evi.Ions[j].Labels.Channel14.Mz = i.Channel14.Mz
				evi.Ions[j].Labels.Channel14.Intensity += i.Channel14.Intensity

				evi.Ions[j].Labels.Channel15.Name = i.Channel15.Name
				evi.Ions[j].Labels.Channel15.CustomName = i.Channel15.CustomName
				evi.Ions[j].Labels.Channel15.Mz = i.Channel15.Mz
				evi.Ions[j].Labels.Channel15.Intensity += i.Channel15.Intensity

				evi.Ions[j].Labels.Channel16.Name = i.Channel16.Name
				evi.Ions[j].Labels.Channel16.CustomName = i.Channel16.CustomName
				evi.Ions[j].Labels.Channel16.Mz = i.Channel16.Mz
				evi.Ions[j].Labels.Channel16.Intensity += i.Channel16.Intensity

				evi.Ions[j].Labels.Channel17.Name = i.Channel17.Name
				evi.Ions[j].Labels.Channel17.CustomName = i.Channel17.CustomName
				evi.Ions[j].Labels.Channel17.Mz = i.Channel17.Mz
				evi.Ions[j].Labels.Channel17.Intensity += i.Channel17.Intensity

				evi.Ions[j].Labels.Channel18.Name = i.Channel18.Name
				evi.Ions[j].Labels.Channel18.CustomName = i.Channel18.CustomName
				evi.Ions[j].Labels.Channel18.Mz = i.Channel18.Mz
				evi.Ions[j].Labels.Channel18.Intensity += i.Channel18.Intensity

				evi.Ions[j].Labels.Channel19.Name = i.Channel19.Name
				evi.Ions[j].Labels.Channel19.CustomName = i.Channel19.CustomName
				evi.Ions[j].Labels.Channel19.Mz = i.Channel19.Mz
				evi.Ions[j].Labels.Channel19.Intensity += i.Channel19.Intensity

				evi.Ions[j].Labels.Channel20.Name = i.Channel20.Name
				evi.Ions[j].Labels.Channel20.CustomName = i.Channel20.CustomName
				evi.Ions[j].Labels.Channel20.Mz = i.Channel20.Mz
				evi.Ions[j].Labels.Channel20.Intensity += i.Channel20.Intensity

				evi.Ions[j].Labels.Channel21.Name = i.Channel21.Name
				evi.Ions[j].Labels.Channel21.CustomName = i.Channel21.CustomName
				evi.Ions[j].Labels.Channel21.Mz = i.Channel21.Mz
				evi.Ions[j].Labels.Channel21.Intensity += i.Channel21.Intensity

				evi.Ions[j].Labels.Channel22.Name = i.Channel22.Name
				evi.Ions[j].Labels.Channel22.CustomName = i.Channel22.CustomName
				evi.Ions[j].Labels.Channel22.Mz = i.Channel22.Mz
				evi.Ions[j].Labels.Channel22.Intensity += i.Channel22.Intensity

				evi.Ions[j].Labels.Channel23.Name = i.Channel23.Name
				evi.Ions[j].Labels.Channel23.CustomName = i.Channel23.CustomName
				evi.Ions[j].Labels.Channel23.Mz = i.Channel23.Mz
				evi.Ions[j].Labels.Channel23.Intensity += i.Channel23.Intensity

				evi.Ions[j].Labels.Channel24.Name = i.Channel24.Name
				evi.Ions[j].Labels.Channel24.CustomName = i.Channel24.CustomName
				evi.Ions[j].Labels.Channel24.Mz = i.Channel24.Mz
				evi.Ions[j].Labels.Channel24.Intensity += i.Channel24.Intensity

				evi.Ions[j].Labels.Channel25.Name = i.Channel25.Name
				evi.Ions[j].Labels.Channel25.CustomName = i.Channel25.CustomName
				evi.Ions[j].Labels.Channel25.Mz = i.Channel25.Mz
				evi.Ions[j].Labels.Channel25.Intensity += i.Channel25.Intensity

				evi.Ions[j].Labels.Channel26.Name = i.Channel26.Name
				evi.Ions[j].Labels.Channel26.CustomName = i.Channel26.CustomName
				evi.Ions[j].Labels.Channel26.Mz = i.Channel26.Mz
				evi.Ions[j].Labels.Channel26.Intensity += i.Channel26.Intensity

				evi.Ions[j].Labels.Channel27.Name = i.Channel27.Name
				evi.Ions[j].Labels.Channel27.CustomName = i.Channel27.CustomName
				evi.Ions[j].Labels.Channel27.Mz = i.Channel27.Mz
				evi.Ions[j].Labels.Channel27.Intensity += i.Channel27.Intensity

				evi.Ions[j].Labels.Channel28.Name = i.Channel28.Name
				evi.Ions[j].Labels.Channel28.CustomName = i.Channel28.CustomName
				evi.Ions[j].Labels.Channel28.Mz = i.Channel28.Mz
				evi.Ions[j].Labels.Channel28.Intensity += i.Channel28.Intensity

				evi.Ions[j].Labels.Channel29.Name = i.Channel29.Name
				evi.Ions[j].Labels.Channel29.CustomName = i.Channel29.CustomName
				evi.Ions[j].Labels.Channel29.Mz = i.Channel29.Mz
				evi.Ions[j].Labels.Channel29.Intensity += i.Channel29.Intensity

				evi.Ions[j].Labels.Channel30.Name = i.Channel30.Name
				evi.Ions[j].Labels.Channel30.CustomName = i.Channel30.CustomName
				evi.Ions[j].Labels.Channel30.Mz = i.Channel30.Mz
				evi.Ions[j].Labels.Channel30.Intensity += i.Channel30.Intensity

				evi.Ions[j].Labels.Channel31.Name = i.Channel31.Name
				evi.Ions[j].Labels.Channel31.CustomName = i.Channel31.CustomName
				evi.Ions[j].Labels.Channel31.Mz = i.Channel31.Mz
				evi.Ions[j].Labels.Channel31.Intensity += i.Channel31.Intensity

				evi.Ions[j].Labels.Channel32.Name = i.Channel32.Name
				evi.Ions[j].Labels.Channel32.CustomName = i.Channel32.CustomName
				evi.Ions[j].Labels.Channel32.Mz = i.Channel32.Mz
				evi.Ions[j].Labels.Channel32.Intensity += i.Channel32.Intensity
			}

			i, ok = phosphoSpectrumMap[k]
			if ok {
				evi.Ions[j].PhosphoLabels = &iso.Labels{}
				evi.Ions[j].PhosphoLabels.Channel1.Name = i.Channel1.Name
				evi.Ions[j].PhosphoLabels.Channel1.CustomName = i.Channel1.CustomName
				evi.Ions[j].PhosphoLabels.Channel1.Mz = i.Channel1.Mz
				evi.Ions[j].PhosphoLabels.Channel1.Intensity += i.Channel1.Intensity

				evi.Ions[j].PhosphoLabels.Channel2.Name = i.Channel2.Name
				evi.Ions[j].PhosphoLabels.Channel2.CustomName = i.Channel2.CustomName
				evi.Ions[j].PhosphoLabels.Channel2.Mz = i.Channel2.Mz
				evi.Ions[j].PhosphoLabels.Channel2.Intensity += i.Channel2.Intensity

				evi.Ions[j].PhosphoLabels.Channel3.Name = i.Channel3.Name
				evi.Ions[j].PhosphoLabels.Channel3.CustomName = i.Channel3.CustomName
				evi.Ions[j].PhosphoLabels.Channel3.Mz = i.Channel3.Mz
				evi.Ions[j].PhosphoLabels.Channel3.Intensity += i.Channel3.Intensity

				evi.Ions[j].PhosphoLabels.Channel4.Name = i.Channel4.Name
				evi.Ions[j].PhosphoLabels.Channel4.CustomName = i.Channel4.CustomName
				evi.Ions[j].PhosphoLabels.Channel4.Mz = i.Channel4.Mz
				evi.Ions[j].PhosphoLabels.Channel4.Intensity += i.Channel4.Intensity

				evi.Ions[j].PhosphoLabels.Channel5.Name = i.Channel5.Name
				evi.Ions[j].PhosphoLabels.Channel5.CustomName = i.Channel5.CustomName
				evi.Ions[j].PhosphoLabels.Channel5.Mz = i.Channel5.Mz
				evi.Ions[j].PhosphoLabels.Channel5.Intensity += i.Channel5.Intensity

				evi.Ions[j].PhosphoLabels.Channel6.Name = i.Channel6.Name
				evi.Ions[j].PhosphoLabels.Channel6.CustomName = i.Channel6.CustomName
				evi.Ions[j].PhosphoLabels.Channel6.Mz = i.Channel6.Mz
				evi.Ions[j].PhosphoLabels.Channel6.Intensity += i.Channel6.Intensity

				evi.Ions[j].PhosphoLabels.Channel7.Name = i.Channel7.Name
				evi.Ions[j].PhosphoLabels.Channel7.CustomName = i.Channel7.CustomName
				evi.Ions[j].PhosphoLabels.Channel7.Mz = i.Channel7.Mz
				evi.Ions[j].PhosphoLabels.Channel7.Intensity += i.Channel7.Intensity

				evi.Ions[j].PhosphoLabels.Channel8.Name = i.Channel8.Name
				evi.Ions[j].PhosphoLabels.Channel8.CustomName = i.Channel8.CustomName
				evi.Ions[j].PhosphoLabels.Channel8.Mz = i.Channel8.Mz
				evi.Ions[j].PhosphoLabels.Channel8.Intensity += i.Channel8.Intensity

				evi.Ions[j].PhosphoLabels.Channel9.Name = i.Channel9.Name
				evi.Ions[j].PhosphoLabels.Channel9.CustomName = i.Channel9.CustomName
				evi.Ions[j].PhosphoLabels.Channel9.Mz = i.Channel9.Mz
				evi.Ions[j].PhosphoLabels.Channel9.Intensity += i.Channel9.Intensity

				evi.Ions[j].PhosphoLabels.Channel10.Name = i.Channel10.Name
				evi.Ions[j].PhosphoLabels.Channel10.CustomName = i.Channel10.CustomName
				evi.Ions[j].PhosphoLabels.Channel10.Mz = i.Channel10.Mz
				evi.Ions[j].PhosphoLabels.Channel10.Intensity += i.Channel10.Intensity

				evi.Ions[j].PhosphoLabels.Channel11.Name = i.Channel11.Name
				evi.Ions[j].PhosphoLabels.Channel11.CustomName = i.Channel11.CustomName
				evi.Ions[j].PhosphoLabels.Channel11.Mz = i.Channel11.Mz
				evi.Ions[j].PhosphoLabels.Channel11.Intensity += i.Channel11.Intensity

				evi.Ions[j].PhosphoLabels.Channel12.Name = i.Channel12.Name
				evi.Ions[j].PhosphoLabels.Channel12.CustomName = i.Channel12.CustomName
				evi.Ions[j].PhosphoLabels.Channel12.Mz = i.Channel12.Mz
				evi.Ions[j].PhosphoLabels.Channel12.Intensity += i.Channel12.Intensity

				evi.Ions[j].PhosphoLabels.Channel13.Name = i.Channel13.Name
				evi.Ions[j].PhosphoLabels.Channel13.CustomName = i.Channel13.CustomName
				evi.Ions[j].PhosphoLabels.Channel13.Mz = i.Channel13.Mz
				evi.Ions[j].PhosphoLabels.Channel13.Intensity += i.Channel13.Intensity

				evi.Ions[j].PhosphoLabels.Channel14.Name = i.Channel14.Name
				evi.Ions[j].PhosphoLabels.Channel14.CustomName = i.Channel14.CustomName
				evi.Ions[j].PhosphoLabels.Channel14.Mz = i.Channel14.Mz
				evi.Ions[j].PhosphoLabels.Channel14.Intensity += i.Channel14.Intensity

				evi.Ions[j].PhosphoLabels.Channel15.Name = i.Channel15.Name
				evi.Ions[j].PhosphoLabels.Channel15.CustomName = i.Channel15.CustomName
				evi.Ions[j].PhosphoLabels.Channel15.Mz = i.Channel15.Mz
				evi.Ions[j].PhosphoLabels.Channel15.Intensity += i.Channel15.Intensity

				evi.Ions[j].PhosphoLabels.Channel16.Name = i.Channel16.Name
				evi.Ions[j].PhosphoLabels.Channel16.CustomName = i.Channel16.CustomName
				evi.Ions[j].PhosphoLabels.Channel16.Mz = i.Channel16.Mz
				evi.Ions[j].PhosphoLabels.Channel16.Intensity += i.Channel16.Intensity

				evi.Ions[j].PhosphoLabels.Channel17.Name = i.Channel17.Name
				evi.Ions[j].PhosphoLabels.Channel17.CustomName = i.Channel17.CustomName
				evi.Ions[j].PhosphoLabels.Channel17.Mz = i.Channel17.Mz
				evi.Ions[j].PhosphoLabels.Channel17.Intensity += i.Channel17.Intensity

				evi.Ions[j].PhosphoLabels.Channel18.Name = i.Channel18.Name
				evi.Ions[j].PhosphoLabels.Channel18.CustomName = i.Channel18.CustomName
				evi.Ions[j].PhosphoLabels.Channel18.Mz = i.Channel18.Mz
				evi.Ions[j].PhosphoLabels.Channel18.Intensity += i.Channel18.Intensity

				evi.Ions[j].PhosphoLabels.Channel19.Name = i.Channel19.Name
				evi.Ions[j].PhosphoLabels.Channel19.CustomName = i.Channel19.CustomName
				evi.Ions[j].PhosphoLabels.Channel19.Mz = i.Channel19.Mz
				evi.Ions[j].PhosphoLabels.Channel19.Intensity += i.Channel19.Intensity

				evi.Ions[j].PhosphoLabels.Channel20.Name = i.Channel20.Name
				evi.Ions[j].PhosphoLabels.Channel20.CustomName = i.Channel20.CustomName
				evi.Ions[j].PhosphoLabels.Channel20.Mz = i.Channel20.Mz
				evi.Ions[j].PhosphoLabels.Channel20.Intensity += i.Channel20.Intensity

				evi.Ions[j].PhosphoLabels.Channel21.Name = i.Channel21.Name
				evi.Ions[j].PhosphoLabels.Channel21.CustomName = i.Channel21.CustomName
				evi.Ions[j].PhosphoLabels.Channel21.Mz = i.Channel21.Mz
				evi.Ions[j].PhosphoLabels.Channel21.Intensity += i.Channel21.Intensity

				evi.Ions[j].PhosphoLabels.Channel22.Name = i.Channel22.Name
				evi.Ions[j].PhosphoLabels.Channel22.CustomName = i.Channel22.CustomName
				evi.Ions[j].PhosphoLabels.Channel22.Mz = i.Channel22.Mz
				evi.Ions[j].PhosphoLabels.Channel22.Intensity += i.Channel22.Intensity

				evi.Ions[j].PhosphoLabels.Channel23.Name = i.Channel23.Name
				evi.Ions[j].PhosphoLabels.Channel23.CustomName = i.Channel23.CustomName
				evi.Ions[j].PhosphoLabels.Channel23.Mz = i.Channel23.Mz
				evi.Ions[j].PhosphoLabels.Channel23.Intensity += i.Channel23.Intensity

				evi.Ions[j].PhosphoLabels.Channel24.Name = i.Channel24.Name
				evi.Ions[j].PhosphoLabels.Channel24.CustomName = i.Channel24.CustomName
				evi.Ions[j].PhosphoLabels.Channel24.Mz = i.Channel24.Mz
				evi.Ions[j].PhosphoLabels.Channel24.Intensity += i.Channel24.Intensity

				evi.Ions[j].PhosphoLabels.Channel25.Name = i.Channel25.Name
				evi.Ions[j].PhosphoLabels.Channel25.CustomName = i.Channel25.CustomName
				evi.Ions[j].PhosphoLabels.Channel25.Mz = i.Channel25.Mz
				evi.Ions[j].PhosphoLabels.Channel25.Intensity += i.Channel25.Intensity

				evi.Ions[j].PhosphoLabels.Channel26.Name = i.Channel26.Name
				evi.Ions[j].PhosphoLabels.Channel26.CustomName = i.Channel26.CustomName
				evi.Ions[j].PhosphoLabels.Channel26.Mz = i.Channel26.Mz
				evi.Ions[j].PhosphoLabels.Channel26.Intensity += i.Channel26.Intensity

				evi.Ions[j].PhosphoLabels.Channel27.Name = i.Channel27.Name
				evi.Ions[j].PhosphoLabels.Channel27.CustomName = i.Channel27.CustomName
				evi.Ions[j].PhosphoLabels.Channel27.Mz = i.Channel27.Mz
				evi.Ions[j].PhosphoLabels.Channel27.Intensity += i.Channel27.Intensity

				evi.Ions[j].PhosphoLabels.Channel28.Name = i.Channel28.Name
				evi.Ions[j].PhosphoLabels.Channel28.CustomName = i.Channel28.CustomName
				evi.Ions[j].PhosphoLabels.Channel28.Mz = i.Channel28.Mz
				evi.Ions[j].PhosphoLabels.Channel28.Intensity += i.Channel28.Intensity

				evi.Ions[j].PhosphoLabels.Channel29.Name = i.Channel29.Name
				evi.Ions[j].PhosphoLabels.Channel29.CustomName = i.Channel29.CustomName
				evi.Ions[j].PhosphoLabels.Channel29.Mz = i.Channel29.Mz
				evi.Ions[j].PhosphoLabels.Channel29.Intensity += i.Channel29.Intensity

				evi.Ions[j].PhosphoLabels.Channel30.Name = i.Channel30.Name
				evi.Ions[j].PhosphoLabels.Channel30.CustomName = i.Channel30.CustomName
				evi.Ions[j].PhosphoLabels.Channel30.Mz = i.Channel30.Mz
				evi.Ions[j].PhosphoLabels.Channel30.Intensity += i.Channel30.Intensity

				evi.Ions[j].PhosphoLabels.Channel31.Name = i.Channel31.Name
				evi.Ions[j].PhosphoLabels.Channel31.CustomName = i.Channel31.CustomName
				evi.Ions[j].PhosphoLabels.Channel31.Mz = i.Channel31.Mz
				evi.Ions[j].PhosphoLabels.Channel31.Intensity += i.Channel31.Intensity

				evi.Ions[j].PhosphoLabels.Channel32.Name = i.Channel32.Name
				evi.Ions[j].PhosphoLabels.Channel32.CustomName = i.Channel32.CustomName
				evi.Ions[j].PhosphoLabels.Channel32.Mz = i.Channel32.Mz
				evi.Ions[j].PhosphoLabels.Channel32.Intensity += i.Channel32.Intensity
			}

		}
	}

	return evi
}

// rollUpProteins gathers PSM info and filters them before summing the instensities to the peptide ION level
func rollUpProteins(evi rep.Evidence, spectrumMap map[id.SpectrumType]iso.Labels, phosphoSpectrumMap map[id.SpectrumType]iso.Labels) rep.Evidence {

	for j := range evi.Proteins {

		evi.Proteins[j].TotalLabels = &iso.Labels{}
		evi.Proteins[j].UniqueLabels = &iso.Labels{}
		evi.Proteins[j].URazorLabels = &iso.Labels{}

		for _, k := range evi.Proteins[j].TotalPeptideIons {
			if k.Protein == evi.Proteins[j].PartHeader {
				for l := range k.Spectra {

					i, ok := spectrumMap[l]
					if ok {
						//evi.Proteins[j].TotalLabels = &iso.Labels{}
						evi.Proteins[j].TotalLabels.Channel1.Name = i.Channel1.Name
						evi.Proteins[j].TotalLabels.Channel1.CustomName = i.Channel1.CustomName
						evi.Proteins[j].TotalLabels.Channel1.Mz = i.Channel1.Mz
						evi.Proteins[j].TotalLabels.Channel1.Intensity += i.Channel1.Intensity

						evi.Proteins[j].TotalLabels.Channel2.Name = i.Channel2.Name
						evi.Proteins[j].TotalLabels.Channel2.CustomName = i.Channel2.CustomName
						evi.Proteins[j].TotalLabels.Channel2.Mz = i.Channel2.Mz
						evi.Proteins[j].TotalLabels.Channel2.Intensity += i.Channel2.Intensity

						evi.Proteins[j].TotalLabels.Channel3.Name = i.Channel3.Name
						evi.Proteins[j].TotalLabels.Channel3.CustomName = i.Channel3.CustomName
						evi.Proteins[j].TotalLabels.Channel3.Mz = i.Channel3.Mz
						evi.Proteins[j].TotalLabels.Channel3.Intensity += i.Channel3.Intensity

						evi.Proteins[j].TotalLabels.Channel4.Name = i.Channel4.Name
						evi.Proteins[j].TotalLabels.Channel4.CustomName = i.Channel4.CustomName
						evi.Proteins[j].TotalLabels.Channel4.Mz = i.Channel4.Mz
						evi.Proteins[j].TotalLabels.Channel4.Intensity += i.Channel4.Intensity

						evi.Proteins[j].TotalLabels.Channel5.Name = i.Channel5.Name
						evi.Proteins[j].TotalLabels.Channel5.CustomName = i.Channel5.CustomName
						evi.Proteins[j].TotalLabels.Channel5.Mz = i.Channel5.Mz
						evi.Proteins[j].TotalLabels.Channel5.Intensity += i.Channel5.Intensity

						evi.Proteins[j].TotalLabels.Channel6.Name = i.Channel6.Name
						evi.Proteins[j].TotalLabels.Channel6.CustomName = i.Channel6.CustomName
						evi.Proteins[j].TotalLabels.Channel6.Mz = i.Channel6.Mz
						evi.Proteins[j].TotalLabels.Channel6.Intensity += i.Channel6.Intensity

						evi.Proteins[j].TotalLabels.Channel7.Name = i.Channel7.Name
						evi.Proteins[j].TotalLabels.Channel7.CustomName = i.Channel7.CustomName
						evi.Proteins[j].TotalLabels.Channel7.Mz = i.Channel7.Mz
						evi.Proteins[j].TotalLabels.Channel7.Intensity += i.Channel7.Intensity

						evi.Proteins[j].TotalLabels.Channel8.Name = i.Channel8.Name
						evi.Proteins[j].TotalLabels.Channel8.CustomName = i.Channel8.CustomName
						evi.Proteins[j].TotalLabels.Channel8.Mz = i.Channel8.Mz
						evi.Proteins[j].TotalLabels.Channel8.Intensity += i.Channel8.Intensity

						evi.Proteins[j].TotalLabels.Channel9.Name = i.Channel9.Name
						evi.Proteins[j].TotalLabels.Channel9.CustomName = i.Channel9.CustomName
						evi.Proteins[j].TotalLabels.Channel9.Mz = i.Channel9.Mz
						evi.Proteins[j].TotalLabels.Channel9.Intensity += i.Channel9.Intensity

						evi.Proteins[j].TotalLabels.Channel10.Name = i.Channel10.Name
						evi.Proteins[j].TotalLabels.Channel10.CustomName = i.Channel10.CustomName
						evi.Proteins[j].TotalLabels.Channel10.Mz = i.Channel10.Mz
						evi.Proteins[j].TotalLabels.Channel10.Intensity += i.Channel10.Intensity

						evi.Proteins[j].TotalLabels.Channel11.Name = i.Channel11.Name
						evi.Proteins[j].TotalLabels.Channel11.CustomName = i.Channel11.CustomName
						evi.Proteins[j].TotalLabels.Channel11.Mz = i.Channel11.Mz
						evi.Proteins[j].TotalLabels.Channel11.Intensity += i.Channel11.Intensity

						evi.Proteins[j].TotalLabels.Channel12.Name = i.Channel12.Name
						evi.Proteins[j].TotalLabels.Channel12.CustomName = i.Channel12.CustomName
						evi.Proteins[j].TotalLabels.Channel12.Mz = i.Channel12.Mz
						evi.Proteins[j].TotalLabels.Channel12.Intensity += i.Channel12.Intensity

						evi.Proteins[j].TotalLabels.Channel13.Name = i.Channel13.Name
						evi.Proteins[j].TotalLabels.Channel13.CustomName = i.Channel13.CustomName
						evi.Proteins[j].TotalLabels.Channel13.Mz = i.Channel13.Mz
						evi.Proteins[j].TotalLabels.Channel13.Intensity += i.Channel13.Intensity

						evi.Proteins[j].TotalLabels.Channel14.Name = i.Channel14.Name
						evi.Proteins[j].TotalLabels.Channel14.CustomName = i.Channel14.CustomName
						evi.Proteins[j].TotalLabels.Channel14.Mz = i.Channel14.Mz
						evi.Proteins[j].TotalLabels.Channel14.Intensity += i.Channel14.Intensity

						evi.Proteins[j].TotalLabels.Channel15.Name = i.Channel15.Name
						evi.Proteins[j].TotalLabels.Channel15.CustomName = i.Channel15.CustomName
						evi.Proteins[j].TotalLabels.Channel15.Mz = i.Channel15.Mz
						evi.Proteins[j].TotalLabels.Channel15.Intensity += i.Channel15.Intensity

						evi.Proteins[j].TotalLabels.Channel16.Name = i.Channel16.Name
						evi.Proteins[j].TotalLabels.Channel16.CustomName = i.Channel16.CustomName
						evi.Proteins[j].TotalLabels.Channel16.Mz = i.Channel16.Mz
						evi.Proteins[j].TotalLabels.Channel16.Intensity += i.Channel16.Intensity

						evi.Proteins[j].TotalLabels.Channel17.Name = i.Channel17.Name
						evi.Proteins[j].TotalLabels.Channel17.CustomName = i.Channel17.CustomName
						evi.Proteins[j].TotalLabels.Channel17.Mz = i.Channel17.Mz
						evi.Proteins[j].TotalLabels.Channel17.Intensity += i.Channel17.Intensity

						evi.Proteins[j].TotalLabels.Channel18.Name = i.Channel18.Name
						evi.Proteins[j].TotalLabels.Channel18.CustomName = i.Channel18.CustomName
						evi.Proteins[j].TotalLabels.Channel18.Mz = i.Channel18.Mz
						evi.Proteins[j].TotalLabels.Channel18.Intensity += i.Channel18.Intensity

						evi.Proteins[j].TotalLabels.Channel19.Name = i.Channel19.Name
						evi.Proteins[j].TotalLabels.Channel19.CustomName = i.Channel19.CustomName
						evi.Proteins[j].TotalLabels.Channel19.Mz = i.Channel19.Mz
						evi.Proteins[j].TotalLabels.Channel19.Intensity += i.Channel19.Intensity

						evi.Proteins[j].TotalLabels.Channel20.Name = i.Channel20.Name
						evi.Proteins[j].TotalLabels.Channel20.CustomName = i.Channel20.CustomName
						evi.Proteins[j].TotalLabels.Channel20.Mz = i.Channel20.Mz
						evi.Proteins[j].TotalLabels.Channel20.Intensity += i.Channel20.Intensity

						evi.Proteins[j].TotalLabels.Channel21.Name = i.Channel21.Name
						evi.Proteins[j].TotalLabels.Channel21.CustomName = i.Channel21.CustomName
						evi.Proteins[j].TotalLabels.Channel21.Mz = i.Channel21.Mz
						evi.Proteins[j].TotalLabels.Channel21.Intensity += i.Channel21.Intensity

						evi.Proteins[j].TotalLabels.Channel22.Name = i.Channel22.Name
						evi.Proteins[j].TotalLabels.Channel22.CustomName = i.Channel22.CustomName
						evi.Proteins[j].TotalLabels.Channel22.Mz = i.Channel22.Mz
						evi.Proteins[j].TotalLabels.Channel22.Intensity += i.Channel22.Intensity

						evi.Proteins[j].TotalLabels.Channel23.Name = i.Channel23.Name
						evi.Proteins[j].TotalLabels.Channel23.CustomName = i.Channel23.CustomName
						evi.Proteins[j].TotalLabels.Channel23.Mz = i.Channel23.Mz
						evi.Proteins[j].TotalLabels.Channel23.Intensity += i.Channel23.Intensity

						evi.Proteins[j].TotalLabels.Channel24.Name = i.Channel24.Name
						evi.Proteins[j].TotalLabels.Channel24.CustomName = i.Channel24.CustomName
						evi.Proteins[j].TotalLabels.Channel24.Mz = i.Channel24.Mz
						evi.Proteins[j].TotalLabels.Channel24.Intensity += i.Channel24.Intensity

						evi.Proteins[j].TotalLabels.Channel25.Name = i.Channel25.Name
						evi.Proteins[j].TotalLabels.Channel25.CustomName = i.Channel25.CustomName
						evi.Proteins[j].TotalLabels.Channel25.Mz = i.Channel25.Mz
						evi.Proteins[j].TotalLabels.Channel25.Intensity += i.Channel25.Intensity

						evi.Proteins[j].TotalLabels.Channel26.Name = i.Channel26.Name
						evi.Proteins[j].TotalLabels.Channel26.CustomName = i.Channel26.CustomName
						evi.Proteins[j].TotalLabels.Channel26.Mz = i.Channel26.Mz
						evi.Proteins[j].TotalLabels.Channel26.Intensity += i.Channel26.Intensity

						evi.Proteins[j].TotalLabels.Channel27.Name = i.Channel27.Name
						evi.Proteins[j].TotalLabels.Channel27.CustomName = i.Channel27.CustomName
						evi.Proteins[j].TotalLabels.Channel27.Mz = i.Channel27.Mz
						evi.Proteins[j].TotalLabels.Channel27.Intensity += i.Channel27.Intensity

						evi.Proteins[j].TotalLabels.Channel28.Name = i.Channel28.Name
						evi.Proteins[j].TotalLabels.Channel28.CustomName = i.Channel28.CustomName
						evi.Proteins[j].TotalLabels.Channel28.Mz = i.Channel28.Mz
						evi.Proteins[j].TotalLabels.Channel28.Intensity += i.Channel28.Intensity

						evi.Proteins[j].TotalLabels.Channel29.Name = i.Channel29.Name
						evi.Proteins[j].TotalLabels.Channel29.CustomName = i.Channel29.CustomName
						evi.Proteins[j].TotalLabels.Channel29.Mz = i.Channel29.Mz
						evi.Proteins[j].TotalLabels.Channel29.Intensity += i.Channel29.Intensity

						evi.Proteins[j].TotalLabels.Channel30.Name = i.Channel30.Name
						evi.Proteins[j].TotalLabels.Channel30.CustomName = i.Channel30.CustomName
						evi.Proteins[j].TotalLabels.Channel30.Mz = i.Channel30.Mz
						evi.Proteins[j].TotalLabels.Channel30.Intensity += i.Channel30.Intensity

						evi.Proteins[j].TotalLabels.Channel31.Name = i.Channel31.Name
						evi.Proteins[j].TotalLabels.Channel31.CustomName = i.Channel31.CustomName
						evi.Proteins[j].TotalLabels.Channel31.Mz = i.Channel31.Mz
						evi.Proteins[j].TotalLabels.Channel31.Intensity += i.Channel31.Intensity

						evi.Proteins[j].TotalLabels.Channel32.Name = i.Channel32.Name
						evi.Proteins[j].TotalLabels.Channel32.CustomName = i.Channel32.CustomName
						evi.Proteins[j].TotalLabels.Channel32.Mz = i.Channel32.Mz
						evi.Proteins[j].TotalLabels.Channel32.Intensity += i.Channel32.Intensity

						//if k.IsNondegenerateEvidence {
						if k.IsUnique {
							//evi.Proteins[j].UniqueLabels = &iso.Labels{}
							evi.Proteins[j].UniqueLabels.Channel1.Name = i.Channel1.Name
							evi.Proteins[j].UniqueLabels.Channel1.CustomName = i.Channel1.CustomName
							evi.Proteins[j].UniqueLabels.Channel1.Mz = i.Channel1.Mz
							evi.Proteins[j].UniqueLabels.Channel1.Intensity += i.Channel1.Intensity

							evi.Proteins[j].UniqueLabels.Channel2.Name = i.Channel2.Name
							evi.Proteins[j].UniqueLabels.Channel2.CustomName = i.Channel2.CustomName
							evi.Proteins[j].UniqueLabels.Channel2.Mz = i.Channel2.Mz
							evi.Proteins[j].UniqueLabels.Channel2.Intensity += i.Channel2.Intensity

							evi.Proteins[j].UniqueLabels.Channel3.Name = i.Channel3.Name
							evi.Proteins[j].UniqueLabels.Channel3.CustomName = i.Channel3.CustomName
							evi.Proteins[j].UniqueLabels.Channel3.Mz = i.Channel3.Mz
							evi.Proteins[j].UniqueLabels.Channel3.Intensity += i.Channel3.Intensity

							evi.Proteins[j].UniqueLabels.Channel4.Name = i.Channel4.Name
							evi.Proteins[j].UniqueLabels.Channel4.CustomName = i.Channel4.CustomName
							evi.Proteins[j].UniqueLabels.Channel4.Mz = i.Channel4.Mz
							evi.Proteins[j].UniqueLabels.Channel4.Intensity += i.Channel4.Intensity

							evi.Proteins[j].UniqueLabels.Channel5.Name = i.Channel5.Name
							evi.Proteins[j].UniqueLabels.Channel5.CustomName = i.Channel5.CustomName
							evi.Proteins[j].UniqueLabels.Channel5.Mz = i.Channel5.Mz
							evi.Proteins[j].UniqueLabels.Channel5.Intensity += i.Channel5.Intensity

							evi.Proteins[j].UniqueLabels.Channel6.Name = i.Channel6.Name
							evi.Proteins[j].UniqueLabels.Channel6.CustomName = i.Channel6.CustomName
							evi.Proteins[j].UniqueLabels.Channel6.Mz = i.Channel6.Mz
							evi.Proteins[j].UniqueLabels.Channel6.Intensity += i.Channel6.Intensity

							evi.Proteins[j].UniqueLabels.Channel7.Name = i.Channel7.Name
							evi.Proteins[j].UniqueLabels.Channel7.CustomName = i.Channel7.CustomName
							evi.Proteins[j].UniqueLabels.Channel7.Mz = i.Channel7.Mz
							evi.Proteins[j].UniqueLabels.Channel7.Intensity += i.Channel7.Intensity

							evi.Proteins[j].UniqueLabels.Channel8.Name = i.Channel8.Name
							evi.Proteins[j].UniqueLabels.Channel8.CustomName = i.Channel8.CustomName
							evi.Proteins[j].UniqueLabels.Channel8.Mz = i.Channel8.Mz
							evi.Proteins[j].UniqueLabels.Channel8.Intensity += i.Channel8.Intensity

							evi.Proteins[j].UniqueLabels.Channel9.Name = i.Channel9.Name
							evi.Proteins[j].UniqueLabels.Channel9.CustomName = i.Channel9.CustomName
							evi.Proteins[j].UniqueLabels.Channel9.Mz = i.Channel9.Mz
							evi.Proteins[j].UniqueLabels.Channel9.Intensity += i.Channel9.Intensity

							evi.Proteins[j].UniqueLabels.Channel10.Name = i.Channel10.Name
							evi.Proteins[j].UniqueLabels.Channel10.CustomName = i.Channel10.CustomName
							evi.Proteins[j].UniqueLabels.Channel10.Mz = i.Channel10.Mz
							evi.Proteins[j].UniqueLabels.Channel10.Intensity += i.Channel10.Intensity

							evi.Proteins[j].UniqueLabels.Channel11.Name = i.Channel11.Name
							evi.Proteins[j].UniqueLabels.Channel11.CustomName = i.Channel11.CustomName
							evi.Proteins[j].UniqueLabels.Channel11.Mz = i.Channel11.Mz
							evi.Proteins[j].UniqueLabels.Channel11.Intensity += i.Channel11.Intensity

							evi.Proteins[j].UniqueLabels.Channel12.Name = i.Channel12.Name
							evi.Proteins[j].UniqueLabels.Channel12.CustomName = i.Channel12.CustomName
							evi.Proteins[j].UniqueLabels.Channel12.Mz = i.Channel12.Mz
							evi.Proteins[j].UniqueLabels.Channel12.Intensity += i.Channel12.Intensity

							evi.Proteins[j].UniqueLabels.Channel13.Name = i.Channel13.Name
							evi.Proteins[j].UniqueLabels.Channel13.CustomName = i.Channel13.CustomName
							evi.Proteins[j].UniqueLabels.Channel13.Mz = i.Channel13.Mz
							evi.Proteins[j].UniqueLabels.Channel13.Intensity += i.Channel13.Intensity

							evi.Proteins[j].UniqueLabels.Channel14.Name = i.Channel14.Name
							evi.Proteins[j].UniqueLabels.Channel14.CustomName = i.Channel14.CustomName
							evi.Proteins[j].UniqueLabels.Channel14.Mz = i.Channel14.Mz
							evi.Proteins[j].UniqueLabels.Channel14.Intensity += i.Channel14.Intensity

							evi.Proteins[j].UniqueLabels.Channel15.Name = i.Channel15.Name
							evi.Proteins[j].UniqueLabels.Channel15.CustomName = i.Channel15.CustomName
							evi.Proteins[j].UniqueLabels.Channel15.Mz = i.Channel15.Mz
							evi.Proteins[j].UniqueLabels.Channel15.Intensity += i.Channel15.Intensity

							evi.Proteins[j].UniqueLabels.Channel16.Name = i.Channel16.Name
							evi.Proteins[j].UniqueLabels.Channel16.CustomName = i.Channel16.CustomName
							evi.Proteins[j].UniqueLabels.Channel16.Mz = i.Channel16.Mz
							evi.Proteins[j].UniqueLabels.Channel16.Intensity += i.Channel16.Intensity

							evi.Proteins[j].UniqueLabels.Channel17.Name = i.Channel17.Name
							evi.Proteins[j].UniqueLabels.Channel17.CustomName = i.Channel17.CustomName
							evi.Proteins[j].UniqueLabels.Channel17.Mz = i.Channel17.Mz
							evi.Proteins[j].UniqueLabels.Channel17.Intensity += i.Channel17.Intensity

							evi.Proteins[j].UniqueLabels.Channel18.Name = i.Channel18.Name
							evi.Proteins[j].UniqueLabels.Channel18.CustomName = i.Channel18.CustomName
							evi.Proteins[j].UniqueLabels.Channel18.Mz = i.Channel18.Mz
							evi.Proteins[j].UniqueLabels.Channel18.Intensity += i.Channel18.Intensity

							evi.Proteins[j].UniqueLabels.Channel19.Name = i.Channel19.Name
							evi.Proteins[j].UniqueLabels.Channel19.CustomName = i.Channel19.CustomName
							evi.Proteins[j].UniqueLabels.Channel19.Mz = i.Channel19.Mz
							evi.Proteins[j].UniqueLabels.Channel19.Intensity += i.Channel19.Intensity

							evi.Proteins[j].UniqueLabels.Channel20.Name = i.Channel20.Name
							evi.Proteins[j].UniqueLabels.Channel20.CustomName = i.Channel20.CustomName
							evi.Proteins[j].UniqueLabels.Channel20.Mz = i.Channel20.Mz
							evi.Proteins[j].UniqueLabels.Channel20.Intensity += i.Channel20.Intensity

							evi.Proteins[j].UniqueLabels.Channel21.Name = i.Channel21.Name
							evi.Proteins[j].UniqueLabels.Channel21.CustomName = i.Channel21.CustomName
							evi.Proteins[j].UniqueLabels.Channel21.Mz = i.Channel21.Mz
							evi.Proteins[j].UniqueLabels.Channel21.Intensity += i.Channel21.Intensity

							evi.Proteins[j].UniqueLabels.Channel22.Name = i.Channel22.Name
							evi.Proteins[j].UniqueLabels.Channel22.CustomName = i.Channel22.CustomName
							evi.Proteins[j].UniqueLabels.Channel22.Mz = i.Channel22.Mz
							evi.Proteins[j].UniqueLabels.Channel22.Intensity += i.Channel22.Intensity

							evi.Proteins[j].UniqueLabels.Channel23.Name = i.Channel23.Name
							evi.Proteins[j].UniqueLabels.Channel23.CustomName = i.Channel23.CustomName
							evi.Proteins[j].UniqueLabels.Channel23.Mz = i.Channel23.Mz
							evi.Proteins[j].UniqueLabels.Channel23.Intensity += i.Channel23.Intensity

							evi.Proteins[j].UniqueLabels.Channel24.Name = i.Channel24.Name
							evi.Proteins[j].UniqueLabels.Channel24.CustomName = i.Channel24.CustomName
							evi.Proteins[j].UniqueLabels.Channel24.Mz = i.Channel24.Mz
							evi.Proteins[j].UniqueLabels.Channel24.Intensity += i.Channel24.Intensity

							evi.Proteins[j].UniqueLabels.Channel25.Name = i.Channel25.Name
							evi.Proteins[j].UniqueLabels.Channel25.CustomName = i.Channel25.CustomName
							evi.Proteins[j].UniqueLabels.Channel25.Mz = i.Channel25.Mz
							evi.Proteins[j].UniqueLabels.Channel25.Intensity += i.Channel25.Intensity

							evi.Proteins[j].UniqueLabels.Channel26.Name = i.Channel26.Name
							evi.Proteins[j].UniqueLabels.Channel26.CustomName = i.Channel26.CustomName
							evi.Proteins[j].UniqueLabels.Channel26.Mz = i.Channel26.Mz
							evi.Proteins[j].UniqueLabels.Channel26.Intensity += i.Channel26.Intensity

							evi.Proteins[j].UniqueLabels.Channel27.Name = i.Channel27.Name
							evi.Proteins[j].UniqueLabels.Channel27.CustomName = i.Channel27.CustomName
							evi.Proteins[j].UniqueLabels.Channel27.Mz = i.Channel27.Mz
							evi.Proteins[j].UniqueLabels.Channel27.Intensity += i.Channel27.Intensity

							evi.Proteins[j].UniqueLabels.Channel28.Name = i.Channel28.Name
							evi.Proteins[j].UniqueLabels.Channel28.CustomName = i.Channel28.CustomName
							evi.Proteins[j].UniqueLabels.Channel28.Mz = i.Channel28.Mz
							evi.Proteins[j].UniqueLabels.Channel28.Intensity += i.Channel28.Intensity

							evi.Proteins[j].UniqueLabels.Channel29.Name = i.Channel29.Name
							evi.Proteins[j].UniqueLabels.Channel29.CustomName = i.Channel29.CustomName
							evi.Proteins[j].UniqueLabels.Channel29.Mz = i.Channel29.Mz
							evi.Proteins[j].UniqueLabels.Channel29.Intensity += i.Channel29.Intensity

							evi.Proteins[j].UniqueLabels.Channel30.Name = i.Channel30.Name
							evi.Proteins[j].UniqueLabels.Channel30.CustomName = i.Channel30.CustomName
							evi.Proteins[j].UniqueLabels.Channel30.Mz = i.Channel30.Mz
							evi.Proteins[j].UniqueLabels.Channel30.Intensity += i.Channel30.Intensity

							evi.Proteins[j].UniqueLabels.Channel31.Name = i.Channel31.Name
							evi.Proteins[j].UniqueLabels.Channel31.CustomName = i.Channel31.CustomName
							evi.Proteins[j].UniqueLabels.Channel31.Mz = i.Channel31.Mz
							evi.Proteins[j].UniqueLabels.Channel31.Intensity += i.Channel31.Intensity

							evi.Proteins[j].UniqueLabels.Channel32.Name = i.Channel32.Name
							evi.Proteins[j].UniqueLabels.Channel32.CustomName = i.Channel32.CustomName
							evi.Proteins[j].UniqueLabels.Channel32.Mz = i.Channel32.Mz
							evi.Proteins[j].UniqueLabels.Channel32.Intensity += i.Channel32.Intensity
						}

						if k.IsURazor {
							//evi.Proteins[j].URazorLabels = &iso.Labels{}
							evi.Proteins[j].URazorLabels.Channel1.Name = i.Channel1.Name
							evi.Proteins[j].URazorLabels.Channel1.CustomName = i.Channel1.CustomName
							evi.Proteins[j].URazorLabels.Channel1.Mz = i.Channel1.Mz
							evi.Proteins[j].URazorLabels.Channel1.Intensity += i.Channel1.Intensity

							evi.Proteins[j].URazorLabels.Channel2.Name = i.Channel2.Name
							evi.Proteins[j].URazorLabels.Channel2.CustomName = i.Channel2.CustomName
							evi.Proteins[j].URazorLabels.Channel2.Mz = i.Channel2.Mz
							evi.Proteins[j].URazorLabels.Channel2.Intensity += i.Channel2.Intensity

							evi.Proteins[j].URazorLabels.Channel3.Name = i.Channel3.Name
							evi.Proteins[j].URazorLabels.Channel3.CustomName = i.Channel3.CustomName
							evi.Proteins[j].URazorLabels.Channel3.Mz = i.Channel3.Mz
							evi.Proteins[j].URazorLabels.Channel3.Intensity += i.Channel3.Intensity

							evi.Proteins[j].URazorLabels.Channel4.Name = i.Channel4.Name
							evi.Proteins[j].URazorLabels.Channel4.CustomName = i.Channel4.CustomName
							evi.Proteins[j].URazorLabels.Channel4.Mz = i.Channel4.Mz
							evi.Proteins[j].URazorLabels.Channel4.Intensity += i.Channel4.Intensity

							evi.Proteins[j].URazorLabels.Channel5.Name = i.Channel5.Name
							evi.Proteins[j].URazorLabels.Channel5.CustomName = i.Channel5.CustomName
							evi.Proteins[j].URazorLabels.Channel5.Mz = i.Channel5.Mz
							evi.Proteins[j].URazorLabels.Channel5.Intensity += i.Channel5.Intensity

							evi.Proteins[j].URazorLabels.Channel6.Name = i.Channel6.Name
							evi.Proteins[j].URazorLabels.Channel6.CustomName = i.Channel6.CustomName
							evi.Proteins[j].URazorLabels.Channel6.Mz = i.Channel6.Mz
							evi.Proteins[j].URazorLabels.Channel6.Intensity += i.Channel6.Intensity

							evi.Proteins[j].URazorLabels.Channel7.Name = i.Channel7.Name
							evi.Proteins[j].URazorLabels.Channel7.CustomName = i.Channel7.CustomName
							evi.Proteins[j].URazorLabels.Channel7.Mz = i.Channel7.Mz
							evi.Proteins[j].URazorLabels.Channel7.Intensity += i.Channel7.Intensity

							evi.Proteins[j].URazorLabels.Channel8.Name = i.Channel8.Name
							evi.Proteins[j].URazorLabels.Channel8.CustomName = i.Channel8.CustomName
							evi.Proteins[j].URazorLabels.Channel8.Mz = i.Channel8.Mz
							evi.Proteins[j].URazorLabels.Channel8.Intensity += i.Channel8.Intensity

							evi.Proteins[j].URazorLabels.Channel9.Name = i.Channel9.Name
							evi.Proteins[j].URazorLabels.Channel9.CustomName = i.Channel9.CustomName
							evi.Proteins[j].URazorLabels.Channel9.Mz = i.Channel9.Mz
							evi.Proteins[j].URazorLabels.Channel9.Intensity += i.Channel9.Intensity

							evi.Proteins[j].URazorLabels.Channel10.Name = i.Channel10.Name
							evi.Proteins[j].URazorLabels.Channel10.CustomName = i.Channel10.CustomName
							evi.Proteins[j].URazorLabels.Channel10.Mz = i.Channel10.Mz
							evi.Proteins[j].URazorLabels.Channel10.Intensity += i.Channel10.Intensity

							evi.Proteins[j].URazorLabels.Channel11.Name = i.Channel11.Name
							evi.Proteins[j].URazorLabels.Channel11.CustomName = i.Channel11.CustomName
							evi.Proteins[j].URazorLabels.Channel11.Mz = i.Channel11.Mz
							evi.Proteins[j].URazorLabels.Channel11.Intensity += i.Channel11.Intensity

							evi.Proteins[j].URazorLabels.Channel12.Name = i.Channel12.Name
							evi.Proteins[j].URazorLabels.Channel12.CustomName = i.Channel12.CustomName
							evi.Proteins[j].URazorLabels.Channel12.Mz = i.Channel12.Mz
							evi.Proteins[j].URazorLabels.Channel12.Intensity += i.Channel12.Intensity

							evi.Proteins[j].URazorLabels.Channel13.Name = i.Channel13.Name
							evi.Proteins[j].URazorLabels.Channel13.CustomName = i.Channel13.CustomName
							evi.Proteins[j].URazorLabels.Channel13.Mz = i.Channel13.Mz
							evi.Proteins[j].URazorLabels.Channel13.Intensity += i.Channel13.Intensity

							evi.Proteins[j].URazorLabels.Channel14.Name = i.Channel14.Name
							evi.Proteins[j].URazorLabels.Channel14.CustomName = i.Channel14.CustomName
							evi.Proteins[j].URazorLabels.Channel14.Mz = i.Channel14.Mz
							evi.Proteins[j].URazorLabels.Channel14.Intensity += i.Channel14.Intensity

							evi.Proteins[j].URazorLabels.Channel15.Name = i.Channel15.Name
							evi.Proteins[j].URazorLabels.Channel15.CustomName = i.Channel15.CustomName
							evi.Proteins[j].URazorLabels.Channel15.Mz = i.Channel15.Mz
							evi.Proteins[j].URazorLabels.Channel15.Intensity += i.Channel15.Intensity

							evi.Proteins[j].URazorLabels.Channel16.Name = i.Channel16.Name
							evi.Proteins[j].URazorLabels.Channel16.CustomName = i.Channel16.CustomName
							evi.Proteins[j].URazorLabels.Channel16.Mz = i.Channel16.Mz
							evi.Proteins[j].URazorLabels.Channel16.Intensity += i.Channel16.Intensity

							evi.Proteins[j].URazorLabels.Channel17.Name = i.Channel17.Name
							evi.Proteins[j].URazorLabels.Channel17.CustomName = i.Channel17.CustomName
							evi.Proteins[j].URazorLabels.Channel17.Mz = i.Channel17.Mz
							evi.Proteins[j].URazorLabels.Channel17.Intensity += i.Channel17.Intensity

							evi.Proteins[j].URazorLabels.Channel18.Name = i.Channel18.Name
							evi.Proteins[j].URazorLabels.Channel18.CustomName = i.Channel18.CustomName
							evi.Proteins[j].URazorLabels.Channel18.Mz = i.Channel18.Mz
							evi.Proteins[j].URazorLabels.Channel18.Intensity += i.Channel18.Intensity

							evi.Proteins[j].URazorLabels.Channel19.Name = i.Channel19.Name
							evi.Proteins[j].URazorLabels.Channel19.CustomName = i.Channel19.CustomName
							evi.Proteins[j].URazorLabels.Channel19.Mz = i.Channel19.Mz
							evi.Proteins[j].URazorLabels.Channel19.Intensity += i.Channel19.Intensity

							evi.Proteins[j].URazorLabels.Channel20.Name = i.Channel20.Name
							evi.Proteins[j].URazorLabels.Channel20.CustomName = i.Channel20.CustomName
							evi.Proteins[j].URazorLabels.Channel20.Mz = i.Channel20.Mz
							evi.Proteins[j].URazorLabels.Channel20.Intensity += i.Channel20.Intensity

							evi.Proteins[j].URazorLabels.Channel21.Name = i.Channel21.Name
							evi.Proteins[j].URazorLabels.Channel21.CustomName = i.Channel21.CustomName
							evi.Proteins[j].URazorLabels.Channel21.Mz = i.Channel21.Mz
							evi.Proteins[j].URazorLabels.Channel21.Intensity += i.Channel21.Intensity

							evi.Proteins[j].URazorLabels.Channel22.Name = i.Channel22.Name
							evi.Proteins[j].URazorLabels.Channel22.CustomName = i.Channel22.CustomName
							evi.Proteins[j].URazorLabels.Channel22.Mz = i.Channel22.Mz
							evi.Proteins[j].URazorLabels.Channel22.Intensity += i.Channel22.Intensity

							evi.Proteins[j].URazorLabels.Channel23.Name = i.Channel23.Name
							evi.Proteins[j].URazorLabels.Channel23.CustomName = i.Channel23.CustomName
							evi.Proteins[j].URazorLabels.Channel23.Mz = i.Channel23.Mz
							evi.Proteins[j].URazorLabels.Channel23.Intensity += i.Channel23.Intensity

							evi.Proteins[j].URazorLabels.Channel24.Name = i.Channel24.Name
							evi.Proteins[j].URazorLabels.Channel24.CustomName = i.Channel24.CustomName
							evi.Proteins[j].URazorLabels.Channel24.Mz = i.Channel24.Mz
							evi.Proteins[j].URazorLabels.Channel24.Intensity += i.Channel24.Intensity

							evi.Proteins[j].URazorLabels.Channel25.Name = i.Channel25.Name
							evi.Proteins[j].URazorLabels.Channel25.CustomName = i.Channel25.CustomName
							evi.Proteins[j].URazorLabels.Channel25.Mz = i.Channel25.Mz
							evi.Proteins[j].URazorLabels.Channel25.Intensity += i.Channel25.Intensity

							evi.Proteins[j].URazorLabels.Channel26.Name = i.Channel26.Name
							evi.Proteins[j].URazorLabels.Channel26.CustomName = i.Channel26.CustomName
							evi.Proteins[j].URazorLabels.Channel26.Mz = i.Channel26.Mz
							evi.Proteins[j].URazorLabels.Channel26.Intensity += i.Channel26.Intensity

							evi.Proteins[j].URazorLabels.Channel27.Name = i.Channel27.Name
							evi.Proteins[j].URazorLabels.Channel27.CustomName = i.Channel27.CustomName
							evi.Proteins[j].URazorLabels.Channel27.Mz = i.Channel27.Mz
							evi.Proteins[j].URazorLabels.Channel27.Intensity += i.Channel27.Intensity

							evi.Proteins[j].URazorLabels.Channel28.Name = i.Channel28.Name
							evi.Proteins[j].URazorLabels.Channel28.CustomName = i.Channel28.CustomName
							evi.Proteins[j].URazorLabels.Channel28.Mz = i.Channel28.Mz
							evi.Proteins[j].URazorLabels.Channel28.Intensity += i.Channel28.Intensity

							evi.Proteins[j].URazorLabels.Channel29.Name = i.Channel29.Name
							evi.Proteins[j].URazorLabels.Channel29.CustomName = i.Channel29.CustomName
							evi.Proteins[j].URazorLabels.Channel29.Mz = i.Channel29.Mz
							evi.Proteins[j].URazorLabels.Channel29.Intensity += i.Channel29.Intensity

							evi.Proteins[j].URazorLabels.Channel30.Name = i.Channel30.Name
							evi.Proteins[j].URazorLabels.Channel30.CustomName = i.Channel30.CustomName
							evi.Proteins[j].URazorLabels.Channel30.Mz = i.Channel30.Mz
							evi.Proteins[j].URazorLabels.Channel30.Intensity += i.Channel30.Intensity

							evi.Proteins[j].URazorLabels.Channel31.Name = i.Channel31.Name
							evi.Proteins[j].URazorLabels.Channel31.CustomName = i.Channel31.CustomName
							evi.Proteins[j].URazorLabels.Channel31.Mz = i.Channel31.Mz
							evi.Proteins[j].URazorLabels.Channel31.Intensity += i.Channel31.Intensity

							evi.Proteins[j].URazorLabels.Channel32.Name = i.Channel32.Name
							evi.Proteins[j].URazorLabels.Channel32.CustomName = i.Channel32.CustomName
							evi.Proteins[j].URazorLabels.Channel32.Mz = i.Channel32.Mz
							evi.Proteins[j].URazorLabels.Channel32.Intensity += i.Channel32.Intensity
						}
					}

					i, ok = phosphoSpectrumMap[l]
					if ok {
						evi.Proteins[j].PhosphoTotalLabels.Channel1.Name = i.Channel1.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel1.CustomName = i.Channel1.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel1.Mz = i.Channel1.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel1.Intensity += i.Channel1.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel2.Name = i.Channel2.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel2.CustomName = i.Channel2.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel2.Mz = i.Channel2.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel2.Intensity += i.Channel2.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel3.Name = i.Channel3.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel3.CustomName = i.Channel3.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel3.Mz = i.Channel3.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel3.Intensity += i.Channel3.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel4.Name = i.Channel4.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel4.CustomName = i.Channel4.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel4.Mz = i.Channel4.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel4.Intensity += i.Channel4.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel5.Name = i.Channel5.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel5.CustomName = i.Channel5.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel5.Mz = i.Channel5.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel5.Intensity += i.Channel5.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel6.Name = i.Channel6.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel6.CustomName = i.Channel6.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel6.Mz = i.Channel6.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel6.Intensity += i.Channel6.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel7.Name = i.Channel7.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel7.CustomName = i.Channel7.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel7.Mz = i.Channel7.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel7.Intensity += i.Channel7.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel8.Name = i.Channel8.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel8.CustomName = i.Channel8.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel8.Mz = i.Channel8.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel8.Intensity += i.Channel8.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel9.Name = i.Channel9.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel9.CustomName = i.Channel9.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel9.Mz = i.Channel9.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel9.Intensity += i.Channel9.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel10.Name = i.Channel10.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel10.CustomName = i.Channel10.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel10.Mz = i.Channel10.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel10.Intensity += i.Channel10.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel11.Name = i.Channel11.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel11.CustomName = i.Channel11.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel11.Mz = i.Channel11.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel11.Intensity += i.Channel11.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel12.Name = i.Channel12.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel12.CustomName = i.Channel12.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel12.Mz = i.Channel12.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel12.Intensity += i.Channel12.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel13.Name = i.Channel13.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel13.CustomName = i.Channel13.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel13.Mz = i.Channel13.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel13.Intensity += i.Channel13.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel14.Name = i.Channel14.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel14.CustomName = i.Channel14.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel14.Mz = i.Channel14.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel14.Intensity += i.Channel14.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel15.Name = i.Channel15.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel15.CustomName = i.Channel15.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel15.Mz = i.Channel15.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel15.Intensity += i.Channel15.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel16.Name = i.Channel16.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel16.CustomName = i.Channel16.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel16.Mz = i.Channel16.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel16.Intensity += i.Channel16.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel17.Name = i.Channel17.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel17.CustomName = i.Channel17.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel17.Mz = i.Channel17.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel17.Intensity += i.Channel17.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel18.Name = i.Channel18.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel18.CustomName = i.Channel18.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel18.Mz = i.Channel18.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel18.Intensity += i.Channel18.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel19.Name = i.Channel19.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel19.CustomName = i.Channel19.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel19.Mz = i.Channel19.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel19.Intensity += i.Channel19.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel20.Name = i.Channel20.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel20.CustomName = i.Channel20.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel20.Mz = i.Channel20.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel20.Intensity += i.Channel20.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel21.Name = i.Channel21.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel21.CustomName = i.Channel21.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel21.Mz = i.Channel21.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel21.Intensity += i.Channel21.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel22.Name = i.Channel22.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel22.CustomName = i.Channel22.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel22.Mz = i.Channel22.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel22.Intensity += i.Channel22.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel23.Name = i.Channel23.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel23.CustomName = i.Channel23.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel23.Mz = i.Channel23.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel23.Intensity += i.Channel23.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel24.Name = i.Channel24.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel24.CustomName = i.Channel24.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel24.Mz = i.Channel24.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel24.Intensity += i.Channel24.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel25.Name = i.Channel25.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel25.CustomName = i.Channel25.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel25.Mz = i.Channel25.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel25.Intensity += i.Channel25.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel26.Name = i.Channel26.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel26.CustomName = i.Channel26.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel26.Mz = i.Channel26.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel26.Intensity += i.Channel26.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel27.Name = i.Channel27.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel27.CustomName = i.Channel27.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel27.Mz = i.Channel27.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel27.Intensity += i.Channel27.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel28.Name = i.Channel28.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel28.CustomName = i.Channel28.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel28.Mz = i.Channel28.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel28.Intensity += i.Channel28.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel29.Name = i.Channel29.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel29.CustomName = i.Channel29.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel29.Mz = i.Channel29.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel29.Intensity += i.Channel29.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel30.Name = i.Channel30.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel30.CustomName = i.Channel30.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel30.Mz = i.Channel30.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel30.Intensity += i.Channel30.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel31.Name = i.Channel31.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel31.CustomName = i.Channel31.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel31.Mz = i.Channel31.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel31.Intensity += i.Channel31.Intensity

						evi.Proteins[j].PhosphoTotalLabels.Channel32.Name = i.Channel32.Name
						evi.Proteins[j].PhosphoTotalLabels.Channel32.CustomName = i.Channel32.CustomName
						evi.Proteins[j].PhosphoTotalLabels.Channel32.Mz = i.Channel32.Mz
						evi.Proteins[j].PhosphoTotalLabels.Channel32.Intensity += i.Channel32.Intensity

						//if k.IsNondegenerateEvidence {
						if k.IsUnique {
							evi.Proteins[j].PhosphoUniqueLabels.Channel1.Name = i.Channel1.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel1.CustomName = i.Channel1.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel1.Mz = i.Channel1.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel1.Intensity += i.Channel1.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel2.Name = i.Channel2.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel2.CustomName = i.Channel2.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel2.Mz = i.Channel2.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel2.Intensity += i.Channel2.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel3.Name = i.Channel3.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel3.CustomName = i.Channel3.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel3.Mz = i.Channel3.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel3.Intensity += i.Channel3.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel4.Name = i.Channel4.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel4.CustomName = i.Channel4.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel4.Mz = i.Channel4.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel4.Intensity += i.Channel4.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel5.Name = i.Channel5.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel5.CustomName = i.Channel5.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel5.Mz = i.Channel5.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel5.Intensity += i.Channel5.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel6.Name = i.Channel6.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel6.CustomName = i.Channel6.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel6.Mz = i.Channel6.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel6.Intensity += i.Channel6.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel7.Name = i.Channel7.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel7.CustomName = i.Channel7.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel7.Mz = i.Channel7.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel7.Intensity += i.Channel7.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel8.Name = i.Channel8.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel8.CustomName = i.Channel8.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel8.Mz = i.Channel8.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel8.Intensity += i.Channel8.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel9.Name = i.Channel9.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel9.CustomName = i.Channel9.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel9.Mz = i.Channel9.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel9.Intensity += i.Channel9.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel10.Name = i.Channel10.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel10.CustomName = i.Channel10.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel10.Mz = i.Channel10.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel10.Intensity += i.Channel10.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel11.Name = i.Channel11.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel11.CustomName = i.Channel11.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel11.Mz = i.Channel11.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel11.Intensity += i.Channel11.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel12.Name = i.Channel12.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel12.CustomName = i.Channel12.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel12.Mz = i.Channel12.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel12.Intensity += i.Channel12.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel13.Name = i.Channel13.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel13.CustomName = i.Channel13.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel13.Mz = i.Channel13.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel13.Intensity += i.Channel13.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel14.Name = i.Channel14.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel14.CustomName = i.Channel14.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel14.Mz = i.Channel14.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel14.Intensity += i.Channel14.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel15.Name = i.Channel15.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel15.CustomName = i.Channel15.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel15.Mz = i.Channel15.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel15.Intensity += i.Channel15.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel16.Name = i.Channel16.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel16.CustomName = i.Channel16.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel16.Mz = i.Channel16.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel16.Intensity += i.Channel16.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel17.Name = i.Channel17.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel17.CustomName = i.Channel17.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel17.Mz = i.Channel17.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel17.Intensity += i.Channel17.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel18.Name = i.Channel18.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel18.CustomName = i.Channel18.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel18.Mz = i.Channel18.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel18.Intensity += i.Channel18.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel19.Name = i.Channel19.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel19.CustomName = i.Channel19.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel19.Mz = i.Channel19.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel19.Intensity += i.Channel19.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel20.Name = i.Channel20.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel20.CustomName = i.Channel20.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel20.Mz = i.Channel20.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel20.Intensity += i.Channel20.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel21.Name = i.Channel21.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel21.CustomName = i.Channel21.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel21.Mz = i.Channel21.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel21.Intensity += i.Channel21.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel22.Name = i.Channel22.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel22.CustomName = i.Channel22.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel22.Mz = i.Channel22.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel22.Intensity += i.Channel22.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel23.Name = i.Channel23.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel23.CustomName = i.Channel23.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel23.Mz = i.Channel23.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel23.Intensity += i.Channel23.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel24.Name = i.Channel24.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel24.CustomName = i.Channel24.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel24.Mz = i.Channel24.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel24.Intensity += i.Channel24.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel25.Name = i.Channel25.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel25.CustomName = i.Channel25.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel25.Mz = i.Channel25.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel25.Intensity += i.Channel25.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel26.Name = i.Channel26.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel26.CustomName = i.Channel26.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel26.Mz = i.Channel26.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel26.Intensity += i.Channel26.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel27.Name = i.Channel27.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel27.CustomName = i.Channel27.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel27.Mz = i.Channel27.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel27.Intensity += i.Channel27.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel28.Name = i.Channel28.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel28.CustomName = i.Channel28.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel28.Mz = i.Channel28.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel28.Intensity += i.Channel28.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel29.Name = i.Channel29.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel29.CustomName = i.Channel29.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel29.Mz = i.Channel29.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel29.Intensity += i.Channel29.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel30.Name = i.Channel30.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel30.CustomName = i.Channel30.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel30.Mz = i.Channel30.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel30.Intensity += i.Channel30.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel31.Name = i.Channel31.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel31.CustomName = i.Channel31.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel31.Mz = i.Channel31.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel31.Intensity += i.Channel31.Intensity

							evi.Proteins[j].PhosphoUniqueLabels.Channel32.Name = i.Channel32.Name
							evi.Proteins[j].PhosphoUniqueLabels.Channel32.CustomName = i.Channel32.CustomName
							evi.Proteins[j].PhosphoUniqueLabels.Channel32.Mz = i.Channel32.Mz
							evi.Proteins[j].PhosphoUniqueLabels.Channel32.Intensity += i.Channel32.Intensity
						}

						if k.IsURazor {
							evi.Proteins[j].PhosphoURazorLabels.Channel1.Name = i.Channel1.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel1.CustomName = i.Channel1.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel1.Mz = i.Channel1.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel1.Intensity += i.Channel1.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel2.Name = i.Channel2.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel2.CustomName = i.Channel2.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel2.Mz = i.Channel2.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel2.Intensity += i.Channel2.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel3.Name = i.Channel3.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel3.CustomName = i.Channel3.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel3.Mz = i.Channel3.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel3.Intensity += i.Channel3.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel4.Name = i.Channel4.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel4.CustomName = i.Channel4.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel4.Mz = i.Channel4.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel4.Intensity += i.Channel4.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel5.Name = i.Channel5.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel5.CustomName = i.Channel5.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel5.Mz = i.Channel5.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel5.Intensity += i.Channel5.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel6.Name = i.Channel6.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel6.CustomName = i.Channel6.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel6.Mz = i.Channel6.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel6.Intensity += i.Channel6.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel7.Name = i.Channel7.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel7.CustomName = i.Channel7.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel7.Mz = i.Channel7.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel7.Intensity += i.Channel7.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel8.Name = i.Channel8.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel8.CustomName = i.Channel8.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel8.Mz = i.Channel8.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel8.Intensity += i.Channel8.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel9.Name = i.Channel9.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel9.CustomName = i.Channel9.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel9.Mz = i.Channel9.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel9.Intensity += i.Channel9.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel10.Name = i.Channel10.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel10.CustomName = i.Channel10.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel10.Mz = i.Channel10.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel10.Intensity += i.Channel10.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel11.Name = i.Channel11.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel11.CustomName = i.Channel11.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel11.Mz = i.Channel11.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel11.Intensity += i.Channel11.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel12.Name = i.Channel12.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel12.CustomName = i.Channel12.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel12.Mz = i.Channel12.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel12.Intensity += i.Channel12.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel13.Name = i.Channel13.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel13.CustomName = i.Channel13.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel13.Mz = i.Channel13.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel13.Intensity += i.Channel13.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel14.Name = i.Channel14.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel14.CustomName = i.Channel14.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel14.Mz = i.Channel14.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel14.Intensity += i.Channel14.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel15.Name = i.Channel15.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel15.CustomName = i.Channel15.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel15.Mz = i.Channel15.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel15.Intensity += i.Channel15.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel16.Name = i.Channel16.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel16.CustomName = i.Channel16.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel16.Mz = i.Channel16.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel16.Intensity += i.Channel16.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel17.Name = i.Channel17.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel17.CustomName = i.Channel17.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel17.Mz = i.Channel17.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel17.Intensity += i.Channel17.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel18.Name = i.Channel18.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel18.CustomName = i.Channel18.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel18.Mz = i.Channel18.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel18.Intensity += i.Channel18.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel19.Name = i.Channel19.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel19.CustomName = i.Channel19.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel19.Mz = i.Channel19.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel19.Intensity += i.Channel19.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel20.Name = i.Channel20.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel20.CustomName = i.Channel20.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel20.Mz = i.Channel20.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel20.Intensity += i.Channel20.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel21.Name = i.Channel21.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel21.CustomName = i.Channel21.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel21.Mz = i.Channel21.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel21.Intensity += i.Channel21.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel22.Name = i.Channel22.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel22.CustomName = i.Channel22.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel22.Mz = i.Channel22.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel22.Intensity += i.Channel22.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel23.Name = i.Channel23.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel23.CustomName = i.Channel23.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel23.Mz = i.Channel23.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel23.Intensity += i.Channel23.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel24.Name = i.Channel24.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel24.CustomName = i.Channel24.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel24.Mz = i.Channel24.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel24.Intensity += i.Channel24.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel25.Name = i.Channel25.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel25.CustomName = i.Channel25.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel25.Mz = i.Channel25.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel25.Intensity += i.Channel25.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel26.Name = i.Channel26.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel26.CustomName = i.Channel26.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel26.Mz = i.Channel26.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel26.Intensity += i.Channel26.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel27.Name = i.Channel27.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel27.CustomName = i.Channel27.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel27.Mz = i.Channel27.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel27.Intensity += i.Channel27.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel28.Name = i.Channel28.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel28.CustomName = i.Channel28.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel28.Mz = i.Channel28.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel28.Intensity += i.Channel28.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel29.Name = i.Channel29.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel29.CustomName = i.Channel29.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel29.Mz = i.Channel29.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel29.Intensity += i.Channel29.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel30.Name = i.Channel30.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel30.CustomName = i.Channel30.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel30.Mz = i.Channel30.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel30.Intensity += i.Channel30.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel31.Name = i.Channel31.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel31.CustomName = i.Channel31.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel31.Mz = i.Channel31.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel31.Intensity += i.Channel31.Intensity

							evi.Proteins[j].PhosphoURazorLabels.Channel32.Name = i.Channel32.Name
							evi.Proteins[j].PhosphoURazorLabels.Channel32.CustomName = i.Channel32.CustomName
							evi.Proteins[j].PhosphoURazorLabels.Channel32.Mz = i.Channel32.Mz
							evi.Proteins[j].PhosphoURazorLabels.Channel32.Intensity += i.Channel32.Intensity
						}
					}
				}
			}
		}
	}

	return evi
}

// NormToTotalProteins calculates the protein level normalization based on total proteins
func NormToTotalProteins(evi rep.Evidence) rep.Evidence {

	var topValue float64
	var channelSum = [32]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var normFactors = [32]float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	// sum TMT singal for each column
	for _, i := range evi.Proteins {
		channelSum[0] += i.URazorLabels.Channel1.Intensity
		channelSum[1] += i.URazorLabels.Channel2.Intensity
		channelSum[2] += i.URazorLabels.Channel3.Intensity
		channelSum[3] += i.URazorLabels.Channel4.Intensity
		channelSum[4] += i.URazorLabels.Channel5.Intensity
		channelSum[5] += i.URazorLabels.Channel6.Intensity
		channelSum[6] += i.URazorLabels.Channel7.Intensity
		channelSum[7] += i.URazorLabels.Channel8.Intensity
		channelSum[8] += i.URazorLabels.Channel9.Intensity
		channelSum[9] += i.URazorLabels.Channel10.Intensity
		channelSum[10] += i.URazorLabels.Channel11.Intensity
		channelSum[11] += i.URazorLabels.Channel12.Intensity
		channelSum[12] += i.URazorLabels.Channel13.Intensity
		channelSum[13] += i.URazorLabels.Channel14.Intensity
		channelSum[14] += i.URazorLabels.Channel15.Intensity
		channelSum[15] += i.URazorLabels.Channel16.Intensity
		channelSum[16] += i.URazorLabels.Channel17.Intensity
		channelSum[17] += i.URazorLabels.Channel18.Intensity
		channelSum[18] += i.URazorLabels.Channel19.Intensity
		channelSum[19] += i.URazorLabels.Channel20.Intensity
		channelSum[20] += i.URazorLabels.Channel21.Intensity
		channelSum[21] += i.URazorLabels.Channel22.Intensity
		channelSum[22] += i.URazorLabels.Channel23.Intensity
		channelSum[23] += i.URazorLabels.Channel24.Intensity
		channelSum[24] += i.URazorLabels.Channel25.Intensity
		channelSum[25] += i.URazorLabels.Channel26.Intensity
		channelSum[26] += i.URazorLabels.Channel27.Intensity
		channelSum[27] += i.URazorLabels.Channel28.Intensity
		channelSum[28] += i.URazorLabels.Channel29.Intensity
		channelSum[29] += i.URazorLabels.Channel30.Intensity
		channelSum[30] += i.URazorLabels.Channel31.Intensity
		channelSum[31] += i.URazorLabels.Channel32.Intensity
	}

	// find the highest value amongst channels
	for _, i := range channelSum {
		if i > topValue {
			topValue = i
		}
	}

	// calculate normalizing factors
	for i := range channelSum {
		normFactors[i] = channelSum[i] / topValue
	}

	// multiply each protein TMT set by the factors to get normalized values
	for _, i := range evi.Proteins {
		i.URazorLabels.Channel1.Intensity *= normFactors[0]
		i.URazorLabels.Channel2.Intensity *= normFactors[1]
		i.URazorLabels.Channel3.Intensity *= normFactors[2]
		i.URazorLabels.Channel4.Intensity *= normFactors[3]
		i.URazorLabels.Channel5.Intensity *= normFactors[4]
		i.URazorLabels.Channel6.Intensity *= normFactors[5]
		i.URazorLabels.Channel7.Intensity *= normFactors[6]
		i.URazorLabels.Channel8.Intensity *= normFactors[7]
		i.URazorLabels.Channel9.Intensity *= normFactors[8]
		i.URazorLabels.Channel10.Intensity *= normFactors[9]
		i.URazorLabels.Channel11.Intensity *= normFactors[10]
		i.URazorLabels.Channel12.Intensity *= normFactors[11]
		i.URazorLabels.Channel13.Intensity *= normFactors[12]
		i.URazorLabels.Channel14.Intensity *= normFactors[13]
		i.URazorLabels.Channel15.Intensity *= normFactors[14]
		i.URazorLabels.Channel16.Intensity *= normFactors[15]
		i.URazorLabels.Channel17.Intensity *= normFactors[16]
		i.URazorLabels.Channel18.Intensity *= normFactors[17]
		i.URazorLabels.Channel19.Intensity *= normFactors[18]
		i.URazorLabels.Channel20.Intensity *= normFactors[19]
		i.URazorLabels.Channel21.Intensity *= normFactors[20]
		i.URazorLabels.Channel22.Intensity *= normFactors[21]
		i.URazorLabels.Channel23.Intensity *= normFactors[22]
		i.URazorLabels.Channel24.Intensity *= normFactors[23]
		i.URazorLabels.Channel25.Intensity *= normFactors[24]
		i.URazorLabels.Channel26.Intensity *= normFactors[25]
		i.URazorLabels.Channel27.Intensity *= normFactors[26]
		i.URazorLabels.Channel28.Intensity *= normFactors[27]
		i.URazorLabels.Channel29.Intensity *= normFactors[28]
		i.URazorLabels.Channel30.Intensity *= normFactors[29]
		i.URazorLabels.Channel31.Intensity *= normFactors[30]
		i.URazorLabels.Channel32.Intensity *= normFactors[31]
	}

	return evi
}
