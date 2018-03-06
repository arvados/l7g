l7g
===

`l7g` is the main codebase for the Lightning system being developed
by Curoverse Research.

The repository contains documents, source code and pipelines for
the various aspects of Lightning.

Code here should be considered "research grade" and is a work in progress.

Overview
---

Lightning is a system based on "genomic tiling".
Genomes are split into small segments, on average roughly 250 base pairs (bp) long,
and these small segments are called "tiles".

For a given population of genomic data, the genomic sequences are tiled with
tiles that have redundant sequences de-duplicated.
Coalescing all unique tiles creates a "Lightning tile library", where a source sequence
from the population pool can be stored by using position references into the
lightning tile library.

A compact representation of a genome can be created by storing arrays of
indexes into the Lightning tile libary referencing their underlying sequence.

A representation of the compact genome representation we've developed is called
"compact genome format" (CGF) that can represent a whole genome in ~30Mb, depending
on the amount of low quality data in the original genome sample.


Directory Structure
---

### cwl-version/

[Common Workflow Language (CWL)](https://github.com/common-workflow-language/common-workflow-language) pipelines
for creating Lightning data.

### doc/

Lightning documentation

### go/

`go` (golang) programs used by Lightning.

### img/

Image directory for pictures.

### prototype/

A directory for the Lightning system prototype.

### proxy/

Authentication for Lightning prototype.

### sandbox/

Subdirectory for experimental code.

### tools/

Source and tools used by Lightning.
