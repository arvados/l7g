#!/bin/bash
#
# This converts a single GFF file.
# Assumptions:
#  - UCSC chromosome naming (e.g. chr1, chrX, chrM, etc.).
#  - The GFF are indexed
#  - The GFF are in "Harvard PGP" format
#
# The conversion takes in an indexed GFF file that converts
# to PASTA then to FastJ on a tilepath by tilepath basis.
# The intermediate PASTA file is not saved.
# The FastJ are compressed and indexed.
#
# The major complication is figuring out the appropriate
# window for the PASTA conversion and restricting afterwards.
#
# USAGE:
#
#   ./convert_gff_to_fastj.sh <GFF_file> <TAGSET> <TILE_ASSEMBLY> <REF_FA> [<REF>]
#
# where <REF> defaults to 'hg19'.
#
# for example:
#
#  ./convert_gff_to_fastj.sh \
#    ./hu826751-GS03052-DNA_B01.gff.gz \
#   $L7G_DIR/tagset.fa/tagset.fa.gz \
#   $L7G_DIR/assembly/assembly.00.hg19.fw.gz \
#   $L7G_DIR/ref/hg19.fa.gz
#
# Will create 863 compressed FastJ files and their indexes.
#


set -a

DEBUG=1

function _q {
  echo $1
  exit 1
}

export gff="$1"
export tagset="$2"
export afn="$3"
export reffa="$4"
export ref="$5"
export out_name="$6"

if [ "$gff" == "" ] ; then
  echo "provide file"
  exit 1
fi

if [ "$tagset" == "" ] || [ "$afn" == "" ] || [ "$reffa" == "" ] ; then
  echo "provide tagset, tile assembly and FASTA reference"
  exit 1
fi

trap "ERROR: $gffInitial $path ; exit" ERR

if [ "$ref" == "" ] ; then
  export ref='hg19'
fi

if [[ "$out_name" == "" ]] ; then
  out_name=`basename $gff .gz`
  out_name=`basename $out_name .gff`
fi

odir="$out_name"
mkdir -p "$odir"

## index file for tile assembly
##
export aidx="$afn.fwi"

for chrom in chr1 chr2 chr3 chr4 chr5 chr6 chr7 chr8 chr9 chr10 chr11 chr12 chr13 chr14 chr15 chr16 chr17 chr18 chr19 chr20 chr21 chr22 chrX chrY chrM ; do

  echo "#### $gff process $chrom"


  while read line
  do

    ## parse tile path from tile assembly index
    ##
    tilepath=`echo "$line" | cut -f1 | cut -f3 -d':'`

    echo ">>> $tilepath"

    ## find 0ref start and end (non-inclusive) of the tile path
    ##
    #export ref_start0=`l7g assembly-range $afn $tilepath | tail -n1 | cut -f2`
    #export ref_end0=`l7g assembly-range $afn $tilepath | tail -n1 | cut -f3`

    export ref_start0=`tile-assembly range $afn $tilepath | tail -n1 | cut -f2`
    export ref_end0=`tile-assembly range $afn $tilepath | tail -n1 | cut -f3`

    ## convert to 1ref (end still non-inclusive)
    ##
    export ref_start1=`expr "$ref_start0" + 1`
    export ref_end1=`expr "$ref_end0" + 1`

    ## Find the actual reference start and end (inclusive) in the GFF file by finding the first
    ## entry in the GFF file that has information
    ##
    export realstart1=`tabix $gff $chrom:$ref_start1-$ref_end0 | head -n1 | cut -f4`
    export realend1=`tabix $gff $chrom:$ref_start1-$ref_end0 | tail -n1 | cut -f5`

    ## default to our reference start and end
    ##
    if [ "$realstart1" == "" ] ; then
      realstart1=$ref_start1
    fi

    if [ "$realend1" == "" ] ; then
      realend1=$ref_end0
    fi

    ## find the "real" number of bp in the GFF file and the
    ## number of reference base pairs for the tilepath.
    ##
    export realdn=`expr "$realend1" - "$realstart1"`
    export dn=`expr "$ref_end0" - "$ref_start0"`

    ## Calculate the window we're actually going to use
    ## when converting to PASTA.
    ## Use the 'real' start of where the GFF begins, unless
    ## it shoots past the reference start, in which case, take
    ## the start oft he reference
    ##
    export window_start1="$realstart1"
    if [ "$realstart1" -ge "$ref_start1" ]
    then
      export realstart1="$ref_start1"
      export window_start1="$ref_start1"
    fi

    export window_start0=`expr "$window_start1" - 1` || true

    ## Do the same for the end of the window
    ##
    export window_end1="$realend1"
    if [ "$ref_end1" -ge "$realend1" ]
    then
      export realend1=$ref_end1
      export window_end1="$ref_end1"
    fi

    export window_end0=`expr "$window_end1" - 1` || true

    ## Finally, convert to PASTA using the window to get the appropriate
    ## reference sequence, and then filter only the portion of the PASTA sequence
    ## we want.
    ##
    pasta -action gff-rotini -start $window_start0 -chrom $chrom \
      -refstream <( refstream $reffa $chrom:$window_start1-$window_end1 ) \
      -i <( cat <( echo -e '\n\n\n' ) <( tabix $gff $chrom:$window_start1-$window_end1 ) ) | \
      pasta -action filter-rotini -start $ref_start0 -n $dn > $odir/$tilepath.pa

    ## Convert from PASTA to FastJ
    ##
    pasta -action rotini-fastj -start $ref_start0 -tilepath $tilepath -chrom $chrom -build $ref \
      -i $odir/$tilepath.pa \
      -assembly <( tile-assembly tilepath $afn $tilepath ) \
      -tag <( cat <( samtools faidx $tagset $tilepath.00 | egrep -v '^>' | tr -d '\n' | fold -w 24 ) <( echo "" ) ) > $odir/$tilepath.fj

    ## Remove the temporary pasta file,
    ## force the compression and index of the
    ## FastJ file.
    ##
    rm $odir/$tilepath.pa
    bgzip -f $odir/$tilepath.fj
    bgzip -r $odir/$tilepath.fj.gz

  done < <( egrep '^'$ref':'$chrom':' $aidx )

done # chrom
