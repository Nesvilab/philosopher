# Philosopher
A data processing toolkit for shotgun proteomics.

![Golang](https://img.shields.io/badge/Go-1.9.2-blue.svg) ![Version](https://img.shields.io/badge/version-beta-blue.svg) [![https://philosopher-toolkit.slack.com](https://img.shields.io/badge/slack-channel-blue.svg)](https://philosopher-toolkit.slack.com?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

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

- Clustering analysis for proteomics results.

- Multi-level detailed reports including peptides, ions and proteins.

- Mass spec identification file conversion with idconvert.


## How to Download
Download the latest version [here](https://github.com/prvst/philosopher/releases/latest)


## How to Use
A simple [tutorial](tutorial.md) is also provided with an extensive example on how to use Philosopher.


## Documentation
Check the [documentation](documentation.md) for more details about the available commands.


## Questions, requests and bug reports
If you have any questions, remarks, requests or if you found a bug, please use the [Issue tracker](https://github.com/prvst/philosopher/issues).


## License
GPL 3
