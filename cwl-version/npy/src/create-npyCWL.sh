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
export makelist=$4
export nthreads=$5

export odir="./npy"
export SHELL=/bin/bash

export tmpdir=`mktemp -d`
mkdir -p $odir

find $indir/*.cgf > cgf_list

function process_tile {
  tilepath=$1

  printf '%03x\n' $tilepath
  h=`printf '%03x' $tilepath`

  rm -f $tmpdir/$h.tmp

  for ifn in `cat cgf_list` ; do

    #echo $ifn
    #$cgb -i $ifn -k -p $tilepath -s 0 -L >> tmp/$h.tmp
    $cgft -b $tilepath $ifn >> $tmpdir/$h.tmp

  done

  #cat tmp/$h.tmp | LD_LIBRARY_PATH=`pwd`/lib/cnpy $band2matrix $h $odir/$h
  cat $tmpdir/$h.tmp | $band2matrix $h $odir/$h

}
export -f process_tile

for tp in {0..862} ; do
  echo $tp
done | parallel --max-procs $nthreads process_tile {}

$makelist cgf_list $odir/names.npy
mv $odir/names.npy $odir/names
