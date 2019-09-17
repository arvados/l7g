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
skippaths=$5

echo "fjdir: $fjdir"
echo "cgft: $cgft"
echo "fjt: $fjt"
echo "cglf: $cglf"

start_tilepath_dec=0
n_tilepath=`expr 863 - $start_tilepath_dec`
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

  if [[ `egrep $tilepath $skippaths` ]] ; then
    echo "# skipping $tilepath"
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

  if [[ `egrep $tilepath $skippaths` ]] ; then
    echo "# skipping $tilepath"
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
for tilepath in `cat $skippaths` ; do

  echo "# creating empty tilepath $tilepath"

  dec_tilepath=`cat <( echo "ibase=16; " )  <( echo "$tilepath" | tr '[:lower:]' '[:upper:]' ) | bc `
  echo -e "[ ]\n[ ]\n[[ ]]\n[[ ]]" | $cgft -e $dec_tilepath $odir/$cgf_fn
done
