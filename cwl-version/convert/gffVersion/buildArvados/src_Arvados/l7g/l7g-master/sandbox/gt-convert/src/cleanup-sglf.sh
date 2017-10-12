#!/bin/bash
#
# Clean up the sglf files in "odata":
#   replace nocalls with 'a'
#   make everything lower case
#   remove duplicates
#   gzip result.
#

d="odata"

for x in `ls $d`; do
  sed 's/n/a/g' $d/$x | tr '[:upper:]' '[:lower:]'  | sort -u > $d/$x.s
  mv $d/$x.s $d/$x
  gzip -f $d/$x
done
