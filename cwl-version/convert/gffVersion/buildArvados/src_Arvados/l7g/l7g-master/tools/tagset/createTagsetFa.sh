#!/bin/bash
#
# Output (to stdout) a FASTA file for the tagset.
#

export tilelibver="00"
export cytofn="$1"
export bwfn="$2"
export reffa="$3"

export bw2bg="bigWigToBedGraph"
export tagsetFa="./src/tagsetFa"

if [[ "$cytofn" == "" ]] || [[ "bwfn" == "" ]] || [[ "$reffa" == "" ]] ; then
  echo "provide cytogenetic band, bigwig file and reference FASTA"
  exit 1
fi

count=0

while IFS='' read -r line || [[ -n "$line" ]]; do
  export chrom=`echo "$line" | cut -f1`
  export start0=`echo "$line" | cut -f2`
  export end0_noninc=`echo "$line" | cut -f3`

  hxp=`printf "%04x" $count`

  echo ">$hxp.$tilelibver"

#  prev=0
#  while read startpos ; do
#    startpos1ref=`expr $startpos + 1`
#    refstream $chrom:$startpos1ref+24
#  done < <( $bw2bg -chrom=$chrom -start=$start0 -end=$end0_noninc $bwfn /dev/stdout | \
#    grep -P '\t1$' | \
#    ./choose_tagset_startpos0_vestigial.py $start0 $end0_noninc | tail -n +2 )

  $bw2bg -chrom=$chrom -start=$start0 -end=$end0_noninc $bwfn /dev/stdout | \
    grep -P '\t1$' | \
    ./choose_tagset_startpos0_vestigial.py $start0 $end0_noninc | \
    tail -n +2 | \
    $tagsetFa -R $reffa -c $chrom

  count=`expr $count + 1`

done < <( cat $cytofn | egrep -v '^#' )
