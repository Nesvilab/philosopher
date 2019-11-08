# Philosopher
A data processing toolkit for shotgun proteomics.

[![Release](https://img.shields.io/github/release/nesvilab/philosopher.svg?style=for-the-badge)](https://github.com/nesvilab/philosopher/releases/latest)
![Golang](https://img.shields.io/badge/Go-1.13.1-blue.svg?style=for-the-badge)

## Features
Philosopher provides easy access to third-party tools and custom algorithms allowing users to develop complete proteomics analyses, from peptide spectrum matching to annotated protein reports. Philosopher is also tuned for Open Search analysis, providing a modified version of the Prophets for peptide validation and protein inference. To date, Philosopher is the only proteomics toolkit that allows you to process and analyze both closed and open search results.


- Database downloading and formatting.

- Peptide Spectrum Matching with Comet.

- Peptide assignment validation with PeptideProphet.

- Multi-level integrative evidence analysis with iProphet.

- PTM site localization with PTMProphet.

- Protein inference with ProteinProphet.

- Open search result validation.

- FDR filtering with custom algorithms.

  - Two-dimensional filtering for simultaneous control of PSM and Protein FDR levels.
  - Sequential FDR estimation for large data sets using filtered PSM and proteins lists.
  - PickedFDR for scalable estimations.
  - Razor peptide assignment for better quantification and interpretation.


- Label-free quantification via spectral counting and MS1 intensities.

- Labeling-based quantification using TMT isobaric tags.

- Clustering analysis for proteomics results.

- Multi-level detailed reports including PSMs, ions, peptides, and proteins.




## How to Download
Download the latest version [here](https://github.com/nesvilab/philosopher/releases/latest)


## How to Use
A simple [tutorial](tutorial.md) is provided with a complete example on Philosopher usage.


## Documentation
Check the [documentation](documentation.md) for more details about the available commands.


## Questions, requests and bug reports
If you have any questions, remarks, or requests, or if you found a bug, please use the [Issue tracker](https://github.com/nesvilab/philosopher/issues).

## Analytics
Check the [analytics](analytics.md) session for more information on the program usage report.

## How to cite
Still in progress!


## About the developers
[Felipe da Veiga leprevost](http://www.leprevost.com.br)
Sarah Haynes
Alexey Nesvizhskii's [nesvilab](http://www.nesvilab.org/)


## License
Since Philosopher is not officially released, the binaries provided here have no license. Once released, the software will be licensed under the GPL 3 license.
