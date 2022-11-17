#!/bin/bash

set -eo pipefail

sampleid="$1"
vcf="$2"
ref="$3"
mask="$4"

haplotypes=(1 2)

for haplotype in ${haplotypes[@]}; do
  bcftools consensus --fasta-ref $ref --haplotype $haplotype --mask $mask $vcf | bgzip -c > "$sampleid"."$haplotype".fa.gz
done
