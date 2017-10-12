#!/bin/bash
#

set -e

VERBOSE=false
#VERBOSE=true

#CHECK_SEQ=1

DP="./dpi"
AIMUKK="../aim_ukk"

ins="0.5"
del="0.5"
sub="0.5"
range="100"

function run_test {
  n=$1
  n_it=$2

  start_seed=100
  let end_seed="$start_seed + $n_it"

  for seed in `seq $start_seed $end_seed` ; do
    score_dp=`$DP < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub ) | head -n1 | cut -f1 -d ' '`
    score_ukk=`$AIMUKK < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub ) | head -n1 | cut -f1 -d' '`

    ma=`$DP < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub ) | md5sum | cut -f1 -d' '`
    mb=`$AIMUKK < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub ) | md5sum | cut -f1 -d' '`

    if [ $VERBOSE == "true" ] ; then
      time $DP < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub ) || true
      time $AIMUKK < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub ) || true
      echo $score_dp $score_ukk
    fi


    if [ "$score_dp" != "$score_ukk" ] || [ "$ma" != "$mb" ] ; then
      echo -e ERROR "scores or sequences do not match for n $n, seed $seed, ins $ins, del $del, sub $sub, range $range"
      ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub
      exit 1
    fi

    ## check sequences reversed
    ##

    score_dp=`$DP < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub | tac ) | head -n1 | cut -f1 -d ' '`
    score_ukk=`$AIMUKK < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub | tac ) | head -n1 | cut -f1 -d' '`

    ma=`$DP < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub | tac ) | md5sum | cut -f1 -d' '`
    mb=`$AIMUKK < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub | tac ) | md5sum | cut -f1 -d' '`

    if [ $VERBOSE == "true" ] ; then
      time $DP < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub | tac ) || true
      time $AIMUKK < <( ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub | tac ) || true
      echo $score_dp $score_ukk
    fi


    if [ "$score_dp" != "$score_ukk" ] || [ "$ma" != "$mb" ] ; then
      echo -e ERROR "scores or sequences do not match for n $n, seed $seed, ins $ins, del $del, sub $sub, range $range"
      ./mkiarray -R $range -n $n -s $seed -I $ins -D $del -U $sub | tac
      exit 1
    fi


  done
}


function run_custom_test {
  fn=$1
  count=0

  while IFS='' read -r line || [[ -n "$line" ]] ; do

    score_dp=`echo "$line" | tr '\t' '\n' | $DP -m 10000 -g 1 | head -n1 | cut -f1 -d' '`
    score_ukk=`echo "$line" | tr '\t' '\n' | $AIMUKK -m 10000 -g 1 | head -n1 | cut -f1 -d' '`

    if [[ $CHECK_SEQ ]] ; then
      ma=`echo "$line" | tr '\t' '\n' | $DP -m 10000 -g 1 | md5sum | cut -f1 -d' '`
      mb=`echo "$line" | tr '\t' '\n' | $AIMUKK -m 10000 -g 1| md5sum | cut -f1 -d' '`

      if [ "$ma" != "$mb" ] ; then
        echo -e ERROR "scores or sequences do not match for sequence pair $count in $fn"
        echo "===="
        echo "$line" | tr '\t' '\n'
        echo "===="
        exit 1
      fi

    else

      if [ "$score_dp" != "$score_ukk" ] ; then
        echo -e ERROR "scores do not match for sequence pair $count in $fn"
        echo "===="
        echo "$line" | tr '\t' '\n'
        echo "===="
        exit 1
      fi

    fi

    count=`expr $count + 1`

  done < <( cat "$fn" | paste - - )
}

echo -n "testing n=100, ins=$ins, del=$del, sub=$sub, range=$range (100 iterations)...."
run_test 100 100
echo "ok"

echo -n "testing n=200, ins=$ins, del=$del, sub=$sub, range=$range (200 iterations)...."
run_test 200 100
echo "ok"

echo -n "testing n=300, ins=$ins, del=$del, sub=$sub, range=$range (300 iterations)...."
run_test 300 100
echo "ok"

echo -n "testing n=400, ins=$ins, del=$del, sub=$sub, range=$range (400 iterations)...."
run_test 400 100
echo "ok"

ins="0.05"
del="0.05"
sub="0.05"

echo -n "testing n=100, ins=$ins, del=$del, sub=$sub, range=$range (100 iterations)...."
run_test 100 100
echo "ok"

echo -n "testing n=200, ins=$ins, del=$del, sub=$sub, range=$range (200 iterations)...."
run_test 200 100
echo "ok"

echo -n "testing n=300, ins=$ins, del=$del, sub=$sub, range=$range (300 iterations)...."
run_test 300 100
echo "ok"

echo -n "testing n=400, ins=$ins, del=$del, sub=$sub, range=$range (400 iterations)...."
run_test 400 100
echo "ok"

ins="0.5"
del="0.5"
sub="0.5"
range="10"

echo -n "testing n=100, ins=$ins, del=$del, sub=$sub, range=$range (100 iterations)...."
run_test 100 100
echo "ok"

echo -n "testing n=200, ins=$ins, del=$del, sub=$sub, range=$range (200 iterations)...."
run_test 200 100
echo "ok"

echo -n "testing n=300, ins=$ins, del=$del, sub=$sub, range=$range (300 iterations)...."
run_test 300 100
echo "ok"

echo -n "testing n=400, ins=$ins, del=$del, sub=$sub, range=$range (400 iterations)...."
run_test 400 100
echo "ok"

ins="0.05"
del="0.05"
sub="0.05"
range="10"

echo -n "testing n=100, ins=$ins, del=$del, sub=$sub, range=$range (100 iterations)...."
run_test 100 100
echo "ok"

echo -n "testing n=200, ins=$ins, del=$del, sub=$sub, range=$range (200 iterations)...."
run_test 200 100
echo "ok"

echo -n "testing n=300, ins=$ins, del=$del, sub=$sub, range=$range (300 iterations)...."
run_test 300 100
echo "ok"

echo -n "testing n=400, ins=$ins, del=$del, sub=$sub, range=$range (400 iterations)...."
run_test 400 100
echo "ok"

echo -n "running custom tests from file 'custom-int.seq'..."
run_custom_test custom-int.seq
echo "ok"

echo "ok, tests passed"
