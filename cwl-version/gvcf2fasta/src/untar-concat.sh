#!/bin/bash

set -e
set -o pipefail

sampleid="$1"
vcftar="$2"

tar -xzf $vcftar -C .
splitvcfdir=`ls`

chroms=(chr1 chr2 chr3 chr4 chr5 chr6 chr7 chr8 chr9 chr10 chr11 chr12 chr13 chr14 chr15 chr16 chr17 chr18 chr19 chr20 chr21 chr22 chrX chrY chrM)
splitvcfs=$(for chrom in ${chroms[@]}; do ls $splitvcfdir/*$chrom\_*gz; done)
echo "splitvcfs: ${splitvcfs[@]}"

bcftools concat ${splitvcfs[@]} -n -O z -o $sampleid.vcf.gz

rm -rf $splitvcfdir
