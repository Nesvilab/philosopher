<p align="center">
  <img height="420" width="593" src="/images/philosopher.png">
</p>

[![Release](https://img.shields.io/github/release/nesvilab/philosopher.svg?color=purple&style=for-the-badge)](https://github.com/Nesvilab/philosopher/releases/latest)
![Golang](https://img.shields.io/badge/Go-1.19.4-blue.svg?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nesvilab/philosopher?style=for-the-badge&color=red&logo=appveyor)](https://goreportcard.com/report/github.com/Nesvilab/philosopher)
![GitHub](https://img.shields.io/github/license/Nesvilab/philosopher?style=for-the-badge)
![](https://img.shields.io/github/downloads/Nesvilab/philosopher/total.svg?color=red&style=for-the-badge)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Nesvilab/philosopher/go.yml?style=for-the-badge)

#### Philosopher is a fast, easy-to-use, scalable, and versatile data analysis software for mass spectrometry-based proteomics. Philosopher is dependency-free and can analyze both traditional database searches and open searches for post-translational modification (PTM) discovery. 

- Database downloading and formatting.

- Peptide-spectrum matching with MSFragger and Comet.

- Peptide assignment validation with PeptideProphet.

- Multi-level integrative analysis with iProphet.

- PTM site localization with PTMProphet.

- Protein inference with ProteinProphet.

- FDR filtering with custom algorithms.

  - Two-dimensional filtering for simultaneous control of PSM and Protein FDR levels.
  - Sequential FDR estimation for large data sets using filtered PSM and proteins lists.

- Label-free quantification via spectral counting and MS1 intensities.

- Label-based quantification using TMT and iTRAQ.

- Quantification based on functional protein groups.

- Multi-level detailed reports for peptides, ions, and proteins.

- Support for REPRINT and MSstats.

## Download
Download the latest version [here](https://github.com/nesvilab/philosopher/releases/latest).


## How to Use
- [Philosopher basics](https://github.com/Nesvilab/philosopher/wiki/Philosopher-Basics) - general usage information
- [Preparing protein databases](https://github.com/Nesvilab/philosopher/wiki/How-to-Prepare-a-Protein-Database) - download and format sequences
- [Simple data analysis](https://github.com/Nesvilab/philosopher/wiki/Simple-Data-Analysis) - basic step-by-step tutorial
- [Using pipeline for TMT analysis](https://github.com/Nesvilab/philosopher/wiki/Pipeline-mode-for-TMT-analysis) - pipeline analysis of a large data set
- [Step-by-step TMT analysis](https://github.com/Nesvilab/philosopher/wiki/Step-by-step-TMT-analysis) - step-by-step tutorial for isobaric quantification of a small data set
- [Open search analysis](https://github.com/Nesvilab/philosopher/wiki/Open-Search-Analysis) - step-by-step tutorial for open searches
- [Step-by-step analysis with Comet](https://github.com/Nesvilab/philosopher/wiki/Step-by-step-analysis-with-Comet) - step-by-step tutorial with Comet search
- [Protein-protein interaction analysis](https://github.com/Nesvilab/philosopher/wiki/REPRINT-Analysis) - analyze AP-MS data for downstream use with REPRINT

## Documentation
See the [documentation](https://github.com/Nesvilab/philosopher/wiki/Home) for more details about the available commands.

## Questions, requests and bug reports
If you have any questions or remarks please use the [Discussion board](https://github.com/Nesvilab/philosopher/discussions). If you want to report a bug, please use the [Issue tracker](https://github.com/nesvilab/philosopher/issues).


## How to cite
da Veiga Leprevost F, Haynes SE, Avtonomov DM, Chang HY, Shanmugam AK, Mellacheruvu D, Kong AT, Nesvizhskii AI. [Philosopher: a versatile toolkit for shotgun proteomics data analysis](https://doi.org/10.1038/s41592-020-0912-y). Nat Methods. 2020 Sep;17(9):869-870. doi: 10.1038/s41592-020-0912-y. PMID: 32669682; PMCID: PMC7509848.

## About the authors, and contributors
[Felipe da Veiga Leprevost (main author)](http://prvst.github.io)

[Sarah Haynes](https://scholar.google.com/citations?user=HtRSUKkAAAAJ&hl=en)

[Guo Ci Teo](https://github.com/guoci)

[Alexey Nesvizhskii's research group](http://www.nesvilab.org/)
