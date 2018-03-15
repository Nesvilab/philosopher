## Compatibility

Philosopher is partially compatible with 16 different operating systems and architectures:

* darwin/386
* darwin/amd64
* linux/386
* linux/amd64
* linux/arm
* freebsd/386
* freebsd/amd64
* openbsd/386
* openbsd/amd64
* windows/386
* windows/amd64
* freebsd/arm
* netbsd/386
* netbsd/amd64
* netbsd/arm
* plan9/386


## System requirements

There is no major restrictions on the system configuration in order to run Philosopher, every analysis will need writing permissions for the workspace and a few GB of disk space for the processed data. Philosopher has a base memory usage of approximately 500MB, the total amount of required memory will depend on how much data you have to analyze. For most medium-size high accuracy experiments, 4GB to 6GB RAM will be enough.

## Why partially ?

Third-party software like the _prophets_ will only work on Linux (Debian and Red Hat based) and Windows systems. The remaining functions will work on the above list. Also, the ProteoWizard tools _msconvert_ and _idconvert_ are only offer to GNU/Linux ad Unix operating systems.
