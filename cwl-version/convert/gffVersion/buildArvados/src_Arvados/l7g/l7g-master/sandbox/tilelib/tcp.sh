#!/bin/bash

bdir="/data-sde/scripts/convert/stage.okg/"
tilepath="0003"

for x in `ls $bdir ` ; do
  echo $x

  mkdir -p data/$x
  cp $bdir/$x/$tilepath.fj.gz data/$x
done
