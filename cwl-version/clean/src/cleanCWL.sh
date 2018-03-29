#!/bin/bash

set -a

DEBUG=0
export MAKE_NEW_DIR=1
export VERBOSE=1

function _q {
  echo $1
  exit 1
}

gvcfDir="$1"
gvcfPrefix="$2"
cleanvcf="$3"

if [ "$gvcfDir" == "" ] ; then
  echo "provide inital directory"
  exit 1
fi

trap "ERROR: $gvcfDir $path ; exit" ERR

newdir='cleaned/'$gvcfPrefix

mkdir -p $newdir

for chrom in chr1 chr2 chr3 chr4 chr5 chr6 chr7 chr8 chr9 chr10 chr11 chr12 chr13 chr14 chr15 chr16 chr17 chr18 chr19 chr20 chr21 chr22 chrX chrY chrM ; do

   gvcfInitial=$gvcfDir'/filtered_'$gvcfPrefix.raw_variants.$chrom.gvcf.gz
   ifnInitial=`basename $gvcfInitial`
   stripped_name=`basename $ifnInitial .gvcf.gz`

   if [[ "$VERBOSE" -eq 1 ]] ; then
     echo "gvcfInitial: $gvcfInitial"
     echo "ifnInitial: $ifnInitial"
     echo "stripped_name: $stripped_name"
   fi

   zcat $gvcfInitial | $cleanvcf | bgzip -c > $newdir'/'$stripped_name'.gvcf.gz'

done #chrom
