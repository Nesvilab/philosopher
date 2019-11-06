This tutorial will show you how to use Philosopher for a complete proteomics data analysis, starting from a raw LC-MS file and ending with quantified protein reports.
**Philosopher can also be used on Windows, though the commands in this tutorial are formatted for GNU/Linux**

## What are the basic steps?
The commands in this tutorial should be executed in a particular order. Consider this the "default" order in which to perform an analysis:

1. Create a workspace
2. Download a database
3. Search with MSFragger
4. PeptideProphet
5. ProteinProphet
6. Filter
7. Quantify
8. Report

## Before we start
For this tutorial, we will use a publicly-available LC-MS data file from a TMT 10-plex phosphorylation-enriched human cell line sample described [in this publication](https://pubs.acs.org/doi/10.1021/acs.jproteome.8b00165). Download the file _06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.raw_ from the dataset [FTP location](ftp://ftp.pride.ebi.ac.uk/pride/data/archive/2018/05/PXD008952) (the full listing is [here](http://proteomecentral.proteomexchange.org/cgi/GetDataset?ID=PXD008952)).

You can choose to convert it to the mzML format or use the .raw file directly and skip the quantification (quantification of .raw files is not supported). A tutorial on raw file conversion can be found [here](https://msfragger.nesvilab.org/tutorial_convert.html).

For additional help on any of the Philosopher commands, you can use the `--help` flag (e.g. `philosopher workspace --help`), which will provide a description of all available flags for each command. For more details, see the [documentation](documentation.md).


### 1. Create a workspace
(Note: use the full path to the Philosopher binary file in place of _philosopher_ in the following steps.)
Place the _06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.raw_ in a new folder, which we will call the 'workspace'. We will create the workspace with the Philosopher __workspace__ command, which will enable the program to store processed data in a special binary format for quick access.

Inside your workspace folder, open a new terminal window and run this command:  
`philosopher workspace --init`

From now on, all steps should be executed inside this same directory.


### 2. Download a protein database
For the first step we will download and format a database file using the _database_ command, but first we need to find the Proteome ID (PID) for our organism. Searching the [UniProt proteome list](http://www.uniprot.org/proteomes), we can see that the _Homo sapiens_ proteome ID is _UP000005640_, so let's prepare the database file with the following command:  
`philosopher database --id UP000005640 --contam`

Philosopher will retreive all reviewed protein sequences from this proteome, add common contaminants, and generate decoy sequences labeled with the tag _rev\__.

You should see that a new file was created in the workspace, and that the file name contains the current date and a _td_ label indicating that this file contains both target and decoy sequences. The database name also includes a _rev_ label indicating that only reviewed protein sequences are included.

### 3. Perform a database search with MSFragger
(Note: use the full path to the MSFragger.jar file in place of _MSFragger.jar_ in the following steps.)

Run `java -jar MSFragger.jar --config` to print three MSFragger parameter files (closed, nonspecific, and open). 

In the _closed_fragger.params_ file, update the _database_name_ parameter to the name of the database file we downloaded in the previous step (e.g. _2019-11-04-td-rev-UP000005640.fas_). Below _variable_mod_02_, add a third variable modification for the TMT isobaric label on the peptide N-terminus: _variable_mod_03 = 229.162932 n^_.

Toward the bottom of the parameter file, change the _add_K_lysine_ value from `0.000000` to `229.162932`.
You can also change the _calibrate_mass_ parameter near the top of the file from `2` to `0` to speed up the search even more.

Launch the search by running:  
`java -Xmx32g -jar MSFragger.jar closed_fragger.params 06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.raw`. (Adjust the `-Xmx` flag to the appropriate amount of RAM for your computer.) 

The search should be done in a few minutes or less. The search hits are now stored in a file called _06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pepXML_.

(Note: for simplicity, we haven't searched for phosphorylation, but you can by adding _variable_mod_04 = 79.966331 STY_ to the MSFragger parameter file.)


### 4. PeptideProphet
The next step is to validate the peptide hits with PeptideProphet:  
`philosopher peptideprophet --database 2019-11-04-td-rev-UP000005640.fas --ppm --accmass --expectscore --decoyprobs --nonparam 06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pepXML`

This will generate a new file called _interact-06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pep.xml_.


### 5. ProteinProphet
Next, perform protein inference and generate a protXML file:  
`philosopher proteinprophet interact-06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pep.xml`


### 6. Filter and estimate FDR
Now we have all necessary files to filter our data using the FDR approach:  
`philosopher filter --razor --pepxml interact-06_CPTAC_TMTS1-NCI7_P_JHUZ_20170509_LUMOS.pep.xml --protxml interact.prot.xml`

The **filter** algorithm can be applied in many different ways, use the `--help` flag and choose the best method to analyze your data. Scoring results will be shown in the console, and all processed data will be stored in your workspace for further analysis.


### 7. Perform label-based quantification
After filtering, label-free or label-based quantification can be performed. For this tutorial we will use `labelquant` to quantify this TMT-10plex sample, but the `freequant` command can be used to perform label-free MS1 quantification (see the [wiki](https://github.com/Nesvilab/philosopher/wiki/Freequant) for details).

Perform the quantification by running:  
`philosopher labelquant --plex 10 --dir .` , where the `.` indicates the current workspace.

Optionally, you can provide an _annotation.txt_ file to indicate the sample/replicate that corresponds to each TMT channel. Make a new text file and fill it according to the following format:

>126 control_1  
127N treated_1  
127C control_2  
128N treated_2  
128C control_3  
129N treated_3  
129C control_4  
130N treated_4  
130C control_5  
131N treated_5  

Save this new file as `annotation.txt`, then run the quantification using the _annotation.txt_ file:  
`philosopher labelquant --plex 10 --dir . --annot annotation.txt`  
(For other types of TMT labeling [6, 11, or 16-plex], use the appropriate `--plex` value.)



### 8. Report the results
Now we can inspect the results by printing the PSM, peptide, and protein reports:  
`philosopher report`


### Backup
As an optional last step, backup your data in case you wish to print the reports again later:  
`philosopher workspace --backup`


### Concluding remarks
We've demonstrated how to run a complete proteomics analysis with TMT quantification using Philosopher. By providing easy access to advanced analysis software and custom processing algorithms, protein reports can be obtained from LC-MS files in just a few minutes.
