#!/usr/bin/python

import sys
import re

if len(sys.argv)!=2:
  print "provide input vcf file to check"
  sys.exit(1)


ifn=sys.argv[1]

line_no=0

with open(ifn) as fp:
  chrom = ""
  pos1ref = 0

  for line in fp:
    line = line.strip()
    line_no+=1

    if len(line)==0: continue
    if line[0]=='#': continue

    field = line.split("\t")

    ch = field[0]
    s1 = int(field[1])
    e1inc = s1

    if (chrom != ch):
      chrom = ch
      pos1ref = s1

    m = re.search(r'END=(\d+)', field[7])
    if (m):
      e1inc = int(m.group(1))

    if e1inc < s1:
      print "ERROR: line_no", line_no, ", (", field[0], field[1], field[2], ")", "END (", e1inc, ") < start (", s1, ")"
      sys.exit(-1)

    if pos1ref > s1:
      print "ERROR: line_no", line_no, ", (", field[0], field[1], field[2], ")", "previous pos (", pos1ref, ") < start (", s1, ")"
      sys.exit(-1)

sys.exit(0)
