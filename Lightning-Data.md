Lightning Data
===

Overview
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

This document will be focusing on the data files and formats
used in the Lightning system.
Please refer to [Lightning Concepts](Lightning-Concepts.md)
for the description of the tiling system and
[Lightning Architecture](Lightning-Architecture.md) for references to the
engineering architecture.

Data Overview
---

Here is a list of the main data components:

* Compact genome format files (CGF), one for each data genomic set
* Tile library, holding the underlying sequence
* Tagset library, holding the tile sequences
* Assembly data, holding the mapping of reference genomes to tile boundaries
* Tile span information

There are some auxiliary data files and formats that are used for
conversion but aren't used for Lightning server queries:

* FastJ
* Simple genome library format

In addition, for the Lightning prototype, there are the ClinVar
database and the untap database provided:

* ClinVar database, holding ClinVar variants and their tile mappings
* untap database, holding phenotype information for individuals from
  the Harvard Personal Genome Project and their mapping to the CGF datasets
  above

Compact Genome Format Files
---

See the [CGF schema](CGF-Schema.md) for a description of the CGF files.
The library and command line tools provided in [cgf](https://github.com/abeconnelly/cgf)
allow for inspection and creation of CGF files.

The current 651 CGF files (433 Thousand Genomes Project, 217 Harvard Personal Genome Project and
hg19) average at about 27MiB each.  Most of that information is fine grained no-call information
needed to reconstruct the original sequence.  The current CGF files can be found at:

* [Compact Genome Format files](https://workbench.su92l.arvadosapi.com/collections/su92l-4zz18-pi1dfpid43okayp)


Tile Library Format
---

The sequence information is stored on a per tile path basis.
All tile variants at a particular tile position are converted
into their twobit format and packed into a tar archive grouped
by tile path.  The tar archive is indexed using [tarindexer](https://github.com/devsnd/tarindexer)
and the tar archive is then `bgzip`'d, creating a `bgzip`
index file in the process.

For example, here are some of the first files in the resulting
tile library format directory:

```
0000.tar.gz      007b.tar.gz.gzi  00f6.tar.tai     0172.tar.gz      01ed.tar.gz.gzi  0268.tar.tai     02e4.tar.gz
0000.tar.gz.gzi  007b.tar.tai     00f7.tar.gz      0172.tar.gz.gzi  01ed.tar.tai     0269.tar.gz      02e4.tar.gz.gzi
0000.tar.tai     007c.tar.gz      00f7.tar.gz.gzi  0172.tar.tai     01ee.tar.gz      0269.tar.gz.gzi  02e4.tar.tai
0001.tar.gz      007c.tar.gz.gzi  00f7.tar.tai     0173.tar.gz      01ee.tar.gz.gzi  0269.tar.tai     02e5.tar.gz
0001.tar.gz.gzi  007c.tar.tai     00f8.tar.gz      0173.tar.gz.gzi  01ee.tar.tai     026a.tar.gz      02e5.tar.gz.gzi
0001.tar.tai     007d.tar.gz      00f8.tar.gz.gzi  0173.tar.tai     01ef.tar.gz      026a.tar.gz.gzi  02e5.tar.tai
0002.tar.gz      007d.tar.gz.gzi  00f8.tar.tai     0174.tar.gz      01ef.tar.gz.gzi  026a.tar.tai     02e6.tar.gz
...
```

The combination of twobit representation and `bgzip` allows for a relatively size compact tile library.

The combination of the tar archive and bgzip archive allow for fast random access to individual tile positions.

Here is an example of getting all tile variants for tile position `01ef.00.1d2e`:

```
$ grep 01ef.00.1d2e 01ef.tar.tai
01ef.00.1d2e.2bit 25138176 230
$ twoBitGulp -i <( bgzip -b 25138176 -s 230 01ef.tar.gz )
>01ef.00.1d2e.000
attagaggtagtgctctacaattttctgttaagttaagatggttgataat
gtcaatatcaaatcttctgtatccttactcatttttaagttggtctatca
attattgaagaaattgtgttaaatcctctaactttcactgtaaatttatt
tctccctatatttctctcagtttttggcaaatggattttgacacattgtt
gttatgattgttacaccttcctgatgcattggctttgttgtcattgtaa

>01ef.00.1d2e.001
attagaggtagtgctctacaattttctgttaagatggttgataatgtcaa
tatcaaatcttctgtatccttactcatttttaagttggtctatcaattat
tgaagaaattgtgttaaatcctctaactttcactgtaaatttatttctcc
ctatatttctctcagtttttggcaaatggattttgacacattgttgttat
gattgttacaccttcctgatgcattggctttgttgtcattgtaa
```

For the 651 samples, 433 thousand genome datasets, 217 Harvard Personal Genome Project datasets and
the hg19 reference, the tile library totals 9GiB.
The latest tile library can be found at:

* [Genome Library Format](https://workbench.su92l.arvadosapi.com/collections/su92l-4zz18-qh9cb6qf94h27gc)

Tagset
---

Lightning tile tag sequences are stored in a 2bit files.  Each
sequence name is identified by it's four hex path string followed
by a period and the 2 digit hex string tile library version
(e..g `01ef.00`).

The tags for a path are concatenated into a long sequence.  The
tag length, currently 24, along with the tile position
indicates how much of the sequence should be skipped to find the
relevant tag.  The first 24 bases in a tile path represent the ending tag of
the first tile.  The last 24 bases in a tile path represent the starting
tag sequence of the last tile in the tile path.
A tile path can be empty and have no tags in them.

For example, here is an illustrative example of how to find the beginning
and end tag for tile position `01ef.00.1d2e`:

```
$ twoBitGulp -i tagset.2bit -name 01ef.00 -w 24 | egrep -v '^>' | head -n 7470 | tail -n 2
atctatggatatttggatttctaa
attagaggtagtgctctacaattt
```

There are roughly 10M tile positions.  The current `tagset.2bit` is ~61MiB in size.
The latest tagset can be found at:

* [tagset.2bit](https://workbench.su92l.arvadosapi.com/collections/su92l-4zz18-v52bn3rrnk5vx4q)

Assembly
---

The mapping between a reference genome and which tile position they map to are kept
in a compressed and indexed text file,
The text file consists of a header followed by fixed with data.  The header is
a greater than symbol (`>`) followed by a triple of the reference name, chromosome
name and 4 digit ASCII hex tile path, delimited by a colon (`:`).
The subsequent data is a tilestep
followed by a non-inclusive, 0 index, ending base position of the tile.
Each tilestep and ending position are on individual lines with a tab (`\t`) immediately
following the 4 (ASCII) hex digit tile step and spaces left padding the (base 10)
ending tile position, padding to make a total of 10 characters.

For example, here is a portion of the `assembly.00.hg19.fw.gz` file:

```
...
202c      34599696
202d      34600000
>hg19:chr1:000c
0000      34600237
0001      34600462
0002      34600715
0003      34600949
...
5240      40099773
5241      40100000
>hg19:chr1:000d
0000      40100225
0001      40100457
0002      40100682
0003      40100907
0004      40101132
...
```

To calculate the beginning of a tile path, the chromosome and reference
position needs to be stored from the previous tile path.
If there is no previous tile path or the tile path has a different chromosome,
the beginning reference position is 0.
If the previous tile path has the same chromosome, then one past the reported value
of the last tile step is the beginning of the current tile path.

In the above example, tile path `0x000d` starts at `chr1`, reference position `40100001`.

Tile steps can be skipped, in which case the tile that the end tag falls in
is reported.
Tile step 0 is always the first tile step in a tile path.
To calculate the span of a tile, the previous tile step needs to be stored.
If the tile step increment is greater than one, the tile step one past the previous
tile step is the anchor tile and the difference between the reported tile step and previously
reported tile step is the span.

For example:

```
...
00b3      59032812
00b4      59373566
>human_g1k_v37:MT:035e
0000           244
0001           467
0002           959
0003          1394
0004          1767
0005          2225
0008          3338
0009          3563
000a          3788
000b          9760
000c          9985
...
```

For tile path `0x35e`, tile step `0x0008` is reported with end reference position `3338` with the previously
reported tile step of `0x0005` and reference position `2225`.
The next tile after tile step `0x0005` is `0x0006` and has span of `3` (`8-5=3`).
In this case, the anchor tile is tile step `0x0006` with span of `3` with the next tile after `0x0006` is `0x0009`.
The tile `0x0006` starts at reference position `201` (`225-24`) and ends (inclusively, 0 reference) on `3338`.


In addition to the fixed width file, a fixed width index file, reminiscent of
a FASTA index file, is used as an index file.  The format is a tab
delimited sequence name, size (in bytes), start (excluding header, in bytes),
line length (in bytes, excluding ending newline) and total line length (including
newline, in bytes).

For example:

```
$ head  assembly.00.hg19.fw.fwi
hg19:chr1:0000  86928   16      15      16
hg19:chr1:0001  185360  86960   15      16
hg19:chr1:0002  113792  272336  15      16
hg19:chr1:0003  120800  386144  15      16
hg19:chr1:0004  209504  506960  15      16
hg19:chr1:0005  160976  716480  15      16
hg19:chr1:0006  241776  877472  15      16
hg19:chr1:0007  211392  1119264 15      16
hg19:chr1:0008  237408  1330672 15      16
hg19:chr1:0009  117776  1568096 15      16
```

The fixed width file is compressed with `bgzip` and indexed, for a total of 3 files per reference tile assembly build.
In theory multiple tile reference assemblies could be provided in the same file but the current version has
each reference tile assembly in distinct files, with three files for `hg19` and three files for `human_g1k_v37`.

Note that the chromosome naming convention is used from the source reference sequence.
For example, `chrM` for `hg19` and `MT` for the `human_g1k_v37` reference.

The current fixed width files total ~56MiB.  The most current assembly files can
be found at:

* [Lightning hg19 Tile Assembly](https://workbench.su92l.arvadosapi.com/collections/su92l-4zz18-rg323w0m5a5ci7n)

Span
---

When variants from an input sequence fall on tags, the resulting tile is extended until the sequence
matches an unaltered tile tag.  The information for how many tags these variant skip is stored in
the span data file.

The span file consists of ASCII CSV file (`,` delimited), with a `4`, `2`, `4`, `3` hex
tile ID, followed by the span number (in decimal).  Only tile span information for tiles
that are other than span of 1 are held.

For example:

```
$ zcat span.gz | head
0000.00.0004.42b,2
0000.00.0004.42c,2
0000.00.0004.42d,2
0000.00.0004.42e,2
0000.00.0005.0da,2
0000.00.0005.0db,2
0000.00.0005.0dc,2
0000.00.0005.0dd,2
0000.00.0005.0de,2
0000.00.0005.0df,2
```

The current compressed `span.gz` file is ~49MiB.  The most current version can be found at:

* [Lightning Tile Span](https://workbench.su92l.arvadosapi.com/collections/su92l-4zz18-sxl2lizzsv4ewq1)


Simple Genome Library Format
---

The simple genome library format (SGLF) is used as a verbose, intermediate format for help
in conversion to the more size compact cousins.

The SGLF files are compressed CSV files, one per tile path.  The columns are the `4`, `2`, `4`, `3`
formatted tile id followed by a `+` and span information (in decimal), followed by the md5sum
of the tile sequence, followed by the sequence itself.

As an example:

```
$ zcat 01ef.sglf.gz  | head
01ef.00.0000.000+1,57f5909d01b7f8edd3f8652e7e610709,ggtgaatgttggctgtggagaatgaatccgaatcacttaggtcaaaagatgactaattccaaacacttttgctctatgcctgtttttatggtggccactctttgctctcaaacagggctcagaagaagagtgccaacaagtttctccacagaggggcactggctggcatccctgtaatacgcggtttgtagagaatgaaagcagctttggttttcttttgtacga
01ef.00.0000.001+1,1fbe5a4c30ceb0054ff6bf995bb6f2f1,ggtgaatgttggctgtggagaatgaatccgaatcacttaggtcaaaagatggctaattccaaacacttttgctctatgcctgtttttatggtggccactctttgctctcaaacagggctcagaagaagagtgccaacaagtttctccacagaggggcactggctggcatccctgtaatacgcggtttgtagagaatgaaagcagctttggttttcttttgtacga
01ef.00.0000.002+1,944d1b3b958a2c8ff4105c816c8b8044,ggtgaatgttggctgtggagaatgaatccgaatcacttaggtcaaaagatgactaattccaaacacttttgctctatgcctgtttttatggtggccactctttgctctcaaacagggctcagaaggagagtgccaacaagtttctccacagaggggcactggctggcatccctgtaatacgcggtttgtagagaatgaaagcagctttggttttcttttgtacga
01ef.00.0000.003+1,9188190557002b3f28d847d707464100,ggtgaatgttggctgtggagaatgaatccgaatcacttaggtcaaaagatgactaactccaaacacttttgctctatgcctgtttttatggtggccactctttgctctcaaacagggctcagaagaagagtgccaacaagtttctccacagaggggcactggctggcatccctgtaatacgcggtttgtagagaatgaaagcagctttggttttcttttgtacga
01ef.00.0001.000+1,c4c6db1a6f52cca5aafc43616246e4c9,cagctttggttttcttttgtacgagtgcacccagttaccggcatgacactatggtttcctcgctcggctttgaaatatagtaaactcacaaaagctactgtcgagttcagaacaaacacagggggtatgtttagttatttgtttttgtagtaaaataaaatatttattatacgcagacacctactgtgaacagagggggaggccactgtgttttatcttgccggtgtcatgtttgacttccattaggga
01ef.00.0001.001+1,d9303363bff612ec0e65108e7b0f189a,cagctttggttttcttttgtacgagtgcacccagttaccagcatgacactatggtttcctcgctcggctttgaaatatagtaaactcacaaaagctactgtcgagttcagaacaaacacagggggtatgtttagttatttgtttttgtagtaaaataaaatatttattatacgcagacacctactgtgaacagagggggaggccactgtgttttatcttgccggtgtcatgtttgacttccattaggga
01ef.00.0001.002+1,0a2d0c1ba8227f4daf56e303955d0f3c,cagctttggttttcttttgtacgagtgcacccagttaccggcatgacactatggtttcctcgctcggctttgaaatatagtaaactcacaaaagctactgtcgagttcagaacaaacacagggggtatgtttagttatttgtttttgtagtaaaataaaatatttattatacacagacacctactgtgaacagagggggaggccactgtgttttatcttgccggtgtcatgtttgacttccattaggga
01ef.00.0001.003+1,0b8c1f117c4467f3734f5d15c716641d,cagctttggttttcttttgtacgagtgcacccagttaccggcatgacactatggtttcctcgctcggctttgaaatatagtaaactcacaaaagctactgtcgagttcagaacagacacagggggtatgtttagttatttgtttttgtagtaaaataaaatatttattatacgcagacacctactgtgaacagagggggaggccactgtgttttatcttgccggtgtcatgtttgacttccattaggga
01ef.00.0001.004+1,75daad328a534dd860f91e219ce158c6,cagctttggttttcttttgtacgagtgcacccagttaccggcatgacactatggtttcctcgctcggctttgaaatatagtaaactcacaaaagctactgtcgagttcagaacaaacacagggggtatgtttagttatttgtttttgtagtaaaataaaatatttattatacgcagacacctactgtgaacagagcgggaggccactgtgttttatcttgccggtgtcatgtttgacttccattaggga
01ef.00.0001.005+1,7719b0d2005abd6955e0afb1c5665d46,cagctttggttttcttttgtacgagtgcacccagttaccggcatgacactatggtttcctcgctcggctttgaaatatagtaaactcacaaaagctactgtcgagttcagaacaaacacagggggtatgtttagttatttgtttttgtagtaaaataaaatatttattatacgcagacacctactgtgaacagagggggaggccactgtgttttatcttgctggtgtcatgtttgacttccattaggga
```

The tile library itself does not hold nocall information.  In construction, a heuristic was
used to populate bases with what was thought to be the most common base found, defaulting to `a`
if the entire population was no-call.

The current size is ~24GiB.  The current SGLF can be found at:

* [Simple Genome Library Format)(https://workbench.su92l.arvadosapi.com/collections/su92l-4zz18-zh3mc5wv4478yhm)


FastJ
---

FastJ is reminiscent of the FASTA file except the header line is replaced by a JSON string.
See the [FastJ](FastJ-Schema.md) for details.

Since FastJ is so verbose, it's used as an intermediate format to construct the tile library and
compact genome format files.


* [FastJ Files](https://workbench.su92l.arvadosapi.com/projects/su92l-j7d0g-fmbjujfq6wy7j1i#Data_collections)

Untap
---

See the [untap](https://github.com/abeconnelly/untap) project for details.

The most current version of the `untap` database can be found at:

* [untap](https://workbench.su92l.arvadosapi.com/collections/su92l-4zz18-ziluyxgz77rkekm)
