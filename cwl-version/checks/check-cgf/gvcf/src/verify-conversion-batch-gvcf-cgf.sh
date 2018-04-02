#!/bin/bash
#
#

VERBOSE=1

cgf_dir="$1"
sglf_dir="$2"
gvcf_dir="$3"

chrom="$4"

afn="$5"
ref_fa="$6"

gvcf_pfx="$7"
gvcf_sfx="$8"

outfile="$9"

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

if [[ "$outfile" == "" ]] ; then
  outifle=/dev/stdout
fi

if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "cgf_dir: $cgf_dir" >> $outfile
  echo "sglf_dir: $sglf_dir" >> $outfile
  echo "gvcf_dir: $gvcf_dir" >> $outfile
  echo "tileassembly: $afn" >> $outfile
  echo "ref_fa: $ref_fa" >> $outfile
  echo "chrom: $chrom" >> $outfile
  echo "gvcf_pfx: $gvcf_pfx" >> $outfile
  echo "gvcf_sfx: $gvcf_sfx" >> $outfile
fi

cgf_fns=$( for base_fn in `ls $cgf_dir/*.cgf` ; do cgf_fn="$cgf_dir/$base_fn" ; echo $base_fn ; done )


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
  echo "tilepath range: 0x$beg_hxp to 0x$end_hxp_inc inclusive [$beg_p-$end_p_inc]" >> $outfile
fi

band_fn=`mktemp`
band_hash=`mktemp`
gvcf_hash=`mktemp`

for cgf_fn in $cgf_fns; do
  for p in `seq $beg_p $end_p_inc` ; do
    cgft -b $p $cgf_fn >> $band_fn
  done
done

if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "executing:" >> $outfile
  echo "tileband-hash -L <( zcat `seq $beg_p $end_p_inc | xargs -n1 -I{} printf " "$sglf_dir/"%04x".sglf.gz {} | sed 's/^ //' ` ) \
    -T $beg_p+$n_p \
    $band_fn" >> $outfile
fi

tileband-hash -L <( zcat `seq $beg_p $end_p_inc | xargs -n1 -I{} printf " "$sglf_dir/"%04x".sglf.gz {} | sed 's/^ //' ` ) \
  -T $beg_p+$n_p \
  $band_fn > $band_hash

for cgf_fn in $cgf_fns; do
  gvcf_fn=`basename $cgf_fn .cgf`
  gvcf_fn=`basename $gvcf_fn .cgfv3`
  gvcf_fn=`basename $gvcf_fn .cgf3`
  #gvcf_fn="$gvcf_dir/$gvcf_pfx$gvcf_fn.raw_variants.$chrom.gvcf.gz"
  gvcf_fn="$gvcf_dir/$gvcf_pfx$gvcf_fn$gvcf_sfx"

  if [[ "$VERBOSE" -eq 1 ]] ; then
    echo "processing gvcf $gvcf_fn" >> $outfile
  fi

  h0=`tabix $gvcf_fn $chrom | \
    pasta -a gvcf-rotini --full-sequence -r <( refstream $ref_fa $chrom ) | \
    pasta -a rotini-alt0 | tr -d '\n' | md5sum | cut -f1 -d' '`
  h1=`tabix $gvcf_fn $chrom | \
    pasta -a gvcf-rotini --full-sequence -r <( refstream $ref_fa $chrom ) | \
    pasta -a rotini-alt1 | tr -d '\n' | md5sum | cut -f1 -d' '`

  echo "$h0 $h1" >> $gvcf_hash
done

x=`cat $band_hash | md5sum | cut -f1 -d' '`
y=`cat $gvcf_hash | md5sum | cut -f1 -d' '`

if [[ "$x" != "$y" ]] ; then
  echo "MISMATCH: $x != $y" >> $outfile
  echo "chrom: $chrom" >> $outfile
  echo "cgf: $cgf_fns" >> $outfile
  echo "band_hash: $band_hash, gvcf_hash: $gvcf_hash" >> $outfile
  echo "band_fn: $band_fn" >> $outfile
  echo "band_hash: $band_hash" >> $outfile
  echo "gvcf_hash: $gvcf_hash" >> $outfile
  diff $band_hash $gvcf_hash >> $outfile
else
  echo "$chrom ok" >> $outfile
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
