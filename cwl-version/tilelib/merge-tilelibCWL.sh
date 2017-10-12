#!/bin/bash

set -e -o pipefail

export srcdir="$1"
export nppdir="$2"
export nthreads="$3"
export mergetilelib="$4"

#export srcdir="/data-sdd/data/sglf"
#export nppdir="/data-sdd/scripts/tilelib/lib-cgi-69"
#export dstdir="/data-sdd/cwl_tiling/tilelib/lib-merge-test"
#export nthreads="16"
#export mergetilelib="/data-sdd/cwl_tiling/tilelib/merge-tilelib"

export SHELL=/bin/bash 

if [[ "$srcdir" == "" ]] ; then
  echo "provide srcdir"
  exit
fi

if [[ "$nppdir" == "" ]] ; then
  echo "provide add dir"
  exit
fi

#if [[ "$dstdir" == "" ]] ; then
#  export dstdir="lib.merge"
#  echo "using $dstdir destination dir"
#fi


export dstdir="lib-merge"
mkdir -p $dstdir

function process {
  tpath=$1
  tfn="$tpath.sglf.gz"

  echo "processing $tpath"

  srcfn="$srcdir/$tfn"
  nppfn="$nppdir/$tfn"
  dstfn="$dstdir/$tfn"

  echo ">>> $srcfn $nppfn $dstfn"

  if [[ -e $nppfn ]] ; then

    $mergetilelib <( zcat $srcfn ) <( zcat $nppfn ) | bgzip -c > $dstfn

  else

    echo "# WARNING: $nppfn does not exist, copying $srcfn to $dstfn"

    cp $srcfn $dstfn

  fi

  echo "  $tpath done"
}
export -f process

#for tfn in `ls $srcdir | head -n 20` ; do
for tfn in `ls $srcdir` ; do
  #echo $tfn

  tpath=`basename $tfn .sglf.gz`
  srcfn="$srcdir/$tfn"
  nppfn="$nppdir/$tfn"
  dstfn="$dstdir/$tfn"

  echo $tpath

done | parallel --max-procs $nthreads process {}

