#!/usr/bin/python
#
# usage:
# ./gt-sglf-to-band.py <( zcat src/012c.sglf.gz | cut -f1,2 -d, ) <( zcat gt/012c.sglf.gz )
#

import sys

if len(sys.argv)<3:
  print "provide sglf-input and sglf-lib."
  print "example:"
  print ""
  print "  /gt-sglf-to-band.py <( zcat src/012c.sglf.gz | cut -f1,2 -d, ) <( zcat tilelib/012c.sglf.gz )"
  print ""
  sys.exit(0)

max_tilestep = 0

srcfn = sys.argv[1]
libfn = sys.argv[2]


# populate the called positions for the genotype file.
# Format is:
#
# <hextileid>+<span> <pos0> <len0> <pos1> <len1> ...
#
# store in info_map
#
info_map = {}
if len(sys.argv)>=3:
  infofn = sys.argv[3]
  with open(infofn) as fp:
    for line in fp:
      line = line.strip()
      if len(line) == 0: continue
      if line[0] == '#': continue
      fields = line.split(" ")
      a = []
      for v in fields[1:]:
        a.append(int(v))
      info_map[fields[0]] = a



sglf_map = {}

with open(libfn) as fp:
  for line in fp:
    a = line.strip().split(",")
    tileid_parts = a[0].split("+")
    tileparts = tileid_parts[0].split(".")

    # add library tile variant to sglf_map of md5sum seq.
    #
    sglf_map[a[1]] = int(tileparts[3], 16)

    tilestep = int(tileparts[2], 16)
    if tilestep > max_tilestep: max_tilestep = tilestep

band0 = []
band1 = []
noc0 = []
noc1 = []

for s in range(max_tilestep+1):
  band0.append(0)
  band1.append(0)

  noc0.append([])
  noc1.append([])

# now go through input sglf and create band,
# filling in with default tile for non-existent tiles
#
with open(srcfn) as fp:
  for line in fp:
    a = line.strip().split(",")
    h = a[1]

    if h not in sglf_map:
      print a[0], a[1], h, "error!"

    #else: print a[0], a[1], sglf_map[h]
    tileid = a[0]

    tileid_parts = a[0].split("+")
    tileparts = tileid_parts[0].split(".")
    band_part = int(tileparts[3], 16)
    tilestep = int(tileparts[2], 16)

    if tilestep >= len(band0):
      print "error!  tilestep", tilestep, ">=", len(band0), "max_tilestep", max_tilestep

    if band_part == 0:
      band0[tilestep] = sglf_map[h]

      if tileid in info_map:
        noc0[tilestep] = info_map[tileid]
    elif band_part == 1:
      band1[tilestep] = sglf_map[h]
      if tileid in info_map:
        noc1[tilestep] = info_map[tileid]
    else:
      print "error! got bad band_part", band_part


sys.stdout.write("[")
for x in band0:
  sys.stdout.write(" " + str(x))
sys.stdout.write("]\n")

sys.stdout.write("[")
for x in band1:
  sys.stdout.write(" " + str(x))
sys.stdout.write("]\n")

sys.stdout.write("[")
for x in noc0:
  sys.stdout.write("[")
  for y in x:
    sys.stdout.write(" " + str(y))
  sys.stdout.write(" ]")
sys.stdout.write("]\n")

sys.stdout.write("[")
for x in noc1:
  sys.stdout.write("[")
  for y in x:
    sys.stdout.write(" " + str(y))
  sys.stdout.write(" ]")
sys.stdout.write("]\n")

