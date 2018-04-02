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
export nthreads=$6

export odir="./npy"
export odir2="./npy-hiq"
export SHELL=/bin/bash

mkdir -p tmp
mkdir -p $odir
mkdir -p $odir2

find $indir/*.cgf > cgf_list

function process_tile {
  tilepath=$1

  printf '%03x\n' $tilepath
  h=`printf '%03x' $tilepath`

  rm -f tmp/$h.tmp

  for ifn in `cat cgf_list` ; do

    #echo $ifn
    #$cgb -i $ifn -k -p $tilepath -s 0 -L >> tmp/$h.tmp
    $cgft -b $tilepath $ifn >> tmp/$h.tmp

  done

  #cat tmp/$h.tmp | LD_LIBRARY_PATH=`pwd`/lib/cnpy $band2matrix $h $odir/$h
  cat tmp/$h.tmp | $band2matrix $h $odir/$h

}
export -f process_tile

for tp in {0..862} ; do
  echo $tp
done | parallel --max-procs $nthreads process_tile {}

$makelist cgf_list $odir/names.npy
mv $odir/names.npy $odir/names

# Filtering out hiq and writing tile variant matrix and 1-hot matrix

$cnvrt2hiq $odir/names $odir $odir2
