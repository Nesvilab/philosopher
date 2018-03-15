Multi-level integrative analysis of shotgun proteomic data


## Version

iProphet v5.01


## Usage

`philosopher iprophet [flags] [files]`


## Flags

`--cat`

Specify file listing peptide categories.

`--decoy`

Specify the decoy tag.

`--inNonsi`

Do not use NSI model.

`--inonrs`

Do not use NRS model.

`--inonsm`

Do not use NSM model.

`--inonsp`

Do not use NSP model.

`--isharpnse`

Use more discriminating model for NSE in SWATH mode (default: Enabled, use NONSE to disable).

`--length`

Use Peptide Length model.

`--minProb`

Specify minimum probability of results to report.

`--nNonss`

Do not use NSS model.

`--nofpkm`

Do not use FPKM model.

`--nonse`

Do not use NSE model.

`--output`

Specify output name (default "iproph.pep.xml").

`--threads`

Specify threads to use (default 1).


## Example

Run iProphet on several files pepXML files using 6 threads.

`philosopher iprophet --threads 6 combined_samples_1.pepxml combined_samples_2.pepxml combined_samples_3.pepxml`


## FAQ

_Do I need TPP installed for running this ?_

No
