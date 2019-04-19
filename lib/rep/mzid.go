package rep

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/prvst/philosopher/lib/dat"
	"github.com/prvst/philosopher/lib/psi"
	"github.com/sirupsen/logrus"
)

// MzIdentMLReport creates a MzIdentML structure to be encoded
func (e Evidence) MzIdentMLReport(version, database string) error {

	var mzid psi.MzIdentML

	t := time.Now()
	var idCounter = 0

	// Header
	mzid.Name = "foo"
	mzid.ID = "Philosopher toolkit"
	mzid.Version = version
	mzid.CreationDate = t.Format(time.ANSIC)

	// CVlist
	mzid.CvList.CV = append(mzid.CvList.CV, psi.CV{ID: "PSI-MS", URI: "https://raw.githubusercontent.com/HUPO-PSI/psi-ms-CV/master/psi-ms.obo", FullName: "PSI-MS"})
	mzid.CvList.CV = append(mzid.CvList.CV, psi.CV{ID: "UNIMOD", URI: "http://www.unimod.org/obo/unimod.obo", FullName: "UNIMOD"})
	mzid.CvList.CV = append(mzid.CvList.CV, psi.CV{ID: "UO", URI: "https://raw.githubusercontent.com/bio-ontology-research-group/unit-ontology/master/unit.obo", FullName: "UNIT-ONTOLOGY"})
	mzid.CvList.CV = append(mzid.CvList.CV, psi.CV{ID: "PRIDE", URI: "https://github.com/PRIDE-Utilities/pride-ontology/blob/master/pride_cv.obo", FullName: "PRIDE"})
	mzid.CvList.Count = len(mzid.CvList.CV)

	// AnalysisSoftwareList
	aa := &psi.AnalysisSoftware{
		ID:      "Philosopher",
		Name:    "Philosopher toolkit",
		URI:     "https://philosopher.nesvilab.org",
		Version: version,
		ContactRole: psi.ContactRole{
			ContactRef: "PS_DEV",
			Role: psi.Role{
				CVParam: psi.CVParam{
					CVRef:     "PSI-MS",
					Accession: "MS:1001267",
					Name:      "software vendor",
				},
			},
		},
		SoftwareName: psi.SoftwareName{
			CVParam: psi.CVParam{
				CVRef:     "PSI-MS",
				Accession: "XXXX",
				Name:      "Philosopher",
			},
		},
		Customizations: psi.Customizations{
			Value: "No customizations",
		},
	}
	mzid.AnalysisSoftwareList.AnalysisSoftware = append(mzid.AnalysisSoftwareList.AnalysisSoftware, *aa)

	//Provider
	provider := &psi.Provider{
		ID: "PROVIDER",
		ContactRole: psi.ContactRole{
			ContactRef: "Philosopher_Author_FVL",
			Role: psi.Role{
				CVParam: psi.CVParam{
					CVRef:     "PSI-MS",
					Accession: "MS:1001271",
					Name:      "researcher",
				},
			},
		},
	}
	mzid.Provider = *provider

	// AuditCollection

	auditCol := &psi.AuditCollection{
		Person: psi.Person{
			ID:        "Philosopher_Author_FVL",
			LastName:  "da Veiga Leprevost",
			FirstName: "Felipe",
			CVParam: []psi.CVParam{
				psi.CVParam{
					Name:      "contact email",
					Value:     "felipevl@umich.edu",
					CVRef:     "PSI-MS",
					Accession: "MS:1000589",
				},
				psi.CVParam{
					Name:      "contact URL",
					Value:     "http://nesvilab.org",
					CVRef:     "PSI-MS",
					Accession: "MS:1000588",
				},
			},
			Affiliation: []psi.Affiliation{
				psi.Affiliation{
					OrganizationRef: "University of Michigan",
				},
			},
		},
		Organization: psi.Organization{
			ID:   "Nesvilab",
			Name: "Proteomics and Integrative Bioinformatics Lab",
			CVParam: []psi.CVParam{
				psi.CVParam{
					Name:      "contact name",
					Value:     "Alexey I. Nesvizhskii",
					CVRef:     "PSI-MS",
					Accession: "MS:1000586",
				},
				psi.CVParam{
					Name:      "contact address",
					Value:     "1301 Catherinse St., Ann Arbor, MI",
					CVRef:     "PSI-MS",
					Accession: "MS:1000587",
				},
				psi.CVParam{
					Name:      "contact URL",
					Value:     "http://nesvilab.org",
					CVRef:     "PSI-MS",
					Accession: "MS:1000588",
				},
				psi.CVParam{
					Name:      "contact email",
					Value:     "nesvi@med.umich.edu",
					CVRef:     "PSI-MS",
					Accession: "MS:1000589",
				},
			},
		},
	}
	mzid.AuditCollection = *auditCol

	// SequenceCollection - DBSequence
	var dtb dat.Base
	dtb.Restore()
	// if len(dtb.Records) < 1 {
	// 	return f, errors.New("Database data not available, interrupting processing")
	// }

	var seqs []psi.DBSequence
	for _, i := range dtb.Records {

		db := &psi.DBSequence{
			ID:                i.ID,
			Accession:         i.ID,
			SearchDatabaseRef: "",
			CVParam: []psi.CVParam{
				psi.CVParam{
					CVRef:     "PSI-MS",
					Accession: "MS:1001088",
					Name:      "protein description",
					Value:     i.Description,
				},
				psi.CVParam{
					CVRef:     "PSI-MS",
					Accession: "MS:1001344",
					Name:      "AA sequence",
				},
			},
			Seq: psi.Seq{
				Value: i.Sequence,
			},
		}

		seqs = append(seqs, *db)
	}
	//mzid.SequenceCollection.DBSequence = seqs
	seqs = nil

	// SequenceCollection - Peptide
	var peps []psi.Peptide
	for _, i := range e.Peptides {

		p := psi.Peptide{
			ID: i.Sequence,
			PeptideSequence: psi.PeptideSequence{
				Value: i.Sequence,
			},
		}

		for _, j := range i.Modifications.Index {
			if j.Name != "Unknown" {
				mod := psi.Modification{
					AvgMassDelta:          j.AverageMass,
					MonoIsotopicMassDelta: j.MonoIsotopicMass,
					Residues:              j.AminoAcid,
					Location:              j.Position,
					CVParam: []psi.CVParam{
						psi.CVParam{
							CVRef:     "UNIMOD",
							Accession: j.ID,
							Name:      j.Name,
						},
					},
				}

				p.Modification = append(p.Modification, mod)
			}
		}

		peps = append(peps, p)
	}
	//mzid.SequenceCollection.Peptide = peps
	peps = nil

	// SequenceCollection - PeptideEvidence
	var pevs []psi.PeptideEvidence
	idCounter = 0
	for _, i := range e.PSM {

		idCounter++

		evi := psi.PeptideEvidence{
			DBSequenceRef: i.ProteinID,
			ID:            fmt.Sprintf("PepEv_%d", idCounter),
			IsDecoy:       strconv.FormatBool(i.IsDecoy),
			PeptideRef:    i.Peptide,
			Pre:           i.PrevAA,
			Post:          i.NextAA,
		}

		pevs = append(pevs, evi)
	}
	//mzid.SequenceCollection.PeptideEvidence = pevs
	pevs = nil

	var sources = make(map[string]uint8)
	for _, i := range e.PSM {
		s := strings.Split(i.Spectrum, ".")
		sources[s[0]]++
	}

	// AnalysisCollection
	idCounter = 0
	ac := &psi.AnalysisCollection{}
	for i := range sources {

		idCounter++

		si := &psi.SpectrumIdentification{
			SpectrumIdentificationListRef:     fmt.Sprintf("SIL_%d", idCounter),
			ID:                                fmt.Sprintf("SpecIdent_%d", idCounter),
			SpectrumIdentificationProtocolRef: fmt.Sprintf("SearchProtocol_%d", idCounter),
			InputSpectra: []psi.InputSpectra{
				psi.InputSpectra{
					SpectraDataRef: i,
				},
			},
			SearchDatabaseRef: []psi.SearchDatabaseRef{
				psi.SearchDatabaseRef{
					SearchDatabaseRef: database,
				},
			},
		}

		ac.SpectrumIdentification = append(ac.SpectrumIdentification, *si)
	}

	ac.ProteinDetection = psi.ProteinDetection{
		ProteinDetectionProtocolRef: "Philosopher_protocol",
		ProteinDetectionListRef:     "Protein Groups",
		ID:                          "Phi_1",
	}
	mzid.AnalysisCollection = *ac

	// AnalysisProtocolCollection
	//TODO

	//DataCollection
	dta := psi.DataCollection{}

	idCounter = 0
	for i := range sources {
		sf := &psi.SourceFile{
			ID:       i,
			Location: i,
			Name:     i,
		}

		dta.Inputs.SourceFile = append(dta.Inputs.SourceFile, *sf)
	}

	sdb := psi.SearchDatabase{
		ID:                   database,
		NumDatabaseSequences: len(dtb.Records),
		Location:             database,
		FileFormat: psi.FileFormat{
			CVParam: psi.CVParam{
				CVRef:     "PSI-MS",
				Accession: "MS:1001348",
				Name:      "FASTA format",
			},
		},
	}
	mzid.DataCollection.Inputs.SearchDatabase[0] = sdb

	// Burn!
	err := mzid.Write()
	if err != nil {
		logrus.Fatal(err)
	}

	return nil
}
