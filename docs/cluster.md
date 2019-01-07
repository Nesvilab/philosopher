Sequence-based protein reports using a clustering approach to assemble a high-level report.


## Usage

`philosopher cluster [flags]`


## Flags

`--id`

UniProt proteome ID for retrieving annotation features.

`--level`

cluster identity level (default 0.9)


## Example

Clustering all reported proteins with at least 85% identity level.

`philosopher cluster --level 0.85`

Clustering all reported proteins and retrieving annotation data from UniProt.

`philosopher cluster --id UP000005640`
