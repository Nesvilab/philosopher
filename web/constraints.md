### Are there any constraints ?

* One of the major source of problems using custom pipelines like Philosopher originates from badly formatted FASTA databases. Philosopher works with NCBI and UniProt databases or a database with headers formatted using the NCBI or [UniProt specification](http://www.uniprot.org/help/fasta-headers). If your FASTA file contains headers with different formatting we strongly suggest you to format them following the guidelines from one of the suggested on line Databases mentioned above.

* Philosopher currently only supports mzML files, if you have other format you will need to convert them to mzML (only important if you opt for doing quantification).
