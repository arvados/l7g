#!/bin/bash

n=862

for i in {0..862} ; do
  tilepath=`printf '%04x' $i`

  echo $tilepath



  continue

  fastj2cgflib -V \
    -t -<( ./verbose_tagset $tilepath )  \
    -f <( find ./data -name $tilepath.fj.gz | xargs zcat )
done
