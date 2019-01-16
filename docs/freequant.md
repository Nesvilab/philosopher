Quantify and compare the protein nature and spectral counting values between 2 or more Philosopher results.


## Usage

`philosopher freequant [files]`


## Flags

`--dir`

Folder path containing the raw files

`--ptw`

Specify the time windows for the peak (minute) (default 0.2).

`--tol`

M/Z tolerance in PPM (default 10).


## Example

Use spectra information from all mz files and, based on the peptide identifications, extract the given intensities. The resulting intensities will be stored inside the workspace, there is no output from this command.

`philosopher freequant --dir mz/`
