This guide will show you how Philosopher can be used for a complete proteomics data analysis. Below you can review a few examples, starting with a converted raw file and finishing with quantified protein reports.

The tutorial will not discuss the peculiarities of each database search parameter or the advantages of each FDR scoring algorithm, the goal here is to demonstrate how simple it is to run an analysis. After mastering how to execute Philosopher, you can work on the more appropriate settings.


## Philosopher Basics
Philosopher is composed by atomic commands that run independently from each other. These commands can be executed separately and in different moments, you can run one command now and the remaining in a different day, for example. When running Philosopher you should not expect to see an output for each command, your data is stored locally in a custom binary format, so each step will update the information aggregated to it. Once that you reach the end of your analysis, the `report` command can be executed to create the report files.


## What is a workspace ?
Philosopher works with the concept of a workspace. A workspace is a directory on your computer that holds the data that should be analyzed. It's OK to have other files in there that are not related to the process, they will not be altered. You can name your folder the way you want (avoid spaces and special characters), but it's advised that you use the same description you are already using for the data set. If you have multiple data setrs, you should arrange them side by side, like the example below:

```
.
└── workdir
    ├── control
    │   └── control.mzML
    ├── test_1
    │   └── test1.mzML
    └── test_2
        └── test2.mzML

4 directories, 3 files
```

In this example we have a directory called `workdir`, and inside we have our directories containing our mass spec. data, each in it's individual folder. Each one for these folder will be converted into a Philosopher `workspace`, and the whole analysis will happen inside them.


### How to organize and process a large data set ?
Some experiments are composed by dozens if not more data sets, but that doesn't mean you will need to process each one of them individually. Philosopher provides a command called [Pipeline](pipeline.md) that allows you to automate the entire analysis. The only task you will need to do is the directory organization, you might consider using a scripting language to assist you.


## What is the analysis Workflow ?
Your data analysis needs to follow a certain logic of execution. Each step will do something different with your data and give you something back. You do not necessarily need to execute all the commands like that, some of the commands like the Prophets can be used with data processed somewhere else, for example. Take this as a "default" order to get your results.

1. Create a Workspace
2. Download or Annotate a database
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


## Commands dependency
Some commands they expect to find processed data, some expect "raw" information, this is how the commands depend on each other:

* Depend on database search results
  * PeptideProphet
* Depend on PeptideProphet
  * InterProphet, PTMProphet, ProteinProphet, Filter
* Depend on Filter
  * Freequant, Labelquant, Abacus, CLuster, Report
* No Dependency
  * Workspace, Database, Comet


### Examples

[Example of a simple data analysis](example_1.md)
