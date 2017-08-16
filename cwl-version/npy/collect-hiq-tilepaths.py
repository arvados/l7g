#!/usr/bin/python
#
# collect all hiq tiles into single vector (and info vector)
#

import os
import re
import sys
import numpy as np

obase="oot"

if len(sys.argv)<2:
  print "provide dir"
  sys.exit(-1)

idir = sys.argv[1]
if len(sys.argv)>2:
  obase = sys.argv[2]

pathlist = []
for fn in os.listdir(idir):
  if not re.match(r'[a-f0-9]{3}$', fn): continue
  pathlist.append(fn)
pathlist.sort()

collect_tilevar = []
tm = np.load(idir + "/" + pathlist[0])
for ii in range(tm.shape[0]):
  collect_tilevar.append( np.array([]) )
collect_tileinfo = np.array([])

for x in pathlist:
  v = np.load(idir + "/" + x)
  vinfo = np.load(idir + "/" + x + "-info")

  print ">>>", x

  for ds in range(v.shape[0]):
    for ii in range(v.shape[1]):
      collect_tilevar[ds] = np.append(collect_tilevar[ds], v[ds][ii])

  for ii in range(vinfo.shape[0]):
    collect_tileinfo = np.append(collect_tileinfo, vinfo[ii])


y = np.ndarray( (len(collect_tilevar), collect_tilevar[0].shape[0]) )
for ii in range(len(collect_tilevar)):
  for jj in range(collect_tilevar[0].shape[0]):
    y[ii][jj] = collect_tilevar[ii][jj]
np.save( obase + "-collect", y )
np.save( obase + "-collect-info", collect_tileinfo )

