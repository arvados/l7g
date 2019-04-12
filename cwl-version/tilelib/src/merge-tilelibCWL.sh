#!/bin/bash

set -e -o pipefail

while getopts "s:n:" option; do
case "$option" in
  s) export srcdir=$OPTARG;;
  n) export newdir=$OPTARG;;
esac
done

export nthreads=${@:$OPTIND:1}
export mergetilelib=${@:$OPTIND+1:1}

export SHELL=/bin/bash

if [[ "$newdir" == "" ]] ; then
  echo "provide add dir"
  exit
fi

export dstdir="lib-merge"
mkdir -p $dstdir

function process {
  tpath=$1
  tfn="$tpath.sglf.gz"

  echo "processing $tpath"

  srcfn="$srcdir/$tfn"
  newfn="$newdir/$tfn"
  dstfn="$dstdir/$tfn"

  echo ">>> $srcfn $newfn $dstfn"

  if [[ -e $newfn ]] ; then

    $mergetilelib <( zcat $srcfn ) <( zcat $newfn ) | bgzip -c > $dstfn

  else

    echo "# WARNING: $newfn does not exist, copying $srcfn to $dstfn"

    cp $srcfn $dstfn

  fi

  echo "  $tpath done"
}
export -f process

if [[ "$srcdir" == "" ]] ; then
  cp $newdir/* $dstdir
else
  for tfn in `ls $newdir` ; do

    tpath=`basename $tfn .sglf.gz`
    srcfn="$srcdir/$tfn"
    newfn="$newdir/$tfn"
    dstfn="$dstdir/$tfn"

    echo $tpath

  done | parallel --max-procs $nthreads process {}
fi
