#!/bin/bash

odir="assay/gff"
mkdir -p $odir

## GFF with snps
##
./pasta -action rstream -param 'p-snp=0.5:ref-seed=11223344:n=1000:seed=1234' > $odir/gff-snp.inp
./pasta -action rotini-gff -i $odir/gff-snp.inp | ./pasta -action gff-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' ) > $odir/gff-snp.out

diff <( ./pasta -action rotini-ref -i $odir/gff-snp.inp ) <( ./pasta -action rotini-ref -i $odir/gff-snp.out )
diff <( ./pasta -action rotini-alt0 -i $odir/gff-snp.inp ) <( ./pasta -action rotini-alt0 -i $odir/gff-snp.out )
diff <( ./pasta -action rotini-alt1 -i $odir/gff-snp.inp ) <( ./pasta -action rotini-alt1 -i $odir/gff-snp.out )

## GFF with indels
##
./pasta -action rstream -param 'p-indel=0.5:p-indel-length=0,3:ref-seed=11223344:n=1000:seed=1234' > $odir/gff-indel.inp
./pasta -action rotini-gff -i $odir/gff-indel.inp | ./pasta -action gff-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' ) > $odir/gff-indel.out

diff <( ./pasta -action rotini-ref -i $odir/gff-indel.inp ) <( ./pasta -action rotini-ref -i $odir/gff-indel.out )
diff <( ./pasta -action rotini-alt0 -i $odir/gff-indel.inp ) <( ./pasta -action rotini-alt0 -i $odir/gff-indel.out )
diff <( ./pasta -action rotini-alt1 -i $odir/gff-indel.inp ) <( ./pasta -action rotini-alt1 -i $odir/gff-indel.out )


## GFF with nocall
##
./pasta -action rstream -param 'p-nocall=0.3:ref-seed=11223344:seed=1234' > $odir/gff-nocall.inp
./pasta -action rotini-gff -i $odir/gff-nocall.inp | ./pasta -action gff-rotini -refstream <( ./pasta -action ref-rstream -param 'ref-seed=11223344:allele=1' )  > $odir/gff-nocall.out


#diff $odir/gff-nocall.inp $odir/gff-nocall.out
#diff <( cat $odir/gff-nocall.inp | tr -d '\n' | fold -w 50 ) <( cat $odir/gff-nocall.out | tr -d '\n'  | fold -w 50 )
diff <( cat $odir/gff-nocall.inp | tr -d '\n' | sed 's/[ACTG]*$//' | fold -w 50 ) <( cat $odir/gff-nocall.out | tr -d '\n' | sed 's/[ACTG]*$//' | fold -w 50 )


## GFF with het nocall
##
refseed="11223344"
altseed="1234"

param_inp="p-indel-nocall=0.5:p-indel=0.5:ref-seed=$refseed:seed=$altseed:p-nocall=0.3"
param_ref="ref-seed=$refseed:allele=1"

./pasta -action rstream -param "$param_inp" > $odir/gff-indel-nocall.inp
./pasta -action rotini-gff -i $odir/gff-indel-nocall.inp | ./pasta -action gff-rotini -refstream <( ./pasta -action ref-rstream -param "$param_ref" ) > $odir/gff-indel-nocall.out

diff <( ./pasta -action rotini-ref -i $odir/gff-indel-nocall.inp ) <( ./pasta -action rotini-ref -i $odir/gff-indel-nocall.out )
diff <( ./pasta -action rotini-alt0 -i $odir/gff-indel-nocall.inp ) <( ./pasta -action rotini-alt0 -i $odir/gff-indel-nocall.out )
diff <( ./pasta -action rotini-alt1 -i $odir/gff-indel-nocall.inp ) <( ./pasta -action rotini-alt1 -i $odir/gff-indel-nocall.out )

## test with a snippet of real data
##
./pasta -action gff-rotini \
  -i assay/data/hu826751.chr5.55091499-55111602.gff \
  -refstream assay/data/hu826751.chr5.55091499-55111602.refstream \
  -start 55091498 > $odir/"hu826751.chr5.55091498-55111601".inp.pa

cat $odir/"hu826751.chr5.55091498-55111601".inp.pa | \
  ./pasta -action rotini-gff | \
  ./pasta -action gff-rotini \
    -i assay/data/hu826751.chr5.55091499-55111602.gff \
    -refstream assay/data/hu826751.chr5.55091499-55111602.refstream \
    -start 55091498 > $odir/"hu826751.chr5.55091498-55111601".out.pa

diff <( ./pasta -action rotini-ref -i $odir/"hu826751.chr5.55091498-55111601".inp.pa ) \
  <( ./pasta -action rotini-ref -i $odir/"hu826751.chr5.55091498-55111601".out.pa )

diff <( ./pasta -action rotini-alt0 -i $odir/"hu826751.chr5.55091498-55111601".inp.pa ) \
  <( ./pasta -action rotini-alt0 -i $odir/"hu826751.chr5.55091498-55111601".out.pa )

diff <( ./pasta -action rotini-alt1 -i $odir/"hu826751.chr5.55091498-55111601".inp.pa ) \
  <( ./pasta -action rotini-alt1 -i $odir/"hu826751.chr5.55091498-55111601".out.pa )


echo ok
exit 0
