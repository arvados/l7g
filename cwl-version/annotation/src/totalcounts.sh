#!/bin/bash

set -e
set -o pipefail

counts=( "$@" )

export allrsid="0"
export allgnomad="0"
export alltotal="0"

cat ${counts[@]}

for count in ${counts[@]}; do
  rsid=`cut -d' ' -f5 $count`
  gnomad=`cut -d' ' -f10 $count`
  total=`cut -d' ' -f2 $count`
  allrsid=`echo $(($allrsid + $rsid))`
  allgnomad=`echo $(($allgnomad + $gnomad))`
  alltotal=`echo $(($alltotal + $total))`
done
rsidpercentage=`awk -v n="$allrsid" -v d="$alltotal" 'BEGIN {print n/d*100}'`
gnomadpercentage=`awk -v n="$allgnomad" -v d="$alltotal" 'BEGIN {print n/d*100}'`

echo "overall: $alltotal total variants, $allrsid variants ($rsidpercentage%) have rsID, $allgnomad variants ($gnomadpercentage%) have gnomad AF"
