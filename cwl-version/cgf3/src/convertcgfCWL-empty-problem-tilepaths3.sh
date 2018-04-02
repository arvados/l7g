#!/bin/bash
#
# convert FastJ to CGF using SGLF library
#

eee() {
  echo "EEE: $@ $? $! $1 $2 $3"
}

trap eee ERR

set -eo pipefail

odir="data"

fjdir=$1
cgft=$2
fjt=$3
cglf=$4
start_tilepath_dec=$5
n_tilepath=$6

echo "fjdir: $fjdir"
echo "cgft: $cgft"
echo "fjt: $fjt"
echo "cglf: $cglf"

## problem tilepaths
#01a8 424
#0241 577
#0265 613
#0285 645
#02ac 684
#030f 783
#031f 799
#035d 861
#0024 36

force_empty_tilepaths_hx="01a8 0241 0265 0285 02ac 030f 031f 035d 0024 0000 0005 0015 0020 0025 0034 0037 003e 0055 0056 0057 0058 005c 0060 0098 00ba 00bb 00bd 00c6 00c7 00e9 00f7 00fd 0116 011e 0129 0133 0136 0146 0157 0159 015a 015c 0162 016f 0175 017f 019a 01a7 01ac 01c5 01d5 01d6 01ef 01f1 01fd 0208 0216 0220 0260 0266 0280 0286 0288 0289 028a 029c 02a1 02a5 02a6 02a7 02a9 02ba 02be 02c0 02c2 02c5 02d1 02d6 02d8 02e6 02e7 02ed 0302 0312 0320 032b 0337 0338 033b 0343 0349 0352 0354 0355 0357 035b 035c" 

if [[ "$start_tilepath_dec" == "" ]] ; then
  start_tilepath_dec=0
fi

if [[ "$n_tilepath" == "" ]] ; then
  n_tilepath=`expr 863 - $start_tilepath_dec`
fi

end_tilepath_dec_inc=`expr $start_tilepath_dec + $n_tilepath - 1`

if [[ "$fjdir" == "" ]] ; then
  exit 1
fi

function list_tilepaths {
  d="$1"
  for fjgz in `ls $d/*.fj.gz | sort` ; do
    tilepath=`basename $fjgz .fj.gz`
    echo $tilepath
  done
}

id=`basename "$fjdir"`
cgf_fn="$id.cgf"

mkdir -p $odir
rm -f $odir/$cgf_fn

echo ">>>> processing $fjdir, creating $odir/$cgf_fn"

ifn="$odir/$cgf_fn"
ofn="$odir/$cgf_fn"

$cgft -C $odir/$cgf_fn
echo header created

for fjgz in `ls $fjdir/*.fj.gz` ; do

  tilepath=`basename $fjgz .fj.gz`
  echo $tilepath

  if [[ `egrep $tilepath <( echo $force_empty_tilepaths_hx )` ]] ; then
    echo "# skipping $tilepath (in $force_empty_tilepaths_hx)"
    continue
  fi

  dec_tilepath=`cat <( echo "ibase=16; " )  <( echo "$tilepath" | tr '[:lower:]' '[:upper:]' ) | bc `

  $fjt -B -L <( zcat $cglf/$tilepath.sglf.gz ) -i <( zcat $fjdir/$tilepath.fj.gz ) | \
    $cgft -e $dec_tilepath $odir/$cgf_fn

  echo path $tilepath appended
done


#DEBUG
#start_tilepath_dec=32
#n_tilepath=5
#end_tilepath_dec_inc=`expr $start_tilepath_dec + $n_tilepath - 1`
#DEBUG


for tilepath in `comm -13 <( list_tilepaths $fjdir ) <( seq $start_tilepath_dec $end_tilepath_dec_inc | xargs -n1 -I{} printf "%04x\n" {} )` ; do

  if [[ `egrep $tilepath <( echo $force_empty_tilepaths_hx )` ]] ; then
    echo "# skipping $tilepath (in $force_empty_tilepaths_hx)"
    continue
  fi

  dec_tilepath=`cat <( echo "ibase=16; " )  <( echo "$tilepath" | tr '[:lower:]' '[:upper:]' ) | bc `

  ## If the FastJ file doesn't exist, use the empty tile from the SGLF and construct
  ## the band information ourselves.
  ##
  empty_tilevar=`zcat $cglf/$tilepath.sglf.gz | awk '{printf("%d %s\n", NR-1, $0)}' | cut -f1 -d' ' | tail -n1`

  echo "# stuffing empty_tilevar $empty_tilevar into tilepath $tilepath"

  echo -e "[ $empty_tilevar]\n[ $empty_tilevar]\n[[ ]]\n[[ ]]" | $cgft -e $dec_tilepath $odir/$cgf_fn

done

# finally, force an empty tilepaths for the "problem" tilepaths
# that we skipped above
#
for tilepath in $force_empty_tilepaths_hx ; do

  echo "# creating empty tilepath $tilepath"

  dec_tilepath=`cat <( echo "ibase=16; " )  <( echo "$tilepath" | tr '[:lower:]' '[:upper:]' ) | bc `
  echo -e "[ ]\n[ ]\n[[ ]]\n[[ ]]" | $cgft -e $dec_tilepath $odir/$cgf_fn
done
