#!/bin/bash
#
# Use local copy of FastJ files.  Run on the 69 CGI 1kg
# FastJ files created from the GFF files.
# example usage:
#
# ./run_cglf_single.sh 00a1
#

set -eo pipefail

tilepath="$1"

dstdir="lib-cgi-69"

if [ "$tilepath" == ""  ] ; then
  echo invalid tilepath
  exit 1
fi

mkdir -p $dstdir

#fastj2cgflib -V -t <( ./verbose_tagset $tilepath ) -f <( find ./data -name "$tilepath.fj.gz" | xargs zcat ) | egrep -v '^#' | cut -f2- -d, | sort | pbgzip -c > $dstdir/$tilepath.sglf.gz
fastj2cgflib -V -t <( ./verbose_tagset $tilepath ) -f <( find -L ./cgi-69-data -name "$tilepath.fj.gz" | xargs zcat ) | egrep -v '^#' | cut -f2- -d, | sort | pbgzip -c > $dstdir/$tilepath.sglf.gz

