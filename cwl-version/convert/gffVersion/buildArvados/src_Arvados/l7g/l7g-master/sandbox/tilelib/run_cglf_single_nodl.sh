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

set -eo pipefail

tilepath="$1"

dstdir="lib"


if [ "$tilepath" == ""  ] ; then
  echo provide tilepath
  exit 1
fi

fastj2cgflib -V -t <( ./verbose_tagset $tilepath ) -f <( find ./data -name "$tilepath.fj.gz" | xargs zcat ) | egrep -v '^#' | cut -f2- -d, | sort | pbgzip -c > $dstdir/$tilepath.sglf.gz
