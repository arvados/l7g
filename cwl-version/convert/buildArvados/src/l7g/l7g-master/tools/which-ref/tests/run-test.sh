#!/bin/bash

bin=../which-ref
cmp_seq="../test-data/human_g1k_v37-MT.fa ../test-data/hg19-chrM.fa"
seq="../test-data/hu826751-chrM.fa"

output=`$bin $cmp_seq $seq`
real_output=`echo "$output" | tr '\n' ' '`
expected_output="min_score: 63 min_idx: 1 name: ../test-data/hg19-chrM.fa "

if [[ "$real_output" != "$expected_output" ]] ; then
  echo "MISMATCH: Expected '$expected_output', got '$real_output'"
  exit -1
fi

cmp_seq="../test-data/hg38-chrM-1ref.tsv ../test-data/hg19-chrM-1ref.tsv"
query="../test-data/hu826751-chrM-1ref.tsv"

output=`$bin -N -C -M $cmp_seq $query`
if [[ "$output" != "1" ]] ; then
  echo "MISMATCH: expected to match index 1 reference, instead got '$output'"
  exit -2
fi

echo "ok"
exit 0
