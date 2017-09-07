#!/bin/bash

export srcdir="$1"
export nppdir="$2"
export dstdir="$3"

if [[ "$srcdir" == "" ]] ; then
  echo "provide srcdir"
  exit
fi

if [[ "$nppdir" == "" ]] ; then
  echo "provide add dir"
  exit
fi

if [[ "$dstdir" == "" ]] ; then
  export dstdir="lib.merge"
  echo "using $dstdir destination dir"
fi

mkdir -p $dstdir

function process {
  tpath=$1
  tfn="$tpath.sglf.gz"

  echo "processing $tpath"

  srcfn="$srcdir/$tfn"
  nppfn="$nppdir/$tfn"
  dstfn="$dstdir/$tfn"

  echo ">>> $srcfn $nppfn $dstfn"

  ./merge-tilelib <( zcat $srcfn ) <( zcat $nppfn ) | bgzip -c > $dstfn

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

done | parallel --max-procs 12 process {}

