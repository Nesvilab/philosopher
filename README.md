# Philosopher
A data processing toolkit for shotgun proteomics.

![Golang](https://img.shields.io/badge/Go-1.8.1-blue.svg)
![Version](https://img.shields.io/badge/version-1.0-blue.svg)
[![https://philosopher-toolkit.slack.com](https://img.shields.io/badge/slack-channel-blue.svg)](https://philosopher-toolkit.slack.com?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Features
Philosopher provides easy access to third-party tools and custom algorithms allowing users to develop proteomics analysis, from Peptide Spectrum Matching to annotated protein reports. Philosopher is also tunned for Open Search analysis, providing a modified version of the prophets for peptide validation and protein inference.

- Database downloading and formatting.

- Peptide Spectrum Matching with Comet.

- Peptide assignment validation with PeptideProphet.

- Multi-level integrative analysis with iProphet.

- Protein inference with ProteinProphet.

- Open Search data validation.

- FDR filtering with custom algorithms.

  - Two-dimensional filtering for simultaneous control of PSM and Protein FDR levels.
  - Sequential FDR estimantion using filtered PSM and proteins lists.
  - PickedFDR for scalable estimations.
  - Razor peptides determination for better quantification and interpretationinterpretation.

- Label-free Quantification (Spectral Count and MS1 intensities).

- Isotope label quantification (TMT).

- Clustering analysis for proteomics results.

- Detailed multi-level reports with optional funtional annotation.


Access [Philosopher website](https://prvst.github.io/philosopher/) for more information on how to use and download the program.
