CGF Schema Draft
=====

Introduction
---

This document is a work in progress.

Compact Genome Format (CGF) is a binary format that encodes a representation of a genome relative to a tile library.

This is version 3 of the CGF format.  The major difference between the previous
format is that it keeps nocall information as separate as possible.


Overview
---

Compact Genome Format (CGF) is a way to represent whole genome information efficiently.
Each CGF file is stored relative to a tiled genome library, with the tiles stored
in each CGF file pointing into the tile library.

The CGF is a compromise between compression and accessibility.

From a high level perspective, CGF consists mainly of:

  - Vector arrays of tile variant ids for each tile position (the 'pointers' into library).
  - Overflow tile variant pointers that can't fit into the fixed width records above.
  - 'NOCALL' information that explicitly holds information of where the no-calls fall within tiles.

Binary CGF Structure
---

```go
Magic             8byte
CGFVersion        String { int32, []char }
LibraryVersion    String { int32, []char }
PathCount         8byte
TileMap           String { int32, []char }
PathStructOffset  []8byte                   // absolute, from beg. of file
PathStruct        []{

  Name          String { int32, []char }
  NTileStep     8byte
  NOverflow     8byte
  NOverflow64   8byte
  ExtraDataSize 8byte

  Loq               []byte
  Span              []byte

  Cache             []8byte

  Overflow          []{
    TileStep      2byte
    TileVariant   [][ 2byte, 2byte ]
  }

  Overflow64        []{
    TileStep      []8byte
    TileVariant   [][ 8byte, 8byte ]
  }

  ExtraData         []byte

  LoqTileStepHomSize      8byte
  LoqTileVariantHomSize   8byte
  LoqTileNocSumHomSize    8byte
  LoqTileNocStartHomSize  8byte
  LoqTileNocLenHomSize    8byte

  LoqTileStepHetSize      8byte
  LoqTileVariantHetSize   8byte
  LoqTileNocSumHetSize    8byte
  LoqTileNocStartHetSize  8byte
  LoqTileNocLenHetSize    8byte

  LoqTileStepHom      sdsl::enc_vector
  LoqTileVariantHom   sdsl::vlc_vector
  LoqTileNocSumHom    sdsl::enc_vector      // running sum of entry count, inclusive
  LoqTileNocStartHom  sdsl::vlc_vector
  LoqTileNocLenHom    sdsl::vlc_vector

  LoqTileStepHet      sdsl::enc_vector
  LoqTileVariantHet   sdsl::vlc_vector
  LoqTileNocSumHet    sdsl::enc_vector      // running sum of entry count, inclusive
  LoqTileNocStartHet  sdsl::vlc_vector
  LoqTileNocLenHet    sdsl::vlc_vector

}
```

Notes
---

* The first 32 bits of each entry in Vector hold a bit to indicate whether
  it's 'canonical' or whether the overflow table should be consulted.
* Barring some complications, a set 'canonical' bit means the tile is canonical (the
  default tile variant for that tile position) and an unset 'canonical' bit
  indicates a non-canonical tile.
* A high quality non-canonical tile can be deduced from consulting the hexit cache
  or the overflow table if the tile variant isn't/can't be stored in the hexit cache.
* Each hexit has 3 values reserved:
  - 0xf - high quality overflow
  - 0x0 - complex
  - 0x1-0xe - lookup in the tilemap
* If there are more than 32/4 = 8 overflow entries, the `Overflow` vector should be  consulted.
* If entries are too big to fit in the `Overflow` vector, `Overflow64` should be consulted.
* To discover anchor tile, this is indicated with the a `Span` bit set and the `Cache` canonical bit not set.
* A non-anchor spanning tile is indicated with the `Span` bit set and the canonical bit set in the `Cache`.
* `Het` low quality fields have the `Variant`, `Sum`, `Start`, `Len` have two entries per 'field'.
* The low quality `LoqTileStep` and `LoqTileNocSum` benefit from delta encoding as they are strictly increasing
  numbers.  The other low quality fields benefit from the variable length encoding as they are more evenly distributed
  within a range.  The variable length encoding is better taken advantage of when these vectors are split out instead
  of interleaved together.
* A future iteration might want the `Loq` and `Span` vectors to be interleaved in the cache to take advantage of data
  locality.
* `ExtraData` is a catch-all to provide a facility to store other data or tiles that aren't in the library, say.
* To indicate a non-anchor spanning tile, the following reserved values are used:
  - `Overflow.TileVariant`: `1<<15 - 1` (e.g. `0xffff`)
  - `Overflow64.TileVariant`: `1<<64 - 1` (e.g. `0xffffffffffffffff`)
  - `LoqTileVariantHom`: `1<<32 - 1` (e.g. `0xffffffff`)
  - `LoqTileVariantHet`: `1<<32 - 1` (e.g. `0xffffffff`)
* Fixed width fields are put at the beginning for ease of parsing

Description
---

### Tile Map

The tile map is stored as a string with a tile map entry stored per line.
For each line, each allele is separated by a `:`.
For each allele, tiles are separated by a `;`.
If the span of a tile is greater than 1, the length of the span is indicated by a `+` followed by the number of "base" tiles it spans.

For example, here is a portion of a tile map:

```
0:0
0:1
1:0
1:1
0:2
2:0
0;0:1+2
1+2:0;0
0+2:0+2
0:3
3:0
1;0:0+2
0+2:1;0
0:4
4:0
1+2:1+2
...
```


### Vector Data

The bulk of the high quality information is stored in the `PathStruct` array.
`PathStruct.Cache` should be around 2.5 MiB (10M tiles, 64 bits per 32 tiles).
Assuming around 4% overflow from the `Cache` vector, the `Overflow` should be in the
range of ~2.5Mb (10M (tile positions) * 0.05 * 3 (entries/position) * 2 (bytes/entry) ~ 2.5Mb).
`Span` and `Loq` both are bit vectors, with a bit reserved for each tile position,
taking up around 1.25Mb each for around 2.5Mb together.
`Overflow64` is there to catch any tiles that can't be represented with the 16 bit `Overflow`
structure.
Until the population exceeds 2^16 the `Overflow64` structure should be 0.

The bulk of the size of any whole genome CGF file will probably be the low quality information.
The low quality information will depend heavily on the technology used to sequence and how
much loq quality data needs to be saved.
Initial tests indicate that this is in the 19Mb region for the CGI sequenced Harvard Personal
Genome whole genome data.


Though this might change in the future, a `Cache` element is chosen to be 64 bits with
the first 32 bits allocated for 'canonical' bits and the last 32 bits to allocated
for the hexit encoding.

A diagram is illustrative:

    /--------- c ----------\/----------------------------- H ------------------------------\
    |    canonical bits    |                       hexit region                            |
    ----------------------------------------------------------------------------------------
    |                      |                                                               |
    |                      | [ hexit_0 ] [ hexit_1 ] ... [ hexit_{k-1} ]  <    unused   >  |
    |                      |                                                               |
    ----------------------------------------------------------------------------------------
                             \___ h ___/ \___ h ___/     \_____ h _____/  \__ H - k*h __/
    \______________________________________ b _____________________________________________/

Here, `b=64`, `s=32`, `h=4` and `H=32`.

### Low Quality Information

The low quality structures hold information on the tile variant and the runs of nocalls within
the tile sequence.
Each of the low quality structures is an encoded array, with each of the encodings chosen to
best fit the type of data (that is, variable length encoding or delta encoding, depending).

The "het" and "hom" indicates whether the low quality information is the same for both
alleles and **doesn't** represent whether the tile variants are the same for both alleles.
In the below, all tile variant arrays have two values per entry to represent each allele.

* `LoqTileStepHom` : tile steps
* `LoqTileVariantHom` : tile variants, two per entry
* `LoqTileNocSumHom` : position in `StartHom` and `LenHom` arrays for this entry
* `LoqTileNocStartHom` : the start position in the tile of the nocall run
* `LoqTileNocLenHom` : the length of the nocall run

* `LoqTileStepHet` : tile steps
* `LoqTileVariantHet` : tile variants, two per entry
* `LoqTileNocSumHet` : position in `StartHom` and `LenHom` arrays for this entry, two per tile position entry
* `LoqTileNocStartHet` : the start position in the tile of the nocall run
* `LoqTileNocLenHet` : the length of the nocall run

Misc
---

```
Read 1 MB sequentially from memory    250,000 ns    250 us    0.25 ms
```

If the high quality data is around 7.5Mb, doing a whole genome concordance should
be in the range of ~4ms (.25ms/Mb * 7.5Mb * 2datasets).

References
---


  - [2bit encoding in closure (*BROKEN*)](http://eigenhombre.com/2013/07/06/a-two-bit-decoder/)
  - [SDSL-lite](https://github.com/simongog/sdsl-lite)
  - [SDSL cheat sheet](http://simongog.github.io/assets/data/sdsl-cheatsheet.pdf)
  - [Latency Numbers Every Programmer Should Know](https://gist.github.com/jboner/2841832)

