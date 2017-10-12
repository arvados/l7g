#!/bin/bash

tilepath="0003"



#fastj2cgflib -o $tilepath.cglf -t <( ./verbose_tagset $tilepath ) -f <( find ./data -name $tilepath.fj.gz | head -n10 | xargs zcat )
#fastj2cgflib -V -o $tilepath.cglf -t $tilepath.verbose_tagset  -f <( find ./data -name $tilepath.fj.gz | xargs zcat )
fastj2cgflib -V -t $tilepath.verbose_tagset  -f <( find ./data -name $tilepath.fj.gz | xargs zcat )
