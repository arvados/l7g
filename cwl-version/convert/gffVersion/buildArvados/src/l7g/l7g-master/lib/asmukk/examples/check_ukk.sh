#!/bin/bash
#

set -e

#VERBOSE=true
VERBOSE=false

DP="./dp"
ASMUKK="../asm_ukk"

function run_test {
  n=$1
  n_it=$2
  p=$3

  start_seed=100
  let end_seed="$start_seed + $n_it"

  for seed in `seq $start_seed $end_seed` ; do
    score_dp=`$DP < <( ./mktest $n $seed $p ) | head -n1 | cut -f1 -d ' '`
    score_ukk=`$ASMUKK < <( ./mktest $n $seed $p ) | head -n1 | cut -f1 -d' '`

    ma=`$DP < <( ./mktest $n $seed $p ) | md5sum | cut -f1 -d' '`
    mb=`$ASMUKK < <( ./mktest $n $seed $p ) | md5sum | cut -f1 -d' '`

    if [ $VERBOSE == "true" ] ; then
      time $DP < <( ./mktest $n $seed $p )
      time $ASMUKK < <( ./mktest $n $seed $p )
      echo $score_dp $score_ukk
    fi


    if [ "$score_dp" != "$score_ukk" ] || [ "$ma" != "$mb" ] ; then
      echo -e ERROR "scores or sequences do not match for n $n, seed $seed, p $p"
      ./mktest $n $seed $p
      exit 1
    fi

  done
}


echo -n "testing n=100, p=0.5 (100 iterations)...."
run_test 100 100 0.5
echo "ok"

echo -n "testing n=200, p=0.5 (100 iterations)...."
run_test 200 100 0.5
echo "ok"

echo -n "testing n=300, p=0.5 (100 iterations)...."
run_test 300 100 0.5
echo "ok"

echo -n "testing n=400, p=0.5 (100 iterations)...."
run_test 400 100 0.5
echo "ok"

echo "ok, tests passed"
