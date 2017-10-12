#!/bin/bash
#

#idir="./npy"
#odir="./npy-hiq"

export idir=$1
export odir=$2
export npyhiq=$3

mkdir -p $odir

export LD_LIBRARY_PATH=`pwd`/lib/cnpy

for tilepath in {0..862}; do
  h=`printf '%03x' $tilepath`
  echo $tilepath $h

  $npyhiq $tilepath $idir/$h $odir/$h
done

