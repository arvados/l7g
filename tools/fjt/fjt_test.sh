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

# test hashing
#
x0hash=`./fjt -b -L $testdir/035e.1.sglf $testdir/hu826751-035e.band | ./fjt -c 0 | tr -d '\n' | md5sum | cut -f1 -d' '`
x1hash=`./fjt -b -L $testdir/035e.1.sglf $testdir/hu826751-035e.band | ./fjt -c 1 | tr -d '\n' | md5sum | cut -f1 -d' '`

y0hash=`./fjt -b -L $testdir/035e.1.sglf $testdir/hu34D5B9-035e.band | ./fjt -c 0 | tr -d '\n' | md5sum | cut -f1 -d' '`
y1hash=`./fjt -b -L $testdir/035e.1.sglf $testdir/hu34D5B9-035e.band | ./fjt -c 1 | tr -d '\n' | md5sum | cut -f1 -d' '`

a=`echo -e "$x0hash $x1hash"'\n'"$y0hash $y1hash" | md5sum | cut -f1 -d' '`
b=`./fjt -H -L $testdir/035e.1.sglf <( cat $testdir/hu826751-035e.band $testdir/hu34D5B9-035e.band ) | md5sum | cut -f1 -d' '`

ok_or_exit "$a" "$b" "Band Hash"

# batch band conversion
#

a=`./fjt -B -L $testdir/$testtilepath.sglf -I <( echo testdata/$testtilepath.fj.gz ) | sed 's/\[ */[/g' | sed 's/ *\]/]/g' | md5sum | cut -f1 -d' '`
b=`cat $testdir/$testtilepath.band | sed 's/\[ */[/g' | sed 's/ *\]/]/g' | md5sum | cut -f1 -d' '`

ok_or_exit "$a" "$b" "Batch Band"


# Testing the tool that checks fastj files for integrity.
# Start with three control cases (i.e. act on healthy fastj)
# that should trigger no error.
#
for i in -T -t -T\ -t ; do
  if cat $testdir/$testtilepath.fj | ./fjt $i &> /dev/null ; then
    :
  else
    # for the controls, this will be the condition that indicates a problem
    #
    echo control with flag $i: fail
    exit 1
  fi
done
echo controls: ok

if cat $testdir/$testtilepath.broken_n.fj | ./fjt -T &> /dev/null ; then
  echo n: fail
  exit 1
fi
echo n: ok

if cat $testdir/$testtilepath.broken_nocallCount.fj | ./fjt -T &> /dev/null ; then
  echo nocallCount: fail
  exit 1
fi
echo noCallCount: ok

if cat $testdir/$testtilepath.broken_startTag.fj | ./fjt -T &> /dev/null ; then
  echo startTag: fail
  exit 1
fi
echo startTag: ok

if cat $testdir/$testtilepath.broken_endTag.fj | ./fjt -T &> /dev/null ; then
  echo endTag: fail
  exit 1
fi
echo endTag: ok

if cat $testdir/$testtilepath.broken_md5sum.fj | ./fjt -T &> /dev/null ; then
  echo md5sum of seq: fail
  exit 1
fi
echo md5sum of seq: ok

# Tests of tileID consistency
#

if cat $testdir/$testtilepath.broken_seedTileLength.fj | ./fjt -T -t &> /dev/null ; then
  echo seedTileLength: fail
  exit 1
fi
echo seedTileLengths: ok

if cat $testdir/$testtilepath.broken_startTile.fj | ./fjt -T -t &> /dev/null ; then
  echo startTile: fail
  exit 1
fi
echo startTile: ok

if cat $testdir/$testtilepath.broken_endTile.fj | ./fjt -T -t &> /dev/null ; then
  echo endTile: fail
  exit 1
fi
echo endTile: ok

if cat $testdir/$testtilepath.duplicate_tile.fj | ./fjt -T -t &> /dev/null ; then
  echo duplicate tile: fail
  exit 1
fi
echo duplicate tile: ok

if cat $testdir/$testtilepath.non_bool_startTile.fj | ./fjt -T -t &> /dev/null ; then
  echo non-boolean startTile: fail
  exit 1
fi
echo non-boolean startTile: ok

exit 0
