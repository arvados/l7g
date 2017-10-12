#!/bin/bash
#
# download FastJ, convert to CGF, delete
#

set -eo pipefail

csv="./data/name-pdh.csv"
stage="./stage"

hid="$1"
if [[ "$hid" == "" ]] ; then
  echo "provide hid"
  exit 1
fi

mkdir -p $stage

name=`grep $hid $csv | csvtool col 1 -`
pdh=`grep $hid $csv | csvtool col 2 -`

echo $name $pdh

mkdir -p $stage/$name
#arv-get $pdh/ $stage/$name/
./convert-fastj-to-cgf.sh $stage/$name
rm -rf $stage/$name

echo "processed $name"
