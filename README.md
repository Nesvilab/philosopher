<p align="center">
  <img height="420" width="593" src="/images/philosopher.png">
</p>

[![Release](https://img.shields.io/github/release/nesvilab/philosopher.svg?color=purple&style=for-the-badge)](https://github.com/Nesvilab/philosopher/releases/latest)
![Golang](https://img.shields.io/badge/Go-1.19.4-blue.svg?style=for-the-badge)
[![Go Report Card](https://goreportcard.com/badge/github.com/Nesvilab/philosopher?style=for-the-badge&color=red&logo=appveyor)](https://goreportcard.com/report/github.com/Nesvilab/philosopher)
![GitHub](https://img.shields.io/github/license/Nesvilab/philosopher?style=for-the-badge)
![](https://img.shields.io/github/downloads/Nesvilab/philosopher/total.svg?color=red&style=for-the-badge)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Nesvilab/philosopher/go.yml?style=for-the-badge)

#### Philosopher is a fast, easy-to-use, scalable, and versatile data analysis software for mass spectrometry-based proteomics. It is also a depencency-free wraper of Trans-Proteomic Pipeline (PeptideProphet, iProphet, PTMProphet, and ProteinProphet).

- Database downloading and formatting.

- Peptide assignment validation with PeptideProphet.

- Multi-level integrative analysis with iProphet.

- PTM site localization with PTMProphet.

- Protein inference with ProteinProphet.

- FDR filtering with custom algorithms.

  - Two-dimensional filtering for simultaneous control of PSM and Protein FDR levels.
  - Sequential FDR estimation for large data sets using filtered PSM and proteins lists.

- Label-free quantification via spectral counting and MS1 intensities.

- Label-based quantification using TMT and iTRAQ.

- Multi-level detailed reports for peptides, ions, and proteins.

- Support for REPRINT and MSstats.


## How to Use
Philosopher is part of [FragPipe](https://fragpipe.nesvilab.org/) which has a user-friendly GUI.

## Documentation
See the [documentation](https://github.com/Nesvilab/philosopher/wiki/Home) for more details about the available commands.

## Questions, requests and bug reports
If you have any questions or remarks please use the [Discussion board](https://github.com/Nesvilab/philosopher/discussions). If you want to report a bug, please use the [Issue tracker](https://github.com/nesvilab/philosopher/issues).


## How to cite
da Veiga Leprevost F, Haynes SE, Avtonomov DM, Chang HY, Shanmugam AK, Mellacheruvu D, Kong AT, Nesvizhskii AI. [Philosopher: a versatile toolkit for shotgun proteomics data analysis](https://doi.org/10.1038/s41592-020-0912-y). Nat Methods. 2020 Sep;17(9):869-870. doi: 10.1038/s41592-020-0912-y. PMID: 32669682; PMCID: PMC7509848.

## About the authors
[Alexey Nesvizhskii's research group](http://www.nesvilab.org/)
