#!/bin/bash
#
#

export tilepath="$1"
export fastj2cgflib="$2"
export datadir="$3"
export verbose_tagset="$4"
export tagset="$5"

export dstdir="lib"

if [ "$tilepath" == ""  ] ; then
  echo provide tilepath
  exit 1
fi

mkdir -p $dstdir

## if we don't have any FastJ for thie tilepath, create an empty tile
## to be inserted into the SGLF
##
fj_count=`find $datadir -name "$tilepath.fj.gz" -type f | wc -l`

if [[ "$fj_count" -eq "0" ]] ; then

  empty_hash=`echo -e -n '' | md5sum | cut -f1 -d' '`
  echo "$tilepath.00.0000.000+1,$empty_hash," | bgzip -c > $dstdir/$tilepath.sglf.gz

else


  $fastj2cgflib -V -t <( $verbose_tagset $tilepath $tagset ) -f <( find $datadir -name "$tilepath.fj.gz" | xargs zcat ) | \
    egrep -v '^#' | \
    cut -f2- -d, | \
    sort | \
    bgzip -c > $dstdir/$tilepath.sglf.gz

fi
