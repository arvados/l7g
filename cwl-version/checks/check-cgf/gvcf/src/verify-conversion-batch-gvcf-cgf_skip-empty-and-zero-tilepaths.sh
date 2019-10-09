#!/bin/bash

VERBOSE=1

export cgf_dir="$1"
export sglf_dir="$2"
export gvcf_dir="$3"
export check_num="$4"

export chrom="$5"

export afn="$6"
export ref="$7"
export ref_fa="$8"

export aidx="${afn%.*}.fwi"

if [[ "$cgf_dir" == "" ]] || \
   [[ "$sglf_dir" == "" ]] || \
   [[ "$gvcf_dir" == "" ]] || \
   [[ "$check_num" == "" ]] || \
   [[ "$chrom" == "" ]] || \
   [[ "$ref" == "" ]] || \
   [[ "$ref_fa" == "" ]] || \
   [[ "$afn" == "" ]] ; then
  echo "usage:"
  echo ""
  echo "  ./verify-conversion-batch-gvcf-cgf_skip-empty-and-zero-tilepaths.sh <cgf_dir> <sglf_dir> <gvcf_dir> <check_num> <chrom> <tileassembly> <ref> <ref_fa>"
  echo ""
  exit 1
fi

if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "## cgf_dir: $cgf_dir"
  echo "## sglf_dir: $sglf_dir"
  echo "## gvcf_dir: $gvcf_dir"
  echo "## check_num: $check_num"
  echo "## chrom: $chrom"
  echo "## tileassembly: $afn"
  echo "## ref: $ref"
  echo "## ref_fa: $ref_fa"
fi

export cgf_fns=$( for base_fn in `ls $cgf_dir/*.cgf` ; do cgf_fn="$cgf_dir/$base_fn" ; echo $base_fn ; done | head -n$check_num )
export rep_cgf=$( for base_fn in `ls $cgf_dir/*.cgf` ; do cgf_fn="$cgf_dir/$base_fn" ; echo $base_fn ; done | head -n1 )

if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "## processing $chrom"
  echo "## $cgf_fns"
fi

####
####

beg_hxp=`egrep ":$chrom:" $aidx | head -n1 | cut -f1 | cut -f3 -d':'`
end_hxp_inc=`egrep ":$chrom:" $aidx | tail -n1 | cut -f1 | cut -f3 -d':'`

beg_p=`cat <( echo "ibase=16;" ) <( echo "$beg_hxp" | tr '[:lower:]' '[:upper:]' ) | bc`
end_p_inc=`cat <( echo "ibase=16;" ) <( echo "$end_hxp_inc" | tr '[:lower:]' '[:upper:]' ) | bc`

n_p=`echo "$end_p_inc - $beg_p + 1" | bc`

if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "## tilepath range: 0x$beg_hxp to 0x$end_hxp_inc inclusive [$beg_p-$end_p_inc]"
fi

mkdir band_hash_dir
band_hash="band_hash_file"

if [[ "$VERBOSE" == "1" ]] ; then
  echo "rep_cgf: $rep_cgf"
fi

export skip_tilepath_regex='xxxx'

for p in `seq $beg_p $end_p_inc` ; do
  hxp=`printf "%04x" $p`

  if [[ "$VERBOSE" == "1" ]] ; then
    echo "## cgft -b $p $rep_cgf | head -n1 | tr -d '[]' | sed 's/^  *//' | sed 's/ *$//' | tr ' ' '\n' | wc -l"
  fi

  c=`cgft -b $p $rep_cgf | head -n1 | tr -d '[]' | sed 's/^  *//' | sed 's/ *$//' | tr ' ' '\n' | wc -l`
  if [[ "$c" -le 1 ]] ; then
    skip_tilepath_regex="$skip_tilepath_regex|$hxp"
    continue
  fi
done

if [[ "$VERBOSE" == "1" ]] ; then
  echo "## skip_tilepath_regex: $skip_tilepath_regex"
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
  echo "## tilepath_list $tilepath_list"
fi

## loading the sglf into memory is the slow part
## so we 'batch' the hashes derived from the sequence
## from the band files on a tilepath basis.
##
for p in `seq $beg_p $end_p_inc` ; do
  hxp=`printf "%04x" $p`
  if [[ "$hxp" =~ $skip_tilepath_regex ]] ; then
    echo "## SKIPPING TILEPATH $hxp (band in $skip_tilepath_regex)"
    continue
  fi

  export band_fn="band_hash_dir/$hxp.band"

  for cgf_fn in $cgf_fns; do
    if [[ "$VERBOSE" -eq 1 ]] ; then
      echo "## cgft -b $p $cgf_fn >> $band_fn"
    fi

    cgft -b $p $cgf_fn >> $band_fn
  done

  if [[ "$VERBOSE" -eq 1 ]] ; then
    echo "## fjt -H -L <( zcat $sglf_dir/$hxp.sglf.gz ) $band_fn >> $band_hash"
  fi

  fjt -H -L <( zcat $sglf_dir/$hxp.sglf.gz ) $band_fn >> $band_hash

done

rm -rf band_hash_dir

###
### create sequence hashes derived from gVCF files
###

## we do the same 'batching' for the gvcf to sequence (to hash)
## conersion so we can easily compare the list tof hashes
##

mkdir gvcf_hash_dir
gvcf_hash="gvcf_hash_file"

while read line ; do

  for cgf_fn in $cgf_fns; do
    dsid=`basename $cgf_fn .cgf`
    gvcf_fn="$gvcf_dir/$dsid.vcf.gz"

    export gvcf="$gvcf_fn"

    if [[ "$VERBOSE" -eq 1 ]] ; then
      echo "## processing gvcf $gvcf_fn"
    fi

    export tilepath=`echo "$line" | cut -f1 | cut -f3 -d':'`
    export tilepath_start0=`l7g assembly-range $afn $tilepath | tail -n1 | cut -f2`
    export tilepath_end0_noninc=`l7g assembly-range $afn $tilepath | tail -n1 | cut -f3`
    export tilepath_len=`expr "$tilepath_end0_noninc" - "$tilepath_start0"` || true

    if [[ "$tilepath" =~ $skip_tilepath_regex ]] ; then
      echo "## SKIPPING TILEPATH $tilepath (in $skip_tilepath_regex)"
      continue
    fi

    if [[ "$VERBOSE" == "1" ]] ; then
      echo "## tilepath: $tilepath"
      echo "## tilepath_start0: $tilepath_start0"
      echo "## tilepath_end0_noninc: $tilepath_end0_noninc"
      echo "## tilepath_len: $tilepath_len"
    fi

    if [[ "$tilepath_len" -eq "0" ]] ; then
      echo "## SKIPPING EMPTY TILEPATH $tilepath (tilepath_len: $tilepath_len)"
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
      echo "## tilepath_start1: $tilepath_start1"
      echo "## tilepath_end1_inc: $tilepath_end1_inc"
      echo "## gvcf_start1: $gvcf_start1"
      echo "## gvcf_end1_inc: $gvcf_end1_inc"
      echo "## gvcf_start0: $gvcf_start0"
      echo "## gvcf_tok_end1_inc: $gvcf_tok_end1_inc"
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
      echo "## SKIPPING EMPTY TILEPATH $tilepath (window: $window_start1-$window_end1_inc)"
      continue
    fi

    if [[ "$VERBOSE" == "1" ]] ; then
      echo "## fin_ent_len: $fin_ent_len"
      echo "## fin_ent_start1: $fin_ent_start1"
      echo "## fin_ent_end1_inc: $fin_ent_end1_inc"
      echo "## window_start0: $window_start0"
      echo "## window_end1_inc: $window_end1_inc"
    fi

    export window_len=`expr "$window_end1_inc" "-" "$window_start1" + 1` || true

    if [[ "$VERBOSE" == "1" ]] ; then
      echo "## refstream $ref_fa "$chrom:$window_start1+$window_len" > gvcf_hash_dir/$tilepath.ref"
    fi

    refstream $ref_fa "$chrom:$window_start1+$window_len" > gvcf_hash_dir/$tilepath.ref

    cat <( echo -e '\n\n\n' ) <( tabix $gvcf $chrom:$window_start1-$window_end1_inc ) > gvcf_hash_dir/$tilepath.gvcf

    if [[ "$VERBOSE" == "1" ]] ; then
      echo "## pasta -action gvcf-rotini -start $window_start0 -chrom $chrom \
        -full-sequence \
        -refstream gvcf_hash_dir/$tilepath.ref \
        -i gvcf_hash_dir/$tilepath.gvcf | \
        pasta -action filter-rotini -start $tilepath_start0 -n $tilepath_len > gvcf_hash_dir/$tilepath.pa"
    fi

    pasta -action gvcf-rotini -start $window_start0 -chrom $chrom \
      -full-sequence \
      -refstream gvcf_hash_dir/$tilepath.ref \
      -i gvcf_hash_dir/$tilepath.gvcf | \
      pasta -action filter-rotini -start $tilepath_start0 -n $tilepath_len > gvcf_hash_dir/$tilepath.pa

		h0=`pasta -action rotini-alt0 -i gvcf_hash_dir/$tilepath.pa | tr -d '\n' | md5sum | cut -f1 -d' '`
		h1=`pasta -action rotini-alt1 -i gvcf_hash_dir/$tilepath.pa | tr -d '\n' | md5sum | cut -f1 -d' '`

    echo "$h0 $h1" >> $gvcf_hash

  done

done < <( egrep '^'$ref':'$chrom':' $aidx )

rm -rf gvcf_hash_dir

x=`cat $band_hash | md5sum | cut -f1 -d' '`
y=`cat $gvcf_hash | md5sum | cut -f1 -d' '`

echo "chrom: $chrom"
echo "cgf: $cgf_fns"
echo "band_hash: $band_hash, gvcf_hash: $gvcf_hash"

if [[ "$x" != "$y" ]] ; then
  echo "MISMATCH: $x != $y"
  diff $band_hash $gvcf_hash
  echo "FAIL"
  exit 1
else
  echo "MATCH: $x = $y"
  echo "PASS"
  exit 0
fi
