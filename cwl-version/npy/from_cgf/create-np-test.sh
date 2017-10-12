#!/bin/bash
#
# Make numpy vectors on a tile path basis
# for the pgp set (

export cgb="/data-sdd/cwl_tiling/npy/cgb"
export odir="./testnpy-filelist"
export indir="/data-sdd/data/cgf-cgi-1kg-69"

mkdir -p tmp
mkdir -p $odir

./create-list $(find $indir/*.cgf) $odir/names.npy
mv $odir/names.npy $odir/names

