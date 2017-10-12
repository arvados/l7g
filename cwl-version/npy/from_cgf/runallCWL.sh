#!/bin/bash

mkdir -p hiq

# filter out hiq tiles
#
./cnvrt-npy-to-hiqCWL.sh

# collect into single vectors, put in 'hiq/' dir with pfx 'hiq'
# input dir 'data-vec-hiq'
#
./collect-hiq-tilepaths.py npy-hiq hiq/hiq

# create simple 1hot arrays
#
./cnvrt-hiq-npy-to-1hot.py hiq/hiq-collect.npy hiq/hiq-1hot-simple.npy

