Statistical validation of protein identification based on peptide assignment to MS/MS spectra


## Version

ProteinProphet v5.01


## Usage

`philosopher proteinprophet [flags] [files]`


## Flags

`--accuracy`

Equivalent to MINPROB0.

`--allpeps`

Consider all possible peptides in the database in the confidence model.

`--asap`

Compute ASAP ratios for protein entries (ASAP must have been run previously on interact dataset).

`--asapprophet`

Compute ASAP ratios for protein entries (ASAP must have been run previously on all input interact datasets with mz/XML raw data format).

`--confem`

Use the EM to compute probability given the confidence.

`--delude`

Do NOT use peptide degeneracy information when assessing proteins.

`--maxppmdiff`

Maximum peptide mass difference in PPM (default 20)

`--fpkm`

Model protein FPKM values.

`--glyc`

Highlight peptide N-glycosylation motif.

`--icat`

Highlight peptide cysteines.

`--instances`

Use Expected Number of Ion Instances to adjust the peptide probabilities prior to NSP adjustment.

`--iprophet`

Input is from iProphet.

`--logprobs`

Use the log of the probabilities in the Confidence calculations.

`--minindep`

Minimum percentage of independent peptides required for a protein (default=0).

`--minprob`

PeptideProphet probability threshold (default=0.05).

`--mufactor`

Fudge factor to scale MU calculation (default 1).

`--nogroupwts`

Check peptide's Protein weight against the threshold (default: check peptide's Protein Group weight against threshold).

`--nooccam`

Non-conservative maximum protein list.

`--noprotlen`

Do not report protein length.

`--normprotlen`

Normalize NSP using Protein Length.

`--output`

Output name (default "interact.prot.xml").

`--protmw`

Get protein mol weights.

`--refresh`

Import manual changes to AAP ratios (after initially using ASAP option).

`--softoccam`

Peptide weights are apportioned equally among proteins within each Protein Group (less conservative protein count estimate).

`--unmapped`

Report results for UNMAPPED proteins.


## Example

Execute a standard analysis on a pepXML file called sample.pepxml.

`philosopher proteinprophet sample.pepxml`


## FAQ

_Do I need TPP installed for running this ?_

No
