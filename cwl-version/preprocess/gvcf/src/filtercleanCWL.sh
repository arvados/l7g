#!/bin/bash

set -a

gvcfInitial="$1"
filtergvcf="$2"
cutoff="$3"
cleanvcf="$4"

echo "$gffInitial"

ifnInitial=`basename $gvcfInitial`
echo "$ifnInitial"

stripped_name=`basename $ifnInitial .vcf.gz`

echo "$stripped_name"

zcat $gvcfInitial | $filtergvcf $cutoff | $cleanvcf | bgzip -c > $stripped_name'.vcf.gz'
tabix -p vcf $stripped_name'.vcf.gz'
