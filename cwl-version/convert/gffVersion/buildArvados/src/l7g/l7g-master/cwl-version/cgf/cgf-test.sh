#!/bin/bash
#
#


#path="01fa"
#path="01f9"
#path="01f8"
#path="01e2"
path="030c"

./bin/cgf -action header -i nop -o test.cgf
./bin/cgf -action append -i <( zcat /data-sde/data/fastj/hu826751-GS03052-DNA_B01/$path.fj.gz ) -path $path -S <( zcat /data-sde/data/sglf/$path.sglf.gz ) -cgf test.cgf -o test.cgf
