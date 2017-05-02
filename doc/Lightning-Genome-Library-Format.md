Lightning Genome Tile Library Format
===

The Lightning Genome Tile Library Format, Genome Library Format (GLF) for short,
is a suite of formats to store and represent the the genome library for Lightning.

To exploit the redundancy in the genome, the genome is 'tiled', breaking the genome
into relatively small contiguous sequences (~250bp).  A population is tiled with
redundant tiles taken out.  A tile library is formed by taking all unique sequences
at a particular tile position and ordering them to create a tile variant for each
tile position.

The compact representation of genomes only stores the variant for each tiled position.
In order to recreate the sequence, a lookup of the underlying sequence needs to occur.
The Lightning genome tile library format is used to do the conversion from a tiled
representation back to sequence.

There are some auxiliary files with the genome library format to make sure lookups against
a reference genome can be done as well as tagset data.

Overview
---

To facilitate different uses of the tile library, different formats are created.
The main formats are the SGLF, the 2bit and the 2bit tar formats.
Each of these three different formats are trying to focus on different aspects
of usability:

* SGLF is verbose but easier to process
* 2bit is smaller but harder for random access
* 2bit tar is smaller and can be efficiently randomly accessed but needs
  specialized knowledge and tools

Auxiliary data structures are used to get out spanning tile information,
tag information and reference position information.

### Terms

Some common terms will be used:

* `Tile ID`: The tile identification, often 'joined dot' notation, with the tile path first,
  the tile library version second, the tile step third and the tile variant fourth.  This
  is often represented as a hex text sequence with 4, 2, 4 and 3 digits respectively.
  For example: `0251.00.341e.009`
* `Extended Tile ID`: The same as the `Tile ID` but with span information as a hex digit
  with a `+` separator.  For example: `0251.00.341e.009+1`



Simple Genome Library Format
---

The Simple Genome Library Format (SGLF) is used as an auxiliary data structure to store
the sequence along with tile path, it's hash and the Tile ID.
Each tile path is split into it's own file and ordered by Tile ID.
Each field is separated by a comma (`,`), with the first field being the extended Tile ID,
the second being the hash of the text sequence and the third being the sequence itself.

Here is a truncated example:

```
...
012e.00.0029.000+1,9cfe1138ccc6b7b55df8ab7998bbb124,ctccatcaaaaatgagtttactctaaatatatggatttatatctagcttctctattctgtttcattaatctatgtgtctgtatttatgtcagtactatgatattttggttactatagctttgtagtgtcatttaaagtcaggtagcatgatacctccagctttcttctttttgttcagattttttttgactattctgagtcttttgtgattccatttaaatttaaggatttctttttctatttatatgaagactatcattggtattttaatagggattgcattgaatctgtagatttatttgagttgtgtggacatat
012e.00.0029.001+1,57ec1da0e03b46a6a32b7027c9cb4e26,ctccatcaaaaatgagtttactctaaatatatggatttatatctagcttctctattctgtttcattaatctatgtgtctgtatttatgtcagtactatgatattttggttactatagctttgtagtgtcatttaaagtcaggtagcatgatacctccagctttcttctttttgttcagattttttttgactattctgagtcttttgtgattccatttaaatttaagggtttctttttctatttatatgaagactatcattggtattttaatagggattgcattgaatctgtagatttatttgagttgtgtggacatat
012e.00.0029.002+1,441e91b43245f6c91d32d87080a8cad0,ctccatcaaaaatgagtttactctaaatatatggatttatatctagcttctctattctgtttcattaatctatgtgtctgtatttatgtcagtactatgatattttggttactatagctttgtagtgtcatttaaagtcaggtagcatgatacctccagctttcttctttttgttcagattttttttgactattctgagtcttttgtgattccatttaaatttaaggatttctttttctatttatatgaagactatcattggtattttaatagggattgcattgaatttgtagatttatttgagttgtgtggacatat
012e.00.0029.003+1,15a3c957a7e795b70aa6017830be7120,ctccatcaaaaatgagtttactctaaatatatggatttatatctagcttctctattctgtttcattaatctatgtgtctgtatttacgtcagtactatgatattttggttactatagctttgtagtgtcatttaaagtcaggtagcatgatacctccagctttcttctttttgttcagattttttttgactattctgagtcttttgtgattccatttaaatttaaggatttctttttctatttatatgaagactatcattggtattttaatagggattgcattgaatctgtagatttatttgagttgtgtggacatat
012e.00.002a.000+1,6acafba77ae5e6d5f08eb3d2733e910c,tttatttgagttgtgtggacatattaacaatataccaaaagtagagaaagactacaaaaaaaggaaactacagaccaatgtctctgatgactatagatgcaaaagttctcaaaaatactagtgaactgaattcaacaacacattaaaaaatgattcatcacaatcaagtgagattcattccagggatacaaagatggtttaacatacacaaattaataaatgtcatacatcacacaacaggatgaataa
012e.00.002a.001+1,30d92913faf83f050e65dcbccb51146d,tttatttgagttgtgtggacatattaacaatataccaaaagtagagaaagactacaaaaaaaggaaactacagaccaatgtctctgatgactatagatgcaaaagttctcaaaaatactagtgaactgaattcaacaacacattaaaaaatgattcatcacaatcaagtgagattcattccagagatacaaagatggtttaacatacacaaattaataaatgtcatacatcacacaacaggatgaataa
012e.00.002a.002+2,5720cfcff975b390061d02709cc122ca,tttatttgagttgtgtggacatattaacaatataccaaaagtagagaaagactacaaaaaaaggaaactacagaccaatgtctctgatgactatagatgcaaaagttctcaaaaatactagtgaactgaattcaacaacacattaaaaaatgattcatcacaatcaagtgagattcattccagggatacaaagatggtttaacatacacaaattaataaatgtcatacaccacacaacaggatgaataacaaaaactatatgatcatctcaatagatgcagagaaaacatttaatgaaattaaacattccgtaatgataaactctcaataagctggatataaaaggaacatagtccaacacaataaaggccatgtatgaaaaattcacagccagtatcatactaaacaagggaaaactagaagcatttcctctaagatcatgaaacaaacaagtatgcccacttttacaacttttattcaagataggactggaagtcctagcaagagcaatcaaaaaagggaaata
...
```

2bit Format
---

For ease of use, each tile path is compressed into it's own 2bit file for use with the suite of 2bit tools.
Each sequence appears with it's extended Tile ID as it's identifier.
Each 2bit file is compressed to save space, often giving 10x space savings probably due to the high redundancy of
storing so many similar sequences.

For example:

```
$ twoBitGulp --list-names --width 50 -i <( zcat 035e.2bit.gz ) | head -n 19
>035e.00.0000.059+1
gatcacaggtctatcaccctattaaccactcacgggagctctccatgcat
ttggtattttcgtctggggggtgtgcacgcgatagcattgcgagacgctg
gagccggagcaccctatgtcgcagtatctgtctttgatccctgcctcatt
ctattatttatcgcacctacgttcaatattacaggcgaacatacctatta
aagtgtgttaattaattaatgcttgtaggacataataataacaa

>035e.00.0000.0a1+2
gatcacaggtctatcaccctattaaccactcacgggagctctccatgcat
ttggtattttcgtttggggggtgtgcacgcgatagcattgcgagacgctg
gagccggagcaccctatgtcgcagtatctgtctttgattcctgccccatt
ctgttatttatcgcacctacgttcaatattacaggcgaacatatctacta
aagtgtgttaattaattaatgcttgtaggacatagtaataacaattgaat
gtctgcacagccgctttccacacagacatcataacaaaaaatttccacca
aacccccccctctccccccgcttctggccacagcacttaaacacatctct
gccaaaccccaaaaacaaagaaccctaacaccagcctaaccagatttcaa
attttatctttaggcggtatgcacttttaacagtcaccccccaactaaca
cattattttcccctcccactc

```

The process to choose how to fill in this data is done when assembling the 2bit library representation from the SGLF files.
The process for filling in sequence was chosen with a focus on space and access rather than any biological motivation.

The fill in process will be discussed later.

2bit, tar format
---

Whereas SGLF was concerned with ease of processing at the cost of space and the 2bit format was concerned with
compatibility and space, the 2bit, tar format is concerned with space and access.
The SGLF files are too big to access quickly and the 2bit files need to load information in on a tile path basis to extract
sequence for a single tile step.
In order to retain the compactness of the combination of 2bit representation and gzip compression, each tile step is compressed
into a 2bit file, bundled into a tar archive on a tile path basis and compressed.
Further, index files are produced to indicate where in the archive each 2bit tile step file is located for fast random access.

As an example, here is a small set of steps to extract the sequence for the extended tile path `035e.00.0018.01c`:

```bash
$ ls 035e*
035e.tar.gz  035e.tar.gz.gzi  035e.tar.tai
$ grep '035e.00.0018' 035e.tar.tai
035e.00.0018.2bit 899584 4635
$ twoBitGulp --name '035e.00.0018.01c' -i <( bgzip -c -b 899584 -s 4635 035e.tar.gz )
>035e.00.0018.01c
accccattaaacgcctggcagccggaagcctattcgcaggatttctcatt
actaacaacatttcccccgcatcccccttccaaacaacaatccccctcta
cctaaaactcacagccctcgctgtcactttcctaggacttctaacagccc
tagacctcaactacctaaccaacaaacttaaaataaactccccactatgc
acattttatttctccaacatactcggattctaccctagcatcacacacc

```

Span Information
---

The 2bit tar format does not include span information in it's sequence name.
To recover this information, a separate text file is used that has a comma separated file (`,`)
with two fields, the first being the tile ID and the second being the span information.

For example:

```bash
$ zgrep '^035e.00.0018.01c' span.gz
035e.00.0018.01c,1
```

The idea is that this is meant to be loaded in once at startup by whatever process needs it.


Assembly Information
---

The associated assembly information files hold how to map tiles back to a reference build.
Currently only `hg19` is supported.


The main file holds the inclusive end locations in base 0 reference position
of each tile position in the tile library to the reference in question.
The format is two fixed width fields separated by a tab, each
field padded with spaces (to the left).  Each block contains a line starting with a
`>` that holds the name of the assembly (e.g.  `hg19`), the chromosome name (`chr`
prefix, e.g. `chr13`) and the path in hex (e.g. `2c5`) all separated by a colon (`:`) delimiter.

For example:

```
...
534c      48101484
534d      48129895
>hg19:chr22:031b
0000       3800000
>hg19:chr22:031c
0000       8300000
>hg19:chr22:031d
0000      12200000
>hg19:chr22:031e
0000      14700000
>hg19:chr22:031f
0000      16054701
0001      16056864
0002      16070005
0003      16083942
...
```

This file is compressed and index with `bgzip`.

An index into the fixed width (uncompressed) file is also created.
The format is reminiscent of the FASTA index format used by htslib ([faidx](http://www.htslib.org/doc/faidx.html)).

From the faidx htslib documentation site:

```
NAME Name of this reference sequence
LENGTH  Total length of this reference sequence, in bases
OFFSET  Offset within the FASTA file of this sequence's first base
LINEBASES The number of bases on each line
LINEWIDTH The number of bytes in each line, including the newline
```

Where `LINEBASES` is replaced with `\n` characters on the line.

For example:

```
...
hg19:chr21:031a 341216  159087461       15      16
hg19:chr22:031b 16      159428694       15      16
hg19:chr22:031c 16      159428727       15      16
hg19:chr22:031d 16      159428760       15      16
hg19:chr22:031e 16      159428793       15      16
hg19:chr22:031f 43920   159428826       15      16
...
```

In the example above, the index line for `hg19:chr22:31f` is saying that it starts
at (uncompressed) byte offset 159428026, has length 43920, with 15 non-return
(`\n`) per line and 16 total characters per line (including `\n`).

For example, to get all information for `hg:chr21:31a`, the following could be done:

```bash
$ begin=`egrep '^hg19:chr21:031a' assembly.00.hg19.fw.fwi | cut -f3`
$ bsize=`egrep '^hg19:chr21:031a' assembly.00.hg19.fw.fwi | cut -f2`
$ bgzip -c -b $begin -s $bsize assembly.00.hg19.fw.gz | head
0000      42600225
0001      42600512
0002      42600739
0003      42600964
0004      42601189
0005      42601414
0006      42601639
0007      42601864
0008      42602089
0009      42602314
```

Tile Library Versioning
---

The library version scheme uses a content versioning mechanism to associate a hash to the tile set used.
Since there are different ways to encode a tile library, a versioning scheme that is encoding agnostic
is desirable.

First, a tile path text manifest is created by making an extended tile ID and sequence hash pair, using a
colon as delimiter (`:`) with each pair being separated by a space.
Each tile path manifest has no trailing newline character.
The tile path manifest is hashed to get a hash per tile path.

The tile library manifest is then created in an analogous way to the tile path manifest by listing
all tile path and tile path hash key pairs, delimited by a colon (`:`), with a space (` `) in between
each tile path, hash key pair.
The hash of the manifest is then taken as the tile library version.

Both are ordered in ascending order of tile step or tile path, respectively.

For example, here is what the beginning of a tile path manifest might look like for tile path `0x35e`:

```
035e.00.0000.000+1:334516bc19d4674fb2dda0e79ffc6bb5 035e.00.0000.001+1:c04a48b618d2df264a4fbff6f1ff326b 035e.00.0000.002+1:108ab0afb8170afbc36f9b3d1d84dfcb 035e.00.0000.003+1:1d9149301e9cea080a52231c3a398998 035e.00.0000.004+1:dce802492cb8903a22c56e773ef807e2 ...
```

Here is what the beginning of the tile library manifest might look like:

```
0000:c903b729f10f71ff44abefa75b6a5533 0001:825523568814f8eac8b41a9113955afb 0002:46c27ad97dae0ce2c4759655411be6d1 0003:eb5030d5538c7f289820d58842502bfc 0004:03d06f41667a28de3dd2f7cacedc9673 ...
```

The hash is `MD5` but this might change in the future.

The library version is done on the 2bit or 2bit tar sequence described above and not on the sequences as they
appear in the SGLF.
No newlines should appear in the manifests.
No spaces should appear in the prefix or trail at the end of the manifest.
Only one space should be in between hash- key pair entries.

### Tiled Genome Version Comparison

Given a tiled genome, one would like to know if the tile library can support conversion.
That is, it would be desirable to know whether the tile library can support conversion to
raw sequence of the tiled genome.

To support this function, a hash of the genome can be done by creating a path manifest of the
interleaved tile IDs and their sequence hashes, then taking the hash of the path manifest.

When querying whether a tile library server supports conversion, the tile genome hash can be sent along
with a list of the tile IDs so that the tile library server can compare it's computed
genome hash with the one received.

For convenience, a small database of known genomes along with their hashes can be loaded on the
tile library server for quick lookup.


### Tile Library Difference

To help keep track of what tiles have been added and to try and keep some provenance of tile versioning,
the tile updates to the tile library, it's new hash and other information is provided in a JSON file.

The structure should be something like:

```javascript
{
  "TileLibrary": {
    "<tile library version hash>": {
      "from": {
        "<previous tile library version hash>": {
          "tiles": [
            { "tileId" : "<extended tile id>", "operation":"<add|delete|change:tileId>", "md5sum":"<sequence hash>" },
            { "tileId" : "<extended tile id>", "operation":"<add|delete|change:tileId>", "md5sum":"<sequence hash>" },
            ...
          ]
        }
      }
    }
  }
}
```

Note that creating the current tile library version can be done by going further back in time and adding the
appropriate tiles.
The previous tile library in the above is only a record of the method the current version
of the tile library was created.

A tile library version hash object can have multiple `from` library versions.

