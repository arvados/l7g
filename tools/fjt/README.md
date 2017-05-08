FastJ Tool (fjt)
===

`fjt` is a tool to manipulate FastJ (text) files.
`fjt` is still experimental so please use with caution.

Quick Start
---

```
$ make
$ ./fjt -h
fjt version: 0.1.0
usage:
  [-c variant]    Concatenate FastJ tiles into sequence.  `variant` is the variant ID to concatenate on.
  [-C]            Output comma separated `tileID`, `hash` and `sequence` (CSV output).
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
$ fjt -c example-tile.fj
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

### CSV output `fjt -C`

To produce the tiles in `sglf` format (simple genome library format), outputing the tiles in CSV can be done.
The output CSV is `tileID`, `hash` and `sequence`.

For example, using the `example-tile.fj` file above would produce:

```
$ fjt -C example-tile.fj
0157.00.0000.0000,d7c151e05a4be1a735af3d1346c90714,agaaaatgccaacatagccagagtgataaattaattctatagaccaacaagtcaaacataagaaaaagttggaaaaatttaccataccaattacattttaattgtagtagtatatctccgtcatctcagcccaaaaactccttaagctgataagcaacttcagcaaggtctcagcatacaaaatcaatgtgcaaaaatcacaagcattccttcacaccaacaatagacaagcagagagccaaatcatgaatgaactcccatttacaatagctacaaagagaataaaatacctaagaatacagttaacaagggatgtgaaggacctcttcaaggagaactacaaaccactgctcaaggaaataatagaggacacaaacaaatggaaaaacgttccatcctcatggataggaagaacaaatatcgtgaaaatggccatactgcccaaagtaattaatagattcattgctattcccatcaaactaccattgacattcttcacagaatcagaaaaaactactttaaatttcatgtagaatcaaagaagaccctgtatagccaagacaatcctaagcataaaaaacaaatggag
0157.00.0000.0001,d7c151e05a4be1a735af3d1346c90714,agaaaatgccaacatagccagagtgataaattaattctatagaccaacaagtcaaacataagaaaaagttggaaaaatttaccataccaattacattttaattgtagtagtatatctccgtcatctcagcccaaaaactccttaagctgataagcaacttcagcaaggtctcagcatacaaaatcaatgtgcaaaaatcacaagcattccttcacaccaacaatagacaagcagagagccaaatcatgaatgaactcccatttacaatagctacaaagagaataaaatacctaagaatacagttaacaagggatgtgaaggacctcttcaaggagaactacaaaccactgctcaaggaaataatagaggacacaaacaaatggaaaaacgttccatcctcatggataggaagaacaaatatcgtgaaaatggccatactgcccaaagtaattaatagattcattgctattcccatcaaactaccattgacattcttcacagaatcagaaaaaactactttaaatttcatgtagaatcaaagaagaccctgtatagccaagacaatcctaagcataaaaaacaaatggag
0157.00.0001.0000,0dad59737cc1b788ca7f5d8e33b6cf90,cctaagcataaaaaacaaatggagacatcatgctacctgacttcaaactatactacagtgctacagtaaccaaaacagcatggtactggtaccaaaacagacatatagaccaaaaaggaacagaacagagacctcagaaataataccacgcatctacaaccatctgatcttcgacaaacctgacaataacaagcagtggggaaaggatctcctatttaataagtggtgctgggaaaactggctagccatatgcagaaaactgaaactggaccccttccttacaccttatacaaatattgagtcaagatagattaga
0157.00.0001.0001,0dad59737cc1b788ca7f5d8e33b6cf90,cctaagcataaaaaacaaatggagacatcatgctacctgacttcaaactatactacagtgctacagtaaccaaaacagcatggtactggtaccaaaacagacatatagaccaaaaaggaacagaacagagacctcagaaataataccacgcatctacaaccatctgatcttcgacaaacctgacaataacaagcagtggggaaaggatctcctatttaataagtggtgctgggaaaactggctagccatatgcagaaaactgaaactggaccccttccttacaccttatacaaatattgagtcaagatagattaga
```


