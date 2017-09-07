#!/bin/bash

function _q {
  echo $1
  exit 1
}


odir="assay/cgivar"
mkdir -p $odir

ifn="cgivar-snp.inp"
ofn="cgivar-snp.out"

## CGI-Var with snps
##
./pasta -action rstream -param 'p-snp=0.5:ref-seed=11223344:n=1000:seed=1234' > $odir/$ifn
./pasta -action rotini-cgivar -i $odir/$ifn \
  | ./pasta -action cgivar-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' ) > $odir/$ofn

diff <( ./pasta -action rotini-ref -i $odir/$ifn ) <( ./pasta -action rotini-ref -i $odir/$ofn ) || _q "cgivar snp ref"
diff <( ./pasta -action rotini-alt0 -i $odir/$ifn ) <( ./pasta -action rotini-alt0 -i $odir/$ofn ) || _q "cgivar snp alt0"
diff <( ./pasta -action rotini-alt1 -i $odir/$ifn ) <( ./pasta -action rotini-alt1 -i $odir/$ofn ) || _q "cgivar snp alt1"

echo ok-snp

## CGI-VAR with indels
##
./pasta -action rstream -param 'p-indel=0.8:p-indel-length=0,3:p-nocall=0:ref-seed=11223344:n=1000:seed=1234' > $odir/cgivar-indel.inp
./pasta -action rotini-cgivar -i $odir/cgivar-indel.inp | \
  ./pasta -action cgivar-rotini \
     -refstream <( ./pasta -action ref-rstream \
     -param 'ref-seed=11223344:allele=1' ) \
     > $odir/cgivar-indel.out

diff <( ./pasta -action rotini-ref -i $odir/cgivar-indel.inp ) <( ./pasta -action rotini-ref -i $odir/cgivar-indel.out ) || _q "cgivar indel ref"
diff <( ./pasta -action rotini-alt0 -i $odir/cgivar-indel.inp ) <( ./pasta -action rotini-alt0 -i $odir/cgivar-indel.out ) || _q "cgivar indel alt0"
diff <( ./pasta -action rotini-alt1 -i $odir/cgivar-indel.inp ) <( ./pasta -action rotini-alt1 -i $odir/cgivar-indel.out ) || _q "cgivar indel alt1"


echo ok-indel

## CGI-VAR with nocall
##
ofn_b="cgivar-nocall"
./pasta -action rstream -param 'p-nocall=0.3:ref-seed=11223344:seed=1234' > $odir/$ofn_b.inp
./pasta -action rotini-cgivar -i $odir/$ofn_b.inp | \
  ./pasta -action cgivar-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' ) \
    > $odir/$ofn_b.out


diff <( ./pasta -action rotini-ref -i $odir/$ofn_b.inp ) <( ./pasta -action rotini-ref -i $odir/$ofn_b.out ) || _q "cgivar nocall ref"
diff <( ./pasta -action rotini-alt0 -i $odir/$ofn_b.inp ) <( ./pasta -action rotini-alt0 -i $odir/$ofn_b.out ) || _q "cgivar nocall alt0"
diff <( ./pasta -action rotini-alt1 -i $odir/$ofn_b.inp ) <( ./pasta -action rotini-alt1 -i $odir/$ofn_b.out ) || _q "cgivar nocall alt1"

echo ok-nocall


## CGI-VAR with indels and nocalls
##
#./pasta -action rstream -param 'p-nocall=0.3:p-indel=0.3:ref-seed=11223344:seed=1234'  > $odir/cgivar-indel-nocall.inp
#./pasta -action rstream -param 'p-nocall=0.3:p-indel=0.3:p-indel-nocall=0.8:ref-seed=11223344:seed=1234'  > $odir/cgivar-indel-nocall.inp
./pasta -action rstream -param 'p-nocall=0.3:p-indel=0.5:p-indel-nocall=0.8:ref-seed=11223344:seed=1234'  > $odir/cgivar-indel-nocall.inp
./pasta -action rotini-cgivar -i $odir/cgivar-indel-nocall.inp | ./pasta -action cgivar-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' ) > $odir/cgivar-indel-nocall.out

diff <( ./pasta -action rotini-ref -i $odir/cgivar-indel-nocall.inp ) \
  <( ./pasta -action rotini-ref -i $odir/cgivar-indel-nocall.out ) || _q "error indel-nocall ref"
diff <( ./pasta -action rotini-alt0 -i $odir/cgivar-indel-nocall.inp ) \
  <( ./pasta -action rotini-alt0 -i $odir/cgivar-indel-nocall.out ) || _q "error indel-nocall alt0"
diff <( ./pasta -action rotini-alt1 -i $odir/cgivar-indel-nocall.inp ) \
  <( ./pasta -action rotini-alt1 -i $odir/cgivar-indel-nocall.out ) || _q "error indel-nocall alt1"

echo ok-indel-nocall
exit 0
