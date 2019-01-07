Protein quantification based on isotope labeling


## Usage

`philosopher labelquant [flags]`


## Flags

`--annot`

Annotation file with custom names for the TMT channels. See below how to format your annotation file.

`--dir`

Folder path containing the raw files.

`--plex`

Number of channels.

`--purity`

Precursor ion purity threshold (default 0.5).

`--tol`

M/Z tolerance in PPM (default "10").

`--uniqueonly`

Report quantification based on only unique peptides.


## Example

A 10-plex TMT analysis:

`philosopher labelquant  --plex 10 --dir mz/`


## FAQ

_How to format my annotation file ?_

The annotation file is a simple text (.txt) file named _annotation.txt_. Each line must contain only 2 words,
the first one is the TMT label and the second one the custom name you want to use. Do not add headers to the file.

_example_

126 control1
127N control2
127C control3
128N treatment_1
128C treatment_1
129N treatment_2
129C treatment_2
130N treatment_3
130C treatment_3
131N pool_1
