#!/bin/bash
#
# copy over a single paths worth of fastj to process
#

tilepath="0003"

basedir="/data-sde/keep/home/FastJ"
odir="data"

for x in `ls $basedir` ; do
  echo $x

  mkdir -p $odir/$x
  cp /data-sde/keep/home/FastJ/$x/$tilepath.fj.gz $odir/$x
done
