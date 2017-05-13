FastJ Tool (fjt)
===

`fjt` is a tool to manipulate FastJ (text) files.
`fjt` is still experimental so please use with caution.

Quick Start
---

```
$ make
$ ./fjt -h
fjt version: 0.1.1
usage:
  fjt [-c variant] [-C] [-v] [-V] [-h] [input]

  [-C]            Output comma separated `extended tileID`, `hash` and `sequence` (CSV output)
  [-B]            Output band format
  [-c variant]    Concatenate FastJ tiles into sequence.  `variant` is the variant ID to concatenate on
  [-L sglf]       Simple genome library format tile path file
  [-i ifn]        input file
  [-p tilepath]   Tile path (in decimal)
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
[ 79 2 0 0 0 0 0 -1 0 0 0 390 0 0 0 0 0 1 0 0 0 0 0 0 26 0 0 1 0 0 -1 34 -1 185 1]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
[[ ][ ][ ][ ][ ][ ][ 903 1 ][ ][ 16 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 96 1 ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ ][ 291 2 ]]
```


### CSV output `fjt -C`

To produce the tiles in `sglf` format (simple genome library format), outputing the tiles in CSV can be done.
The output CSV is `extended tileID`, `hash` and `sequence`.

For example, using the `035e.fj` file included in `testdata` would produce:

```
$ ./fjt -C testdata/035e.fj  | head -n4
035e.00.0000.000+1,7346f663d221ed28c112df86eb5986ef,gatcacaggtctatcaccctattaaccactcacgggagctctccatgcatttggtattttcgtctggggggcgtgcacgcgatagcattgcgggacgctggagccggagcaccctatgtcgcagtatctgtctttgattcctgcctcattctattatttatcgcacctacgttcaatattacaggcgaacatacctactaaagtgtgttaattaattaatgcttgtaggacataataataacaa
035e.00.0000.001+1,7346f663d221ed28c112df86eb5986ef,gatcacaggtctatcaccctattaaccactcacgggagctctccatgcatttggtattttcgtctggggggcgtgcacgcgatagcattgcgggacgctggagccggagcaccctatgtcgcagtatctgtctttgattcctgcctcattctattatttatcgcacctacgttcaatattacaggcgaacatacctactaaagtgtgttaattaattaatgcttgtaggacataataataacaa
035e.00.0001.000+1,1cadbbf41d5898b9e37ecdfd1d751f4e,gcttgtaggacataataataacaattgaatgtctgcacagccgctttccacacagacatcataacaaaaaatttccaccaaacccccccccctctccccccgcttctggccacagcacttaaacacatctctgccaaaccccaaaaacaaagaaccctaacaccagcctaaccagatttcaaattttatctttaggcggtatgcacttttaacagtcaccccccaactaacacattattttcccctcccactc
035e.00.0001.001+1,d69eb38a5317a1f387c49c84a003df36,gcttgtaggacataataataacaattgaatgtctgcacagccgctttccacacagacatcataacaaaaaatttccaccaaaccccccccctctccccccgcttctggccacagcacttaaacacatctctgccaaaccccaaaaacaaagaaccctaacaccagcctaaccagatttcaaattttatctttaggcggtatgcacttttaacagtcaccccccaactaacacattattttcccctcccactc
```


