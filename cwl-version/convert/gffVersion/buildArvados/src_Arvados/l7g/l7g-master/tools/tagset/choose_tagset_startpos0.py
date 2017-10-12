#!/usr/bin/python
#
# An example tagset creation.
# Input a bedGraph file and output the start positions of tags, 0 reference.
# This program is meant to be run on a tile path by tile path basis and
# chromosome boundaries shouldn't be crossed.
#
# This program expects a bedGraph file as input that holds a value for 1 for
# where the tags should be able to start.
# Input is 0 reference and in the following format:
#
#   chr start_0ref  end_0ref_noninc value
#
# where each field is separated by a TAB.
#
# example usage:
#
#     bigWigToBedGraph -chrom=chr1 -start=0 -end=2300000 /path/to/wgEncodeCrgMapabilityAlign24mer.bw /dev/stdout | \
#      ./choose_tagset_startpos0.py
#
# bounds can be used from e.g. cytogenetic bands
#
# NOTE: this is provided as an example of how to generate the start positions of the tag reference.
# For generating the actual tagset in current use, please use 'choose_tagset_start0_vestigial.py'.
#

import sys

prevpos0 = 0
endpos0_noninc = -1
if len(sys.argv)>1:
  prevpos0 = int(sys.argv[1])
  if len(sys.argv)>2:
    endpos0_noninc = int(sys.argv[2])

taglen=24
midwindowlen=200
windowlen = 2*taglen + midwindowlen

for line in sys.stdin:
  line = line.strip()
  if len(line)==0: continue

  parts = line.split("\t")
  if len(parts) != 4: continue

  chrom       = parts[0]
  start0      = int(parts[1])
  end0_noninc = int(parts[2])
  val         = float(parts[3])

  # make sure the value is 1 (a 'unique' sequence of 24)
  #
  if abs(1.0-val) > (1.0/65536.0) : continue

  # don't allow for end tiles to be less than the window length
  #
  if endpos0_noninc > 0:
    de = endpos0_noninc - start0
    if de < windowlen: continue

  # If the start position is within the windowlen,
  # emit the start position (0ref).
  #
  dn = start0 - prevpos0 + 24
  if dn >= windowlen:
    print start0
    prevpos0 = start0
    continue

  # otehrwise emit the position that appears in the
  # middle of the run (0ref).
  #
  dn = end0_noninc - prevpos0 + 24
  if dn >= windowlen:
    prevpos0 = prevpos0 + windowlen - 24
    print prevpos0
    continue


