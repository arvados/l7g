#!/bin/bash
#
# Make numpy vectors on a tile path basis
# for the pgp set (

export cgb="./cgb"
export odir="./data-vec-pgp"

mkdir -p tmp
mkdir -p $odir

function process_tile {
  tilepath=$1

  h=`printf '%03x' $tilepath`

  rm -f tmp/$h.tmp

  for huid in `cat list-pgp` ; do

    ifn="cgf/$huid.cgf"

    echo $huid $ifn

    $cgb -i $ifn -k -p $tilepath -s 0 -L >> tmp/$h.tmp

  done

  cat tmp/$h.tmp | LD_LIBRARY_PATH=`pwd`/lib/cnpy ./band-to-matrix-npy $h $odir/$h

}
export -f process_tile

for tp in {0..862} ; do
  echo $tp
done | parallel --max-procs 10 process_tile {}

./create-list list-pgp $odir/names.npy
mv $odir/names.npy $odir/names

