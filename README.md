# Philosopher
A data processing toolkit for shotgun proteomics.

![Golang](https://img.shields.io/badge/Go-1.8.0-blue.svg) ![Version](https://img.shields.io/badge/version-rc1-blue.svg)


## Features
Philosopher provides easy access to third-party tools and custom algorithms allowing users to develop proteomics analysis, from Peptide Spectrum Matching to annotated protein reports.

- Database downloading and formatting.

- Peptide Spectrum Matching with Comet.

- Peptide assignment validation with PeptideProphet.

- Multi-level integrative analysis with iProphet.

- Protein inference with ProteinProphet.

- FDR filtering with custom algorithms.

  - Two-dimensional filtering for simultaneous control of PSM and Protein FDR levels.
  - Sequential FDR estimantion using filtered PSM and proteins lists.

- Label-free Quantification (Spectral Count and MS1 intensities)

- Isotope lable quantification (TMT)

- Detailed protein reports with optional protein annotation.
