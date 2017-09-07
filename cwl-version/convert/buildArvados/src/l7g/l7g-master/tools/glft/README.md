Genome Library Format Tool (glft)
===

A tool to manipulate Genome Tile Libraries (GLF).

This tool should be considered experimental and a work in progress.
Please use with caution.

Quick Start
---

```bash
git clone https://github.com/curoverse/l7g
cd l7g/tools/glft
make
```

```
$ ./glft
glft version: 0.1.0
usage:
  [-L]        Tile Library Version
  [-P]        Tile Library Path Version. From an SGLF file, get path library version.
  [-H]        Use hash as reported in SGLF file (default recalculate hash by sequence)
  [-v]        Verbose
  [-V]        Version
  [-h]        Help
```

Example Usage
---

```bash
$ zcat $LIGHTNING_DATA_DIR/sglf/022a.sglf.gz |  ./glft -P
717a750b70fbd719b3435146afe8bc2a
$ zcat $LIGHTNING_DATA_DIR/sglf/022a.sglf.gz | ./glft -P -H
717a750b70fbd719b3435146afe8bc2a
```
