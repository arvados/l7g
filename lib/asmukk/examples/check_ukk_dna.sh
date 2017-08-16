#!/bin/bash
#

set -e

#VERBOSE=true
VERBOSE=false

DP="./dp_dna"
ASMUKK="../asm_ukk_dna"

ins="0.5"
del="0.5"
sub="0.5"
noc="0.5"

function run_test {
  n=$1
  n_it=$2

  start_seed=100
  let end_seed="$start_seed + $n_it"

  for seed in `seq $start_seed $end_seed` ; do
    score_dp=`$DP < <( ./mkseq -n $n -s $seed -I $ins -D $del -U $sub -N $noc ) | head -n1 | cut -f1 -d ' '`
    score_ukk=`$ASMUKK < <( ./mkseq -n $n -s $seed -I $ins -D $del -U $sub -N $noc  ) | head -n1 | cut -f1 -d' '`

    ma=`$DP < <( ./mkseq -n $n -s $seed -I $ins -D $del -U $sub -N $noc ) | md5sum | cut -f1 -d' '`
    mb=`$ASMUKK < <( ./mkseq -n $n -s $seed -I $ins -D $del -U $sub -N $noc ) | md5sum | cut -f1 -d' '`

    if [ $VERBOSE == "true" ] ; then
      time $DP < <( ./mkseq -n $n -s $seed -I $ins -D $del -U $sub -N $noc ) || true
      time $ASMUKK < <( ./mkseq -n $n -s $seed -I $ins -D $del -U $sub -N $noc ) || true
      echo $score_dp $score_ukk
    fi


    if [ "$score_dp" != "$score_ukk" ] || [ "$ma" != "$mb" ] ; then
      echo -e ERROR "scores or sequences do not match for n $n, seed $seed, ins $ins, del $del, sub $sub, noc $noc"
      ./mkseq -n $n -s $seed -I $ins -D $del -U $sub -N $noc
      exit 1
    fi

  done
}


echo -n "testing n=100, ins=$ins, del=$del, sub=$sub, noc=$noc (100 iterations)...."
run_test 100 100
echo "ok"

echo -n "testing n=200, ins=$ins, del=$del, sub=$sub, noc=$noc (200 iterations)...."
run_test 200 100
echo "ok"

echo -n "testing n=300, ins=$ins, del=$del, sub=$sub, noc=$noc (300 iterations)...."
run_test 300 100
echo "ok"

echo -n "testing n=400, ins=$ins, del=$del, sub=$sub, noc=$noc (400 iterations)...."
run_test 400 100
echo "ok"

echo "ok, tests passed"
