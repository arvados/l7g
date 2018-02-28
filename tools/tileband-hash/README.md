`tileband-hash`
===

Tileband hash is a tool to validate compact genome format (`CGF`) files
expressed in "band format" by hashing the implied sequence.
`tileband-hash` loads in a tile library into memory at startup, expressed
in simple genome library format (`SGLF`), and then loads all band information.

The output is two columns of hashes, one line per input band with
the first column for allele 1 and the second column for allele 2.

The hash of the implied sequence can then be compared with the implied
sequence of the source file to make sure the `CGF` conversion is valid.



Quick Start
---

```
$ cat test-data/hu826751-035e.band
[ 81 6 0 0 0 0 0 -1 0 0 0 411 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 1 0 0 -1 35 -1 194 1]
[ 81 2 0 0 0 0 0 -1 0 0 0 412 0 0 0 0 0 1 0 0 0 0 0 0 29 0 0 1 0 0 -1 35 -1 194 1]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
$ cat test-data/hu34D5B9-035e.band
[ 28 17 -1 0 0 0 0 -1 0 0 0 186 16 -1 6 0 0 0 0 0 5 0 0 0 7 0 0 0 1 0 -1 0 2 58 13]
[ 28 17 -1 0 0 0 0 -1 0 0 0 186 16 -1 6 0 0 0 0 0 5 0 0 0 7 0 0 0 1 0 -1 0 2 58 13]
[[ 0 72 ][ ][ ][ ][ 346 1 ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ 763 1 1003 1 2593 1 2968 1 3262 1 4485 1 4807 1 5094 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
[[ 0 72 ][ ][ ][ ][ 346 1 ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ 763 1 1003 1 2593 1 2968 1 3262 1 4485 1 4807 1 5094 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
$ ./tileband-hash \
  -L test-data/035e.1.sglf \
  -T 862 \
  <( cat test-data/hu826751-035e.band test-data/hu34D5B9-035e.band )
5bdfef6f0d40023ce59170477d9a5903 a50f886b344227e123745333f7d88885
01b1662617d9ee8fc1ba7749a8dea27f 01b1662617d9ee8fc1ba7749a8dea27f
```

Usage
---

```
tileband-hash version: 0.1.1
usage:
  tileband-hash [-n N] [-L sglf_stream] [-T tilepaths] [-v] [-V] [-h] bands

  -L sglf_stream  SGLF stream
  -T tilepaths    decimal list of tilepaths (e.g. '752+2', '752-753', '752,753')
  [-n N]          number of datasets to convert (input band count must be a 4x this number, default 1)
  [-W]            store sequence as ASCII in mem (default 2bit)
  [-v]            Verbose
  [-V]            Version
  [-h]            Help
```

Notes
---

The hash function is `MD5`.
The sequence that is hashed is an ASCII sequence with the letter `n` filled in for nocall sequence,
as expressed by the last two arrays in the band format.
