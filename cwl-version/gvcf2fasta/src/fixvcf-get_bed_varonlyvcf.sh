#!/bin/bash

set -eo pipefail

sampleid="$1"
vcf="$2"
gqcutoff="$3"
genomebed="$4"

bcftools view --trim-alt-alleles $vcf | egrep -v "\*|<NON_REF>" | tee \
  >( /gvcf_regions/gvcf_regions.py --min_GQ $gqcutoff - > "$sampleid".bed ) \
  >( rtg vcffilter -i - -o - --remove-overlapping | awk '{if ($5 != ".") print $0}' | bgzip -c > "$sampleid"_varonly.vcf.gz ) \
  > /dev/null

bedtools subtract -a $genomebed -b "$sampleid".bed > "$sampleid"_nocall.bed
rm "$sampleid".bed
tabix "$sampleid"_varonly.vcf.gz
