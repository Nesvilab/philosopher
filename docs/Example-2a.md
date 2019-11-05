This tutorial will show you how to use Philosopher for a complete open search proteomics analysis. 
**Philosopher can also be used on Windows, though the commands in this tutorial are formatted for GNU/Linux**

## What are the basic steps?
The commands in this tutorial should be executed in a particular order. Consider this the "default" order in which to perform an analysis:

1. Create a workspace
2. Download a database
3. Search with MSFragger
4. PeptideProphet
5. ProteinProphet
6. Filter
7. Report


## Before we start
As an example data set, we will use a publicly-available LC-MS data file from a human protein sample described [in this publication](https://www.ncbi.nlm.nih.gov/pubmed?term=29718670). Download the file _06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.raw_ from the dataset [FTP location](ftp://ftp.pride.ebi.ac.uk/pride/data/archive/2018/05/PXD008952) (the full listing is [here](http://proteomecentral.proteomexchange.org/cgi/GetDataset?ID=PXD008952)).

You can choose to use the .raw spectral format, or you can convert it to the mzML format (needed for quantification). A tutorial on raw file conversion can be found [here](https://msfragger.nesvilab.org/tutorial_convert.html).

If you need help with the commands, you can run them using the `--help` flag (e.g. `philosopher workspace --help`), which will provide a description of all available flags for each command. For more details, see the [documentation](documentation.md).


### 1. Create a workspace
(Note: use the full path to the Philosopher binary file in place of _philosopher_ in the following steps.)
Place the _06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.raw_ in a new folder, which we will call the 'workspace'. We will create the workspace with the Philosopher __workspace__ command, which will enable the program to store processed data in a special binary format for quick access.

Inside your workspace folder, open a new terminal window and run this command:
`philosopher workspace --init`

From now on, all steps should be executed inside this same directory.


### 2. Download a protein database
For the first step we will download and format a database file using the _database_ command, but first we need to find the Proteome ID (PID) for our organism. Searching the [UniProt proteome list](http://www.uniprot.org/proteomes), we can see that the _Homo sapiens_ proteome ID is _UP000005640_, so let's prepare the database file with the following command:

`philosopher database --id UP000005640 --reviewed --contam`

Philosopher will retreive all reviewed protein sequences from this proteome, add common contaminants, and generate decoy sequences labeled with the tag _rev\__.

You should see that a new file was created in the workspace, and that the file name contains the current date and a _td_ label indicating that this file contains both target and decoy sequences. The database file also includes a _rev_ label indicating only reviewed sequences are included.)

### 3. Perform a database search with MSFragger
(Note: use the full path to the MSFragger.jar file in place of _MSFragger.jar_ in the following steps.)

Run `java -jar MSFragger.jar --config` to print three MSFragger parameter files (closed, nonspecific, and open). 

In the _open_fragger.params_ file, update the _database_name_ parameter to the name of the database file we downloaded in the previous step (e.g. _2019-11-04-td-rev-UP000005640.fas_). You can also change the _calibrate_mass_ parameter from 2 to 0 to speed up the search even more.

Launch the search by running: `java -Xmx32g -jar MSFragger.jar open_fragger.params 06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.raw`. (Adjust the `-Xmx` flag to the appropriate amount of RAM for your computer.) 

The search should be done in a few minutes or less. The search hits are now stored in a file called _06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pepXML_.


### 4. PeptideProphet
The next step is to validate the peptide hits with PeptideProphet:

`philosopher peptideprophet --database 2019-11-04-td-rev-UP000005640.fas --nonparam --expectscore --decoyprobs --masswidth 1000.0 --clevel -2 06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pepXML`

This will generate a new file called _interact-06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pep.xml_.


### 5. ProteinProphet
Next, perform protein inference and generate a protXML file:

`philosopher proteinprophet --maxppmdiff 2000000 interact-06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pep.xml`


### 6. Filter and estimate FDR
Now we have all necessary files to filter our data using the FDR approach:

`philosopher filter --pepxml interact-06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pep.xml --protxml interact.prot.xml`

The **filter** algorithm can be applied in many different ways, use the `--help` flag and choose the best method to analyze your data. Scoring results will be shown in the console, and all processed data will be stored in your workspace for further analysis.


### 7. Report the results
Now we can inspect the experiment results by printing the PSM, peptide, and protein reports:

`philosopher report`


### Backup
As an optional last step, backup your data in case you wish to print the reports again later.

`philosopher workspace --backup`


### Concluding remarks
We've demonstrated how to run a complete proteomics analysis using Philosopher. By providing easy access to advanced analysis software and custom processing algorithms, protein reports can be obtained from LC-MS files in just a few minutes.
