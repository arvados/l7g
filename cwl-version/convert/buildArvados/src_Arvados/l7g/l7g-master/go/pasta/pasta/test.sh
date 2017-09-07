#!/bin/bash

expect0="ref	0	1	a
alt	1	1	cc/-;-
ref	1	3	gc
alt	3	3	-/aa;-
ref	3	6	tcc
alt	6	7	a/t;a
ref	7	8	c
alt	8	8	t/cc;-
ref	8	11	acc
alt	11	13	-/ccaa;aa"

#z=`./pasta -action interleave -i <( echo -n 'aSSgctccacdacc!!' ) -i <( echo -n 'agcQQtcc.@cSSaccSSaa' )  | ./pasta -action rotini -i - -F`
z=`./pasta -action interleave -i <( echo -n 'aSSgctccacdacc!!' ) -i <( echo -n 'agcQQtcc.@cSSaccSSaa' )  | ./pasta -action rotini-diff -i - -F`

if [ "$expect0" != "$z" ]
then
  echo ERROR: got
  echo "$z"
  echo expected:
  echo "$expect0"
  exit 1
fi


expect1="ref	0	1	.
alt	1	1	cc/-;-
ref	1	3	.
alt	3	3	-/aa;-
ref	3	6	.
alt	6	7	a/t;a
ref	7	8	.
alt	8	8	t/cc;-
ref	8	11	.
alt	11	13	-/ccaa;aa"

#z=`./pasta -action interleave -i <( echo -n 'aSSgctccacdacc!!' ) -i <( echo -n 'agcQQtcc.@cSSaccSSaa' )  | ./pasta -action rotini -i -`
z=`./pasta -action interleave -i <( echo -n 'aSSgctccacdacc!!' ) -i <( echo -n 'agcQQtcc.@cSSaccSSaa' )  | ./pasta -action rotini-diff -i -`

if [ "$expect1" != "$z" ]
then
  echo ERROR: got
  echo "$z"
  echo expected:
  echo "$expect1"
  exit 1
fi

expect2="aaS.S.ggcc.Q.Qttcccc..a@ccdS.Saacccc.S.S!a!a"
#z=`./pasta -action interleave -i <( echo -n 'aSSgctccacdacc!!' ) -i <( echo -n 'agcQQtcc.@cSSaccSSaa' )`
z=`./pasta -action interleave -i <( echo -n 'aSSgctccacdacc!!' ) -i <( echo -n 'agcQQtcc.@cSSaccSSaa' )`

if [ "$expect2" != "$z" ]
then
  echo ERROR: got
  echo "$z"
  echo expected:
  echo "$expect2"
  exit 1
fi

echo Tests passed
