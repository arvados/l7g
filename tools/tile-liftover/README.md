Tile Liftover
===

A tool to create 'assembly' coordinate information
for a given tagset and haplotype stream.

This tool has minimal features and is a work in progress.

Quick Start
---

```
git clone https://github.com/curoverse/l7g
cd l7g/tools/tile-liftover
make
```

```
$ ./tile-liftover
tile-liftover version: 0.1.0
usage:
  -T tagset       tagset stream
  -p tilepath     tilepath
  [-R ref]        reference stream (stdin default)
  [-s start]      start position (0 reference, 0 default)
  [-c chrom]      chromosome
  [-N refname]    reference name (default 'hg19')
  [-v]            verbose
  [-V]            version
  [-h]            help
```


Example Usage
---

Create a Lightning assembly for a single Lightning tile path:

```
$ refstream chrM | ./tile-liftover -T <( refstream $LIGHTNING_DATA_DIR/tagset.fa/tagset.fa.gz 035e.00 ) -c chrM -p 862 -N hg19
>hg19:chrM:035e
0000           244
0001           469
0002           961
0003          1396
0004          1769
0005          2227
0006          2722
0007          3114
0008          3339
0009          3564
000a          3789
000b          9761
000c          9986
000d         10214
000e         10603
000f         10924
0010         11294
0011         11790
0012         12021
0013         12360
0014         12632
0015         12894
0016         13490
0017         13715
0018         13940
0019         14328
001a         14582
001b         14877
001c         15112
001d         15339
001e         15600
001f         15835
0020         16077
0021         16302
0022         16571
```

