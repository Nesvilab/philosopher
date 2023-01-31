package bio

import (
	"errors"

	"github.com/Nesvilab/philosopher/lib/msg"
)

// AminoAcid representation
type AminoAcid struct {
	Code            string
	ShortName       string
	Name            string
	MonoIsotopeMass float64
	AverageMass     float64
}

// OligoPeptide is an array of Aminoacids
type OligoPeptide []AminoAcid

// New return the correct information for the give aminoacid
func New(name string) AminoAcid {

	var aa AminoAcid

	switch name {

	case "Alanine":
		aa = AminoAcid{Code: "A", ShortName: "Ala", Name: "Alanine", MonoIsotopeMass: 71.037113805, AverageMass: 71.0779}
	case "Arginine":
		aa = AminoAcid{Code: "R", ShortName: "Arg", Name: "Arginine", MonoIsotopeMass: 156.101111050, AverageMass: 156.18568}
	case "Asparagine":
		aa = AminoAcid{Code: "N", ShortName: "Asn", Name: "Asparagine", MonoIsotopeMass: 114.042927470, AverageMass: 114.10264}
	case "Aspartic Acid":
		aa = AminoAcid{Code: "D", ShortName: "Asp", Name: "Aspartic Acid", MonoIsotopeMass: 115.026943065, AverageMass: 115.0874}
	case "Cysteine":
		aa = AminoAcid{Code: "C", ShortName: "Cys", Name: "Cysteine", MonoIsotopeMass: 103.009184505, AverageMass: 103.1429}
	case "Glutamine":
		aa = AminoAcid{Code: "E", ShortName: "Glu", Name: "Glutamine", MonoIsotopeMass: 129.042593135, AverageMass: 129.11398}
	case "Glutamic Acid":
		aa = AminoAcid{Code: "Q", ShortName: "Gln", Name: "Glutamic Acid", MonoIsotopeMass: 128.058577540, AverageMass: 128.12922}
	case "Glycine":
		aa = AminoAcid{Code: "G", ShortName: "Gly", Name: "Glycine", MonoIsotopeMass: 57.021463735, AverageMass: 57.05132}
	case "Histidine":
		aa = AminoAcid{Code: "H", ShortName: "His", Name: "Histidine", MonoIsotopeMass: 137.058911875, AverageMass: 137.13928}
	case "Isoleucine":
		aa = AminoAcid{Code: "I", ShortName: "Ile", Name: "Isoleucine", MonoIsotopeMass: 113.084064015, AverageMass: 113.15764}
	case "Leucine":
		aa = AminoAcid{Code: "L", ShortName: "Leu", Name: "Leucine", MonoIsotopeMass: 113.084064015, AverageMass: 113.15764}
	case "Lysine":
		aa = AminoAcid{Code: "K", ShortName: "Lys", Name: "Lysine", MonoIsotopeMass: 128.094963050, AverageMass: 128.17228}
	case "Methionine":
		aa = AminoAcid{Code: "M", ShortName: "Met", Name: "Methionine", MonoIsotopeMass: 131.040484645, AverageMass: 131.19606}
	case "Phenylalanine":
		aa = AminoAcid{Code: "F", ShortName: "Phe", Name: "Phenylalanine", MonoIsotopeMass: 147.068413945, AverageMass: 147.17386}
	case "Proline":
		aa = AminoAcid{Code: "P", ShortName: "Pro", Name: "Proline", MonoIsotopeMass: 97.052763875, AverageMass: 97.11518}
	case "Serine":
		aa = AminoAcid{Code: "S", ShortName: "Ser", Name: "Serine", MonoIsotopeMass: 87.032028435, AverageMass: 87.0773}
	case "Threonine":
		aa = AminoAcid{Code: "T", ShortName: "Thr", Name: "Threonine", MonoIsotopeMass: 101.047678505, AverageMass: 101.10388}
	case "Tryptophan":
		aa = AminoAcid{Code: "W", ShortName: "Trp", Name: "Tryptophan", MonoIsotopeMass: 186.079312980, AverageMass: 186.2099}
	case "Tyrosine":
		aa = AminoAcid{Code: "Y", ShortName: "Tyr", Name: "Tyrosine", MonoIsotopeMass: 163.063328575, AverageMass: 163.17326}
	case "Valine":
		aa = AminoAcid{Code: "V", ShortName: "Val", Name: "Valine", MonoIsotopeMass: 99.068413945, AverageMass: 99.13106}
	default:
		msg.Custom(errors.New("amino acid not found"), "warning")
		return aa
	}

	return aa
}
