Lightning Concepts
---

Lightning is a system that allows for efficient access to
large large scale population genomic data with a focus
on clinical and research use.
The primary use of Lightning is for populations human genomic
data but this could be extended to encompass other organisms.
Genomes are stored in a compressed format that is a compromise
between size and accessibility.
Additional data data sources, such as phenotype data from
the Harvard Personal Genome Project and variant data from
the ClinVar database, are added for practical use.

The Lightning system is a combination of a conceptual way
to think about genomes (genomic tiling),
the internal representation of genomes for efficient access
(the compact genome format and auxiliary data) and the
software that manages access to the data.

This document will be focusing on some of the concepts
motivating the rest of the architecture and data formats.
Please refer to [Lightning Architecture](Lightning-Architecture.md)
for a description on the Lightning architecture.
and
[Lightning Data](Lightning-Data.md) for references to the
data structures used.


