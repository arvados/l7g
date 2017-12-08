#!/bin/bash
#
# Make numpy vectors on a tile path basis
# for the pgp set (


#export cgb="/data-sdd/cwl_tiling/npy/cgb"
#export indir="/data-sdd/data/cgf-test"
#export band2matrix="/data-sdd/cwl_tiling/npy/band-to-matrix-npy"
#export makelist="/data-sdd/cwl_tiling/npy/create-list"
#export cpnylib="/data-sdd/cwl_tiling/npy/lib/cnpy"

#export cgb=$1
export cgft=$1
export indir=$2
export band2matrix=$3
export cnvrt2hiq=$4
export makelist=$5
export odir=$4

#export nthreads=$6

#export odir="./npy"
export odir2="./npy-hiq"
export SHELL=/bin/bash

mkdir -p tmp
#mkdir -p $odir
mkdir -p $odir2

find $indir/*.cgf > cgf_list

$makelist cgf_list $odir2/names.npy
mv $odir2/names.npy $odir2/names

# Filtering out hiq and writing tile variant matrix and 1-hot matrix

$cnvrt2hiq $odir2/names $odir $odir2
