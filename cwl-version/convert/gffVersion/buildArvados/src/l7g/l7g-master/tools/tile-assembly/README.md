tile-assembly tool
===

A tool to help extract information from the Tile Assembly files.

To compile, make sure `htslib` is installed the install location is updated
in the `Makefile`.

Quick Start
---

Assuming the shared libraries for `htslib` are setup properly:

```
$ ./tile-assembly range $TILE_ASSEMBLY_FILE 035e
#nstep  beg     end     chrom_name      ref_name
34      0       16571   chrM    hg19
```

```
$ ./tile-assembly tilepath $TILE_ASSEMBLY_FILE 035e
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
