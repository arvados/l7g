#!/bin/bash

#ref="/data-sdd/data/ref/human_g1k_v37.fa.gz"
ref="../testdata/ref.fa.gz"
igvcf="small.gvcf"

../vcfbed2homref \
  -r ../testdata/ref.fa.gz \
  -b ../testdata/small.bed \
  ../testdata/small.vcf > "$igvcf"

chrom="1"
s1=238902
e1inc=1032960
s0=`echo "$s1 - 1" | bc `


pasta -i <( cat $igvcf | grep -P '^1\t' ) \
  -a gvcf-rotini \
  -r <( refstr $ref $chrom:$s1-$e1inc ) \
  -s $s0 \
  --chrom $chrom > small-$chrom.pa


chrom="2"
s1=12188
e1inc=547281
s0=`echo "$s1 - 1" | bc `

pasta -i <( cat $igvcf | grep -P '^2\t' ) \
  -a gvcf-rotini \
  -r <( refstr $ref $chrom:$s1-$e1inc ) \
  -s $s0 \
  --chrom $chrom > small-$chrom.pa



