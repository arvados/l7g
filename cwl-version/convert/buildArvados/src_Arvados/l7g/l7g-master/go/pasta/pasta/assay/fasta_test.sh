#!/bin/bash


function _q {
  echo $1
  exit 1
}

####
####


a0=`./pasta -action pasta-fasta -i <( ./pasta -action rstream -param 'allele=1:seed=1234:ref-seed=11223344:n=125:p-snp=0.8' ) | md5sum | cut -f1 -d' ' `
a1=`./pasta -action pasta-fasta -i <( ./pasta -action rstream -param 'allele=1:seed=1234:ref-seed=11223344:n=125:p-snp=0.8' ) | ./pasta -action fasta-pasta -refstream <( ./pasta -action ref-rstream -param 'allele=1:n=125:ref-seed=11223344' ) | sed 's/^>.*//' | ./pasta -action pasta-fasta | md5sum | cut -f1 -d' '`

diff <( echo $a0 ) <( echo $a1 ) || _q "alt streams do not match"

####
####

echo ok
