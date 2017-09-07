#!/bin/bash

band_dir="odata.band"
icgf="test.cgf3"

for p in {0..862}; do
  hxp=`printf '%04x' $p`

  cgf3_m5=`cgft -b $p $icgf | md5sum | cut -f1 -d' '`
  orig_m5=`md5sum $band_dir/$hxp.band | cut -f1 -d ' '`

  echo $p $hxp, $cgf3_m5 $orig_m5

  if [[ "$cgf3_m5" != "$orig_m5" ]] ; then
    echo "MISMATCH" $p $hxp "$orig_m5" "!=" "$cgf3_m5"
    exit
  fi

done

echo "ok"
