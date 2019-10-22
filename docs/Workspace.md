## Workspace management

Philosopher uses a workspace, where all processed and parsed data is stored locally in a custom binary format for more efficient data handling.

### Usage

`philosopher workspace [flags]`

## Flags

`--analytics`

reports when a workspace is created for usage estimation (default true). See the [analytics page](Analytics.md) for more information.

`--init`

Initialize the workspace. Should be executed on the directory where the experimental files are.

`--backup`

Create a backup of the experiment meta data in .zip format.

`--clean`

Remove the workspace and all meta data. Experimental file are kept intact.
