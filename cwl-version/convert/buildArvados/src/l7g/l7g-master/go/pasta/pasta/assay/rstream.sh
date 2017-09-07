#!/bin/bash

mkdir -p assay

function _q {
  echo $1
  exit 1
}

bdir="assay/data"
mkdir -p $bdir

# diploid stream
#
./pasta -action rstream -F -param 'allele=2:n=10000:seed=1234:p-snp=0.3:p-snp-locked=0.5:seed=1234' > $bdir/a.0
./pasta -action rotini-diff -i $bdir/a.0 -F | ./pasta -action diff-rotini -i - > $bdir/b.0
diff $bdir/a.0 $bdir/b.0 || _q "mismatch 0 ($bdir/a.0, $bdir/b.0)"

# two diploid streams concatenated
#
cat <( ./pasta -action rstream -param 'pos=0:n=100:seed=1234' ) <( ./pasta -action rstream -param 'pos=200:chrom=chr1:n=100:seed=4321' ) | sed '/^$/d' > $bdir/a.1
./pasta -action rotini-diff -i $bdir/a.1 -F | ./pasta -action diff-rotini -i - | sed '/^$/d' > $bdir/b.1
diff $bdir/a.1 $bdir/b.1 || _q "mismatch 1 ($bdir/a.1, $bdir/b.1)"

# nocall (locked)
#
./pasta -action rstream -F -param 'allele=2:n=10000:seed=1234:p-nocall=0.3:p-nocall-locked=1.0:seed=1234' > $bdir/a.2
./pasta -action rotini-diff -i $bdir/a.2 -F | ./pasta -action diff-rotini -i - > $bdir/b.2
diff $bdir/a.2 $bdir/b.2 || _q "mismatch 2 ($bdir/a.2, $bdir/b.2)"


# test indels
#
ofn_b="$bdir/indel"
./pasta -action rstream -param 'p-indel=0.5:p-indel-locked=0.8:p-indel-length=0,3:seed=1234' > $ofn_b.inp
./pasta -action rotini-ref -i $ofn_b.inp > $ofn_b.inp.ref
./pasta -action rotini-alt0 -i $ofn_b.inp > $ofn_b.inp.alt0
./pasta -action rotini-alt1 -i $ofn_b.inp > $ofn_b.inp.alt1


./pasta -action rotini-diff -i $ofn_b.inp -F | ./pasta -action diff-rotini > $ofn_b.out
./pasta -action rotini-ref -i $ofn_b.out > $ofn_b.out.ref
./pasta -action rotini-alt0 -i $ofn_b.out > $ofn_b.out.alt0
./pasta -action rotini-alt1 -i $ofn_b.out > $ofn_b.out.alt1

#diff <( ./pasta -action rotini-diff -i $ofn_b.inp -F ) <( ./pasta -action rotini-diff -i $ofn_b.out -F ) || ( echo "indel diff mismatch" && exit 1 )
diff $ofn_b.inp.ref $ofn_b.out.ref || ( echo "indel ref mismatch" && exit 1 )
diff $ofn_b.inp.alt0 $ofn_b.out.alt0 || ( echo "indel alt0 mismatch" && exit 1 )
diff $ofn_b.inp.alt1 $ofn_b.out.alt1 || ( echo "indel alt1 mismatch" && exit 1 )


## snp and nocall
#
ofn_b="$bdir/snp_nocall"
./pasta -action rstream -param 'p-snp=0.8:p-snp-nocall=0.5:seed=1234:p-snp-locked=0.0' > $ofn_b.inp

./pasta -action rotini-ref -i $ofn_b.inp > $ofn_b.inp.ref
./pasta -action rotini-alt0 -i $ofn_b.inp > $ofn_b.inp.alt0
./pasta -action rotini-alt1 -i $ofn_b.inp > $ofn_b.inp.alt1

./pasta -action rotini-diff -i $ofn_b.inp -F --full-nocall-sequence | ./pasta -action diff-rotini > $ofn_b.out
./pasta -action rotini-ref -i $ofn_b.out > $ofn_b.out.ref
./pasta -action rotini-alt0 -i $ofn_b.out > $ofn_b.out.alt0
./pasta -action rotini-alt1 -i $ofn_b.out > $ofn_b.out.alt1

#diff <( ./pasta -action rotini-diff -i $ofn_b.inp -F ) <( ./pasta -action rotini-diff -i $ofn_b.out -F ) || ( echo "indel diff mismatch" && exit 1 )
diff $ofn_b.inp.ref $ofn_b.out.ref || ( echo "indel ref mismatch" && exit 1 )
diff $ofn_b.inp.alt0 $ofn_b.out.alt0 || ( echo "indel alt0 mismatch" && exit 1 )
diff $ofn_b.inp.alt1 $ofn_b.out.alt1 || ( echo "indel alt1 mismatch" && exit 1 )


## indel and nocall
#
ofn_b="$bdir/indel_nocall"
./pasta -action rstream -param 'p-indel=0.5:p-indel-nocall=0.5:seed=1234:n=5000' > $ofn_b.inp

./pasta -action rotini-ref -i $ofn_b.inp > $ofn_b.inp.ref
./pasta -action rotini-alt0 -i $ofn_b.inp > $ofn_b.inp.alt0
./pasta -action rotini-alt1 -i $ofn_b.inp > $ofn_b.inp.alt1

./pasta -action rotini-diff -i $ofn_b.inp -F --full-nocall-sequence | ./pasta -action diff-rotini > $ofn_b.out
./pasta -action rotini-ref -i $ofn_b.out > $ofn_b.out.ref
./pasta -action rotini-alt0 -i $ofn_b.out > $ofn_b.out.alt0
./pasta -action rotini-alt1 -i $ofn_b.out > $ofn_b.out.alt1

#diff <( ./pasta -action rotini-diff -i $ofn_b.inp -F --full-nocall-sequence ) <( ./pasta -action rotini-diff -i $ofn_b.out -F --full-nocall-sequence) ||_q "indel diff mismatch"
diff $ofn_b.inp.ref $ofn_b.out.ref || _q "indel ref mismatch"
diff $ofn_b.inp.alt0 $ofn_b.out.alt0 || _q  "indel alt0 mismatch"
diff $ofn_b.inp.alt1 $ofn_b.out.alt1 || _q "indel alt1 mismatch"


## test random ref streams match up
#
ofn_b="$bdir/seed"
parama="p-snp=0.8:p-indel=0.3:ref-seed=11223344:seed=1234"
paramb="ref-seed=11223344:allele=1"

./pasta -action rstream -param "$parama" > $ofn_b.inp
./pasta -action ref-rstream -param "$paramb" > $ofn_b.ref
diff <( ./pasta -action rotini-ref -i $ofn_b.inp ) $ofn_b.ref || _q "ref streams don't match"


## nocall-indel testing
##
ofn_b="$bdir/indel_nocall2"
refseed="11223344"
seed="1234"
param="p-indel-nocall=0.3:p-indel=0.3:ref-seed=$refseed:seed=$seed"
./pasta -action rstream -param "$param" > $ofn_b.inp
./pasta -action rotini-diff -i $ofn_b.inp -F | ./pasta -action diff-rotini -refstream <( ./pasta -action ref-rstream -param "ref-seed=$refseed:allele=1" ) > $ofn_b.out

diff <( ./pasta -action rotini-ref -i $ofn_b.inp ) <( ./pasta -action rotini-ref -i $ofn_b.out ) || _q "nocall-indel2 ref mismatch"
diff <( ./pasta -action rotini-alt0 -i $ofn_b.inp ) <( ./pasta -action rotini-alt0 -i $ofn_b.out ) || _q "nocall-indel2 alt0 mismatch"
diff <( ./pasta -action rotini-alt1 -i $ofn_b.inp ) <( ./pasta -action rotini-alt1 -i $ofn_b.out ) || _q "nocall-indel2 alt1 mismatch"


## Everything passed
#
echo ok
exit 0
