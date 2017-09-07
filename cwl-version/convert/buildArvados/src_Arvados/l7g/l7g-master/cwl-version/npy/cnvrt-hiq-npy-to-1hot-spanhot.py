#!/usr/bin/python
#
# Converts from integer interleaved tilepath to double one-hot with nan.
# The new 1-hot numpy arrays don't have sub arrays in them.
# The non-anchor spanning tiles will have a 1-hot representation
# that is `anchor-tile * 2^{max spanning tile}` so that it gets counted
# in any dot product.
#
# The general strategy is to do a first pass to find the maximum value
# for each tile step position.  Go through again keeping a running vector
# of the anchor tile variant id and the current spanning tile position.
# For example, most will have spanning tile position 0 because they won't
# be spanning tiles but a few will have 1 or even 2 or more.
#
# For non-anchor spanning tiles, the value should be:
#
#    SpanningVal_{step} = MaxTileVariant_{step} * MaxVal_{step}^{step - AnchorStepPos}
#
# For example, if we have (a partial knot?) [3, -1, -1], the encodings might be
#
#   3 -> [ 0, 0, 0, 1 ]
#  -1 -> [ 0, 0,            0, 0, 0, 1 ]
#  -1 -> [ 0, 0, 0, 0, 0,   0, 0, 0, 0, 0,     0, 0, 0, 1 ]
#
# if the max tile variant was 1 for the first non-anchor spanning tile position and
# the maximum tile variant was 4 for the second non-anchor spanning tile position.
#
## still need to verify
## before we forget:
## it looks like there's a bug when trying to encode spanning tiles of the form:
##
##            _ allele a pos in vector
##           |
##           |    _ allele b pos in vector
##           |   |
##           |   |   _ allele a 1hot length
##           |   |  |
##           |   |  |  _ allele b 1hot length
##           |   |  | |
## 1hotlen: 262 263 2 2  (0,0) (1,1) (1,0) (0,0) (1,0)
## 1hotlen: 264 265 2 2  (-1,-1) (0,0) (0,-1) (-1,-1) (0,-1)
## 1hotlen: 266 267 1 1  (0,0) (0,0) (0,0) (0,0) (0,0)
##
## ... when you have something like [ (1,0), (0,-1) ] ?
##
## turns out it was the decoding in the check.  encoding looks to be good.
##
## The addition here (hiq) should be another 'current' vector which holds
## a t/f value of whether we've seen a non-spanning tile without a jump.


import os
import sys
import numpy as np
import json
import re

def one_hot(v,n):
  oh = np.zeros(n)

  # debug
  if (v >= len(oh)):
    print "BAD, v:", v, "n:", n
    return []

  oh[v] = 1
  return oh

ifn="input.npy"
ofn="output"


if len(sys.argv)>1:
  ifn=sys.argv[1]
  if len(sys.argv)>2:
    ofn=sys.argv[2]
else:
  print "provide ifn ofn"
  sys.exit(0)

x = np.load(ifn)
ifn_info = re.sub("\.npy", "", ifn) + "-info.npy"
tilevec_info = np.load(ifn_info)

dataset = []
for ds in range(x.shape[0]):
  dataset.append( np.array([]) )

info_vec = np.array([])

prev_maxstepsize = 0

N_DS = x.shape[0]
N_TILEVEC_IDX = x.shape[1]

#print ">>>>>", "ds:", N_DS, "tilestep:", N_TILEVEC_IDX

# Do an initial pass to find the maximum values for
# each tile step position
#
basevec = []
curtilevarid = [ [], []]
anchorpos = [ [], [] ]
valid_tilevec = [ [], [] ]

for ds in range(N_DS):
  curtilevarid[0].append(x[ds][0])
  curtilevarid[1].append(x[ds][1])
  anchorpos[0].append(0)
  anchorpos[1].append(0)
  valid_tilevec[0].append(False)
  valid_tilevec[1].append(False)

prev_tilestep = -2

# find maximum per index stored in basevec
#
for idx in range(0, N_TILEVEC_IDX, 2):

  tilestep = tilevec_info[idx]

  for ds in range(N_DS):
    for allele in range(0,2):
      if x[ds][idx+allele] == -2:
        curtilevarid[allele][ds] = 0
        valid_tilevec[allele][ds] = False
      elif x[ds][idx+allele] == -1:

        if tilestep != (prev_tilestep+1):
          valid_tilevec[allele][ds] = False

        #print "keeping", "ds:", ds, "ts:", idx+allele, "v:", x[ds][idx+allele], "curtilevarid:", curtilevarid[allele][ds]
        pass
      elif x[ds][idx+allele] >= 0:
        curtilevarid[allele][ds] = x[ds][idx+allele]

        valid_tilevec[allele][ds] = True

        #print "setting", "a:", allele, "ds:", ds, "ctvid:", curtilevarid[allele][ds]


  curmax = 0
  for ds in range(N_DS):

    #curmax = max(curmax, x[ds][idx])
    #curmax = max(curmax, x[ds][idx+1])

    if valid_tilevec[0][ds]:
      curmax = max(curmax, curtilevarid[0][ds]+1)

    if valid_tilevec[1][ds]:
      curmax = max(curmax, curtilevarid[1][ds]+1)

  basevec.append(curmax)
  basevec.append(curmax)

  prev_tilestep = tilestep

  ## DEBUG
  #zz = ""
  #for ds in range(N_DS):
  #  zz += " (" + str(x[ds][idx]) + "," + str(x[ds][idx+1]) + ")"
  #print "adding basevec", idx, idx+1, basevec[idx], basevec[idx+1], zz

prev_tilestep = -2
for ds in range(N_DS):
  valid_tilevec[0][ds] = False
  valid_tilevec[1][ds] = False

onehotvecsize = []
for idx in range(0, N_TILEVEC_IDX, 2):
  tilestep = tilevec_info[idx]
  onehotvecsize.append(basevec[idx])
  onehotvecsize.append(basevec[idx+1])

for idx in range(0, N_TILEVEC_IDX, 2):

  tilestep = tilevec_info[idx]

  # update curtilevar and anchorpos
  #
  for ds in range(N_DS):
    for allele in range(0,2):
      if x[ds][idx+allele] >= 0:
        anchorpos[allele][ds] = idx + allele
        curtilevarid[allele][ds] = x[ds][idx+allele]

        valid_tilevec[allele][ds] = True

      elif x[ds][idx+allele] == -2:
        anchorpos[allele][ds] = idx+allele
        curtilevarid[allele][ds] = -2

        valid_tilevec[allele][ds] = False

      ## if it's a non-anchor spanning tile and we've skipped a tile step
      ## then we know to make it as invalid.  If it's a non-anchor spanning
      ## tile but the prev_tilestep is one away from the current tilestep,
      ## then we inherit whatever the value was, valid or no.
      ##
      elif (x[ds][idx+allele] == -1) and (tilestep != (prev_tilestep+1)):
        valid_tilevec[allele][ds] = False

  # max sure to take the maximum needed base for the current position
  #
  for ds in range(N_DS):
    for allele in range(0,2):
      if x[ds][idx+allele] == -1:
        if curtilevarid[allele][ds] >= basevec[idx+allele]:
          basevec[idx+allele] = curtilevarid[allele][ds]

  # force interleaved basevec to be the same
  #
  basevec[idx] = max(basevec[idx], basevec[idx+1])
  basevec[idx+1] = max(basevec[idx], basevec[idx+1])

  onehotvecsize[idx] = basevec[idx]
  onehotvecsize[idx+1] = basevec[idx+1]

  # Now that we have the basevec, we can calculate the onehotvecsize
  # for our current tile step position.
  # The only condition we have to concern ourselves here
  # is when it's a spanning tile.  The above default
  # population will take care of most cases when we have
  # either no-calls or high quality non-spanning tiles.
  # When we have a high-quality spanning tile, we need
  # to adjust the length to number of spanning tiles deep
  # it's in multiplied by the 'base'.
  #
  for ds in range(N_DS):
    for allele in range(0,2):

      #DEBUG
      if (ds == 383) and (idx == 9340):
        print ">>>> basevec:", basevec[idx], basevec[idx+1]

      # we can skip the invalid tiles
      #
      if not valid_tilevec[allele][ds]:
        continue

      #DEBUG
      if (ds == 383) and (idx == 9340):
        print ">>>> cp0"

      # non-anchor spanning tile
      #
      if x[ds][idx+allele] == -1:

        #DEBUG
        if (ds == 383) and (idx == 9340):
          print ">>>> cp1", "ds:", ds, "idx:", idx, "allele:", allele, "x..:", x[ds][idx+allele]

        # if our current tile is not a nocall, use it to calculate the 1hot vector size
        #
        if (curtilevarid[allele][ds] >= 0):
          #dstep = int( (idx + allele - anchorpos[allele][ds])/2 )
          anch_idx = anchorpos[allele][ds]

          #DEBUG
          if (ds == 383) and (idx == 9340):
            print ">>>> cp2", "anch_idx:", anch_idx, "tilestep:", tilestep, "ds:", ds, "idx:", idx, "allele:", allele, "x..:", x[ds][idx+allele]

          dstep = int( (tilestep - tilevec_info[anch_idx])/2 )
          onehotvecsize[idx+allele] = max( onehotvecsize[idx+allele], (basevec[idx+allele])*(dstep+1) )

  # force interleaved 1hot lengths to be the same
  #
  onehotvecsize[idx] = max(onehotvecsize[idx], onehotvecsize[idx+1])
  onehotvecsize[idx+1] = max(onehotvecsize[idx], onehotvecsize[idx+1])

  # DEBUG
  #zz = ""
  #for ds in range(N_DS):
  #  zz += " (" + str(x[ds][idx]) + "," + str(x[ds][idx+1]) + ")"
  #print "1hotlen:", idx, idx+1, onehotvecsize[idx], onehotvecsize[idx+1], "(", basevec[idx], basevec[idx+1], ")",  zz

  # onehotvecsize now holds the proper length for the one hot encoding.  Use it
  # to populate the final dataset
  #

  for ds in range(N_DS):
    for allele in range(0,2):
      if x[ds][idx+allele] >= 0:

        oha = one_hot(x[ds][idx+allele], onehotvecsize[idx+allele])

        if len(oha)==0:
          print "BAD?, ds:", ds, "idx:", idx, "allele:", allele, "x...:", x[ds][idx+allele], "sz:", onehotvecsize[idx+allele]
          quit()

        dataset[ds] = np.append(dataset[ds], oha)

        anchorpos[allele][ds] = idx + allele
        curtilevarid[allele][ds] = x[ds][idx+allele]

        #print " #Oha", "ds:", ds, idx+allele, json.dumps(oha.tolist())

      elif x[ds][idx+allele] == -2:

        anchorpos[allele][ds] = idx+allele
        curtilevarid[allele][ds] = -2
        for ii in range(int(onehotvecsize[idx+allele])):
          dataset[ds] = np.append(dataset[ds], np.nan)

        #print " #noc", "ds:", ds, idx+allele, " 0"*(onehotvecsize[idx+allele])

      else: # x[ds][idx] == -1

        if curtilevarid[allele][ds] == -2:
          for ii in range(int(onehotvecsize[idx+allele])):
            dataset[ds] = np.append(dataset[ds], 0)

          #print " #noa", "ds:", ds, idx+allele, " 0"*(onehotvecsize[idx+allele])

        # an invalid non-anchor spanning tile, fill in with 0s
        #
        elif not valid_tilevec[allele][ds]:
          for ii in range(int(onehotvecsize[idx+allele])):
            dataset[ds] = np.append(dataset[ds], 0)

        else:

          #print ">>> dstep", dstep, "ds", ds, "idx", idx, "a:", allele, "anchor", anchorpos[allele][ds]

          dstep = int( (idx + allele - anchorpos[allele][ds])/2 )
          oha = one_hot( (basevec[idx+allele])*dstep + curtilevarid[allele][ds], onehotvecsize[idx+allele])
          dataset[ds] = np.append(dataset[ds], oha)

          #print " #oha", "ds:", ds, idx+allele, json.dumps(oha.tolist())

  for allele in range(0,2):
    n = onehotvecsize[idx+allele]
    b = basevec[idx+allele]
    q = float(int(n/b))
    for ii in range(int(onehotvecsize[idx+allele])):
      rem = float(int(ii/b))
      #info_vec = np.append(info_vec, float(idx) + float(allele) + (rem/q))
      info_vec = np.append(info_vec, float(tilestep) + (rem/q))

      #print "  ## info_vec", idx+allele, ii, n, b, q, rem, float(idx) + float(allele) + (rem/q)


  prev_tilestep = tilestep

#print dataset
#for ds in range(len(dataset)):
#  print ds, dataset[ds].shape

#for ii in range(len(dataset)):
#  print ii, dataset[ii].shape

y = np.ndarray( (len(dataset), dataset[0].shape[0]) )

for ii in range(len(dataset)):
  for jj in range(dataset[0].shape[0]):
    y[ii][jj] = dataset[ii][jj]

np.save(ofn + "-info", info_vec)
np.save(ofn, y)

