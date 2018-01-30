FastJ Tool (fjt)
===

`fjt` is a tool to manipulate FastJ (text) files.
`fjt` is still experimental so please use with caution.

Quick Start
---

```
$ make
$ ./fjt -h
fjt version: 0.1.2
usage:
  fjt [-c variant] [-C] [-v] [-V] [-h] [input]

  [-C]            Output comma separated `extended tileID`, `hash` and `sequence` (CSV output)
  [-B]            Output band format
  [-b]            input band format and output FastJ (requires '-L sglf' option)
  [-H]            batch hash of input bands (requires '-L sglf' option)
  [-c variant]    Concatenate FastJ tiles into sequence.  `variant` is the variant ID to concatenate on
  [-L sglf]       Simple genome library format tile path file
  [-i ifn]        input file
  [-I ifn_list]   file containing gziped list of FastJ files to convert
  [-p tilepath]   Tile path (in decimal)
  [-U]            do not sort output (for use with -C option)
  [-v]            Verbose
  [-V]            Version
  [-h]            Help
```

Commands
---

### Concatenate `fjt -C`

Concatenate a series of FastJ tiles, filtered by Variant ID.

For example, given the four tiles saved in a file `example-tiles.fj`:

```
>{"tileID":"0157.00.0000.000","md5sum":"d7c151e05a4be1a735af3d1346c90714","tagmask_md5sum":"d7c151e05a4be1a735af3d1346c90714","locus":[{"build":"hg19 chr7 0 54000589"}],"n":589,"seedTileLength":1,"startTile":true,"endTile":false,"startSeq":"agaaaatgccaacatagccagagt","endSeq":"cctaagcataaaaaacaaatggag","startTag":"","endTag":"cctaagcataaaaaacaaatggag","nocallCount":0,"notes":[]}
agaaaatgccaacatagccagagtgataaattaattctatagaccaacaa
gtcaaacataagaaaaagttggaaaaatttaccataccaattacatttta
attgtagtagtatatctccgtcatctcagcccaaaaactccttaagctga
taagcaacttcagcaaggtctcagcatacaaaatcaatgtgcaaaaatca
caagcattccttcacaccaacaatagacaagcagagagccaaatcatgaa
tgaactcccatttacaatagctacaaagagaataaaatacctaagaatac
agttaacaagggatgtgaaggacctcttcaaggagaactacaaaccactg
ctcaaggaaataatagaggacacaaacaaatggaaaaacgttccatcctc
atggataggaagaacaaatatcgtgaaaatggccatactgcccaaagtaa
ttaatagattcattgctattcccatcaaactaccattgacattcttcaca
gaatcagaaaaaactactttaaatttcatgtagaatcaaagaagaccctg
tatagccaagacaatcctaagcataaaaaacaaatggag

>{"tileID":"0157.00.0000.001","md5sum":"d7c151e05a4be1a735af3d1346c90714","tagmask_md5sum":"d7c151e05a4be1a735af3d1346c90714","locus":[{"build":"hg19 chr7 0 54000589"}],"n":589,"seedTileLength":1,"startTile":true,"endTile":false,"startSeq":"agaaaatgccaacatagccagagt","endSeq":"cctaagcataaaaaacaaatggag","startTag":"","endTag":"cctaagcataaaaaacaaatggag","nocallCount":0,"notes":[ ]}
agaaaatgccaacatagccagagtgataaattaattctatagaccaacaa
gtcaaacataagaaaaagttggaaaaatttaccataccaattacatttta
attgtagtagtatatctccgtcatctcagcccaaaaactccttaagctga
taagcaacttcagcaaggtctcagcatacaaaatcaatgtgcaaaaatca
caagcattccttcacaccaacaatagacaagcagagagccaaatcatgaa
tgaactcccatttacaatagctacaaagagaataaaatacctaagaatac
agttaacaagggatgtgaaggacctcttcaaggagaactacaaaccactg
ctcaaggaaataatagaggacacaaacaaatggaaaaacgttccatcctc
atggataggaagaacaaatatcgtgaaaatggccatactgcccaaagtaa
ttaatagattcattgctattcccatcaaactaccattgacattcttcaca
gaatcagaaaaaactactttaaatttcatgtagaatcaaagaagaccctg
tatagccaagacaatcctaagcataaaaaacaaatggag

>{"tileID":"0157.00.0001.000","md5sum":"0dad59737cc1b788ca7f5d8e33b6cf90","tagmask_md5sum":"0dad59737cc1b788ca7f5d8e33b6cf90","locus":[{"build":"hg19 chr7 54000565 54000881"}],"n":316,"seedTileLength":1,"startTile":false,"endTile":false,"startSeq":"cctaagcataaaaaacaaatggag","endSeq":"aatattgagtcaagatagattaga","startTag":"cctaagcataaaaaacaaatggag","endTag":"aatattgagtcaagatagattaga","nocallCount":0,"notes":[]}
cctaagcataaaaaacaaatggagacatcatgctacctgacttcaaacta
tactacagtgctacagtaaccaaaacagcatggtactggtaccaaaacag
acatatagaccaaaaaggaacagaacagagacctcagaaataataccacg
catctacaaccatctgatcttcgacaaacctgacaataacaagcagtggg
gaaaggatctcctatttaataagtggtgctgggaaaactggctagccata
tgcagaaaactgaaactggaccccttccttacaccttatacaaatattga
gtcaagatagattaga

>{"tileID":"0157.00.0001.001","md5sum":"0dad59737cc1b788ca7f5d8e33b6cf90","tagmask_md5sum":"0dad59737cc1b788ca7f5d8e33b6cf90","locus":[{"build":"hg19 chr7 54000565 54000881"}],"n":316,"seedTileLength":1,"startTile":false,"endTile":false,"startSeq":"cctaagcataaaaaacaaatggag","endSeq":"aatattgagtcaagatagattaga","startTag":"cctaagcataaaaaacaaatggag","endTag":"aatattgagtcaagatagattaga","nocallCount":0,"notes":[ ]}
cctaagcataaaaaacaaatggagacatcatgctacctgacttcaaacta
tactacagtgctacagtaaccaaaacagcatggtactggtaccaaaacag
acatatagaccaaaaaggaacagaacagagacctcagaaataataccacg
catctacaaccatctgatcttcgacaaacctgacaataacaagcagtggg
gaaaggatctcctatttaataagtggtgctgggaaaactggctagccata
tgcagaaaactgaaactggaccccttccttacaccttatacaaatattga
gtcaagatagattaga

```

The following command would produce the concatenated sequence for variant `001`:

```
$ fjt -c 1 testdata/example-tile.fj
agaaaatgccaacatagccagagtgataaattaattctatagaccaacaa
gtcaaacataagaaaaagttggaaaaatttaccataccaattacatttta
attgtagtagtatatctccgtcatctcagcccaaaaactccttaagctga
taagcaacttcagcaaggtctcagcatacaaaatcaatgtgcaaaaatca
caagcattccttcacaccaacaatagacaagcagagagccaaatcatgaa
tgaactcccatttacaatagctacaaagagaataaaatacctaagaatac
agttaacaagggatgtgaaggacctcttcaaggagaactacaaaccactg
ctcaaggaaataatagaggacacaaacaaatggaaaaacgttccatcctc
atggataggaagaacaaatatcgtgaaaatggccatactgcccaaagtaa
ttaatagattcattgctattcccatcaaactaccattgacattcttcaca
gaatcagaaaaaactactttaaatttcatgtagaatcaaagaagaccctg
tatagccaagacaatcctaagcataaaaaacaaatggagacatcatgcta
cctgacttcaaactatactacagtgctacagtaaccaaaacagcatggta
ctggtaccaaaacagacatatagaccaaaaaggaacagaacagagacctc
agaaataataccacgcatctacaaccatctgatcttcgacaaacctgaca
ataacaagcagtggggaaaggatctcctatttaataagtggtgctgggaa
aactggctagccatatgcagaaaactgaaactggaccccttccttacacc
ttatacaaatattgagtcaagatagattaga
```

### Band Format `fjt -B`

Produce the tile information in 'band format'.
To produce the 'band format', an `SGLF` file needs to be specified.

For example, on the `testdata` provided, here is the result:

```bash
$ ./fjt -B testdata/035e.fj -L testdata/035e.sglf
[ 79 8 0 0 0 0 0 -1 0 0 0 389 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 1 0 0 -1 34 -1 185 1]
[ 79 8 0 0 0 0 0 -1 0 0 0 389 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 1 0 0 -1 34 -1 185 1]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
```

### Batch Band Format `fjt -B -I ...`

Produce 'band format' output for a series of input gziped FastJ files.
A stream with the list of gziped FastJ files is provided with the `-I` option.

```bash
$ ./fjt -B -I <( echo -e 'testdata/035e.fj.gz\ntestdata/035e.fj.gz' ) -L testdata/035e.sglf
[ 79 8 0 0 0 0 0 -1 0 0 0 389 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 1 0 0 -1 34 -1 185 1]
[ 79 8 0 0 0 0 0 -1 0 0 0 389 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 1 0 0 -1 34 -1 185 1]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
[ 79 8 0 0 0 0 0 -1 0 0 0 389 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 1 0 0 -1 34 -1 185 1]
[ 79 8 0 0 0 0 0 -1 0 0 0 389 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 1 0 0 -1 34 -1 185 1]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
```

Often, the SGLF file is large and is the bottleneck to processing.
Batching tries to speed up the band format conversion process by
loading the SGLF file once, at the beginning, into memory and then
processing each FastJ conversion serially.
This method also keeps the SGLF information in a 2bit representation
to try and reduce the memory footprint.


### CSV Output `fjt -C`

Output comma separated fields with `extended tileID`, `hash` and `sequence`.
The `hash` is the MD5SUM of the ASCII sequence with no trailing newline.
Note this is not `SGLF` output as the sequence can have `n`'s for portions of the
sequence that have not been confidently called.

For example, using the `035e.fj` file included in `testdata` would produce:

```
$ ./fjt -C testdata/035e.fj  | head -n4
035e.00.0021.000+1,fc72f8955140a5423023488dcd374b94,cacccaagtattgactcacccatcaacaaccgctatgtatttcgtacattactgccagcctccatgaatattgtacggtaccataaatacttgaccacctgtagtacataaaaacccaatccacatcaaaccccccccccccatgcttacaagcaagtacagcaatcaaccttcaactatcacacatcaactgcaactccaaagccacccctcacccactaggataccaacaaacctacccacccttaa
035e.00.0021.001+1,fc72f8955140a5423023488dcd374b94,cacccaagtattgactcacccatcaacaaccgctatgtatttcgtacattactgccagcctccatgaatattgtacggtaccataaatacttgaccacctgtagtacataaaaacccaatccacatcaaaccccccccccccatgcttacaagcaagtacagcaatcaaccttcaactatcacacatcaactgcaactccaaagccacccctcacccactaggataccaacaaacctacccacccttaa
035e.00.0022.000+1,5c5eb5980f21f9232dab6bc4fa066bd2,accaacaaacctacccacccttaacagtacatagtacataaagtcatttaccgtacatagcacattacagtcaaatcccttctcgtccccatggatgacccccctcagataggggtcccttgaccaccatcctccgtgaaatcaatatcccgcacaagagtgctactctcctcgctccgggcccataacacttgggggtagctaaagtgaactgtatccgacatctggttcctacttcagggccataaagcctaaatagcccacacgttccccttaaataagacatcacgann
035e.00.0022.001+1,5c5eb5980f21f9232dab6bc4fa066bd2,accaacaaacctacccacccttaacagtacatagtacataaagtcatttaccgtacatagcacattacagtcaaatcccttctcgtccccatggatgacccccctcagataggggtcccttgaccaccatcctccgtgaaatcaatatcccgcacaagagtgctactctcctcgctccgggcccataacacttgggggtagctaaagtgaactgtatccgacatctggttcctacttcagggccataaagcctaaatagcccacacgttccccttaaataagacatcacgann
```

### FastJ Output From Band Format `fjt -b`

Take in 'band format' and an SGLF file to produce a FastJ file.

For example:

```
$ ./fjt -b -L testdata/035e.sglf testdata/035e.band | head -n 40
>{"tileID":"035e.00.0000.000","md5sum":"7346f663d221ed28c112df86eb5986ef","tagmask_md5sum":"7346f663d221ed28c112df86eb5986ef","locus":[ ],"n":244,"seedTileLength":1,"startTile":true,"endTile":false,"startSeq":"gatcacaggtctatcaccctatta","endSeq":"gcttgtaggacataataataacaa","startTag":"gatcacaggtctatcaccctatta","endTag":"gcttgtaggacataataataacaa","nocallCount":0,"notes":[ ]}
gatcacaggtctatcaccctattaaccactcacgggagctctccatgcat
ttggtattttcgtctggggggcgtgcacgcgatagcattgcgggacgctg
gagccggagcaccctatgtcgcagtatctgtctttgattcctgcctcatt
ctattatttatcgcacctacgttcaatattacaggcgaacatacctacta
aagtgtgttaattaattaatgcttgtaggacataataataacaa

>{"tileID":"035e.00.0001.000","md5sum":"1cadbbf41d5898b9e37ecdfd1d751f4e","tagmask_md5sum":"1cadbbf41d5898b9e37ecdfd1d751f4e","locus":[ ],"n":253,"seedTileLength":1,"startTile":false,"endTile":false,"startSeq":"gcttgtaggacataataataacaa","endSeq":"acacattattttcccctcccactc","startTag":"gcttgtaggacataataataacaa","endTag":"acacattattttcccctcccactc","nocallCount":0,"notes":[ ]}
gcttgtaggacataataataacaattgaatgtctgcacagccgctttcca
cacagacatcataacaaaaaatttccaccaaacccccccccctctccccc
cgcttctggccacagcacttaaacacatctctgccaaaccccaaaaacaa
agaaccctaacaccagcctaaccagatttcaaattttatctttaggcggt
atgcacttttaacagtcaccccccaactaacacattattttcccctccca
ctc

>{"tileID":"035e.00.0002.000","md5sum":"1d46f1fff282a060f2e3b28592daa12e","tagmask_md5sum":"1d46f1fff282a060f2e3b28592daa12e","locus":[ ],"n":516,"seedTileLength":1,"startTile":false,"endTile":false,"startSeq":"acacattattttcccctcccactc","endSeq":"gtaaagagtgttttagatcacccc","startTag":"acacattattttcccctcccactc","endTag":"gtaaagagtgttttagatcacccc","nocallCount":0,"notes":[ ]}
acacattattttcccctcccactcccatactactaatctcatcaatacaa
cccccgcccatcctacccagcacacacacaccgctgctaaccccataccc
cgaaccaaccaaaccccaaagacaccccccacagtttatgtagcttacct
cctcaaagcaatacactgaaaatgtttagacgggctcacatcaccccata
aacaaataggtttggtcctagcctttctattagctcttagtaagattaca
catgcaagcatccccgttccagtgagttcaccctctaaatcaccacgatc
aaaagggacaagcatcaagcacgcagcaatgcagctcaaaacgcttagcc
tagccacacccccacgggaaacagcagtgattaacctttagcaataaacg
aaagtttaactaagctatactaaccccagggttggtcaatttcgtgccag
ccaccgcggtcacacgattaacccaagtcaatagaagccggcgtaaagag
tgttttagatcacccc

>{"tileID":"035e.00.0003.000","md5sum":"7c146e32156d9e4d3297d06705d72e4c","tagmask_md5sum":"7c146e32156d9e4d3297d06705d72e4c","locus":[ ],"n":459,"seedTileLength":1,"startTile":false,"endTile":false,"startSeq":"gtaaagagtgttttagatcacccc","endSeq":"tctaccccagaaaactacgatagc","startTag":"gtaaagagtgttttagatcacccc","endTag":"tctaccccagaaaactacgatagc","nocallCount":0,"notes":[ ]}
gtaaagagtgttttagatcaccccctccccaataaagctaaaactcacct
gagttgtaaaaaactccagttgacacaaaatagactacgaaagtggcttt
aacatatctgaacacacaatagctaagacccaaactgggattagataccc
cactatgcttagccctaaacctcaacagttaaatcaacaaaactgctcgc
cagaacactacgagccacagcttaaaactcaaaggacctggcggtgcttc
atatccctctagaggagcctgttctgtaatcgataaaccccgatcaacct
caccacctcttgctcagcctatataccgccatcttcagcaaaccctgatg
aaggctacaaagtaagcgcaagtacccacgtaaagacgttaggtcaaggt
gtagcccatgaggtggcaagaaatgggctacattttctaccccagaaaac
tacgatagc

```

The different alleles are not guaranteed to appear interleaved.
