Lightning Concepts
---

Lightning is a system that allows efficient storage of and access to
large scale population genomic data with a focus
on clinical and research use.
Lightning's primary focus is human genomic
data, but this could be extended to encompass other organisms.
Genomes are stored in a compressed format that compromises
between size and accessibility.

The Lightning system is a combination of a conceptual way
to think about genomes (genomic tiling),
the internal representation of genomes for efficient access
(the compact genome format and auxiliary data), and
software that manages access and analysis of that data.

Tiling is an efficient representation for genomic data that 
enables fast queries and machine learning.  It abstracts a called 
genome by partitioning it into overlapping variable length shorter
sequences, known as tiles.  A tile is a genomic sequence that is braced 
on either side by 24 base (24-mer) “tags.” Our choice of tags (“tag-set”)
partition the human reference genome into 10,655,006 tiles (the median tile 
is 250 bases long and the average is 314.5 bases per tile).   

The unique sequences for a tile position are called tile variants. The set 
of all positions and all tile variants is called the tile library. An individual 
called genome can then be easily represented as an integer array referencing the 
tile library. Each position in the array corresponds to a tile position and references 
the tile variant observed at that tile position for that individual. 

The major benefits of tiling are:  1) A set of genomes can be represented as a
numerical matrix so that we can easily use them with “out of the box” machine learning
(ML) and large data methods.  2) Tiled data represents the full genome including homozygous
reference calls and both phases. Therefore, we know if regions are confidently called as
reference or have variants.  3) Tiling is reference and sequencing technology independent.
This makes it possible to harmonize different studies that use a different reference
(e.g. GrCh37 vs GrCh38), different sequencing technologies, and/or to integrate 
genome, exome and microarray data.     4) Tiled data is compact and scalable. 
The human reference genome becomes ~10M tiles vs 3B bases. Tile data is stored in 
compact genome formatted (CGF) files which average 30-50 MB per genome.




