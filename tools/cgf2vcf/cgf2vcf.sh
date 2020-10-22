#!/bin/bash

set -aeo pipefail

export makeheadergenomebed="$1"
export cgf2vcf="$2"
export steps2bed="$3"
export assembly="$4"
export annotationlib="$5"
export cgf="$6"

export sampleid=`basename $cgf .cgf`
export assemblyindex="${assembly%.*}.fwi"
export paths=$( cat $assemblyindex | cut -f1 | cut -d':' -f3 )

## make vcf header and genome bed
$makeheadergenomebed $assembly --outtype header --sampleid $sampleid > $sampleid.vcf
$makeheadergenomebed $assembly --outtype genomebed > genome.bed

> "$sampleid"_notcovered.bed

for path in $paths; do
  echo "## processing path $path"
  $cgf2vcf $path $assembly $annotationlib $cgf \
    --nocall "$sampleid"_"$path"_nocall.txt \
    --unannotated "$sampleid"_"$path"_unannotated.txt \
    --pathskipped "$sampleid"_"$path"_pathskipped.txt >> $sampleid.vcf
  pathskipped=`tr -d '\n' < "$sampleid"_"$path"_pathskipped.txt`
  if [ $pathskipped = "True" ]; then
    $steps2bed $path $assembly >> "$sampleid"_notcovered.bed
  else
    sort <( cat "$sampleid"_"$path"_nocall.txt "$sampleid"_"$path"_unannotated.txt ) > "$sampleid"_"$path"_notcovered.txt
    $steps2bed $path $assembly --stepsfile "$sampleid"_"$path"_notcovered.txt >> "$sampleid"_notcovered.bed
    rm "$sampleid"_"$path"_notcovered.txt
  fi
  rm "$sampleid"_"$path"_nocall.txt "$sampleid"_"$path"_unannotated.txt "$sampleid"_"$path"_pathskipped.txt
done

bedtools subtract -a genome.bed -b "$sampleid"_notcovered.bed > "$sampleid"_covered.bed
rm genome.bed "$sampleid"_notcovered.bed
bgzip $sampleid.vcf
tabix $sampleid.vcf.gz
