#!/bin/bash

VERBOSE=1

merge_exec="../merge-sglf"
idir="./data"

tf0="$idir/035e-small.sglf"
tf1="$idir/035e-small-decimated.sglf"
tf2="$idir/035e-small-decimated2.sglf"
tf3="$idir/035e-small-half.sglf"

tf4="$idir/035e-small-top5-bot5.sglf"
tf5="$idir/035e-small-mid20.sglf"

tf6="$idir/035e-small-top8-bot8.sglf"

##

a=`md5sum $tf0 | cut -f1 -d' '`
b=`$merge_exec $tf0 $tf1 | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: $tf0 merged with $tf1 does not give back $tf0 ($a != $b)"
  exit 1
fi

if [[ "$VERBOSE" == 1 ]] ; then
  echo "merge with decimated: ok"
fi

##

a=`cat $tf0 | cut -f2 -d, | sort -u | md5sum | cut -f1 -d' '`
b=`$merge_exec $tf1 $tf2 | cut -f2 -d, | sort -u | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: $tf1 merged with $tf2 does not give back sequence hashes of $tf0 ($a != $b)"
  diff <( cat $tf0 | cut -f2 -d, | sort -u ) <( ../merge-sglf $tf1 $tf2 | cut -f2 -d, | sort -u )
  exit 1
fi

if [[ "$VERBOSE" == 1 ]] ; then
  echo "decimated merge hash match: ok"
fi

##

a=`cat $tf0 | cut -f2,3 -d, | sort -u | md5sum | cut -f1 -d' '`
b=`$merge_exec $tf0 $tf3 | cut -f2,3 -d, | sort -u | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: $tf0 merged with $tf3 does not give back sequence hashes of $tf0 ($a != $b)"
  exit 1
fi

a=`cat $tf0 | cut -f2,3 -d, | sort -u | md5sum | cut -f1 -d' '`
b=`$merge_exec $tf3 $tf0 | cut -f2,3 -d, | sort -u | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: $tf3 merged with $tf0 does not give back sequence hashes of $tf3 ($a != $b)"
  exit 1
fi

if [[ "$VERBOSE" == 1 ]] ; then
  echo "half merge hash-sequence match: ok"
fi

##

a=`cat $tf0 |  md5sum | cut -f1 -d' '`
b=`$merge_exec $tf0 <( echo -e "\n" ) | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: merge with empty sglf failed"
  exit 1
fi

a=`cat $tf0 |  cut -f2 -d, | sort | md5sum | cut -f1 -d' '`
b=`$merge_exec <( echo -e "\n" ) $tf0  | cut -f2 -d, | sort | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: merge with empty sglf failed"
  exit 1
fi

if [[ "$VERBOSE" == 1 ]] ; then
  echo "empty merge: ok"
fi


##

a=`cat $tf0 |  cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`
b=`$merge_exec <( cat $tf0 | head -n 29 )  <( cat $tf0 | tr '\n' 'Z' | sed 's/Z$//' | tr 'Z' '\n' ) | cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: no newline at EOF failed to merge"
  exit 1
fi

a=`cat $tf0 |  cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`
b=`$merge_exec <( cat $tf0 | tr '\n' 'Z' | sed 's/Z$//' | tr 'Z' '\n' )  <( cat $tf0 | head -n29 ) | cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: no newline at EOF failed to merge"
  exit 1
fi

if [[ "$VERBOSE" == 1 ]] ; then
  echo "no newline at eof merge: ok"
fi

##

a=`cat $tf0 |  cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`
b=`$merge_exec $tf4 $tf5 | cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: top-mid merge failed: $tf4, $tf5"
  exit 1
fi

a=`cat $tf0 |  cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`
b=`$merge_exec $tf5 $tf4 | cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: top-mid merge failed: $tf5, $tf4"
  exit 1
fi

a=`cat $tf0 |  cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`
b=`$merge_exec $tf6 $tf5 | cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: top-mid merge failed: $tf6, $tf5"
  exit 1
fi

a=`cat $tf0 |  cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`
b=`$merge_exec $tf5 $tf6 | cut -f2,3 -d, | sort | md5sum | cut -f1 -d' '`

if [[ "$a" != "$b" ]] ; then
  echo "ERROR: top-mid merge failed: $tf5, $tf6"
  exit 1
fi

if [[ "$VERBOSE" == 1 ]] ; then
  echo "top-bottom-mid merge: ok"
fi

##


if [[ "$VERBOSE" == 1 ]] ; then
  echo "ok"
fi
exit 0
