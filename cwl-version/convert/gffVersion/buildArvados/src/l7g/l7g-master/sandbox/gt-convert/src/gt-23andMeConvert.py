#!/usr/bin/python
#
# Covert a 23andMe file to
# a space delimited file
# consisting of:
#
# chrom pos0ref tilepath tilestep gt rsid pos1ref
#
# chrom - ucsc chromosome name (chr1-chr22, chrX, chrY, chrM)
# pos0ref - position in chromosome, 0 reference
# tilepath - tilepath this falls on
# tilestep - tilestep this falls on
# gt - the genotype as reported by 23andMe
# rsid - the rsid as reported by 23andMe
# pos1ref - the position as reported by 23andMe (1 reference)
#
# For example, this might be output:
#
# chr1 82153 0000 000a AA rs4477212 82154
# chr1 752565 0000 003f AG rs3094315 752566
# chr1 752720 0000 003f AG rs3131972 752721
# chr1 776545 0000 0065 -- rs12124819 776546
# chr1 798958 0000 0085 GG rs11240777 798959
# ...
#
# to run:
#
#    ./gt-23andMeConvert.py 23andMe-genotyping.txt > gt-list.ssv
#
# The tile assembly file is hard coded to be
# `/data-sdd/data/l7g/assembly/assembly.00.hg19.fw.gz`.
#
# notes:
# * I and D are for insertion and deletions.  These need a library to look up the appropriate
#   types of insertion and deletions
# * A '-' represents a no-call
# * The variants are unphased
# * This script assumes hg19


import gzip
import sys

if len(sys.argv)<2:
  print "provide genotyping (23andMe) file"
  sys.exit(1)

#ifn = "./data/23andMe-genotyping.txt"
ifn = sys.argv[1]

tile_assembly_fn = "/data-sdd/data/l7g/assembly/assembly.00.hg19.fw.gz"

fp = gzip.open(tile_assembly_fn)

pos_map = {}
assembly_map = {}

uniq_chrom = {}

chrom = 'unk'
prev_chrom = ''
tilepath = '-1'
tilestep = '-1'
prev_pos = 0

# make tile assembly map
# key is <chrom>:<pos>
# non header lines in assembly are formatted as:
#
# <tilestep>\t<ws>end_pos_exclusive
#
for line in fp:
  line = line.strip()
  if len(line)==0: continue

  if line[0] == '>':
    fields = line.split(":")
    chrom = fields[1]
    tilepath = fields[2]

    if prev_chrom != chrom:
      prev_pos = 0

    # There is no start tag for tiles at the
    # beginning of a tile path.  Increment
    # so calculation works out below.
    #
    prev_pos += 24

    uniq_chrom[chrom] = True

    if not (chrom in pos_map):
      pos_map[chrom] = []

    prev_chrom = chrom

    continue

  fields = line.split("\t")
  fields[0] = fields[0].strip()
  fields[1] = fields[1].strip()
  tilestep = fields[0]
  end_pos_excl = int(fields[1])

  pos_map[chrom].append(end_pos_excl)
  assembly_map[ chrom + ":" + str(end_pos_excl) ] = [ tilepath , tilestep , prev_pos-24, end_pos_excl ]

  prev_pos = end_pos_excl

fp.close()

#print "# assembly read..."
print "#chrom,pos0ref,tilepath,tilestep,gt,hg19_start+len,span_flag,rsid,pos1ref"

count=0

cur_chrom = "unk"
idx = 0

prev_emit_chrom = 'unk'
prev_emit_pos = 0

fp = open(ifn)
for line in fp:
  line = line.strip()
  if line[0] == '#': continue

  fields = line.split("\t")
  rsid = fields[0]
  chrom = "chr" + fields[1]
  pos1ref = fields[2]
  gt = fields[3]

  # skip indels and nocalls
  #
  if len(gt)>0:
    if gt[0] == 'I' or gt[0] == 'D' or gt[0] == '-':
      continue
    if len(gt)>1:
      if gt[1] == 'I' or gt[1] == 'I' or gt[1] == '-':
        continue

  if chrom ==  "chrMT": chrom = "chrM"

  pos0ref = int(pos1ref)-1

  if chrom != cur_chrom:
    cur_chrom = chrom
    idx = 0

  #print "## chrom", chrom, "idx", idx, "pos0ref", pos0ref, "pos_map", pos_map[cur_chrom][idx]

  while (idx < len(pos_map[cur_chrom])) and (pos_map[cur_chrom][idx] <= pos0ref):
    idx+=1

  #print "## >> chrom", chrom, "idx", idx, "pos0ref", pos0ref, "pos_map", pos_map[cur_chrom][idx]

  key = cur_chrom + ":" + str(pos_map[cur_chrom][idx])
  a = assembly_map[key]

  tilepath = a[0]
  tilestep = a[1]
  prev_pos = a[2]
  end_pos_excl = a[3]

  tile_len = end_pos_excl - prev_pos

  span_flag = ''
  if (pos0ref - prev_pos) < 24:
    span_flag += '*'
  else:
    span_flag += '.'

  from_end = end_pos_excl - pos0ref
  if from_end < 24:
    span_flag += '*'
  else:
    span_flag += '.'

  # remove duplicate entries...
  #
  if (chrom != prev_emit_chrom) or (prev_emit_pos != pos0ref):
    print chrom, pos0ref, tilepath, tilestep, gt, str(prev_pos) + "+" + str(tile_len), span_flag, rsid, pos1ref

  count+=1
  #if count>10: break

  prev_emit_chrom = chrom
  prev_emit_pos = pos0ref

fp.close()
