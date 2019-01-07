Peptide Spectrum Matching using the Comet database search engine.


## Version
Comet release 2017.01 rev. 4


## Usage

`philosopher comet [flags] [files]`


## Flags

`--print`

Print the default parameter file to the local directory.

`--param`

Points the location of the parameter file for the analysis.


## Example

Execute Comet using the parameter file _param.txt_, and using 2 converted raw files _sample1_ and _sample2_. Comet will analyze them one at a time using the specified parameters.

`philosopher comet --param param.txt sample1.mzml sample2.mzml`
