#!/usr/bin/python
#
# Take in vcf and bed file and compare
# against location in the `out_vcf` file.
#
# All start positions in the input VCF and BED
# file should appear in the output vcf.
#

import sys
import re

if (len(sys.argv) !=4 ):
  print "provide input vcf and bed file and output vcf file"
  sys.exit(1)

vcf = sys.argv[1]
bed = sys.argv[2]
out_vcf = sys.argv[3]

bed_start_lookup = {}
in_end_lookup = {}

in_vcf_lookup = {}
in_vcf_end_lookup = {}
out_vcf_lookup = {}
out_vcf_end_lookup = {}

with open(bed) as fp:
  for line in fp:
    line = line.strip()
    field = line.split("\t")
    chrom = field[0]
    s0 = int(field[1])
    e0noninc = int(field[2])

    s1 = s0+1
    e1inc = e0noninc

    bed_start_lookup[ str(chrom) + ":" + str(s1) ] = 1
    in_end_lookup[ str(chrom) + ":" + str(e1inc) ] = 1


with open(vcf) as fp:
  for line in fp:
    if len(line)==0: continue
    if line[0]=='#': continue
    field = line.split("\t")
    chrom = field[0]
    s1 = field[1]
    in_vcf_lookup[chrom + ":" + str(s1) ] = 1
    m = re.search(r'END=(\d+)', field[7])
    if (m):
      e1inc = m.group(1)
      in_end_lookup[chrom +":" + e1inc] = 1

with open(out_vcf) as fp:
  for line in fp:
    if len(line)==0: continue;
    if line[0] == '#': continue
    field = line.split("\t")
    chrom = field[0]
    s1 = field[1]
    out_vcf_lookup[chrom + ":" + str(s1) ] = 1
    m = re.search(r'END=(\d+)', field[7])
    if (m):
      e1inc = m.group(1)
      out_vcf_end_lookup[chrom +":" + e1inc] = 1


for x in bed_start_lookup:
  if x in out_vcf_lookup:
    pass
  else:
    print x, "NOT FOUND in outvcf but found in bed"
    sys.exit(-1)

for x in in_vcf_lookup:
  if x in out_vcf_lookup:
    #print x, "outvcf + invcf: ok"
    pass
  else:
    print x, "NOT FOUND in outvcf but found invcf"
    sys.exit(-1)

for x in in_end_lookup:
  if x in out_vcf_end_lookup:
    #print x, "outvcf end + bed bed: ok"
    pass
  else:
    print x, "NOT FOUND in outvcf end but found in bed end"
    sys.exit(-1)

sys.exit(0)
