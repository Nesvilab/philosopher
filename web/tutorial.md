This guide will show you how Philosopher can be used for a complete proteomics data analysis. We are going to review a complete example, starting with a converted raw file and finishing with quantified protein reports.

The tutorial will not discuss the peculiarities of each database search parameter or the advantages of each FDR scoring algorithm, the goal here is to demonstrate how simple it is to run an analysis. After mastering how to execute Philosopher, you can work on the more appropriate settings.

## What are the Basic Steps ?
Your data analysis needs to follow a certain logic of execution. Each step will do something different with your data and give you something back. You do not necessarily need to execute all the commands like that, some of the commands like the Prophets can be used with data processed somewhere else, for example. Take this as a "default" order to get your results.

1. Create a Workspace
2. Annotate or download a database
3. Database search using Comet
4. PeptideProphet
5. InterProphet or PTMProphet (optional) <-- these will update your PeptideProphet results
6. ProteinProphet
7. Filter
8. Freequant (MS1 quantification, optional)
9. Labelquant (TMT quantification, optional)
9. Report
10. Backup and Clean

Philosopher also provides some extra commands in case you need to analyze your data in a different way:

1. Cluster (Protein report based on protein clusters, optional)
2. Abacus  (Combined analysis of LC-MS/MS and cohorts, optional)


## Before we start
As an example data set, I am going to use a _Pyrococcus furiosus_ sample described [here](http://pubs.acs.org/doi/abs/10.1021/pr300055q). To begin with the analysis we will need a converted RAW, so you can download the file available via [ProteomeXchange](http://proteomecentral.proteomexchange.org/cgi/GetDataset?ID=PXD001077) and use your favorite converter to get a mzXML file. __Remember to convert the files using 64-bit encoding and peak picking__.

If you need help with the commands, you can run them using the `--help` flag, this will bring a description of all the available flags for each command. For more details check the [Documentation](documentation.md).


### 1. Creating a Workspace
A Workspace is a directory containing your data to be analyzed. To create one you need to execute the __workspace__ command followed by the `--init` parameter, this will transform your current directory into a Philosopher workspace, all processed data will be sorted in a special binary format for easy and fast access.

`philosopher workspace --init`

From now on, all the following steps should be executed inside the workspace.


### 2. Creating a protein database
For the first step we will download and format a database file using the _database_ command, but first we need to find the Proteome ID (PID) for our organism. Searching the [UniProt proteome list](http://www.uniprot.org/proteomes), we can see that _Pyrococcus furiosus_ proteome ID is **UP000001013**, so let's prepare the database file:

`philosopher database --id UP000001013 --contam --reviewed`

We are only going to need reviewed sequences and we are adding decoys with the _rev_ tag. After executing the line above, you will see that a new file was created on you local folder, the file name contains the current date, the _td_ indicates that this file is a target-decoy database. The database is also converted into a proper data structure inside the workspace in order to be used later.


### 3. Database Search using Comet
For the database search, first we need to print a parameter file:

`philosopher comet --print`

After setting your parameter file (you can use the publication cited above as a reference), we can launch the search:

`philosopher comet --param comet.params Velos005137.mzML`

The search should be done in a few minutes since this is a small data file. As a result you now have a new file called _Velos005137.pep.xml_


### 4. PeptideProphet
The next step is to run PeptideProphet to validate the peptide assignments:

`philosopher peptideprophet --database 2016-11-28-td-186497.fas --accmass --decoy rev_ --decoyprobs --nonparam Velos005137.pep.xml`

You will see on the terminal screen the same output you normally see when running PeptideProphet from TPP. If everything runs OK, you will have a new file called _interact-Velos005137.pep.xml_ as a result.


### 5. ProteinProphet
Now we can run the protein inference step and generate a protXML file:

`philosopher proteinprophet interact-Velos005137.pep.xml`


### 6. Filtering and estimating FDR levels
Now we have all necessary files to filter our data using the FDR approach:

`philosopher filter --pepxml interact-Velos005137.pep.xml`

Running the above command with only a pepXML will give you the current levels for the pepXML file only, if you include a protXML file, Philosopher will use protein inference information to make the FDR score more precise:

`philosopher filter --pepxml interact-Velos005137.pep.xml --protxml interact-Velos005137.prot.xml`

The _filter_ algorithm can be applied in different ways, using different strategies or even applying razor peptides. here it's up to you to choose the best method to analyze your data. The program will output in the screen a log with the scoring results, all the processed data will be stored inside your workspace for further analysis.


### 7. Reporting
Now we can inspect the experiment results by printing the PSM, peptide and protein reports:

`philosopher report`


### Backup
Finally, for the last step, backup your data just in case you wish to print the reports again later.

`philosopher workspace --backup`


### Concluding remarks
Here I just demonstrated how easy it is to run a complete proteomics analysis using Philosopher. By providing access to advanced analysis software, combined with custom processing algorithms, it is possible to go from converted raw files to protein reports in a few minutes.
