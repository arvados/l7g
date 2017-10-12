#!/bin/bash

set -a

base="$1"
ref="$2"
tag="$3"
refM0="$4"
tagM0="$5"
refM1="$6"
tagM1="$7"
refM2="$8"
tagM2="$9"
whichref="${10}"
tmpdir="${11}"
numberify="${12}"

#export base="/data-sdd/cwl_tiling/testdir"
export files="${base}/*.gz"
#export pasta="/data-sdd/cwl_tiling/convert/buildArvados/dest_local/pasta"
#export refstream="/data-sdd/cwl_tiling/convert/buildArvados/dest_local/refstream"
#export ref="/data-sdd/cwl_tiling/convert/ref/hg19.fa"
#export tag="chrM"
#export refM0="/data-sdd/data/ref/human_g1k_v37.fa.gz"
#export tagM0="MT"
#export refM1="/data-sdd/data/ref/hg38.fa.gz"
#export tagM1="chrM"
#export refM2="/data-sdd/data/ref/hg19.fa.gz"
#export tagM2="chrM"
#export whichref="/data-sdd/cwl_tiling/convert/buildArvados/dest_local/which-ref"
#export tmpdir="tmp"
#export numberify="/data-sdd/cwl_tiling/convert/buildArvados/dest_local/numberify"

for f in $files;
do

  export f
  echo "Processing $f..."

  mkdir -p $tmpdir
  mkdir -p "ChrMref0"
  mkdir -p "ChrMref1"
  mkdir -p "ChrMref2"

  ifnInitial=`basename $f`
  extensionfull=${ifnInitial#*.}
  extensionfull='.'$extensionfull
  extensionshort=${extensionfull%.*}
  stripped_nameInitial=`basename $ifnInitial $extensionfull`
  gff=$tmpdir'/cleaned'$stripped_nameInitial$extensionshort
  
  #zcat $f | head -n2 > $gff
  #zcat $f | egrep -v '^#' >> $gff
 
  zcat $f | head -n2 | egrep '^#' > $gff
  zcat $f | egrep -v '^#' >> $gff

  bgzip -i $gff
  gff=$gff'.gz'
  tabix $gff
  export gff

  output=$($whichref -M -C -N \
 <( samtools faidx $refM0 $tagM0 | egrep -v '^>' | tr '[:upper:]' '[:lower:]' | tr -d '\n' | sed 's/\(.\)/\1\n/g' | egrep -v '^$' | $numberify 1 ) \
 <( samtools faidx $refM1 $tagM1 | tr '[:upper:]' '[:lower:]' | tr -d '\n' | sed 's/\(.\)/\1\n/g' | egrep -v '^$' | $numberify 1 ) \
 <( samtools faidx $refM2 $tagM2 | egrep -v '^>' | tr '[:upper:]' '[:lower:]' | tr -d '\n' | sed 's/\(.\)/\1\n/g' | egrep -v '^$' | $numberify 1 ) \
 <( tabix $gff $tag | grep ref_allele | grep ref_allele | cut -f4,9 | sed 's/^\([^\t]*\)\t.*ref_allele \(.\).*/\1\t\2/') ) 

  export newdir='ChrMref'$output
  mv $gff $newdir
  mv $gff'.tbi' $newdir
  mv $gff'.gzi' $newdir
  echo 'moved to '$newdir

done
