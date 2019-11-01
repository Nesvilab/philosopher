This tutorial will show you how to use Philosopher for complete proteomics data analysis. We will start with a converted raw file and finishing with quantified protein reports.

The tutorial will not discuss the peculiarities of each database search parameter or the advantages of each FDR scoring algorithm, the goal here is to demonstrate how simple it is to run an analysis. After mastering how to execute Philosopher, you can work on the more appropriate settings.

**Philosopher can also be used on Windows, though the commands in this tutorial are formatted for GNU/Linux**

## What are the basic steps?
The commands in this tutorial should be executed in a particular order. You do not necessarily need to execute all the commands, some (like the Prophets) can be used with data processed somewhere else. Consider this the "default" order in which to perform an analysis:

1. Create a workspace
2. Annotate or download a database
3. Database search using MSFragger
4. PeptideProphet
5. ProteinProphet
6. Filter
7. Report


Philosopher provides additional commands in case you need to analyze your data in a different way:

1. Freequant (label-free quantification based on precursor intensity, performed after filtering)
2. Labelquant (quantification from isobaric labels, performed after filtering)
3. Cluster (Protein report based on protein clusters, optional)
4. Abacus  (Combined analysis of LC-MS/MS and cohorts, optional)


## Before we start
As an example data set, we will use a publicly-available LC-MS data file from a _Pyrocuccus furiosus_ (extremophile species of Archea) sample described [in this publication](https://pubs.acs.org/doi/abs/10.1021/pr300055q). Download the file _Velos005137.raw_ from the dataset [FTP location](ftp://ftp.pride.ebi.ac.uk/pride/data/archive/2014/06/PXD001077) (the full listing is [here](http://proteomecentral.proteomexchange.org/cgi/GetDataset?ID=PXD001077)).

You can choose to use the .raw spectral format, or you can convert it to the mzML format (which is needed for quantification). A tutorial on raw file conversion can be found [here](https://msfragger.nesvilab.org/tutorial_convert.html).

If you need help with the commands, you can run them using the `--help` flag (e.g. `philosopher workspace --help`), which will provide a description of all available flags for each command. For more details, see the [documentation](documentation.md).


### 1. Create a workspace
(Note: use the full path to the Philosopher binary file in place of _philosopher_ in the following steps.)
Place the _Velos005137.raw_ in a new folder, which we will call the 'workspace'. We will create the workspace with the Philosopher __workspace__ command, which will enable the program to store processed data in a special binary format for quick access.

Inside your workspace folder, open a new terminal window and run this command:
`philosopher workspace --init`

From now on, all steps should be executed inside this same directory.


### 2. Download a protein database
For the first step we will download and format a database file using the _database_ command, but first we need to find the Proteome ID (PID) for our organism. Searching the [UniProt proteome list](http://www.uniprot.org/proteomes), we can see that the _Pyrococcus furiosus_ proteome ID is _UP000001013_, so let's prepare the database file with the following command:

`philosopher database --id UP000001013 --contam`

Philosopher will retreive all protein sequences from this proteome, add common contaminants, and generate decoy sequences labeled with the tag _rev\__.

You should see that a new file was created in the workspace, and that the file name contains the current date and a _td_ label indicating that this file contains both target and decoy sequences. (For well-studied organisms, such as _Homo sapiens_ and _Mus musculus_, the `--reviewed` parameter should be used to include only reviewed sequences, these database files will contain a _rev_ label.)

### 3. Perform a database search with MSFragger
(Note: use the full path to the MSFragger.jar file in place of _MSFragger.jar_ in the following steps.)

Run `java -jar MSFragger.jar --config` to print three MSFragger parameter files (closed, nonspecific, and open). 

In the _closed_fragger.params_ file, update the _database_name_ parameter to the name of the database file we downloaded in the previous step (e.g. _2019-10-31-td-UP000001013.fas_). You can also change the _calibrate_mass_ parameter from 2 to 0 to speed up the search even more.

Launch the search by running: `java -Xmx32g -jar MSFragger.jar closed_fragger.params Velos005137.pepXML`. (Adjust the `-Xmx` flag to the appropriate amount of RAM for your computer.) 

The search should be done in a few minutes or less. The search hits are now stored in a file called _TGR_02603.pepXML_.


### 4. PeptideProphet
The next step is to validate the peptide hits with PeptideProphet:

`philosopher peptideprophet --database 2019-10-30-td-UP000001013.fas --ppm --accmass --expectscore --decoyprobs --nonparam Velos005137.pepXML`

This will generate a new file called _interact-TGR_02603.pep.xml_.


### 5. ProteinProphet
Next, perform protein inference and generate a protXML file:

`philosopher proteinprophet interact-Velos005137.pep.xml`


### 6. Filter and estimate FDR
Now we have all necessary files to filter our data using the FDR approach:

`philosopher filter --pepxml interact-Velos005137.pep.xml`

Running the above command with only a pepXML will give you the current levels for the pepXML file only. If you include a protXML file, Philosopher will use protein inference information to make the FDR score more precise:

`philosopher filter --pepxml interact-Velos005137.pep.xml --protxml interact.prot.xml`

The **filter** algorithm can be applied in many different ways, use the `--help` flag and choose the best method to analyze your data. Scoring results will be shown in the console, and all processed data will be stored in your workspace for further analysis.


### 7. Report the results
Now we can inspect the experiment results by printing the PSM, peptide and protein reports:

`philosopher report`


### Backup
As an optional last step, backup your data in case you wish to print the reports again later.

`philosopher workspace --backup`


### Concluding remarks
We've demonstrated how to run a complete proteomics analysis using Philosopher. By providing easy access to advanced analysis software and custom processing algorithms, protein reports can be obtained from raw LC-MS files in just a few minutes.
