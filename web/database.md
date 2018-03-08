The _Database_ command can be used to prepare a protein FASTA file for the following analysis. The database can be downloaded using the UniProt Proteome ID or by using a pre-formed FASTA file.


## Usage

`philosopher database [flags]`


## Flags

`--add`

Allows you to add one or more custom sequences to your database. Only UniProt FASTA formatting is allowed.

`--annotate`

Process a ready-to-use database.

`--contam`

Add 116 common contaminants found in LC-MS/MS experiments. More information can be found [here](http://www.thegpm.org/crap/).

`--custom`

Skips the downloading of a fresh database and use an existing one instead. The custom file will also be used to create decoys and contaminants if desired.

`--enzyme`

The name of the enzyme for the digestion. The options are:

* trypsin
* lys_c
* lys_n
* chymotrypsin

The default option is _trypsin_.

`--id`

The Proteome ID used to find and download an organism proteome. See below where to find the ID.

`--isoform`

Allows isoform sequences to be added to the download.

`--prefix`

Decoy prefix to be added, default is _rev__.

`--reviewed`

Download only reviewed sequences from Swiss-Prot.


## Examples

Download a complete human proteome snapshot without isoforms, using Trypsin for protein digestion and adding contaminants.

`philosopher database --id UP000005640 --contam`

Download the reviewed version of the human proteome, containing isoforms and contaminants.

`philospher database --id UP000005640 --reviewed --contam`

Prepare a custom protein FASTA file for the analysis (skip the download).

`philosopher database --custom protein.fas --contam`

Download the complete human proteome and add external sequences from another FASTA file.

`philosopher database --id UP000005640 -add spikes.fas`

This example will download all reference sequences from the Human proteome, contaminants will be added to the database with the rev_ prefix and the resulting file will be digested using Lys C.

`philosopher database --contam --enzyme lys_c --reviewed --id UP000005640`


## FAQ

_Where can I find the UniProt Proteome ID for my organism?_

The list of all existing UniProt Identifiers can be found [here](http://www.uniprot.org/proteomes/). Before using Philosopher, you need to search the UniProt website for the correct id.
