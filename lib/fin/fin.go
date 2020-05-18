package fin

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unicode/utf16"

	"philosopher/lib/msg"
)

// interface shared by all data objects in the raw file
type reader interface {
	Read(io.Reader, Version)
}

////

// RawFileInfo contains the addresses of the different RunHeaders,
// (header of the data that each connected instrument produced)
// also the acquisition date
type RawFileInfo struct {
	Preamble InfoPreamble
	Heading1 PascalString
	Heading2 PascalString
	Heading3 PascalString
	Heading4 PascalString
	Heading5 PascalString
	Unknown1 PascalString
}

// InfoPreamble struct
type InfoPreamble struct {
	Methodfilepresent uint32
	Year              uint16
	Month             uint16
	Weekday           uint16
	Day               uint16
	Hour              uint16
	Minute            uint16
	Second            uint16
	Millisecond       uint16

	Unknown1        uint32
	DataAddr32      uint32
	NControllers    uint32
	NControllers2   uint32
	Unknown2        uint32
	Unknown3        uint32
	RunHeaderAddr32 []uint32
	Unknown4        []uint32
	Unknown5        []uint32
	Padding1        [764]byte //760 bytes, 756 bytes in v57

	DataAddr      uint64
	Unknown6      uint64
	RunHeaderAddr []uint64
	Unknown7      []uint64
	Padding2      [1032]byte //1024 bytes, 1008 bytes in v64
}

// PascalString struct
type PascalString struct {
	Length int32
	Text   []uint16
}

// AutoSamplerInfo comes from the sampling device
type AutoSamplerInfo struct {
	Preamble AutoSamplerPreamble
	Text     PascalString
}

// AutoSamplerPreamble struct
type AutoSamplerPreamble struct {
	Unknown1      uint32
	Unknown2      uint32
	NumberOfWells uint32
	Unknown3      uint32
	Unknown4      uint32
	Unknown15     uint32
}

// SequencerRow contains more information about what the autosampler did
type SequencerRow struct {
	Injection  InjectionData
	Unknown1   PascalString
	Unknown2   PascalString
	ID         PascalString
	Comment    PascalString
	Userlabel1 PascalString
	Userlabel2 PascalString
	Userlabel3 PascalString
	Userlabel4 PascalString
	Userlabel5 PascalString
	Instmethod PascalString
	Procmethod PascalString
	Filename   PascalString
	Path       PascalString

	Vial     PascalString
	Unknown3 PascalString
	Unknown4 PascalString
	Unknown5 uint32

	Unknown6  PascalString
	Unknown7  PascalString
	Unknown8  PascalString
	Unknown9  PascalString
	Unknown10 PascalString
	Unknown11 PascalString
	Unknown12 PascalString
	Unknown13 PascalString
	Unknown14 PascalString
	Unknown15 PascalString
	Unknown16 PascalString
	Unknown17 PascalString
	Unknown18 PascalString
	Unknown19 PascalString
	Unknown20 PascalString
}

// InjectionData struct
type InjectionData struct {
	Unknown1                    uint32
	Rownumber                   uint32
	Unknown2                    uint32
	Vial                        [6]uint16 //utf-16
	Injectionvolume             float64
	SampleWeight                float64
	SampleVolume                float64
	InternationalStandardAmount float64
	Dilutionfactor              float64
}

// RunHeader contain all data addresses for data that a certain machine
// connected to the Mass Spectrometer (including the MS itself)
// has acquired. Also SN data is available
type RunHeader struct {
	Address           uint64
	SampleInfo        SampleInfo
	Filename1         filename
	Filename2         filename
	Filename3         filename
	Filename4         filename
	Filename5         filename
	Filename6         filename
	Unknown1          float64
	Unknown2          float64
	Filename7         filename
	Filename8         filename
	Filename9         filename
	FilenameA         filename
	FilenameB         filename
	FilenameC         filename
	FilenameD         filename
	ScantrailerAddr32 uint32
	ScanparamsAddr32  uint32
	Unknown3          uint32
	Unknown4          uint32
	Nsegs             uint32
	Unknown5          uint32
	Unknown6          uint32
	OwnAddr32         uint32
	Unknown7          uint32
	Unknown8          uint32
	ScanindexAddr     uint64
	DataAddr          uint64
	InstlogAddr       uint64
	ErrorlogAddr      uint64
	Unknown9          uint64
	ScantrailerAddr   uint64
	ScanparamsAddr    uint64
	Unknown10         uint32
	OwnAddr           uint64
	Unknown11         uint32
	Unknown12         uint32
	Unknown13         uint32
	Unknown14         uint32
	Unknown15         uint32
	Unknown16         uint32
	Unknown17         uint32
	Unknown18         uint32
	Unknown19         uint32
	Unknown20         uint32
	Unknown21         uint32
	Unknown22         uint32
	Unknown23         uint32
	Unknown24         uint32
	Unknown25         uint32
	Unknown26         uint32
	Unknown27         uint32
	Unknown28         uint32
	Unknown29         uint32
	Unknown30         uint32
	Unknown31         uint32
	Unknown32         uint32
	Unknown33         uint32
	Unknown34         uint32
}

//SampleInfo contains some other info
type SampleInfo struct {
	Unknown1        uint32
	Unknown2        uint32
	FirstScanNumber uint32
	LastScanNumber  uint32
	InstlogLength   uint32
	Unknown3        uint32
	Unknown4        uint32
	ScanindexAddr   uint32 //unused in 64-bit versions
	DataAddr        uint32
	InstlogAddr     uint32
	ErrorlogAddr    uint32
	Unknown5        uint32
	MaxSignal       float64
	Lowmz           float64
	Highmz          float64
	Starttime       float64
	Endtime         float64
	Unknown6        [56]byte
	Tag1            [44]uint16
	Tag2            [20]uint16
	Tag3            [160]uint16
}

type filename [260]uint16

// FileHeader is the Thermo fileheaders most valuable piece of info is the file version.
// It determines the reading strategy for some data structures that changed over time
// 1356 bytes
type FileHeader struct {
	Magic      uint16    // 2 bytes
	Signature  signature // 18 bytes
	Unknown1   uint32    // 4 bytes
	Unknown2   uint32    // 4 bytes
	Unknown3   uint32    // 4 bytes
	Unknown4   uint32    // 4 bytes
	Version    Version   // 4 bytes
	AuditStart AuditTag  // 112 bytes
	AuditEnd   AuditTag  // 112 bytes
	Unknown5   uint32    // 4 bytes
	Unknown6   [60]byte  // 60 bytes
	Tag        headertag // 1028 bytes
}

// AuditTag : 112 bytes
type AuditTag struct {
	Time     uint64   // 8 bytes Windows 64-bit timestamp
	Tag1     audittag // 50 bytes
	Tag2     audittag
	Unknown1 uint32 // 4 bytes
}

type audittag [25]uint16

type headertag [514]uint16

type signature [9]uint16

// Version ...
type Version uint32

// InstID : 142 bytes
type InstID struct {
	Address  uint64
	Unknown1 [8]byte      // 8 bytes
	Unknown2 uint32       // 4 bytes
	Model1   PascalString // 40 bytes
	Model2   PascalString // 40 bytes
	Serial   PascalString // 20 bytes
	Version  PascalString // 14 bytes
	Tag1     PascalString // 4 bytes
	Tag2     PascalString // 4 bytes
	Tag3     PascalString // 4 bytes
	Tag4     PascalString // 4 bytes
}

// Profile is the structure containing the profile-mode points
type Profile struct {
	FirstValue float64
	Step       float64
	PeakCount  uint32
	Nbins      uint32
	Chunks     []ProfileChunk
}

// ProfileChunk points are collected in chunks with adjacent signal points
type ProfileChunk struct {
	Firstbin uint32
	Nbins    uint32
	Fudge    float32
	Signal   []float32
}

// ScanDataPacket is an MS scan packet, containing Centroid Peak or Profile intensities
type ScanDataPacket struct {
	Header         PacketHeader
	Profile        Profile
	PeakList       PeakList
	DescriptorList []PeakDescriptor
	Unknown        []float32
	Triplets       []float32
}

//PeakDescriptor is a struct containing more info about the peaks
type PeakDescriptor struct {
	Index  uint16
	Flags  uint8
	Charge uint8
}

// PeakList is the data structure holding the peaks
type PeakList struct {
	Count uint32
	Peaks []CentroidedPeak
}

// CentroidedPeak ia a peak itself
type CentroidedPeak struct {
	Mz        float32
	Abundance float32
}

// PacketHeader is a Header containing info about how many peaks/profile points were registered
type PacketHeader struct {
	Unknown1           uint32
	ProfileSize        uint32
	PeaklistSize       uint32
	Layout             uint32
	DescriptorListSize uint32
	UnknownStreamSize  uint32
	TripletStreamSize  uint32
	Unknown2           uint32
	Lowmz              float32
	Highmz             float32
}

// ScanDataPackets ...
type ScanDataPackets []ScanDataPacket

// TrailerLength : I currently have no idea what TrailerLength is
type TrailerLength uint32

// ScanEvent are encoded headers of the MS scans, their Preamble
// contain the MS level, type of ionization etc. Events themselves contain range, and
// conversion parameters from Hz to m/z
type ScanEvent struct {
	Preamble [132]uint8 //128 bytes from v63 on, 120 in v62, 80 in v57, 41 below that
	//Preamble[6] == ms-level
	//Preamble[40] == analyzer

	Nprecursors uint32
	Reaction    []Reaction
	Unknown1    [13]uint32
	MZrange     [3]FractionCollector
	Nparam      uint32
	Unknown2    [4]float64
	A           float64
	B           float64
	C           float64
}

// Reaction ...
type Reaction struct {
	Precursormz float64
	Unknown1    float64
	Energy      float64
	Unknown2    uint32
	Unknown3    uint32
}

// FractionCollector ...
type FractionCollector struct {
	Lowmz  float64
	Highmz float64
}

// ScanEvents ...
type ScanEvents []ScanEvent

// ScanIndexEntry is the scan index entries are a list of pointers to the scans
// other important information is the scan time
type ScanIndexEntry struct {
	Offset32       uint32
	Index          uint32
	Scanevent      uint16
	Scansegment    uint16
	Next           uint32
	Unknown1       uint32
	DataPacketSize uint32
	Time           float64
	Totalcurrent   float64
	Baseintensity  float64
	Basemz         float64
	Lowmz          float64
	Highmz         float64
	Offset         uint64
	Unknown2       uint32
	Unknown3       uint32
}

// ScanIndex ...
type ScanIndex []ScanIndexEntry

// CIndexEntry entries for Chromatography data
type CIndexEntry struct {
	Offset32 uint32
	Index    uint32
	Event    uint16
	Unknown1 uint16
	Unknown2 uint32
	Unknown3 uint32
	Unknown4 uint32
	Unknown5 float64
	Time     float64
	Unknown6 float64
	Unknown7 float64
	Value    float64

	Offset uint64
}

// CIndexEntries ...
type CIndexEntries []CIndexEntry

// CDataPacket are the data from Chromatography machines
type CDataPacket struct { //16 bytes
	Value float64
	Time  float64
}

// CDataPackets ...
type CDataPackets []CDataPacket

// Peak represents an ion peak
type Peak struct {
	Mz float64
	I  float32
}

// A Spectrum is a collection of peaks
type Spectrum []Peak

// Spectrum implements sort.Interface for []Peak based on m/z
func (a Spectrum) Len() int           { return len(a) }
func (a Spectrum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Spectrum) Less(i, j int) bool { return a[i].Mz < a[j].Mz }

// Scan represents the peak acquisition event of the mass spectrometer
type Scan struct {
	Analyzer      Analyzer
	MSLevel       uint8
	Polarity      []byte
	Mode          []byte
	Type          []byte
	LowMz         float64
	HighMz        float64
	BaseMz        float64
	BaseIntensity float64
	TotalCurrent  float64
	// Spectrum is a function forcing the read of a spectrum,
	// which is "delayed" for efficiency reasons. If it was not delayed
	// and Spectrum were a data structure, it would always have to
	// be read, which is very expensive. Now if only another property of
	// Scan (cheaper to obtain) is requested, resources are saved.
	Spectrum func(centroided ...bool) Spectrum
	// PrecursorMzs is only filled with mz values at MSx scans.
	//PrecursorMzs []float64
	Fragment []Fragment
	Time     float64
}

// Analyzer is the mass analyzer
type Analyzer int

// Fragment represents an MSn scan
type Fragment struct {
	// PrecursorMzs is only filled with mz values at MSx scans.
	PrecursorMzs   []float64
	ColisionEnergy []float64
}

////

// RawData is a top level representantio of a Thermo Finnigan file
type RawData struct {
	File            *os.File
	FileName        string
	Version         Version
	ScanCount       uint64
	StartTime       float64
	EndTime         float64
	Manufacturer    string
	Model           string
	SoftwareName    string
	SoftwareVersion string
	Ionization      []string
	Analyzer        []string
	Detector        []string
	Scanevents      ScanEvents
	Scanindex       ScanIndex
}

// ProcessRaw calls other low level functions and fill out RawData struct
func (rd *RawData) ProcessRaw(f string) {

	file, e := os.Open(f)
	if e != nil {
		msg.ReadFile(e, "fatal")
	}

	// Read headers for file version and RunHeader addresses.
	info, ver := readHeaders(file)

	// RunHeader
	var rh RunHeader
	rh.Retrieve(file, info, ver)

	// InstrumentID
	var inst InstID
	inst.Retrieve(file, rh.Address, ver)

	// For later conversion of frequency values to m/z, we need a ScanEvent for each Scan.
	// The list of them starts an uint32 later than ScantrailerAddr
	nScans := uint64(rh.SampleInfo.LastScanNumber - rh.SampleInfo.FirstScanNumber + 1)

	scanevents := make(ScanEvents, nScans)
	readBetween(file, rh.ScantrailerAddr+4, rh.ScanparamsAddr, ver, scanevents)

	// read all scanindexentries at once
	scanindex := make(ScanIndex, nScans)
	readBetween(file, rh.ScanindexAddr, rh.ScantrailerAddr, ver, scanindex)

	// make the offsets absolute in the file instead of relative to the data address
	for i := range scanindex {
		scanindex[i].Offset += rh.DataAddr
	}

	rd.File = file
	rd.FileName = f
	rd.Version = ver
	rd.ScanCount = nScans
	rd.StartTime = rh.SampleInfo.Starttime
	rd.EndTime = rh.SampleInfo.Endtime
	rd.Manufacturer = "Thermo Scientific"
	rd.Model = fmt.Sprintf("%s", inst.Model2)
	rd.Ionization = ScanEventData(scanevents, "ionization")
	rd.Analyzer = ScanEventData(scanevents, "analyzer")
	rd.Detector = ScanEventData(scanevents, "detector")
	rd.SoftwareName = "XCalibur"
	rd.SoftwareVersion = fmt.Sprintf("%s", inst.Tag1)
	rd.Scanevents = scanevents
	rd.Scanindex = scanindex

	return
}

// ScanEventData ...
func ScanEventData(se ScanEvents, info string) (list []string) {

	var infoMap = make(map[uint8]string)
	var inRaw = make(map[string]uint8)

	if strings.EqualFold(info, "analyzer") {

		infoMap[0] = "ITMS"
		infoMap[1] = "TQMS"
		infoMap[2] = "SQMS"
		infoMap[3] = "TOFMS"
		infoMap[4] = "FTMS"
		infoMap[5] = "Sector"
		infoMap[6] = "undefined"

		for i := range se {
			v, ok := infoMap[se[i].Preamble[40]]
			if ok {
				_, okInRaw := inRaw[v]
				if okInRaw {
					inRaw[v] = 0
				} else {
					inRaw[v] = 0
					list = append(list, v)
				}
			}
		}

	} else if strings.EqualFold(info, "detector") {

		infoMap[0] = "Valid"
		infoMap[1] = "Undefined"

		for i := range se {
			v, ok := infoMap[se[i].Preamble[3]]
			if ok {
				_, okInRaw := inRaw[v]
				if okInRaw {
					inRaw[v] = 0
				} else {
					inRaw[v] = 0
					list = append(list, v)
				}
			}
		}

	} else if strings.EqualFold(info, "ionization") {

		infoMap[0] = "EI"
		infoMap[1] = "CI"
		infoMap[2] = "FABI"
		infoMap[3] = "ESI"
		infoMap[4] = "APCI"
		infoMap[5] = "NSI"
		infoMap[6] = "TSI"
		infoMap[7] = "FDI"
		infoMap[8] = "MALDI"
		infoMap[9] = "GDI"
		infoMap[10] = "Undefined"

		for i := range se {
			v, ok := infoMap[se[i].Preamble[11]]
			if ok {
				_, okInRaw := inRaw[v]
				if okInRaw {
					inRaw[v] = 0
				} else {
					inRaw[v] = 0
					list = append(list, v)
				}
			}
		}

	} else if strings.EqualFold(info, "scantype") {

		infoMap[0] = "Full"
		infoMap[1] = "Zoom"
		infoMap[2] = "SIM"
		infoMap[3] = "SRM"
		infoMap[4] = "CRM"
		infoMap[5] = "undefined"
		infoMap[6] = "Q1"
		infoMap[7] = "Q3"

		for i := range se {
			v, ok := infoMap[se[i].Preamble[7]]
			if ok {
				list = append(list, v)
			}
		}

	}

	return list
}

// Close closes the RAW file
func (rd *RawData) Close() error {
	return rd.File.Close()
}

// AllScans is a convenience function that runs over all spectra in the raw file
// On every encountered MS Scan, the function fun is called
func (rd *RawData) AllScans(fun func(scan Scan)) {

	for i := 1; i <= rd.NScans(); i++ {
		fun(rd.Scan(i))
	}

	return
}

// Scan returns the scan at the scan number in argument
func (rd *RawData) Scan(sn int) (scan Scan) {

	if sn < 1 || sn > rd.NScans() {
		log.Print("Scan Number ", sn, " is out of bounds [1, ", rd.NScans(), "]")
		return
	}

	var polMap = make(map[uint8][]byte)
	polMap[0] = []byte("-") // "Negative"
	polMap[1] = []byte("+") // "Positive"
	//polMap[2] = []byte("Undefined")

	var modeMap = make(map[uint8][]byte)
	modeMap[0] = []byte("Centroided")
	modeMap[1] = []byte("Profile")
	//polMap[2] = []byte("Undefined")

	var typeMap = make(map[uint8][]byte)
	typeMap[0] = []byte("Full")
	typeMap[1] = []byte("Zoom")
	typeMap[2] = []byte("SIM")
	typeMap[3] = []byte("SRM")
	typeMap[4] = []byte("CRM")
	typeMap[5] = []byte("Unknown")
	typeMap[6] = []byte("Q1")
	typeMap[7] = []byte("Q3")

	scan.Time = rd.Scanindex[sn-1].Time
	scan.MSLevel = rd.Scanevents[sn-1].Preamble[6]
	scan.Polarity = polMap[rd.Scanevents[sn-1].Preamble[4]]
	scan.Type = typeMap[rd.Scanevents[sn-1].Preamble[7]]
	scan.Mode = modeMap[rd.Scanevents[sn-1].Preamble[5]]
	scan.Analyzer = Analyzer(rd.Scanevents[sn-1].Preamble[40])
	scan.LowMz = rd.Scanindex[sn-1].Lowmz
	scan.HighMz = rd.Scanindex[sn-1].Highmz
	scan.BaseMz = rd.Scanindex[sn-1].Basemz
	scan.BaseIntensity = rd.Scanindex[sn-1].Baseintensity
	scan.TotalCurrent = rd.Scanindex[sn-1].Totalcurrent

	for j := range rd.Scanevents[sn-1].Reaction {
		var frag Fragment
		frag.PrecursorMzs = make([]float64, len(rd.Scanevents[sn-1].Reaction))
		frag.ColisionEnergy = make([]float64, len(rd.Scanevents[sn-1].Reaction))
		frag.PrecursorMzs[j] = rd.Scanevents[sn-1].Reaction[j].Precursormz
		frag.ColisionEnergy[j] = rd.Scanevents[sn-1].Reaction[j].Energy
		scan.Fragment = append(scan.Fragment, frag)
	}

	//scan.Fragment = frag

	// scan.PrecursorMzs = make([]float64, len(rd.Scanevents[sn-1].Reaction))
	// for j := range rd.Scanevents[sn-1].Reaction {
	// 	scan.PrecursorMzs[j] = rd.Scanevents[sn-1].Reaction[j].Precursormz
	// }

	scan.Spectrum = func(centroided ...bool) Spectrum { return rd.spectrum(sn, centroided...) }

	return
}

// NScans returns the number of scans in the index
func (rd *RawData) NScans() int {
	return len(rd.Scanindex)
}

// Read only the initial header part of the file (for the juicy addresses)
func readHeaders(rs io.ReadSeeker) (RawFileInfo, Version) {

	hdr := new(FileHeader)
	info := new(RawFileInfo)

	//save position in file after reading, we need to sequentially
	//read some things in order to get to actual byte addresses
	pos := readAt(rs, 0, 0, hdr)

	ver := hdr.Version

	pos = readAt(rs, pos, ver, new(SequencerRow))
	pos = readAt(rs, pos, 0, new(AutoSamplerInfo))
	readAt(rs, pos, ver, info)

	return *info, ver
}

// Retrieve ...
func (data *RunHeader) Retrieve(rs io.ReadSeeker, info RawFileInfo, ver Version) {

	// read runheaders until we have a non-empty Scantrailer Address
	// indicating it is the runheader for a MS device (not a chromatography device)
	for i := 0; i < len(info.Preamble.RunHeaderAddr) && data.ScantrailerAddr == 0; i++ {
		data.Address = readAt(rs, info.Preamble.RunHeaderAddr[i], ver, data)
	}

	if data.ScantrailerAddr == 0 {
		msg.Custom(errors.New(""), "fatal")
	}

	return
}

// Retrieve ...
func (data *InstID) Retrieve(rs io.ReadSeeker, pos uint64, ver Version) {
	data.Address = readAt(rs, pos, ver, data)
	return
}

//// Readers
func (data *RunHeader) Read(r io.Reader, v Version) {

	binaryread(r, &data.SampleInfo)

	if v == 57 || v == 60 || v == 62 || v == 63 {
		binaryread(r, &data.Filename1)
		binaryread(r, &data.Filename2)
		binaryread(r, &data.Filename3)
		binaryread(r, &data.Filename4)
		binaryread(r, &data.Filename5)
		binaryread(r, &data.Filename6)
		binaryread(r, &data.Unknown1)
		binaryread(r, &data.Unknown2)
		binaryread(r, &data.Filename7)
		binaryread(r, &data.Filename8)
		binaryread(r, &data.Filename9)
		binaryread(r, &data.FilenameA)
		binaryread(r, &data.FilenameB)
		binaryread(r, &data.FilenameC)
		binaryread(r, &data.FilenameD)
		binaryread(r, &data.ScantrailerAddr32)
		binaryread(r, &data.ScanparamsAddr32)
		binaryread(r, &data.Unknown3)
		binaryread(r, &data.Unknown4)
		binaryread(r, &data.Nsegs)
		binaryread(r, &data.Unknown5)
		binaryread(r, &data.Unknown6)
		binaryread(r, &data.OwnAddr32)
		binaryread(r, &data.Unknown7)
		binaryread(r, &data.Unknown8)
	}

	if v >= 64 {
		binaryread(r, &data.Filename1)
		binaryread(r, &data.Filename2)
		binaryread(r, &data.Filename3)
		binaryread(r, &data.Filename4)
		binaryread(r, &data.Filename5)
		binaryread(r, &data.Filename6)
		binaryread(r, &data.Unknown1)
		binaryread(r, &data.Unknown2)
		binaryread(r, &data.Filename7)
		binaryread(r, &data.Filename8)
		binaryread(r, &data.Filename9)
		binaryread(r, &data.FilenameA)
		binaryread(r, &data.FilenameB)
		binaryread(r, &data.FilenameC)
		binaryread(r, &data.FilenameD)
		binaryread(r, &data.ScantrailerAddr32)
		binaryread(r, &data.ScanparamsAddr32)
		binaryread(r, &data.Unknown3)
		binaryread(r, &data.Unknown4)
		binaryread(r, &data.Nsegs)
		binaryread(r, &data.Unknown5)
		binaryread(r, &data.Unknown6)
		binaryread(r, &data.OwnAddr32)
		binaryread(r, &data.Unknown7)
		binaryread(r, &data.Unknown8)

		binaryread(r, &data.ScanindexAddr)
		binaryread(r, &data.DataAddr)
		binaryread(r, &data.InstlogAddr)
		binaryread(r, &data.ErrorlogAddr)
		binaryread(r, &data.Unknown9)
		binaryread(r, &data.ScantrailerAddr)
		binaryread(r, &data.ScanparamsAddr)
		binaryread(r, &data.Unknown10)
		binaryread(r, &data.OwnAddr)

		binaryread(r, &data.Unknown11)
		binaryread(r, &data.Unknown12)
		binaryread(r, &data.Unknown13)
		binaryread(r, &data.Unknown14)
		binaryread(r, &data.Unknown15)
		binaryread(r, &data.Unknown16)
		binaryread(r, &data.Unknown17)
		binaryread(r, &data.Unknown18)
		binaryread(r, &data.Unknown19)
		binaryread(r, &data.Unknown20)
		binaryread(r, &data.Unknown21)
		binaryread(r, &data.Unknown22)
		binaryread(r, &data.Unknown23)
		binaryread(r, &data.Unknown24)
		binaryread(r, &data.Unknown25)
		binaryread(r, &data.Unknown26)
		binaryread(r, &data.Unknown27)
		binaryread(r, &data.Unknown28)
		binaryread(r, &data.Unknown29)
		binaryread(r, &data.Unknown30)
		binaryread(r, &data.Unknown31)
		binaryread(r, &data.Unknown32)
		binaryread(r, &data.Unknown33)
		binaryread(r, &data.Unknown34)
	}

	// binaryread(r, &data.Unknown35)
	// binaryread(r, &data.Unknown36)
	// binaryread(r, &data.Device)
	// binaryread(r, &data.Model)
	// binaryread(r, &data.SN)
	// binaryread(r, &data.SWVer)
	// binaryread(r, &data.Tag1)
	// binaryread(r, &data.Tag2)
	// binaryread(r, &data.Tag3)
	// binaryread(r, &data.Tag4)
}

func (data *InstID) Read(r io.Reader, v Version) {
	//binaryread(r, data)
	binaryread(r, &data.Unknown1)
	binaryread(r, &data.Unknown2)
	binaryread(r, &data.Model1)
	binaryread(r, &data.Model2)
	binaryread(r, &data.Serial)
	binaryread(r, &data.Version)
	binaryread(r, &data.Tag1)
	binaryread(r, &data.Tag2)
	binaryread(r, &data.Tag3)
	binaryread(r, &data.Tag4)
}

func (data CDataPackets) Read(r io.Reader, v Version) {
	for i := range data {
		data[i].Read(r, v)
	}
}

func (data *TrailerLength) Read(r io.Reader, v Version) {
	binaryread(r, data)
}

func (data ScanDataPackets) Read(r io.Reader, v Version) {
	for i := range data {
		data[i].Read(r, v)
	}
}

func (data *ScanDataPacket) Read(r io.Reader, v Version) {
	binaryread(r, &data.Header)

	if data.Header.ProfileSize > 0 {
		binaryread(r, &data.Profile.FirstValue)
		binaryread(r, &data.Profile.Step)
		binaryread(r, &data.Profile.PeakCount)
		binaryread(r, &data.Profile.Nbins)

		data.Profile.Chunks = make([]ProfileChunk, data.Profile.PeakCount)

		for i := uint32(0); i < data.Profile.PeakCount; i++ {
			binaryread(r, &data.Profile.Chunks[i].Firstbin)
			binaryread(r, &data.Profile.Chunks[i].Nbins)
			if data.Header.Layout > 0 {
				binaryread(r, &data.Profile.Chunks[i].Fudge)
			}
			data.Profile.Chunks[i].Signal = make([]float32, data.Profile.Chunks[i].Nbins)
			binaryread(r, data.Profile.Chunks[i].Signal)
		}
	}

	if data.Header.PeaklistSize > 0 {
		binaryread(r, &data.PeakList.Count)
		data.PeakList.Peaks = make([]CentroidedPeak, data.PeakList.Count)
		binaryread(r, data.PeakList.Peaks)
	}

	data.DescriptorList = make([]PeakDescriptor, data.Header.DescriptorListSize)
	binaryread(r, data.DescriptorList)

	data.Unknown = make([]float32, data.Header.UnknownStreamSize)
	binaryread(r, data.Unknown)

	data.Triplets = make([]float32, data.Header.TripletStreamSize)
	binaryread(r, data.Triplets)
}

func (data ScanEvents) Read(r io.Reader, v Version) {
	for i := range data {
		data[i].Read(r, v)
	}
}

func (data *ScanEvent) Read(r io.Reader, v Version) {
	if v < 66 {
		switch {
		case v < 57:
			binaryread(r, data.Preamble[:41])
		case v >= 57 && v < 62:
			binaryread(r, data.Preamble[:80])
		case v >= 62 && v < 63:
			binaryread(r, data.Preamble[:120])
		case v >= 63:
			binaryread(r, data.Preamble[:128])
		}
		binaryread(r, &data.Nprecursors)
		data.Reaction = make([]Reaction, data.Nprecursors)
		for i := range data.Reaction {
			binaryread(r, &data.Reaction[i])
		}

		binaryread(r, &data.Unknown1[0])
		binaryread(r, &data.MZrange[0])
		binaryread(r, &data.Nparam)

		switch data.Nparam {
		case 4:
			binaryread(r, &data.Unknown2[0])
			binaryread(r, &data.A)
			binaryread(r, &data.B)
			binaryread(r, &data.C)
		case 7:
			binaryread(r, data.Unknown2[0:2])
			binaryread(r, &data.A)
			binaryread(r, &data.B)
			binaryread(r, &data.C)
			binaryread(r, data.Unknown2[2:4])
		}

		binaryread(r, data.Unknown1[1:3])
	} else { //v66
		binaryread(r, &data.Preamble)
		binaryread(r, &data.Unknown1[0])
		binaryread(r, &data.Nprecursors) //this is just a guess according to Gene Selkov
		if data.Preamble[10] == 1 {      //ms2 (dependent scan)
			data.Reaction = make([]Reaction, data.Nprecursors)
			for i := range data.Reaction {
				binaryread(r, &data.Reaction[i])
			}
			binaryread(r, data.Unknown2[0:2])
			binaryread(r, data.Unknown1[1:4])
			binaryread(r, &data.MZrange[0])
			binaryread(r, &data.Nparam)
		} else { //ms1
			binaryread(r, &data.MZrange[0])
			binaryread(r, data.Unknown1[1:5])
			binaryread(r, &data.MZrange[1])
			binaryread(r, data.Unknown1[5:8])
			binaryread(r, &data.MZrange[2])
			binaryread(r, &data.Nparam)
		}
		binaryread(r, data.Unknown2[2:4])
		binaryread(r, &data.A)
		binaryread(r, &data.B)
		binaryread(r, &data.C)
		binaryread(r, data.Unknown1[8:13])
	}
}

func (data ScanIndex) Read(r io.Reader, v Version) {
	for i := range data {
		data[i].Read(r, v)
	}
}

func (data *ScanIndexEntry) Read(r io.Reader, v Version) {
	if v == 66 {
		binaryread(r, data)
	} else if v == 64 {
		binaryread(r, &data.Offset32)
		binaryread(r, &data.Index) //starts from 0
		binaryread(r, &data.Scanevent)
		binaryread(r, &data.Scansegment)
		binaryread(r, &data.Next)
		binaryread(r, &data.Unknown1)
		binaryread(r, &data.DataPacketSize)
		binaryread(r, &data.Time)
		binaryread(r, &data.Totalcurrent)
		binaryread(r, &data.Baseintensity)
		binaryread(r, &data.Basemz)
		binaryread(r, &data.Lowmz)
		binaryread(r, &data.Highmz)
		binaryread(r, &data.Offset)
	} else {
		binaryread(r, &data.Offset32)
		binaryread(r, &data.Index) //starts from 0
		binaryread(r, &data.Scanevent)
		binaryread(r, &data.Scansegment)
		binaryread(r, &data.Next)
		binaryread(r, &data.Unknown1)
		binaryread(r, &data.DataPacketSize)
		binaryread(r, &data.Time)
		binaryread(r, &data.Totalcurrent)
		binaryread(r, &data.Baseintensity)
		binaryread(r, &data.Basemz)
		binaryread(r, &data.Lowmz)
		binaryread(r, &data.Highmz)

		data.Offset = uint64(data.Offset32)
	}
}

func (data CIndexEntries) Read(r io.Reader, v Version) {
	for i := range data {
		data[i].Read(r, v)
	}
}

func (data *CIndexEntry) Read(r io.Reader, v Version) {
	switch {
	case v < 64:
		binaryread(r, &data.Offset32)
		binaryread(r, &data.Index)
		binaryread(r, &data.Event)
		binaryread(r, &data.Unknown1)
		binaryread(r, &data.Unknown2)
		binaryread(r, &data.Unknown3)
		binaryread(r, &data.Unknown4)
		binaryread(r, &data.Unknown5)
		binaryread(r, &data.Time)
		binaryread(r, &data.Unknown6)
		binaryread(r, &data.Unknown7)
		binaryread(r, &data.Value)

		data.Offset = uint64(data.Offset32)
	default:
		binaryread(r, data)
	}
}

func (data *CDataPacket) Read(r io.Reader, v Version) {
	binaryread(r, data)
}

func (data *RawFileInfo) Read(r io.Reader, v Version) {
	binaryread(r, &data.Preamble.Methodfilepresent)
	binaryread(r, &data.Preamble.Year)
	binaryread(r, &data.Preamble.Month)
	binaryread(r, &data.Preamble.Weekday)
	binaryread(r, &data.Preamble.Day)
	binaryread(r, &data.Preamble.Hour)
	binaryread(r, &data.Preamble.Minute)
	binaryread(r, &data.Preamble.Second)
	binaryread(r, &data.Preamble.Millisecond)

	if v >= 57 {
		binaryread(r, &data.Preamble.Unknown1)
		binaryread(r, &data.Preamble.DataAddr32)
		binaryread(r, &data.Preamble.NControllers)
		binaryread(r, &data.Preamble.NControllers2)
		binaryread(r, &data.Preamble.Unknown2)
		binaryread(r, &data.Preamble.Unknown3)
		if v < 64 {
			data.Preamble.RunHeaderAddr32 = make([]uint32, data.Preamble.NControllers)
			data.Preamble.Unknown4 = make([]uint32, data.Preamble.NControllers)
			data.Preamble.Unknown5 = make([]uint32, data.Preamble.NControllers)
			for i := range data.Preamble.RunHeaderAddr32 {
				binaryread(r, &data.Preamble.RunHeaderAddr32[i])
				binaryread(r, &data.Preamble.Unknown4[i])
				binaryread(r, &data.Preamble.Unknown5[i])
			}

			data.Preamble.RunHeaderAddr = make([]uint64, data.Preamble.NControllers)
			for i := range data.Preamble.RunHeaderAddr {
				data.Preamble.RunHeaderAddr[i] = uint64(data.Preamble.RunHeaderAddr32[i])
			}

			if v == 57 {
				binaryread(r, data.Preamble.Padding1[:756-12*data.Preamble.NControllers])
			} else {
				binaryread(r, data.Preamble.Padding1[:760-12*data.Preamble.NControllers])
			}
		} else {
			binaryread(r, &data.Preamble.Padding1)
		}

	}
	if v >= 64 {
		binaryread(r, &data.Preamble.DataAddr)
		binaryread(r, &data.Preamble.Unknown6)

		data.Preamble.RunHeaderAddr = make([]uint64, data.Preamble.NControllers)
		data.Preamble.Unknown7 = make([]uint64, data.Preamble.NControllers)
		for i := range data.Preamble.RunHeaderAddr {
			binaryread(r, &data.Preamble.RunHeaderAddr[i])
			binaryread(r, &data.Preamble.Unknown7[i])
		}
		if v < 66 {
			binaryread(r, data.Preamble.Padding2[:1016-16*data.Preamble.NControllers])
		} else {
			binaryread(r, data.Preamble.Padding2[:1032-16*data.Preamble.NControllers])
		}
	}

	binaryread(r, &data.Heading1)
	binaryread(r, &data.Heading2)
	binaryread(r, &data.Heading3)
	binaryread(r, &data.Heading4)
	binaryread(r, &data.Heading5)
	binaryread(r, &data.Unknown1)
}

func (data *FileHeader) Read(r io.Reader, v Version) {
	binaryread(r, data)
}

func (data *SequencerRow) Read(r io.Reader, v Version) {
	binaryread(r, &data.Injection)

	binaryread(r, &data.Unknown1)
	binaryread(r, &data.Unknown2)
	binaryread(r, &data.ID)
	binaryread(r, &data.Comment)
	binaryread(r, &data.Userlabel1)
	binaryread(r, &data.Userlabel2)
	binaryread(r, &data.Userlabel3)
	binaryread(r, &data.Userlabel4)
	binaryread(r, &data.Userlabel5)
	binaryread(r, &data.Instmethod)
	binaryread(r, &data.Procmethod)
	binaryread(r, &data.Filename)
	binaryread(r, &data.Path)

	if v >= 57 {
		binaryread(r, &data.Vial)
		binaryread(r, &data.Unknown3)
		binaryread(r, &data.Unknown4)
		binaryread(r, &data.Unknown5)
	}
	if v >= 60 {
		binaryread(r, &data.Unknown6)
		binaryread(r, &data.Unknown7)
		binaryread(r, &data.Unknown8)
		binaryread(r, &data.Unknown9)
		binaryread(r, &data.Unknown10)
		binaryread(r, &data.Unknown11)
		binaryread(r, &data.Unknown12)
		binaryread(r, &data.Unknown13)
		binaryread(r, &data.Unknown14)
		binaryread(r, &data.Unknown15)
		binaryread(r, &data.Unknown16)
		binaryread(r, &data.Unknown17)
		binaryread(r, &data.Unknown18)
		binaryread(r, &data.Unknown19)
		binaryread(r, &data.Unknown20)
	}
}

func (data *AutoSamplerInfo) Read(r io.Reader, v Version) {
	binaryread(r, &data.Preamble)
	binaryread(r, &data.Text)
}

////

// Wrapper around binary.Read, reads both PascalStrings and structs from r
func binaryread(r io.Reader, data interface{}) {
	switch v := data.(type) {
	case *PascalString:
		binary.Read(r, binary.LittleEndian, &v.Length)
		v.Text = make([]uint16, v.Length)
		binary.Read(r, binary.LittleEndian, &v.Text)
	default:
		binary.Read(r, binary.LittleEndian, v)
	}
}

// For a version v Thermo File, starting at position pos, reads
// data, and returns the position in the file afterwards
func readAt(rs io.ReadSeeker, pos uint64, v Version, data reader) uint64 {

	spos, e1 := rs.Seek(int64(pos), 0)
	if e1 != nil {
		msg.Custom(errors.New("error seeking file"), "fatal")
	}

	data.Read(rs, v)

	spos, e2 := rs.Seek(0, 1)
	if e2 != nil {
		msg.Custom(errors.New("error determining position in file"), "fatal")
	}

	return uint64(spos)
}

// Copies the range in memory and then fills the Reader
// This tested faster than bufio or just reading away
func readBetween(rs io.ReadSeeker, begin uint64, end uint64, v Version, data reader) {

	_, e := rs.Seek(int64(begin), 0)
	if e != nil {
		msg.Custom(errors.New("error seeking file"), "fatal")
	}

	// may fail because of memory requirements
	b := make([]byte, end-begin)
	io.ReadFull(rs, b)

	data.Read(bytes.NewReader(b), v)

	return
}

// Spectrum returns an ms.Spectrum belonging to the scan number in argument
func (rd *RawData) spectrum(sn int, centroided ...bool) (s Spectrum) {

	//read Scan Packet for the scan
	scn := new(ScanDataPacket)
	readBetween(rd.File, rd.Scanindex[sn-1].Offset, rd.Scanindex[sn-1].Offset+uint64(rd.Scanindex[sn-1].DataPacketSize), 0, scn)

	if len(centroided) > 0 && centroided[0] || scn.Profile.PeakCount <= 0 {
		//Save the Centroided Peaks, they also occur in profile scans but
		//overlap with profiles, Thermo always does centroiding just for fun
		for i := uint32(0); i < scn.PeakList.Count; i++ {
			s = append(s,
				Peak{
					Mz: float64(scn.PeakList.Peaks[i].Mz),
					I:  scn.PeakList.Peaks[i].Abundance,
				})
		}
	} else {
		//convert Hz values into m/z and save the profile peaks
		for i := uint32(0); i < scn.Profile.PeakCount; i++ {
			for j := uint32(0); j < scn.Profile.Chunks[i].Nbins; j++ {
				tmpmz := rd.Scanevents[sn-1].Convert(scn.Profile.FirstValue+
					float64(scn.Profile.Chunks[i].Firstbin+j)*scn.Profile.Step) +
					float64(scn.Profile.Chunks[i].Fudge)
				s = append(s, Peak{Mz: tmpmz, I: scn.Profile.Chunks[i].Signal[j]})
			}
		}
	}

	return
}

// Chromatography Experimental: read out chromatography data from a connected instrument
func (rd *RawData) Chromatography(instr int) (cdata CDataPackets) {
	info, ver := readHeaders(rd.File)

	if uint32(instr) > info.Preamble.NControllers-1 {
		log.Print(instr, " is higher than number of extra controllers: ", info.Preamble.NControllers-1)
		return
	}

	rh := new(RunHeader)
	readAt(rd.File, info.Preamble.RunHeaderAddr[instr], ver, rh)
	//The ScantrailerAddr has to be 0. in other words: we're not looking at the MS runheader
	if rh.ScantrailerAddr != 0 {
		log.Println("You selected the MS instrument, no chromatography data can be read.")
		return
	}

	//The instrument RunHeader contains an interesting address: DataAddr
	//There is another address ScanIndexAddr, which points to CIndexEntry
	//containers at ScanIndexAddr. Less data can be read for now

	nScan := uint64(rh.SampleInfo.LastScanNumber - rh.SampleInfo.FirstScanNumber + 1)
	cdata = make(CDataPackets, nScan)
	for i := uint64(0); i < nScan; i++ {
		readAt(rd.File, rh.DataAddr+i*16, ver, &cdata[i]) //16 bytes of CDataPacket
	}
	return cdata
}

// Convert Hz values to m/z
func (data ScanEvent) Convert(v float64) float64 {
	switch data.Nparam {
	case 4:
		return data.A + data.B/v + data.C/v/v
	case 5, 7:
		return data.A + data.B/v/v + data.C/v/v/v/v
	default:
		return v
	}
}

// String decoders
func (t filename) String() string {
	return string(utf16.Decode(t[:]))
}

func (t audittag) String() string {
	return string(utf16.Decode(t[:]))
}

func (t headertag) String() string {
	return string(utf16.Decode(t[:]))
}

func (t signature) String() string {
	return string(utf16.Decode(t[:]))
}

func (t PascalString) String() string {
	return string(utf16.Decode(t.Text[:]))
}
