Automatic execution of consecutive analysis steps


## Usage

`philosopher pipeline [flags] [folders]`


## Flags

`--config`

Configuration file for the pipeline execution.

`--print`

Print the pipeline configuration file.


## Example


`philosopher pipeline --config philosopher.yaml folder1/ folder2/ folder3`


## FAQ

_How do I set the steps I want to run ?_

Your configuration file (philosopher.yaml) contains the steps you want to execute and
their respective parameters. You need to change the steps from _no_ to _yes_ in order to include them
into your analysis. After that, check each section parameters and set them accordingly
to your analysis.
