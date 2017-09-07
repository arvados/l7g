#!/bin/bash
#
# Download (via arv-get, in parallele, through ./pget.sh)
# all the relevant fastj files into the data directory.
# Once collected, run fastj2cgflib to create the library.
# output to $dstdir/PATH.sglf.gz
#
# example usage:
#
# ./run_cglf_single.sh 00a1
#

#set -eo pipefail

export tilepath="$1"
export dstdir="$2"
export fastj2cgflib="$3"
export datadir="$4"
export verbose_tagset="$5"
export tagset="$6"

#export dstdir="lib"
#export fastj2cgflib="/data-sdd/cwl_tiling/tilelib/fastj2cgflib"
#export datadir="/data-sdd/cwl_tiling/tilelib/data"
#export verbose_tagset="/data-sdd/cwl_tiling/tilelib/verbose_tagset"
#export tagset="/data-sdd/data/l7g/tagset.fa/tagset.fa.gz"

if [ "$tilepath" == ""  ] ; then
  echo provide tilepath
  exit 1
fi

mkdir -p $dstdir

#./pget.sh $tilepath

$fastj2cgflib -V -t <( $verbose_tagset $tilepath $tagset ) -f <( find $datadir -name "$tilepath.fj.gz" | xargs zcat ) | egrep -v '^#' | cut -f2- -d, | sort | bgzip -c > $dstdir/$tilepath.sglf.gz

#pushd data
#rm -f */$tilepath.fj.gz
#popd
