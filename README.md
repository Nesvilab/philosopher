# Philosopher
A data processing toolkit for shotgun proteomics.

[![Release](https://img.shields.io/github/release/nesvilab/philosopher.svg?color=purple&style=for-the-badge)](https://github.com/nesvilab/philosopher/releases/latest)
![Golang](https://img.shields.io/badge/Go-1.13.1-blue.svg?style=for-the-badge)
![CircleCI](https://img.shields.io/circleci/build/github/Nesvilab/philosopher-source?style=for-the-badge&token=8fc2895a5f10b2507a73d672ad89d50798b417ad)
![](https://img.shields.io/github/downloads/Nesvilab/philosopher/total.svg?color=red&style=for-the-badge)


## Features
Philosopher provides easy access to third-party tools and custom algorithms allowing users to develop proteomics analysis, from Peptide Spectrum Matching to annotated protein reports. Philosopher is also tunned for Open Search analysis, providing a modified version of the prophets for peptide validation and protein inference. To this date, Philosopher is the only proteomics toolkit that allows you to process and analyze close and open search results.

- Mass spec data file conversion with msconvert.

- Database downloading and formatting.

- Peptide Spectrum Matching with Comet.

- Peptide assignment validation with PeptideProphet.

- Multi-level integrative analysis with iProphet.

- PTM site localization with PTMProphet.

- Protein inference with ProteinProphet.

- FDR filtering with custom algorithms.

  - Two-dimensional filtering for simultaneous control of PSM and Protein FDR levels.
  - Sequential FDR estimation for large data sets using filtered PSM and proteins lists.

- Label-free quantification via Spectral counting and MS1 Quantification.

- Labeling-based quantification using TMT isobaric tags.

- Quantification based on functional protein groups.

- Multi-level detailed reports including peptides, ions and proteins.

- Mass spec identification file conversion with idconvert.


## How to Download
Download the latest version [here](https://github.com/nesvilab/philosopher/releases/latest)


## How to Use
A simple [tutorial](https://github.com/Nesvilab/philosopher/wiki) is also provided with an extensive example on how to use Philosopher.


## Documentation
Check the [documentation](https://github.com/Nesvilab/philosopher/wiki/Philosopher-Basics) for more details about the available commands.


## Questions, requests and bug reports
If you have any questions, remarks, or requests, or if you found a bug, please use the [Issue tracker](https://github.com/nesvilab/philosopher/issues).


## How to cite
Still in progress!


## About the developers
[Felipe da Veiga leprevost](http://www.leprevost.com.br)

Sarah Haynes

Alexey Nesvizhskii's [nesvilab](http://www.nesvilab.org/)



## Questions, requests and bug reports
If you have any questions, remarks, requests or if you found a bug, please use the [Issue tracker](https://github.com/nesvilab/philosopher/issues).


## License
GPL 3
