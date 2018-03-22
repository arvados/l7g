#!/bin/bash
#
#

VERBOSE=1
SKIP_BAND=0

export cgf_dir="$1"
export sglf_dir="$2"
export gvcf_dir="$3"

export chrom="$4"

export afn="$5"
export ref_fa="$6"

export gvcf_pfx="$7"
export gvcf_sfx="$8"

export outfile="$9"

export ref=`basename $ref_fa .fa.gz`
export aidx="$afn.fwi"

if [[ "$outfile" == "" ]] ; then
  outfile=/dev/stdout
else
  rm -rf $outfile
fi


if [[ "$cgf_dir" == "" ]] || \
   [[ "$sglf_dir" == "" ]] || \
   [[ "$sglf_dir" == "" ]] || \
   [[ "$chrom" == "" ]] || \
   [[ "$ref_fa" == "" ]] || \
   [[ "$afn" == "" ]] ; then
  echo "usage:" >> $outfile
  echo "" >> $outfile
  echo "  ./verify-conversion-batch-gvcf-cgf.sh <cgf_dir> <sglf_dir> <gvcf_dir> <tileassembly> <ref.fa> [chrom] [gvcf_prefix] [gvcf_suffix}" >> $outfile
  echo "" >> $outfile
  exit -1
fi

if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "## cgf_dir: $cgf_dir" >> $outfile
  echo "## sglf_dir: $sglf_dir" >> $outfile
  echo "## gvcf_dir: $gvcf_dir" >> $outfile
  echo "## tileassembly: $afn" >> $outfile
  echo "## ref_fa: $ref_fa" >> $outfile
  echo "## chrom: $chrom" >> $outfile
  echo "## gvcf_pfx: $gvcf_pfx" >> $outfile
  echo "## gvcf_sfx: $gvcf_sfx" >> $outfile
fi

export cgf_fns=$( for base_fn in `ls $cgf_dir/*.cgf` ; do cgf_fn="$cgf_dir/$base_fn" ; echo $base_fn ; done )
export rep_cgf=$( for base_fn in `ls $cgf_dir/*.cgf` ; do cgf_fn="$cgf_dir/$base_fn" ; echo $base_fn ; done | head -n1 )

if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "## processing $chrom" >> $outfile
  echo "## $cgf_fns" >> $outfile
fi

####
####

beg_hxp=`egrep ":$chrom:" $afn.fwi | head -n1 | cut -f1 | cut -f3 -d':'`
end_hxp_inc=`egrep ":$chrom:" $afn.fwi | tail -n1 | cut -f1 | cut -f3 -d':'`

beg_p=`cat <( echo "ibase=16;" ) <( echo "$beg_hxp" | tr '[:lower:]' '[:upper:]' ) | bc`
end_p_inc=`cat <( echo "ibase=16;" ) <( echo "$end_hxp_inc" | tr '[:lower:]' '[:upper:]' ) | bc`

n_p=`echo "$end_p_inc - $beg_p + 1" | bc`

if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "## tilepath range: 0x$beg_hxp to 0x$end_hxp_inc inclusive [$beg_p-$end_p_inc]" >> $outfile
fi

band_fn=`mktemp`
band_hash=`mktemp`
gvcf_hash=`mktemp`

if [[ "$VERBOSE" == "1" ]] ; then
  echo "rep_cgf: $rep_cgf" >> "$outfile"
fi

export skip_tilepath_regex='xxxx'

for p in `seq $beg_p $end_p_inc` ; do
  hxp=`printf "%04x" $p`

  if [[ "$VERBOSE" == "1" ]] ; then
    echo "## cgft -b $p $rep_cgf | head -n1 | tr -d '[]' | sed 's/^  *//' | sed 's/ *$//' | tr ' ' '\n' | wc -l" >> "$outfile"
  fi

  c=`cgft -b $p $rep_cgf | head -n1 | tr -d '[]' | sed 's/^  *//' | sed 's/ *$//' | tr ' ' '\n' | wc -l`
  if [[ "$c" -le 1 ]] ; then
    skip_tilepath_regex="$skip_tilepath_regex|$hxp"
    continue
  fi
done

if [[ "$VERBOSE" == "1" ]] ; then
  echo "## skip_tilepath_regex: $skip_tilepath_regex" >> "$outfile"
fi

tilepath_list=$( cat <( echo "ibase=16;" ) \
  <( seq $beg_p $end_p_inc | \
     xargs -n1 -I{} printf "%04x\n" {} | \
     grep -v -P "$skip_tilepath_regex" | \
     tr '[:lower:]' '[:upper:]' ) | \
  bc | \
  tr '\n' ',' | \
  sed 's/^,//' | sed 's/,$//' )

if [[ "$VERBOSE" == "1" ]] ; then
  echo "## tilepath_list $tilepath_list" >> "$outfile"
fi

## loading the sglf into memory is the slow part
## so we 'batch' the hashes derived from the sequence
## fromt eh band files on a tilepath basis.
##
rm -f "$band_hash"
for p in `seq $beg_p $end_p_inc` ; do
  rm -f $band_fn

  hxp=`printf "%04x" $p`
  if [[ "$hxp" =~ $skip_tilepath_regex ]] ; then
    echo "## SKIPPING TILEPATH $hxp (band in $skip_tilepath_regex)" >> "$outfile"
    continue
  fi

  for cgf_fn in $cgf_fns; do
    cgft -b $p $cgf_fn >> $band_fn
  done

  #tileband-hash -L <( zcat `seq $beg_p $end_p_inc | xargs -n1 -I{} printf $sglf_dir/"%04x".sglf.gz"\n" {} | egrep -v 031f | tr '\n' ' ' ` ) \
  #  -T $beg_p+$n_p \
  tileband-hash -L <( zcat $sglf_dir/$hxp.sglf.gz ) \
    -T $p \
    $band_fn >> $band_hash

done

###
### create sequence hashes derived from gVCF files
###

gvcf_tdir=`mktemp -d`

## we do the same 'batching' for the gvcf to sequence (to hash)
## conersion so we can easily compare the lis tof hashes
##
while read line ; do

  for cgf_fn in $cgf_fns; do
    dsid=`basename $cgf_fn .cgf`
    dsid=`basename $dsid .cgfv3`
    dsid=`basename $dsid .cgf3`
    gvcf_fn="$gvcf_dir/$dsid/$gvcf_pfx$dsid$gvcf_sfx"

    export gvcf="$gvcf_fn"
    if [[ ! -e "$gvcf_fn.tbi" ]] ; then

      tgvcf=$gvcf_tdir/`basename $gvcf_fn`

      ## reuse indexed gvcf is we've already copied it locally
      ## and indexed.
      ##
      if [[ ! -e "$tgvcf.tbi" ]] ; then

        if [[ "$VERBOSE" -eq 1 ]] ; then
          echo "## copying gVCF locally and indexing ($tgvcf)" >> $outfile
        fi

        #cp $gvcf_fn $tgvcf
        ln -s $gvcf_fn $tgvcf
        tabix $tgvcf

      else
        if [[ "$VERBOSE" -eq 1 ]] ; then echo "## reusing index file ($tgvcf)" fi >> $outfile ; fi
      fi

      gvcf="$tgvcf"
    fi


    if [[ "$VERBOSE" -eq 1 ]] ; then
      echo "## processing gvcf $gvcf_fn" >> $outfile
    fi

    export tilepath=`echo "$line" | cut -f1 | cut -f3 -d':'`
    export tilepath_start0=`tile-assembly range $afn $tilepath | tail -n1 | cut -f2`
    export tilepath_end0_noninc=`tile-assembly range $afn $tilepath | tail -n1 | cut -f3`
    export tilepath_len=`expr "$tilepath_end0_noninc" - "$tilepath_start0"` || true

    if [[ "$tilepath" =~ $skip_tilepath_regex ]] ; then
      echo "## SKIPPING TILEPATH $tilepath (in $skip_tilepath_regex)" >> "$outfile"
      continue
    fi

    if [[ "$VERBOSE" == "1" ]] ; then
      echo "## tilepath: $tilepath" >> "$outfile"
      echo "## tilepath_start0: $tilepath_start0" >> "$outfile"
      echo "## tilepath_end0_noninc: $tilepath_end0_noninc" >> "$outfile"
      echo "## tilepath_len: $tilepath_len" >> "$outfile"
    fi

    if [[ "$tilepath_len" -eq "0" ]] ; then
      echo "## SKIPPING EMPTY TILEPATH $tilepath (tilepath_len: $tilepath_len)" >> "$outfile"
      continue
    fi

    export tilepath_start1=`expr "$tilepath_start0" + 1` || true
    export tilepath_end1_inc=`expr "$tilepath_start0" + "$tilepath_len" + 1` || true

    ## find the bounds for gvcf snippet in question
    ##
    export gvcf_start1=`tabix $gvcf $chrom:$tilepath_start1-$tilepath_end1_inc | head -n1 | cut -f2`
    export gvcf_end1_inc=`tabix $gvcf $chrom:$tilepath_start1-$tilepath_end1_inc | tail -n1 | cut -f2`

    export gvcf_start0="$tilepath_start0"
    if [[ "$gvcf_start1" != "" ]] ; then gvcf_start0=`expr "$gvcf_start1" - 1` || true ; fi

    gvcf_tok_end1_inc=`tabix $gvcf $chrom:$tilepath_start1-$tilepath_end1_inc | grep END | tail -n1 | cut -f8 | cut -f2 -d'='`
    if [[ "$gvcf_tok_end1_inc" != "" ]] && [[ "$gvcf_tok_end1_inc" -gt "$gvf_end1_inc" ]] ; then
      gvcf_end1_inc="$gvcf_tok_end1_inc"
    fi

    if [[ "$VERBOSE" == "1" ]] ; then
      echo "## tilepath_start1: $tilepath_start1" >> "$outfile"
      echo "## tilepath_end1_inc: $tilepath_end1_inc" >> "$outfile"
      echo "## gvcf_start1: $gvcf_start1" >> "$outfile"
      echo "## gvcf_end1_inc: $gvcf_end1_inc" >> "$outfile"
      echo "## gvcf_start0: $gvcf_start0" >> "$outfile"
      echo "## gvcf_tok_end1_inc: $gvcf_tok_end1_inc" >> "$outfile"
    fi

    ## The window under consideration starts at the minimum
    ## of the tilepath start and the gvcf start
    ##
    export window_start0="$tilepath_start0"
    if [[ "$gvcf_start0" -lt "$window_start0" ]] ; then window_start0="$gvcf_start0" ; fi
    export window_start1=`expr "$window_start0" + 1` || true

    ## now find the end of the window by taking the maximum
    ## of the tilepath end, the last gvcf reported position
    ##

    export window_end1_inc="$tilepath_end0_noninc"
    if [[ "$gvcf_end1_inc" != "" ]] && [[ "$gvcf_end1_inc" -gt "$window_end1_inc" ]] ; then
      window_end1_inc="$gvcf_end1_inc"
    fi

    ## take the maximum of the 'END' field or the length of the reference sequence
    ## in the 'REF' column.
    ##
    fin_ent_len=`tabix $gvcf $chrom:$tilepath_start1-$tilepath_end1_inc | tail -n1 | cut -f4 | tr -d '\n' | wc -c`
    fin_ent_start1=`tabix $gvcf $chrom:$tilepath_start1-$tilepath_end1_inc | tail -n1 | cut -f2 | tr -d '\n' `
    fin_ent_end1_inc=`expr $fin_ent_start1 + $fin_ent_len - 1` || true

    if [[ "$fin_ent_end1_inc" != "" ]] && [[ "$fin_ent_end1_inc" -gt "$window_end1_inc" ]] ; then
      window_end1_inc="$fin_ent_end1_inc"
    fi

    if [[ "$window_start1" -gt "$window_end1_inc" ]] ; then
      echo "## SKIPPING EMPTY TILEPATH $tilepath (window: $window_start1-$window_end1_inc)" >> "$outfile"
      continue
    fi

    if [[ "$VERBOSE" == "1" ]] ; then
      echo "## fin_ent_len: $fin_ent_len" >> "$outfile"
      echo "## fin_ent_start1: $fin_ent_start1" >> "$outfile"
      echo "## fin_ent_end1_inc: $fin_ent_end1_inc" >> "$outfile"
      echo "## window_start0: $window_start0" >> "$outfile"
      echo "## window_end1_inc: $window_end1_inc" >> "$outfile"
    fi

    export tdir=`mktemp -d`

    export window_len=`expr "$window_end1_inc" "-" "$window_start1" + 1` || true

    if [[ "$VERBOSE" == "1" ]] ; then
      echo "## refstream $ref_fa "$chrom:$window_start1+$window_len" > $tdir/$tilepath.ref" >> "$outfile"
    fi

    refstream $ref_fa "$chrom:$window_start1+$window_len" > $tdir/$tilepath.ref

    cat <( echo -e '\n\n\n' ) <( tabix $gvcf $chrom:$window_start1-$window_end1_inc ) > $tdir/$tilepath.gvcf

    if [[ "$VERBOSE" == "1" ]] ; then
      echo "## pasta -action gvcf-rotini -start $window_start0 -chrom $chrom \
        -full-sequence \
        -refstream $tdir/$tilepath.ref \
        -i $tdir/$tilepath.gvcf | \
        pasta -action filter-rotini -start $tilepath_start0 -n $tilepath_len > $tdir/$tilepath.pa" >> "$outfile"
    fi

    pasta -action gvcf-rotini -start $window_start0 -chrom $chrom \
      -full-sequence \
      -refstream $tdir/$tilepath.ref \
      -i $tdir/$tilepath.gvcf | \
      pasta -action filter-rotini -start $tilepath_start0 -n $tilepath_len > $tdir/$tilepath.pa

		h0=`pasta -action rotini-alt0 -i $tdir/$tilepath.pa | tr -d '\n' | md5sum | cut -f1 -d' '`
		h1=`pasta -action rotini-alt1 -i $tdir/$tilepath.pa | tr -d '\n' | md5sum | cut -f1 -d' '`

    echo "$h0 $h1" >> $gvcf_hash

    rm -rf $tdir
  done

done < <( egrep '^'$ref':'$chrom':' $aidx )

if [[ "$gvcf_tdir" != "" ]] ; then
  rm -rf $gvcf_tdir
fi


x=`cat $band_hash | md5sum | cut -f1 -d' '`
y=`cat $gvcf_hash | md5sum | cut -f1 -d' '`

if [[ "$x" != "$y" ]] ; then
  echo "chrom: $chrom" >> $outfile
  echo "cgf: $cgf_fns" >> $outfile
  echo "band_hash: $band_hash, gvcf_hash: $gvcf_hash" >> $outfile
  echo "band_fn: $band_fn" >> $outfile
  echo "band_hash: $band_hash" >> $outfile
  echo "gvcf_hash: $gvcf_hash" >> $outfile
  echo "MISMATCH: $x != $y" >> $outfile
  diff $band_hash $gvcf_hash >> $outfile
else
  echo "## $chrom ok" >> $outfile
  echo "ok" >> $outfile
fi

rm -f $band_fn
rm -f $band_hash
rm -f $gvcf_hash

####
####


if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "finished" >> $outfile
fi

exit 0
