#!/bin/bash
#

idir="./data-vec-pgp"
odir="./data-vec-pgp-hiq"

mkdir -p $odir

export LD_LIBRARY_PATH=`pwd`/lib/cnpy

for tilepath in {0..862}; do
  h=`printf '%03x' $tilepath`
  echo $tilepath $h

  ./npy-to-hiq $tilepath $idir/$h $odir/$h
done

