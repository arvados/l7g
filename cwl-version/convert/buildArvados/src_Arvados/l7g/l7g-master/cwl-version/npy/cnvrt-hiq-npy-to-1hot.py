#!/usr/bin/python
#
# Converts from integer interleaved tilepath to double one-hot with nan.
# The new 1-hot numpy arrays don't have sub arrays in them.
#

import re
import os
import sys
import numpy as np

def one_hot(v,n):
  oh = np.zeros(n+1)
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

ifn_info = re.sub("\.npy$", "-info.npy", ifn)


x = np.load(ifn)
x_info = np.load(ifn_info)

x_info_tilepath = []
x_info_tilepos = []

for ii in range(x_info.shape[0]):
  t = int(x_info[ii])
  tpos = t & 0xffffff
  tpath = (t >> (4*6)) & 0xffff

  #print hex(t), hex(tpath), hex(tpos), "(", tpath, float(tpos)/2.0, ")"

dataset = []
for ds in range(x.shape[0]):
  dataset.append( np.array([]) )

info_vec = np.array([])

# shape[1] holds interleaved tile position (x.shape[1] = 2*(# tile positions))
#
for apos in range(0, x.shape[1], 2):

  # shape[0] holds dataset (for each person, say).
  #
  maxstepsize = 0.0

  # collect max step size for the current (interleaved) tile position
  #
  for ds in range(x.shape[0]):
    maxstepsize = max(maxstepsize, x[ds][apos])
    maxstepsize = max(maxstepsize, x[ds][apos+1])

  for ii in range(int(maxstepsize+1)):
    #info_vec = np.append(info_vec, float(apos)/2.0)
    info_vec = np.append(info_vec, x_info[apos])

  for ii in range(int(maxstepsize+1)):
    #info_vec = np.append(info_vec, (float(apos)/2.0) + 0.5)
    info_vec = np.append(info_vec, x_info[apos+1])

  # Now we can construct the 1-hot vector for the position
  # (and both alleles).
  #
  for ds in range(x.shape[0]):

    # first allele
    #
    if x[ds][apos] == -2:
      for ii in range(maxstepsize+1):
        dataset[ds] = np.append(dataset[ds], np.nan)

    elif x[ds][apos] == -1:
      dataset[ds] = np.append(dataset[ds], np.zeros(maxstepsize+1))
    else:
      oha = one_hot(x[ds][apos], maxstepsize)
      dataset[ds] = np.append(dataset[ds], oha)

    # second allele
    #
    if x[ds][apos+1] == -2:
      for ii in range(maxstepsize+1):
        dataset[ds] = np.append(dataset[ds], np.nan)

    elif x[ds][apos+1] == -1:
      dataset[ds] = np.append(dataset[ds], np.zeros(maxstepsize+1))
    else:
      oha = one_hot(x[ds][apos+1], maxstepsize)
      dataset[ds] = np.append(dataset[ds], oha)

#print dataset
#for ds in range(len(dataset)):
#  print ds, dataset[ds].shape

y = np.ndarray( (len(dataset), dataset[0].shape[0]) )

for ii in range(len(dataset)):
  for jj in range(dataset[0].shape[0]):
    y[ii][jj] = dataset[ii][jj]

np.save(ofn + "-info", info_vec)
np.save(ofn, y)

