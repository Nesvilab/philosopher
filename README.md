# Philosopher
A data processing toolkit for shotgun proteomics.

![Golang](https://img.shields.io/badge/Go-1.8.0-blue.svg) ![Version](https://img.shields.io/badge/version-rc5-blue.svg)


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


## How to Download
Download the latest version [here](https://github.com/prvst/philosopher/releases/latest)


## How to Use
A simple [tutorial](tutorial.md) is also provided with an extensive example on how to use Philosopher.


## Documentation
Check the [documentation](documentation.md) for more details about the available commands.


## Questions, requests and bug reports
If you have any questions, remarks, requests or if you found a bug, please use the [Issue tracker](https://github.com/prvst/philosopher/issues).

## How to cite
Still in progress!

## About the developers
Another great work from Alexey Nesvizhskii's [nesvilab](http://www.nesvilab.org/)

## License
GPL 3
