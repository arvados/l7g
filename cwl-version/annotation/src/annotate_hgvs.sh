#!/bin/bash

set -aeo pipefail

export get_hgvs="$1"

export path="$2"
export ref="$3"
export tilelib="$4"
export varnum="$5"
export assembly="$6"

export assemblyindex="${assembly%.*}.fwi"

export length=`grep -P ":$path\t" $assemblyindex | cut -f2`
export offset=`grep -P ":$path\t" $assemblyindex | cut -f3`

export steps=`bgzip -b $offset -s $length $assembly | cut -f1`

for step in ${steps[@]}; do
  echo "## annotating path $path step $step"
  $get_hgvs $path $step $ref $tilelib $varnum $assembly >> $path.csv
done
