#!/bin/bash
#
# b37 has some differences than hg19.  Lift it over
# to a new tile sassembly for use.
#

tagset="$1"
ref="$2"
hg19cyto="$3"
libver="00"

if [[ "$tagset" == "" ]] ||
   [[ "$ref" == "" ]] ||
   [[ $hg19cyto == "" ]] ; then
  echo "provide tagset, fasta reference and cytogenetic band files"
  exit 1
fi

if [[ "$libver" == "" ]] ; then
  libver="00"
fi

tilepath=0

refname=` basename "$ref" .fa.gz `
refname=` basename "$refname" .fasta.gz `

stagedir="./stage"
odir="./output"

mkdir -p $odir

mkdir -p $stagedir
while IFS='' read -r line  || [[ -n "$line" ]]; do
  chrom=`echo "$line" | cut -f1`
  pos0ref_start=`echo "$line" | cut -f2`
  pos0ref_end_noninc=`echo "$line" | cut -f3`

  altchrom=`echo "$chrom" | sed 's/^chr//' `

  if [[ "$chrom" == "chrM" ]] ; then continue ; fi

  pos1ref_start=`expr 1 + $pos0ref_start`
  pos1ref_end_inc=`expr 1 + $pos0ref_end_noninc - 1`



  hxp=`printf "%04x" $tilepath`

  echo $chrom $altchrom $pos0ref_start $pos0ref_end_noninc ... $pos1ref_start $pos1ref_end_inc "($tilepath $hxp)"

  tile-liftover -T <( refstream $tagset $hxp.$libver | tr -d '\n' ) \
    -p $tilepath \
    -R <( refstream $ref $altchrom:$pos1ref_start-$pos1ref_end_inc | tr -d '\n' ) \
    -s $pos0ref_start \
    -N $refname \
    -c $altchrom > $stagedir/$hxp.liftover

  tilepath=`expr 1 + $tilepath`

done < <( cat $hg19cyto | egrep -v '^#' )

echo "processing MT"
p=862
hxp="035e"
altchrom="MT"

tile-liftover -T <( refstream $tagset $hxp.$libver | tr -d '\n' ) \
  -p $tilepath \
  -R <( refstream $ref $altchrom | tr -d '\n' ) \
  -N $refname \
  -c $altchrom > $stagedir/$hxp.liftover

assembly_name="assembly.$libver.$refname"

for f in `find $stagedir -name '*.liftover' | sort`; do
  cat $f >> $odir/$assembly_name.fw
done

pushd $odir
cat $assembly_name.fw | tile-assembly-index - > $assembly_name.fw.fwi
bgzip -i $assembly_name.fw
ln -s $assembly_name.fw.fwi $assembly_name.fw.gz.fwi
popd

mv $odir/* .
rm -rf $stagedir
rm -rf $odir


