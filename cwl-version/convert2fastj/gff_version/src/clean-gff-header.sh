#!/bin/bash
#
# Copyright (C) The Arvados Authors. All rights reserved.
#
# SPDX-License-Identifier: AGPL-3.0
#
####
#
# The Harvard PGP GFF files sometimes have spurious
# 'header' lines interspersed that need to be taken out,
# which tabix has trouble with when indexing.
# Often, the GFF files are also not indexed.
#
# This script takes out all headers, creates a single GFF
# file and then indexes it with tabix and bgzip.
#

VERBOSE=1

igff="$1"
ogff="$2"

if [[ "$igff" == "" ]] ; then
  echo ""
  echo "usage:"
  echo ""
  echo "  clean-gff-header.sh <inputgff> [<outputgff>]"
  echo ""
  exit -1
fi

if [[ "$ogff" == "" ]] ; then
  ogff=`basename $igff`
fi

if [[ "$VERBOSE" -eq 1 ]] ; then
  echo "# igff: $igff"
  echo "# ogff: $ogff"
fi

zcat $igff | egrep -v '^#' | bgzip -c > $ogff

exit 0
