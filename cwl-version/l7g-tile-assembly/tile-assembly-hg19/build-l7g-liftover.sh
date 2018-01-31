#!/bin/bash

VERBOSE=1
SKIP_GEN=0

libver="00"

tagset="$1"
reffa="$2"
cytoband="$3"
refname="$4"
chromlist="$5"

tilepath=0

if [[ "$tagset" == "" ]] || [[ "$reffa" == "" ]] || [[ "$cytoband" == "" ]] ; then
  echo "provide tagset, reffa, cytoband"
  exit 1
fi

if [[ "$chromlist" == "" ]] ; then
  chromlist=` cat <( seq 1 22 | sed 's/^/chr/' | tr '\n' ' ' ) <( echo chrX chrY chrM ) `
fi

if [[ "$refname" == "" ]] ; then
  refname=`basename $reffa .gz`
  refname=`basename $refname .fa`
  refname=`basename $refname .fasta`
fi

echo tagset $tagset
echo reffa $reffa
echo cytoband $cytoband
echo refname $refname

mkdir -p stage

## normalize chromosome names
##c
function norm_chrom {
  inpchrom=$1

  if [[ "$inpchrom" =~ ^[0-9] ]] ; then
    echo chr$inpchrom
    return
  fi

  if [[ "$inpchrom" =~ ^[XY]$ ]] ; then
    echo "chr$inpchrom"
    return
  fi

  if [[ "$inpchrom" =~ ^MT$ ]] ; then
    echo chrM
    return
  fi

  echo $inpchrom
}

if [[ "$SKIP_GEN" == "0" ]] ; then

for inpchrom in $chromlist ; do

  chrom=`norm_chrom $inpchrom`

  tilepath=`egrep -v '^#' $cytoband | nl -v 0 --number-width=1 | grep -P '^[0-9]*\t'$chrom'\t' | head -n1 | cut -f1`

  while IFS='' read -r line  || [[ -n "$line" ]]; do
    #chrom=`echo "$line" | cut -f1`

    pos0ref_start=`echo "$line" | cut -f2`
    pos0ref_end_noninc=`echo "$line" | cut -f3`

    pos1ref_start=`expr 1 + $pos0ref_start`
    pos1ref_end_inc=`expr 1 + $pos0ref_end_noninc - 1`

    hxp=`printf "%04x" $tilepath`

    if [[ "$VERBOSE" == "1" ]] ; then
      echo $chrom $pos0ref_start $pos0ref_end_noninc ... $pos1ref_start $pos1ref_end_inc "($tilepath $hxp)"
      echo "---"
      echo "chrom $chrom"
      echo "pos0ref_start $pos0ref_start"
      echo "pos0ref_end_noninc $pos0ref_end_noninc"
      echo "pos1ref_start $pos1ref_start"
      echo "pos1ref_end_inc $pos1ref_end_inc"
      echo "tilepath $tilepath"
      echo "hxp $hxp"
      echo "cmd: "
      echo "  tile-liftover -T <( samtools faidx $tagset $hxp.$libver | egrep -v '^>' | tr -d '\n' ) \
      -p $tilepath \
      -R <( samtools faidx $reffa $inpchrom:$pos1ref_start-$pos1ref_end_inc | egrep -v '^>' | tr -d '\n' | tr '[:upper:]' '[:lower:]' ) \
      -s $pos0ref_start \
      -N $refname \
      -c $chrom > stage/$hxp.liftover"
      echo "---"

    fi

    tile-liftover -T <( samtools faidx $tagset $hxp.$libver | egrep -v '^>' | tr -d '\n' ) \
      -p $tilepath \
      -R <( samtools faidx $reffa $inpchrom:$pos1ref_start-$pos1ref_end_inc | egrep -v '^>' | tr -d '\n' | tr '[:upper:]' '[:lower:]' ) \
      -s $pos0ref_start \
      -N $refname \
      -c $chrom > stage/$hxp.liftover

    tilepath=`expr 1 + $tilepath`

  done < <( cat $cytoband | egrep -v '^#' | grep -P '^'$chrom'\t'  )

done

fi

afn="assembly.$libver.$refname.fw"

cat stage/*.liftover | bgzip -c > $afn.gz
bgzip -r $afn.gz
zcat $afn | tile-assembly-index - > $afn.fwi

cp $afn.fwi $afn.gz.fwi

if [[ "$SKIP_GEN" == "0" ]] ; then
  rm -rf stage
fi

