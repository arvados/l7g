#!/bin/bash

VERBOSE=1

sglf_check="../sglf-check"
idir="./data"

badid0="$idir/035e-bad-tileid-0.sglf"
badid1="$idir/035e-bad-tileid-1.sglf"
badid2="$idir/035e-bad-tileid-2.sglf"
badid3="$idir/035e-bad-tileid-3.sglf"

badseq="$idir/035e-badseq.sglf"
badhash="$idir/035e-badhash.sglf"

ok0="$idir/035e-a0.sglf"
ok1="$idir/035e-a1.sglf"
ok2="$idir/035e-b0.sglf"
ok3="$idir/035e-b1.sglf"

emptyseq="$idir/035e-empty-seq.sglf"

a=`$sglf_check $badid0`
if [[ "$?" == "0" ]] ; then
  echo "failed to catch bad tileid in $badid0, got $a"
  exit 1
fi

##

a=`$sglf_check $badid1`
if [[ "$?" == "0" ]] ; then
  echo "failed to catch bad tileid in $badid1, got $a"
  exit 1
fi

##

a=`$sglf_check $badid2`
if [[ "$?" == "0" ]] ; then
  echo "failed to catch bad tileid in $badid2, got $a"
  exit 1
fi

##

a=`$sglf_check $badid3`
if [[ "$?" == "0" ]] ; then
  echo "failed to catch bad tileid in $badid3, got $a"
  exit 1
fi

##

a=`$sglf_check $badseq`
if [[ "$?" == "0" ]] ; then
  echo "failed to catch bad tileid in $badseq, got $a"
  exit 1
fi

##

a=`$sglf_check $badhash`
if [[ "$?" == "0" ]] ; then
  echo "failed to catch bad tileid in $badhash, got $a"
  exit 1
fi

##

a=`$sglf_check $ok0`
if [[ "$?" != "0" ]] ; then
  echo "ERROR: $ok0 failed test: $a"
  exit 1
fi


##

a=`$sglf_check $ok1`
if [[ "$?" != "0" ]] ; then
  echo "ERROR: $ok1 failed test: $a"
  exit 1
fi


##

a=`$sglf_check $ok2`
if [[ "$?" != "0" ]] ; then
  echo "ERROR: $ok2 failed test: $a"
  exit 1
fi


##

a=`$sglf_check $ok3`
if [[ "$?" != "0" ]] ; then
  echo "ERROR: $ok3 failed test: $a"
  exit 1
fi

##

a=`$sglf_check $emptyseq`
if [[ "$?" != "0" ]] ; then
  echo "ERROR: $emptyseq failed test: $a"
  exit 1
fi


##

echo "ok"
