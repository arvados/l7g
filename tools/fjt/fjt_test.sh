#!/bin/bash
#
# some simple tests.
#
# tests include:
#  * CSV output
#  * Band output
#  * Sequence concatenation
#

set -eo pipefail

verbose=1
testdir="testdata"
testtilepath="035e"

function ok_or_exit {
  a="$1"
  b="$2"
  name="$3"

  if [[ "$a" != "$b" ]] ; then
    if [[ "$verbose" == "1" ]] ; then
      echo "$name: $a != $b"
    fi
    exit 1
  fi


  if [[ "$verbose" == "1" ]] ; then
    echo "$name: ok"
  fi

}


# test 0
# Testing CSV output
#
a=`./fjt -C $testdir/$testtilepath.fj | sort | md5sum | cut -f1 -d' '`
b=`sort $testdir/$testtilepath.csv | md5sum | cut -f1 -d' '`

ok_or_exit "$a" "$b" "CSV"

# test 1
# testing band output
#
a=`./fjt -B -L $testdir/$testtilepath.sglf testdata/$testtilepath.fj | sed 's/\[ */[/g' | sed 's/ *\]/]/g' | md5sum | cut -f1 -d' '`
b=`cat $testdir/$testtilepath.band | sed 's/\[ */[/g' | sed 's/ *\]/]/g' | md5sum | cut -f1 -d' '`

ok_or_exit "$a" "$b" "Band"

# test 2
#
a=`./fjt -c 0 $testdir/$testtilepath.fj | tr -d '\n' | md5sum | cut -f1 -d' '`
b=`cat $testdir/$testtilepath.0.seq | tr -d '\n' | md5sum | cut -f1 -d ' ' `

ok_or_exit "$a" "$b" "Seq"

# test 3
#
a=`./fjt -b -L $testdir/$testtilepath.sglf $testdir/$testtilepath.band | ./fjt -c 0 | tr -d '\n' | md5sum | cut -f1 -d' '`
b=`cat $testdir/$testtilepath.0.seq | tr -d '\n' | md5sum | cut -f1 -d' '`

ok_or_exit "$a" "$b" "FastJ Coversion"

exit 0
