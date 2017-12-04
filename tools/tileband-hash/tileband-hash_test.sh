#!/bin/bash

a=`./tileband-hash -L test-data/035e.1.sglf -T 862 <( cat test-data/hu826751-035e.band test-data/hu34D5B9-035e.band ) | md5sum | cut -f1 -d' '`
b=`cat test-data/expect-test0.txt | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: mismatch: expected $b got $a"
  exit -1
fi

echo "ok"
exit 0
