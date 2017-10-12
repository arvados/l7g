#!/usr/bin/python

import sys

if len(sys.argv)<=1:
  print "provide input file"
  sys.exit()

ifn=sys.argv[1]

prev_tilepath = ""

ofn = ""
ofp = ""

count = 0

with open(ifn) as fp:
  for line in fp:
    if len(line)==0: continue
    if line[0] == '#': continue

    tilepath = line[0:4]

    if prev_tilepath != tilepath:
      print tilepath

      if prev_tilepath != "":
        ofp.close()

      #ofn = "odata/" + tilepath
      ofn = "odata/" + tilepath + ".sglf"
      ofp = open(ofn, "w")


    ofp.write(line)

    #count += 1
    #if count > 100:
    #  ofp.close()
    #  sys.exit()

    prev_tilepath = tilepath



ofp.close()

