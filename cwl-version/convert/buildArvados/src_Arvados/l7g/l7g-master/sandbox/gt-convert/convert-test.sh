#!/bin/bash

set -eo pipefail
mkdir -p .tmp

fncgf="test.cgf3"

#rm -rf odata odata.merge
#mkdir -p odata
#mkdir -p odata.merge
#
#merg="/data-sdd/scripts/tilelib/merge-tilelib.sh"
#tilelib="/data-sdd/data/sglf"
#
#gtfn="$1"
#
#if [[ "$gtfn" == "" ]] ; then
#  echo "provide genotyping file"
#  exit
#fi
#
#tfn=`tempfile -d .tmp`
#tsglf=`tempfile -d .tmp`

#tfn='.tmp/fileP5cW0l'
#tsglf='.tmp/filejYHlNl'

#
#echo "creating intermediary tile-variant position ssv file..."
#./src/gt-23andMeConvert.py $gtfn > $tfn
#
#echo "creating sglf file from ssv..."
##./src/gt-ssv-tile.py $tfn > $tsglf
#./src/gt-ssv-tile-single.py $tfn 0 > $tsglf.0
#./src/gt-ssv-tile-single.py $tfn 1 > $tsglf.1
#
### DOUBLE CHECK
###
#./src/cap-tile-check.py $tsglf.0 > $tsglf.0.check
#./src/cap-tile-check.py $tsglf.1 > $tsglf.1.check
#pushd .tmp
#../scripts/filt-clean `basename $tsglf.0.check`
#../scripts/filt-clean `basename $tsglf.1.check`
#popd
#
#cat $tsglf.0 $tsglf.1 | sort > $tsglf
#
#./src/create-noc-info.py $tsglf > $tsglf.info
#
#echo "splitting sglf..."
#./src/split-gt-ssv-sglf.py $tsglf # into odata
#
#echo "cleaning up..."
#./src/cleanup-sglf.sh
#
#echo "merging from tilelib" $tilelib
#pushd src
#./merge-tilelib.sh $tilelib ../odata ../tilelib.merge
#popd

#echo "converting to band format..."
#for itilepath in {0..862} ; do
#  hxtilepath=`printf "%04x.sglf.gz" $itilepath`
#  ofn_band=`basename $hxtilepath .sglf.gz`
#  ./src/gt-sglf-to-band.py \
#    <( zcat odata/$hxtilepath | cut -f1,2 -d, ) \
#    <( zcat tilelib.merge/$hxtilepath ) \
#    $tsglf.info > odata.band/${ofn_band}.band
#done
#

rm -f $fncgf
cgft -C $fncgf

echo "creating $fncgf"
for x in `ls odata.band`; do
  hxname=`basename $x .band`
  tilepath=`echo $((16#$hxname))`

  echo $x $hxname $tilepath

  cgft -e $tilepath $fncgf < odata.band/$hxname.band
done

exit

rm -f $tfn
rm -f $tsglf
rm -f $tsglf.info
rm -f $tsglf.0
rm -f $tsglf.0.check
rm -f $tsglf.1
rm -f $tsglf.1.check
