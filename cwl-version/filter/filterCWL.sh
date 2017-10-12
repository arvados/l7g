#!/bin/bash

set -a

gffDir="$1"
gffPrefix="$2"
filter_gvcf="$3"

#gffDir='/data-sdd/cwl_tiling/convert/checkfastJ/keep/by_id/su92l-4zz18-t8ami4w0swhseki'
#gffPrefix='A-UPN-UP000241-CL-UPN-1437'
#filter_gvcf='/data-sdd/cwl_tiling/filter/filter-gvcf'

newdir='filtered'/$gffPrefix

mkdir -p $newdir

for chrom in chr1 chr2 chr3 chr4 chr5 chr6 chr7 chr8 chr9 chr10 chr11 chr12 chr13 chr14 chr15 chr16 chr17 chr18 chr19 chr20 chr21 chr22 chrX chrY chrM ; do
#for chrom in chr1; do

   gffInitial="$gffDir/$gffPrefix.raw_variants.$chrom.g.vcf.gz"
   echo "$gffInitial"

   ifnInitial=`basename $gffInitial`
   echo "$ifnInitial"

   stripped_name=`basename $ifnInitial .g.vcf.gz`

   echo "$stripped_name"

   zcat $gffInitial | $filter_gvcf 30 | gzip -c > $newdir'/filtered_'$stripped_name'.gvcf.gz'

   
done # chrom 
