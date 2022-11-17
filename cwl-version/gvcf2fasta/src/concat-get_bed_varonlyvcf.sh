#!/bin/bash

set -eo pipefail

sampleid="$1"
splitvcfdir="$2"
gqcutoff="$3"
genomebed="$4"

chroms=(chr1 chr2 chr3 chr4 chr5 chr6 chr7 chr8 chr9 chr10 chr11 chr12 chr13 chr14 chr15 chr16 chr17 chr18 chr19 chr20 chr21 chr22 chrX chrY chrM)
splitvcfs=$(for chrom in ${chroms[@]}; do ls $splitvcfdir/*$chrom\.*gz; done)
echo "splitvcfs: ${splitvcfs[@]}"

bcftools concat ${splitvcfs[@]} -n | bcftools view --trim-alt-alleles | egrep -v "\*|<NON_REF>" | tee \
  >( /gvcf_regions/gvcf_regions.py --min_GQ $gqcutoff - > "$sampleid".bed ) \
  >( awk '{if ($5 != ".") print $0}' | bgzip -c > "$sampleid"_varonly.vcf.gz ) \
  > /dev/null

bedtools subtract -a $genomebed -b "$sampleid".bed > "$sampleid"_nocall.bed
rm "$sampleid".bed
tabix "$sampleid"_varonly.vcf.gz
