#!/bin/bash
#
# wrapper script to parallelize fastj to cgf conversion
#

parallel -P 15 ./process-single.sh ::: `csvtool col 1 data/name-pdh.csv | head -n15`
