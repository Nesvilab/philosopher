A custom algorithm for MS/MS data filtering and multi-level false discovery rate estimation.


## Usage

`philosopher filter [flags]`


## Flags


`--ion`

Peptide ion FDR level (default 0.01).

`--mapmodels`

Map modifications acquired by an open search.

`--models`

Print model distribution for the analyzed pepXML.

`--pep`

Peptide FDR level (default 0.01)

`--pepProb`

top peptide probability threshold for the FDR filtering (default 0.7)

`--pepxml`

Path to the pepXML file or the path to a directory containing a set of pepXML files.

`--picked`

Apply the picked FDR algorithm before the protein scoring.

`--prot`

Protein FDR level (default 0.01).

`--protProb`

Protein probability threshold for the FDR filtering (not used with the razor algorithm) (default 0.5).

`--protxml`

Path to the protXML file path.

`--psm`

PSM FDR level (default 0.01).

`--razor`

Uses razor peptides for protein FDR scoring.

`--sequential`

Alternative algorithm that estimates FDR using both filtered PSM and Protein lists to boost identifications.

`--tag`

The decoy prefix used on decoy sequences (default "rev_").

`--weight`

Threshold for defining peptide uniqueness (default 1).


## Examples

Process a single pepXML file using the standard filter values.

`philosopher filter --pepxml interact.pepxml`

Process several pepXML files and one protXML file using the standard filter values.

`philosopher filter --pepxml results/ --protxml interact.protxml`

Process all pepXML files and the protXML file called _sample.protxml_ found in the _result_ folder. The program will employ the razor algorithm during the FDR filtering and at the end will save the analyzed data in binary format for posterior processing and reporting.

`philosopher filter --pepxml results/ --protxml results/sample.protxml --razor`
