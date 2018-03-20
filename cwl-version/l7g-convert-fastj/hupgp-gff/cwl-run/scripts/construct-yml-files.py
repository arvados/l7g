#!/usr/bin/python

import sys

hg19_list = []
g1k_list = []

pfx = """
script:
  class: File
  path: ../src/convert_gff_to_fastj.sh

tagset:
  class: File
  path: keep:c102dcaa86972313f32f0ccff37358a0+174/tagset.fa.gz

"""

#tileassembly:
#  class: File
#  path: keep:98c5e71956730c36cc89bb25b99fe58b+965/assembly.00.hg19.fw.gz
#
#refFaFn:
#  class: File
#  path: keep:320d7a6717aa7b555a09e52976ba4a02+18534/hg19.fa.gz
#
#gffFns:
#  - class: File
#    path: keep:8122d87ab406a6475c6826ae8bcfded2+56754/hu826751-GS03052-DNA_B01.gff.gz
#  - class: File
#    path: keep:8122d87ab406a6475c6826ae8bcfded2+56754/hu0211D6-GS01175-DNA_E02.gff.gz



for line in sys.stdin:
  line = line.strip()
  if len(line) == 0 : continue

  fields = line.split(",")
  if fields[1] == "hg19":
    hg19_list.append(fields[0])
  elif fields[1] == "g1k_v37":
    g1k_list.append(fields[0])
  else:
    print "ERROR: unknown line", line


ofp = open("scatter_hupgp-gff-to-fastj_hg19.yml", "w")
ofp.write(pfx)
ofp.write("tileassembly:\n")
ofp.write("  class: File\n")
ofp.write("  path: keep:98c5e71956730c36cc89bb25b99fe58b+965/assembly.00.hg19.fw.gz\n")
ofp.write("\n")

ofp.write("refFaFn:\n")
ofp.write("  class: File\n")
ofp.write("  path: keep:320d7a6717aa7b555a09e52976ba4a02+18534/hg19.fa.gz\n")
ofp.write("\n")

ofp.write("gffFns:\n")
for ds in hg19_list:
  ofp.write("  - class: File\n")
  ofp.write("    path: keep:8122d87ab406a6475c6826ae8bcfded2+56754/" + ds + "\n")
ofp.close()

###
###

ofp = open("scatter_hupgp-gff-to-fastj_human_g1k_v37.yml", "w")
ofp.write(pfx)
ofp.write("tileassembly:\n")
ofp.write("  class: File\n")
ofp.write("  path: keep:98c5e71956730c36cc89bb25b99fe58b+965/assembly.00.human_g1k_v37.fw.gz\n")
ofp.write("\n")

ofp.write("refFaFn:\n")
ofp.write("  class: File\n")
ofp.write("  path: keep:320d7a6717aa7b555a09e52976ba4a02+18534/human_g1k_v37.fasta.gz\n")
ofp.write("\n")

ofp.write("gffFns:\n")
for ds in g1k_list:
  ofp.write("  - class: File\n")
  ofp.write("    path: keep:8122d87ab406a6475c6826ae8bcfded2+56754/" + ds + "\n")
ofp.close()


