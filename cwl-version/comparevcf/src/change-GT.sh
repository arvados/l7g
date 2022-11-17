#!/bin/bash

set -eo pipefail

sampleid=$1
suffix=$2
vcf=$3
header=$4

cat $header <(bgzip -dc $vcf | egrep -v ^# | awk '{if ($4 != $5) print $1 "\t" $2 "\t" $3 "\t" $4 "\t" $5 "\t" $6 "\t" $7 "\t" $8 "\tGT\t0/1"}') | bgzip -c > "$sampleid"_"$suffix".vcf.gz
tabix "$sampleid"_"$suffix".vcf.gz
