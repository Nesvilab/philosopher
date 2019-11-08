package bio_test

import (
	. "github.com/nesvilab/philosopher/lib/bio"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Bio", func() {

	Context("Amino acid instantiation", func() {

		It("Alanine", func() {
			a := New("Alanine")
			Expect(a.Code).To(Equal("A"))
			Expect(a.ShortName).To(Equal("Ala"))
			Expect(a.Name).To(Equal("Alanine"))
			Expect(a.MonoIsotopeMass).To(Equal(71.037113805))
			Expect(a.AverageMass).To(Equal(71.0779))
		})

		It("Arginine", func() {
			a := New("Arginine")
			Expect(a.Code).To(Equal("R"))
			Expect(a.ShortName).To(Equal("Arg"))
			Expect(a.Name).To(Equal("Arginine"))
			Expect(a.MonoIsotopeMass).To(Equal(156.101111050))
			Expect(a.AverageMass).To(Equal(156.18568))
		})

		It("Asparagine", func() {
			a := New("Asparagine")
			Expect(a.Code).To(Equal("N"))
			Expect(a.ShortName).To(Equal("Asn"))
			Expect(a.Name).To(Equal("Asparagine"))
			Expect(a.MonoIsotopeMass).To(Equal(114.042927470))
			Expect(a.AverageMass).To(Equal(114.10264))
		})

		It("Aspartic Acid", func() {
			a := New("Aspartic Acid")
			Expect(a.Code).To(Equal("D"))
			Expect(a.ShortName).To(Equal("Asp"))
			Expect(a.Name).To(Equal("Aspartic Acid"))
			Expect(a.MonoIsotopeMass).To(Equal(115.026943065))
			Expect(a.AverageMass).To(Equal(115.0874))
		})

		It("Cysteine", func() {
			a := New("Cysteine")
			Expect(a.Code).To(Equal("C"))
			Expect(a.ShortName).To(Equal("Cys"))
			Expect(a.Name).To(Equal("Cysteine"))
			Expect(a.MonoIsotopeMass).To(Equal(103.009184505))
			Expect(a.AverageMass).To(Equal(103.1429))
		})

		It("Glutamine", func() {
			a := New("Glutamine")
			Expect(a.Code).To(Equal("E"))
			Expect(a.ShortName).To(Equal("Glu"))
			Expect(a.Name).To(Equal("Glutamine"))
			Expect(a.MonoIsotopeMass).To(Equal(129.042593135))
			Expect(a.AverageMass).To(Equal(129.11398))
		})

		It("Glutamic Acid", func() {
			a := New("Glutamic Acid")
			Expect(a.Code).To(Equal("Q"))
			Expect(a.ShortName).To(Equal("Gln"))
			Expect(a.Name).To(Equal("Glutamic Acid"))
			Expect(a.MonoIsotopeMass).To(Equal(128.058577540))
			Expect(a.AverageMass).To(Equal(128.12922))
		})

		It("Glycine", func() {
			a := New("Glycine")
			Expect(a.Code).To(Equal("G"))
			Expect(a.ShortName).To(Equal("Gly"))
			Expect(a.Name).To(Equal("Glycine"))
			Expect(a.MonoIsotopeMass).To(Equal(57.021463735))
			Expect(a.AverageMass).To(Equal(57.05132))
		})

		It("Histidine", func() {
			a := New("Histidine")
			Expect(a.Code).To(Equal("H"))
			Expect(a.ShortName).To(Equal("His"))
			Expect(a.Name).To(Equal("Histidine"))
			Expect(a.MonoIsotopeMass).To(Equal(137.058911875))
			Expect(a.AverageMass).To(Equal(137.13928))
		})

		It("Isoleucine", func() {
			a := New("Isoleucine")
			Expect(a.Code).To(Equal("I"))
			Expect(a.ShortName).To(Equal("Ile"))
			Expect(a.Name).To(Equal("Isoleucine"))
			Expect(a.MonoIsotopeMass).To(Equal(113.084064015))
			Expect(a.AverageMass).To(Equal(113.15764))
		})

		It("Leucine", func() {
			a := New("Leucine")
			Expect(a.Code).To(Equal("L"))
			Expect(a.ShortName).To(Equal("Leu"))
			Expect(a.Name).To(Equal("Leucine"))
			Expect(a.MonoIsotopeMass).To(Equal(113.084064015))
			Expect(a.AverageMass).To(Equal(113.15764))
		})

		It("Lysine", func() {
			a := New("Lysine")
			Expect(a.Code).To(Equal("K"))
			Expect(a.ShortName).To(Equal("Lys"))
			Expect(a.Name).To(Equal("Lysine"))
			Expect(a.MonoIsotopeMass).To(Equal(128.094963050))
			Expect(a.AverageMass).To(Equal(128.17228))
		})

		It("Methionine", func() {
			a := New("Methionine")
			Expect(a.Code).To(Equal("M"))
			Expect(a.ShortName).To(Equal("Met"))
			Expect(a.Name).To(Equal("Methionine"))
			Expect(a.MonoIsotopeMass).To(Equal(131.040484645))
			Expect(a.AverageMass).To(Equal(131.19606))
		})

		It("Phenylalanine", func() {
			a := New("Phenylalanine")
			Expect(a.Code).To(Equal("F"))
			Expect(a.ShortName).To(Equal("Phe"))
			Expect(a.Name).To(Equal("Phenylalanine"))
			Expect(a.MonoIsotopeMass).To(Equal(147.068413945))
			Expect(a.AverageMass).To(Equal(147.17386))
		})

		It("Proline", func() {
			a := New("Proline")
			Expect(a.Code).To(Equal("P"))
			Expect(a.ShortName).To(Equal("Pro"))
			Expect(a.Name).To(Equal("Proline"))
			Expect(a.MonoIsotopeMass).To(Equal(97.052763875))
			Expect(a.AverageMass).To(Equal(97.11518))
		})

		It("Serine", func() {
			a := New("Serine")
			Expect(a.Code).To(Equal("S"))
			Expect(a.ShortName).To(Equal("Ser"))
			Expect(a.Name).To(Equal("Serine"))
			Expect(a.MonoIsotopeMass).To(Equal(87.032028435))
			Expect(a.AverageMass).To(Equal(87.0773))
		})

		It("Threonine", func() {
			a := New("Threonine")
			Expect(a.Code).To(Equal("T"))
			Expect(a.ShortName).To(Equal("Thr"))
			Expect(a.Name).To(Equal("Threonine"))
			Expect(a.MonoIsotopeMass).To(Equal(101.047678505))
			Expect(a.AverageMass).To(Equal(101.10388))
		})

		It("Tryptophan", func() {
			a := New("Tryptophan")
			Expect(a.Code).To(Equal("W"))
			Expect(a.ShortName).To(Equal("Trp"))
			Expect(a.Name).To(Equal("Tryptophan"))
			Expect(a.MonoIsotopeMass).To(Equal(186.079312980))
			Expect(a.AverageMass).To(Equal(186.2099))
		})

		It("Tyrosine", func() {
			a := New("Tyrosine")
			Expect(a.Code).To(Equal("Y"))
			Expect(a.ShortName).To(Equal("Tyr"))
			Expect(a.Name).To(Equal("Tyrosine"))
			Expect(a.MonoIsotopeMass).To(Equal(163.063328575))
			Expect(a.AverageMass).To(Equal(163.17326))
		})

		It("Valine", func() {
			a := New("Valine")
			Expect(a.Code).To(Equal("V"))
			Expect(a.ShortName).To(Equal("Val"))
			Expect(a.Name).To(Equal("Valine"))
			Expect(a.MonoIsotopeMass).To(Equal(99.068413945))
			Expect(a.AverageMass).To(Equal(99.13106))
		})

		It("Invalid", func() {
			New("Foobar")
		})
	})

	Describe("Bio::che", func() {
		Context("Contant values", func() {
			It("Proton", func() {
				p := Proton
				Expect(p).To(Equal(1.007276467))
			})
		})
	})

	Describe("Bio::enz", func() {

		Context("Enzyme instantiation", func() {

			It("Trypsin", func() {
				var e Enzyme
				e.Synth("Trypsin")
				Expect(e.Name).To(Equal("trypsin"))
				Expect(e.Pattern).To(Equal("KR[^P]"))
				Expect(e.Join).To(Equal("KR"))
			})

			It("Lys_c", func() {
				var e Enzyme
				e.Synth("Lys_c")
				Expect(e.Name).To(Equal("lys_c"))
				Expect(e.Pattern).To(Equal("K[^P]"))
				Expect(e.Join).To(Equal("K"))
			})

			It("Lys_n", func() {
				var e Enzyme
				e.Synth("Lys_n")
				Expect(e.Name).To(Equal("lys_n"))
				Expect(e.Pattern).To(Equal("K"))
				Expect(e.Join).To(Equal("K"))
			})

			It("Chymotrypsin", func() {
				var e Enzyme
				e.Synth("Chymotrypsin")
				Expect(e.Name).To(Equal("chymotrypsin"))
				Expect(e.Pattern).To(Equal("FWYL[^P]"))
				Expect(e.Join).To(Equal("K"))
			})

			It("Glu_c", func() {
				var e Enzyme
				e.Synth("Glu_c")
				Expect(e.Name).To(Equal("glu_c"))
				Expect(e.Pattern).To(Equal("DE[^P]"))
				Expect(e.Join).To(Equal("K"))
			})

		})
	})

})
