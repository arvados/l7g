#!/bin/bash

set -e
set -o pipefail

vcfdir="$1"

vcfchr1=`ls $vcfdir/*.chr1.*`
sample=`basename $vcfchr1 | cut -d '.' -f 1`
chroms=(chr1 chr2 chr3 chr4 chr5 chr6 chr7 chr8 chr9 chr10 chr11 chr12 chr13 chr14 chr15 chr16 chr17 chr18 chr19 chr20 chr21 chr22 chrX chrY chrM)

files=$(for chrom in ${chroms[@]}; do echo "$vcfdir/$sample.raw_variants.$chrom.g.vcf.gz"; done)
echo "files: ${files[@]}"

bcftools concat ${files[@]} -n -O z -o $sample.vcf.gz
tabix $sample.vcf.gz
