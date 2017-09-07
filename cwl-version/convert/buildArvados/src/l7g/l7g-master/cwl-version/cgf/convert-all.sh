#!/bin/bash
#
# wrapper script to parallelize fastj to cgf conversion
#

# Processes are getting killed (OOM).  I think this can balloon up
# to multiple gigs from the intermediate representaiton, presumably.
# Let's try 8 and see if that fixes the majority of the problems.
# We can always run the outliers after the fact.
#
#parallel -P 15 ./process-single.sh ::: `csvtool col 1 data/name-pdh.csv`
parallel -P 8 ./process-single.sh ::: `csvtool col 1 data/name-pdh.csv`
