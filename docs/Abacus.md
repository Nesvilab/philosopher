Aggregate data from multiple experiments and adjusts label-free quantification to accurately account for peptides shared across multiple proteins


## Usage

`philosopher abacus [flags] [folders]`


## flags

`--peptide`

Global level peptide report (needs a combined.pep.xml file).

`--protein`

Global level protien report (needs a combined.prot.xml file).

`--labels`

Indicates whether the data sets include TMT labels or not.

`--pepProb`

Minimum peptide probability (default 0.5).

`--picked`

Apply the picked FDR algorithm before the protein scoring.

`--prtProb`

Minimum protein probability (default 0.9).

`--razor`

Use razor peptides for protein FDR scoring.

`--tag`

Decoy tag (default "rev_").

`--uniqueonly`

Report TMT quantification based on only unique peptides.

`--reprint`

Create abacus reports using the Reprint format


## Example

Aggregating data from 3 different experiments, in 3 different workspaces

`philosopher abacus control/ treatment_1/ treatment_2/`

Aggregating data from 3 different experiments, in 3 different workspaces and using a pre existing protXML combined file.

`philosopher abacus --protein control/ treatment_1/ treatment_2/`


## FAQ

_What exactly do I need to do before running Abacus ?_

You need to work on each individual experiment workspace before running Abacus. Each folder containing individual experimental data must be converted to a Workspace and must have its data analyzed by the filter command.

_I don't have a combined pepXML file, how do I get one?_

You need to execute iProphet using all interact.pep.xml files from each individual folder you are analyzing.

_I don't have a combined protXML file, how do I get one?_

You need to execute ProteinProphet using all pepXML files from each individual folder you are analyzing.

_Where should I execute the abacus command ?_

The command should be execute one level above the experimental data

_This seems to be a lot of work, isn't there any workaround ?_

Yes, take a look at the [Pipeline](pipeline.md) command.
