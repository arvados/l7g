#!/bin/sh

REFERENCE=$1
CGIVAR=$2

cgatools mkvcf --beta --reference $REFERENCE --include-no-calls --field-names GT,GQ,DP,AD --source-names masterVar --master-var $CGIVAR || true
