#!/usr/bin/python
#
# Tagset start position for current tagset (2017-06-28)
#
# example usage:
#
#     bigWigToBedGraph -chrom=chr1 -start=0 -end=2300000 /path/to/wgEncodeCrgMapabilityAlign24mer.bw /dev/stdout | \
#      ./choose_tagset_startpos0_vestigial.py
#

import sys

prevpos0 = 0
seqEnd = -1
if len(sys.argv)>1:
  prevpos0 = int(sys.argv[1])
  if len(sys.argv)>2:
    seqEnd = int(sys.argv[2])

minTileDistance = 200

mergeLastFlag = False

minEndSeqPos = prevpos0 + minTileDistance
nextTagStart = prevpos0
prevTagStart = -1

for line in sys.stdin:
  line = line.strip()
  if len(line)==0: continue

  parts = line.split("\t")
  if len(parts) != 4: continue

  chrom = parts[0]
  start0 = int(parts[1])
  end0_noninc = int(parts[2])
  val = float(parts[3])

  # make sure the value is 1 (a 'unique' sequence of 24)
  #
  if abs(1.0-val) > (1.0/65536.0) : continue

  if end0_noninc < minEndSeqPos: continue
  prevTagStart = nextTagStart
  nextTagStart = start0+1
  if start0 < minEndSeqPos: nextTagStart = minEndSeqPos+1
  if (seqEnd - nextTagStart) <= minTileDistance:
    mergeLastFlag = True
    break

  print prevTagStart

  minEndSeqPos = nextTagStart + minTileDistance + 24

if not mergeLastFlag: prevTagStart = nextTagStart
print prevTagStart
