#!/usr/bin/python

import sys
import os

if len(sys.argv)>1:
  ifn = sys.argv[1]
else:
  print "provide input"
  sys.exit(-1)

fp = open(ifn)
for line in fp:
  line = line.strip()
  if len(line)==0: continue
  if line[0] == '#': continue
  fields = line.split(",")
  seq = fields[2]

  sys.stdout.write(fields[0])
  for pos in range(len(seq)):
    if seq[pos:pos+1] == 'A' or seq[pos:pos+1] == 'C' or seq[pos:pos+1] == 'G' or seq[pos:pos+1] == 'T':
      sys.stdout.write(" " + str(pos) + " 1")
  sys.stdout.write("\n")

fp.close()
