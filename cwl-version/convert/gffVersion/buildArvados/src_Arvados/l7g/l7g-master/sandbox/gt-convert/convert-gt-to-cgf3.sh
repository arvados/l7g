#!/bin/bash

set -eo pipefail
mkdir -p .tmp

rm -rf odata odata.merge odata.band
mkdir -p odata
mkdir -p odata.merge
mkdir -p odata.band

merg="/data-sdd/scripts/tilelib/merge-tilelib.sh"
tilelib="/data-sdd/data/sglf"

gtfn="$1"
ocgf="out.cgf3"

if [[ "$gtfn" == "" ]] ; then
  echo "provide genotyping file"
  exit
fi

if [[ "$2" != "" ]] ; then
  ocgf="$2"
fi

tfn=`tempfile -d .tmp`
tsglf=`tempfile -d .tmp`

echo "# creating intermediary tile-variant position ssv file..."
./src/gt-23andMeConvert.py $gtfn > $tfn

echo "# creating sglf file from ssv..."
#./src/gt-ssv-tile.py $tfn > $tsglf
./src/gt-ssv-tile-single.py $tfn 0 > $tsglf.0
./src/gt-ssv-tile-single.py $tfn 1 > $tsglf.1

## DOUBLE CHECK
##
./src/cap-tile-check.py $tsglf.0 > $tsglf.0.check
./src/cap-tile-check.py $tsglf.1 > $tsglf.1.check
pushd .tmp
../scripts/filt-clean `basename $tsglf.0.check`
../scripts/filt-clean `basename $tsglf.1.check`
popd

cat $tsglf.0 $tsglf.1 | sort > $tsglf

./src/create-noc-info.py $tsglf > $tsglf.info

echo "# splitting sglf..."
./src/split-gt-ssv-sglf.py $tsglf # into odata

echo "# cleaning up..."
./src/cleanup-sglf.sh

echo "# merging from tilelib" $tilelib
pushd src
./merge-tilelib.sh $tilelib ../odata ../tilelib.merge
popd

echo "# converting to band format..."
for itilepath in {0..862} ; do
  hxp=`printf "%04x" $itilepath`
  hxtilepath=`printf "%04x.sglf.gz" $itilepath`
  ofn_band=`basename $hxtilepath .sglf.gz`

  #echo $hxp

  if [[ -f "odata/$hxtilepath" ]] ; then
    ./src/gt-sglf-to-band.py \
      <( zcat odata/$hxtilepath | cut -f1,2 -d, ) \
      <( zcat tilelib.merge/$hxtilepath ) \
      <( egrep '^'$hxp'\.' $tsglf.info ) > odata.band/${ofn_band}.band
  else

    #echo "## using empty.file"

    ./src/gt-sglf-to-band.py \
      <( echo -n "" ) \
      <( zcat tilelib.merge/$hxtilepath ) \
      <( egrep '^'$hxp'\.' $tsglf.info ) > odata.band/${ofn_band}.band
  fi
done

rm -f "$ocgf"
cgft -C "$ocgf"

echo "# creating $ocgf"
for x in `ls odata.band`; do
  hxname=`basename $x .band`
  tilepath=`echo $((16#$hxname))`

  echo "# adding band $x ($hxname,$tilepath)"

  cgft -e $tilepath $ocgf < odata.band/$hxname.band
done

band_dir="odata.band"
icgf="$ocgf"

echo "# CHECKING cgf"
for p in {0..862}; do
  hxp=`printf '%04x' $p`

  cgf3_m5=`cgft -b $p $icgf | md5sum | cut -f1 -d' '`
  orig_m5=`md5sum $band_dir/$hxp.band | cut -f1 -d ' '`

  echo "#" $p $hxp, $cgf3_m5 $orig_m5

  if [[ "$cgf3_m5" != "$orig_m5" ]] ; then
    echo "MISMATCH" $p $hxp "$orig_m5" "!=" "$cgf3_m5"
    exit
  fi

done

echo "## FINISHED: $ocgf"

## DEBUG
echo "## NOT CLEANING UP $tfn $tsglf"
exit

rm -f $tfn
rm -f $tsglf
rm -f $tsglf.info
rm -f $tsglf.0
rm -f $tsglf.0.check
rm -f $tsglf.1
rm -f $tsglf.1.check
