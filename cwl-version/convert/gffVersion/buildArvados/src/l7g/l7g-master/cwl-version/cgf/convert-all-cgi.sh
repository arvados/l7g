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
#parallel -P 8 ./convert-fastj-to-cgf-cgi-69.sh ::: `find /data-sde/scripts/convert/stage -maxdepth 1`
#parallel -P 8 ./convert-fastj-to-cgf-cgi-69.sh ::: `find /data-sde/scripts/convert/stage/ -maxdepth 1 -type d | egrep -v '\/$' | head -n 2`
parallel -P 5 ./convert-fastj-to-cgf-cgi-69.sh ::: `find /data-sde/scripts/convert/stage/ -maxdepth 1 -type d | egrep -v '\/$' `
