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
cgf=$2
cglf=$3


if [[ "$fjdir" == "" ]] ; then
  echo "provide fjdir"
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

$cgf -action header -i nop -o $odir/$cgf_fn
echo header created

for fjgz in `ls $fjdir/*.fj.gz` ; do

  tilepath=`basename $fjgz .fj.gz`
  # echo $tilepath

  $cgf -action append -i <( zcat $fjdir/$tilepath.fj.gz ) -path $tilepath -S <( zcat $cglf/$tilepath.sglf.gz ) -cgf $ifn -o $ofn
  echo path $tilepath appended
done
