#!/bin/bash

set -e

odir="assay/gvcf"
mkdir -p $odir

## GVCF with snps
##
./pasta -action rstream -param 'p-snp=0.5:ref-seed=11223344:n=1000:seed=1234' > $odir/gvcf-snp.inp
./pasta -action rotini-gvcf -i $odir/gvcf-snp.inp | ./pasta -action gvcf-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' ) > $odir/gvcf-snp.out

diff <( ./pasta -action rotini-ref -i $odir/gvcf-snp.inp ) <( ./pasta -action rotini-ref -i $odir/gvcf-snp.out )
diff <( ./pasta -action rotini-alt0 -i $odir/gvcf-snp.inp ) <( ./pasta -action rotini-alt0 -i $odir/gvcf-snp.out )
diff <( ./pasta -action rotini-alt1 -i $odir/gvcf-snp.inp ) <( ./pasta -action rotini-alt1 -i $odir/gvcf-snp.out )

echo ok-snp

## GVCF with indels
##
./pasta -action rstream -param 'p-indel=0.8:p-indel-length=0,3:p-nocall=0:ref-seed=11223344:n=1000:seed=1234' > $odir/gvcf-indel.inp
./pasta -action rotini-gvcf -i $odir/gvcf-indel.inp | ./pasta -action gvcf-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' ) > $odir/gvcf-indel.out

diff <( ./pasta -action rotini-ref -i $odir/gvcf-indel.inp ) <( ./pasta -action rotini-ref -i $odir/gvcf-indel.out )
diff <( ./pasta -action rotini-alt0 -i $odir/gvcf-indel.inp ) <( ./pasta -action rotini-alt0 -i $odir/gvcf-indel.out )
diff <( ./pasta -action rotini-alt1 -i $odir/gvcf-indel.inp ) <( ./pasta -action rotini-alt1 -i $odir/gvcf-indel.out )


echo ok-indel

## GVCF with nocall
##
ofn_b="gvcf-nocall"
./pasta -action rstream -param 'p-nocall=0.3:ref-seed=11223344:seed=1234' > $odir/$ofn_b.inp
./pasta -action rotini-gvcf -i $odir/$ofn_b.inp | ./pasta -action gvcf-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' )  > $odir/$ofn_b.out


diff <( ./pasta -action rotini-ref -i $odir/$ofn_b.inp ) <( ./pasta -action rotini-ref -i $odir/$ofn_b.out )
diff <( ./pasta -action rotini-alt0 -i $odir/$ofn_b.inp ) <( ./pasta -action rotini-alt0 -i $odir/$ofn_b.out )
diff <( ./pasta -action rotini-alt1 -i $odir/$ofn_b.inp ) <( ./pasta -action rotini-alt1 -i $odir/$ofn_b.out )

echo ok-nocall


## GVCF with indels and nocalls
##
#./pasta -action rstream -param 'p-nocall=0.3:p-indel=0.3:ref-seed=11223344:seed=1234'  > $odir/gvcf-indel-nocall.inp
#./pasta -action rstream -param 'p-nocall=0.3:p-indel=0.3:p-indel-nocall=0.8:ref-seed=11223344:seed=1234'  > $odir/gvcf-indel-nocall.inp
./pasta -action rstream -param 'p-nocall=0.3:p-indel=0.5:p-indel-nocall=0.8:ref-seed=11223344:seed=1234'  > $odir/gvcf-indel-nocall.inp
./pasta -action rotini-gvcf -i $odir/gvcf-indel-nocall.inp | ./pasta -action gvcf-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' ) > $odir/gvcf-indel-nocall.out

diff <( ./pasta -action rotini-ref -i $odir/gvcf-indel-nocall.inp ) <( ./pasta -action rotini-ref -i $odir/gvcf-indel-nocall.out ) || echo "error ref"
diff <( ./pasta -action rotini-alt0 -i $odir/gvcf-indel-nocall.inp ) <( ./pasta -action rotini-alt0 -i $odir/gvcf-indel-nocall.out ) || echo "error alt0"
diff <( ./pasta -action rotini-alt1 -i $odir/gvcf-indel-nocall.inp ) <( ./pasta -action rotini-alt1 -i $odir/gvcf-indel-nocall.out ) || echo "error alt1"

echo ok-indel-nocall

# test when there's an alt at the beginning of the stream
#

echo 'cttccttctttccttccctccctctttcctttcctttccttccctccctctttcctttcctttccttccctccctctttcctttccctccctcctcctttcctttcctttcctatcctttccctccatcctcctttccctcccctcccctcctctcccttccccttcccttcccttcctttcctttccttttttctttttctttcagactgagtctccctttgtcgcccaggctggagtgcagttgtgcaatctcagctcactgcaacctccgcctcctgggtttcaaatgattctcctgcctcactctcccaagtagctgggattatagctatgtgccacgacaccaggctaatttttgtattttaagtagagacagggtttcaccatgttggccaggctgatctcgaactccttacctcaagtgatccacctgcctcagcctcccaaaatgctaggatttcaggcgtaagccaccactcctggccccttagttactt' > $odir/ref-test0.inp
echo 'cttccttctttccttccctccctctttcctttcctttccttccctccctctttcctttcctttccttccctccctctttcctttcctttccttccctccctctttcctttccctccctcctcctttcctttcctttcctatcctttccctccatcctcctttccctcccctcccctcctctcccttccccttcccttcccttcctttcctttccttttttctttttctttcagactgagtctccctttgtcgcccaggctggagtgcagttgtgcaatctcagctcactgcaacctccgcctcctgggtttcaaatgattctcctgcctcactctcccaagtagctgggattatagctatgtgccacgacaccaggctaatttttgtattttaagtagagacagggtttcaccatgttggccaggctgatctcgaactccttacctcaagtgatccacctgcctcagcctcccaaaatgctaggatttcaggcgtaagccaccactcctggccccttagttactt' > $odir/alt-test0.inp

echo 'SSddddSSSSddddSSddddddSSSSddddSSSSSSddSSSSSSddccddttttccccddttttSSccttttttccccttttccccccttccccccttccttttttccccttttttccccttttttccccttttccccccttccccccttccttttttccccttttttccccttttttccccttttccccccttccccccttccttttttccccttttttccccccttccccccttccccttccccttttttccccttttttccccttttttccccttaattccccttttttccccccttccccaattccccttccccttttttccccccttccccccccttccccccccttccccttccttccccccttttccccccccttttccccccttttccccccttttccccttttttccccttttttccccttttttttttttccttttttttttccttttttccaaggaaccttggaaggttccttccccccttttttggttccggccccccaaggggccttggggaaggttggccaaggttttggttggccaaaattccttccaaggccttccaaccttggccaaaaccccttccccggccccttccccttggggggttttttccaaaaaattggaattttccttccccttggccccttccaaccttccttccccccaaaaggttaaggccttggggggaattttaattaaggccttaattggttggccccaaccggaaccaaccccaaggggccttaaaattttttttttggttaattttttttaaaaggttaaggaaggaaccaaggggggttttttccaaccccaattggttttggggccccaaggggccttggaattccttccggaaaaccttccccttttaaccccttccaaaaggttggaattccccaaccccttggccccttccaaggccccttccccccaaaaaaaattggccttaaggggaattttttccaaggggccggttaaaaggccccaaccccaaccttccccttggggccccccccttttaaggttttaacctttt' > $odir/snippet0.pa


diff <( ./pasta -action rotini-gvcf -i $odir/snippet0.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test0.inp | ./pasta -action rotini-ref ) <( ./pasta -action rotini-ref -i $odir/snippet0.pa ) || echo "snippet0 ref failed"
diff <( ./pasta -action rotini-gvcf -i $odir/snippet0.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test0.inp | ./pasta -action rotini-alt0 ) <( ./pasta -action rotini-alt0 -i $odir/snippet0.pa ) || echo "snippet0 alt0 failed"
diff <( ./pasta -action rotini-gvcf -i $odir/snippet0.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test0.inp | ./pasta -action rotini-alt1 ) <( ./pasta -action rotini-alt1 -i $odir/snippet0.pa ) || echo "snippet0 alt1 failed"

echo ok-snippet0

# test when there's an alt at the beginning of the stream and the ending alt base is different from reference
#

echo 'attccttctttccttccctccctctttcctttcctttccttccctccctctttcctttcctttccttccctccctctttcctttccctccctcctcctttcctttcctttcctatcctttccctccatcctcctttccctcccctcccctcctctcccttccccttcccttcccttcctttcctttccttttttctttttctttcagactgagtctccctttgtcgcccaggctggagtgcagttgtgcaatctcagctcactgcaacctccgcctcctgggtttcaaatgattctcctgcctcactctcccaagtagctgggattatagctatgtgccacgacaccaggctaatttttgtattttaagtagagacagggtttcaccatgttggccaggctgatctcgaactccttacctcaagtgatccacctgcctcagcctcccaaaatgctaggatttcaggcgtaagccaccactcctggccccttagttactt' > $odir/ref-test1.inp
echo 'cttccttctttccttccctccctctttcctttcctttccttccctccctctttcctttcctttccttccctccctctttcctttcctttccttccctccctctttcctttccctccctcctcctttcctttcctttcctatcctttccctccatcctcctttccctcccctcccctcctctcccttccccttcccttcccttcctttcctttccttttttctttttctttcagactgagtctccctttgtcgcccaggctggagtgcagttgtgcaatctcagctcactgcaacctccgcctcctgggtttcaaatgattctcctgcctcactctcccaagtagctgggattatagctatgtgccacgacaccaggctaatttttgtattttaagtagagacagggtttcaccatgttggccaggctgatctcgaactccttacctcaagtgatccacctgcctcagcctcccaaaatgctaggatttcaggcgtaagccaccactcctggccccttagttactt' > $odir/alt-test1.inp

echo 'SSddddSSSSddddSSddddddSSSSddddSSSSSSddSSSSSSddaaddttttccccddttttSSccttttttccccttttccccccttccccccttccttttttccccttttttccccttttttccccttttccccccttccccccttccttttttccccttttttccccttttttccccttttccccccttccccccttccttttttccccttttttccccccttccccccttccccttccccttttttccccttttttccccttttttccccttaattccccttttttccccccttccccaattccccttccccttttttccccccttccccccccttccccccccttccccttccttccccccttttccccccccttttccccccttttccccccttttccccttttttccccttttttccccttttttttttttccttttttttttccttttttccaaggaaccttggaaggttccttccccccttttttggttccggccccccaaggggccttggggaaggttggccaaggttttggttggccaaaattccttccaaggccttccaaccttggccaaaaccccttccccggccccttccccttggggggttttttccaaaaaattggaattttccttccccttggccccttccaaccttccttccccccaaaaggttaaggccttggggggaattttaattaaggccttaattggttggccccaaccggaaccaaccccaaggggccttaaaattttttttttggttaattttttttaaaaggttaaggaaggaaccaaggggggttttttccaaccccaattggttttggggccccaaggggccttggaattccttccggaaaaccttccccttttaaccccttccaaaaggttggaattccccaaccccttggccccttccaaggccccttccccccaaaaaaaattggccttaaggggaattttttccaaggggccggttaaaaggccccaaccccaaccttccccttggggccccccccttttaaggttttaacctttt' > $odir/snippet1.pa


diff <( ./pasta -action rotini-gvcf -i $odir/snippet1.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test1.inp | ./pasta -action rotini-ref ) <( ./pasta -action rotini-ref -i $odir/snippet1.pa ) || echo "snippet1 ref failed"
diff <( ./pasta -action rotini-gvcf -i $odir/snippet1.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test1.inp | ./pasta -action rotini-alt0 ) <( ./pasta -action rotini-alt0 -i $odir/snippet1.pa ) || echo "snippet1 alt0 failed"
diff <( ./pasta -action rotini-gvcf -i $odir/snippet1.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test1.inp | ./pasta -action rotini-alt1 ) <( ./pasta -action rotini-alt1 -i $odir/snippet1.pa ) || echo "snippet1 alt1 failed"

echo 'ok-snippet1'

# test when there's an alt at the end of the stream
#

echo 'aaaattccaaaaaaggggaaaaccaaaattttccccaaaattccccaaggttggttggggccccccttttccaattttttggttccttggggttggttggccttttggccttggggggaaggggggttttccttaaaattccaaggttttttccccccccccaaggaaccccccccaattttccccccttggttggccccttccttccccttccttccccttccccaaccccaaggaaaaaaggaattttaaggccaaggggaaggccccccaaggttttccccttggttggggggttttttttaaaattggttggttggaaaaggggccaaggggttggggaattttttggggaaaaggccttaaggggttaattaaaaaaaattttaattggaaggggttggttccttccttaaccaaggggttggaaccttggggaaaattccttggaaccaaggccaattggggttggttggttaaccttaaggggttttggggaattggggttggttccccccaaccccaaaaaaccttttaattggttccccccttccccttggggaaaaccccccaaaaaaccttccttggaaccccttttaattttttggggaaaaaaccaattggggttccaaccttggccaaggaattttttaaaattccaaggttttaaaaggaattggaaaaggttccttttaaccaaggttggggggccccttttttaaaattttccaaaattaattggttccttggggccaattccccccaaaattggaaggaaaaggaaggaaaaggaaggaaccccaaaaggaaaaggccttggaattaaEE' > $odir/snippet2.pa

./pasta -action rotini-ref -i $odir/snippet2.pa > $odir/ref-test2.inp

diff <( ./pasta -action rotini-gvcf -i $odir/snippet2.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test2.inp | ./pasta -action rotini-ref ) \
  <( ./pasta -action rotini-ref -i $odir/snippet2.pa ) || echo "snippet2 ref failed"

diff <( ./pasta -action rotini-gvcf -i $odir/snippet2.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test2.inp | ./pasta -action rotini-alt0 ) \
  <( ./pasta -action rotini-alt0 -i $odir/snippet2.pa ) || echo "snippet2 alt0 failed"

diff <( ./pasta -action rotini-gvcf -i $odir/snippet2.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test2.inp | ./pasta -action rotini-alt1 ) \
  <( ./pasta -action rotini-alt1 -i $odir/snippet2.pa ) || echo "snippet2 alt1 failed"

echo 'ok-snippet2'

# There was a problem with an alt->noc->ref sequence that was giving a weird (and invalid) GT value.
# The GT sequence was 1/0/2 for a bi-allelic stream. Test the following pasta sequence to make sure
# it's good now.

echo '.ZAaaaaaaaaaaaaaaaaaaaaaaa.Saaaaaaaaccaaccttaaggaaaaggaaaaaaaaccaattggggggttggaaggccttttttttccttggttaaaaccaattggggaattggttaaaaggaaaaaaaaggggttttttttccttaaaaccaaaattggaaccttccaaaaaaaattccccaaggaattggttaaggttttaaaaaaggaaaaaaaaggaattttggaaccaaaaggttttttggaattttaaccaaaaaaaaggttccccaattggaattaaaaggccaaccaaggttccaaaaaaaattggccaaaaaattggaaccttttttggttaaggggaaggaaaaggggttaattttttggccaaggccaattaattaattccaaggaaaaggttaaaaaaggggggttttaaaattaattccccccttggaattaatt' > $odir/snippet3.pa

./pasta -action rotini-ref -i $odir/snippet3.pa > $odir/ref-test3.inp

diff <( ./pasta -action rotini-gvcf -i $odir/snippet3.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test3.inp | ./pasta -action rotini-ref ) \
  <( ./pasta -action rotini-ref -i $odir/snippet3.pa ) || echo "snippet3 ref failed"

diff <( ./pasta -action rotini-gvcf -i $odir/snippet3.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test3.inp | ./pasta -action rotini-alt0 ) \
  <( ./pasta -action rotini-alt0 -i $odir/snippet3.pa ) || echo "snippet3 alt0 failed"

diff <( ./pasta -action rotini-gvcf -i $odir/snippet3.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test3.inp | ./pasta -action rotini-alt1 ) \
  <( ./pasta -action rotini-alt1 -i $odir/snippet3.pa ) || echo "snippet3 alt1 failed"

echo 'ok-snippet3'


echo '!!EEaaEEEEtt77aattttttaaaaggttggccttttaattttttttttttttaaaaccccccaaaattttaaaattttaaggaaggccttttttttttaaTTAATTAATTAAAAaaccaattaaccaaaaccaaccaattaattaaaaaattaaccaaccaaggaaccaaggaaccaaggaaaaggaattttccaaggccaaccttttggttaaaaggaattttttttttccaattttttggccccaaggttttttccttttaaaattttggggaattggaaccttggggccttttccaaggggggttggggaaggccccccttttggggaaaaggaaaaccaaaaggggccttggggggaaaaaaggccttttggggttttttccttaaggggggccccaaaaaattaaaaggccaaggccttggaaaaggggccaaaaaaggaaccaaggaaggttccttttaaaaaaaattttaaaaggggaatt' > $odir/snippet4.pa

./pasta -action rotini-ref -i $odir/snippet4.pa > $odir/ref-test4.inp

diff <( ./pasta -action rotini-gvcf -i $odir/snippet4.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test4.inp | ./pasta -action rotini-ref ) \
  <( ./pasta -action rotini-ref -i $odir/snippet4.pa ) || echo "snippet4 ref failed"

diff <( ./pasta -action rotini-gvcf -i $odir/snippet4.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test4.inp | ./pasta -action rotini-alt0 ) \
  <( ./pasta -action rotini-alt0 -i $odir/snippet4.pa ) || echo "snippet4 alt0 failed"

diff <( ./pasta -action rotini-gvcf -i $odir/snippet4.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test4.inp | ./pasta -action rotini-alt1 ) \
  <( ./pasta -action rotini-alt1 -i $odir/snippet4.pa ) || echo "snippet4 alt1 failed"

echo 'ok-snippet4'

## snippet 5

echo '!!77aa7gaaggaaggaaggaaggaaggaaggaaaaaattaaccttccttggggttttttccaaggttttttccttttggttttaaaaaattaaccccccaaaaggggaattttaaaaaattggggaattaaaaaattaaaattttaattttttggaaggaaccaaggaaggggccaaggttttttccccaaggaaggttttttaaaattaaccccaaggaaaaggttttttttttccccaaaaaaggggccccttaaggggttccaaggttggggttggccaaccccaaggggaaccttggggccttggggttttaaaaccaaggttaaaattccccccttaaggaaggaattttttttttaaaaggggaaggttccccccttggccccttccaaccccccttccttaaaaaaaaaattttccttaaaaggccccaaggttggaaggg' > $odir/snippet5.pa

./pasta -action rotini-ref -i $odir/snippet5.pa > $odir/ref-test5.inp

diff <( ./pasta -action rotini-gvcf -i $odir/snippet5.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test5.inp | ./pasta -action rotini-ref ) \
  <( ./pasta -action rotini-ref -i $odir/snippet5.pa ) || echo "snippet5 ref failed"

diff <( ./pasta -action rotini-gvcf -i $odir/snippet5.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test5.inp | ./pasta -action rotini-alt0 ) \
  <( ./pasta -action rotini-alt0 -i $odir/snippet5.pa ) || echo "snippet5 alt0 failed"

diff <( ./pasta -action rotini-gvcf -i $odir/snippet5.pa | ./pasta -action gvcf-rotini -refstream $odir/ref-test5.inp | ./pasta -action rotini-alt1 ) \
  <( ./pasta -action rotini-alt1 -i $odir/snippet5.pa ) || echo "snippet5 alt1 failed"

echo 'ok-snippet5'

## custom small gvcf snippet
## edge case when whole tilepath is nocall (with
## a spurious gvcf line) and then filtered
##
## full-sequence is needed
##
## array-data/ref-snippet5.seq made with:
##
##    refstream /data-sdd/data/ref/hg38.fa.gz chrY:17999991+100
##

gvcf_snippet="chrY 18000000 . C A,<NON_REF> 0.2 . . GT ./."

t5s_a=`./pasta -full-sequence \
  -action gvcf-rotini \
  -i <( echo "$gvcf_snippet" | tr ' ' '\t' ) \
  -r "assay-data/ref-snippet6.seq" \
  -start 17999990 | \
  ./pasta -a filter-rotini -start 17999995 -n 10 | \
  ./pasta -a rotini-ref | tr -d '\n' | \
  md5sum | cut -f1 -d' '`

t5s_b=`cat assay-data/ref-snippet6.seq | \
  tr -d '\n' | \
  head -c 15 | tail -c 10 | tr -d '\n' | \
  md5sum | cut -f1 -d' '`

if [[ "$t5s_a" == "$t5s_b" ]] ; then
  echo "ok-snippet6"
else
  echo "FAIL: snippet6"
fi




## custom small gvcf snippet
##

export inpgvcf="assay-data/gvcf-snippet.gvcf"
export inpref="assay-data/gvcf-snippet-ref.seq"

diff <( ./pasta -action gvcf-rotini -i $inpgvcf -r $inpref | ./pasta -action rotini-ref | tr -d '\n' | md5sum  ) \
  <( cat $inpref | tr -d '\n' | md5sum ) || echo "gvcf-snippet failed"

diff <( ./pasta -action gvcf-rotini -i $inpgvcf -r $inpref | ./pasta -action rotini-alt0 | tr -d '\n' | md5sum ) \
  <( ./pasta -a gvcf-rotini -i $inpgvcf -r $inpref | ./pasta -a rotini-gvcf | ./pasta -a gvcf-rotini -r $inpref | ./pasta -a rotini-alt0 | tr -d '\n' | md5sum )

diff <( ./pasta -action gvcf-rotini -i $inpgvcf -r $inpref | ./pasta -action rotini-alt1 | tr -d '\n' | md5sum ) \
  <( ./pasta -a gvcf-rotini -i $inpgvcf -r $inpref | ./pasta -a rotini-gvcf | ./pasta -a gvcf-rotini -r $inpref | ./pasta -a rotini-alt1 | tr -d '\n' | md5sum )

echo "ok-custom-gvcf"

exit 0

#diff $odir/gvcf-nocall.inp $odir/gvcf-nocall.out
#diff <( cat $odir/gvcf-nocall.inp | tr -d '\n' | fold -w 50 ) <( cat $odir/gvcf-nocall.out | tr -d '\n'  | fold -w 50 )
diff <( cat $odir/gvcf-nocall.inp | tr -d '\n' | sed 's/[ACTG]*$//' | fold -w 50 ) <( cat $odir/gvcf-nocall.out | tr -d '\n' | sed 's/[ACTG]*$//' | fold -w 50 )


## GVCF with het nocall
##
refseed="11223344"
altseed="1234"

param_inp="p-indel-nocall=0.5:p-indel=0.5:ref-seed=$refseed:seed=$altseed:p-nocall=0.3"
param_ref="ref-seed=$refseed:allele=1"

./pasta -action rstream -param "$param_inp" > $odir/gvcf-indel-nocall.inp
./pasta -action rotini-gvcf -i $odir/gvcf-indel-nocall.inp | ./pasta -action gvcf-rotini -refstream <( ./pasta -action ref-rstream -param "$param_ref" ) > $odir/gvcf-indel-nocall.out

diff <( ./pasta -action rotini-ref -i $odir/gvcf-indel-nocall.inp ) <( ./pasta -action rotini-ref -i $odir/gvcf-indel-nocall.out )
diff <( ./pasta -action rotini-alt0 -i $odir/gvcf-indel-nocall.inp ) <( ./pasta -action rotini-alt0 -i $odir/gvcf-indel-nocall.out )
diff <( ./pasta -action rotini-alt1 -i $odir/gvcf-indel-nocall.inp ) <( ./pasta -action rotini-alt1 -i $odir/gvcf-indel-nocall.out )

echo ok
exit 0
