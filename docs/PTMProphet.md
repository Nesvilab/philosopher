PTM site localization


## Version

PTMProphet v5.01


## Usage

`philosopher ptmprophet [flags] [file]`


## Flags

`--keepold`

Retain old PTMProphet results in the pepXML file

`--massdiffmode`

Use the Mass Difference and localize

`--minprob`

Use specified minimum probability to evaluate peptides

`--mztol`

Use specified +/- MS2 mz tolerance on site specific ions (default 0.1)

`--output`

Output prefix file name

`--ppmtol`

Use specified +/- MS1 ppm tolerance on peptides which may have a slight offset depending on search parameters (default 1)

`--verbose`

Produce Warnings to help troubleshoot potential PTM shuffling or mass difference issues

## Example

Execute a standard analysis on a pepXML file called sample.pepxml.

`philosopher proteinprophet sample.pepxml`


## FAQ

_Do I need TPP installed for running this ?_

No
