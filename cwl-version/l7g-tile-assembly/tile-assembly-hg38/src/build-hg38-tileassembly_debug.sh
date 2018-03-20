#!/bin/bash
#

export SHELL='/bin/bash'

export tadb="$1"
export taidx="$tadb.fwi"
export tagset="$2"
export hg38="$3"

export chroms="chr1 chr2 chr3 chr4 chr5 chr6 chr7 chr8 chr9 chr10 chr11 chr12 chr13 chr14 chr15 chr16 chr17 chr18 chr19 chr20 chr21 chr22 chrX chrY chrM"
#export chroms="chr22"
export N_PARALLEL=10
#export N_PARALLEL=1

export ta="tile-assembly"

export dest_ref="hg38"
export build="hg38"

export SKIP_LIFTOVER="0"
export VERBOSE_FLAG="1"

## go past the end to try and pick up tags that might not have fallen
# before the cytoband cutoff in a previous reference build
#
export BUFFER_DEFAULT="40000000"

if [[ "$hg38" == "" ]] || [[ "$tadb" == "" ]] || [[ "$taidx" == "" ]] || [[ "$tagset" == "" ]] ; then

  echo ""
  echo "must provide:"
  echo "  L7G_TILEASSEMBLY        - Lightning tile assembly (e.g. assembly.00.hg19.fw.gz)"
  echo "  L7G_TAGSET              - Lightning tile tagset (e.g. tagset.fa.gz)"
  echo "  hg38                    - location of reference hg38 FASTA file (indexed)"
  echo ""

  exit 1
fi

if [[ "$4" != "" ]] ; then
  chroms="$4"
fi

if [[ "$5" != "" ]] ; then
  export BUFFER_DEFAULT="$5"
fi

if [[ "$6" != "" ]] ; then
  export $N_PARALLEL="$6"
fi

export buffer_bp="$BUFFER_DEFAULT"

mkdir -p stage

function process_chrom {
  chrom=$1

  echo ">>> $chrom $taidx"

  prev_chrom="unk"
  prev_end_0ref=0

  last_hxp=`grep ":$chrom:" $taidx | cut -f1 | cut -f3 -d':' | tail -n1`

  echo "???? $prev_chrom $prev_end_0ref $last_hxp"
  echo "grep ":$chrom:" $taidx | cut -f1 | cut -f3 -d':'"
  grep ":$chrom:" $taidx | cut -f1 | cut -f3 -d':'

  while IFS='' read -r hxp || [[ -n "$hxp" ]] ; do
    p=`cat <( echo "ibase=16;" ) <( echo $hxp | tr '[:lower:]' '[:upper:]' ) | bc`

    ## DEBUG
    echo ">>> $p"

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

      ## DEBUG
      echo "buffer_bp $buffer_bp, prev_chrom: $prev_chrom, chrom $chrom"
      echo "skipping..."

      continue
    fi

    ref_start=`expr $prev_end_0ref + 1`

    echo $hxp $chrom $nstep $beg0ref $end0ref_noninc $last_tile_len

    refrange="$ref_start+$n_bp"
    if [[ "$hxp" == "$last_hxp" ]] ; then
      refrange="$ref_start"
      last_tile_len="-1"
    fi

    echo "tile-liftover -p $p \
      -T <( refstream $tagset $hxp.00 | tr -d '\n' | fold -w 24 ) \
      -R <( refstream $hg38 $chrom:$refrange ) \
      -s $prev_end_0ref \
			-N $build \
      -c $chrom \
      -M $last_tile_len > stage/$hxp.assembly"

    tile-liftover -p $p \
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

if [[ "$SKIP_LIFTOVER" != "1" ]] ; then

  echo "$chroms" "($BUFFER_DEFAULT)"
  echo "$chroms" | tr ' ' '\n' | parallel --no-notice --max-procs $N_PARALLEL process_chrom {}

fi

## fill in unprocessed tilpaths and collect them into final lifted over tile assembly
## file.
##

idir="stage"
odir="out-data"

prev_chrom="unk"
chrom="unk"

dst_tafn="$odir/assembly.00.$build.fw"
dst_ta_idx_fn="$dst_tafn.fwi"

dst_ta_gz_fn="$dst_tafn.gz"
dst_ta_gz_idx_fn="$dst_tafn.gz.fwi"

mkdir -p $odir
rm -f "$dst_tafn"

if [[ "$VERBOSE_FLAG" == "1" ]] ; then
  echo "# FILLING IN"
fi

# fill in missing assembly files
#
export prev_chrom="unk"
for p in {0..862} ; do
  hxp=`printf "%04x" $p`

  chrom=`egrep ":$hxp" $taidx | cut -f1 | cut -f2 -d':' `

  if [[ ! -e "$idir/$hxp.assembly" ]] || [[ ! -s "$idir/$hxp.assembly" ]] ; then

    export prev_addr="0"
    if [[ "$prev_chrom" == "$chrom" ]] ; then

      if [[ "$VERBOSE" == "1" ]] ; then
        echo "## start of chromosome with no assembly file, using 0 prev_addr"
      fi

      prev_addr=`tail -n1 "$dst_tafn" | cut -f2 | tr -d ' '`
    fi

    echo ">$build:$chrom:$hxp" >> "$dst_tafn"
    printf "%04x\t%10i\n" 0 $prev_addr >> "$dst_tafn"


    if [[ "$VERBOSE_FLAG" == "1" ]] ; then
      echo "## filling in for $hxp ($build:$chrom:$hxp to $dst_tafn)"
      echo "## " `printf "%04x\t%10i\n" 0 $prev_addr` "to $dst_tafn"
    fi

		continue
  fi

  cat $idir/$hxp.assembly >> "$dst_tafn"
done


tile-assembly-index "$dst_tafn" > "$dst_ta_idx_fn"
bgzip -i -f "$dst_tafn"

pushd $odir
rm -f `basename "$dst_ta_gz_idx_fn"`
#ln -s `basename "$dst_ta_idx_fn"` `basename "$dst_ta_gz_idx_fn"`
cp `basename "$dst_ta_idx_fn"` `basename "$dst_ta_gz_idx_fn"`
popd

mv out-data/assembly.* .
rmdir out-data
rm -rf stage
