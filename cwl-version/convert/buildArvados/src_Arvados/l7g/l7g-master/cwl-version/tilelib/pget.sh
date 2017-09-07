#!/bin/bash
#
# we do a simple version of
# parallelism by waiting after every
# batch of PN arv-get invocations (currently
# set to 30).
#
# arv-get gives oom when not doing this wait.

set -e

tilepath="$1"

if [ "$tilepath" == "" ] ; then
  echo "provide tilepath"
  exit 1
fi

tvar=0
mkdir -p data

PN=30
pcount=0

while read line ; do
  if [ "$line" == "" ] ; then
    continue
  fi

  id=`echo "$line" | cut -f1`
  pdh=`echo "$line" | cut -f2`

  #echo $line '-->' $id','$pdh '('$tvar')'


  let tvar="$tvar + 1"
  let pcount="$pcount + 1"

  mkdir -p data/$id

  echo "data/$id/$tilepath.fj.gz"

  #echo "arv-get $line/$tilepath.fj.gz ./data/$id/$tilepath.fj.gz"
  arv-get $pdh/$tilepath.fj.gz ./data/$id/$tilepath.fj.gz &

  if [ "$pcount" -gt $PN ] ; then
    pcount=0
    wait
  fi

done < <( cat id_pdh.tsv | egrep -v '^#' )

wait

