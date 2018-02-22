#!/bin/bash

VERBOSE=1

fjcsv2sglf="../fjcsv2sglf"

#$fjcsv2sglf ../test-data/035e.tagset ../test-data/test0/035e.fjcsv


pos8count=`$fjcsv2sglf ../test-data/035e.tagset ../test-data/test1/035e.fjcsv | cut -f1 -d, | grep -P '^035e\.00\.0008\.' | wc -l | cut -f1 -d' '`
pos9count=`$fjcsv2sglf ../test-data/035e.tagset ../test-data/test1/035e.fjcsv | cut -f1 -d, | grep -P '^035e\.00\.0009\.' | wc -l | cut -f1 -d' '`

if [[ "$pos8count" == "0" ]] || [[ "$pos9count" == "0" ]] ; then
  echo "FAIL: could not find tile position 035e.00.0008 or 035e.00.0009 (md5 hash collision test)"
  exit 1
fi

if [[ "$VERBOSE" == "1" ]] ; then
  echo "ok md5-collision-test"
fi

####
####

if [[ "$VERBOSE" == "1" ]] ; then
  echo ok
fi

exit 0
