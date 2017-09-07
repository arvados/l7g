#!/bin/bash
#

function _q {
  echo $1
  exit 1
}

odir="assay/fastj"
fjfilter="../../fjtools/fjfilter"

if [[ ! -e "$fjfilter" ]] ; then
  echo ""
  echo "please compile 'fjfilter' in ../../fjtools:"
  echo ""
  echo "  pushd ../../fjtools ; go build fjfilter.go ; popd"
  echo ""
  exit 1
fi

mkdir -p $odir

if [[ "$L7G_TILEASSEMBLY" == "" ]] || [[ "$L7G_TAGSET" == "" ]] ; then
  echo ""
  echo "The following global variables must be defined:"
  echo "  L7G_TILEASSEMBLY  - path to the Lightning tile assembly"
  echo "  L7G_TAGSET        - path to the Lightning tagset (FASTA file)"
  exit 1
fi

if [[ "$hg19" == "" ]] ; then
  echo ""
  echo "Please provide global variable 'hg19' that points to the hg19 reference (FASTA)"
  echo ""
  exit 1
fi

afn="$L7G_TILEASSEMBLY"
aidx="$L7G_TILEASSEMBLY.fwi"
tagdir="$L7G_TAGSET"

ref="hg19"
chrom="chr5"
path="00fa"


inpgff="assay-data/hu826751-chr5-62900001-63200001.gff.gz"

reffa="$hg19"

ucpath=`echo $path | tr '[:lower:]' '[:upper:]'`
prevpath=`echo -e "ibase=16\n$ucpath - 1" | bc -q  | tr '[:upper:]' '[:lower:]'`
prevpath=`printf "%04x" $prevpath`

st0=`l7g assembly $afn $prevpath | tail -n1 | cut -f2`
en0=`l7g assembly $afn $path | tail -n1 | cut -f2`
dn=`expr $en0 - $st0`

st1=`expr $st0 + 1`
en1=`expr $en0 + 1`

realstart1=`tabix $inpgff $chrom:$st1-$en1 | head -n1 | cut -f4`
realend1=`tabix $inpgff $chrom:$st1-$en1 | tail -n1 | cut -f5`
realdn=`expr $realend1 - $realstart1 + 1`

realstart0=`expr $realstart1 - 1`

tabix $inpgff $chrom:$realstart1-$realend1 | \
  ./pasta -action gff-rotini \
    -refstream <( refstream $reffa $chrom:$realstart1-$realend1 ) \
    -start $realstart0 | \
  ./pasta -action filter-rotini -start $st0 -n $dn | fold -w 50 | \
  egrep -v '^>' > $odir/inp_$path.pa

cat $odir/inp_$path.pa | \
  ./pasta -action rotini-fastj -start $st0 -tilepath $path -chrom $chrom -build $ref \
  -assembly <( l7g assembly $afn $path ) \
    -tag <( samtools faidx $tagdir $path.00 | egrep -v '^>' | tr -d '\n' | fold -w 24 ) > $odir/inp_$path.fj

st0=`echo $st0`

###
### step 00e4 is a non-trivial 'knot': One tile has seedTileLength of 2
### and the other of two tiles of length 1.
###
inpfj="$odir/inp_$path.fj"

x0="62958521"
y0="62959015"

x1=`expr $x0 + 1`
y1=`expr $y0 + 1`

step_a="00e4"
step_b="00e5"


./pasta -action fastj-rotini -start $x0 \
  -i <( $fjfilter -i $inpfj -s $path.$step_b -e $path.$step_b ) \
  -assembly <( l7g assembly $afn $path | egrep -A2 '^'$step_a ) \
  -refstream <( samtools faidx $reffa $chrom:$x1-$y1 | egrep -v '^>' | tr '[:upper:]' '[:lower:]' | cat <( echo ">P{$x0}" ) - ) > $odir/out_$path.$step_a.pa

cat <( $fjfilter -i $inpfj -s $path.0.e5.0 -e $path.0.e4.0 | egrep -v '^>' ) \
  <( $fjfilter -i $inpfj -s $path.0.e6.0 -e $path.0.e5.0 | egrep -v '^>' | tr -d '\n' | fold -w 24 | tail -n +2 ) | \
  tr -d '\n' | fold -w 50 > $odir/$path.00e4.seq

diff <( refstream $chrom:$x1-$y0 | egrep -v '^>' | tr -d '\n' | tr '[:upper:]' '[:lower:]' | fold -w 50 ) \
  <( ./pasta -action rotini-ref -i $odir/out_$path.$step_a.pa | tr -d '\n' | fold -w 50 ) || _q "error: ref path $path step $step_a difference"

diff $odir/$path.00e4.seq \
  <( ./pasta -action rotini-alt0 -i $odir/out_$path.$step_a.pa | tr -d '\n' | fold -w 50 ) || _q "error: path $path step $step_a alt0 difference"

diff <( $fjfilter -i $odir/inp_$path.fj -s $path.0.$step_b.1 -e $path.0.$step_b.1 | egrep -v '^>' | tr -d '\n' | fold -w 50 ) \
  <( ./pasta -action rotini-alt1 -i $odir/out_$path.$step_a.pa | tr -d '\n' | fold -w 50 ) || _q "error: path $path step $step_a atl1 difference"


###
### single path test
###

step_check=false

if [ "$step_check" = true ]
then

for x in {1..1181}
do
  f=`expr $x + 1`
  step=`printf "%04x" $x`
  step_p1=`printf "%04x" $f`


  m1=`egrep "$path\.00\.$step\.00[01]" $inpfj | egrep 'seedTileLength" *: *1,' | wc -l`
  m2=`egrep "$path\.00\.$step\.00[01]" $inpfj | egrep 'seedTileLength" *: *2,' | wc -l`

  ## skip over non-trivial tiles that are spanning.  That is, only
  ## consider tiles that are each of seedTileLength 1.
  ##
  if [ "$m1" != 2 ] && [ "$m2" != 2 ]
  then
    echo "#SKIPPING $step"
    continue
  fi

  echo "#step" $step

  #egrep "$path\.00\.$step\.00[01]" $inpfj | sed 's/^>//' | jq -r '.locus[0].build' | head -n1

  ref_st0=`egrep "$path\.00\.$step\.00[01]" $inpfj | sed 's/^>//' | jq -r '.locus[0].build' | head -n1 | cut -f3 -d' '`
  ref_en0=`egrep "$path\.00\.$step\.00[01]" $inpfj | sed 's/^>//' | jq -r '.locus[0].build' | head -n1 | cut -f4 -d' '`

  ref_st1=`expr $ref_st0 + 1`
  ref_en1=`expr $ref_en0 + 1`

  ./pasta -action fastj-rotini -start $ref_st0 \
    -i <( $fjfilter -i $inpfj -s $path.$step_p1 -e $path.$step ) \
    -assembly <( l7g assembly $afn $path | egrep -A1 '^'$step ) \
    -refstream <( samtools faidx $reffa $chrom:$ref_st1-$ref_en0 | egrep -v '^>' | tr '[:upper:]' '[:lower:]' | cat <( echo ">P{$ref_st0}" ) - ) > $odir/out_$path.$step.pa

  diff <( refstream $chrom:$ref_st1-$ref_en0 | egrep -v '^>' | tr -d '\n' | tr '[:upper:]' '[:lower:]' | fold -w 50 ) \
    <( ./pasta -action rotini-ref -i $odir/out_$path.$step.pa | tr -d '\n' | fold -w 50 ) || _q "error: ref path $path step $step difference"

  diff <( $fjfilter -i $odir/inp_$path.fj -s $path.0.$step_p1.0 -e $path.0.$step.0 | egrep -v '^>' | tr -d '\n' | fold -w 50 ) \
    <( ./pasta -action rotini-alt0 -i $odir/out_$path.$step.pa | tr -d '\n' | fold -w 50 ) || _q "error: path $path step $step alt0 difference"

  diff <( $fjfilter -i $odir/inp_$path.fj -s $path.0.$step_p1.1 -e $path.0.$step.1 | egrep -v '^>' | tr -d '\n' | fold -w 50 ) \
    <( ./pasta -action rotini-alt1 -i $odir/out_$path.$step.pa | tr -d '\n' | fold -w 50 ) || _q "error: path $path step $step atl1 difference"

done

fi


###
###
###

./pasta -action fastj-rotini -i $odir/inp_$path.fj -assembly <( l7g assembly $afn $path ) \
  -refstream <( samtools faidx $reffa $chrom:$st1-$en1 | egrep -v '^>' | tr '[:upper:]' '[:lower:]' | cat <( echo ">P{$st0}" ) - ) > $odir/out_$path.pa

diff <( ./pasta -action rotini-ref -i $odir/inp_$path.pa ) <( ./pasta -action rotini-ref -i $odir/out_$path.pa ) || _q "error: ref difference"
diff <( ./pasta -action rotini-alt0 -i $odir/inp_$path.pa ) <( ./pasta -action rotini-alt0 -i $odir/out_$path.pa ) || _q "error: alt0 difference"
diff <( ./pasta -action rotini-alt1 -i $odir/inp_$path.pa ) <( ./pasta -action rotini-alt1 -i $odir/out_$path.pa ) || _q "error: alt1 difference"

echo ok
