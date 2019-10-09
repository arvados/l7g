#!/bin/bash
#
#

export tilepathmin="$1"
export tilepathmax="$2"
export fjcsv2sglf="$3"
export datadir="$4"
export fjt="$5"
export tagset="$6"

export dstdir="."


for tilepath_dec in `seq $tilepathmin $tilepathmax` ; 

do

  tilepath=$(printf "%04x" $tilepath_dec) 
  echo $tilepath_dec
  echo $tilepath 
  ## if we don't have any FastJ for thie tilepath, create an empty tile
  ## to be inserted into the SGLF
  ##
  fj_count=`find $datadir -name "$tilepath.fj.gz" -type f | wc -l`

  if [[ "$fj_count" -eq "0" ]] ; then

    empty_hash=`echo -e -n '' | md5sum | cut -f1 -d' '`
    echo "$tilepath.00.0000.000+1,$empty_hash," | bgzip -c > $dstdir/$tilepath.sglf.gz

  else

  find $datadir -name "$tilepath.fj.gz" -type f | \
    xargs zcat | \
    $fjt -C -U | \
    $fjcsv2sglf <( samtools faidx $tagset $tilepath.00 | egrep -v '^>' | tr -d '\n' | fold -w 24 ) | \
    bgzip -c > $dstdir/$tilepath.sglf.gz

  fi
done
