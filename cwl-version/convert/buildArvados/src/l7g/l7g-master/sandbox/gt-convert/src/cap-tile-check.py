#!/usr/bin/python
#
# check the output of gt-ssv-tile-single.py
# to make sure we can get back our genotype data
#

import os
import sys
import gzip
import json
import md5

ifn = "oot"
if len(sys.argv)>1:
  ifn = sys.argv[1]

tile_assembly_fn = "/data-sdd/data/l7g/assembly/assembly.00.hg19.fw.gz"


pos_map = {}
assembly_map = {}
tilestep_map = {}

uniq_chrom = {}

chrom = 'unk'
prev_chrom = ''
tilepath = '-1'
tilestep = '-1'
prev_pos = 24
end_pos_excl = -1

# make tile assembly map
# key is <chrom>:<pos>
# non header lines in assembly are formatted as:
#
# <tilestep>\t<ws>end_pos_excllusive
#
fp = gzip.open(tile_assembly_fn)
for line in fp:
  line = line.strip()
  if len(line)==0: continue

  if line[0] == '>':
    fields = line.split(":")
    chrom = fields[1]
    tilepath = fields[2]

    if chrom != prev_chrom:
      prev_pos = 0
    end_pos_excl = -1

    prev_pos += 24

    uniq_chrom[chrom] = True

    if not (chrom in pos_map):
      pos_map[chrom] = []

    print line
    continue

  fields = line.split("\t")
  fields[0] = fields[0].strip()
  fields[1] = fields[1].strip()
  tilestep = fields[0]
  end_pos_excl = int(fields[1])

  pos_map[chrom].append(end_pos_excl)
  assembly_map[ chrom + ":" + str(end_pos_excl) ] = [ tilepath , tilestep , prev_pos-24, end_pos_excl ]
  tilestep_map[ tilepath + ".00." + tilestep  ]  = [ tilepath , tilestep , prev_pos-24, end_pos_excl, chrom ]


  prev_pos = end_pos_excl
  prev_chrom = chrom

fp.close()

print "#..."

alt_map = {}
alt_map["A"] = True
alt_map["C"] = True
alt_map["G"] = True
alt_map["T"] = True

fp = open(ifn)
for line in fp:
  line = line.strip()
  if len(line)==0: continue
  if line[0] == '#': continue

  fields = line.split(",")
  tilepos = fields[0]
  tilehash = fields[1]
  tileseq = fields[2]

  tp = ".".join(tilepos.split(".")[0:3])

  ent = tilestep_map[tp]
  spos = ent[2]
  epos = ent[3]
  chrom = ent[4]
  n = epos - spos


  for pos in range(len(tileseq)):
    if tileseq[pos:pos+1] in alt_map:
      pos1ref = spos+pos+1
      print chrom, pos1ref, tileseq[pos:pos+1], ent[0], ent[1]


fp.close()

