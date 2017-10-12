#!/bin/bash
#
cgb="./cgb"
odir="./data"

mkdir -p tmp
mkdir -p $odir


#for tilepath in {862..862} ; do
for tilepath in {50,862} ; do

  h=`printf '%03x' $tilepath`

  rm -f tmp/$h.tmp

  for huid in `cat list` ; do

    ifn="cgf/$huid.cgf"

    echo $huid $ifn

    $cgb -i $ifn -k -p $tilepath -s 0 -L >> tmp/$h.tmp

  done

  cat tmp/$h.tmp | LD_LIBRARY_PATH=`pwd`/lib/cnpy ./band-to-matrix-npy $h $odir/$h
  ./cnvrt-npy-to-1hot.py $odir/$h $odir/$h-1hot
  mv $odir/$h-1hot.npy $odir/$h-1hot
  mv $odir/$h-1hot-info.npy $odir/$h-1hot-info

  ./cnvrt-npy-to-1hot-spanhot.py $odir/$h $odir/$h-1hot-span
  mv $odir/$h-1hot-span.npy $odir/$h-1hot-span
  mv $odir/$h-1hot-span-info.npy $odir/$h-1hot-span-info


done

./collect-tilepaths.py $odir $odir/data

pushd $odir
for ds in `ls data-collect*` ; do
  bn=`basename $ds .npy`
  mv $ds $bn
done
popd

./create-list list $odir/names.npy
mv $odir/names.npy $odir/names


pushd $odir
rm -f l7g-tile.npz
#zip l7g-tile.npz [0-9]* *.npy
zip l7g-tile.npz [a-f0-9][a-f0-9][a-f0-9]* names data-collect*
popd

