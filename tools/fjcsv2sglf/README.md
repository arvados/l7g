fjcsv2sglf
===

A program to create SGLF files from CSV FastJ files.
This program is only meant to be used on a tile path at a time.

Quick Start
---

```
find /path/to/fastj -type f -name 035e.fj.gz | \
  xargs -n1 zcat | \
  fjt -C -U | \
  ./fjcsv2sglf <( tagset 035e ) | \
  bgzip -c > 035e.sglf.gz
```

Where `fjt` takes in a FastJ stream and spits out a CSV representation of the FastJ tile IDs, MD5 sequence digest and FastJ tile sequences.

For example:

```
zcat /path/to/hu826751-GS03052-DNA_B01/035e.fj.gz | \
  fjt -C -U | \
  head
035e.00.0000.000+1,7346f663d221ed28c112df86eb5986ef,gatcacaggtctatcaccctattaaccactcacgggagctctccatgcatttggtattttcgtctggggggcgtgcacgcgatagcattgcgggacgctggagccggagcaccctatgtcgcagtatctgtctttgattcctgcctcattctattatttatcgcacctacgttcaatattacaggcgaacatacctactaaagtgtgttaattaattaatgcttgtaggacataataataacaa
035e.00.0000.001+1,7346f663d221ed28c112df86eb5986ef,gatcacaggtctatcaccctattaaccactcacgggagctctccatgcatttggtattttcgtctggggggcgtgcacgcgatagcattgcgggacgctggagccggagcaccctatgtcgcagtatctgtctttgattcctgcctcattctattatttatcgcacctacgttcaatattacaggcgaacatacctactaaagtgtgttaattaattaatgcttgtaggacataataataacaa
035e.00.0001.000+1,1cadbbf41d5898b9e37ecdfd1d751f4e,gcttgtaggacataataataacaattgaatgtctgcacagccgctttccacacagacatcataacaaaaaatttccaccaaacccccccccctctccccccgcttctggccacagcacttaaacacatctctgccaaaccccaaaaacaaagaaccctaacaccagcctaaccagatttcaaattttatctttaggcggtatgcacttttaacagtcaccccccaactaacacattattttcccctcccactc
035e.00.0001.001+1,d69eb38a5317a1f387c49c84a003df36,gcttgtaggacataataataacaattgaatgtctgcacagccgctttccacacagacatcataacaaaaaatttccaccaaaccccccccctctccccccgcttctggccacagcacttaaacacatctctgccaaaccccaaaaacaaagaaccctaacaccagcctaaccagatttcaaattttatctttaggcggtatgcacttttaacagtcaccccccaactaacacattattttcccctcccactc
035e.00.0002.000+1,1d46f1fff282a060f2e3b28592daa12e,acacattattttcccctcccactcccatactactaatctcatcaatacaacccccgcccatcctacccagcacacacacaccgctgctaaccccataccccgaaccaaccaaaccccaaagacaccccccacagtttatgtagcttacctcctcaaagcaatacactgaaaatgtttagacgggctcacatcaccccataaacaaataggtttggtcctagcctttctattagctcttagtaagattacacatgcaagcatccccgttccagtgagttcaccctctaaatcaccacgatcaaaagggacaagcatcaagcacgcagcaatgcagctcaaaacgcttagcctagccacacccccacgggaaacagcagtgattaacctttagcaataaacgaaagtttaactaagctatactaaccccagggttggtcaatttcgtgccagccaccgcggtcacacgattaacccaagtcaatagaagccggcgtaaagagtgttttagatcacccc
035e.00.0002.001+1,1d46f1fff282a060f2e3b28592daa12e,acacattattttcccctcccactcccatactactaatctcatcaatacaacccccgcccatcctacccagcacacacacaccgctgctaaccccataccccgaaccaaccaaaccccaaagacaccccccacagtttatgtagcttacctcctcaaagcaatacactgaaaatgtttagacgggctcacatcaccccataaacaaataggtttggtcctagcctttctattagctcttagtaagattacacatgcaagcatccccgttccagtgagttcaccctctaaatcaccacgatcaaaagggacaagcatcaagcacgcagcaatgcagctcaaaacgcttagcctagccacacccccacgggaaacagcagtgattaacctttagcaataaacgaaagtttaactaagctatactaaccccagggttggtcaatttcgtgccagccaccgcggtcacacgattaacccaagtcaatagaagccggcgtaaagagtgttttagatcacccc
...
```

Building
---

```
make
```

Description
---

Simple Genome Library Format (SGLF) encodes Lightning libraries as a CSV file with the extended tile ID, MD5 digest of the sequence
and the sequence itself.  For example:

```
035e.00.0000.000+1,334516bc19d4674fb2dda0e79ffc6bb5,gatcacaggtctatcaccctattaaccactcacgggagctctccatgcatttggtattttcgtctggggggtgtgcacgcgatagcattgcgagacgctggagccggagcaccctatgtcgcagtatctgtctttgattcctgcctcattctattatttatcgcacctacgttcaatattacaggcgaacatacctactaaagtgtgttaattaattaatgcttgtaggacataataataacaa
035e.00.0000.001+1,c04a48b618d2df264a4fbff6f1ff326b,gatcacaggtctatcaccctattaaccactcacgggagctctccatgcatttggtattttcgtctggggggtgtgcacgcgatagcattgcgagacgctggagccggagcaccctatgtcgcagtatctgtctttgattcctgcctcattccattatttatcgcacctacgttcaatattacaggcgaacatacctactaaagtgtgttaattaattaatgcttgtaggacataataataacaa
035e.00.0000.002+1,108ab0afb8170afbc36f9b3d1d84dfcb,gatcacaggtctatcaccctattaaccactcacgggagctctccatgcatttggtattttcgtctggggggtgtgcacgcgatagcattgcgagacgctggagccggagcaccctatgtcgcagtatctgtctttgattcctgccccattctattatttatcgcacctacgttcaatattacaggcgaacatacctactaaagtgtgttaattaattaatgcttgtaggacataataataacaa
035e.00.0000.003+1,1d9149301e9cea080a52231c3a398998,gatcacaggtctatcaccctattaaccactcacgggagctctccatgcatttggtattttcgtctggggggtgtgcacgcgatagcattgcgagacgctggagccggagcaccctatgtcgcagtatctgtctttgattcctgccccattccattatttatcgcacctacgttcaatattacaggcgaacatacctactaaagtgtgttaattaattaatgcttgtaggacataataataacaa
035e.00.0000.004+1,dce802492cb8903a22c56e773ef807e2,gatcacaggtctatcaccctattaaccactcacgggagctctccatgcatttggtattttcgtctggggggtgtgcacgcgatagcattgcgagacgctggagccggagcaccctatgtcgcagtatctgtctttgattcctgcctcattctattatttatcgcacctacgttcaatattacaggcgaacatacctaccaaagtgtgttaattaattaatgcttgtaggacataataataacaa
...
```

SGLF has a requirement that the sequence contains only `a`, `c`, `g` and `t`, with all 'no-calls' (`n` or `N`) being replaced with some sequence.

`fjcsv2sglf` gathers a collection of FastJ tiles and attempts to fill in sequence for the no-calls it encounters.
The heuristic to fill in uncalled bases is as follows:

* If the no-call lands on a tag, fill in with the tag sequence
* Otherwise, compare sequences on the same tile step and tile span with the same sequence length and fill in each base with
  the most frequently occuring base at that position, counting `n` bases as `a`.

Only unique sequences are considered when doing the base frequency heuristic.

Once all sequences have been filled in, do a final de-duplication and order by the original sequence frequency.

Comments
---

Running the `fjt` tool converts FastJ to the CSV representation.
By deafult, `fjt` sorts the output but by using the unsorted option, `-U`, this saves
memory and time by allowing `fjt` to pass through the FastJ directly.

The internal representation that `fjcsv2sglf` uses is a twobit representation to save space.
Even with this, the in-memory tile library can get extremely large and so the SGLF files might
need to be split up and merged at a later step.

The heuristic to fill in was chosen more for it's algorithmic benefits rather than it's biological implications.
The only requirement for the derived SGLF sequences is that the original FastJ sequences match some SGLF
sequences without the requirement they be unique.
The fill in heuristic might change in the future and shouldn't be relied on.

Though the frequency order is not necessary,
the order of the SGLF sequences has implications for compactness of representation in CGF.

The tagset sequences are necessary for `fjcsv2sglf` in order to be able to fill in the nocall sequences that
appear on tags.

This tool is meant to be as close as possible as a drop in replacement for `fastj2cgflib`.


License
---

Copyright Curoverse Research, licensed under AGPLv3
