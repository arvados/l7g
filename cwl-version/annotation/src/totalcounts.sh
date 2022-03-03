#!/bin/bash

set -e
set -o pipefail

counts=( "$@" )

export allgnomad="0"
export alltotal="0"

cat ${counts[@]}

for count in ${counts[@]}; do
  gnomad=`cut -d' ' -f2 $count`
  total=`cut -d' ' -f5 $count`
  allgnomad=`echo $(($allgnomad + $gnomad))`
  alltotal=`echo $(($alltotal + $total))`
done
percentage=`awk -v n="$allgnomad" -v d="$alltotal" 'BEGIN {print n/d*100}'`

echo "overall: $allgnomad out of $alltotal variants ($percentage%) have gnomad AF"
