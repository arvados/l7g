#!/bin/bash

expect0="Unk	ref	0	1	a
Unk	alt	1	1	cc/-;-
Unk	ref	1	3	gc
Unk	alt	3	3	-/aa;-
Unk	ref	3	6	tcc
Unk	alt	6	7	a/t;a
Unk	ref	7	8	c
Unk	alt	8	8	t/cc;-
Unk	ref	8	11	acc
Unk	alt	11	13	-/ccaa;aa"

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


expect1="Unk	ref	0	1	.
Unk	alt	1	1	cc/-;-
Unk	ref	1	3	.
Unk	alt	3	3	-/aa;-
Unk	ref	3	6	.
Unk	alt	6	7	a/t;a
Unk	ref	7	8	.
Unk	alt	8	8	t/cc;-
Unk	ref	8	11	.
Unk	alt	11	13	-/ccaa;aa"

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
