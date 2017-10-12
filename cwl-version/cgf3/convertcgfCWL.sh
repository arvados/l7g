#!/bin/bash
#
# convert FastJ to CGF using SGLF library
#

eee() {
  echo "EEE: $@ $? $! $1 $2 $3"
}

trap eee ERR

set -eo pipefail

#cgf="/data-sdd/cwl_tiling/cgf/cgf"
#cglf="/data-sdd/data/sglf"
#fjdir="/data-sdd/cwl_tiling/tilelib/data/HG00731-200-37"
odir="data"

fjdir=$1
cgft=$2
fjt=$3
cglf=$4

echo "fjdir: $fjdir"
echo "cgft: $cgft"
echo "fjt: $fjt"
echo "cglft: $cglft"


if [[ "$fjdir" == "" ]] ; then
#  echo "provide fjdir"
  exit 1
fi

id=`basename "$fjdir"`

cgf_fn="$id.cgf"

mkdir -p $odir
rm -f $odir/$cgf_fn

#mkdir -p log

echo ">>>> processing $fjdir, creating $odir/$cgf_fn"

ifn="$odir/$cgf_fn"
ofn="$odir/$cgf_fn"

#$cgf -action header -i nop -o $odir/$cgf_fn 
$cgft -C $odir/$cgf_fn
echo header created

for fjgz in `ls $fjdir/*.fj.gz` ; do

  tilepath=`basename $fjgz .fj.gz`
  echo $tilepath

  dec_tilepath=`cat <( echo "ibase=16; " )  <( echo "$tilepath" | tr '[:lower:]' '[:upper:]' ) | bc `

  #$cgf -action append -i <( zcat $fjdir/$tilepath.fj.gz ) -path $tilepath -S <( zcat $cglf/$tilepath.sglf.gz ) -cgf $ifn -o $ofn > errorlog 
  $fjt -B -L <( zcat $cglf/$tilepath.sglf.gz ) -i <( zcat $fjdir/$tilepath.fj.gz ) | \
    $cgft -e $dec_tilepath $odir/$cgf_fn

#  tail -n 5 errorlog >&1
  echo path $tilepath appended
done

#tail errorlog >&1
