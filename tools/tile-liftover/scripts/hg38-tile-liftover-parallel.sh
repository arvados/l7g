#!/bin/bash
#

export hg38="$hg38"
export chroms="chr1 chr2 chr3 chr4 chr5 chr6 chr7 chr8 chr9 chr10 chr11 chr12 chr13 chr14 chr15 chr16 chr17 chr18 chr19 chr20 chr21 chr22 chrX chrY chrM"
export ta="../../tile-assembly/tile-assembly"
export tadb="$L7G_TILEASSEMBLY"
export taidx="$L7G_TILEASSEMBLY_INDEX"
export tagset="$L7G_TAGSET"

export dest_ref="hg38"
export build="hg38"

if [[ "$hg38" == "" ]] || [[ "$tadb" == "" ]] || [[ "$taidx" == "" ]] || [[ "$tagset" == "" ]] ; then

  echo ""
  echo "must provide exported variables:"
  echo "  hg38                    - location of reference hg38 FASTA file (indexed)"
  echo "  L7G_TILEASSEMBLY        - Lightning tile assembly (e.g. assembly.00.hg19.fw.gz)"
  echo "  L7G_TILEASSEMBLY_INDEX  - Lightning tile assembly index (e.g. assembly.00.hg19.fw.fwi)"
  echo "  L7G_TAGSET              - Lightning tile tagset (e.g. tagset.fa.gz)"
  echo ""

  exit 1
fi

if [[ "$1" != "" ]] ; then
  chroms="$1"
fi

## go past the end to try and pick up tags that might not have fallen
# before the cytoband cutoff in a previous reference build
#
export BUFFER_DEFAULT="40000000"
export buffer_bp="$BUFFER_DEFAULT"

if [[ "$2" != "" ]] ; then
  BUFFER_DEFAULT="$2"
fi

mkdir -p stage

function process_chrom {
  chrom=$1

  echo ">>> $chrom $taidx"

  prev_chrom="unk"
  prev_end_0ref=0

  last_hxp=`grep ":$chrom:" $taidx | cut -f1 | cut -f3 -d':' | tail -n1`

  while IFS='' read -r hxp || [[ -n "$hxp" ]] ; do
    p=`cat <( echo "ibase=16;" ) <( echo $hxp | tr '[:lower:]' '[:upper:]' ) | bc`


    nstep=`$ta range $tadb $hxp | egrep -v '^#' | cut -f1 |  tr -d '\n'`
    beg0ref=`$ta range $tadb $hxp | egrep -v '^#' | cut -f2 | tr -d '\n'`
    end0ref_noninc=`$ta range $tadb $hxp | egrep -v '^#' | cut -f3 | tr -d '\n'`
    chrom=`$ta range $tadb $hxp | egrep -v '^#' | cut -f4 | tr -d '\n'`
    orig_ref=`$ta range $tadb $hxp | egrep -v '^#' | cut -f5 | tr -d '\n'`

    n_bp=`expr $end0ref_noninc - $beg0ref + $buffer_bp`
    beg1ref=`expr $beg0ref + 1`

    last_tile_len="225"

    if [[ "$prev_chrom" != "$chrom" ]] ; then
      prev_end_0ref=0
    fi

    if [[ "$nstep" == "1" ]] ; then
      buffer_bp=`expr $buffer_bp + $end0ref_noninc - $beg0ref`
      prev_chrom="$chrom"
      continue
    fi

    ref_start=`expr $prev_end_0ref + 1`


    echo $hxp $chrom $nstep $beg0ref $end0ref_noninc $last_tile_len

    refrange="$ref_start+$n_bp"
    if [[ "$hxp" == "$last_hxp" ]] ; then
      refrange="$ref_start"
      last_tile_len="-1"
    fi

    echo "../tile-liftover -p $p \
      -T <( refstream $tagset $hxp.00 | tr -d '\n' | fold -w 24 ) \
      -R <( refstream $hg38 $chrom:$refrange ) \
      -s $prev_end_0ref \
			-N $build \
      -c $chrom \
      -M $last_tile_len > stage/$hxp.assembly"

    ../tile-liftover -p $p \
      -T <( refstream $tagset $hxp.00 | tr -d '\n' | fold -w 24 ) \
      -R <( refstream $hg38 $chrom:$refrange ) \
      -s $prev_end_0ref \
			-N $build \
      -c $chrom \
      -M $last_tile_len > stage/$hxp.assembly

    prev_chrom="$chrom"
    assembly_end=`tail -n1 stage/$hxp.assembly | cut -f2 | tr -d ' '`
    if [[ "$assembly_end" == "" ]] ; then
      continue
    fi

    prev_end_0ref=`tail -n1 stage/$hxp.assembly | cut -f2 | tr -d ' '`

    buffer_bp="$BUFFER_DEFAULT"

  done < <( grep ":$chrom:" $taidx | cut -f1 | cut -f3 -d':' )

}
export -f process_chrom

function echo_test {
  echo $1
}
export -f echo_test

echo "$chroms" "($BUFFER_DEFAULT)"
echo "$chroms" | tr ' ' '\n' | parallel --no-notice --max-procs 10 process_chrom {}

## fill in unprocessed tilpaths and collect them into final lifted over tile assembly
## file.
##

idir="stage"
odir="out-data"

prev_chrom="unk"
chrom="unk"

dst_tafn="$odir/assembly.00.$build.fw"
dst_ta_idx_fn="$dst_tafn.fwi"

mkdir -p $odir
rm -f "$dst_tafn"

# fill in missing assembly files
#
for p in {0..862} ; do
  hxp=`printf "%04x" $p`

  if [[ ! -e "$idir/$hxp.assembly" ]] || [[ ! -s "$idir/$hxp.assembly" ]] ; then
    prev_addr=`tail -n1 "$dst_tafn" | cut -f2 | tr -d ' '`
    echo ">$build:$chrom:$hxp" >> "$dst_tafn"
    printf "%04x\t%10i\n" 0 $prev_addr >> "$dst_tafn"
		continue
  fi

  cat $idir/$hxp.assembly >> "$dst_tafn"

  prev_chrom="$chrom"
  chrom=`head -n1 $idir/$hxp.assembly | cut -f2 -d':' `

done


../tile-assembly-index "$dst_tafn" > "$dst_ta_idx_fn"
bgzip -f "$dst_tafn"
