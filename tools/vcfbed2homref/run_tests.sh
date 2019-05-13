#!/bin/bash

export verbose=1

tdir="testdata"
ref="$tdir/ref.fa.gz"
bed="$tdir/small.bed"
vcf="$tdir/small.vcf"

ref_b="$tdir/bed-bound-error.fa.gz"
bed_b="$tdir/bed-bound-error.bed"
vcf_b="$tdir/bed-bound-error.vcf.gz"

function ok_or_exit {
  a="$1"
  b="$2"
  name="$3"

  if [[ "$a" != "$b" ]] ; then
    if [[ "$verbose" == "1" ]] ; then
      echo "$name: $a != $b"
    fi
    exit 1
  fi


  if [[ "$verbose" == "1" ]] ; then
    echo "$name: ok"
  fi

}


## test that the headers looks the same.
## for some reason, htslib adds in an extra header line
## so take it out.
## Thsi might be too brittle and we might need to rethink
## this test.
##
a=`./vcfbed2homref -r $ref \
  -b $bed \
  $vcf | \
  grep -P '^#' | \
  grep -v -P '^##FILTER=<ID=PASS,Description="All filters passed">' | \
  md5sum | cut -f1 -d' '`
b=`grep -P '^#' $vcf | md5sum | cut -f1 -d' '`

ok_or_exit "$a" "$b" "header check"

###


## Make sure removing the added lines gives back the original
## VCF content
##
a=`./vcfbed2homref -r $ref \
  -b $bed \
  $vcf | \
  grep -v -P '^#' | \
  grep -v  'NON_REF' | \
  md5sum | cut -f1 -d' '`
b=`grep -v -P '^#' $vcf | md5sum | cut -f1 -d' '`

ok_or_exit "$a" "$b" "original content check"

###

./tests/check-vcf-bed.py $vcf $bed <( ./vcfbed2homref -r $ref -b $bed $vcf )
r=$?

if [[ "$r" -ne 0 ]] ;then
  echo "ERROR: start/end check"
  exit 1
elif [[ "$verbose" -eq 1 ]] ; then
  echo "start/end check: ok"
fi

##

./tests/check-vcf-bounds.py <( ./vcfbed2homref -r $ref_b -b $bed_b $vcf_b )
r=$?

if [[ "$r" -ne 0 ]] ;then
  echo "ERROR: bad vcf conversion $vcf_b (mangled regions)"
  exit 1
elif [[ "$verbose" -eq 1 ]] ; then
  echo "region-check: ok"
fi

##

if [[ "$verbose" -eq 1 ]] ; then
  echo "ok"
fi

exit 0
