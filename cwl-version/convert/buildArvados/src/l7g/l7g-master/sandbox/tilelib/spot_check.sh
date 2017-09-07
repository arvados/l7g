#!/bin/bash
#
# Check to make sure last step is in hg19.
# This is a rudamentary check to make
# sure the sglf files didn't get truncated.
#

ref_fj_dir="/data-sde/data/fastj/hg19"

for x in `ls lib/*.gz` ; do
  p=`basename $x .sglf.gz`
  #echo $p $x
  sglf_last_tilestep=`zcat $x | tail -n1 | cut -f1 -d',' | cut -f1-3 -d'.'`
  reffj_last_tilestep=`zgrep '^>' $ref_fj_dir/$p.fj.gz | tail -n1 | sed 's/^>//' | jq -r -c '.tileID' | cut -f1-3 -d'.'`

  if [[ "$sglf_last_tilestep" != "$reffj_last_tilestep" ]] ; then
    echo "MISMATCH $p $sglf_last_tilestep != $reffj_last_tilestep"
  else
    echo "ok $p"
  fi

done
