Lightning Tile
===

Introduction
---

In order to facilitate ease of use we've developed a "Lightning Tile Array" format
to access tile encoded genomes.
This document describes the format.

Overview
---

There are two formats under consideration for encoding Lightning Tile Arrays:
Simple one hot encoded Lightning tile arrays and span encoded one hot
Lightning tile arrays.
Auxiliary arrays are also provided to facilitate recreating tile path and positions
for each encoded array position.

Simple One-Hot Encoded Lightning Tile Arrays
---

For each tile position, the tile variant is converted to a one-hot representation
and added to the array.  Different alleles are interleaved so that two different
alleles for the same tile position have their one-hot representation placed
next to each other.  Each position is a float value.

No-call tiles are represented with a `NaN` for each element in the one-hot representation
for that tile position.
Non anchor spanning tiles are represented with a `0` for each element in
the one-hot representation for that tile position.

For example, if a tiled genome has the following tile variants, with `-1` representing
a non-anchor spanning tile and a `-2` representing a no-call tile:

```
[ [ ..., 0, 1, 0, 2, -1, 0, -2, 0, ...], [ ..., 0, 0, 0, 2, -1, 0, -2, 0, ... ] ]
```

This could be converted into a one-hot representation as follows:

```
[
  ...,

  1.0,
  1.0,

  0.0, 1.0,
  1.0, 0.0,

  1.0,
  1.0,

  0.0, 0.0, 1.0,
  0.0, 0.0, 1.0,

  0.0,
  0.0,

  1.0,
  1.0,

  NaN,
  NaN,

  1.0,
  1.0,

  ...

]
```

With spaces and newlines added in for clarity.

In a sense, the one-hot encoding is flattened for each tile position and allele, alternating alleles, until
the whole one-hot array is populated.
The size of the flattened one-hot encoding is taken to be one more
than the maximum tile variant value at that tile position, with a minimum one-hot encoding size of 1, in
the special case all values are no-calls.

As another example, with two encoded genomes, with the same special values as above:

```
[
  [ [ ..., 0, 1, ... ], [ ..., 0, 0, ... ] ],
  [ [ ..., 3, 0, ... ], [ ..., 2, 0, ... ] ]
]
```

could be encoded as:

```
[
  [
    ...,

    1.0, 0.0, 0.0, 0.0,
    1.0, 0.0, 0.0, 0.0,

    0.0, 1.0,
    1.0, 0.0,

    ...
  ],
  [
    ...,

    0.0, 0.0, 0.0, 1.0,
    0.0, 0.0, 1.0, 0.0,

    1.0, 0.0,
    1.0, 0.0,

    ...
  ]
]
```


Simple One-Hot Informational Lightning Tile Array Data
---

An informational data array is also provided to map
each position in the simple one hot encoding to a
lightning tile path, position and allele.

For each position in the simple one-hot encoded array,
there is a corresponding position in the informational
array.
Each value in the one-hot encoded array has value
`(tilepath * 16^4) + (tilestep) + (allele/2)`.

For example, in the above example, if the 12 one-hot encoded
values represented tile path `0x35e`, starting at tile position `0x2f`
and ending at tile position `0x30`, the corresponding informational
array would look like:

```
[
   ...,

   909.0, 909.0, 909.0, 909.0,
   909.5, 909.5, 909.5, 909.5,

   910.0, 910.0,
   910.5, 910.5,

   ...
]
```

With the simple one-hot encoded tile vector and the
informational vector, the original tile variant id vector can
be recreated.


--


Span Encoded One-Hot Lightning Tile Arrays
---

The span encoded one-hot lightning tile arrays are similar to the above
simple one-hot encoding with the difference in how non-anchor spanning
tiles are represented.
Each non-anchor tile is given the one-hot encoded value for the anchor
tile, shifted by the `d * M`, where `d` is the distance from the anchor
tile and `M` is the maximum of all the anchor tile variant ids for spanning
tile in that position and all the variant ids for the tile position.

With this way of encoding non-anchor spanning tiles, non-anchor spanning tiles
will be included in any tile concordance count that occurs.
The previous simple one-hot encoding will only count a matching spanning tile
once, whereas the span encoded one-hot version will add the length of the spanning
run for matching spanning tiles.

For the non-anchor, non-spanning tiles, the encoding is similar to the above:
for each tile position, the tile variant is converted to a one-hot representation
and added to the array.  Different alleles are interleaved so that two different
alleles for the same tile position have their one-hot representation placed
next to each other.  Each position is a float value.

No-call tiles are represented with a `NaN` for each element in the one-hot representation
for that tile position.

For example, if a tiled genome has the following tile variants, with `-1` representing
a non-anchor spanning tile and a `-2` representing a no-call tile:

```
[ [ ..., 0, 1, 0, 2, -1, 0, -2, 0, ...], [ ..., 0, 0, 0, 2, -1, 0, -2, 0, ... ] ]
```

This could be converted into a one-hot representation as follows:

```
[
  ...,

  1.0,
  1.0,

  0.0, 1.0,
  1.0, 0.0,

  1.0,
  1.0,

  0.0, 0.0, 1.0,
  0.0, 0.0, 1.0,

  0.0, 0.0, 1.0,
  0.0, 0.0, 1.0,

  1.0,
  1.0,

  NaN,
  NaN,

  1.0,
  1.0,

  ...

]
```

With spaces and newlines added in for clarity.

As another example, with two encoded genomes, with the same special values as above:

```
[
  [ [ ..., 0, 1, 1, -1, 0, 0, 0, ... ], [ ..., 0, 0, 0, 3, 3, -1, 0, ... ] ],
  [ [ ..., 3, 0, 1, -1, 0, 0, 0, ... ], [ ..., 2, 0, 0, 2, 3, -1, 1, ... ] ]
]
```

could be encoded as:

```
[
  [
    ...,

    1.0, 0.0, 0.0, 0.0,
    1.0, 0.0, 0.0, 0.0,

    0.0, 1.0,
    1.0, 0.0,

    0.0, 1.0,
    0.0, 1.0,

    0.0, 0.0,    0.0, 1.0,
    0.0, 0.0,    0.0, 1.0,

    1.0, 0.0, 0.0, 0.0,
    1.0, 0.0, 0.0, 0.0,

    1.0, 0.0, 0.0, 0.0,    0.0, 0.0, 0.0, 0.0,
    1.0, 0.0, 0.0, 0.0,    0.0, 0.0, 0.0, 0.0,

    1.0, 0.0,
    1.0, 0.0,

    ...
  ],
  [
    ...,

    0.0, 0.0, 0.0, 1.0,
    0.0, 0.0, 1.0, 0.0,

    1.0, 0.0,
    1.0, 0.0,

    1.0, 0.0,
    1.0, 0.0,

    1.0, 0.0,    0.0, 0.0,
    1.0, 0.0,    0.0, 0.0,

    0.0, 0.0, 0.0, 1.0,
    0.0, 0.0, 0.0, 1.0,

    0.0, 0.0, 0.0, 0.0,    0.0, 0.0, 0.0, 1.0,
    0.0, 0.0, 0.0, 0.0,    0.0, 0.0, 0.0, 1.0,

    1.0, 0.0,
    0.0, 1.0,

    ...
  ]
]
```

Span Encoded One-Hot Informational Lightning Tile Array Data
---

An informational data array is also provided to map
each position in the span encoded one hot array to a
lightning tile path, position, allele and span position,
if applicable.

For each position in the span encoded one-hot array,
there is a corresponding position in the informational
array.
Each value in the one-hot encoded array has value
`(tilepath * 16^4) + (tilestep*2) + (allele) + (K/maxK)`.
Here, `K` is the "tranche" of one-hot encoding that represents
whether the one-hot encoding position represents a spanning
tile.
`maxK` is the maximum of the anchor spanning tile value
and all tile variant values at that tile position.

For example, in the above example, if the 52 one-hot encoded
values represented tile path `0x35e`, starting at tile position `0x2f`
and ending at tile position `0x35`, the corresponding informational
array would look like:

```
[
   ...,

   956.0, 956.0, 956.0, 956.0,
   957.0, 957.0, 957.0, 957.0,

   958.0, 958.0,
   959.0, 959.0,

   960.0, 960.0,
   961.0, 961.0,

   962.0, 962.0,    962.5, 962.5,
   963.0, 963.0,    963.5, 963.5,

   964.0, 964.0, 964.0, 964.0,
   965.0, 965.0, 965.0, 965.0,

   966.0, 966.0, 966.0, 966.0,    966.5, 966.5, 966.5, 966.5,
   967.0, 967.0, 967.0, 967.0,    967.5, 967.5, 967.5, 967.5,

   968.0, 968.0,
   969.0, 969.0,

   ...
]
```

With the span encoded one-hot tile vector and the span encoded
informational vector, the original tile variant id vector can
be recreated.

