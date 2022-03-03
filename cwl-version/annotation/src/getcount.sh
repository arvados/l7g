#!/bin/bash

set -e
set -o pipefail

sample=$1
vcf=$2

total=`zcat $vcf | awk '!(/^#/)' | wc -l`
gnomad=`zcat $vcf | awk '(!(/^#/) && /AF/)' | wc -l`
percentage=`awk -v n="$gnomad" -v d="$total" 'BEGIN {print n/d*100}'`

echo "$sample: $gnomad out of $total variants ($percentage%) have gnomad AF"
